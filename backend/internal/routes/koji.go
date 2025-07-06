package routes

import (
	"penguin-backend/internal/handlers"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cache"
)

// SetupKojiRoutes は工事関連のルートを設定します
func SetupKojiRoutes(api fiber.Router, handler *handlers.KojiHandler) {
	kojies := api.Group("/kojies")

	// 工事一覧取得 - 中程度のキャッシュ
	kojies.Get("/", 
		cache.New(cache.Config{
			Expiration:   1 * time.Minute,     // 1分間キャッシュ
			CacheHeader:  "X-Koji-Cache",
			CacheControl: true,
			KeyGenerator: func(c fiber.Ctx) string {
				// filterクエリパラメータを含めてキャッシュキーを生成
				filter := c.Query("filter", "")
				return "kojies:list:" + filter
			},
		}),
		handler.GetKojies,
	)

	// 工事取得 - 長めのキャッシュ（個別工事は変更頻度が低い）
	kojies.Get("/:path", 
		cache.New(cache.Config{
			Expiration:   3 * time.Minute,     // 3分間キャッシュ
			CacheHeader:  "X-Koji-Detail-Cache",
			CacheControl: true,
			KeyGenerator: func(c fiber.Ctx) string {
				path := c.Params("path")
				return "kojies:detail:" + path
			},
		}),
		handler.GetByPath,
	)

	// 工事更新 - キャッシュなし（更新系）
	kojies.Put("/", handler.Update)

	// 管理ファイル名変更 - キャッシュなし（更新系）
	kojies.Put("/managed-files", handler.RenameManagedFile)
}