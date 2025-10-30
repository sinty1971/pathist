package services

import (
	"fmt"
	"grpc-backend/internal/models"
	"os"
	"path/filepath"
	"penguin-backend/internal/utils"
	"slices"
	"strings"
)

// CompanyServiceOld は会社データ操作を処理します
type CompanyServiceOld struct {
	// container はトップコンテナのインスタンス	RootService *RootService
	container *Container

	// 会社フォルダー
	TargetFolder string

	// データベースサービス
	DatabaseService *PersistService[*models.Company]
}

// Cleanup はサービスをクリーンアップする
func (cs *CompanyServiceOld) Cleanup() error {
	return nil
}

// GetPersistPath は会社の属性データのパスを返す
func (cs *CompanyServiceOld) GetPersistPath(company *models.Company) string {
	return filepath.Join(company.TargetFolder, cs.DatabaseService.persistFilename)
}

// BuildWithOption は opt でCompanyServiceを初期化します
// rs はルートサービス
// opts はオプション
func (cs *CompanyServiceOld) Initialize(container *Container, serviceOptions *Options) (*CompanyServiceOld, error) {

	// コンテナを設定
	cs.container = container

	// CompanyCategoryReverseMapを初期化
	models.CompanyCategoryReverseMap = make(map[string]models.CompanyCategoryIndex)
	for code, category := range models.CompanyCategoryMap {
		models.CompanyCategoryReverseMap[category] = code
	}
	// targetFolderのパスチェック
	targetFolder, err := utils.CleanAbsPath(serviceOptions.CompanyServiceManagedFolder)
	if err != nil {
		return nil, err
	}

	// targetFolder がアクセス可能かチェック
	fi, err := os.Stat(targetFolder)
	if err != nil {
		return nil, err
	}

	// フォルダーではない場合はエラー
	if !fi.IsDir() {
		return nil, fmt.Errorf("フォルダーではありません: %s", targetFolder)
	}

	// 会社一覧のフォルダー名を設定
	cs.TargetFolder = targetFolder

	// 会社フォルダー基準のDatabaseServiceを初期化
	cs.DatabaseService = NewPersistService[*models.Company](serviceOptions.PersistFilename)

	return cs, nil
}

// GetCompany は指定されたパスから会社を取得する
// folderName: 会社フォルダー名

func (cs *CompanyServiceOld) GetCompany(folderName string) (*models.Company, error) {
	// 会社データモデルを作成
	company, err := models.NewCompany(folderName)
	if err != nil {
		return nil, fmt.Errorf("会社データモデルの作成に失敗しました: %v", err)
	}

	// データベースファイルと同期
	err = cs.SyncDatabaseFile(company)
	if err != nil {
		return nil, fmt.Errorf("データベースファイルとの同期に失敗しました: %v", err)
	}

	return company, nil
}

type GetCompaniesOptions struct {
	SyncDatabaseFile bool
}

const (
	SyncDatabaseFile = true
)

func GetCompaniesOptionsFunc(opts ...GetCompaniesOptions) GetCompaniesOptions {
	opt := GetCompaniesOptions{
		SyncDatabaseFile: true,
	}
	for _, o := range opts {
		opt = o
	}
	return opt
}

// GetCompanies は会社一覧を取得する
func (cs *CompanyServiceOld) GetCompanies(opts ...GetCompaniesOptions) []models.Company {
	opt := GetCompaniesOptionsFunc(opts...)

}

// SyncDatabaseFile は会社の属性データを読み込み、会社データモデルに反映する
func (cs *CompanyServiceOld) SyncDatabaseFile(target *models.Company) error {
	// 会社の属性データを読み込む
	database, err := cs.DatabaseService.Load(target)
	if err != nil {
		// ファイルが存在しない場合は新規ファイルとして非同期で保存
		cs.DatabaseService.Save(target)
		return nil
	}

	// データベースにしか保持されない内容をtargetに反映
	target.LegalName = database.LegalName
	target.PostalCode = database.PostalCode
	target.Address = database.Address
	target.Phone = database.Phone
	target.Email = database.Email
	target.Website = database.Website
	target.Tags = database.Tags

	return nil
}

// GetCompanyByID は指定されたIDの会社を取得します
func (cs *CompanyServiceOld) GetCompanyByID(id string) (*models.Company, error) {

	// 会社一覧を取得、データベースファイルと同期は行わない
	opt := GetCompaniesOptions{SyncDatabaseFile: false}
	companies := cs.GetCompanies(opt)

	// slices.IndexFuncを使用して高速化
	idx := slices.IndexFunc(companies, func(c models.Company) bool {
		return c.ID == id
	})
	if idx == -1 {
		return nil, fmt.Errorf("ID %s の会社が見つかりません", id)
	}

	company, err := cs.GetCompany(companies[idx].TargetFolder)
	if err != nil {
		return nil, err
	}
	return company, nil
}

// GetCompanyByName は会社名で会社インスタンスを取得します（大文字小文字を区別しない）
func (cs *CompanyServiceOld) GetCompanyByName(name string) (*models.Company, error) {

	// 会社一覧を取得、データベースファイルと同期は行わない
	opt := GetCompaniesOptions{SyncDatabaseFile: false}
	companies := cs.GetCompanies(opt)

	// slices.IndexFuncを使用して高速化
	idx := slices.IndexFunc(companies, func(c models.Company) bool {
		return strings.ToLower(c.ShortName) == name ||
			strings.ToLower(c.LegalName) == name
	})
	if idx == -1 {
		return nil, fmt.Errorf("会社名 %s の会社が見つかりません", name)
	}

	// データベースファイルと同期
	dbCompany, err := cs.GetCompany(companies[idx].TargetFolder)
	if err != nil {
		return nil, err
	}

	return dbCompany, nil
}

// Update は会社情報を更新し、必要に応じてフォルダー名も変更します
// updateCompany: 更新対象の会社データモデル
func (cs *CompanyServiceOld) Update(updateCompany *models.Company) error {
	// 変更前のフォルダー名を保持しておく
	prevFolderPath := updateCompany.TargetFolder

	// 新しいフォルダー名を生成して変更が必要かチェック
	updateCompany.UpdateIdentity()

	// フォルダー名の変更が必要な場合のみ処理
	if updateCompany.TargetFolder != prevFolderPath {
		// ファイルの移動が必要
		err := os.Rename(prevFolderPath, updateCompany.TargetFolder)
		if err != nil {
			return err
		}
	}

	// 更新後の会社情報をデータベースファイルに反映
	return cs.DatabaseService.Save(updateCompany)
}

// GetCategories は会社のカテゴリー一覧を配列形式で取得
func (cs *CompanyServiceOld) GetCategories() map[models.CompanyCategoryIndex]string {
	return models.CompanyCategoryMap
}

// Categories は会社カテゴリーを配列形式で取得します。
func (cs *CompanyServiceOld) Categories() []models.CompanyCategory {
	keys := make([]models.CompanyCategoryIndex, 0, len(models.CompanyCategoryMap))
	for code := range models.CompanyCategoryMap {
		keys = append(keys, code)
	}
	slices.Sort(keys)

	results := make([]models.CompanyCategory, 0, len(keys))
	for _, code := range keys {
		results = append(results, models.CompanyCategory{Index: code, Label: models.CompanyCategoryMap[code]})
	}
	return results
}
