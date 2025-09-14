package endpoints

import (
	"penguin-backend/internal/services"

	"github.com/gofiber/fiber/v3"
)

// HTTPService はHTTPハンドラー機能を持つサービスのインターフェース
type HTTPService interface {
	services.Service
	// RegisterRoutes はルートを登録する
	RegisterRoutes(router fiber.Router)
}

// BaseHTTPService はHTTPサービスの基底実装
type BaseHTTPService struct {
	rootService *services.RootService
}

// JSONResponse は統一的なJSONレスポンスを返す
func (s *BaseHTTPService) JSONResponse(c fiber.Ctx, data interface{}) error {
	return c.JSON(data)
}

// JSONError は統一的なエラーレスポンスを返す
func (s *BaseHTTPService) JSONError(c fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"error": message,
	})
}
