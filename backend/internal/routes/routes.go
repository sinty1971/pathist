package routes

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
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

func registerFileEndpoints(api huma.API, container app.ServiceContainer) {
	if container.File == nil {
		return
	}

	type input struct {
		Path string `query:"path" doc:"取得するディレクトリの相対パス" example:"2 工事"`
	}
	type output struct {
		Body []models.FileInfo `json:"" doc:"ファイルエントリ一覧"`
	}

	huma.Register(api, huma.Operation{
		OperationID: "list-files",
		Method:      http.MethodGet,
		Path:        "/files",
		Summary:     "ファイルエントリ一覧の取得",
		Description: "指定されたパスからファイルとフォルダーの一覧を取得します",
		Tags:        []string{"ファイル管理"},
	}, func(ctx context.Context, in *input) (*output, error) {
		var (
			entries []models.FileInfo
			err     error
		)
		if in.Path == "" {
			entries, err = container.File.GetFileInfos()
		} else {
			entries, err = container.File.GetFileInfos(in.Path)
		}
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to read directory", err)
		}
		return &output{Body: entries}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-file-base-path",
		Method:      http.MethodGet,
		Path:        "/files/base-path",
		Summary:     "ファイルサービスの基準パス取得",
		Description: "ファイルサービスが参照する基準パスを返します",
		Tags:        []string{"ファイル管理"},
	}, func(ctx context.Context, _ *struct{}) (*struct {
		Body struct {
			BasePath string `json:"basePath" doc:"ファイルサービスの基準パス"`
		}
	}, error) {
		resp := &struct {
			Body struct {
				BasePath string `json:"basePath" doc:"ファイルサービスの基準パス"`
			}
		}{}
		resp.Body.BasePath = container.File.TargetFolder
		return resp, nil
	})
}

func registerCompanyEndpoints(api huma.API, container app.ServiceContainer) {
	if container.Company == nil {
		return
	}

	huma.Register(api, huma.Operation{
		OperationID: "list-companies",
		Method:      http.MethodGet,
		Path:        "/companies",
		Summary:     "会社一覧の取得",
		Description: "会社フォルダーの一覧を取得します",
		Tags:        []string{"会社管理"},
	}, func(ctx context.Context, _ *struct{}) (*struct {
		Body []models.Company `json:"" doc:"会社一覧"`
	}, error) {
		companies := container.Company.GetCompanies()
		return &struct {
			Body []models.Company `json:"" doc:"会社一覧"`
		}{Body: companies}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-company",
		Method:      http.MethodGet,
		Path:        "/companies/{id}",
		Summary:     "会社詳細の取得",
		Description: "指定されたIDの会社詳細を取得します",
		Tags:        []string{"会社管理"},
	}, func(ctx context.Context, in *struct {
		ID string `path:"id" example:"toyota" doc:"会社ID"`
	}) (*struct {
		Body *models.Company `json:"" doc:"会社詳細"`
	}, error) {
		company, err := container.Company.GetCompanyByID(in.ID)
		if err != nil {
			return nil, huma.Error404NotFound("company not found", err)
		}
		return &struct {
			Body *models.Company `json:"" doc:"会社詳細"`
		}{Body: company}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "update-company",
		Method:      http.MethodPut,
		Path:        "/companies",
		Summary:     "会社情報の更新",
		Description: "会社情報を更新し、必要に応じてフォルダー名を変更します",
		Tags:        []string{"会社管理"},
	}, func(ctx context.Context, in *struct {
		Body models.Company `json:""`
	}) (*struct {
		Body models.Company `json:""`
	}, error) {
		if err := container.Company.Update(&in.Body); err != nil {
			return nil, huma.Error500InternalServerError("failed to update company", err)
		}
		return &struct {
			Body models.Company `json:""`
		}{Body: in.Body}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "list-company-categories",
		Method:      http.MethodGet,
		Path:        "/companies/categories",
		Summary:     "カテゴリー一覧の取得",
		Description: "カテゴリー一覧を取得します",
		Tags:        []string{"会社管理"},
	}, func(ctx context.Context, _ *struct{}) (*struct {
		Body []models.CompanyCategoryInfo `json:"" doc:"カテゴリー一覧"`
	}, error) {
		categories := container.Company.Categories()
		return &struct {
			Body []models.CompanyCategoryInfo `json:"" doc:"カテゴリー一覧"`
		}{Body: categories}, nil
	})
}

func registerKojiEndpoints(api huma.API, container app.ServiceContainer) {
	if container.Koji == nil {
		return
	}

	huma.Register(api, huma.Operation{
		OperationID: "list-kojies",
		Method:      http.MethodGet,
		Path:        "/kojies",
		Summary:     "工事一覧の取得",
		Description: "工事一覧を取得します。filter=recentで最近の工事のみを取得できます",
		Tags:        []string{"工事管理"},
	}, func(ctx context.Context, in *struct {
		Filter string `query:"filter" doc:"フィルター (recent: 最近の工事のみ)"`
	}) (*struct {
		Body []models.Koji `json:"" doc:"工事一覧"`
	}, error) {
		kojies := container.Koji.GetKojies()
		return &struct {
			Body []models.Koji `json:"" doc:"工事一覧"`
		}{Body: kojies}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-koji",
		Method:      http.MethodGet,
		Path:        "/kojies/{path}",
		Summary:     "工事詳細の取得",
		Description: "指定されたフォルダー名の工事を取得します",
		Tags:        []string{"工事管理"},
	}, func(ctx context.Context, in *struct {
		Path string `path:"path" doc:"Kojiフォルダー名"`
	}) (*struct {
		Body *models.Koji `json:"" doc:"工事詳細"`
	}, error) {
		decoded, err := url.QueryUnescape(in.Path)
		if err != nil {
			decoded = in.Path
		}
		koji, err := container.Koji.GetKoji(decoded)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to fetch koji", err)
		}
		return &struct {
			Body *models.Koji `json:"" doc:"工事詳細"`
		}{Body: koji}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "update-koji",
		Method:      http.MethodPut,
		Path:        "/kojies",
		Summary:     "工事情報の更新",
		Description: "工事のファイル情報を更新します",
		Tags:        []string{"工事管理"},
	}, func(ctx context.Context, in *struct {
		Body models.Koji `json:""`
	}) (*struct {
		Body models.Koji `json:""`
	}, error) {
		if err := container.Koji.Update(&in.Body); err != nil {
			return nil, huma.Error500InternalServerError("failed to update koji", err)
		}
		return &struct {
			Body models.Koji `json:""`
		}{Body: in.Body}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "rename-koji-standard-files",
		Method:      http.MethodPut,
		Path:        "/kojies/standard-files",
		Summary:     "標準ファイル名の変更",
		Description: "StandardFileの名前を変更し、最新の工事データを返します",
		Tags:        []string{"工事管理"},
	}, func(ctx context.Context, in *struct {
		Body struct {
			Koji     models.Koji `json:"koji"`
			Currents []string    `json:"currents"`
		}
	}) (*struct {
		Body interface{} `json:"" doc:"更新後の工事データ"`
	}, error) {
		folderName := in.Body.Koji.GetFolderName()
		if folderName == "" {
			return &struct {
				Body interface{} `json:"" doc:"更新後の工事データ"`
			}{Body: []string{}}, nil
		}
		updated, err := container.Koji.GetKoji(folderName)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to load updated koji", err)
		}
		return &struct {
			Body interface{} `json:"" doc:"更新後の工事データ"`
		}{Body: updated}, nil
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
