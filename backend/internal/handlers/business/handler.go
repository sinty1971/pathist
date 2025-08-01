package business

import (
	"penguin-backend/internal/services"

	"github.com/gofiber/fiber/v3"
)

// BusinessHandler はビジネス関連の全てのハンドラーメソッドを提供します
type BusinessHandler struct {
	businessService *services.BusinessService
}

// NewBusinessHandler は新しいBusinessHandlerを作成します
func NewBusinessHandler(businessService *services.BusinessService) *BusinessHandler {
	return &BusinessHandler{
		businessService: businessService,
	}
}

// GetAttributeFilename godoc
// @Summary      属性ファイル名の取得
// @Description  属性ファイル名を取得します
// @Tags         ビジネス管理
// @Produce      json
// @Success      200 {object} map[string]string "属性ファイル名"
// @Router       /business/attribute-filename [get]
func (bh *BusinessHandler) GetAttributeFilename(c fiber.Ctx) error {
	attributeFilename := bh.businessService.DatabaseFilename
	return c.JSON(fiber.Map{"attributeFilename": attributeFilename})
}
