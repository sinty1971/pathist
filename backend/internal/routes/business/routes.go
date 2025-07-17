package business

import (
	"penguin-backend/internal/handlers/business"
	"penguin-backend/internal/services"

	"github.com/gofiber/fiber/v3"
)

// SetupBusinessRoutes はビジネス関連のルートを設定します
func SetupBusinessRoutes(api fiber.Router, businessService *services.BusinessService) {
	// 一つのBusinessHandlerを初期化
	businessHandler := business.NewBusinessHandler(businessService)

	// 各機能のルートを設定
	setupBaseRoutes(api, businessHandler)
	setupFileRoutes(api, businessHandler)
	setupCompanyRoutes(api, businessHandler)
	setupKojiRoutes(api, businessHandler)
}