package handlers

import (
	"fmt"
	"net/url"
	"penguin-backend/internal/models"
	"penguin-backend/internal/services"

	"github.com/gofiber/fiber/v3"
)

type BusinessHandler struct {
	businessService *services.BusinessService
}

func NewBusinessHandler(businessService *services.BusinessService) *BusinessHandler {
	return &BusinessHandler{
		businessService: businessService,
	}
}

// GetBusinessBasePath godoc
// @Summary      ビジネスベースパスの取得
// @Description  ビジネスベースパスを取得します
// @Tags         ビジネス管理
// @Produce      json
// @Success      200 {object} map[string]string "ビジネスベースパス"
// @Router       /business/base-path [get]
func (bh *BusinessHandler) GetBusinessBasePath(c fiber.Ctx) error {
	businessBasePath := bh.businessService.FileService.BasePath
	return c.JSON(fiber.Map{"businessBasePath": businessBasePath})
}

// GetAttributeFilename godoc
// @Summary      属性ファイル名の取得
// @Description  属性ファイル名を取得します
// @Tags         ビジネス管理
// @Produce      json
// @Success      200 {object} map[string]string "属性ファイル名"
// @Router       /business/attribute-filename [get]
func (bh *BusinessHandler) GetAttributeFilename(c fiber.Ctx) error {
	attributeFilename := bh.businessService.AttributeFilename
	return c.JSON(fiber.Map{"attributeFilename": attributeFilename})
}

// GetFiles godoc
// @Summary      ファイルエントリ一覧の取得
// @Description  指定されたパスからファイルとフォルダーの一覧を取得します
// @Tags         ファイル管理
// @Produce      json
// @Param        path query string false "取得するディレクトリのパス" default(2 工事)
// @Success      200 {array} models.FileInfo "正常なレスポンス"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /business/files [get]
func (bh *BusinessHandler) GetFiles(c fiber.Ctx) error {
	fsPath := c.Query("path", "")

	fileInfos, err := bh.businessService.FileService.GetFileInfos(fsPath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to read directory",
			"message": err.Error(),
		})
	}

	return c.JSON(fileInfos)
}

// GetCompanyByID godoc
// @Summary      会社詳細の取得
// @Description  指定されたIDの会社詳細を取得します
// @Tags         会社管理
// @Produce      json
// @Param        id path string true "会社ID"
// @Success      200 {object} models.Company "正常なレスポンス"
// @Failure      404 {object} map[string]string "会社が見つからない"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /business/companies/{id} [get]
func (bh *BusinessHandler) GetCompanyByID(c fiber.Ctx) error {
	company, err := bh.businessService.CompanyService.GetCompanyByID(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Company not found",
			"message": err.Error(),
		})
	}
	return c.JSON(company)
}

// GetCompanies godoc
// @Summary      会社一覧の取得
// @Description  会社フォルダーの一覧を取得します
// @Tags         会社管理
// @Produce      json
// @Success      200 {array} models.Company "正常なレスポンス"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /business/companies [get]
func (bh *BusinessHandler) GetCompanies(c fiber.Ctx) error {
	companies := bh.businessService.CompanyService.GetCompanies()
	return c.JSON(companies)
}

// UpdateCompanyRequest 会社情報更新リクエスト
type UpdateCompanyRequest struct {
	OldFolderName string         `json:"oldFolderName"` // 変更前のフォルダー名
	Company       models.Company `json:"company"`       // 更新後の会社情報
}

// UpdateCompany godoc
// @Summary      会社情報の更新
// @Description  会社情報を更新し、必要に応じてフォルダー名を変更します
// @Tags         会社管理
// @Accept       json
// @Produce      json
// @Param        body body UpdateCompanyRequest true "会社更新リクエスト"
// @Success      200 {object} models.Company "更新後の会社データ"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /business/companies [put]
func (bh *BusinessHandler) UpdateCompany(c fiber.Ctx) error {
	// リクエストボディから編集された会社を取得
	var request UpdateCompanyRequest
	if err := c.Bind().Body(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "リクエストボディの解析に失敗しました",
			"message": err.Error(),
		})
	}

	// 会社の情報を更新
	err := bh.businessService.CompanyService.Update(request.OldFolderName, &request.Company)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "会社データの保存に失敗しました",
			"message": err.Error(),
		})
	}

	return c.JSON(request.Company)
}

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
// @Router       /business/kojies [get]
func (bh *BusinessHandler) GetKojies(c fiber.Ctx) error {
	// クエリパラメータでフィルターを取得
	filter := c.Query("filter")

	// filter=recentの場合は最近の工事のみを取得
	if filter == "recent" {
		kojies := bh.businessService.KojiService.GetRecentKojies()
		return c.JSON(kojies)
	} else {
		fmt.Printf("filter: %s\n", filter)
		// フィルターが指定されていない場合は最近の工事を返す（互換性のため）
		kojies := bh.businessService.KojiService.GetRecentKojies()
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
// @Router       /business/kojies [put]
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

// RenameKojiMangedFile godoc
// @Summary      ManagedFileの名前を変更
// @Description  ManagedFileの名前を変更します。
// @Tags         工事管理
// @Accept       json
// @Produce      json
// @Param        body body RenameManagedFileRequest true "工事データと管理ファイル"
// @Success      200 {object} models.Koji "更新後の工事データ"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /business/kojies/managed-files [put]
func (bh *BusinessHandler) RenameKojiManagedFile(c fiber.Ctx) error {
	// リクエストボディから編集された工事を取得
	var request RenameManagedFileRequest
	if err := c.Bind().Body(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "リクエストボディの解析に失敗しました",
			"message": err.Error(),
		})
	}

	// ファイル名を変更
	renamedFiles := bh.businessService.KojiService.RenameManagedFile(request.Koji, request.Currents)

	// ファイル名変更後、最新の工事データを取得して返す
	if request.Koji.FileInfo.Name != "" {
		updatedKoji, err := bh.businessService.KojiService.GetKoji(request.Koji.FileInfo.Name)
		if err == nil {
			return c.JSON(updatedKoji)
		}
	}

	// 取得できない場合は変更されたファイル名リストを返す（後方互換性）
	return c.JSON(renamedFiles)
}
