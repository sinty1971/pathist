package models

import (
	"errors"
	"path/filepath"
	grpcv1 "server-grpc/gen/grpc/v1"
	"server-grpc/internal/core"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"
)

type Koji struct {
	// Koji メッセージ本体
	*grpcv1.Koji

	// Common 共通モデルフィールド
	Pathist *core.Pathist
}

// NewKoji FolderNameからKojiを作成します（高速化版）
func NewKoji() *Koji {

	koji := &Koji{}
	koji.Koji = grpcv1.Koji_builder{}.Build()
	koji.Pathist = core.NewPathist(koji, core.ConfigMap["KojiPersistFilename"])

	return koji
}

func (m *Koji) GetProtoMessage() proto.Message {
	if m == nil {
		return nil
	}
	return m.Koji
}

// ParseFrom は target から工事開始日・会社名・現場名を取得
func (m *Koji) ParseFrom(pathistFolder string) error {

	var (
		start        = new(Timestamp)
		companyName  string
		locationName string
	)

	// フォルダー名を取得
	foldername := core.GetBaseName(pathistFolder)

	// ファイル名から工事開始日の取得と日付除外文字列の取得
	dateRemoved, err := ParseTimestamp(foldername, start)
	if err != nil {
		return err
	}

	if dateRemoved == "" {
		return errors.New("フォルダー名が工事フォルダーの書式に合致していません")
	}

	// 最初のスペースで分割（最適化）
	if idx := strings.Index(dateRemoved, " "); idx > 0 {
		companyName = dateRemoved[:idx]
		if idx+1 < len(dateRemoved) {
			locationName = dateRemoved[idx+1:]
		}
	} else {
		companyName = dateRemoved
	}

	m.SetPathistFolder(pathistFolder)
	m.SetStart(start.Timestamp)
	m.SetCompanyName(companyName)
	m.SetLocationName(locationName)

	// IDの設定
	id, err := m.Pathist.GenerateId()
	if err != nil {
		return err
	}
	m.SetId(id)
	return nil
}

// GenerateKojiStatus はプロジェクトステータスを判定する
func GenerateKojiStatus(start *Timestamp, end *Timestamp) string {
	if start == nil || end == nil {
		return "不明"
	}

	now := time.Now()

	if !start.IsValid() || !end.IsValid() {
		return "不明"
	} else if now.Before(start.AsTime()) {
		return "予定"
	} else if now.After(end.AsTime()) {
		return "完了"
	} else {
		return "進行中"
	}
}

func (m *Koji) ImportFrom(src *Koji) (*Koji, error) {
	if m == nil || src == nil {
		return nil, errors.New("koji or updatedKoji is nil")
	}

	// 管理フォルダーは変更しない
	src.SetPathistFolder(m.GetPathistFolder())

	// 永続化サービスの設定を引き継ぐ
	// updatedKoji.PersistFilename = obj.PersistFilename

	return src, nil
}

// UpdateFolderPath は工事フォルダー名を更新します
//
// # Id 及び Target の情報は無視されます
//
// TODO: 不完全です、実際にはまだ更新処理していません
func (m *Koji) UpdateFolderPath(src *Koji) bool {
	if src == nil {
		return false
	}

	// 開始日が無効な場合の早期リターン
	start := &Timestamp{Timestamp: src.GetStart()}
	if !start.IsValid() {
		return false
	}

	startString, err := start.FormatTime("2006-0102")
	if err != nil {
		return false
	}

	companyName := src.GetCompanyName()
	locationName := src.GetLocationName()

	// 事前に容量を計算してstrings.Builderを初期化（再アロケーション回避）
	// 日付(9文字) + スペース(1文字) + 会社名 + スペース(1文字) + 現場名 の概算
	var builder strings.Builder
	builder.Grow(len(startString) + 1 + len(companyName) + 1 + len(locationName))

	// 日付部分を手動構築（YYYY-MMDD形式）
	builder.WriteString(startString)

	// 会社名と現場名を追加
	builder.WriteByte(' ')
	builder.WriteString(companyName)
	builder.WriteByte(' ')
	builder.WriteString(locationName)

	dir := filepath.Dir(m.GetPathistFolder())
	if dir == "." {
		return false
	}
	prevTarget := m.GetPathistFolder()
	target := builder.String()
	src.SetPathistFolder(filepath.Join(dir, target))

	return prevTarget != target
}
