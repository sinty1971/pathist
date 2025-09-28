package endpoints

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"penguin-backend/internal/models"
	"penguin-backend/internal/services"
)

type GetCompaniesRequest struct{}

type GetCompaniesResponse struct {
	Body []models.Company `json:"" doc:"会社一覧"`
}

type GetCompanyRequest struct {
	ID string `path:"id" example:"toyota" doc:"会社ID"`
}

type GetCompanyResponse struct {
	Body *models.Company `json:"" doc:"会社詳細"`
}

type PutCompanyRequest struct {
	Body models.Company `json:""`
}

type PutCompanyResponse struct {
	Body models.Company `json:""`
}

type GetCompanyCategoriesRequest struct{}

type GetCompanyCategoriesResponse struct {
	Body []models.CompanyCategoryInfo `json:"" doc:"カテゴリー一覧"`
}

// registerCompanyEndpoints は会社管理エンドポイントを登録します。
func registerCompanyEndpoints(api huma.API, container services.Container) {
	cs := container.CompanyService
	if cs == nil {
		return
	}

	if cs.DatabaseService == nil {
		return
	}

	// 会社一覧の取得
	huma.Register(api, huma.Operation{
		OperationID: "get-companies",
		Method:      http.MethodGet,
		Path:        "/companies",
		Summary:     "会社一覧の取得",
		Description: "会社フォルダーの一覧を取得します",
		Tags:        []string{"会社管理"},
	}, func(ctx context.Context, _ *GetCompaniesRequest) (*GetCompaniesResponse, error) {
		return &GetCompaniesResponse{Body: cs.GetCompanies()}, nil
	})

	// 会社詳細の取得
	huma.Register(api, huma.Operation{
		OperationID: "get-company",
		Method:      http.MethodGet,
		Path:        "/companies/{id}",
		Summary:     "会社詳細の取得",
		Description: "指定されたIDの会社詳細を取得します",
		Tags:        []string{"会社管理"},
	}, func(ctx context.Context, in *GetCompanyRequest) (*GetCompanyResponse, error) {
		company, err := cs.GetCompanyByID(in.ID)
		if err != nil {
			return nil, huma.Error404NotFound("company not found", err)
		}
		return &GetCompanyResponse{Body: company}, nil
	})

	// 会社情報の更新
	huma.Register(api, huma.Operation{
		OperationID: "put-company",
		Method:      http.MethodPut,
		Path:        "/companies",
		Summary:     "会社情報の更新",
		Description: "会社情報を更新し、必要に応じてフォルダー名を変更します",
		Tags:        []string{"会社管理"},
	}, func(ctx context.Context, in *PutCompanyRequest) (*PutCompanyResponse, error) {
		if err := cs.Update(&in.Body); err != nil {
			return nil, huma.Error500InternalServerError("failed to update company", err)
		}
		return &PutCompanyResponse{Body: in.Body}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-company-categories",
		Method:      http.MethodGet,
		Path:        "/companies/categories",
		Summary:     "カテゴリー一覧の取得",
		Description: "カテゴリー一覧を取得します",
		Tags:        []string{"会社管理"},
	}, func(ctx context.Context, _ *GetCompanyCategoriesRequest) (*GetCompanyCategoriesResponse, error) {
		return &GetCompanyCategoriesResponse{Body: cs.Categories()}, nil
	})
}
