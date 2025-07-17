package services

import (
	"fmt"
	"penguin-backend/internal/models"
	"slices"
	"strings"
)

// CompanyService は会社データ操作を処理します
type CompanyService struct {
	FileService     *FileService
	DatabaseService *DatabaseService[*models.Company]
}

// NewCompanyService は新しいCompanyServiceを作成します
func NewCompanyService(businessFileService *FileService, folderName string) (*CompanyService, error) {
	// 会社一覧フォルダーのフルパスを取得
	folderPath, err := businessFileService.GetFullpath(folderName)
	if err != nil {
		return nil, err
	}

	// 会社一覧フォルダー基準のFileServiceを初期化
	fileService, err := NewFileService(folderPath)
	if err != nil {
		return nil, err
	}

	// 会社フォルダー基準のDatabaseServiceを初期化
	databaseService := NewDatabaseService[*models.Company](fileService, ".detail.yaml")

	// CompanyServiceのインスタンスを作成
	return &CompanyService{
		FileService:     fileService,
		DatabaseService: databaseService,
	}, nil
}

// GetCompany は指定されたパスから会社を取得する
// folderName: 会社フォルダー名
func (cs *CompanyService) GetCompany(folderName string) (models.Company, error) {
	// 会社データモデルを作成
	company, err := models.NewCompany(folderName)
	if err != nil {
		return models.Company{}, fmt.Errorf("会社データモデルの作成に失敗しました: %v", err)
	}

	// データベースファイルと同期
	err = cs.SyncDatabaseFile(&company)
	if err != nil {
		return models.Company{}, fmt.Errorf("データベースファイルとの同期に失敗しました: %v", err)
	}

	return company, nil
}

// GetCompanies は会社一覧を取得する
func (cs *CompanyService) GetCompanies() []models.Company {
	// 会社フォルダー一覧をデータベースファイルと同期せずに取得
	companies := cs.GetCompaniesNoSyncDatabaseFile()

	for _, company := range companies {
		// データベースファイルと同期（エラーは無視）
		_ = cs.SyncDatabaseFile(&company)
	}

	return companies
}

// GetCompaniesNoSyncDatabaseFile は会社一覧を取得する、データベースへのアクセスは行わない
func (cs *CompanyService) GetCompaniesNoSyncDatabaseFile() []models.Company {
	// ファイルシステムから会社フォルダー一覧を取得
	fileInfos, err := cs.FileService.GetFileInfos()
	if err != nil || len(fileInfos) == 0 {
		return []models.Company{}
	}

	companies := make([]models.Company, len(fileInfos))
	count := 0
	for _, fileInfo := range fileInfos {
		// 会社データモデルを作成、これはデータベースアクセスを行いません
		company, err := models.NewCompany(fileInfo.Name)
		if err != nil {
			continue
		}

		companies[count] = company
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
		go func() {
			cs.DatabaseService.Save(target)
		}()
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
	// slices.IndexFuncを使用して高速化
	companies := cs.GetCompaniesNoSyncDatabaseFile()

	idx := slices.IndexFunc(companies, func(c models.Company) bool {
		return c.ID == id
	})
	if idx == -1 {
		return nil, fmt.Errorf("ID %s の会社が見つかりません", id)
	}

	company, err := cs.GetCompany(companies[idx].FolderName)
	if err != nil {
		return nil, err
	}
	return &company, nil
}

// GetCompanyByName は会社名で会社を取得します（大文字小文字を区別しない）
func (cs *CompanyService) GetCompanyByName(name string) (*models.Company, error) {

	// 会社一覧を取得、データベースファイルと同期は行わない
	companies := cs.GetCompaniesNoSyncDatabaseFile()

	// slices.IndexFuncを使用して高速化
	idx := slices.IndexFunc(companies, func(c models.Company) bool {
		return strings.ToLower(c.ShortName) == name ||
			strings.ToLower(c.LegalName) == name
	})
	if idx == -1 {
		return nil, fmt.Errorf("会社名 %s の会社が見つかりません", name)
	}

	// データベースファイルと同期
	databasedCompany, err := cs.GetCompany(companies[idx].FolderName)
	if err != nil {
		return nil, err
	}

	return &databasedCompany, nil
}

// Update は会社情報を更新し、必要に応じてフォルダー名も変更します
// target: 更新対象の会社データモデル
// target.FolderName: 変更前のフォルダー名を指定すること
func (cs *CompanyService) Update(target *models.Company) error {
	// 変更前のフォルダー名を取得
	prevFolderName := target.FolderName

	// 会社の業種とショート名から新しいアイデンティティフィールドを更新
	err := target.UpdateIdentity()
	if err != nil {
		return err
	}

	// エラーがない場合はファイルの移動が必要
	err = cs.FileService.MoveFile(prevFolderName, target.FolderName)
	if err != nil {
		return err
	}

	// 更新後の会社情報をデータベースファイルに反映
	return cs.DatabaseService.Save(target)
}

// Categories は会社のカテゴリー一覧をマップ形式で取得
func (cs *CompanyService) Categories() map[string]string {
	return models.GetAllCompanyCategoriesMap()
}
