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

// GetEntries godoc
// @Summary      最近の工事プロジェクト一覧の取得
// @Description  最近の工事プロジェクト一覧を取得します。
// @Tags         プロジェクト管理
// @Accept       json
// @Produce      json
// @Success      200 {object} []models.Project "最近のプロジェクト一覧"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /project/get/recent [get]
func (h *ProjectHandler) GetRecentProjects(c *fiber.Ctx) error {
	// ProjectServiceを使用して最近のプロジェクト一覧を取得
	projects := h.projectService.GetRecentProjects()

	return c.JSON(projects)
}

// Save godoc
// @Summary      指定されたプロジェクト情報のYAML保存
// @Description  []models.ProjectをYAMLファイルに保存します
// @Tags         プロジェクト管理
// @Accept       json
// @Produce      json
// @Param        body body []models.KoujiEntry false "工事データ（オプション）"
// @Success      200 {object} map[string]string "成功メッセージ"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /kouji/save [post]
func (h *ProjectHandler) Save(c *fiber.Ctx) error {
	// リクエストボディから編集されたプロジェクトを取得
	var project models.Project
	if err := c.BodyParser(&project); err != nil {
		// ボディが空の場合は、現在のデータを取得
		project = h.projectService.GetProject(project.ID)
	}

	// KoujiServiceを使用して工事プロジェクトを保存
	err := h.koujiService.SaveToYAML(entries)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "工事データの保存に失敗しました",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":     "工事データを正常に保存しました",
		"output_path": h.koujiService.Database,
		"count":       len(entries),
	})
}

// SearchManageFiles は工事フォルダー内で管理されているファイルを検索する
// 戻り値はそれらのファイルが規則に沿っているかを返す
// func (h *KoujiHandler) SearchManageFiles(koujiFolderPath string) error {
// 	// 工事フォルダーのパスを取得
// 	koujiFolderPath, err := h.fileSystemService.BuildPath(koujiFolderPath)
// 	if err != nil {
// 		return err
// 	}

// 	// データベースから工事一覧を取得
// 	entryies, err := h.fileSystemService.GetFileEntries(koujiFolderPath)
// 	if err != nil {
// 		return err
// 	}

// 	// 工事フォルダー内で管理されているファイルを検索する

// }
