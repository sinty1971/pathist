package services

import (
	"errors"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"penguin-backend/internal/models"
	"runtime"
	"strings"
	"sync"
)

// FileService はファイルシステムを管理するサービス
type FileService struct {
	// BasePath はファイルシステムの基底ディレクトリ
	BasePath string `json:"base_path" yaml:"base_path" example:"~/penguin"`
}

// NewFileService FileServiceを作成する
func NewFileService(basePath string) (*FileService, error) {
	// ホームディレクトリに展開
	if strings.HasPrefix(basePath, "~/") {
		usr, err := user.Current()
		if err != nil {
			return nil, err
		}
		basePath = filepath.Join(usr.HomeDir, basePath[2:])
	}

	// 絶対パスに変換（チェック）
	absPath, err := filepath.Abs(basePath)
	if err != nil {
		return nil, err
	}
	basePath = absPath

	return &FileService{
		BasePath: basePath,
	}, nil
}

// FileInfoResult 並列処理用の結果構造体
type FileInfoResult struct {
	FileInfo *models.FileInfo
	Index    int
	Error    error
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

	if len(entries) == 0 {
		return []models.FileInfo{}, nil
	}

	// 並列処理用のワーカー数を決定（CPU数の2倍、最大16）
	numWorkers := min(min(runtime.NumCPU()*2, 16), len(entries))

	// チャンネルとワーカーグループを設定
	jobs := make(chan int, len(entries))
	results := make(chan FileInfoResult, len(entries))
	var wg sync.WaitGroup

	// ワーカーを起動
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range jobs {
				entry := entries[index]
				entryFullpath := filepath.Join(fullpath, entry.Name())

				fileInfo, err := models.NewFileInfo(entryFullpath)
				results <- FileInfoResult{
					FileInfo: fileInfo,
					Index:    index,
					Error:    err,
				}
			}
		}()
	}

	// ジョブを送信
	for i := range entries {
		jobs <- i
	}
	close(jobs)

	// ワーカーの完了を待つ
	go func() {
		wg.Wait()
		close(results)
	}()

	// 結果を収集（元の順序を保持）
	fileInfos := make([]models.FileInfo, 0, len(entries))
	resultMap := make(map[int]*models.FileInfo, len(entries))

	for result := range results {
		if result.Error == nil && result.FileInfo != nil {
			resultMap[result.Index] = result.FileInfo
		}
	}

	// 元の順序でファイル情報を構築
	for i := range len(entries) {
		if fileInfo, exists := resultMap[i]; exists {
			fileInfos = append(fileInfos, *fileInfo)
		}
	}

	return fileInfos, nil
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
