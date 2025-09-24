package routes

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v2"

	"penguin-backend/internal/routes/business"
	"penguin-backend/internal/services"
)

// apiRootDescription はルートエンドポイントの説明
const apiRootDescription = "Penguin Backend API"

// SetupRoutes はすべてのルートを設定します
func SetupRoutes(app *fiber.App, api huma.API, rootService *services.RootService) {
	huma.Register(api, huma.Operation{
		OperationID: "get-root",
		Method:      http.MethodGet,
		Path:        "/",
		Summary:     "API ルート情報取得",
		Description: "API のバージョンやドキュメントURLを返します",
		Tags:        []string{"system"},
	}, func(ctx context.Context, _ *struct{}) (*struct {
		Body struct {
			Message string `json:"message" example:"Penguin Backend API"`
			Version string `json:"version" example:"1.0.0"`
			Docs    string `json:"docs" example:"/docs"`
		}
	}, error) {
		resp := &struct {
			Body struct {
				Message string `json:"message" example:"Penguin Backend API"`
				Version string `json:"version" example:"1.0.0"`
				Docs    string `json:"docs" example:"/docs"`
			}
		}{}
		resp.Body.Message = apiRootDescription
		resp.Body.Version = "1.0.0"
		resp.Body.Docs = "/docs"
		return resp, nil
	})

	// 既存 Fiber ルート（段階的移行用）
	apiGroup := app.Group("/api")
	officeAPI := apiGroup.Group("/toyotachikurul")

	service := rootService.GetService("BusinessService")
	if service != nil {
		if businessService, ok := (*service).(*services.BusinessService); ok {
			business.SetupBusinessRoutes(officeAPI, businessService)
		}
	}
}
