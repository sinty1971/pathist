package endpoints

import (
	"context"
	"net/http"
	"net/url"

	"github.com/danielgtaylor/huma/v2"

	app "penguin-backend/internal/app"
	"penguin-backend/internal/models"
)

type GetKojiesRequest struct {
	Filter string `query:"filter" doc:"フィルター (recent: 最近の工事のみ)"`
}

type GetKojiesResponse struct {
	Body []models.Koji `json:"" doc:"工事一覧"`
}

type GetKojiRequest struct {
	Path string `path:"path" doc:"Kojiフォルダー名"`
}

type GetKojiResponse struct {
	Body *models.Koji `json:"" doc:"工事詳細"`
}

type PutKojiRequest struct {
	Body models.Koji `json:""`
}

type PutKojiResponse struct {
	Body models.Koji `json:""`
}

type PutKojiStandardFilesRequest struct {
	Body struct {
		Koji     models.Koji `json:"koji"`
		Currents []string    `json:"currents"`
	} `json:""`
}

type PutKojiStandardFilesResponse struct {
	Body interface{} `json:"" doc:"更新後の工事データ"`
}

func registerKojiEndpoints(api huma.API, container app.ServiceContainer) {
	if container.Koji == nil {
		return
	}

	huma.Register(api, huma.Operation{
		OperationID: "get-kojies",
		Method:      http.MethodGet,
		Path:        "/kojies",
		Summary:     "工事一覧の取得",
		Description: "工事一覧を取得します。filter=recentで最近の工事のみを取得できます",
		Tags:        []string{"工事管理"},
	}, func(ctx context.Context, in *GetKojiesRequest) (*GetKojiesResponse, error) {
		_ = in // 現状はフィルター未対応
		kojies := container.Koji.GetKojies()
		return &GetKojiesResponse{Body: kojies}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-koji",
		Method:      http.MethodGet,
		Path:        "/kojies/{path}",
		Summary:     "工事詳細の取得",
		Description: "指定されたフォルダー名の工事を取得します",
		Tags:        []string{"工事管理"},
	}, func(ctx context.Context, in *GetKojiRequest) (*GetKojiResponse, error) {
		decoded, err := url.QueryUnescape(in.Path)
		if err != nil {
			decoded = in.Path
		}
		koji, err := container.Koji.GetKoji(decoded)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to fetch koji", err)
		}
		return &GetKojiResponse{Body: koji}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "put-koji",
		Method:      http.MethodPut,
		Path:        "/kojies",
		Summary:     "工事情報の更新",
		Description: "工事のファイル情報を更新します",
		Tags:        []string{"工事管理"},
	}, func(ctx context.Context, in *PutKojiRequest) (*PutKojiResponse, error) {
		if err := container.Koji.Update(&in.Body); err != nil {
			return nil, huma.Error500InternalServerError("failed to update koji", err)
		}
		return &PutKojiResponse{Body: in.Body}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "put-koji-standard-files",
		Method:      http.MethodPut,
		Path:        "/kojies/standard-files",
		Summary:     "標準ファイル名の変更",
		Description: "StandardFileの名前を変更し、最新の工事データを返します",
		Tags:        []string{"工事管理"},
	}, func(ctx context.Context, in *PutKojiStandardFilesRequest) (*PutKojiStandardFilesResponse, error) {
		folderName := in.Body.Koji.GetFolderName()
		if folderName == "" {
			return &PutKojiStandardFilesResponse{Body: []string{}}, nil
		}
		updated, err := container.Koji.GetKoji(folderName)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to load updated koji", err)
		}
		return &PutKojiStandardFilesResponse{Body: updated}, nil
	})
}
