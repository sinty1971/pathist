package models

import (
	"time"
)

type Member struct {
	ID          ID         `json:"id" yaml:"-" example:"ABC1234"`
	Name        string     `json:"name" yaml:"name" example:"first middle Last"`
	FolderPath  string     `json:"folderPath" yaml:"folder_path" example:"/home/shin/penguin/豊田築炉/1 会社/豊田築炉工業/社員/山田太郎"`
	CompanyName string     `json:"companyName" yaml:"company_name" example:"豊田築炉工業株式会社"`
	Position    string     `json:"position,omitempty" yaml:"position" example:"部長"`
	Department  string     `json:"department,omitempty" example:"営業部"`
	Email       string     `json:"email,omitempty" example:"yamada@example.com"`
	Phone       string     `json:"phone,omitempty" example:"090-1234-5678"`
	JoinDate    *time.Time `json:"joinDate,omitempty" yaml:"join_date" example:"2020-04-01T00:00:00Z"`

	// 補助ファイルフィールド
	Database *Repository[*Member] `json:"-" yaml:"-"`
	Assets   []FileInfo           `json:"assets" yaml:"assets"`
}

// GetFolderPath Databaseインターフェースの実装
func (mm *Member) GetFolderPath() string {
	return mm.FolderPath
}
