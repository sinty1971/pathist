package main

import (
	"log"
	"penguin-backend/internal/handlers"
	"penguin-backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	_ "penguin-backend/docs"

	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title Penguin ファイルシステム管理API
// @version 1.0.0
// @description ファイルエントリの管理と閲覧のためのAPI
// @host localhost:8080
// @BasePath /api
// @schemes http
func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Swagger documentation
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// FileServiceを作成
	fileService, err := services.NewFileService("~/penguin")
	if err != nil {
		log.Fatal(err)
	}

	// ProjectServiceを作成
	projectService, err := services.NewProjectService("~/penguin/豊田築炉/2-工事")
	if err != nil {
		log.Fatal(err)
	}

	// Create handlers
	fileServiceHandler := handlers.NewFileHandler(fileService)
	projectHandler := handlers.NewProjectHandler(projectService)

	api := app.Group("/api")

	// File entries routes
	api.Get("/file/fileinfos", fileServiceHandler.GetFileInfos)

	// Project routes
	api.Get("/project/recent", projectHandler.GetRecentProjects)
	api.Post("/project/update", projectHandler.Update)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Penguin Backend API",
			"version": "1.0.0",
			"docs":    "/swagger/index.html",
		})
	})

	log.Println("Server starting on :8080")
	log.Println("API documentation available at http://localhost:8080/swagger/index.html")
	log.Fatal(app.Listen(":8080"))
}
