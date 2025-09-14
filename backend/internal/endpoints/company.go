package endpoints

import (
	"penguin-backend/internal/models"
	"penguin-backend/internal/services"

	"github.com/gofiber/fiber/v3"
)

type CompanyEndPoint struct {
	services.CompanyService
}

// RegisterRoutes は会社管理のルートを登録する
func (ce *CompanyEndPoint) RegisterRoutes(router fiber.Router) {
	companies := router.Group("/companies")
	companies.Get("/", ce.HandleGetCompanies)
	companies.Get("/categories", ce.HandleGetCategories)
	companies.Get("/:id", ce.HandleGetCompanyByID)
	companies.Put("/", ce.HandleUpdateCompany)
}

// HandleGetCompanies godoc
// @Summary      会社一覧の取得
// @Description  会社フォルダーの一覧を取得します
// @Tags         会社管理
// @Produce      json
// @Success      200 {array} models.Company "正常なレスポンス"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /companies [get]
func (ce *CompanyEndPoint) HandleGetCompanies(c fiber.Ctx) error {
	companies := ce.GetCompanies()
	return c.JSON(companies)
}

// HandleGetCompanyByID godoc
// @Summary      会社詳細の取得
// @Description  指定されたIDの会社詳細を取得します
// @Tags         会社管理
// @Produce      json
// @Param        id path string true "会社ID"
// @Success      200 {object} models.Company "正常なレスポンス"
// @Failure      404 {object} map[string]string "会社が見つからない"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /companies/{id} [get]
func (ce *CompanyEndPoint) HandleGetCompanyByID(c fiber.Ctx) error {
	company, err := ce.GetCompanyByID(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Company not found",
			"message": err.Error(),
		})
	}
	return c.JSON(company)
}

// HandleUpdateCompany godoc
// @Summary      会社情報の更新
// @Description  会社情報を更新し、必要に応じてフォルダー名を変更します
// @Tags         会社管理
// @Accept       json
// @Produce      json
// @Param        body body models.Company true "会社更新リクエスト"
// @Success      200 {object} models.Company "更新後の会社データ"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /companies [put]
func (ce *CompanyEndPoint) HandleUpdateCompany(c fiber.Ctx) error {
	var updateCompany models.Company
	if err := c.Bind().JSON(&updateCompany); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err := ce.Update(&updateCompany)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(updateCompany)
}

// HandleGetCategories godoc
// @Summary      カテゴリー一覧の取得
// @Description  カテゴリー一覧を取得します
// @Tags         会社管理
// @Produce      json
// @Success      200 {array} models.CompanyCategoryInfo "正常なレスポンス"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /companies/categories [get]
func (ce *CompanyEndPoint) HandleGetCategories(c fiber.Ctx) error {
	categories := ce.GetCategories()
	return c.JSON(categories)
}
