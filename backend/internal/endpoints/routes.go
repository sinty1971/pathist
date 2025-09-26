package endpoints

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v2"

	app "penguin-backend/internal/app"
	"penguin-backend/internal/models"
)

// apiRootDescription はルートエンドポイントの説明
const apiRootDescription = "Penguin Backend API"

// SetupRoutes はすべてのルートを設定します
func SetupRoutes(app *fiber.App, api huma.API, container app.ServiceContainer) {
	_ = app // Fiber ルーター移行の名残。物理ルート追加時に再利用可。
	registerRoot(api)
	registerFileEndpoints(api, container)
	registerCompanyEndpoints(api, container)
	registerKojiEndpoints(api, container)
	registerIDSyncEndpoints(api)
}

func registerRoot(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "get-root",
		Method:      http.MethodGet,
		Path:        "/",
		Summary:     "API ルート情報取得",
		Description: "API のバージョンやドキュメントURLを返します",
		Tags:        []string{"system"},
	}, func(ctx context.Context, _ *struct{}) (*struct {
		Body struct {
			Message string `json:"message" example:"Penguin Backend API"`
			Version string `json:"version" example:"1.0.0"`
			Docs    string `json:"docs" example:"/docs"`
		}
	}, error) {
		resp := &struct {
			Body struct {
				Message string `json:"message" example:"Penguin Backend API"`
				Version string `json:"version" example:"1.0.0"`
				Docs    string `json:"docs" example:"/docs"`
			}
		}{}
		resp.Body.Message = apiRootDescription
		resp.Body.Version = "1.0.0"
		resp.Body.Docs = "/docs"
		return resp, nil
	})
}

