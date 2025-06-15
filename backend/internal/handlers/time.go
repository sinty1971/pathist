package handlers

import (
	"penguin-backend/internal/models"
	"penguin-backend/internal/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

type TimeHandler struct{}

func NewTimeHandler() *TimeHandler {
	return &TimeHandler{}
}

// ParseTime godoc
// @Summary      タイムスタンプの解析
// @Description  様々な日時文字列フォーマットを解析します
// @Tags         時刻処理
// @Accept       json
// @Produce      json
// @Param        request body models.TimeParseRequest true "解析する時刻文字列"
// @Success      200 {object} models.TimeParseResponse "正常なレスポンス"
// @Failure      400 {object} map[string]string "不正なリクエスト"
// @Router       /time/parse [post]
func (h *TimeHandler) ParseTime(c *fiber.Ctx) error {
	var req models.TimeParseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
	}

	parsedTime, err := utils.ParseTime(req.TimeString)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Failed to parse time",
			"message": err.Error(),
		})
	}

	parsedTimestamp := models.NewTimestamp(parsedTime)
	response := models.ConvertDateTimeInstant(parsedTimestamp)
	response.Original = req.TimeString
	return c.JSON(response)
}

// GetSupportedFormats godoc
// @Summary      サポートされる時刻フォーマット一覧
// @Description  サポートされているすべての日時フォーマットの一覧を取得します
// @Tags         時刻処理
// @Accept       json
// @Produce      json
// @Success      200 {object} models.SupportedFormatsResponse "正常なレスポンス"
// @Router       /time/formats [get]
func (th *TimeHandler) GetSupportedFormats(c *fiber.Ctx) error {
	formats := []models.TimeFormat{
		{Name: "RFC3339Nano", Pattern: time.RFC3339Nano, Example: "2006-01-02T15:04:05.999999999Z07:00"},
		{Name: "RFC3339", Pattern: time.RFC3339, Example: "2006-01-02T15:04:05Z07:00"},
		{Name: "ISO8601 with ns", Pattern: "2006-01-02T15:04:05.999999999", Example: "2006-01-02T15:04:05.123456789"},
		{Name: "ISO8601", Pattern: "2006-01-02T15:04:05", Example: "2006-01-02T15:04:05"},
		{Name: "DateTime", Pattern: "2006-01-02 15:04:05", Example: "2006-01-02 15:04:05"},
		{Name: "Date", Pattern: "2006-01-02", Example: "2006-01-02"},
		{Name: "CompactDate", Pattern: "20060102", Example: "20060102"},
		{Name: "SlashDate", Pattern: "2006/01/02", Example: "2006/01/02"},
		{Name: "DotDate", Pattern: "2006.01.02", Example: "2006.01.02"},
		{Name: "Special Compact", Pattern: "YYYY-MMDD", Example: "1971-0618"},
	}

	return c.JSON(models.SupportedFormatsResponse{
		Formats: formats,
	})
}
