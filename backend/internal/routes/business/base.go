package business

import (
	"penguin-backend/internal/handlers/business"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
)

// setupBaseRoutes はベース関連のルートを設定します
func setupBaseRoutes(api fiber.Router, businessHandler *business.BusinessHandler) {
	// ビジネスベースパス取得 - キャッシュあり
	api.Get("/base-path",
		cache.New(cache.Config{
			Expiration:   1 * time.Minute, // 1分間キャッシュ
			CacheHeader:  "X-Business-Base-Path-Cache",
			CacheControl: true,
		}),
		businessHandler.GetBusinessBasePath,
	)

	// 属性ファイル名取得 - キャッシュあり
	api.Get("/attribute-filename",
		cache.New(cache.Config{
			Expiration:   1 * time.Minute, // 1分間キャッシュ
			CacheHeader:  "X-Attribute-Filename-Cache",
			CacheControl: true,
		}),
		businessHandler.GetAttributeFilename,
	)
}
