package business

import (
	"penguin-backend/internal/models"

	"github.com/gofiber/fiber/v3"
)

// GetCompanyByID godoc
// @Summary      会社詳細の取得
// @Description  指定されたIDの会社詳細を取得します
// @Tags         会社管理
// @Produce      json
// @Param        id path string true "会社ID"
// @Success      200 {object} models.Company "正常なレスポンス"
// @Failure      404 {object} map[string]string "会社が見つからない"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /business/companies/{id} [get]
func (bh *BusinessHandler) GetCompanyByID(c fiber.Ctx) error {
	company, err := bh.businessService.CompanyService.GetCompanyByID(c.Params("id"))
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
// @Produce      json
// @Success      200 {array} models.Company "正常なレスポンス"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /business/companies [get]
func (bh *BusinessHandler) GetCompanies(c fiber.Ctx) error {
	companies := bh.businessService.CompanyService.GetCompanies()
	return c.JSON(companies)
}

// UpdateCompany godoc
// @Summary      会社情報の更新
// @Description  会社情報を更新し、必要に応じてフォルダー名を変更します
// @Tags         会社管理
// @Accept       json
// @Produce      json
// @Param        body body models.Company true "会社更新リクエスト"
// @Success      200 {object} models.Company "更新後の会社データ"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /business/companies [put]
func (bh *BusinessHandler) UpdateCompany(c fiber.Ctx) error {
	// リクエストボディから編集された会社を取得
	var company models.Company
	if err := c.Bind().Body(&company); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "リクエストボディの解析に失敗しました",
			"message": err.Error(),
		})
	}

	// 会社の情報を更新
	err := bh.businessService.CompanyService.Update(&company)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "会社データの保存に失敗しました",
			"message": err.Error(),
		})
	}

	return c.JSON(company)
}

// GetCategories godoc
// @Summary      カテゴリー一覧の取得
// @Description  カテゴリー一覧を取得します
// @Tags         会社管理
// @Produce      json
// @Success      200 {array} models.CompanyCategoryInfo "正常なレスポンス"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /business/companies/categories [get]
func (bh *BusinessHandler) GetCategories(c fiber.Ctx) error {
	return c.JSON(bh.businessService.CompanyService.Categories())
}
