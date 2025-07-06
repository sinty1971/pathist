package models

import (
	"strings"
	"time"
)

// Koji は追加のメタデータを持つ工事プロジェクトフォルダーを表します
// @Description 拡張属性を持つ工事プロジェクトフォルダー情報
type Koji struct {
	// 基本のFileInfo構造体を埋め込み
	FileInfo

	// 計算フィールド
	ID     string `json:"id,omitempty" yaml:"-" example:"TC618"`
	Status string `json:"status,omitempty" yaml:"-" example:"進行中"`

	// パス名からの固有フィールド
	CompanyName  string    `json:"company_name,omitempty" yaml:"company_name" example:"豊田築炉"`
	LocationName string    `json:"location_name,omitempty" yaml:"location_name" example:"名和工場"`
	StartDate    Timestamp `json:"start_date,omitempty" yaml:"start_date"`

	// 属性ファイルフィールド
	EndDate      Timestamp     `json:"end_date,omitempty" yaml:"end_date"`
	Description  string        `json:"description,omitempty" yaml:"description" example:"工事関連の資料とドキュメント"`
	Tags         []string      `json:"tags,omitempty" yaml:"tags" example:"['工事', '豊田築炉', '名和工場']"`
	ManagedFiles []ManagedFile `json:"managed_files" yaml:"managed_files"`
}

// GetFileInfo AttributeServiceで使用するためのメソッド
func (k *Koji) GetFileInfo() FileInfo {
	return k.FileInfo
}

// SetFileInfo AttributeServiceで使用するためのメソッド
func (k *Koji) SetFileInfo(fileInfo FileInfo) {
	k.FileInfo = fileInfo
}

// NewKoji FileInfoからKojiを作成します（高速化版）
func NewKoji(fileInfo FileInfo) (Koji, error) {
	// ファイル名から日付を取得と残りの文字列を取得
	var nameWithoutDate string
	startDate, err := ParseTimestamp(fileInfo.Name, &nameWithoutDate)
	if err != nil {
		return Koji{}, err
	}

	// nameWithoutDate の解析を最適化
	companyName, locationName := parseKojiName(nameWithoutDate)

	// Koji名を生成（一度だけ）
	kojiName := fileInfo.Name

	// Kojiインスタンスを作成（構造体リテラルで一度に初期化）
	koji := Koji{
		FileInfo: fileInfo,

		ID:           GenerateKojiID(startDate, companyName, locationName),
		Status:       DetermineKojiStatus(startDate),
		CompanyName:  companyName,
		LocationName: locationName,
		StartDate:    startDate,
		EndDate:      startDate,
		Description:  buildDescription(companyName, locationName),
		Tags:         buildTags(companyName, locationName, startDate),
	}

	// Nameフィールドを設定（FileInfoに含まれない場合）
	if koji.Name == "" {
		koji.Name = kojiName
	}

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

// GenerateKojiID は工事IDを生成する
func GenerateKojiID(startDate Timestamp, companyName string, locationName string) string {
	startDateStr, err := startDate.Format("20060102")
	if err != nil {
		return ""
	}

	// 文字列結合を最適化（strings.Builderを使用）
	var builder strings.Builder
	builder.Grow(len(startDateStr) + len(companyName) + len(locationName))
	builder.WriteString(startDateStr)
	builder.WriteString(companyName)
	builder.WriteString(locationName)

	id := NewIDFromString(builder.String())
	return id.Len5()
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
func (km *Koji) UpdateID() {
	km.ID = GenerateKojiID(km.StartDate, km.CompanyName, km.LocationName)
}

// UpdateStatus はプロジェクトステータスを更新する
func (km *Koji) UpdateStatus() {
	km.Status = DetermineKojiStatus(km.StartDate, km.EndDate)
}
