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

	// PersistFilename は永続化サービス用のファイル名
	PersistFilename string
}

// NewCompany インスタンス作成と初期化を行います
func NewCompany() *Company {

	// インスタンス作成と初期化
	return &Company{
		Company:         grpcv1.Company_builder{}.Build(),
		PersistFilename: core.ConfigMap["CompanyPersistFilename"],
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
	categoryIndex := CompanyCategoryIndex(obj.GetCategoryIndex())
	folderName := strconv.Itoa(int(categoryIndex)) + " " + obj.GetShortName()
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

	// CompanyCategoryIndexの取得
	number, err := strconv.Atoi(string(folderName[0]))
	if err != nil {
		return err
	}
	categoryIndex := CompanyCategoryIndex(number)
	if err := categoryIndex.Error(); err != nil {
		return err
	}

	// 会社名と関連名の取得
	var companyName string
	var relatedName string

	// 会社フォルダー名の解析
	companyParts := strings.Split(folderName[2:], " ")
	if len(companyParts) == 0 || companyParts[0] == "" {
		return errors.New("会社名が取得できません")
	}
	companyPart0 := companyParts[0]
	if idx := strings.Index(companyPart0, "-"); idx > -1 {
		// ハイフン以前の文字列を会社名とする
		companyName = companyPart0[:idx]
		// ハイフン以降の文字列を関連文字列とする
		relatedName = companyPart0[idx+1:]
	} else {
		companyName = companyPart0
	}

	// ID,ManagedFolder,Category,ShortNameの設定
	obj.SetId(GenerateCompanyId(companyName))
	obj.SetManagedFolder(managedFolder)
	obj.SetCategoryIndex(categoryIndex.ToInt32())
	obj.SetShortName(companyName)

	// タグの設定
	core.ModelAddInsideTags(obj, []string{
		companyName,
		CompanyCategoryMap[categoryIndex],
		relatedName}...,
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

// GetPersistData は永続化対象のオブジェクトを取得します
// Persistable インターフェースの実装
func (obj *Company) GetPersistData() (map[string]any, error) {
	return map[string]any{
		"inside_ideal_path":     obj.GetInsideIdealPath(),
		"inside_legal_name":     obj.GetInsideLegalName(),
		"inside_postal_code":    obj.GetInsidePostalCode(),
		"inside_address":        obj.GetInsideAddress(),
		"inside_phone":          obj.GetInsidePhone(),
		"inside_email":          obj.GetInsideEmail(),
		"inside_website":        obj.GetInsideWebsite(),
		"inside_tags":           obj.GetInsideTags(),
		"inside_required_files": obj.GetInsideRequiredFiles(),
	}, nil
}

// SetPersistData は永続化対象のオブジェクトを設定します
// Persistable インターフェースの実装
func (obj *Company) SetPersistData(persistData map[string]any) error {
	// 文字列フィールドのセッターのマップ
	setterMap := map[string]func(string){
		"inside_ideal_path":  obj.SetInsideIdealPath,
		"inside_legal_name":  obj.SetInsideLegalName,
		"inside_postal_code": obj.SetInsidePostalCode,
		"inside_address":     obj.SetInsideAddress,
		"inside_phone":       obj.SetInsidePhone,
		"inside_email":       obj.SetInsideEmail,
		"inside_website":     obj.SetInsideWebsite,
		"inside_tags":        obj.setRawInsideTags,
	}

	// 配列フィールドの処理
	if v, ok := persistData["inside_required_files"].([]any); ok {
		files := make([]*grpcv1.FileInfo, 0, len(v))
		for _, file := range v {
			if fileMap, ok := file.(map[string]any); ok {
				fileInfo := grpcv1.FileInfo_builder{}.Build()
				fileData, _ := json.Marshal(fileMap)
				json.Unmarshal(fileData, fileInfo)
				files = append(files, fileInfo)
			}
		}
		obj.SetInsideRequiredFiles(files)
	}

	// デフォルトの文字列フィールド設定処理を呼び出し
	return core.DefaultSetPersistData(persistData, setterMap)
}

func (obj *Company) setRawInsideTags(raw string) {
	var tags []string
	if err := json.Unmarshal([]byte(raw), &tags); err != nil {
		return
	}
	core.ModelAddInsideTags(obj, tags...)
}

// MarshalJSON は JSON 用のシリアライズを行います。
func (obj Company) MarshalJSON() ([]byte, error) {
	return core.DefaultMarshalJSON(obj.GetPersistData)
}

// UnmarshalYAML は YAML からの復元を行います。
func (obj *Company) UnmarshalYAML(unmarshal func(any) error) error {
	return core.DefaultUnmarshalYAML(unmarshal, obj.SetPersistData)
}

// MarshalYAML は YAML 用のシリアライズを行います。
func (obj Company) MarshalYAML() (any, error) {
	return core.DefaultMarshalYAML(obj.GetPersistData)
}

// UnmarshalJSON は JSON からの復元を行います。
func (obj *Company) UnmarshalJSON(json []byte) error {
	return core.DefaultUnmarshalJSON(json, obj.SetPersistData)
}
