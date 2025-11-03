package models

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"slices"
	"strings"

	grpcv1 "backend-grpc/gen/grpc/v1"
	"backend-grpc/internal/persist"
	"backend-grpc/internal/utils"
)

// Company は gRPC grpc.v1.Company メッセージの拡張版です。
type Company struct {
	// Company メッセージ本体
	*grpcv1.Company

	// insideFPS はファイル永続化サービス用のヘルパー
	// FPS: File Persist Service
	insideFPS persist.FilePersistService[*Company]
}

// GetFilePersistPath は永続化ファイルのパスを取得します
// Persistable インターフェースの実装
func (c *Company) GetFilePersistPath() string {
	return filepath.Join(c.GetManagedFolder(), c.insideFPS.PersistFilename)
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
	// フォルダー名を取得
	var folderName string
	if lastSlash := strings.LastIndex(managedFolder, "/"); lastSlash != -1 {
		folderName = managedFolder[lastSlash+1:]
	}

	// ファイル名の[0-9] [会社名]の規則を解析
	result, err := ParseCompanyName(folderName)
	if err != nil {
		return nil, err
	}

	company := &Company{
		Company: grpcv1.Company_builder{
			Id:            GenerateCompanyId(result.ShortName),
			ManagedFolder: managedFolder,
			ShortName:     result.ShortName,
			CategoryIndex: int32(result.Category),

			InsideIdealPath:     "",
			InsideLegalName:     result.ShortName,
			InsideTags:          append([]string{"会社"}, result.Tags...),
			InsideRequiredFiles: []*grpcv1.FileInfo{},
			InsidePostalCode:    "",
			InsideAddress:       "",
			InsidePhone:         "",
			InsideEmail:         "",
			InsideWebsite:       "",
		}.Build(),
	}

	// IDを更新とリターン
	return company, nil
}

func GenerateCompanyId(shortName string) string {
	return utils.GenerateIdFromString(shortName)
}

// parseCompanyNameResult は ParseCompanyName の戻り値を表します
type parseCompanyNameResult struct {
	Category  CompanyCategoryIndex
	ShortName string
	Tags      []string
}

// ParseCompanyName は"[0-9] [会社名]"形式のファイル名を解析します
// 会社名内にハイフンが含まれている場合は関連会社名または業務内容として処理します
// 戻り値: [0]業種, [1]会社名, error
func ParseCompanyName(name string) (parseCompanyNameResult, error) {
	result := parseCompanyNameResult{}

	// 最小長チェック（"0 X"の最小3文字）
	if len(name) < 3 {
		return result, errors.New("ファイル名が短すぎます")
	}

	// 2番目の文字がスペースかチェック
	if name[1] != ' ' {
		return result, errors.New("2番目の文字がスペースではありません")
	}

	// 最初の文字がCompanyCategoryCodeかチェック
	code := CompanyCategoryIndex(name[0])
	if !(&code).IsValid() {
		return result, errors.New("最初の文字が業種コードではありません")
	}

	// 会社名部分を取得
	companyPart := name[2:]
	if companyPart == "" {
		return result, errors.New("会社名が空です")
	}

	// 業種を取得
	result.Category = code

	// 会社名内にハイフンが含まれているかチェック
	var relatedInfo string
	if hyphenIndex := strings.Index(companyPart, "-"); hyphenIndex != -1 {
		// ハイフン以降の文字列を取得
		relatedInfo = companyPart[hyphenIndex+1:]
		result.ShortName = companyPart[:hyphenIndex]
	} else {
		// ハイフンがない場合は全体を会社名とする
		result.ShortName = companyPart
	}

	// タグを作成
	result.Tags = []string{CompanyCategoryMap[result.Category], result.ShortName}
	if relatedInfo != "" {
		result.Tags = append(result.Tags, relatedInfo)
	}

	return result, nil
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

func (c *Company) Update(updatedCompany *Company) (*Company, error) {
	if c == nil || updatedCompany == nil {
		return nil, errors.New("company or updatedCompany is nil")
	}

	// 管理フォルダーは変更しない
	updatedCompany.SetManagedFolder(c.GetManagedFolder())

	// 永続化サービスの設定を引き継ぐ
	updatedCompany.insideFPS = c.insideFPS

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
