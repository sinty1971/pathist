package endpoints

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"penguin-backend/internal/models"
	"penguin-backend/internal/services"
)

type GetFilesRequest struct {
	Path string `query:"path" doc:"取得するディレクトリの相対パス" example:"2 工事"`
}

type GetFilesResponse struct {
	Body []models.FileInfo `json:"" doc:"ファイルエントリ一覧"`
}

type GetFileBasePathRequest struct{}

type GetFileBasePathResponse struct {
	Body struct {
		BasePath string `json:"basePath" doc:"ファイルサービスの基準パス"`
	} `json:""`
}

func registerFileEndpoints(api huma.API, container services.Container) {
	// FileServiceが無ければ終了
	fs := container.FileService
	if fs == nil {
		return
	}

	huma.Register(api, huma.Operation{
		OperationID: "get-files",
		Method:      http.MethodGet,
		Path:        "/files",
		Summary:     "ファイルエントリ一覧の取得",
		Description: "指定されたパスからファイルとフォルダーの一覧を取得します",
		Tags:        []string{"ファイル管理"},
	}, func(ctx context.Context, in *GetFilesRequest) (*GetFilesResponse, error) {
		var (
			entries []models.FileInfo
			err     error
		)
		if in.Path == "" {
			entries, err = fs.GetFileInfos()
		} else {
			entries, err = fs.GetFileInfos(in.Path)
		}
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to read directory", err)
		}
		return &GetFilesResponse{Body: entries}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-file-base-path",
		Method:      http.MethodGet,
		Path:        "/files/base-path",
		Summary:     "ファイルサービスの基準パス取得",
		Description: "ファイルサービスが参照する基準パスを返します",
		Tags:        []string{"ファイル管理"},
	}, func(ctx context.Context, _ *GetFileBasePathRequest) (*GetFileBasePathResponse, error) {
		resp := &GetFileBasePathResponse{}
		resp.Body.BasePath = fs.TargetFolder
		return resp, nil
	})
}
