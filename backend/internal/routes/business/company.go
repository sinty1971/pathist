package business

import (
	"penguin-backend/internal/handlers/business"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cache"
)

// setupCompanyRoutes は会社関連のルートを設定します
func setupCompanyRoutes(api fiber.Router, businessHandler *business.BusinessHandler) {
	// カテゴリー一覧取得 - キャッシュあり
	api.Get("/companies/categories",
		cache.New(cache.Config{
			Expiration:   5 * time.Minute, // 5分間キャッシュ
			CacheHeader:  "X-Categories-Cache",
			CacheControl: true,
			KeyGenerator: func(c fiber.Ctx) string {
				return "categories:list"
			},
		}),
		businessHandler.GetCategories,
	)

	// 会社一覧取得 - キャッシュあり
	api.Get("/companies",
		cache.New(cache.Config{
			Expiration:   10 * time.Second, // 10秒間キャッシュ
			CacheHeader:  "X-Company-Cache",
			CacheControl: true,
			KeyGenerator: func(c fiber.Ctx) string {
				return "companies:list"
			},
		}),
		businessHandler.GetCompanies,
	)

	// 会社詳細取得 - キャッシュあり
	api.Get("/companies/:id",
		cache.New(cache.Config{
			Expiration:   10 * time.Second, // 10秒間キャッシュ
			CacheHeader:  "X-Company-Detail-Cache",
			CacheControl: true,
			KeyGenerator: func(c fiber.Ctx) string {
				id := c.Params("id")
				return "companies:detail:" + id
			},
		}),
		businessHandler.GetCompanyByID,
	)

	// 会社情報更新 - キャッシュなし（更新系）
	api.Put("/companies", businessHandler.UpdateCompany)
}