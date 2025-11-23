package models

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	grpcv1 "backend-grpc/gen/grpc/v1"
	"backend-grpc/internal/ext"
)

// Company は gRPC grpc.v1.Company メッセージの拡張版です。
type Company struct {
	// Company メッセージ本体
	*grpcv1.Company

	// PersistFilename は永続化サービス用のファイル名
	PersistFilename string
}

// GetPersistPath は永続化ファイルのパスを取得します
// Persistable インターフェースの実装
func (c *Company) GetPersistPath() string {
	return filepath.Join(c.GetManagedFolder(), c.PersistFilename)
}

// GetPersistInfo は永続化対象のオブジェクトを取得します
// Persistable インターフェースの実装
func (c *Company) GetPersistInfo() any {
	// Companyモデル自体を返すことで、MarshalYAMLメソッドが使われる
	return c
}

// SetPersistInfo は永続化対象のオブジェクトを設定します
// Persistable インターフェースの実装
func (c *Company) SetPersistInfo(obj any) {
	// Companyモデルが渡される場合
	if company, ok := obj.(*Company); ok {
		c.Company = company.Company
		c.PersistFilename = company.PersistFilename
		return
	}

	// grpcv1.Companyが直接渡される場合（後方互換性）
	if company, ok := obj.(*grpcv1.Company); ok {
		c.Company = company
	}
}

// NewCompany 会社フォルダーパス名からCompanyを作成します
func NewCompany(managedFolder string) (*Company, error) {

	// フォルダー名を解析
	company, err := parseCompany(managedFolder)
	if err != nil {
		return nil, err
	}

	// Companyインスタンスを初期化
	company.SetId(GenerateCompanyId(company.GetShortName()))
	company.SetManagedFolder(managedFolder)
	company.SetInsideIdealPath("")
	company.SetInsideLegalName(company.GetShortName())
	company.SetInsideRequiredFiles([]*grpcv1.FileInfo{})
	company.SetInsidePostalCode("")
	company.SetInsideAddress("")
	company.SetInsidePhone("")
	company.SetInsideEmail("")
	company.SetInsideWebsite("")
	company.PersistFilename = ext.ConfigMap["CompanyPersistFilename"]

	return &company, nil
}

// GenerateCompanyId は会社の短縮名から一意の会社IDを生成します
func GenerateCompanyId(shortName string) string {
	return ext.GenerateIdFromString(shortName)
}

// parseCompany は"[0-9] [会社名]"形式のファイル名を解析します
// 会社名内にハイフンが含まれている場合は関連会社名または業務内容として処理します
// 戻り値Companyは: Cateory, ShortName, Tags のみ設定されます
func parseCompany(managedFolder string) (Company, error) {

	// フォルダー名を取得
	var folderName string
	if lastSlash := strings.LastIndex(managedFolder, "/"); lastSlash != -1 {
		folderName = managedFolder[lastSlash+1:]
	}

	company := Company{
		Company: grpcv1.Company_builder{}.Build(),
	}

	// 最小長チェック（"0 X"の最小3文字）
	if len(folderName) < 3 {
		return company, errors.New("ファイル名が短すぎます")
	}

	// 2番目の文字がスペースかチェック
	if folderName[1] != ' ' {
		return company, errors.New("2番目の文字がスペースではありません")
	}

	// 最初の文字がCompanyCategoryCodeかチェック
	i, err := strconv.Atoi(string(folderName[0]))
	if err != nil {
		return company, err
	}
	idx := CompanyCategoryIndex(i)
	if !(&idx).IsValid() {
		return company, errors.New("最初の文字が業種コードではありません")
	}
	// 業種を取得
	company.SetCategoryIndex(int32(idx))

	// 会社名部分を取得
	companyPart := folderName[2:]
	if companyPart == "" {
		return company, errors.New("会社名が空です")
	}

	// 会社名内にハイフンが含まれているかチェック
	var relatedInfo string
	if hyphenIndex := strings.Index(companyPart, "-"); hyphenIndex != -1 {
		// ハイフン以降の文字列を取得
		relatedInfo = companyPart[hyphenIndex+1:]
		company.SetShortName(companyPart[:hyphenIndex])
	} else {
		// ハイフンがない場合は全体を会社名とする
		company.SetShortName(companyPart)
	}

	// タグを作成
	category := CompanyCategoryMap[idx]
	tags := []string{category, company.GetShortName()}
	company.SetInsideTags(tags)
	if relatedInfo != "" {
		company.SetInsideTags(append(company.GetInsideTags(), relatedInfo))
	}

	return company, nil
}

