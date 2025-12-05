package models

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"strconv"
	"strings"

	grpcv1 "backend-grpc/gen/grpc/v1"
	"backend-grpc/internal/core"
)

// Company は gRPC grpc.v1.Company メッセージの拡張版です。
type Company struct {
	// Company メッセージ本体
	*grpcv1.Company

	// Persist 永続化用フィールド
	*Persist
}

// NewCompany インスタンス作成と初期化を行います
func NewCompany() *Company {

	// インスタンス作成と初期化
	return &Company{
		Company: grpcv1.Company_builder{}.Build(),
	}
}

// GenerateCompanyId は会社の短縮名から一意の会社IDを生成します
func GenerateCompanyId(shortName string) string {
	return core.GenerateIdFromString(shortName)
}

// NeedsRenameManagedFolder ファイル名変更の必要があるかチェック
func (obj *Company) NeedsRenameManagedFolder(newCompany *Company) bool {
	return (obj.GetShortName() != newCompany.GetShortName()) ||
		(obj.GetManagedFolder() != newCompany.GetManagedFolder()) ||
		(obj.GetCategoryIndex() != newCompany.GetCategoryIndex())
}

// CreateManagedFolder は会社の管理フォルダー名をパラメーターから生成します
func (obj *Company) CreateManagedFolder() string {
	base := filepath.Dir(obj.GetManagedFolder())
	categoryIndex := strconv.Itoa(int(obj.GetCategoryIndex()))
	folderName := categoryIndex + " " + obj.GetShortName()
	return filepath.Join(base, folderName)
}

// ParseFromManagedFolder は"[0-9] [会社名]"形式のファイル名となっているパスを解析します
// 会社名内のハイフン（含まれる場合）以前の文字列を会社名、ハイフン以降の文字列を関連名として扱います
// 戻り値Companyは: Id, ManagedFolder, Cateory, ShortName, Tags のみ設定されます
func (obj *Company) ParseFromManagedFolder(managedFolders ...string) error {

	// パスを結合
	managedFolder := filepath.Join(managedFolders...)

	// 引数 managedFolder からフォルダー名取得とチェック
	// "[0-9] [会社名]"の解析
	folderName := filepath.Base(managedFolder)
	if len(folderName) < 3 {
		return errors.New("managedFolderのファイル名形式が無効です（長さが短い）")
	} else if folderName[1] != ' ' {
		// 2番目の文字がスペースかチェック
		return errors.New("managedFolderのファイル名形式が無効です")
	}

	// カテゴリー情報の取得
	var catIndex int
	var err error

	if catIndex, err = strconv.Atoi(string(folderName[0])); err != nil {
		return err
	}
	if err := ErrorCompanyCategoryIndex(catIndex); err != nil {
		return err
	}

	// 会社フォルダー名の解析
	companyParts := strings.Split(folderName[2:], " ")
	if len(companyParts) == 0 || companyParts[0] == "" {
		return errors.New("会社名が取得できません")
	}

	// 会社名の解析（ハイフンで分割）
	shortName := companyParts[0]
	var relationTag string
	if idx := strings.Index(shortName, "-"); idx > -1 {
		// ハイフン以前の文字列を会社名とする
		shortName = shortName[:idx]
		// ハイフン以降の文字列を関連文字列とする
		relationTag = shortName[idx+1:]
	}

	// ID,ManagedFolder,Category,ShortNameの設定
	obj.SetId(GenerateCompanyId(shortName))
	obj.SetManagedFolder(managedFolder)
	obj.SetCategoryIndex(int32(catIndex))
	obj.SetShortName(shortName)

	// タグの設定
	AddInsideTags(obj, []string{
		shortName,
		CompanyCategoryMap[catIndex],
		relationTag}...,
	)

	return nil
}

// Update は会社情報を更新します
func (obj *Company) Update(updatedCompany *Company) (*Company, error) {
	if updatedCompany == nil {
		return nil, errors.New("updatedCompany is nil")
	}

	// 管理フォルダーは変更しない
	updatedCompany.SetManagedFolder(obj.GetManagedFolder())

	// 永続化サービスの設定を引き継ぐ
	updatedCompany.PersistFilename = obj.PersistFilename

	return updatedCompany, nil
}

// Persistable インターフェースの実装
//

// GetPersistPath は永続化ファイルのパスを取得します
func (obj *Company) GetPersistPath() string {
	return filepath.Join(obj.GetManagedFolder(), obj.PersistFilename)
}

// PersistMap は永続化用のフィールドマップを返します
func (obj *Company) PersistMap() map[string]PersistFunc {
	return map[string]PersistFunc{
		"inside_ideal_path":     PersistStringFunc(obj.GetInsideIdealPath, obj.SetInsideIdealPath),
		"inside_legal_name":     PersistStringFunc(obj.GetInsideLegalName, obj.SetInsideLegalName),
		"inside_postal_code":    PersistStringFunc(obj.GetInsidePostalCode, obj.SetInsidePostalCode),
		"inside_address":        PersistStringFunc(obj.GetInsideAddress, obj.SetInsideAddress),
		"inside_phone":          PersistStringFunc(obj.GetInsidePhone, obj.SetInsidePhone),
		"inside_email":          PersistStringFunc(obj.GetInsideEmail, obj.SetInsideEmail),
		"inside_website":        PersistStringFunc(obj.GetInsideWebsite, obj.SetInsideWebsite),
		"inside_tags":           PersistStringSliceFunc(obj.GetInsideTags, obj.SetInsideTags),
		"inside_required_files": PersistFileInfoSliceFunc(obj.GetInsideRequiredFiles, obj.SetInsideRequiredFiles),
	}
}

// MarshalJSON は Marshaler インターフェースの実装で、JSON 用のシリアライズを行います。
func (obj *Company) MarshalJSON() ([]byte, error) {
	val, err := Serialize(obj)
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(val)
}

// UnmarshalJSON は JSON からの復元を行います。
func (obj *Company) UnmarshalJSON(raw []byte) error {
	return core.UnmarshalJSON(obj, raw)
}

// MarshalYAML は YAML 用のシリアライズを行います。
func (obj Company) MarshalYAML() (any, error) {
	return core.MarshalYAML(&obj)
}

// UnmarshalYAML は YAML からの復元を行います。
func (obj *Company) UnmarshalYAML(unmarshal func(any) error) error {
	return core.UnmarshalYAML(obj, unmarshal)
}
