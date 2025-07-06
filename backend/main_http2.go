// HTTP/2版のメインファイル - main.goとは別に実行する
// 実行方法: go run main_http2.go
package main

import (
	"crypto/tls"
	"log"
	"net"
	"penguin-backend/internal/routes"
	"penguin-backend/internal/services"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"

	_ "penguin-backend/docs"
)

// @title Penguin ファイルシステム管理API (HTTP/2)
// @version 1.0.0
// @description ファイルエントリの管理と閲覧のためのAPI (HTTP/2 + HTTPS対応)
// @servers.url https://localhost:8443/api
func main() {
	app := fiber.New(fiber.Config{
		// HTTP/2サポートを有効化
		EnableIPValidation: true,
		ServerHeader:       "Penguin-Backend/1.0-HTTP2",
		AppName:           "Penguin Backend API HTTP/2",
		
		// パフォーマンス設定
		ReadTimeout:       time.Second * 15,
		WriteTimeout:      time.Second * 15,
		IdleTimeout:       time.Second * 60,
		ReadBufferSize:    4096,
		WriteBufferSize:   4096,
		
		// エラーハンドリング
		ErrorHandler: func(c fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware - 順序が重要（圧縮 → キャッシュ → CORS → ログ）
	
	// 1. 圧縮ミドルウェア（最初に適用）
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // パフォーマンス重視
	}))
	
	// 2. キャッシュミドルウェア（無効化）
	// app.Use(cache.New(cache.Config{
	// 	Next: func(c fiber.Ctx) bool {
	// 		// POST、PUT、DELETE、PATCHはキャッシュしない
	// 		return c.Method() != fiber.MethodGet
	// 	},
	// 	Expiration:   30 * time.Second,    // 30秒間キャッシュ
	// 	CacheHeader:  "X-Cache",           // キャッシュヘッダー
	// 	CacheControl: true,                // Cache-Controlヘッダーを追加
	// }))
	
	// 3. CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))
	
	// 4. ログ（最後に適用）
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} - ${latency} - HTTP/2\n",
		TimeFormat: "2006-01-02 15:04:05",
	}))

	// containerServiceを作成
	var err error
	sc := &services.ServiceContainer{}

	sc.BusinessService, err = services.NewBusinessDataService("~/penguin/豊田築炉", ".detail.yaml")
	if err != nil {
		log.Fatal(err)
	}

	defer sc.Cleanup()

	// ルートを設定
	routes.SetupRoutes(app, sc)

	log.Println("🚀 HTTP/2 + HTTPS Server starting on :8443")
	log.Println("📖 API documentation: https://localhost:8443/swagger/index.html")
	log.Println("🔒 Using self-signed certificate (cert.pem + key.pem)")
	
	log.Println("🌟 Features enabled:")
	log.Println("  ✅ HTTP/2 (h2) - Fiber v3 auto-enables")
	log.Println("  ✅ TLS 1.2+")
	log.Println("  ✅ Gzip compression")
	log.Println("  ✅ Intelligent caching")
	log.Println("  ✅ CORS")

	// HTTP/2 + HTTPSで起動
	// 証明書を読み込み
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Fatalf("証明書の読み込みに失敗: %v", err)
	}

	// TLS設定
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"h2", "http/1.1"}, // HTTP/2を優先
	}

	// リスナーを作成
	ln, err := net.Listen("tcp", ":8443")
	if err != nil {
		log.Fatalf("リスナーの作成に失敗: %v", err)
	}

	// TLSリスナーでラップ
	tlsLn := tls.NewListener(ln, tlsConfig)

	log.Fatal(app.Listener(tlsLn))
}