// AddTag は会社のタグリストにタグを追加します
func (c *Company) AddTag(tag string) {
	if c == nil {
		return
	}
	if slices.Contains(c.GetInsideTags(), tag) {
		return // タグは既に存在します
	}
	c.SetInsideTags(append(c.GetInsideTags(), tag))
}

// Update は会社情報を更新します
func (c *Company) Update(updatedCompany *Company) (*Company, error) {
	if c == nil || updatedCompany == nil {
		return nil, errors.New("company or updatedCompany is nil")
	}

	// 管理フォルダーは変更しない
	updatedCompany.SetManagedFolder(c.GetManagedFolder())

	// 永続化サービスの設定を引き継ぐ
	updatedCompany.PersistFilename = c.PersistFilename

	return updatedCompany, nil
}

// MarshalYAML は YAML 用のシリアライズを行います。
func (c Company) MarshalYAML() (any, error) {
	// Getterを使って明示的にマップ化
	return c.MarshalMap(), nil
}

// UnmarshalYAML は YAML からの復元を行います。
func (c *Company) UnmarshalYAML(unmarshal func(any) error) error {
	var m map[string]any
	if err := unmarshal(&m); err != nil {
		return err
	}

	if err := c.UnmarshalMap(m); err != nil {
		return err
	}
	return nil
}

// MarshalJSON は JSON 用のシリアライズを行います。
func (c Company) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.MarshalMap())
}

// UnmarshalJSON は JSON からの復元を行います。
func (c *Company) UnmarshalJSON(data []byte) error {
	var m map[string]any

	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	if err := c.UnmarshalMap(m); err != nil {
		return err
	}

	return nil
}

func (c *Company) MarshalMap() map[string]any {
	return map[string]any{
		"inside_ideal_path":     c.GetInsideIdealPath(),
		"inside_legal_name":     c.GetInsideLegalName(),
		"inside_postal_code":    c.GetInsidePostalCode(),
		"inside_address":        c.GetInsideAddress(),
		"inside_phone":          c.GetInsidePhone(),
		"inside_email":          c.GetInsideEmail(),
		"inside_website":        c.GetInsideWebsite(),
		"inside_tags":           c.GetInsideTags(),
		"inside_required_files": c.GetInsideRequiredFiles(),
	}
}

func (c *Company) UnmarshalMap(m map[string]any) error {
	if v, ok := m["inside_ideal_path"].(string); ok {
		c.SetInsideIdealPath(v)
	}
	if v, ok := m["inside_legal_name"].(string); ok {
		c.SetInsideLegalName(v)
	}
	if v, ok := m["inside_postal_code"].(string); ok {
		c.SetInsidePostalCode(v)
	}
	if v, ok := m["inside_address"].(string); ok {
		c.SetInsideAddress(v)
	}
	if v, ok := m["inside_phone"].(string); ok {
		c.SetInsidePhone(v)
	}
	if v, ok := m["inside_email"].(string); ok {
		c.SetInsideEmail(v)
	}
	if v, ok := m["inside_website"].(string); ok {
		c.SetInsideWebsite(v)
	}
	if v, ok := m["inside_tags"].([]any); ok {
		tags := make([]string, 0, len(v))
		for _, tag := range v {
			if tagStr, ok := tag.(string); ok {
				tags = append(tags, tagStr)
			}
		}
		c.SetInsideTags(tags)
	}
	if v, ok := m["inside_required_files"].([]any); ok {
		files := make([]*grpcv1.FileInfo, 0, len(v))
		for _, file := range v {
			if fileMap, ok := file.(map[string]any); ok {
				fileInfo := &grpcv1.FileInfo{}
				fileData, _ := json.Marshal(fileMap)
				json.Unmarshal(fileData, fileInfo)
				files = append(files, fileInfo)
			}
		}
		c.SetInsideRequiredFiles(files)
	}

	return nil
}
