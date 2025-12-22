package models

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"google.golang.org/protobuf/proto"

	grpcv1 "server-grpc/gen/grpc/v1"
	"server-grpc/internal/core"
)

// Company は gRPC grpc.v1.Company メッセージの拡張版です。
type Company struct {
	// Company メッセージ本体
	*grpcv1.Company

	// Model Pathist 共通モデル
	Pathist *core.Pathist
}

// NewCompany インスタンス作成と初期化を行います
func NewCompany() *Company {

	// インスタンス作成と初期化
	company := &Company{}
	company.Company = grpcv1.Company_builder{}.Build()
	company.Pathist = core.NewPathist(company, core.ConfigMap["CompanyPersistFilename"])

	return company
}

func (m *Company) GetProtoMessage() proto.Message {
	if m == nil {
		return nil
	}
	return m.Company
}

// ParseFromTarget は"[0-9] [会社名]"形式のファイル名となっているパスを解析します
// 会社名内のハイフン（含まれる場合）以前の文字列を会社名、ハイフン以降の文字列を関連名として扱います
// 戻り値Companyは: Id, Target, Cateory, ShortName, Tags のみ設定されます
func (m *Company) ParseFrom(pathistFolder ...string) error {

	// パスを結合
	folder := filepath.Join(pathistFolder...)

	// 引数 target からフォルダー名取得とチェック
	// "[0-9] [会社名]"の解析
	folderName := filepath.Base(folder)
	if len(folderName) < 3 {
		return errors.New("targetのファイル名形式が無効です（長さが短い）")
	} else if folderName[1] != ' ' {
		// 2番目の文字がスペースかチェック
		return errors.New("targetaのファイル名形式が無効です")
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

	// Target,Category,ShortNameの設定
	m.SetPathistFolder(folder)
	m.SetCategoryIndex(int32(catIndex))
	m.SetShortName(shortName)

	// IDの設定、targetの設定が終了した後に実行
	id, err := m.Pathist.GenerateId()
	if err != nil {
		return err
	}
	m.SetId(id)
	return nil
}

// ImportFrom は会社情報を更新します
// 必要に応じて管理フォルダー名の変更も行います
func (m *Company) ImportFrom(src *Company) error {

	// 引数チェック
	if src == nil {
		return errors.New("src Company is nil")
	}

	// 新しいパラメータを元に管理フォルダーパスを生成
	newPathistFolder := GenerateCompanyPathistFolder(
		filepath.Dir(m.GetPathistFolder()),
		src.GetCategoryIndex(),
		src.GetShortName(),
	)

	// ファイル名変更の必要がある場合は管理フォルダー名を更新
	if newPathistFolder != m.GetPathistFolder() {

		// フォルダー名変更
		if err := os.Rename(m.GetPathistFolder(), newPathistFolder); err != nil {
			return err
		}
	}

	// Persist情報の更新
	return m.Pathist.ImportPersists(src.Pathist)
}

// GenerateCompanyTarget はパラメータをもとに管理フォルダー名変更します
// base: 基本パス(原則として　O:/.../1 会社 などの親フォルダー)
// idx: カテゴリーインデックス
// name: 省略会社名
func GenerateCompanyPathistFolder(base string, idx int32, name string) string {
	folderName := strconv.Itoa(int(idx)) + " " + name
	return filepath.Join(base, folderName)
}
