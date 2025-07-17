package business

import (
	"github.com/gofiber/fiber/v3"
)

// GetFiles godoc
// @Summary      ファイルエントリ一覧の取得
// @Description  指定されたパスからファイルとフォルダーの一覧を取得します
// @Tags         ファイル管理
// @Produce      json
// @Param        path query string false "取得するディレクトリのパス" default(2 工事)
// @Success      200 {array} models.FileInfo "正常なレスポンス"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /business/files [get]
func (bh *BusinessHandler) GetFiles(c fiber.Ctx) error {
	fsPath := c.Query("path", "")

	fileInfos, err := bh.businessService.FileService.GetFileInfos(fsPath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to read directory",
			"message": err.Error(),
		})
	}

	return c.JSON(fileInfos)
}