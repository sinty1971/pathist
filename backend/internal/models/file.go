package models

import (
	"errors"
	"os"
	"penguin-backend/internal/utils"
)

// FileInfo はファイルまたはディレクトリを表す
type FileInfo struct {
	// 対象ファイルまたはディレクトリのフルパス
	TargetPath string `json:"targetPath" yaml:"target_path" example:"/penguin/豊田築炉/2-工事/2025-0618 豊田築炉 名和工場/工事.xlsx"`

	// 標準ファイルまたはディレクトリのフルパス
	StandardPath string `json:"standardPath" yaml:"standard_path" example:"/penguin/豊田築炉/2-工事/2025-0618 豊田築炉 名和工場.xlsx"`

	// ディレクトリかどうか
	IsDirectory bool `json:"isDirectory" yaml:"-" example:"true"`
	// ファイルサイズ
	Size int64 `json:"size" yaml:"-" example:"4096"`
	// 更新日時
	ModifiedTime Timestamp `json:"modifiedTime" yaml:"-" example:"2025-06-18T12:00:00Z"`
}

// NewFileInfo はファイルのフルパスからFileInfo構造体を作成します
// paths はファイルのフルパスを指定します
// 引数が１つの場合は、ファイルのフルパスを指定します
// 引数が２つの場合は、ファイルのフルパスと標準ファイルのフルパスを指定します
func NewFileInfo(paths ...string) (*FileInfo, error) {
	pathsLen := len(paths)
	if pathsLen == 0 || pathsLen > 2 {
		return nil, errors.New("引数が１つまたは２つ必要です")
	}

	// 実際のフルパスと標準ファイルのフルパスを取得
	fileInfo := &FileInfo{}
	var err error

	if pathsLen >= 1 {
		fileInfo.TargetPath, err = utils.CleanAbsPath(paths[0])
		if err != nil {
			return nil, err
		}
		osFileInfo, err := os.Stat(fileInfo.TargetPath)
		if err != nil {
			return nil, err
		}
		fileInfo.IsDirectory = osFileInfo.IsDir()
		fileInfo.Size = osFileInfo.Size()
		fileInfo.ModifiedTime = Timestamp{Time: osFileInfo.ModTime()}
	}

	if pathsLen == 2 {
		fileInfo.StandardPath, err = utils.CleanAbsPath(paths[1])
		if err != nil {
			return nil, err
		}
	}

	return fileInfo, nil
}
