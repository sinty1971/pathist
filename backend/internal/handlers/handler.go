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

// GetBusinessBasePath godoc
// @Summary      ビジネスベースパスの取得
// @Description  ビジネスベースパスを取得します
// @Tags         ビジネス管理
// @Produce      json
// @Success      200 {object} map[string]string "ビジネスベースパス"
// @Router       /business/base-path [get]
func (h *BusinessHandler) GetBusinessBasePath(c fiber.Ctx) error {
	basePath := h.businessService.FolderPath
	return c.JSON(fiber.Map{"basePath": basePath})
}

// GetAttributeFilename godoc
// @Summary      属性ファイル名の取得
// @Description  属性ファイル名を取得します
// @Tags         ビジネス管理
// @Produce      json
// @Success      200 {object} map[string]string "属性ファイル名"
// @Router       /business/attribute-filename [get]
func (h *BusinessHandler) GetAttributeFilename(c fiber.Ctx) error {
	attributeFilename := h.businessService.DatabaseFilename
	return c.JSON(fiber.Map{"attributeFilename": attributeFilename})
}

// GetFiles godoc
// @Summary      ファイルエントリ一覧の取得
// @Description  指定されたパスからファイルとフォルダーの一覧を取得します
// @Tags         ファイル管理
// @Produce      json
// @Param        path query string false "取得するディレクトリのパス"
// @Success      200 {array} models.FileInfo "正常なレスポンス"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /business/files [get]
func (h *BusinessHandler) GetFiles(c fiber.Ctx) error {
	fsPath := c.Query("path", "")

	fileInfos, err := h.businessService.RootService.FileService.GetFileInfos(fsPath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to read directory",
			"message": err.Error(),
		})
	}

	return c.JSON(fileInfos)
}
