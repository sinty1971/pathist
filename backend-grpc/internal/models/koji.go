package models

import (
	grpcv1 "backend-grpc/gen/grpc/v1"
	"backend-grpc/internal/ext"
	"errors"
	"math/big"
	"path/filepath"
	"strings"
	"time"
)

type Koji struct {
	// Koji メッセージ本体
	*grpcv1.Koji

	// persistFilename は永続化サービス用のファイル名
	persistFilename string
}

// GetPersistPath は永続化ファイルのパスを取得します
// Persistable インターフェースの実装
func (k *Koji) GetPersistPath() string {
	return filepath.Join(k.GetManagedFolder(), k.persistFilename)
}

// GetPersistInfo は永続化対象のオブジェクトを取得します
// Persistable インターフェースの実装
func (k *Koji) GetPersistInfo() any {
	return k.Koji
}

// SetPersistInfo は永続化対象のオブジェクトを設定します
// Persistable インターフェースの実装
func (k *Koji) SetPersistInfo(obj any) {
	if koji, ok := obj.(*grpcv1.Koji); ok {
		k.Koji = koji
	}
}

// NewKoji FolderNameからKojiを作成します（高速化版）
func NewKoji(managedFolder string) (*Koji, error) {

	// withoutDate の解析を最適化
	start, companyName, locationName, err := parseKojiManagedFolder(managedFolder)
	if err != nil {
		return nil, err
	}

	// Kojiインスタンスを作成（構造体リテラルで一度に初期化）
	koji := &Koji{
		Koji: grpcv1.Koji_builder{
			Id:                  GenerateKoujiId(start, companyName, locationName),
			ManagedFolder:       managedFolder,
			Start:               start.Timestamp,
			CompanyName:         companyName,
			LocationName:        locationName,
			Status:              GenerateKojiStatus(start, start),
			InsideEnd:           start.Timestamp,
			InsideDescription:   generateDescription(companyName, locationName),
			InsideTags:          generateTags(companyName, locationName, start),
			InsideRequiredFiles: []*grpcv1.FileInfo{},
		}.Build(),
	}

	return koji, nil
}

// parseKojiManagedFolder は managedFolder から工事開始日・会社名・現場名を取得
func parseKojiManagedFolder(managedFolder string) (*Timestamp, string, string, error) {

	var (
		start        *Timestamp
		companyName  string
		locationName string
		err          error
	)

	// フォルダー名を取得
	foldername := ext.GetBaseName(managedFolder)

	// ファイル名から工事開始日の取得と日付除外文字列の取得
	var withoutDate string
	start, err = ParseTimestamp(foldername, &withoutDate)
	if err != nil {
		return nil, "", "", err
	}

	if withoutDate == "" {
		return nil, "", "", errors.New("フォルダー名が工事フォルダーの書式に合致していません")
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

	return start, companyName, locationName, nil
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

// generateTags はタグ配列を効率的に構築
func generateTags(companyName, locationName string, startDate *Timestamp) []string {
	// 容量を事前確保（通常5-6個のタグ）
	tags := make([]string, 0, 6)
	tags = append(tags, "Koji", "工事")

	if companyName != "" {
		tags = append(tags, companyName)
	}
	if locationName != "" {
		tags = append(tags, locationName)
	}
	if year, err := startDate.FormatTime("2006"); err == nil {
		tags = append(tags, year)
	}

	return tags
}

// GenerateKojiId は工事IDを生成する
func GenerateKoujiId(start *Timestamp, companyName, locationName string) string {

	// ID用のバイト配列を生成
	startBytes := big.NewInt(int64(start.GetSeconds())).Bytes()
	companyBytes := []byte(companyName)
	locationBytes := []byte(locationName)
	n := len(startBytes) + len(companyBytes) + len(locationBytes)
	bytes := make([]byte, n)
	offset := copy(bytes, startBytes)
	offset += copy(bytes[offset:], companyBytes)
	copy(bytes[offset:], locationBytes)

	// バイト配列からハッシュ文字列IDを生成
	return ext.GenerateIdFromBytes(bytes)
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

func (k *Koji) Update(updatedKoji *Koji) (*Koji, error) {
	if k == nil || updatedKoji == nil {
		return nil, errors.New("koji or updatedKoji is nil")
	}

	// 管理フォルダーは変更しない
	updatedKoji.SetManagedFolder(k.GetManagedFolder())

	// 永続化サービスの設定を引き継ぐ
	updatedKoji.persistFilename = k.persistFilename

	return updatedKoji, nil
}

// UpdateFolderPath は工事フォルダー名を更新します
// ID 及び ManagedFolder の情報は無視されます
// TODO: 不完全です、実際にはまだ更新処理していません
func (k *Koji) UpdateFolderPath(src *Koji) bool {
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

	dir := filepath.Dir(k.GetManagedFolder())
	if dir == "." {
		return false
	}
	prevManagedFolder := k.GetManagedFolder()
	managedFolder := builder.String()
	src.SetManagedFolder(filepath.Join(dir, managedFolder))

	return prevManagedFolder != managedFolder
}