func registerIDSyncEndpoints(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "generate-koji-id",
		Method:      http.MethodPost,
		Path:        "/id-sync/generate-koji",
		Summary:     "工事IDの生成",
		Description: "開始日・会社名・場所名から工事IDを生成します",
		Tags:        []string{"ID同期"},
	}, func(ctx context.Context, in *struct {
		Body struct {
			StartDate    string `json:"startDate" doc:"開始日 (ISO 8601)"`
			CompanyName  string `json:"companyName" doc:"会社名"`
			LocationName string `json:"locationName" doc:"場所名"`
		}
	}) (*struct {
		Body struct {
			ID       string    `json:"id"`
			IDType   string    `json:"idType"`
			Length   int       `json:"length"`
			Method   string    `json:"method"`
			Source   string    `json:"source"`
			Created  time.Time `json:"created"`
			IsValid  bool      `json:"isValid"`
			ErrorMsg string    `json:"errorMsg"`
		}
	}, error) {
		startDate, err := time.Parse(time.RFC3339, in.Body.StartDate)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid start date", err)
		}
		id := models.NewIDFromString(startDate.Format("20060102") + in.Body.CompanyName + in.Body.LocationName).Len5()
		resp := &struct {
			Body struct {
				ID       string    `json:"id"`
				IDType   string    `json:"idType"`
				Length   int       `json:"length"`
				Method   string    `json:"method"`
				Source   string    `json:"source"`
				Created  time.Time `json:"created"`
				IsValid  bool      `json:"isValid"`
				ErrorMsg string    `json:"errorMsg"`
			}
		}{}
		resp.Body.ID = id
		resp.Body.IDType = "koji"
		resp.Body.Length = len(id)
		resp.Body.Method = "BLAKE2b-256 + Base32"
		resp.Body.Source = fmt.Sprintf("%s_%s_%s", startDate.Format("20060102"), in.Body.CompanyName, in.Body.LocationName)
		resp.Body.Created = time.Now()
		resp.Body.IsValid = true
		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "generate-path-id",
		Method:      http.MethodPost,
		Path:        "/id-sync/generate-path",
		Summary:     "パスIDの生成",
		Description: "フルパスから Len7 ID を生成します",
		Tags:        []string{"ID同期"},
	}, func(ctx context.Context, in *struct {
		Body struct {
			FullPath string `json:"fullPath" doc:"パスID生成用フルパス"`
		}
	}) (*struct {
		Body struct {
			ID       string    `json:"id"`
			IDType   string    `json:"idType"`
			Length   int       `json:"length"`
			Method   string    `json:"method"`
			Source   string    `json:"source"`
			Created  time.Time `json:"created"`
			IsValid  bool      `json:"isValid"`
			ErrorMsg string    `json:"errorMsg"`
		}
	}, error) {
		if strings.TrimSpace(in.Body.FullPath) == "" {
			return nil, huma.Error400BadRequest("fullPath is required", nil)
		}
		id := models.NewIDFromString(in.Body.FullPath).Len7()
		resp := &struct {
			Body struct {
				ID       string    `json:"id"`
				IDType   string    `json:"idType"`
				Length   int       `json:"length"`
				Method   string    `json:"method"`
				Source   string    `json:"source"`
				Created  time.Time `json:"created"`
				IsValid  bool      `json:"isValid"`
				ErrorMsg string    `json:"errorMsg"`
			}
		}{}
		resp.Body.ID = id
		resp.Body.IDType = "path"
		resp.Body.Length = len(id)
		resp.Body.Method = "BLAKE2b-256 + Base32"
		resp.Body.Source = in.Body.FullPath
		resp.Body.Created = time.Now()
		resp.Body.IsValid = true
		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "validate-id",
		Method:      http.MethodPost,
		Path:        "/id-sync/validate",
		Summary:     "IDの検証",
		Description: "提供されたIDが正しいかどうかを検証します",
		Tags:        []string{"ID同期"},
	}, func(ctx context.Context, in *struct {
		Body struct {
			ID           string `json:"id"`
			StartDate    string `json:"startDate"`
			CompanyName  string `json:"companyName"`
			LocationName string `json:"locationName"`
			FullPath     string `json:"fullPath"`
		}
	}) (*struct {
		Body struct {
			IsValid    bool      `json:"isValid"`
			ExpectedID string    `json:"expectedId"`
			ProvidedID string    `json:"providedId"`
			NeedsSync  bool      `json:"needsSync"`
			ErrorMsg   string    `json:"errorMsg"`
			Method     string    `json:"method"`
			ComparedAt time.Time `json:"comparedAt"`
		}
	}, error) {
		resp := &struct {
			Body struct {
				IsValid    bool      `json:"isValid"`
				ExpectedID string    `json:"expectedId"`
				ProvidedID string    `json:"providedId"`
				NeedsSync  bool      `json:"needsSync"`
				ErrorMsg   string    `json:"errorMsg"`
				Method     string    `json:"method"`
				ComparedAt time.Time `json:"comparedAt"`
			}
		}{}
		resp.Body.ProvidedID = in.Body.ID
		resp.Body.ComparedAt = time.Now()

		if strings.TrimSpace(in.Body.FullPath) != "" {
			expected := models.NewIDFromString(in.Body.FullPath).Len7()
			resp.Body.ExpectedID = expected
			resp.Body.IsValid = in.Body.ID == expected
			resp.Body.NeedsSync = !resp.Body.IsValid
			resp.Body.Method = "path"
		} else {
			startDate, err := time.Parse(time.RFC3339, in.Body.StartDate)
			if err != nil {
				return nil, huma.Error400BadRequest("invalid start date", err)
			}
			expected := models.NewIDFromString(startDate.Format("20060102") + in.Body.CompanyName + in.Body.LocationName).Len5()
			resp.Body.ExpectedID = expected
			resp.Body.IsValid = in.Body.ID == expected
			resp.Body.NeedsSync = !resp.Body.IsValid
			resp.Body.Method = "koji"
		}

		if !resp.Body.IsValid {
			resp.Body.ErrorMsg = fmt.Sprintf("ID mismatch: %s != %s", resp.Body.ProvidedID, resp.Body.ExpectedID)
		}

		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "bulk-convert-ids",
		Method:      http.MethodPost,
		Path:        "/id-sync/bulk-convert",
		Summary:     "一括ID変換",
		Description: "複数のアイテムを一括でID変換します",
		Tags:        []string{"ID同期"},
	}, func(ctx context.Context, in *struct {
		Body struct {
			Items []struct {
				StartDate    string `json:"startDate"`
				CompanyName  string `json:"companyName"`
				LocationName string `json:"locationName"`
				FullPath     string `json:"fullPath"`
			}
		}
	}) (*struct {
		Body struct {
			Results []struct {
				ID       string    `json:"id"`
				IDType   string    `json:"idType"`
				Length   int       `json:"length"`
				Method   string    `json:"method"`
				Source   string    `json:"source"`
				Created  time.Time `json:"created"`
				IsValid  bool      `json:"isValid"`
				ErrorMsg string    `json:"errorMsg"`
			}
			TotalCount   int       `json:"totalCount"`
			SuccessCount int       `json:"successCount"`
			ErrorCount   int       `json:"errorCount"`
			ProcessedAt  time.Time `json:"processedAt"`
		}
	}, error) {
		resp := &struct {
			Body struct {
				Results []struct {
					ID       string    `json:"id"`
					IDType   string    `json:"idType"`
					Length   int       `json:"length"`
					Method   string    `json:"method"`
					Source   string    `json:"source"`
					Created  time.Time `json:"created"`
					IsValid  bool      `json:"isValid"`
					ErrorMsg string    `json:"errorMsg"`
				}
				TotalCount   int       `json:"totalCount"`
				SuccessCount int       `json:"successCount"`
				ErrorCount   int       `json:"errorCount"`
				ProcessedAt  time.Time `json:"processedAt"`
			}
		}{}
		resp.Body.ProcessedAt = time.Now()
		resp.Body.TotalCount = len(in.Body.Items)

		for _, item := range in.Body.Items {
			result := struct {
				ID       string    `json:"id"`
				IDType   string    `json:"idType"`
				Length   int       `json:"length"`
				Method   string    `json:"method"`
				Source   string    `json:"source"`
				Created  time.Time `json:"created"`
				IsValid  bool      `json:"isValid"`
				ErrorMsg string    `json:"errorMsg"`
			}{}
			if strings.TrimSpace(item.FullPath) != "" {
				id := models.NewIDFromString(item.FullPath).Len7()
				result.ID = id
				result.IDType = "path"
				result.Length = len(id)
				result.Method = "BLAKE2b-256 + Base32"
				result.Source = item.FullPath
				result.Created = time.Now()
				result.IsValid = true
				resp.Body.SuccessCount++
			} else {
				startDate, err := time.Parse(time.RFC3339, item.StartDate)
				if err != nil {
					result.IDType = "koji"
					result.IsValid = false
					result.ErrorMsg = "開始日の解析に失敗しました: " + err.Error()
					result.Created = time.Now()
					resp.Body.ErrorCount++
					resp.Body.Results = append(resp.Body.Results, result)
					continue
				}
				id := models.NewIDFromString(startDate.Format("20060102") + item.CompanyName + item.LocationName).Len5()
				result.ID = id
				result.IDType = "koji"
				result.Length = len(id)
				result.Method = "BLAKE2b-256 + Base32"
				result.Source = fmt.Sprintf("%s_%s_%s", startDate.Format("20060102"), item.CompanyName, item.LocationName)
				result.Created = time.Now()
				result.IsValid = true
				resp.Body.SuccessCount++
			}
			resp.Body.Results = append(resp.Body.Results, result)
		}

		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-id-mapping",
		Method:      http.MethodGet,
		Path:        "/id-sync/mapping",
		Summary:     "IDマッピングテーブルの取得",
		Description: "フルパスIDとLen7 IDのマッピングテーブルを取得します",
		Tags:        []string{"ID同期"},
	}, func(ctx context.Context, in *struct {
		Paths string `query:"paths" doc:"カンマ区切りのパスリスト"`
	}) (*struct {
		Body map[string]string `json:""`
	}, error) {
		if strings.TrimSpace(in.Paths) == "" {
			return nil, huma.Error400BadRequest("paths is required", nil)
		}
		mapping := make(map[string]string)
		for _, raw := range strings.Split(in.Paths, ",") {
			trimmed := strings.TrimSpace(raw)
			if trimmed == "" {
				continue
			}
			mapping[trimmed] = models.NewIDFromString(trimmed).Len7()
		}
		return &struct {
			Body map[string]string `json:""`
		}{Body: mapping}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-id-sync-config",
		Method:      http.MethodGet,
		Path:        "/id-sync/config",
		Summary:     "同期設定の取得",
		Description: "ID同期システムの設定情報を取得します",
		Tags:        []string{"ID同期"},
	}, func(ctx context.Context, _ *struct{}) (*struct {
		Body map[string]interface{} `json:""`
	}, error) {
		config := map[string]interface{}{
			"version":       "1.0.0",
			"algorithm":     "BLAKE2b-256",
			"encoding":      "Base32",
			"charset":       models.RadixTable,
			"kojiIdLength":  5,
			"pathIdLength":  7,
			"fullIdLength":  25,
			"maxCollisions": 32768,
			"supportedFormats": []string{
				"ISO8601",
				"RFC3339",
				"YYYY-MM-DD",
			},
			"lastUpdated": time.Now(),
		}
		return &struct {
			Body map[string]interface{} `json:""`
		}{Body: config}, nil
	})
}
