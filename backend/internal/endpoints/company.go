package endpoints

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	app "penguin-backend/internal/app"
	"penguin-backend/internal/models"
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

func registerCompanyEndpoints(api huma.API, container app.ServiceContainer) {
	if container.Company == nil {
		return
	}

	huma.Register(api, huma.Operation{
		OperationID: "get-companies",
		Method:      http.MethodGet,
		Path:        "/companies",
		Summary:     "会社一覧の取得",
		Description: "会社フォルダーの一覧を取得します",
		Tags:        []string{"会社管理"},
	}, func(ctx context.Context, _ *GetCompaniesRequest) (*GetCompaniesResponse, error) {
		companies := container.Company.GetCompanies()
		return &GetCompaniesResponse{Body: companies}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-company",
		Method:      http.MethodGet,
		Path:        "/companies/{id}",
		Summary:     "会社詳細の取得",
		Description: "指定されたIDの会社詳細を取得します",
		Tags:        []string{"会社管理"},
	}, func(ctx context.Context, in *GetCompanyRequest) (*GetCompanyResponse, error) {
		company, err := container.Company.GetCompanyByID(in.ID)
		if err != nil {
			return nil, huma.Error404NotFound("company not found", err)
		}
		return &GetCompanyResponse{Body: company}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "put-company",
		Method:      http.MethodPut,
		Path:        "/companies",
		Summary:     "会社情報の更新",
		Description: "会社情報を更新し、必要に応じてフォルダー名を変更します",
		Tags:        []string{"会社管理"},
	}, func(ctx context.Context, in *PutCompanyRequest) (*PutCompanyResponse, error) {
		if err := container.Company.Update(&in.Body); err != nil {
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
		categories := container.Company.Categories()
		return &GetCompanyCategoriesResponse{Body: categories}, nil
	})
}
