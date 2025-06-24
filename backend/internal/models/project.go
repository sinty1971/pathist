package models

import (
	"fmt"
	"strings"
	"time"
)

// Project represents a construction project folder with additional metadata
// @Description Construction project folder information with extended attributes
type Project struct {
	// Embed the base FileInfo struct
	FileInfo

	// Calculated fields
	ID     string `json:"id,omitempty" yaml:"-" example:"TC618"`
	Status string `json:"status,omitempty" yaml:"-" example:"進行中"`

	// Additional fields specific to Project folders
	CompanyName  string    `json:"company_name,omitempty" yaml:"company_name" example:"豊田築炉"`
	LocationName string    `json:"location_name,omitempty" yaml:"location_name" example:"名和工場"`
	StartDate    Timestamp `json:"start_date,omitempty" yaml:"start_date"`

	// Detail
	DetailFile  string    `json:"detail_file,omitempty" yaml:"detail_file" example:".detail.yaml"`
	EndDate     Timestamp `json:"end_date,omitempty" yaml:"end_date"`
	Description string    `json:"description,omitempty" yaml:"description" example:"工事関連の資料とドキュメント"`
	Tags        []string  `json:"tags,omitempty" yaml:"tags" example:"['工事', '豊田築炉', '名和工場']"`

	// Managed items
	ProjectFile string `json:"project_file,omitempty" yaml:"project_file" example:"2025-0618 豊田築炉 名和工場.xlsx"`
}

// NewProject FileInfoからProjectを作成します
func NewProject(fileInfo FileInfo) (Project, error) {

	// ファイル名から日付を取得と残りの文字列を取得
	var nameWithoutDate string
	startDate, err := ParseTimestamp(fileInfo.Name, &nameWithoutDate)
	if err != nil {
		return Project{}, err
	}

	// nameWithoutDate は "豊田築炉 名和工場 詳細" のような形式の文字列
	// companyName は会社名(ex. 豊田築炉)
	// locationName は会社名以降の文字列(ex. "名和工場 詳細")
	parts := strings.Split(nameWithoutDate, " ")
	companyName := parts[0]
	var locationName string
	if len(parts) > 1 {
		locationName = strings.Join(parts[1:], " ")
	}

	// Projectインスタンスを作成
	project := Project{
		FileInfo: fileInfo,

		ID:     GetProjectID(startDate, companyName, locationName),
		Status: DetermineProjectStatus(startDate, startDate),

		CompanyName:  companyName,
		LocationName: locationName,
		StartDate:    startDate,

		DetailFile:  ".detail.yaml",
		EndDate:     startDate,
		Description: companyName + "の" + locationName + "における工事情報",
		Tags:        []string{"Project", "工事", companyName, locationName, startDate.Time.Format("2006")}, // Include year as tag
	}

	return project, nil
}

// DeterminKoujiID は工事IDを決定する
func GetProjectID(startDate Timestamp, companyName string, locationName string) string {
	startDateStr, err := startDate.Format("20060102")
	if err != nil {
		return ""
	}
	idSource := fmt.Sprintf("%s%s%s", startDateStr, companyName, locationName)
	id := NewIDFromString(idSource)
	return id.Len5()
}

// DetermineProjectStatus determines the project status based on the date
func DetermineProjectStatus(startDate Timestamp, endDate Timestamp) string {
	if startDate.Time.IsZero() {
		return "不明"
	}

	now := time.Now()

	if now.Before(startDate.Time) {
		return "予定"
	} else if now.After(endDate.Time) {
		return "完了"
	} else {
		return "進行中"
	}
}
