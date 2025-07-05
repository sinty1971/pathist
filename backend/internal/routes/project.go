package routes

import (
	"penguin-backend/internal/handlers"

	"github.com/gofiber/fiber/v3"
)

// SetupProjectRoutes はプロジェクト関連のルートを設定します
func SetupProjectRoutes(api fiber.Router, handler *handlers.ProjectHandler) {
	project := api.Group("/project")

	// プロジェクト取得
	project.Get("/get/:path", handler.GetByPath)

	// プロジェクト一覧取得
	project.Get("/recent", handler.GetRecent)

	// プロジェクト更新
	project.Post("/update", handler.Update)

	// 管理ファイル名変更
	project.Post("/rename-managed-file", handler.RenameManagedFile)
}