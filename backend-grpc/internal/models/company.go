package models

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	grpcv1 "backend-grpc/gen/grpc/v1"
	"backend-grpc/internal/core"

	"google.golang.org/protobuf/encoding/protojson"
)

// Company は gRPC grpc.v1.Company メッセージの拡張版です。
type Company struct {
	// Company メッセージ本体
	*grpcv1.Company

	// Persist 永続化用フィールド
	*core.Persist
}

// NewCompany インスタンス作成と初期化を行います
func NewCompany() *Company {

	// インスタンス作成と初期化
	company := &Company{}
	company.Company = grpcv1.Company_builder{}.Build()
	company.Persist = core.NewPersister(company, core.ConfigMap["CompanyPersistFilename"])

	return company
}

// ParseFromManagedFolder は"[0-9] [会社名]"形式のファイル名となっているパスを解析します
// 会社名内のハイフン（含まれる場合）以前の文字列を会社名、ハイフン以降の文字列を関連名として扱います
// 戻り値Companyは: Id, ManagedFolder, Cateory, ShortName, Tags のみ設定されます
func (m *Company) ParseFromManagedFolder(managedFolders ...string) error {

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
	nameParts := strings.Split(folderName[2:], " ")
	if len(nameParts) == 0 || nameParts[0] == "" {
		return errors.New("会社名が取得できません")
	}

	// 会社名の解析（ハイフンで分割）
	shortName := nameParts[0]
	if idx := strings.Index(nameParts[0], "-"); idx > -1 {
		// ハイフン以前の文字列を会社名とする
		shortName = nameParts[0][:idx]
		// ハイフン以降の文字列を関連文字列とする
	}

	// ID,ManagedFolder,Category,ShortNameの設定
	m.SetManagedFolder(managedFolder)
	m.SetCategoryIndex(int32(catIndex))
	m.SetShortName(shortName)

	// 会社IDの生成と設定
	id := m.GenerateId(shortName)
	m.SetId(id)

	return nil
}

// Update は会社情報を更新します
// 必要に応じて管理フォルダー名の変更も行います
func (m *Company) Update(newCompany *Company) error {

	// 引数チェック
	if newCompany == nil {
		return errors.New("updatedCompany is nil")
	}

	// 新しいパラメータを元に管理フォルダーパスを生成
	newManagedFolder := newCompany.GenerateManagedFolder(
		filepath.Dir(m.GetManagedFolder()),
		newCompany.GetCategoryIndex(),
		newCompany.GetShortName(),
	)

	// ファイル名変更の必要がある場合は管理フォルダー名を更新
	if newManagedFolder != m.GetManagedFolder() {

		// フォルダー名変更
		if err := os.Rename(m.GetManagedFolder(), newManagedFolder); err != nil {
			return err
		}
	}

	// Persist情報の更新
	m.UpdatePersists(newCompany.Persist)

	return nil
}

// GenerateId は会社の短縮名から一意の会社IDを生成します
func (m *Company) GenerateId(shortName string) string {
	return core.GenerateIdFromString(shortName)
}

// GenerateManagedFolder はパラメータをもとに管理フォルダー名変更します
// base: 基本パス(原則として　O:/.../1 会社 などの親フォルダー)
// idx: カテゴリーインデックス
// name: 省略会社名
func (m *Company) GenerateManagedFolder(base string, idx int32, name string) string {
	folderName := strconv.Itoa(int(idx)) + " " + name
	return filepath.Join(base, folderName)
}

// Persiser インターフェースの実装
// GetPersistFolder は永続化フォルダーのパスを取得します
func (m *Company) GetPersistFolder() string {
	return m.Company.GetManagedFolder()
}

// MarshalJSON()  は永続化用のメッセージのJSONを取得します
// JSON Marshaler インターフェースの実装
func (m *Company) MarshalJSON() ([]byte, error) {
	opts := protojson.MarshalOptions{
		Multiline: true,
		Indent:    "  ",
	}
	return opts.Marshal(m.Company)
}

// UnmarshalJSON() は永続化用のメッセージのJSONを設定します
// JSON Unmarshaler インターフェースの実装
func (m *Company) UnmarshalJSON(b []byte) error {
	opts := protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}
	return opts.Unmarshal(b, m.Company)
}
