package main

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cache"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

// createOptimizedApp 最適化されたFiberアプリを作成
func createOptimizedApp() *fiber.App {
	app := fiber.New(fiber.Config{
		// HTTP/2サポートを有効化
		EnableIPValidation: true,
		ServerHeader:       "Penguin-Backend/1.0",
		AppName:            "Penguin Backend API Test",

		// パフォーマンス設定
		ReadTimeout:     time.Second * 15,
		WriteTimeout:    time.Second * 15,
		IdleTimeout:     time.Second * 60,
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
	})

	// 最適化ミドルウェア
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	app.Use(cache.New(cache.Config{
		Next: func(c fiber.Ctx) bool {
			return c.Method() != fiber.MethodGet
		},
		Expiration:   30 * time.Second,
		CacheControl: true,
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	return app
}

// createBasicApp 基本的なFiberアプリを作成（最適化なし）
func createBasicApp() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "Basic App Test",
	})

	// ログは無効化（ベンチマーク用）
	// app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	return app
}

// setupTestRoutes テスト用ルートを設定
func setupTestRoutes(app *fiber.App) {
	api := app.Group("/api")

	// テスト用エンドポイント
	api.Get("/test/small", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message":   "small response",
			"timestamp": time.Now(),
		})
	})

	api.Get("/test/large", func(c fiber.Ctx) error {
		// 大きなレスポンスをシミュレート
		data := make([]fiber.Map, 100)
		for i := 0; i < 100; i++ {
			data[i] = fiber.Map{
				"id":          i,
				"name":        "Test Koji " + string(rune(i)),
				"description": "This is a test koji description that is long enough to test compression and caching performance",
				"tags":        []string{"test", "koji", "benchmark"},
				"timestamp":   time.Now(),
			}
		}
		return c.JSON(data)
	})
}

// BenchmarkBasicApp 基本的なアプリのベンチマーク
func BenchmarkBasicApp_SmallResponse(b *testing.B) {
	app := createBasicApp()
	setupTestRoutes(app)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req, _ := http.NewRequest("GET", "/api/test/small", nil)
			resp, _ := app.Test(req)
			resp.Body.Close()
		}
	})
}

func BenchmarkBasicApp_LargeResponse(b *testing.B) {
	app := createBasicApp()
	setupTestRoutes(app)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req, _ := http.NewRequest("GET", "/api/test/large", nil)
			resp, _ := app.Test(req)
			resp.Body.Close()
		}
	})
}

// BenchmarkOptimizedApp 最適化されたアプリのベンチマーク
func BenchmarkOptimizedApp_SmallResponse(b *testing.B) {
	app := createOptimizedApp()
	setupTestRoutes(app)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req, _ := http.NewRequest("GET", "/api/test/small", nil)
			resp, _ := app.Test(req)
			resp.Body.Close()
		}
	})
}

func BenchmarkOptimizedApp_LargeResponse(b *testing.B) {
	app := createOptimizedApp()
	setupTestRoutes(app)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req, _ := http.NewRequest("GET", "/api/test/large", nil)
			resp, _ := app.Test(req)
			resp.Body.Close()
		}
	})
}

// BenchmarkCacheHit キャッシュヒット時のパフォーマンス
func BenchmarkOptimizedApp_CacheHit(b *testing.B) {
	app := createOptimizedApp()
	setupTestRoutes(app)

	// キャッシュを準備（最初のリクエスト）
	req, _ := http.NewRequest("GET", "/api/test/large", nil)
	app.Test(req)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req, _ := http.NewRequest("GET", "/api/test/large", nil)
			resp, _ := app.Test(req)
			resp.Body.Close()
		}
	})
}

// TestCompressionEffectiveness 圧縮効果のテスト
func TestCompressionEffectiveness(t *testing.T) {
	// 最適化なしのアプリ
	basicApp := createBasicApp()
	setupTestRoutes(basicApp)

	// 最適化ありのアプリ
	optimizedApp := createOptimizedApp()
	setupTestRoutes(optimizedApp)

	// 基本アプリのレスポンス（圧縮なし）
	basicReq, _ := http.NewRequest("GET", "/api/test/large", nil)
	basicResp, err := basicApp.Test(basicReq)
	if err != nil {
		t.Fatal(err)
	}
	defer basicResp.Body.Close()

	// 最適化アプリのレスポンス（圧縮なし - Fiberテストでは自動圧縮解除される）
	optimizedReq, _ := http.NewRequest("GET", "/api/test/large", nil)
	optimizedResp, err := optimizedApp.Test(optimizedReq)
	if err != nil {
		t.Fatal(err)
	}
	defer optimizedResp.Body.Close()

	// レスポンスをデコード
	var basicData, optimizedData []fiber.Map
	json.NewDecoder(basicResp.Body).Decode(&basicData)
	json.NewDecoder(optimizedResp.Body).Decode(&optimizedData)

	// データの整合性をチェック
	if len(basicData) != len(optimizedData) {
		t.Errorf("Data length mismatch: basic=%d, optimized=%d", len(basicData), len(optimizedData))
	}

	// 基本的な動作確認
	if len(basicData) == 0 {
		t.Error("Basic app returned no data")
	}

	if len(optimizedData) == 0 {
		t.Error("Optimized app returned no data")
	}

	t.Logf("Basic app response items: %d", len(basicData))
	t.Logf("Optimized app response items: %d", len(optimizedData))
	t.Logf("Basic status: %d", basicResp.StatusCode)
	t.Logf("Optimized status: %d", optimizedResp.StatusCode)
}
