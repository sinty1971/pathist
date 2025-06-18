package handlers

import (
	"penguin-backend/internal/models"
	"penguin-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

// KoujiHandler 工事関連のHTTPリクエストを処理するハンドラー
type KoujiHandler struct {
	fileSystemService *services.FileSystemService
	koujiService      *services.KoujiService
}

// NewKoujiHandler 新しいKoujiHandlerインスタンスを作成します
func NewKoujiHandler(fsService *services.FileSystemService, koujiService *services.KoujiService) *KoujiHandler {
	return &KoujiHandler{
		fileSystemService: fsService,
		koujiService:      koujiService,
	}
}

// GetKoujiEntries godoc
// @Summary      工事プロジェクト一覧の取得
// @Description  指定されたパスから工事プロジェクトフォルダーの一覧を取得します。
// @Description  各工事プロジェクトには会社名、現場名、工事開始日などの詳細情報が含まれます。
// @Tags         工事管理
// @Accept       json
// @Produce      json
// @Param        path query string false "工事フォルダーのパス" default(~/penguin/豊田築炉/2-工事)
// @Success      200 {object} models.KoujiEntriesResponse "工事プロジェクト一覧"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /kouji-entries [get]
func (h *KoujiHandler) GetKoujiEntries(c *fiber.Ctx) error {
	// KoujiServiceを使用して工事エントリを取得
	koujiEntries := h.koujiService.GetKoujiEntries()

	totalSize := int64(0)
	for _, kouji := range koujiEntries {
		totalSize += kouji.FileEntry.Size
	}

	return c.JSON(models.KoujiEntriesResponse{
		KoujiEntries: koujiEntries,
		Count:        len(koujiEntries),
		TotalSize:    totalSize,
	})
}

// SaveKoujiEntries godoc
// @Summary      工事情報のYAML保存
// @Description  工事プロジェクト情報をYAMLファイルに保存します
// @Tags         工事管理
// @Accept       json
// @Produce      json
// @Param        body body []models.KoujiEntry false "工事データ（オプション）"
// @Success      200 {object} map[string]string "成功メッセージ"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /kouji-entries/save [post]
func (h *KoujiHandler) SaveKoujiEntries(c *fiber.Ctx) error {
	// リクエストボディから編集された工事データを取得
	var entries []models.KoujiEntry
	if err := c.BodyParser(&entries); err != nil {
		// ボディが空の場合は、現在のデータを取得
		entries = h.koujiService.GetKoujiEntries()
	}
	
	// KoujiServiceを使用して工事プロジェクトを保存
	err := h.koujiService.SaveKoujiEntries(entries)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to save kouji entries",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":     "工事データを正常に保存しました",
		"output_path": h.koujiService.Database,
		"count":       len(entries),
	})
}
