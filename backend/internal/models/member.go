package models

import (
	"time"
)

type Member struct {
	ID          ID                     `json:"id" yaml:"-" example:"ABC1234"`
	Name        string                 `json:"name" example:"first middle Last"`
	FolderPath  string                 `json:"folderPath" example:"/home/shin/penguin/豊田築炉/1 会社/豊田築炉工業/社員/山田太郎"`
	CompanyID   ID                     `json:"companyId" example:"cm5dtcx6o0000fgpgqh0k43kw"`
	CompanyName string                 `json:"companyName" example:"豊田築炉工業株式会社"`
	Position    string                 `json:"position,omitempty" example:"部長"`
	Department  string                 `json:"department,omitempty" example:"営業部"`
	Email       string                 `json:"email,omitempty" example:"yamada@example.com"`
	Phone       string                 `json:"phone,omitempty" example:"090-1234-5678"`
	JoinDate    *time.Time             `json:"joinDate,omitempty" example:"2020-04-01T00:00:00Z"`
	Detail      map[string]interface{} `json:"detail,omitempty"`
	CreatedAt   time.Time              `json:"createdAt" example:"2024-01-01T00:00:00Z"`
	UpdatedAt   time.Time              `json:"updatedAt" example:"2024-01-01T00:00:00Z"`
}
