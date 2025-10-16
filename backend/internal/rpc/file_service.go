package rpc

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	penguinv1 "penguin-backend/gen/penguin/v1"
	penguinv1connect "penguin-backend/gen/penguin/v1/penguinv1connect"
	"penguin-backend/internal/models"
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
	req *connect.Request[penguinv1.ListFilesRequest]) (
	*connect.Response[penguinv1.ListFilesResponse],
	error) {
	_ = ctx

	var (
		entries []penguinv1.FileInfo
		err     error
	)

	targetPath := req.Msg.GetPath()

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

	if path == "" {
		entries, err = h.fileService.GetFileInfos()
	} else {
		entries, err = h.fileService.GetFileInfos(path)
	}
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	files := make([]*penguinv1.FileInfo, 0, len(entries))
	for i := range entries {
		files = append(files, convertModelFileInfo(&entries[i]))
	}

	res := penguinv1.ListFilesResponse_builder{
		Files: files,
	}.Build()

	return connect.NewResponse(res), nil
}

func getFileInfos(targetPath string) ([]penguinv1.FileInfo, error) {
}

func (h *FileServiceHandler) GetFileBasePath(ctx context.Context, _ *connect.Request[penguinv1.GetFileBasePathRequest]) (*connect.Response[penguinv1.GetFileBasePathResponse], error) {
	_ = ctx

	res := penguinv1.GetFileBasePathResponse_builder{
		BasePath: h.fileService.TargetFolder,
	}.Build()

	return connect.NewResponse(res), nil
}
