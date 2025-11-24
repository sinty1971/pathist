package models

import (
	"encoding/json"
	"errors"
	"log"
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
func (h *Company) GetPersistPath() string {
	return filepath.Join(h.GetManagedFolder(), h.PersistFilename)
}

// GenerateCompanyId は会社の短縮名から一意の会社IDを生成します
func GenerateCompanyId(shortName string) string {
	return ext.GenerateIdFromString(shortName)
}

// ParseFromManagedFolder は"[0-9] [会社名]"形式のファイル名となっているパスを解析します
// 会社名内のハイフン（含まれる場合）以前の文字列を会社名、ハイフン以降の文字列を関連名として扱います
// 戻り値Companyは: Id, ManagedFolder, Cateory, ShortName, Tags のみ設定されます
func (h *Company) ParseFromManagedFolder(managedFolder string) error {

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
		h.SetManagedFolder(managedFolder)
	}

	// CompanyCategoryIndexの取得
	var categoryIndex CompanyCategoryIndex
	if number, err := strconv.Atoi(string(folderName[0])); err != nil {
		return err
	} else if categoryIndex = CompanyCategoryIndex(number); categoryIndex.Error() != nil {
		return categoryIndex.Error()
	} else {
		h.SetCategoryIndex(categoryIndex.ToInt32())
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
	h.SetId(GenerateCompanyId(companyName))
	h.SetShortName(companyName)

	// タグの設定
	h.AddInsideTags([]string{
		companyName,
		CompanyCategoryMap[categoryIndex],
		relatedName}...,
	)

	return nil
}

// AddInsideTags は会社のタグリストにタグを追加します
func (h *Company) AddInsideTags(tags ...string) {
	insideTags := h.GetInsideTags()
	for _, tag := range tags {
		if tag != "" && !slices.Contains(insideTags, tag) {
			insideTags = append(insideTags, tag)
		}
	}
	h.SetInsideTags(tags)
}

// Update は会社情報を更新します
func (h *Company) Update(updatedCompany *Company) (*Company, error) {
	if h == nil || updatedCompany == nil {
		return nil, errors.New("company or updatedCompany is nil")
	}

	// 管理フォルダーは変更しない
	updatedCompany.SetManagedFolder(h.GetManagedFolder())

	// 永続化サービスの設定を引き継ぐ
	updatedCompany.PersistFilename = h.PersistFilename

	return updatedCompany, nil
}

// GetPersistData は永続化対象のオブジェクトを取得します
// Persistable インターフェースの実装
func (h *Company) GetPersistData() map[string]any {
	return map[string]any{
		"inside_ideal_path":     h.GetInsideIdealPath(),
		"inside_legal_name":     h.GetInsideLegalName(),
		"inside_postal_code":    h.GetInsidePostalCode(),
		"inside_address":        h.GetInsideAddress(),
		"inside_phone":          h.GetInsidePhone(),
		"inside_email":          h.GetInsideEmail(),
		"inside_website":        h.GetInsideWebsite(),
		"inside_tags":           h.GetInsideTags(),
		"inside_required_files": h.GetInsideRequiredFiles(),
	}
}

// MarshalJSON は JSON 用のシリアライズを行います。
func (h Company) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.GetPersistData())
}

// MarshalYAML は YAML 用のシリアライズを行います。
func (h Company) MarshalYAML() (any, error) {
	// Getterを使って明示的にマップ化
	return h.GetPersistData(), nil
}

// SetPersistData は永続化対象のオブジェクトを設定します
// Persistable インターフェースの実装
func (h *Company) SetPersistData(persistData map[string]any) {
	// 文字列フィールドの一括処理
	stringFields := map[string]func(string){
		"inside_ideal_path":  h.SetInsideIdealPath,
		"inside_legal_name":  h.SetInsideLegalName,
		"inside_postal_code": h.SetInsidePostalCode,
		"inside_address":     h.SetInsideAddress,
		"inside_phone":       h.SetInsidePhone,
		"inside_email":       h.SetInsideEmail,
		"inside_website":     h.SetInsideWebsite,
	}
	for key, setter := range stringFields {
		if v, ok := persistData[key].(string); ok {
			setter(v)
		} else {
			log.Printf("%v の項目がファイル %s にありません", v, h.PersistFilename)
		}
	}

	if tags, ok := persistData["inside_tags"].([]string); ok {
		h.SetInsideTags(tags)
	}
	if v, ok := persistData["inside_required_files"].([]any); ok {
		files := make([]*grpcv1.FileInfo, 0, len(v))
		for _, file := range v {
			if fileMap, ok := file.(map[string]any); ok {
				fileInfo := &grpcv1.FileInfo{}
				fileData, _ := json.Marshal(fileMap)
				json.Unmarshal(fileData, fileInfo)
				files = append(files, fileInfo)
			}
		}
		h.SetInsideRequiredFiles(files)
	}
}

// UnmarshalYAML は YAML からの復元を行います。
func (h *Company) UnmarshalYAML(unmarshal func(any) error) error {
	var m map[string]any
	if err := unmarshal(&m); err != nil {
		return err
	}

	h.SetPersistData(m)
	return nil
}

// UnmarshalJSON は JSON からの復元を行います。
func (h *Company) UnmarshalJSON(data []byte) error {
	var m map[string]any

	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	h.SetPersistData(m)

	return nil
}
