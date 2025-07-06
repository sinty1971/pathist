package handlers

import (
	"fmt"
	"net/url"
	"penguin-backend/internal/models"
	"penguin-backend/internal/services"

	"github.com/gofiber/fiber/v3"
)

// KojiHandler 工事関連のHTTPリクエストを処理するハンドラー
type KojiHandler struct {
	kojiService *services.KojiService
}

// NewKojiHandler 新しいKojiHandlerインスタンスを作成します
func NewKojiHandler(kojiService *services.KojiService) *KojiHandler {
	return &KojiHandler{
		kojiService: kojiService,
	}
}

// GetKoji godoc
// @Summary      Kojiの取得
// @Description  Kojiを取得します。
// @Tags         工事管理
// @Produce      json
// @Param        path path string true "Kojiフォルダーのファイル名"
// @Success      200 {object} models.Koji "Koji"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /kojies/{path} [get]
func (h *KojiHandler) GetByPath(c fiber.Ctx) error {
	// パスパラメータを取得
	path := c.Params("path")

	// URLデコード
	decodedPath, err := url.QueryUnescape(path)
	if err != nil {
		fmt.Printf("URL decode error: %v\n", err)
		decodedPath = path // デコードに失敗した場合は元のパスを使用
	}

	// KojiServiceを使用してKojiを取得
	koji, err := h.kojiService.GetKoji(decodedPath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "工事の取得に失敗しました",
			"message": err.Error(),
		})
	}

	return c.JSON(koji)
}

// GetKojies godoc
// @Summary      工事一覧の取得
// @Description  工事一覧を取得します。filter=recentで最近の工事のみを取得できます。
// @Tags         工事管理
// @Produce      json
// @Param        filter query string false "フィルター (recent: 最近の工事のみ)"
// @Success      200 {object} []models.Koji "工事一覧"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /kojies [get]
func (h *KojiHandler) GetKojies(c fiber.Ctx) error {
	// クエリパラメータでフィルターを取得
	filter := c.Query("filter")

	// filter=recentの場合は最近の工事のみを取得
	if filter == "recent" {
		kojies := h.kojiService.GetRecentKojies()
		return c.JSON(kojies)
	} else {
		fmt.Printf("filter: %s\n", filter)
		// フィルターが指定されていない場合は最近の工事を返す（互換性のため）
		kojies := h.kojiService.GetRecentKojies()
		return c.JSON(kojies)

	}
}

// Update godoc
// @Summary      更新対象フォルダー（Koji.FileInfo）のファイル情報をKojiのメンバ変数で更新
// @Description  更新対象フォルダー（Koji.FileInfo）のファイル情報をKojiのメンバ変数で更新します。
// @Tags         工事管理
// @Accept       json
// @Produce      json
// @Param        body body models.Koji true "工事データ"
// @Success      200 {object} models.Koji "更新後の工事データ"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /kojies [put]
func (h *KojiHandler) Update(c fiber.Ctx) error {
	// リクエストボディから編集された工事を取得
	var koji models.Koji
	if err := c.Bind().Body(&koji); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "リクエストボディの解析に失敗しました",
			"message": err.Error(),
		})
	}

	// kojiのファイル情報を更新（.detail.yamlも更新）
	err := h.kojiService.Update(&koji)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "工事データの保存に失敗しました",
			"message": err.Error(),
		})
	}

	return c.JSON(koji)
}

// RenameManagedFileRequest リクエスト用の構造体
type RenameManagedFileRequest struct {
	Koji     models.Koji `json:"koji"`
	Currents []string    `json:"currents"`
}

// RenameMangedFile godoc
// @Summary      ManagedFileの名前を変更
// @Description  ManagedFileの名前を変更します。
// @Tags         工事管理
// @Accept       json
// @Produce      json
// @Param        body body RenameManagedFileRequest true "工事データと管理ファイル"
// @Success      200 {object} []string "更新後のファイル名リスト"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /kojies/managed-files [put]
func (h *KojiHandler) RenameManagedFile(c fiber.Ctx) error {
	// リクエストボディから編集された工事を取得
	var request RenameManagedFileRequest
	if err := c.Bind().Body(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "リクエストボディの解析に失敗しました",
			"message": err.Error(),
		})
	}

	renamedFiles := h.kojiService.RenameManagedFile(request.Koji, request.Currents)

	return c.JSON(renamedFiles)
}