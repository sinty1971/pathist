package models

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	grpcv1 "backend-grpc/gen/grpc/v1"
	exts "backend-grpc/internal/extentions"
)

// Company は gRPC grpc.v1.Company メッセージの拡張版です。
type Company struct {
	// Company メッセージ本体
	*grpcv1.Company

	// persist はファイル永続化サービス用のヘルパー
	// FPS: File Persist Service
	persist exts.ObjectPersistService[*Company]
}

// GetPersistPath は永続化ファイルのパスを取得します
// Persistable インターフェースの実装
func (c *Company) GetPersistPath() string {
	return filepath.Join(c.GetManagedFolder(), c.persist.PersistFilename)
}

// GetObject は永続化対象のオブジェクトを取得します
// Persistable インターフェースの実装
func (c *Company) GetObject() any {
	return c.Company
}

// SetObject は永続化対象のオブジェクトを設定します
// Persistable インターフェースの実装
func (c *Company) SetObject(obj any) {
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

	return &company, nil
}

// GenerateCompanyId は会社の短縮名から一意の会社IDを生成します
func GenerateCompanyId(shortName string) string {
	return exts.GenerateIdFromString(shortName)
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
	updatedCompany.persist = c.persist

	return updatedCompany, nil
}

// MarshalYAML は YAML 用のシリアライズを行います。
func (c Company) MarshalYAML() (any, error) {
	return c.Company, nil
}

// UnmarshalYAML は YAML からの復元を行います。
func (c *Company) UnmarshalYAML(unmarshal func(any) error) error {
	enc := &grpcv1.Company{}
	if err := unmarshal(enc); err != nil {
		return err
	}
	c.Company = enc
	return nil
}

// MarshalJSON は JSON 用のシリアライズを行います。
func (c Company) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Company)
}

// UnmarshalJSON は JSON からの復元を行います。
func (c *Company) UnmarshalJSON(data []byte) error {
	enc := &grpcv1.Company{}
	if err := json.Unmarshal(data, enc); err != nil {
		return err
	}
	c.Company = enc
	return nil
}
