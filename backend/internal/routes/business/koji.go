package business

import (
	"penguin-backend/internal/handlers/business"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
)

// setupKojiRoutes は工事関連のルートを設定します
func setupKojiRoutes(api fiber.Router, businessHandler *business.BusinessHandler) {
	// 工事一覧取得 - 中程度のキャッシュ
	api.Get("/kojies",
		cache.New(cache.Config{
			Expiration:   5 * time.Second, // 5秒間キャッシュ
			CacheHeader:  "X-Koji-Cache",
			CacheControl: true,
			KeyGenerator: func(c *fiber.Ctx) string {
				// filterクエリパラメータを含めてキャッシュキーを生成
				filter := c.Query("filter", "")
				return "kojies:list:" + filter
			},
		}),
		businessHandler.GetKojies,
	)

	// 工事取得 - キャッシュあり
	api.Get("/kojies/:path",
		cache.New(cache.Config{
			Expiration:   5 * time.Second, // 5秒間キャッシュ
			CacheHeader:  "X-Koji-Detail-Cache",
			CacheControl: true,
			KeyGenerator: func(c *fiber.Ctx) string {
				path := c.Params("path")
				return "kojies:detail:" + path
			},
		}),
		businessHandler.GetKojiByPath,
	)

	// 工事更新 - キャッシュなし（更新系）
	api.Put("/kojies", businessHandler.UpdateKoji)

	// 標準ファイル名変更 - キャッシュなし（更新系）
	api.Put("/kojies/standard-files", businessHandler.RenameKojiStandardFiles)
}
