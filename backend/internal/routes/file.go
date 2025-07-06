package routes

import (
	"penguin-backend/internal/handlers"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cache"
)

// SetupFileRoutes はファイル関連のルートを設定します
func SetupFileRoutes(api fiber.Router, handler *handlers.FileHandler) {
	files := api.Group("/files")

	// ファイル情報取得 - 長めのキャッシュ（ファイルシステムの変更は頻繁でない）
	files.Get("/", 
		cache.New(cache.Config{
			Expiration:   2 * time.Minute,     // 2分間キャッシュ
			CacheHeader:  "X-File-Cache",      // カスタムヘッダー
			CacheControl: true,
			KeyGenerator: func(c fiber.Ctx) string {
				// パスクエリパラメータを含めてキャッシュキーを生成
				path := c.Query("path", "")
				return "files:" + path
			},
		}),
		handler.GetFileInfos,
	)
}