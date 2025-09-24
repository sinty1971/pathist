package handlers

import (
	"fmt"
	"strings"
	"time"

	"penguin-backend/internal/models"
	"penguin-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

// IDSyncHandler はID同期関連のHTTPハンドラーを提供
type IDSyncHandler struct {
	businessService *services.BusinessService
}

// NewIDSyncHandler はIDSyncHandlerの新しいインスタンスを作成
func NewIDSyncHandler(businessService *services.BusinessService) *IDSyncHandler {
	return &IDSyncHandler{
		businessService: businessService,
	}
}

// ID生成リクエストの構造体
type IDGenerationRequest struct {
	StartDate    string `json:"startDate"`    // ISO 8601形式
	CompanyName  string `json:"companyName"`  // 会社名
	LocationName string `json:"locationName"` // 場所名
	FullPath     string `json:"fullPath"`     // フルパス（パスID生成用）
}

// ID生成レスポンスの構造体
type IDGenerationResponse struct {
	ID       string    `json:"id"`       // 生成されたID
	IDType   string    `json:"idType"`   // IDの種類（koji/path）
	Length   int       `json:"length"`   // ID長
	Method   string    `json:"method"`   // 生成方法
	Source   string    `json:"source"`   // 生成元データ
	Created  time.Time `json:"created"`  // 生成日時
	IsValid  bool      `json:"isValid"`  // 有効性
	ErrorMsg string    `json:"errorMsg"` // エラーメッセージ
}

// ID検証リクエストの構造体
type IDValidationRequest struct {
	ID           string `json:"id"`           // 検証対象のID
	StartDate    string `json:"startDate"`    // ISO 8601形式
	CompanyName  string `json:"companyName"`  // 会社名
	LocationName string `json:"locationName"` // 場所名
	FullPath     string `json:"fullPath"`     // フルパス
}

// ID検証レスポンスの構造体
type IDValidationResponse struct {
	IsValid    bool      `json:"isValid"`    // 有効性
	ExpectedID string    `json:"expectedId"` // 期待されるID
	ProvidedID string    `json:"providedId"` // 提供されたID
	NeedsSync  bool      `json:"needsSync"`  // 同期が必要かどうか
	ErrorMsg   string    `json:"errorMsg"`   // エラーメッセージ
	Method     string    `json:"method"`     // 検証方法
	ComparedAt time.Time `json:"comparedAt"` // 比較日時
}

// 一括変換リクエストの構造体
type BulkConversionRequest struct {
	Items []IDGenerationRequest `json:"items"` // 変換対象のアイテム
}

// 一括変換レスポンスの構造体
type BulkConversionResponse struct {
	Results      []IDGenerationResponse `json:"results"`      // 変換結果
	TotalCount   int                    `json:"totalCount"`   // 総数
	SuccessCount int                    `json:"successCount"` // 成功数
	ErrorCount   int                    `json:"errorCount"`   // エラー数
	ProcessedAt  time.Time              `json:"processedAt"`  // 処理日時
}

// GenerateKojiID godoc
// @Summary      工事IDの生成
// @Description  開始日・会社名・場所名から工事IDを生成します
// @Tags         ID同期
// @Accept       json
// @Produce      json
// @Param        request body IDGenerationRequest true "ID生成リクエスト"
// @Success      200 {object} IDGenerationResponse "正常なレスポンス"
// @Failure      400 {object} map[string]string "リクエストエラー"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /api/id-sync/generate-koji [post]
func (h *IDSyncHandler) GenerateKojiID(c *fiber.Ctx) error {
	var req IDGenerationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "リクエストの解析に失敗しました"})
	}

	// 開始日をパース
	startDate, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "開始日の形式が正しくありません"})
	}

	// 工事IDを生成
	id, err := generateKojiID(startDate, req.CompanyName, req.LocationName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "ID生成に失敗しました"})
	}

	response := IDGenerationResponse{
		ID:       id,
		IDType:   "koji",
		Length:   len(id),
		Method:   "BLAKE2b-256 + Base32",
		Source:   fmt.Sprintf("%s_%s_%s", startDate.Format("20060102"), req.CompanyName, req.LocationName),
		Created:  time.Now(),
		IsValid:  true,
		ErrorMsg: "",
	}

	return c.JSON(response)
}

