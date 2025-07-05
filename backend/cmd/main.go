package main

import (
	"log"
	"penguin-backend/internal/routes"
	"penguin-backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	_ "penguin-backend/docs"
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

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// containerServiceを作成
	var err error
	sc := &services.ServiceContainer{}

	sc.BusinessService, err = services.NewBusinessDataService("~/penguin/豊田築炉", ".detail.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// sc.MediaService, err := services.NewMediaDataService("~/penguin/homes/sinty/media", ".detail.yaml")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	defer sc.Cleanup()

	// ルートを設定
	routes.SetupRoutes(app, sc)

	log.Println("Server starting on :8080")
	log.Println("API documentation available at http://localhost:8080/swagger/index.html")
	log.Fatal(app.Listen(":8080"))
}
