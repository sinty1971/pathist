package routes

import (
	"penguin-backend/internal/handlers"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cache"
)

// SetupBusinessRoutes はビジネス関連のルートを設定します
func SetupBusinessRoutes(api fiber.Router, handler *handlers.BusinessHandler) {

	// ビジネスベースパス取得 - キャッシュあり
	api.Get("/base-path",
		cache.New(cache.Config{
			Expiration:   1 * time.Minute, // 1分間キャッシュ
			CacheHeader:  "X-Business-Base-Path-Cache",
			CacheControl: true,
		}),
		handler.GetBusinessBasePath,
	)

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
		handler.GetFiles,
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
		handler.GetCompanies,
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
		handler.GetCompanyByID,
	)

	// 会社情報更新 - キャッシュなし（更新系）
	api.Put("/companies", handler.UpdateCompany)

	// 工事一覧取得 - 中程度のキャッシュ
	api.Get("/kojies",
		cache.New(cache.Config{
			Expiration:   5 * time.Second, // 5秒間キャッシュ
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

	// 工事取得 - キャッシュあり
	api.Get("/kojies/:path",
		cache.New(cache.Config{
			Expiration:   5 * time.Second, // 5秒間キャッシュ
			CacheHeader:  "X-Koji-Detail-Cache",
			CacheControl: true,
			KeyGenerator: func(c fiber.Ctx) string {
				path := c.Params("path")
				return "kojies:detail:" + path
			},
		}),
		handler.GetKojiByPath,
	)

	// 工事更新 - キャッシュなし（更新系）
	api.Put("/kojies", handler.UpdateKoji)

	// 管理ファイル名変更 - キャッシュなし（更新系）
	api.Put("/kojies/managed-files", handler.RenameKojiManagedFile)

}
