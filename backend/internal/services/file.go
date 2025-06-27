package services

import (
	"errors"
	"io"
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

// CopyFile はファイルまたはディレクトリをコピーする
// srcが相対パス、dstが絶対パスの場合にも対応
func (s *FileService) CopyFile(src, dst string) error {
	var srcFullpath, dstFullpath string
	var err error

	// srcが絶対パスかチェック
	if filepath.IsAbs(src) {
		srcFullpath = src
	} else {
		srcFullpath, err = s.GetFullpath(src)
		if err != nil {
			return err
		}
	}

	// dstが絶対パスかチェック
	if filepath.IsAbs(dst) {
		dstFullpath = dst
	} else {
		dstFullpath, err = s.GetFullpath(dst)
		if err != nil {
			return err
		}
	}

	// コピー元の存在確認
	srcInfo, err := os.Stat(srcFullpath)
	if err != nil {
		return err
	}

	// ディレクトリの場合
	if srcInfo.IsDir() {
		return s.copyDir(srcFullpath, dstFullpath)
	}

	// ファイルの場合
	return s.copyFile(srcFullpath, dstFullpath)
}

// copyFile はファイルをコピーする内部関数
func (s *FileService) copyFile(src, dst string) error {
	// コピー元ファイルを開く
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// コピー先のディレクトリが存在しない場合は作成
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return err
	}

	// コピー先ファイルを作成
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// ファイル内容をコピー
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	// ファイル権限をコピー
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	
	return os.Chmod(dst, srcInfo.Mode())
}

// copyDir はディレクトリを再帰的にコピーする内部関数
func (s *FileService) copyDir(src, dst string) error {
	// コピー元ディレクトリの情報を取得
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// コピー先ディレクトリを作成
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	// ディレクトリ内のエントリを読み取り
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// 各エントリを処理
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// サブディレクトリの場合、再帰的にコピー
			if err := s.copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// ファイルの場合、ファイルをコピー
			if err := s.copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// MoveFile はファイルを移動する
func (s *FileService) MoveFile(src, dst string) error {
	srcFullpath, err := s.GetFullpath(src)
	if err != nil {
		return err
	}
	dstFullpath, err := s.GetFullpath(dst)
	if err != nil {
		return err
	}

	// 移動先のディレクトリが存在するかチェック
	if _, err := os.Stat(srcFullpath); os.IsNotExist(err) {
		return errors.New("移動元のファイル/ディレクトリが存在しません: " + src)
	}

	// 移動先の親ディレクトリを作成（必要に応じて）
	dstParent := filepath.Dir(dstFullpath)
	if err := os.MkdirAll(dstParent, 0755); err != nil {
		return err
	}

	return os.Rename(srcFullpath, dstFullpath)
}

// DeleteFile はファイルを削除する
func (s *FileService) DeleteFile(path string) error {
	fullpath, err := s.GetFullpath(path)
	if err != nil {
		return err
	}

	return os.Remove(fullpath)
}
