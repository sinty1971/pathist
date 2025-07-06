package handlers

import (
	"penguin-backend/internal/services"

	"github.com/gofiber/fiber/v3"
)

type FileHandler struct {
	FileService *services.FileService
}

func NewFileHandler(fsService *services.FileService) *FileHandler {
	return &FileHandler{
		FileService: fsService,
	}
}

// GetFileInfos godoc
// @Summary      ファイルエントリ一覧の取得
// @Description  指定されたパスからファイルとフォルダーの一覧を取得します
// @Tags         ファイル管理
// @Produce      json
// @Param        path query string false "取得するディレクトリのパス" default(2 工事)
// @Success      200 {array} models.FileInfo "正常なレスポンス"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /files [get]
func (h *FileHandler) GetFileInfos(c fiber.Ctx) error {
	fsPath := c.Query("path", "")

	fileInfos, err := h.FileService.GetFileInfos(fsPath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to read directory",
			"message": err.Error(),
		})
	}

	return c.JSON(fileInfos)
}
