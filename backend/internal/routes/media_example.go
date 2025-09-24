package routes

// このファイルは将来のMediaDataService実装時の参考例です

/*
import (
	"penguin-backend/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

// SetupMediaRoutes はメディア関連のルートを設定します
func SetupMediaRoutes(api fiber.Router, handler *handlers.MediaHandler) {
	media := api.Group("/media")

	// メディアファイル一覧の取得
	// @Summary メディアファイル一覧取得
	// @Description 指定されたパスのメディアファイル一覧を取得
	// @Tags media
	// @Accept json
	// @Produce json
	// @Param path query string false "メディアパス"
	// @Success 200 {array} models.MediaFile
	// @Router /media/list [get]
	media.Get("/list", handler.ListMediaFiles)

	// メディアファイルのアップロード
	// @Summary メディアファイルアップロード
	// @Description メディアファイルをアップロード
	// @Tags media
	// @Accept multipart/form-data
	// @Produce json
	// @Param file formData file true "アップロードファイル"
	// @Success 200 {object} models.MediaFile
	// @Router /media/upload [post]
	media.Post("/upload", handler.UploadMediaFile)

	// メディアファイルのメタデータ取得
	// @Summary メディアメタデータ取得
	// @Description 指定されたメディアファイルのメタデータを取得
	// @Tags media
	// @Accept json
	// @Produce json
	// @Param id path string true "メディアファイルID"
	// @Success 200 {object} models.MediaMetadata
	// @Router /media/metadata/{id} [get]
	media.Get("/metadata/:id", handler.GetMediaMetadata)

	// サムネイル生成
	// @Summary サムネイル生成
	// @Description 指定されたメディアファイルのサムネイルを生成
	// @Tags media
	// @Accept json
	// @Produce json
	// @Param id path string true "メディアファイルID"
	// @Success 200 {object} models.Thumbnail
	// @Router /media/thumbnail/{id} [post]
	media.Post("/thumbnail/:id", handler.GenerateThumbnail)
}
*/
