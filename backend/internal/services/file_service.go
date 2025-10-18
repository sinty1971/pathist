package services

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	penguinv1 "penguin-backend/gen/penguin/v1"
	penguinv1connect "penguin-backend/gen/penguin/v1/penguinv1connect"
	"penguin-backend/internal/adapters"
	"penguin-backend/internal/services"

	"connectrpc.com/connect"
)

// FileService の実装

// FileServiceHandler exposes FileService operations via Connect handlers.
type FileServiceHandler struct {
	penguinv1connect.UnimplementedFileServiceHandler
	// handlers は任意のgrpcサービスハンドラーへの参照
	handlers *Handlers

	// TargetPath はファイルサービスの絶対パスフォルダー
	TargetPath string `json:"targetPath" yaml:"target_path" example:"/penguin/豊田築炉"`
}

func NewFileServiceHandler(service *services.FileService) *FileServiceHandler {
	return &FileServiceHandler{}
}

func (h *FileServiceHandler) ListFiles(
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

				fi, _ := adapters.NewFileInfo(fullpath)
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
func (fs *FileServiceHandler) GetAbsPathFrom(relPath string) (string, error) {

	// 絶対パスがある場合はエラーを返す
	if strings.HasPrefix(relPath, "~/") || filepath.IsAbs(relPath) {
		return "", errors.New("絶対パスは使用できません")
	}

	return filepath.Join(fs.TargetPath, relPath), nil
}
