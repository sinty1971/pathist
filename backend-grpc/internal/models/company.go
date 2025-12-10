package models

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	grpcv1 "backend-grpc/gen/grpc/v1"
	"backend-grpc/internal/core"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
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
	m.SetId(GenerateCompanyId(shortName))
	m.SetManagedFolder(managedFolder)
	m.SetCategoryIndex(int32(catIndex))
	m.SetShortName(shortName)

	return nil
}

// Update は会社情報を更新します
// 必要に応じて管理フォルダー名の変更も行います
func (m *Company) Update(updatedCompany *Company) error {

	// 引数チェック
	if updatedCompany == nil {
		return errors.New("updatedCompany is nil")
	}

	// ファイル名変更の必要がある場合は管理フォルダー名を更新
	updatedCompany.UpdateManagedFolderWithParams()
	if m.GetManagedFolder() != updatedCompany.GetManagedFolder() {

		// フォルダー名変更
		if err := os.Rename(m.GetManagedFolder(), updatedCompany.GetManagedFolder()); err != nil {
			return err
		}
	}

	// Inside情報の更新
	m.UpdateInsideParams(updatedCompany)

	return nil
}

// CreateManagedFolder はパラメータをもとに管理フォルダー名変更します
func (m *Company) UpdateManagedFolderWithParams() {
	base := filepath.Dir(m.GetManagedFolder())
	categoryIndex := strconv.Itoa(int(m.GetCategoryIndex()))
	folderName := categoryIndex + " " + m.GetShortName()
	managedFolder := filepath.Join(base, folderName)
	m.SetManagedFolder(managedFolder)
}

func (m *Company) UpdateInsideParams(updatedCompany *Company) {
	// Inside情報の更新
	m.SetInsideAddress(updatedCompany.GetInsideAddress())
	m.SetInsideWebsite(updatedCompany.GetInsideWebsite())
	m.SetInsideEmail(updatedCompany.GetInsideEmail())
	m.SetInsideTel(updatedCompany.GetInsideTel())
	m.SetInsideLegalName(updatedCompany.GetInsideLegalName())

	m.SetShortName(updatedCompany.GetShortName())
}

// Persiser インターフェースの実装

// PersistBytes は永続化用のメッセージを取得します
func (m *Company) GetPersistBytes() ([]byte, error) {
	return proto.Marshal(m.Company)
}

// SetPersistBytes は永続化用のバイトデータを設定します
func (m *Company) SetPersistBytes(b []byte) error {
	company := &grpcv1.Company{}
	if err := proto.Unmarshal(b, company); err != nil {
		return err
	}

	// protobuf リフレクションを使ってフィールドにアクセス
	msg := company.ProtoReflect()
	fields := msg.Descriptor().Fields()

	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		fieldName := field.Name() // フィールド名（例: "id", "short_name"）

		// Getter: 値を取得
		value := msg.Get(field)

		// フィールドの型に応じて処理
		switch field.Kind() {
		case protoreflect.StringKind:
			strValue := value.String()
			_ = strValue // 使用例
		case protoreflect.Int32Kind:
			int32Value := int32(value.Int())
			_ = int32Value // 使用例
			// 必要に応じて他の型も処理

		}

		// Setter: 値を設定（例）
		// msg.Set(field, protoreflect.ValueOfString("new value"))

		_ = fieldName // 使用例
	}

	m.Company = company
	return nil
}

// GenerateCompanyId は会社の短縮名から一意の会社IDを生成します
func GenerateCompanyId(shortName string) string {
	return core.GenerateIdFromString(shortName)
}