// GeneratePathID godoc
// @Summary      パスIDの生成
// @Description  フルパスからLen7 IDを生成します
// @Tags         ID同期
// @Accept       json
// @Produce      json
// @Param        request body IDGenerationRequest true "ID生成リクエスト"
// @Success      200 {object} IDGenerationResponse "正常なレスポンス"
// @Failure      400 {object} map[string]string "リクエストエラー"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /api/id-sync/generate-path [post]
func (h *IDSyncHandler) GeneratePathID(c *fiber.Ctx) error {
	var req IDGenerationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "リクエストの解析に失敗しました"})
	}

	if req.FullPath == "" {
		return c.Status(400).JSON(fiber.Map{"error": "フルパスが指定されていません"})
	}

	// パスIDを生成
	id := models.NewIDFromString(req.FullPath).Len7()

	response := IDGenerationResponse{
		ID:       id,
		IDType:   "path",
		Length:   len(id),
		Method:   "BLAKE2b-256 + Base32",
		Source:   req.FullPath,
		Created:  time.Now(),
		IsValid:  true,
		ErrorMsg: "",
	}

	return c.JSON(response)
}

// ValidateID godoc
// @Summary      IDの検証
// @Description  提供されたIDが正しいかどうかを検証します
// @Tags         ID同期
// @Accept       json
// @Produce      json
// @Param        request body IDValidationRequest true "ID検証リクエスト"
// @Success      200 {object} IDValidationResponse "正常なレスポンス"
// @Failure      400 {object} map[string]string "リクエストエラー"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /api/id-sync/validate [post]
func (h *IDSyncHandler) ValidateID(c *fiber.Ctx) error {
	var req IDValidationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "リクエストの解析に失敗しました"})
	}

	var expectedID string
	var method string

	if req.FullPath != "" {
		// パスIDの検証
		expectedID = models.NewIDFromString(req.FullPath).Len7()
		method = "path"
	} else {
		// 工事IDの検証
		startDate, err := time.Parse(time.RFC3339, req.StartDate)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "開始日の形式が正しくありません"})
		}

		expectedID, err = generateKojiID(startDate, req.CompanyName, req.LocationName)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "期待されるID生成に失敗しました"})
		}
		method = "koji"
	}

	response := IDValidationResponse{
		IsValid:    req.ID == expectedID,
		ExpectedID: expectedID,
		ProvidedID: req.ID,
		NeedsSync:  req.ID != expectedID,
		Method:     method,
		ComparedAt: time.Now(),
	}

	if !response.IsValid {
		response.ErrorMsg = fmt.Sprintf("IDが一致しません: %s != %s", req.ID, expectedID)
	}

	return c.JSON(response)
}

