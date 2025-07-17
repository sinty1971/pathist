package business

import (
	"penguin-backend/internal/handlers/business"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cache"
)

// setupFileRoutes はファイル関連のルートを設定します
func setupFileRoutes(api fiber.Router, businessHandler *business.BusinessHandler) {
	// ファイル情報取得 - キャッシュあり
	api.Get("/files",
		cache.New(cache.Config{
			Expiration:   3 * time.Second, // 3秒間キャッシュ
			CacheHeader:  "X-File-Cache",  // カスタムヘッダー
			CacheControl: true,
			KeyGenerator: func(c fiber.Ctx) string {
				// パスクエリパラメータを含めてキャッシュキーを生成
				path := c.Query("path", "")
				return "files:" + path
			},
		}),
		businessHandler.GetFiles,
	)
}