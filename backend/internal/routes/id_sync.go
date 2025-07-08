package routes

import (
	"penguin-backend/internal/handlers"

	"github.com/gofiber/fiber/v3"
)

// SetupIDSyncRoutes はID同期関連のルートを設定します
func SetupIDSyncRoutes(router fiber.Router, handler *handlers.IDSyncHandler) {
	// 工事IDの生成
	router.Post("/generate-koji", handler.GenerateKojiID)

	// パスIDの生成
	router.Post("/generate-path", handler.GeneratePathID)

	// IDの検証
	router.Post("/validate", handler.ValidateID)

	// 一括ID変換
	router.Post("/bulk-convert", handler.BulkConvert)

	// IDマッピングテーブルの取得
	router.Get("/mapping", handler.GetIDMapping)

	// 同期設定の取得
	router.Get("/config", handler.SyncConfiguration)
}