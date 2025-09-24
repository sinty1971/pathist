package business

import (
	"penguin-backend/internal/models"
	"penguin-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

// BusinessHandler はビジネス関連の全てのハンドラーメソッドを提供します。
type BusinessHandler struct {
	businessService *services.BusinessService
}

// NewBusinessHandler は新しいBusinessHandlerを作成します。
func NewBusinessHandler(businessService *services.BusinessService) *BusinessHandler {
	return &BusinessHandler{businessService: businessService}
}

// GetBusinessBasePath はビジネスベースパスを返します。
func (h *BusinessHandler) GetBusinessBasePath(c *fiber.Ctx) error {
	basePath := h.businessService.FolderPath
	return c.JSON(fiber.Map{"basePath": basePath})
}

// GetAttributeFilename は属性ファイル名を返します。
func (h *BusinessHandler) GetAttributeFilename(c *fiber.Ctx) error {
	attributeFilename := h.businessService.DatabaseFilename
	return c.JSON(fiber.Map{"attributeFilename": attributeFilename})
}

// GetFiles はファイル一覧を取得します。
func (h *BusinessHandler) GetFiles(c *fiber.Ctx) error {
	fsPath := c.Query("path", "")

	if h.businessService.FileService == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "FileService not configured",
			"message": "FileService is unavailable",
		})
	}

	fileInfos, err := h.businessService.FileService.GetFileInfos(fsPath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to read directory",
			"message": err.Error(),
		})
	}

	return c.JSON(fileInfos)
}

// Service は内部で利用するサービス参照を返します。
func (h *BusinessHandler) Service() *services.BusinessService {
	return h.businessService
}

// PostProcessFiles はファイル一覧を加工するための拡張ポイントです。
func (h *BusinessHandler) PostProcessFiles(files []models.FileInfo) []models.FileInfo {
	return files
}
