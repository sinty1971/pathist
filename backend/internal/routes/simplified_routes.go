package routes

import (
	"penguin-backend/internal/services"

	"github.com/gofiber/fiber/v3"
)

// SetupSimplifiedRoutes はハンドラー層を省略した簡略化ルートを設定
func SetupSimplifiedRoutes(app *fiber.App, rootService *services.RootService) {
	// API group
	api := app.Group("/api")
	businessApi := api.Group("/business")

	// サービスから直接ルートを登録
	if service := rootService.GetService("CompanyService"); service != nil {
		if httpService, ok := (*service).(services.HTTPService); ok {
			httpService.RegisterRoutes(businessApi)
		}
	}

	if service := rootService.GetService("KojiService"); service != nil {
		if httpService, ok := (*service).(services.HTTPService); ok {
			httpService.RegisterRoutes(businessApi)
		}
	}

	if service := rootService.GetService("FileService"); service != nil {
		if httpService, ok := (*service).(services.HTTPService); ok {
			httpService.RegisterRoutes(businessApi)
		}
	}
}

// SetupAutoRoutes は登録されたサービスから自動的にルートを設定
func SetupAutoRoutes(app *fiber.App, rootService *services.RootService) {
	api := app.Group("/api")
	businessApi := api.Group("/business")

	// RootServiceに登録されている全サービスを自動で処理
	for serviceName, service := range rootService.GetAllServices() {
		if httpService, ok := service.(services.HTTPService); ok {
			// サービス名に基づいてグループを作成
			group := businessApi.Group("/" + serviceName)
			httpService.RegisterRoutes(group)
		}
	}
}