// BulkConvert godoc
// @Summary      一括ID変換
// @Description  複数のアイテムを一括でID変換します
// @Tags         ID同期
// @Accept       json
// @Produce      json
// @Param        request body BulkConversionRequest true "一括変換リクエスト"
// @Success      200 {object} BulkConversionResponse "正常なレスポンス"
// @Failure      400 {object} map[string]string "リクエストエラー"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /api/id-sync/bulk-convert [post]
func (h *IDSyncHandler) BulkConvert(c *fiber.Ctx) error {
	var req BulkConversionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "リクエストの解析に失敗しました"})
	}

	results := make([]IDGenerationResponse, 0, len(req.Items))
	successCount := 0
	errorCount := 0

	for _, item := range req.Items {
		var result IDGenerationResponse

		if item.FullPath != "" {
			// パスID生成
			id := models.NewIDFromString(item.FullPath).Len7()
			result = IDGenerationResponse{
				ID:      id,
				IDType:  "path",
				Length:  len(id),
				Method:  "BLAKE2b-256 + Base32",
				Source:  item.FullPath,
				Created: time.Now(),
				IsValid: true,
			}
			successCount++
		} else {
			// 工事ID生成
			startDate, err := time.Parse(time.RFC3339, item.StartDate)
			if err != nil {
				result = IDGenerationResponse{
					ID:       "",
					IDType:   "koji",
					IsValid:  false,
					ErrorMsg: "開始日の解析に失敗しました: " + err.Error(),
					Created:  time.Now(),
				}
				errorCount++
			} else {
				id, err := generateKojiID(startDate, item.CompanyName, item.LocationName)
				if err != nil {
					result = IDGenerationResponse{
						ID:       "",
						IDType:   "koji",
						IsValid:  false,
						ErrorMsg: "ID生成に失敗しました: " + err.Error(),
						Created:  time.Now(),
					}
					errorCount++
				} else {
					result = IDGenerationResponse{
						ID:      id,
						IDType:  "koji",
						Length:  len(id),
						Method:  "BLAKE2b-256 + Base32",
						Source:  fmt.Sprintf("%s_%s_%s", startDate.Format("20060102"), item.CompanyName, item.LocationName),
						Created: time.Now(),
						IsValid: true,
					}
					successCount++
				}
			}
		}

		results = append(results, result)
	}

	response := BulkConversionResponse{
		Results:      results,
		TotalCount:   len(req.Items),
		SuccessCount: successCount,
		ErrorCount:   errorCount,
		ProcessedAt:  time.Now(),
	}

	return c.JSON(response)
}

// GetIDMapping godoc
// @Summary      IDマッピングテーブルの取得
// @Description  フルパスIDとLen7 IDのマッピングテーブルを取得します
// @Tags         ID同期
// @Accept       json
// @Produce      json
// @Param        paths query string false "カンマ区切りのパスリスト"
// @Success      200 {object} map[string]string "IDマッピングテーブル"
// @Failure      400 {object} map[string]string "リクエストエラー"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /api/id-sync/mapping [get]
func (h *IDSyncHandler) GetIDMapping(c *fiber.Ctx) error {
	pathsParam := c.Query("paths")
	if pathsParam == "" {
		return c.Status(400).JSON(fiber.Map{"error": "パスリストが指定されていません"})
	}

	paths := strings.Split(pathsParam, ",")
	mapping := make(map[string]string)

	for _, path := range paths {
		path = strings.TrimSpace(path)
		if path != "" {
			len7ID := models.NewIDFromString(path).Len7()
			mapping[path] = len7ID
		}
	}

	return c.JSON(mapping)
}

// ヘルパー関数: 工事IDを生成
func generateKojiID(startDate time.Time, companyName, locationName string) (string, error) {
	// 日付をYYYYMMDD形式に変換
	dateStr := startDate.Format("20060102")

	// 結合文字列を作成
	combined := dateStr + companyName + locationName

	// IDを生成
	id := models.NewIDFromString(combined).Len5()

	return id, nil
}

// ヘルパー関数: 工事フォルダー名を生成
func generateKojiFolderName(startDate time.Time, companyName, locationName string) string {
	year := startDate.Year()
	month := int(startDate.Month())
	day := startDate.Day()

	// YYYY-MMDD形式
	dateStr := fmt.Sprintf("%d-%02d%02d", year, month, day)

	return fmt.Sprintf("%s %s %s", dateStr, companyName, locationName)
}

// SyncConfiguration godoc
// @Summary      同期設定の取得
// @Description  ID同期システムの設定情報を取得します
// @Tags         ID同期
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]interface{} "同期設定情報"
// @Router       /api/id-sync/config [get]
func (h *IDSyncHandler) SyncConfiguration(c *fiber.Ctx) error {
	config := map[string]interface{}{
		"version":       "1.0.0",
		"algorithm":     "BLAKE2b-256",
		"encoding":      "Base32",
		"charset":       models.RadixTable,
		"kojiIdLength":  5,
		"pathIdLength":  7,
		"fullIdLength":  25,
		"maxCollisions": 32768, // 32^5 / 1000 の概算
		"supportedFormats": []string{
			"ISO8601",
			"RFC3339",
			"YYYY-MM-DD",
		},
		"lastUpdated": time.Now(),
	}

	return c.JSON(config)
}
