package handlers

import (
	"penguin-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

type FileSystemHandler struct {
	FileSystemService *services.FileSystemService
}

func NewFileSystemHandler(fsService *services.FileSystemService) *FileSystemHandler {
	return &FileSystemHandler{
		FileSystemService: fsService,
	}
}

// GetFileEntries godoc
// @Summary      ファイルエントリ一覧の取得
// @Description  指定されたパスからファイルとフォルダーの一覧を取得します
// @Tags         ファイル管理
// @Accept       json
// @Produce      json
// @Param        path query string false "取得するディレクトリのパス" default(~/penguin)
// @Success      200 {object} models.FileEntriesListResponse "正常なレスポンス"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /file-entries [get]
func (h *FileSystemHandler) GetFileEntries(c *fiber.Ctx) error {
	fsPath := c.Query("path", "~/penguin")

	fileEntries, err := h.FileSystemService.GetFileEntries(fsPath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to read directory",
			"message": err.Error(),
		})
	}

	return c.JSON(fileEntries)
}
