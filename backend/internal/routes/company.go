package routes

import (
	"penguin-backend/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

// SetupCompanyRoutes 会社関連のルートを設定
func SetupCompanyRoutes(api fiber.Router, handler *handlers.CompanyHandler) {
	company := api.Group("/company")

	// 会社一覧取得
	company.Get("/list", handler.GetCompanies)

	// 会社詳細取得
	company.Get("/:id", handler.GetCompany)
}
