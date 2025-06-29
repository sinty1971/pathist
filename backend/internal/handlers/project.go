package handlers

import (
	"fmt"
	"net/url"
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

// GetProject godoc
// @Summary      Projectの取得
// @Description  Projectを取得します。
// @Tags         プロジェクト管理
// @Accept       json
// @Produce      json
// @Param        path path string true "Projectフォルダーのファイル名"
// @Success      200 {object} models.Project "Project"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /project/get/{path} [get]
func (h *ProjectHandler) GetProject(c *fiber.Ctx) error {
	// パスパラメータを取得
	path := c.Params("path")
	fmt.Printf("Handler received raw path: '%s'\n", path)

	// URLデコード
	decodedPath, err := url.QueryUnescape(path)
	if err != nil {
		fmt.Printf("URL decode error: %v\n", err)
		decodedPath = path // デコードに失敗した場合は元のパスを使用
	}
	fmt.Printf("Handler decoded path: '%s'\n", decodedPath)

	// ProjectServiceを使用してProjectを取得
	project, err := h.projectService.GetProject(decodedPath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "工事の取得に失敗しました",
			"message": err.Error(),
		})
	}

	return c.JSON(project)
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

// Update godoc
// @Summary      Project.FileInfoで指定されたファイル情報をProjectのメンバ変数で更新
// @Description  Project.FileInfoで指定されたファイル情報をProjectのメンバ変数で更新します。
// @Tags         工事 更新
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

// RenameManagedFileRequest リクエスト用の構造体
type RenameManagedFileRequest struct {
	Project  models.Project `json:"project"`
	Currents []string       `json:"currents"`
}

// Rename MangedFile godoc
// @Summary      ManagedFileの名前を変更
// @Description  ManagedFileの名前を変更します。
// @Tags         プロジェクト管理
// @Accept       json
// @Produce      json
// @Param        body body RenameManagedFileRequest true "工事データと管理ファイル"
// @Success      200 {object} []string "更新後のファイル名リスト"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /project/rename-managed-file [post]
func (h *ProjectHandler) RenameManagedFile(c *fiber.Ctx) error {
	// リクエストボディから編集されたプロジェクトを取得
	var request RenameManagedFileRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "リクエストボディの解析に失敗しました",
			"message": err.Error(),
		})
	}

	renamedFiles := h.projectService.RenameManagedFile(request.Project, request.Currents)

	return c.JSON(renamedFiles)
}
