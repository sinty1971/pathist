package routes

import (
	"penguin-backend/internal/handlers"

	"github.com/gofiber/fiber/v3"
)

// SetupFileRoutes はファイル関連のルートを設定します
func SetupFileRoutes(api fiber.Router, handler *handlers.FileHandler) {
	file := api.Group("/file")

	// ファイル情報取得
	file.Get("/fileinfos", handler.GetFileInfos)
}