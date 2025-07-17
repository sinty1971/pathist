package models

import (
	"os"
	"strings"
)

// FileInfo はファイルまたはディレクトリを表す
// @Description ファイルまたはディレクトリの情報
type FileInfo struct {
	// FileService.BasePathからの相対パス
	Path string `json:"path" yaml:"-" example:"豊田築炉/2-工事/2025-0618 豊田築炉 名和工場"`
	// Name of the file or folder
	Name string `json:"name" yaml:"name" example:"2025-0618 豊田築炉 名和工場"`
	// Whether this item is a directory
	IsDirectory bool `json:"isDirectory" yaml:"is_directory" example:"true"`
	// Size of the file in bytes
	Size int64 `json:"size" yaml:"size" example:"4096"`
	// Last modification time
	ModifiedTime Timestamp `json:"modifiedTime" yaml:"modified_time"`
}

// NewFileInfo フルパスからFileInfo構造体を作成します
func NewFileInfo(fullpath string) (*FileInfo, error) {
	osFileInfo, err := os.Stat(fullpath)
	if err != nil {
		return nil, err
	}

	return &FileInfo{
		Name:         osFileInfo.Name(),
		Path:         fullpath,
		IsDirectory:  osFileInfo.IsDir(),
		Size:         osFileInfo.Size(),
		ModifiedTime: Timestamp{Time: osFileInfo.ModTime()},
	}, nil
}

// IsExist ファイルが存在するかどうかをチェック
func (f *FileInfo) IsExist() bool {
	_, err := os.Stat(f.Path)
	return err == nil
}

// IsExcel エクセルファイルかどうかをチェック
func (f *FileInfo) IsExcel() bool {
	if f.IsDirectory {
		return false
	}

	excelSuffix := []string{".xlsx", ".xls"}
	nameLower := strings.ToLower(f.Name)

	for _, suffix := range excelSuffix {
		if strings.HasSuffix(nameLower, suffix) {
			return true
		}
	}
	return false
}
