package routes

import (
	"net/http"
	"strings"

	_ "penguin-backend/docs" // swaggo generated docs
	"penguin-backend/internal/handlers"
	"penguin-backend/internal/services"

	"github.com/gofiber/fiber/v3"
)

// SetupRoutes はすべてのルートを設定します
func SetupRoutes(app *fiber.App, container *services.ServiceContainer) {
	// Swagger UI - カスタム実装（Fiber v3対応）+ OpenAPI 3.0対応
	app.Get("/swagger/*", func(c fiber.Ctx) error {
		path := strings.TrimPrefix(c.Path(), "/swagger")

		// OpenAPI 3.0 仕様ファイル
		if path == "/openapi-v3.json" {
			return c.SendFile("../schemas/openapi-v3.json")
		}

		if path == "/openapi-v3.yaml" {
			return c.SendFile("../schemas/openapi-v3.yaml")
		}

		// デフォルトページ
		if path == "" || path == "/" || path == "/index.html" {
			return c.SendFile("./templates/swagger.html")
		}

		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Not found",
		})
	})

	// Root endpoint
	app.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Penguin Backend API",
			"version": "1.0.0",
			"docs":    "/swagger/index.html",
		})
	})

	// API group
	api := app.Group("/api")

	// Business Data Services のルートを設定
	if container.BusinessService != nil {
		// Initialize handlers
		businessHandler := handlers.NewBusinessHandler(container.BusinessService)

		// Setup routes for each domain
		SetupBusinessRoutes(api.Group("/business"), businessHandler)
	}

	// Media Data Services のルートを設定（将来実装）
	// if container.MediaData != nil {
	//     mediaHandler := handlers.NewMediaHandler(container.MediaData)
	//     SetupMediaRoutes(api, mediaHandler)
	// }
}
