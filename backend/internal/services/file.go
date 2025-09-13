package services

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"penguin-backend/internal/models"
	"penguin-backend/internal/utils"
	"runtime"
	"strings"
	"sync"
)

// FileService はファイルシステムを管理するサービス
type FileService struct {
	// Container はトップコンテナのインスタンス
	Container *RootService

	// BaseFolder はファイルシステムの絶対パスフォルダー
	BaseFolder string `json:"baseFolder" yaml:"base_folder" example:"/penguin/豊田築炉"`
}

// ResetWithContext FileServiceを引数のコンテナに設定し、基準フォルダーを設定する
// buildContext はファイルサービスの基準フォルダー
// 戻り値はファイルサービスのインスタンス
func (fs *FileService) BuildWithOption(rs *RootService, folderPath string) error {
	// コンテナを設定
	fs.Container = rs.RootService

	// 絶対パスに展開
	cleanBaseFolder, err := utils.CleanAbsPath(folderPath)
	if err != nil {
		return err
	}

	// 基準フォルダーにアクセスできるかチェック
	osFileInfo, err := os.Stat(cleanBaseFolder)
	if err != nil {
		return err
	}

	// フォルダーかどうかをチェック
	if !osFileInfo.IsDir() {
		return errors.New("フォルダーではありません")
	}

	// 基準フォルダーを設定
	fs.BaseFolder = cleanBaseFolder

	return nil
}

// GetFileInfos 指定されたディレクトリパス内のファイル情報を取得
// target はディレクトリパス、省略時はBaseFolder
func (fs *FileService) GetFileInfos(target ...string) ([]models.FileInfo, error) {
	// パスの結合
	targetPath := filepath.Join(target...)

	// ターゲットフォルダーの絶対パスを取得
	targetFolder, err := fs.JoinBasePath(targetPath)
	if err != nil {
		return nil, err
	}

	// ファイルエントリ配列を取得
	targetEntries, err := os.ReadDir(targetFolder)
	if err != nil {
		return nil, err
	}
	// ファイルエントリが0の場合は空配列を返す
	if len(targetEntries) == 0 {
		return []models.FileInfo{}, nil
	}

	// チャンネルとワーカーグループを設定
	jobs := make(chan int, len(targetEntries))
	fis := make(chan *models.FileInfo, len(targetEntries))
	// 並列処理用のワーカー数を決定（CPU数の2倍、最大16）
	numWorkers := min(min(runtime.NumCPU()*2, 16), len(targetEntries))
	// ワーカーグループを設定
	var wg sync.WaitGroup
	// ワーカーを起動
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range jobs {
				entry := targetEntries[index]
				fullpath := filepath.Join(targetFolder, entry.Name())

				fileInfo, _ := models.NewFileInfo(fullpath)
				if fileInfo != nil {
					fis <- fileInfo
				}
			}
		}()
	}

	// ジョブを送信
	for i := range targetEntries {
		jobs <- i
	}
	close(jobs)

	// ワーカーの完了を待つ
	go func() {
		wg.Wait()
		close(fis)
	}()

	// 結果を収集
	fileInfos := make([]models.FileInfo, 0, len(targetEntries))
	for fileInfo := range fis {
		fileInfos = append(fileInfos, *fileInfo)
	}

	return fileInfos, nil
}

// JoinBasePath BaseFolderに可変長引数の相対パスを追加した絶対パスを返す
// 引数が絶対パスの場合はエラーを返す
func (fs *FileService) JoinBasePath(target ...string) (string, error) {
	if len(target) == 0 {
		return fs.BaseFolder, nil
	}

	// BasePathに相対パスを追加したパスを返す
	targetPath := filepath.Join(target...)

	if strings.HasPrefix(targetPath, "~/") {
		// 接頭語 "~/" がある場合はエラーを返す
		return "", errors.New("接頭語 \"~\" は使用できません")

	} else if filepath.IsAbs(targetPath) {
		// 絶対パスがある場合はエラーを返す
		return "", errors.New("絶対パスは使用できません")
	}

	return filepath.Join(fs.BaseFolder, targetPath), nil
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
		srcFullpath, err = s.JoinBasePath(src)
		if err != nil {
			return err
		}
	}

	// dstが絶対パスかチェック
	if filepath.IsAbs(dst) {
		dstFullpath = dst
	} else {
		dstFullpath, err = s.JoinBasePath(dst)
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
	srcFullpath, err := s.JoinBasePath(src)
	if err != nil {
		return err
	}
	dstFullpath, err := s.JoinBasePath(dst)
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
	fullpath, err := s.JoinBasePath(path)
	if err != nil {
		return err
	}

	return os.Remove(fullpath)
}
