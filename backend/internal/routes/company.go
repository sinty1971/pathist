package routes

import (
	"penguin-backend/internal/handlers"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cache"
)

// SetupCompanyRoutes 会社関連のルートを設定
func SetupCompanyRoutes(api fiber.Router, handler *handlers.CompanyHandler) {
	companies := api.Group("/companies")

	// 会社一覧取得 - 長めのキャッシュ（会社情報は変更頻度が低い）
	companies.Get("/", 
		cache.New(cache.Config{
			Expiration:   5 * time.Minute,     // 5分間キャッシュ
			CacheHeader:  "X-Company-Cache",
			CacheControl: true,
			KeyGenerator: func(c fiber.Ctx) string {
				return "companies:list"
			},
		}),
		handler.GetCompanies,
	)

	// 会社詳細取得 - 最も長いキャッシュ（個別会社情報は最も変更頻度が低い）
	companies.Get("/:id", 
		cache.New(cache.Config{
			Expiration:   10 * time.Minute,    // 10分間キャッシュ
			CacheHeader:  "X-Company-Detail-Cache",
			CacheControl: true,
			KeyGenerator: func(c fiber.Ctx) string {
				id := c.Params("id")
				return "companies:detail:" + id
			},
		}),
		handler.GetByID,
	)

	// 会社情報更新 - キャッシュなし（更新系）
	companies.Put("/", handler.Update)
}
