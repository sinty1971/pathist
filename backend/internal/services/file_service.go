package services

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	penguinv1 "penguin-backend/gen/penguin/v1"
	penguinv1connect "penguin-backend/gen/penguin/v1/penguinv1connect"
	"penguin-backend/internal/models"

	"connectrpc.com/connect"
)

// FileService の実装

// FileServiceHandler exposes FileService operations via Connect handlers.
type FileServiceHandler struct {
	penguinv1connect.UnimplementedFileServiceHandler
	// handlers は任意のgrpcサービスハンドラーへの参照
	handlers *RootHandler

	// TargetPath はファイルサービスの絶対パスフォルダー
	TargetPath string `json:"targetPath" yaml:"target_path" example:"/penguin/豊田築炉"`
}

func NewFileServiceHandler(rootHandler *RootHandler, targetPath string) *FileServiceHandler {
	return &FileServiceHandler{
		handlers:   rootHandler,
		TargetPath: targetPath,
	}
}

func (h *FileServiceHandler) ListFileInfos(
	ctx context.Context,
	req *connect.Request[penguinv1.ListFileInfosRequest]) (
	*connect.Response[penguinv1.ListFileInfosResponse],
	error) {
	// コンテキストを無視
	_ = ctx

	// 変数定義
	var (
		response *connect.Response[penguinv1.ListFileInfosResponse]
		dirs     []os.DirEntry
		fis      []*penguinv1.FileInfo = []*penguinv1.FileInfo{}
		err      error
	)

	// リクエスト情報の取得
	relPath := req.Msg.GetPath()

	// 絶対パスを取得
	targetPath, err := h.GetAbsPathFrom(relPath)
	if err != nil {
		return nil, err
	}

	// ファイルエントリ配列を取得
	dirs, err = os.ReadDir(targetPath)
	if err != nil {
		return nil, err
	}
	// ファイルエントリが0の場合は空配列を返す

	if len(dirs) == 0 {
		response.Msg.SetFileInfos(fis)
		return response, nil
	}

	// チャンネルとワーカーグループを設定
	jobsChan := make(chan int, len(dirs))
	fisChan := make(chan *penguinv1.FileInfo, len(dirs))
	// 並列処理用のワーカー数を決定（CPU数の2倍、最大16）
	numWorkers := min(min(runtime.NumCPU()*2, 16), len(dirs))
	// ワーカーグループを設定
	var wg sync.WaitGroup
	// ワーカーを起動
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range jobsChan {
				dir := dirs[index]
				fullpath := filepath.Join(targetPath, dir.Name())

				fi, _ := models.NewFileInfo(fullpath)
				if fi != nil {
					fisChan <- fi
				}
			}
		}()
	}

	// ジョブを送信
	for i := range dirs {
		jobsChan <- i
	}
	close(jobsChan)

	// ワーカーの完了を待つ
	go func() {
		wg.Wait()
		close(fisChan)
	}()

	// 結果を収集
	fis = make([]*penguinv1.FileInfo, 0, len(dirs))
	for fi := range fisChan {
		fis = append(fis, fi)
	}

	// レスポンスを更新して返す
	response.Msg.SetFileInfos(fis)
	return response, nil
}

// GetAbsPathFrom BasePathに引数の相対パスを追加した絶対パスを返す
func (h *FileServiceHandler) GetAbsPathFrom(relPath string) (string, error) {

	// 絶対パスがある場合はエラーを返す
	if strings.HasPrefix(relPath, "~/") || filepath.IsAbs(relPath) {
		return "", errors.New("絶対パスは使用できません")
	}

	return filepath.Join(h.TargetPath, relPath), nil
}

// CopyFile はファイルまたはディレクトリをコピーする
func (h *FileServiceHandler) CopyFile(relSrc, relDst string) error {
	var absSrc, absDst string
	var err error

	// relSrcがパスチェック及び絶対パス変換
	absSrc, err = h.GetAbsPathFrom(relSrc)
	if err != nil {
		return err
	}

	// relDstのパスチェック及び絶対パス変換
	absDst, err = h.GetAbsPathFrom(relDst)
	if err != nil {
		return err
	}

	// コピー元の存在確認
	srcOsFi, err := os.Stat(absSrc)
	if err != nil {
		return err
	}

	// ディレクトリの場合
	if srcOsFi.IsDir() {
		return h.absCopyDir(absSrc, absDst)
	}

	// ファイルの場合
	return h.absCopyFile(absSrc, absDst)
}

// absCopyFile はファイルをコピーする内部関数
func (h *FileServiceHandler) absCopyFile(absSrc, absDst string) error {
	// コピー元ファイルを開く
	srcFile, err := os.Open(absSrc)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// コピー先のディレクトリが存在しない場合は作成
	dstDir := filepath.Dir(absDst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return err
	}

	// コピー先ファイルを作成
	dstFile, err := os.Create(absDst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// ファイル内容をコピー
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	// ファイル権限をコピー
	srcInfo, err := os.Stat(absSrc)
	if err != nil {
		return err
	}

	return os.Chmod(absDst, srcInfo.Mode())
}

// absCopyDir はディレクトリを再帰的にコピーする内部関数
func (h *FileServiceHandler) absCopyDir(absSrc, absDst string) error {
	// コピー元ディレクトリの情報を取得
	srcInfo, err := os.Stat(absSrc)
	if err != nil {
		return err
	}

	// コピー先ディレクトリを作成
	if err := os.MkdirAll(absDst, srcInfo.Mode()); err != nil {
		return err
	}

	// ディレクトリ内のエントリを読み取り
	entries, err := os.ReadDir(absSrc)
	if err != nil {
		return err
	}

	// 各エントリを処理
	for _, entry := range entries {
		srcPath := filepath.Join(absSrc, entry.Name())
		dstPath := filepath.Join(absDst, entry.Name())

		if entry.IsDir() {
			// サブディレクトリの場合、再帰的にコピー
			if err := h.absCopyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// ファイルの場合、ファイルをコピー
			if err := h.absCopyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// MoveFile はファイルを移動する
func (h *FileServiceHandler) MoveFile(relSrc, relDst string) error {
	absSrc, err := h.GetAbsPathFrom(relSrc)
	if err != nil {
		return err
	}
	absDst, err := h.GetAbsPathFrom(relDst)
	if err != nil {
		return err
	}

	// 移動先のディレクトリが存在するかチェック
	if _, err := os.Stat(absSrc); os.IsNotExist(err) {
		return errors.New("移動元のファイル/ディレクトリが存在しません: " + relSrc)
	}

	// 移動先の親ディレクトリを作成（必要に応じて）
	dstParent := filepath.Dir(absDst)
	if err := os.MkdirAll(dstParent, 0755); err != nil {
		return err
	}

	return os.Rename(absSrc, absDst)
}

// DeleteFile はファイルを削除する
func (h *FileServiceHandler) DeleteFile(relPath string) error {
	absPath, err := h.GetAbsPathFrom(relPath)
	if err != nil {
		return err
	}

	return os.Remove(absPath)
}
