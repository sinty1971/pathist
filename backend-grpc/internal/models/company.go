package models

import (
	"encoding/json"
	"errors"
	"os"
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

// NewCompany インスタンス作成と初期化を行います
func NewCompany() *Company {

	// インスタンス作成と初期化
	return &Company{
		Company:         grpcv1.Company_builder{}.Build(),
		PersistFilename: ext.ConfigMap["CompanyPersistFilename"],
	}
}

// GetPersistPath は永続化ファイルのパスを取得します
// Persistable インターフェースの実装
func (obj *Company) GetPersistPath() string {
	return filepath.Join(obj.GetManagedFolder(), obj.PersistFilename)
}

// GenerateCompanyId は会社の短縮名から一意の会社IDを生成します
func GenerateCompanyId(shortName string) string {
	return ext.GenerateIdFromString(shortName)
}

// ParseFromManagedFolder は"[0-9] [会社名]"形式のファイル名となっているパスを解析します
// 会社名内のハイフン（含まれる場合）以前の文字列を会社名、ハイフン以降の文字列を関連名として扱います
// 戻り値Companyは: Id, ManagedFolder, Cateory, ShortName, Tags のみ設定されます
func (obj *Company) ParseFromManagedFolder(managedFolders ...string) error {

	// パスを結合
	managedFolder := filepath.Join(managedFolders...)

	// 引数 managedFolder のチェックとフォルダー名を取得
	var folderName string
	folderParts := strings.Split(managedFolder, string(os.PathSeparator))
	if idx := len(folderParts) - 1; idx < 0 {
		return errors.New("managedFolderの値が無効です")
	} else if folderName = folderParts[idx]; len(folderName) < 3 {
		return errors.New("managedFolderのファイル名形式が無効です（長さが短い）")
	} else if folderName[1] != ' ' {
		// 2番目の文字がスペースかチェック
		return errors.New("managedFolderのファイル名形式が無効です")
	} else {
		obj.SetManagedFolder(managedFolder)
	}

	// CompanyCategoryIndexの取得
	var categoryIndex CompanyCategoryIndex
	if number, err := strconv.Atoi(string(folderName[0])); err != nil {
		return err
	} else if categoryIndex = CompanyCategoryIndex(number); categoryIndex.Error() != nil {
		return categoryIndex.Error()
	} else {
		obj.SetCategoryIndex(categoryIndex.ToInt32())
	}

	// 会社名と関連名の取得
	var companyName string
	var relatedName string

	// 会社フォルダー名の解析
	if companyPart := strings.Split(folderName[2:], " ")[0]; companyPart == "" {
		return errors.New("会社名が取得できません")
	} else if idx := strings.Index(companyPart, "-"); idx > -1 {
		// ハイフン以前の文字列を会社名とする
		companyName = companyPart[:idx]
		// ハイフン以降の文字列を関連文字列とする
		relatedName = companyPart[idx+1:]
	} else {
		companyName = companyPart
	}

	// 会社IDと会社名の設定
	obj.SetId(GenerateCompanyId(companyName))
	obj.SetShortName(companyName)

	// タグの設定
	obj.AddInsideTags([]string{
		companyName,
		CompanyCategoryMap[categoryIndex],
		relatedName}...,
	)

	return nil
}

// AddInsideTags は会社のタグリストにタグを追加します
func (obj *Company) AddInsideTags(tags ...string) {
	insideTags := obj.GetInsideTags()
	for _, tag := range tags {
		if tag == "" || slices.Contains(insideTags, tag) {
			continue
		}
		insideTags = append(insideTags, tag)
	}
	obj.SetInsideTags(insideTags)
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

// LoadPersistData は永続化対象のオブジェクトを読み込みます
func (obj *Company) LoadPersistData() error {
	return ext.DefaultLoadPersitData(obj)
}

// SavePersistData は永続化対象のオブジェクトを保存します
func (obj *Company) SavePersistData() error {
	return ext.DefaultSavePersistData(obj)
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
		"inside_tags":        obj.altSetInsideTags,
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
	return ext.DefaultSetPersistData(persistData, setterMap)
}

func (obj *Company) altSetInsideTags(raw string) {
	var tags []string
	if err := json.Unmarshal([]byte(raw), &tags); err != nil {
		return
	}
	obj.AddInsideTags(tags...)
}

// MarshalJSON は JSON 用のシリアライズを行います。
func (obj Company) MarshalJSON() ([]byte, error) {
	return ext.DefaultMarshalJSON(obj.GetPersistData)
}

// UnmarshalYAML は YAML からの復元を行います。
func (obj *Company) UnmarshalYAML(unmarshal func(any) error) error {
	return ext.DefaultUnmarshalYAML(unmarshal, obj.SetPersistData)
}

// MarshalYAML は YAML 用のシリアライズを行います。
func (obj Company) MarshalYAML() (any, error) {
	return ext.DefaultMarshalYAML(obj.GetPersistData)
}

// UnmarshalJSON は JSON からの復元を行います。
func (obj *Company) UnmarshalJSON(json []byte) error {
	return ext.DefaultUnmarshalJSON(json, obj.SetPersistData)
}
