package services

import (
	"fmt"
	"os"
	"path/filepath"
	"penguin-backend/internal/models"
	"slices"
	"strings"
)

// CompanyService は会社データ操作を処理します
type CompanyService struct {
	// RootService はトップコンテナのインスタンス
	RootService *ContainerService

	// 工事一覧フォルダー
	FolderPath string

	// データベースサービス
	DatabaseService *DatabaseService[*models.Company]
}

// BuildWithOption は opt でCompanyServiceを初期化します
// folderName は会社一覧の絶対パスフォルダー名
// databaseFilename はデータベースファイルのファイル名
func (cs *CompanyService) BuildWithOption(opt ContainerOption, folderPath string, databaseFilename string) error {

	// CompanyCategoryReverseMapを初期化
	models.CompanyCategoryReverseMap = make(map[string]models.CompanyCategoryIndex)
	for code, category := range models.CompanyCategoryMap {
		models.CompanyCategoryReverseMap[category] = code
	}

	// folderPath がアクセス可能かチェック
	fi, err := os.Stat(folderPath)
	if err != nil {
		return err
	}

	// フォルダーではない場合はエラー
	if !fi.IsDir() {
		return fmt.Errorf("フォルダーではありません: %s", folderPath)
	}

	// ルートサービスを設定
	cs.RootService = opt.RootService

	// 会社一覧のフォルダー名を設定
	cs.FolderPath = folderPath

	// 会社フォルダー基準のDatabaseServiceを初期化
	cs.DatabaseService = NewDatabaseService[*models.Company](databaseFilename)

	return nil
}

// GetCompany は指定されたパスから会社を取得する
// folderName: 会社フォルダー名
func (cs *CompanyService) GetCompany(folderName string) (*models.Company, error) {
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
func (cs *CompanyService) GetCompanies(opts ...GetCompaniesOptions) []models.Company {
	opt := GetCompaniesOptionsFunc(opts...)

	// ファイルシステムから会社フォルダー一覧を取得
	entries, err := os.ReadDir(cs.FolderPath)
	if err != nil || len(entries) == 0 {
		return []models.Company{}
	}

	// 会社データモデルを作成
	companies := make([]models.Company, len(entries))
	count := 0
	for _, entry := range entries {
		// 会社データモデルを作成、これはデータベースアクセスを行いません
		entryPath := filepath.Join(cs.FolderPath, entry.Name())
		company, err := models.NewCompany(entryPath)
		if err != nil {
			continue
		}

		if opt.SyncDatabaseFile {
			_ = cs.SyncDatabaseFile(company)
		}

		companies[count] = *company
		count++
	}

	return companies[:count]
}

// SyncDatabaseFile は会社の属性データを読み込み、会社データモデルに反映する
func (cs *CompanyService) SyncDatabaseFile(target *models.Company) error {
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
func (cs *CompanyService) GetCompanyByID(id string) (*models.Company, error) {

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

	company, err := cs.GetCompany(companies[idx].FolderPath)
	if err != nil {
		return nil, err
	}
	return company, nil
}

// GetCompanyByName は会社名で会社インスタンスを取得します（大文字小文字を区別しない）
func (cs *CompanyService) GetCompanyByName(name string) (*models.Company, error) {

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
	dbCompany, err := cs.GetCompany(companies[idx].FolderPath)
	if err != nil {
		return nil, err
	}

	return dbCompany, nil
}

// Update は会社情報を更新し、必要に応じてフォルダー名も変更します
// target: 更新対象の会社データモデル
// target.FolderName: 変更前のフォルダー名を指定すること
func (cs *CompanyService) Update(target *models.Company) error {
	// 変更前のフォルダー名を取得
	prevFolderPath := target.FolderPath

	// 新しいフォルダー名を生成して変更が必要かチェック
	target.UpdateIdentity()

	// フォルダー名の変更が必要な場合のみ処理
	if target.FolderPath != prevFolderPath {
		// ファイルの移動が必要
		err := os.Rename(prevFolderPath, target.FolderPath)
		if err != nil {
			return err
		}
	}

	// 更新後の会社情報をデータベースファイルに反映
	return cs.DatabaseService.Save(target)
}

// Categories は会社のカテゴリー一覧を配列形式で取得
func (cs *CompanyService) Categories() map[models.CompanyCategoryIndex]string {
	return models.CompanyCategoryMap
}
