package handlers

import (
	"penguin-backend/internal/models"
	"penguin-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

// ProjectHandler 工事関連のHTTPリクエストを処理するハンドラー
type ProjectHandler struct {
	projectService *services.ProjectService
}

// NewProjectHandler 新しいProjectHandlerインスタンスを作成します
func NewProjectHandler(projectService *services.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
	}
}

// GetRecentProjects godoc
// @Summary      最近の工事一覧の取得
// @Description  最近の工事一覧を取得します。
// @Tags         プロジェクト管理
// @Accept       json
// @Produce      json
// @Success      200 {object} []models.Project "最近のプロジェクト一覧"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /project/recent [get]
func (h *ProjectHandler) GetRecentProjects(c *fiber.Ctx) error {
	// ProjectServiceを使用して最近のプロジェクト一覧を取得
	projects := h.projectService.GetRecentProjects()

	return c.JSON(projects)
}

// Save godoc
// @Summary      指定されたプロジェクト情報のYAML保存
// @Description  models.ProjectをYAMLファイルに保存します
// @Tags         プロジェクト管理
// @Accept       json
// @Produce      json
// @Param        body body models.Project true "工事データ"
// @Success      200 {object} models.Project "更新後の工事データ"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /project/update [post]
func (h *ProjectHandler) Update(c *fiber.Ctx) error {
	// リクエストボディから編集されたプロジェクトを取得
	var project models.Project
	if err := c.BodyParser(&project); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "リクエストボディの解析に失敗しました",
			"message": err.Error(),
		})
	}

	// projectのファイル情報を更新（.detail.yamlも更新）
	err := h.projectService.UpdateProjectFileInfo(&project)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "工事データの保存に失敗しました",
			"message": err.Error(),
		})
	}

	return c.JSON(project)
}
