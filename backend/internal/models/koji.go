package models

import (
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Koji は追加のメタデータを持つ工事プロジェクトフォルダーを表します
// @Description 拡張属性を持つ工事プロジェクトフォルダー情報
type Koji struct {
	// ID, Status is calculated from FolderName.
	ID     string `json:"id" yaml:"-" example:"TC618"`
	Status string `json:"status,omitempty" yaml:"-" example:"進行中"`

	// Identity fields. FolderName is linked StartDate, CompanyName and LocationName.
	FolderPath   string    `json:"folderPath" yaml:"-" example:"~/penguin/豊田築炉/2-工事/2025-0618 豊田築炉 名和工場"`
	StartDate    Timestamp `json:"startDate,omitempty" yaml:"-"`
	CompanyName  string    `json:"companyName,omitempty" yaml:"-" example:"豊田築炉"`
	LocationName string    `json:"locationName,omitempty" yaml:"-" example:"名和工場"`

	// Database Service field

	// Database file fields
	Database    *Repository[*Koji] `json:"-" yaml:"-"`
	EndDate     Timestamp          `json:"endDate" yaml:"end_date"`
	Description string             `json:"description,omitempty" yaml:"description" example:"工事関連の資料とドキュメント"`
	Tags        []string           `json:"tags,omitempty" yaml:"tags" example:"['工事', '豊田築炉', '名和工場']"`

	// 必須ファイルフィールド
	RequiredFiles []FileInfo `json:"requiredFiles" yaml:"required_files"`
}

// GetFolderPath Persistableインターフェースの実装
func (km *Koji) GetFolderPath() string {
	return km.FolderPath
}

func (km *Koji) GetFolderName() string {
	// フォルダー名を取得
	var folderName string
	if lastSlash := strings.LastIndex(km.FolderPath, "/"); lastSlash != -1 {
		folderName = km.FolderPath[lastSlash+1:]
	}

	return folderName
}

// NewKoji FolderNameからKojiを作成します（高速化版）
func NewKoji(folderPath string) (*Koji, error) {

	koji := &Koji{
		FolderPath: folderPath,
	}

	// フォルダー名を取得
	folderName := koji.GetFolderName()

	// ファイル名から日付を取得と残りの文字列を取得
	nameWithoutDate := folderName
	startDate, err := ParseTimestamp(folderPath, &nameWithoutDate)
	if err != nil {
		return nil, err
	}

	// nameWithoutDate の解析を最適化
	companyName, locationName := parseKojiName(nameWithoutDate)

	// Kojiインスタンスを作成（構造体リテラルで一度に初期化）
	koji.StartDate = startDate
	koji.CompanyName = companyName
	koji.LocationName = locationName
	koji.EndDate = startDate
	koji.Description = buildDescription(companyName, locationName)
	koji.Tags = buildTags(companyName, locationName, startDate)
	koji.RequiredFiles = []FileInfo{}

	// IDとステータスの更新
	koji.UpdateID()
	koji.UpdateStatus()

	return koji, nil
}

// parseKojiName は名前文字列を会社名と場所名に分割（高速化）
func parseKojiName(nameWithoutDate string) (companyName, locationName string) {
	if nameWithoutDate == "" {
		return "", ""
	}

	// 最初のスペースで分割（最適化）
	if spaceIndex := strings.Index(nameWithoutDate, " "); spaceIndex > 0 {
		companyName = nameWithoutDate[:spaceIndex]
		if spaceIndex+1 < len(nameWithoutDate) {
			locationName = nameWithoutDate[spaceIndex+1:]
		}
	} else {
		companyName = nameWithoutDate
	}

	return companyName, locationName
}

// buildDescription は説明文を効率的に構築
func buildDescription(companyName, locationName string) string {
	if companyName == "" {
		return "工事情報"
	}
	if locationName == "" {
		return companyName + "の工事情報"
	}
	return companyName + "の" + locationName + "における工事情報"
}

// buildTags はタグ配列を効率的に構築
func buildTags(companyName, locationName string, startDate Timestamp) []string {
	// 容量を事前確保（通常5-6個のタグ）
	tags := make([]string, 0, 6)
	tags = append(tags, "Koji", "工事")

	if companyName != "" {
		tags = append(tags, companyName)
	}
	if locationName != "" {
		tags = append(tags, locationName)
	}
	if !startDate.Time.IsZero() {
		tags = append(tags, startDate.Time.Format("2006"))
	}

	return tags
}

// DetermineKojiStatus はプロジェクトステータスを判定する
func DetermineKojiStatus(startDate Timestamp, endDate ...Timestamp) string {
	if startDate.Time.IsZero() {
		return "不明"
	}

	now := time.Now()

	if now.Before(startDate.Time) {
		return "予定"
	}

	// endDateが指定されている場合のみチェック
	if len(endDate) > 0 && !endDate[0].Time.IsZero() && now.After(endDate[0].Time) {
		return "完了"
	}

	return "進行中"
}

// UpdateID は工事IDを更新する
func (km *Koji) UpdateID() bool {
	prevID := km.ID
	km.ID = NewIDFromString(km.FolderPath).Len5()

	return prevID != km.ID
}

// UpdateStatus はプロジェクトステータスを更新する
func (km *Koji) UpdateStatus() bool {
	prevStatus := km.Status
	km.Status = DetermineKojiStatus(km.StartDate, km.EndDate)

	return prevStatus != km.Status
}

// UpdateFolderPath は工事フォルダー名を更新します
func (km *Koji) UpdateFolderPath() bool {
	if km.StartDate.Time.IsZero() {
		// 開始日が無効な場合の早期リターン
		return false
	}

	t := km.StartDate.Time

	// 日付文字列を手動で構築（Format関数より高速）
	year := t.Year()
	month := int(t.Month())
	day := t.Day()

	// 事前に容量を計算してstrings.Builderを初期化（再アロケーション回避）
	// 日付(9文字) + スペース(1文字) + 会社名 + スペース(1文字) + 現場名 の概算
	estimatedSize := 9 + 1 + len(km.CompanyName) + 1 + len(km.LocationName)
	var builder strings.Builder
	builder.Grow(estimatedSize)

	// 日付部分を手動構築（YYYY-MMDD形式）
	builder.WriteString(strconv.Itoa(year))
	builder.WriteByte('-')

	if month < 10 {
		builder.WriteByte('0')
	}
	builder.WriteString(strconv.Itoa(month))

	if day < 10 {
		builder.WriteByte('0')
	}
	builder.WriteString(strconv.Itoa(day))

	// 会社名と現場名を追加
	builder.WriteByte(' ')
	builder.WriteString(km.CompanyName)
	builder.WriteByte(' ')
	builder.WriteString(km.LocationName)

	dir := filepath.Dir(km.FolderPath)
	if dir == "." {
		return false
	}
	prevFolderName := km.GetFolderName()
	folderName := builder.String()
	km.FolderPath = filepath.Join(dir, folderName)

	return prevFolderName != folderName
}
