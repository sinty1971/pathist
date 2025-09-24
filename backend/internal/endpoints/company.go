package endpoints

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"penguin-backend/internal/models"
	"penguin-backend/internal/services"
)

// CompanyEndpoint は会社ドメインの操作を提供します。
type CompanyEndpoint struct {
	service *services.CompanyService
}

// NewCompanyEndpoint はサービスを注入したエンドポイントを返します。
func NewCompanyEndpoint(service *services.CompanyService) *CompanyEndpoint {
	return &CompanyEndpoint{service: service}
}

// RegisterCompanyEndpoints は Huma API に会社関連の操作を登録します。
func RegisterCompanyEndpoints(api huma.API, endpoint *CompanyEndpoint) {
	if endpoint == nil || endpoint.service == nil {
		return
	}

	registerGetCompanies(api, endpoint)
	registerGetCompanyByID(api, endpoint)
	registerUpdateCompany(api, endpoint)
	registerGetCategories(api, endpoint)
}

func registerGetCompanies(api huma.API, endpoint *CompanyEndpoint) {
	type output struct {
		Body []models.Company `json:"" doc:"会社一覧"`
	}

	huma.Register(api, huma.Operation{
		OperationID: "list-companies",
		Method:      http.MethodGet,
		Path:        "/companies",
		Summary:     "会社一覧の取得",
		Description: "会社フォルダーの一覧を取得します",
		Tags:        []string{"会社管理"},
	}, func(ctx context.Context, _ *struct{}) (*output, error) {
		companies := endpoint.service.GetCompanies()
		return &output{Body: companies}, nil
	})
}

func registerGetCompanyByID(api huma.API, endpoint *CompanyEndpoint) {
	type input struct {
		ID string `path:"id" example:"toyota" doc:"会社ID"`
	}
	type output struct {
		Body *models.Company `json:"" doc:"会社詳細"`
	}

	huma.Register(api, huma.Operation{
		OperationID: "get-company",
		Method:      http.MethodGet,
		Path:        "/companies/{id}",
		Summary:     "会社詳細の取得",
		Description: "指定されたIDの会社詳細を取得します",
		Tags:        []string{"会社管理"},
	}, func(ctx context.Context, in *input) (*output, error) {
		company, err := endpoint.service.GetCompanyByID(in.ID)
		if err != nil {
			return nil, huma.Error404NotFound("company not found", err)
		}
		return &output{Body: company}, nil
	})
}

func registerUpdateCompany(api huma.API, endpoint *CompanyEndpoint) {
	type input struct {
		Body models.Company `json:""`
	}
	type output struct {
		Body models.Company `json:""`
	}

	huma.Register(api, huma.Operation{
		OperationID: "update-company",
		Method:      http.MethodPut,
		Path:        "/companies",
		Summary:     "会社情報の更新",
		Description: "会社情報を更新し、必要に応じてフォルダー名を変更します",
		Tags:        []string{"会社管理"},
	}, func(ctx context.Context, in *input) (*output, error) {
		if err := endpoint.service.Update(&in.Body); err != nil {
			return nil, huma.Error500InternalServerError("failed to update company", err)
		}
		return &output{Body: in.Body}, nil
	})
}

func registerGetCategories(api huma.API, endpoint *CompanyEndpoint) {
	type output struct {
		Body []models.CompanyCategoryInfo `json:"" doc:"カテゴリー一覧"`
	}

	huma.Register(api, huma.Operation{
		OperationID: "list-company-categories",
		Method:      http.MethodGet,
		Path:        "/companies/categories",
		Summary:     "カテゴリー一覧の取得",
		Description: "カテゴリー一覧を取得します",
		Tags:        []string{"会社管理"},
	}, func(ctx context.Context, _ *struct{}) (*output, error) {
		categories := endpoint.service.Categories()
		return &output{Body: categories}, nil
	})
}
