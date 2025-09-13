package business

import (
	"fmt"
	"net/url"
	"penguin-backend/internal/models"

	"github.com/gofiber/fiber/v3"
)

// GetKoji godoc
// @Summary      Kojiの取得
// @Description  Kojiを取得します。
// @Tags         工事管理
// @Produce      json
// @Param        path path string true "Kojiフォルダーのファイル名"
// @Success      200 {object} models.Koji "Koji"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /kojies/{path} [get]
func (bh *BusinessHandler) GetKojiByPath(c fiber.Ctx) error {
	// パスパラメータを取得
	path := c.Params("path")

	// URLデコード
	decodedPath, err := url.QueryUnescape(path)
	if err != nil {
		fmt.Printf("URL decode error: %v\n", err)
		decodedPath = path // デコードに失敗した場合は元のパスを使用
	}

	// KojiServiceを使用してKojiを取得
	koji, err := bh.businessService.KojiService.GetKoji(decodedPath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "工事の取得に失敗しました",
			"message": err.Error(),
		})
	}

	return c.JSON(koji)
}

// GetKojies godoc
// @Summary      工事一覧の取得
// @Description  工事一覧を取得します。filter=recentで最近の工事のみを取得できます。
// @Tags         工事管理
// @Produce      json
// @Param        filter query string false "フィルター (recent: 最近の工事のみ)"
// @Success      200 {array} models.Koji "工事一覧"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /kojies [get]
func (bh *BusinessHandler) GetKojies(c fiber.Ctx) error {
	// クエリパラメータでフィルターを取得
	filter := c.Query("filter")

	// filter=recentの場合は最近の工事のみを取得
	if filter == "recent" {
		kojies := bh.businessService.KojiService.GetKojies()
		return c.JSON(kojies)
	} else {
		fmt.Printf("filter: %s\n", filter)
		// フィルターが指定されていない場合は最近の工事を返す（互換性のため）
		kojies := bh.businessService.KojiService.GetKojies()
		return c.JSON(kojies)

	}
}

// UpdateKoji godoc
// @Summary      更新対象フォルダー（Koji.FileInfo）のファイル情報をKojiのメンバ変数で更新
// @Description  更新対象フォルダー（Koji.FileInfo）のファイル情報をKojiのメンバ変数で更新します。
// @Tags         工事管理
// @Accept       json
// @Produce      json
// @Param        body body models.Koji true "工事データ"
// @Success      200 {object} models.Koji "更新後の工事データ"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /kojies [put]
func (bh *BusinessHandler) UpdateKoji(c fiber.Ctx) error {
	// リクエストボディから編集された工事を取得
	var koji models.Koji
	if err := c.Bind().Body(&koji); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "リクエストボディの解析に失敗しました",
			"message": err.Error(),
		})
	}

	// kojiのファイル情報を更新（.detail.yamlも更新）
	err := bh.businessService.KojiService.Update(&koji)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "工事データの保存に失敗しました",
			"message": err.Error(),
		})
	}

	return c.JSON(koji)
}

// RenameManagedFileRequest リクエスト用の構造体
type RenameManagedFileRequest struct {
	Koji     models.Koji `json:"koji"`
	Currents []string    `json:"currents"`
}

// RenameKojiStandardFiles godoc
// @Summary      StandardFileの名前を変更
// @Description  StandardFileの名前を変更します。
// @Tags         工事管理
// @Accept       json
// @Produce      json
// @Param        body body RenameManagedFileRequest true "工事データと管理ファイル"
// @Success      200 {object} models.Koji "更新後の工事データ"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /kojies/standard-files [put]
func (bh *BusinessHandler) RenameKojiStandardFiles(c fiber.Ctx) error {
	// リクエストボディから編集された工事を取得
	var request RenameManagedFileRequest
	if err := c.Bind().Body(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "リクエストボディの解析に失敗しました",
			"message": err.Error(),
		})
	}

	// TODO: RenameStandardFileメソッドの実装が必要
	// renamedFiles := bh.businessService.KojiService.RenameStandardFile(request.Koji, request.Currents)

	// ファイル名変更後、最新の工事データを取得して返す
	if request.Koji.GetFolderName() != "" {
		updatedKoji, err := bh.businessService.KojiService.GetKoji(request.Koji.GetFolderName())
		if err == nil {
			return c.JSON(updatedKoji)
		}
	}

	// 取得できない場合は空の配列を返す（後方互換性）
	return c.JSON([]string{})
}
