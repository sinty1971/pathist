package routes

import (
	"penguin-backend/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

// SetupProjectRoutes はプロジェクト関連のルートを設定します
func SetupProjectRoutes(api fiber.Router, handler *handlers.ProjectHandler) {
	project := api.Group("/project")

	// 単一プロジェクトの取得
	// @Summary プロジェクト取得
	// @Description 指定されたパスのプロジェクトを取得
	// @Tags project
	// @Accept json
	// @Produce json
	// @Param path path string true "プロジェクトパス"
	// @Success 200 {object} models.Project
	// @Router /project/get/{path} [get]
	project.Get("/get/:path", handler.GetProject)

	// 最近のプロジェクト一覧
	// @Summary 最近のプロジェクト一覧取得
	// @Description 最近の工事プロジェクトの一覧を取得
	// @Tags project
	// @Accept json
	// @Produce json
	// @Success 200 {array} models.Project
	// @Router /project/recent [get]
	project.Get("/recent", handler.GetRecentProjects)

	// プロジェクトの更新
	// @Summary プロジェクト更新
	// @Description 工事プロジェクト情報を更新
	// @Tags project
	// @Accept json
	// @Produce json
	// @Param project body models.Project true "プロジェクト情報"
	// @Success 200 {object} models.Project
	// @Router /project/update [post]
	project.Post("/update", handler.Update)

	// 管理ファイルの名前変更
	// @Summary 管理ファイル名変更
	// @Description プロジェクトの管理ファイル名を変更
	// @Tags project
	// @Accept json
	// @Produce json
	// @Param body body object true "リクエストボディ"
	// @Success 200 {array} string
	// @Router /project/rename-managed-file [post]
	project.Post("/rename-managed-file", handler.RenameManagedFile)
}