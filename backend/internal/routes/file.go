package routes

import (
	"penguin-backend/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

// SetupFileRoutes はファイル関連のルートを設定します
func SetupFileRoutes(api fiber.Router, handler *handlers.FileHandler) {
	file := api.Group("/file")

	// ファイル情報の取得
	// @Summary ファイル情報一覧取得
	// @Description 指定されたパスのファイル情報一覧を取得
	// @Tags file
	// @Accept json
	// @Produce json
	// @Param path query string false "ファイルパス"
	// @Success 200 {array} models.FileInfo
	// @Router /file/fileinfos [get]
	file.Get("/fileinfos", handler.GetFileInfos)
}