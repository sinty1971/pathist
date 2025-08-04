package main

import (
	"crypto/tls"
	"flag"
	"log"
	"os"
	"penguin-backend/internal/routes"
	"penguin-backend/internal/services"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"

	_ "penguin-backend/docs"
)

// @title Penguin ファイルシステム管理API
// @version 1.0.0
// @description ファイルエントリの管理と閲覧のためのAPI
// @servers.url http://localhost:8080/api
func main() {
	// コマンドライン引数の解析
	var (
		useHTTP2 = flag.Bool("http2", false, "Enable HTTP/2 with HTTPS")
		port     = flag.String("port", "8080", "Port to listen on")
		certFile = flag.String("cert", "cert.pem", "TLS certificate file")
		keyFile  = flag.String("key", "key.pem", "TLS private key file")
	)
	flag.Parse()

	// 使用方法の表示
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		log.Println("Usage:")
		log.Println("  go run cmd/main.go                    # HTTP/1.1 on port 8080")
		log.Println("  go run cmd/main.go -http2             # HTTP/2 + HTTPS on port 8080")
		log.Println("  go run cmd/main.go -port=9080         # HTTP/1.1 on port 9080")
		log.Println("  go run cmd/main.go -http2 -port=8443  # HTTP/2 + HTTPS on port 8443")
		os.Exit(0)
	}
	// サーバー設定
	var serverHeader, appName string
	if *useHTTP2 {
		serverHeader = "Penguin-Backend/1.0-HTTP2"
		appName = "Penguin Backend API HTTP/2"
	} else {
		serverHeader = "Penguin-Backend/1.0"
		appName = "Penguin Backend API"
	}

	app := fiber.New(fiber.Config{
		// 動的サーバー設定
		ServerHeader: serverHeader,
		AppName:      appName,

		// パフォーマンス設定
		ReadTimeout:     time.Second * 15,
		WriteTimeout:    time.Second * 15,
		IdleTimeout:     time.Second * 60,
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,

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
		// Level: compress.LevelBestCompression, // 圧縮率重視の場合
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
	var logFormat string
	if *useHTTP2 {
		logFormat = "[${time}] ${status} - ${method} ${path} - ${latency} - HTTP/2\n"
	} else {
		logFormat = "[${time}] ${status} - ${method} ${path} - ${latency}\n"
	}
	app.Use(logger.New(logger.Config{
		Format:     logFormat,
		TimeFormat: "2006-01-02 15:04:05",
	}))

	const defaultDatabaseFilename = ".detail.yaml"
	const defaultFileFolderPath = "~/penguin"
	const defaultBusinessFolderPath = "~/penguin/豊田築炉"
	const defaultCompanyFolderPath = "~/penguin/豊田築炉/1 会社"
	const defaultKojiFolderPath = "~/penguin/豊田築炉/2 工事"

	// ServiceContainerを作成
	cs := services.NewContainerService()

	// コンテナオプションを作成
	opt := services.ContainerOption{
		RootService: cs,
	}

	// ファイルサービスをリセット
	cs.FileService.BuildWithOption(opt, defaultFileFolderPath)
	// ビジネスサービスをリセット
	cs.BusinessService.BuildWithOption(opt, defaultBusinessFolderPath)
	// 会社サービスをリセット
	cs.BusinessService.CompanyService.BuildWithOption(opt, defaultCompanyFolderPath, defaultDatabaseFilename)
	// 工事サービスをリセット
	cs.BusinessService.KojiService.BuildWithContext(opt, defaultKojiFolderPath, defaultDatabaseFilename)

	// sc.MediaService, err := services.NewMediaDataService("~/penguin/homes/sinty/media", ".detail.yaml")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	defer cs.Cleanup()

	// ルートを設定
	routes.SetupRoutes(app, cs)

	// サーバー起動メッセージ
	if *useHTTP2 {
		log.Printf("🚀 HTTP/2 + HTTPS Server starting on :%s", *port)
		log.Printf("📖 API documentation: https://localhost:%s/swagger/index.html", *port)
		log.Printf("🔒 Using TLS certificate: %s", *certFile)
		log.Println("🌟 Features enabled:")
		log.Println("  ✅ HTTP/2 (h2) - Fiber v3 auto-enables")
		log.Println("  ✅ TLS 1.2+")
		log.Println("  ✅ Gzip compression")
		log.Println("  ✅ CORS")

		// HTTP/2 + HTTPS で起動
		log.Fatal(app.Listen("0.0.0.0:"+*port, fiber.ListenConfig{
			CertFile:      *certFile,
			CertKeyFile:   *keyFile,
			TLSMinVersion: tls.VersionTLS12,
		}))
	} else {
		log.Printf("🚀 HTTP/1.1 Server starting on :%s", *port)
		log.Printf("📖 API documentation: http://localhost:%s/swagger/index.html", *port)
		log.Println("🌟 Features enabled:")
		log.Println("  ✅ HTTP/1.1")
		log.Println("  ✅ Gzip compression")
		log.Println("  ✅ CORS")

		// HTTP/1.1で起動
		log.Fatal(app.Listen("0.0.0.0:" + *port))
	}
}
