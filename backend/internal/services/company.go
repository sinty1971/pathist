package services

import (
	"fmt"
	"penguin-backend/internal/models"
	"strings"
)

// CompanyService handles company data operations
type CompanyService struct {
	FileService      *FileService
	AttributeService *AttributeService[*models.Company]
	FolderName       string
}

// NewCompanyService creates a new CompanyService
func NewCompanyService(businessFileService *FileService, folderName string) (*CompanyService, error) {
	cs := &CompanyService{}

	// フォルダー名を設定
	cs.FolderName = folderName

	// フォルダーのフルパスの取得
	folderPath, err := businessFileService.GetFullpath(folderName)
	if err != nil {
		return nil, err
	}

	// FileServiceを初期化
	cs.FileService, err = NewFileService(folderPath)
	if err != nil {
		return nil, err
	}

	// AttributeServiceを初期化
	cs.AttributeService = NewAttributeService[*models.Company](cs.FileService, ".detail.yaml")

	return cs, nil
}

// GetCompany は指定されたパスから会社を取得する
func (cs *CompanyService) GetCompany(folderName string) (models.Company, error) {
	// 会社データモデルを作成
	company, err := models.NewCompany(folderName)
	if err != nil {
		return models.Company{}, fmt.Errorf("会社データモデルの作成に失敗しました: %v", err)
	}

	// 属性ファイルと同期
	err = cs.SyncAttributeFile(&company)
	if err != nil {
		return models.Company{}, fmt.Errorf("属性ファイルとの同期に失敗しました: %v", err)
	}

	return company, nil
}

// SyncAttributeFile は会社の属性データを読み込み、会社データモデルに反映する
func (cs *CompanyService) SyncAttributeFile(company *models.Company) error {
	// 会社の属性データを読み込む
	attribute, err := cs.AttributeService.Load(company)
	if err != nil {
		// ファイルが存在しない場合は新規ファイルとして非同期で保存
		go func() {
			cs.AttributeService.Save(company)
		}()
		return nil
	}

	// 属性ファイルにしか保持されない内容をcompanyに反映
	company.LongName = attribute.LongName
	company.PostalCode = attribute.PostalCode
	company.Address = attribute.Address
	company.Phone = attribute.Phone
	company.Email = attribute.Email
	company.Website = attribute.Website
	company.Tags = attribute.Tags

	// フォルダー名から解析した情報を属性に反映
	attribute.ShortName = company.ShortName
	attribute.BusinessType = company.BusinessType

	// データに変更がある場合のみ非同期で保存（パフォーマンス重視）
	if cs.hasCompanyChanged(*attribute, *company) {
		go func() {
			cs.AttributeService.Save(attribute)
		}()
	}

	return nil
}

// hasCompanyChanged 会社データに変更があるかチェック
func (cs *CompanyService) hasCompanyChanged(attribute, current models.Company) bool {
	return attribute.ShortName != current.ShortName ||
		attribute.BusinessType != current.BusinessType
}

// GetCompanies は会社一覧を取得する
func (cs *CompanyService) GetCompanies() []models.Company {
	// ファイルシステムから会社フォルダー一覧を取得
	fileInfos, err := cs.FileService.GetFileInfos()
	if err != nil {
		return []models.Company{}
	}

	companies := make([]models.Company, 0, len(fileInfos))
	for _, fileInfo := range fileInfos {
		// フォルダー名のみを使用
		company, err := models.NewCompany(fileInfo.Name)
		if err != nil {
			continue
		}

		// 属性ファイルと同期（エラーは無視）
		cs.SyncAttributeFile(&company)

		companies = append(companies, company)
	}

	return companies
}

// GetCompanyByID retrieves a company by its ID
func (cs *CompanyService) GetCompanyByID(id string) (*models.Company, error) {
	companies := cs.GetCompanies()

	for _, company := range companies {
		if company.ID == id {
			return &company, nil
		}
	}

	return nil, fmt.Errorf("company with ID %s not found", id)
}

// GetCompanyByName retrieves a company by its name (case-insensitive)
func (cs *CompanyService) GetCompanyByName(name string) (*models.Company, error) {
	companies := cs.GetCompanies()

	name = strings.TrimSpace(strings.ToLower(name))

	for _, company := range companies {
		if strings.ToLower(company.ShortName) == name ||
			strings.ToLower(company.LongName) == name {
			return &company, nil
		}
	}

	return nil, nil // Return nil if not found (not an error)
}

// Update は会社情報を更新し、必要に応じてフォルダー名を変更します
// oldFolderName: 変更前のフォルダー名
// company: 更新後の会社情報
func (cs *CompanyService) Update(oldFolderName string, company *models.Company) error {
	// 会社のショート名と業種から新しいフォルダー名を生成
	newFolderName := fmt.Sprintf("%s %s", company.BusinessType.Code(), company.ShortName)

	// 新しいフォルダー名と古いフォルダー名を比較
	if newFolderName != oldFolderName {
		// フォルダー名の変更
		err := cs.FileService.MoveFile(oldFolderName, newFolderName)
		if err != nil {
			return err
		}

		// FolderNameを更新
		company.FolderName = newFolderName
	}

	// UpdateFolderNameを呼び出してFolderNameを確実に更新
	company.UpdateFolderName()

	// 計算が必要な項目の更新
	company.ID = models.NewIDFromString(company.ShortName).Len5()

	// 更新後の会社情報を属性ファイルに反映
	return cs.AttributeService.Save(company)
}
