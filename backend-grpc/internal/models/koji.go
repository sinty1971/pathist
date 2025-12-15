package models

import (
	grpcv1 "backend-grpc/gen/grpc/v1"
	"backend-grpc/internal/core"
	"errors"
	"path/filepath"
	"strings"
	"time"
)

type Koji struct {
	// Koji メッセージ本体
	*grpcv1.Koji

	// Common 共通モデルフィールド
	*core.PathistModel
}

// NewKoji FolderNameからKojiを作成します（高速化版）
func NewKoji() *Koji {

	koji := &Koji{}
	koji.Koji = grpcv1.Koji_builder{}.Build()
	koji.PathistModel = core.NewPathistModel(core.ConfigMap["KojiPersistFilename"])

	return koji
}

// ParseKojiTarget は target から工事開始日・会社名・現場名を取得
func (m *Koji) ParseKojiTarget(target string) error {

	var (
		start        *Timestamp
		companyName  string
		locationName string
		err          error
	)

	// フォルダー名を取得
	foldername := core.GetBaseName(target)

	// ファイル名から工事開始日の取得と日付除外文字列の取得
	var withoutDate string
	start, err = ParseTimestamp(foldername, &withoutDate)
	if err != nil {
		return err
	}

	if withoutDate == "" {
		return errors.New("フォルダー名が工事フォルダーの書式に合致していません")
	}

	// 最初のスペースで分割（最適化）
	if spaceIndex := strings.Index(withoutDate, " "); spaceIndex > 0 {
		companyName = withoutDate[:spaceIndex]
		if spaceIndex+1 < len(withoutDate) {
			locationName = withoutDate[spaceIndex+1:]
		}
	} else {
		companyName = withoutDate
	}

	m.SetTarget(target)
	m.SetStart(start.Timestamp)
	m.SetCompanyName(companyName)
	m.SetLocationName(locationName)
	m.SetPersistDescription(companyName + " " + locationName)

	// IDの設定
	m.SetMessageId()

	return nil
}

// generateDescription は説明文を効率的に構築
func generateDescription(companyName, locationName string) string {
	if companyName == "" {
		return "工事情報"
	}
	if locationName == "" {
		return companyName + "の工事情報"
	}
	return companyName + "の" + locationName + "における工事情報"
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

func (obj *Koji) Update(updatedKoji *Koji) (*Koji, error) {
	if obj == nil || updatedKoji == nil {
		return nil, errors.New("koji or updatedKoji is nil")
	}

	// 管理フォルダーは変更しない
	updatedKoji.SetTarget(obj.GetTarget())

	// 永続化サービスの設定を引き継ぐ
	// updatedKoji.PersistFilename = obj.PersistFilename

	return updatedKoji, nil
}

// UpdateFolderPath は工事フォルダー名を更新します
// ID 及び Target の情報は無視されます
// TODO: 不完全です、実際にはまだ更新処理していません
func (obj *Koji) UpdateFolderPath(src *Koji) bool {
	if src == nil {
		return false
	}

	// 開始日が無効な場合の早期リターン
	start := &Timestamp{Timestamp: src.GetStart()}
	if start.IsValid() == false {
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

	dir := filepath.Dir(obj.GetTarget())
	if dir == "." {
		return false
	}
	prevTarget := obj.GetTarget()
	target := builder.String()
	src.SetTarget(filepath.Join(dir, target))

	return prevTarget != target
}
