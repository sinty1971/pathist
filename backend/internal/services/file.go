package services

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"
	"penguin-backend/internal/models"
	"strings"
)

// FileService is a service for managing the file system
type FileService struct {
	// BasePath is the base directory of the file system
	BasePath string `json:"base_path" yaml:"base_path" example:"~/penguin"`
}

// NewFileService creates a new FileService
func NewFileService(basePath string) (*FileService, error) {
	// Expand ~ to home directory
	if strings.HasPrefix(basePath, "~/") {
		usr, err := user.Current()
		if err != nil {
			return nil, err
		}
		basePath = filepath.Join(usr.HomeDir, basePath[2:])
	} else {
		absPath, err := filepath.Abs(basePath)
		if err != nil {
			return nil, err
		}
		basePath = absPath
	}

	return &FileService{
		BasePath: basePath,
	}, nil
}

// GetFullpath BasePathに可変長引数の相対パスを追加したパスを返す
// 絶対パスの場合はエラーを返す
func (s *FileService) GetFullpath(joinPaths ...string) (string, error) {
	if len(joinPaths) == 0 {
		return s.BasePath, nil
	}

	// BasePathに相対パスを追加したパスを返す
	joinedPath := filepath.Join(joinPaths...)

	if strings.HasPrefix(joinedPath, "~/") {
		// 接頭語 "~/" がある場合はエラーを返す
		return "", errors.New("接頭語 \"~\" は使用できません")

	} else if filepath.IsAbs(joinedPath) {
		// 絶対パスがある場合はエラーを返す
		return "", errors.New("絶対パスは使用できません")
	}

	return filepath.Join(s.BasePath, joinedPath), nil
}

// GetFileInfos 指定されたディレクトリパス内のファイル情報を取得
// 取得したファイル情報は実体のファイル情報を取得して判定したもの
func (s *FileService) GetFileInfos(joinPaths ...string) ([]models.FileInfo, error) {
	// 絶対パスを取得
	fullpath, err := s.GetFullpath(joinPaths...)
	if err != nil {
		return nil, err
	}

	// ファイルエントリ配列を取得
	entries, err := os.ReadDir(fullpath)
	if err != nil {
		return nil, err
	}

	// ファイル情報を作成
	fileInfos := make([]models.FileInfo, len(entries))
	var count int
	for _, entry := range entries {
		// ファイルパスを取得
		entryFullpath := filepath.Join(fullpath, entry.Name())

		// FileInfoを作成
		fileInfo, err := models.NewFileInfo(entryFullpath)
		if err != nil {
			// ファイル情報の取得に失敗した場合はスキップ
			continue
		}

		// ファイル情報の登録
		fileInfos[count] = *fileInfo
		count++
	}

	return fileInfos[:count], nil
}
