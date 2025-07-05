package handlers

import (
	"penguin-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

type CompanyHandler struct {
	CompanyService *services.CompanyService
}

func NewCompanyHandler(companyService *services.CompanyService) *CompanyHandler {
	return &CompanyHandler{
		CompanyService: companyService,
	}
}

// GetCompanyByID godoc
// @Summary      会社詳細の取得
// @Description  指定されたIDの会社詳細を取得します
// @Tags         会社管理
// @Accept       json
// @Produce      json
// @Param        id path string true "会社ID"
// @Success      200 {object} models.Company "正常なレスポンス"
// @Failure      404 {object} map[string]string "会社が見つからない"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /company/{id} [get]
func (h *CompanyHandler) GetCompany(c *fiber.Ctx) error {
	company, err := h.CompanyService.GetCompanyByID(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Company not found",
			"message": err.Error(),
		})
	}
	return c.JSON(company)
}

// GetCompanies godoc
// @Summary      会社一覧の取得
// @Description  会社フォルダーの一覧を取得します
// @Tags         会社管理
// @Accept       json
// @Produce      json
// @Success      200 {array} models.Company "正常なレスポンス"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /company/list [get]
func (h *CompanyHandler) GetCompanies(c *fiber.Ctx) error {
	companies := h.CompanyService.GetCompanies()
	return c.JSON(companies)
}
