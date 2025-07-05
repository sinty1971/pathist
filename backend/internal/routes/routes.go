package routes

import (
	"penguin-backend/internal/handlers"
	"penguin-backend/internal/services"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// SetupRoutes はすべてのルートを設定します
func SetupRoutes(app *fiber.App, container *services.ServiceContainer) {
	// Swagger documentation
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Root endpoint
	app.Get("/", func(c *fiber.Ctx) error {
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
		fileHandler := handlers.NewFileHandler(container.BusinessService.FileService)
		projectHandler := handlers.NewProjectHandler(container.BusinessService.ProjectService)
		companyHandler := handlers.NewCompanyHandler(container.BusinessService.CompanyService)

		// Setup routes for each domain
		SetupFileRoutes(api, fileHandler)
		SetupProjectRoutes(api, projectHandler)
		SetupCompanyRoutes(api, companyHandler)
	}

	// Media Data Services のルートを設定（将来実装）
	// if container.MediaData != nil {
	//     mediaHandler := handlers.NewMediaHandler(container.MediaData)
	//     SetupMediaRoutes(api, mediaHandler)
	// }
}
