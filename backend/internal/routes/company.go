package routes

import (
	"penguin-backend/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupCompanyRoutes(api fiber.Router, handler *handlers.CompanyHandler) {
	company := api.Group("/company")

	// 会社一覧を取得
	// @Summary 会社一覧を取得
	// @Description 会社一覧を取得
	// @Tags company
	// @Accept json
	// @Produce json
	// @Success 200 {array} models.Company
	// @Router /company/all [get]
	company.Get("/all", handler.GetCompanies)

	// 会社を取得
	// @Summary 会社を取得
	// @Description 会社を取得
	// @Tags company
	// @Accept json
	// @Produce json
	// @Success 200 {object} models.Company
	// @Router /company/{id} [get]
	company.Get("/:id", handler.GetCompany)
}
