package main

import (
	"crypto/tls"
	"flag"
	"log"
	"os"
	appbootstrap "penguin-backend/internal/app"
	"penguin-backend/internal/endpoints"
	"penguin-backend/internal/huma/fiberv2"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// @title Penguin ファイル情報管理API
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

	// Middleware - 順序が重要（圧縮 → キャッシュ → CORS → ログ）

	// 1. 圧縮ミドルウェア（最初に適用）
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // パフォーマンス重視
		// Level: compress.LevelBestCompression, // 圧縮率重視の場合
	}))

	// 2. キャッシュミドルウェア（無効化）
	// app.Use(cache.New(cache.Config{
	// 	Next: func(c *fiber.Ctx) bool {
	// 		// POST、PUT、DELETE、PATCHはキャッシュしない
	// 		return c.Method() != fiber.MethodGet
	// 	},
	// 	Expiration:   30 * time.Second,    // 30秒間キャッシュ
	// 	CacheHeader:  "X-Cache",           // キャッシュヘッダー
	// 	CacheControl: true,                // Cache-Controlヘッダーを追加
	// }))

	// 3. CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
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

	// 5. サービスを初期化（手書き DI）
	servicesContainer, err := appbootstrap.InitializeServices()
	if err != nil {
		log.Fatalf("failed to initialize services: %v", err)
	}

	rs := servicesContainer.Root
	defer rs.Cleanup()

	// 6. OpenAPI関連の設定
	config := huma.DefaultConfig("Penguin ファイルシステム管理API", "1.0.0")
	config.OpenAPI.Info.Description = "ファイルエントリの管理と閲覧のためのAPI"
	serverProtocol := "http"
	if *useHTTP2 {
		serverProtocol = "https"
	}
	config.OpenAPI.Servers = []*huma.Server{
		{URL: serverProtocol + "://localhost:" + *port + "/api"},
	}
	config.OpenAPIPath = "/openapi"
	config.DocsPath = ""
	config.SchemasPath = "/schemas"

	// 7. APIルーティングの設定
	apiGroup := app.Group("/api")
	api := fiberv2.NewWithGroup(app, apiGroup, config)

	apiGroup.Get("/docs", func(c *fiber.Ctx) error {
		c.Type("html", "utf-8")
		return c.SendString(`<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="referrer" content="same-origin" />
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no" />
    <title>Penguin Backend API Reference</title>
    <link href="https://unpkg.com/@stoplight/elements@9.0.0/styles.min.css" rel="stylesheet" />
    <script src="https://unpkg.com/@stoplight/elements@9.0.0/web-components.min.js" integrity="sha256-Tqvw1qE2abI+G6dPQBc5zbeHqfVwGoamETU3/TSpUw4=" crossorigin="anonymous"></script>
  </head>
  <body style="height: 100vh;">
    <elements-api
      apiDescriptionUrl="/api/openapi.yaml"
      router="hash"
      layout="sidebar"
      tryItCredentialsPolicy="same-origin"
    />
  </body>
</html>`)
	})

	app.Get("/docs", func(c *fiber.Ctx) error {
		return c.Redirect("/api/docs", fiber.StatusTemporaryRedirect)
	})

	// サービスのエンドポイントセットアップ
	endpoints.SetupRoutes(app, api, *servicesContainer)

	// 8. サーバー起動メッセージ
	if *useHTTP2 {
		log.Printf("🚀 HTTP/2 + HTTPS Server starting on :%s", *port)
		log.Printf("📖 API documentation: https://localhost:%s/api/docs", *port)
		log.Printf("🔒 Using TLS certificate: %s", *certFile)
		log.Println("🌟 Features enabled:")
		log.Println("  ✅ HTTP/2 (h2)")
		log.Println("  ✅ TLS 1.2+")
		log.Println("  ✅ Gzip compression")
		log.Println("  ✅ CORS")

		cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("failed to load TLS certificate: %v", err)
		}
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
			NextProtos:   []string{"h2", "http/1.1"},
		}
		listener, err := tls.Listen("tcp", "0.0.0.0:"+*port, tlsConfig)
		if err != nil {
			log.Fatalf("failed to start TLS listener: %v", err)
		}
		log.Fatal(app.Listener(listener))
	} else {
		log.Printf("🚀 HTTP/1.1 Server starting on :%s", *port)
		log.Printf("📖 API documentation: http://localhost:%s/api/docs", *port)
		log.Println("🌟 Features enabled:")
		log.Println("  ✅ HTTP/1.1")
		log.Println("  ✅ Gzip compression")
		log.Println("  ✅ CORS")

		// HTTP/1.1で起動
		log.Fatal(app.Listen("0.0.0.0:" + *port))
	}
}
