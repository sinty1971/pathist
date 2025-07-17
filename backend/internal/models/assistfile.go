package models

// AssistFile は補助ファイルの情報を表します
type AssistFile struct {
	// FileService.BasePathからの相対パス
	CurrentPath string `json:"currentPath" yaml:"current_path" example:"豊田築炉/2-工事/2025-0618 豊田築炉 名和工場"`
	// 現在のファイル名
	Current string `json:"current" yaml:"current" example:"工事.xlsx"`

	// FileService.BasePathからの相対パス
	DesiredPath string `json:"desiredPath" yaml:"desired_path" example:"豊田築炉/2-工事/2025-0618 豊田築炉 名和工場"`
	// 希望するファイル名
	Desired string `json:"desired" yaml:"desired" example:"2025-0618 豊田築炉 名和工場.xlsx"`
}
