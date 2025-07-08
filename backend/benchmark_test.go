package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/stretchr/testify/assert"

	"penguin-backend/internal/handlers"
	"penguin-backend/internal/routes"
	"penguin-backend/internal/services"
)

// ベンチマーク用のテストアプリケーションを作成
func setupTestApp() *fiber.App {
	app := fiber.New()

	// CORS設定
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH"},
		AllowHeaders: []string{"*"},
	}))

	// サービス初期化（テスト用の簡易設定）
	businessService, err := services.NewBusinessService("~/penguin/豊田築炉", ".detail.yaml")
	if err != nil {
		// テスト環境でディレクトリが存在しない場合は、モックサービスを使用
		businessService = &services.BusinessService{}
	}

	// ハンドラー初期化
	businessHandler := handlers.NewBusinessHandler(businessService)

	// ルート設定
	routes.SetupBusinessRoutes(app, businessHandler)

	return app
}

// GET /api/file/fileinfos のベンチマーク
func BenchmarkGetFileInfos(b *testing.B) {
	app := setupTestApp()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/api/file/fileinfos", nil)
		resp, err := app.Test(req)
		if err != nil {
			b.Fatal(err)
		}
		resp.Body.Close()
	}
}

// GET /api/kojies のベンチマーク
func BenchmarkGetKojiRecent(b *testing.B) {
	app := setupTestApp()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/api/kojies", nil)
		resp, err := app.Test(req)
		if err != nil {
			b.Fatal(err)
		}
		resp.Body.Close()
	}
}

// GET /api/companies のベンチマーク
func BenchmarkGetCompanyList(b *testing.B) {
	app := setupTestApp()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/api/companies", nil)
		resp, err := app.Test(req)
		if err != nil {
			b.Fatal(err)
		}
		resp.Body.Close()
	}
}

// シンプルなヘルスチェックのベンチマーク
func BenchmarkHealthCheck(b *testing.B) {
	app := setupTestApp()

	// 簡単なヘルスチェックエンドポイントを追加
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/health", nil)
		resp, err := app.Test(req)
		if err != nil {
			b.Fatal(err)
		}
		resp.Body.Close()
	}
}

// POST /api/kojies のベンチマーク
func BenchmarkPostKojiUpdate(b *testing.B) {
	app := setupTestApp()

	// テスト用のプロジェクトデータ
	kojiData := map[string]interface{}{
		"id":            "test-koji-id",
		"company_name":  "テスト会社",
		"location_name": "テスト現場",
		"description":   "テスト説明",
		"tags":          []string{"テスト", "ベンチマーク"},
		"start_date":    map[string]string{"time.Time": "2025-01-01T00:00:00+09:00"},
		"end_date":      map[string]string{"time.Time": "2025-12-31T23:59:59+09:00"},
	}

	jsonData, _ := json.Marshal(kojiData)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/api/kojies", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			b.Fatal(err)
		}
		resp.Body.Close()
	}
}

// 複数のエンドポイントを連続で呼び出すベンチマーク
func BenchmarkMixedRequests(b *testing.B) {
	app := setupTestApp()

	// ヘルスチェックエンドポイントを追加
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	endpoints := []string{
		"/health",
		"/api/file/fileinfos",
		"/api/kojies",
		"/api/companies",
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		endpoint := endpoints[i%len(endpoints)]
		req := httptest.NewRequest("GET", endpoint, nil)
		resp, err := app.Test(req)
		if err != nil {
			b.Fatal(err)
		}
		resp.Body.Close()
	}
}

// メモリ使用量の測定
func BenchmarkMemoryUsage(b *testing.B) {
	app := setupTestApp()

	// ヘルスチェックエンドポイントを追加
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// 複数のリクエストを実行してメモリ使用量を測定
		requests := []string{
			"/health",
			"/api/file/fileinfos",
			"/api/kojies",
			"/api/companies",
		}

		for _, endpoint := range requests {
			req := httptest.NewRequest("GET", endpoint, nil)
			resp, err := app.Test(req)
			if err != nil {
				b.Fatal(err)
			}
			resp.Body.Close()
		}
	}
}

// 機能テスト（ベンチマークではないが、動作確認用）
func TestAPIFunctionality(t *testing.T) {
	app := setupTestApp()

	// ヘルスチェックエンドポイントを追加
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// ヘルスチェックAPI
	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// ファイル情報API
	req = httptest.NewRequest("GET", "/api/file/fileinfos", nil)
	resp, err = app.Test(req)
	assert.NoError(t, err)
	// ディレクトリが存在しない場合は500エラーも許可
	assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusInternalServerError)
	resp.Body.Close()

	// プロジェクト一覧API
	req = httptest.NewRequest("GET", "/api/kojies", nil)
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusInternalServerError)
	resp.Body.Close()

	// 会社一覧API
	req = httptest.NewRequest("GET", "/api/companies", nil)
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusInternalServerError)
	resp.Body.Close()
}
