package models

import (
	"errors"
	"slices"
	"strings"
)

// Company は工事会社の基本情報をファイル名から取得したモデルを表します
// @Description 工事会社の基本情報をファイル名から取得したモデル
type Company struct {
	// ID is generated from FolderName.
	ID string `json:"id" yaml:"-" example:"TC001"`

	// Identity fields. FolderName is linked ShortName and Category.
	FolderName string          `json:"folderName" yaml:"-" example:"0 豊田築炉"`
	ShortName  string          `json:"shortName,omitempty" yaml:"-" example:"豊田築炉"`
	Category   CompanyCategory `json:"category" yaml:"-" example:"4"`

	// Database file fields
	LegalName  string   `json:"legalName,omitempty" yaml:"legal_name" example:"有限会社 豊田築炉"`
	PostalCode string   `json:"postalCode,omitempty" yaml:"postal_code" example:"456-0001"`
	Address    string   `json:"address,omitempty" yaml:"address" example:"愛知県名古屋市熱田区三本松町1-1"`
	Phone      string   `json:"phone,omitempty" yaml:"phone" example:"052-681-8111"`
	Email      string   `json:"email,omitempty" yaml:"email" example:"info@toyotachikuro.jp"`
	Website    string   `json:"website,omitempty" yaml:"website" example:"https://www.toyotachikuro.jp"`
	Tags       []string `json:"tags,omitempty" yaml:"tags" example:"['元請け', '製造業']"`

	// 補助ファイルフィールド
	Assists []AssistFile `json:"assists" yaml:"assists"`
}

// GetFolderName フォルダー名を取得します
// DatabaseServiceで使用
func (c *Company) GetFolderName() string {
	return c.FolderName
}

// NewCompany フォルダーパスからCompanyを作成します
func NewCompany(folderPath string) (Company, error) {
	// フォルダー名を取得（'/'が含まれている場合は最後の要素を取得）
	folderName := folderPath
	if lastSlash := strings.LastIndex(folderPath, "/"); lastSlash != -1 {
		folderName = folderPath[lastSlash+1:]
	}

	// ファイル名の[0-9] [会社名]の規則を解析
	parts, err := ParseCompanyName(folderName)
	if err != nil {
		return Company{}, err
	}

	// 基本タグを作成
	tags := []string{"会社", parts[0], parts[1]}

	// 会社名内にハイフンが含まれているときは関連会社名または業務内容をタグに追加
	companyPart := folderName[2:] // "0 " を除いた部分
	if hyphenIndex := strings.Index(companyPart, "-"); hyphenIndex != -1 {
		// ハイフン以降の文字列を取得
		relatedInfo := companyPart[hyphenIndex+1:]
		if relatedInfo != "" {
			tags = append(tags, relatedInfo)
		}
	}

	companyCategory := ParseCompanyCategoryFromCode(folderName[:1])

	return Company{
		ID:         NewIDFromString(parts[1]).Len5(),
		FolderName: folderName, // フォルダー名のみを格納
		LegalName:  parts[1],
		ShortName:  parts[1],
		Category:   companyCategory,
		Tags:       tags,
	}, nil
}

// ParseCompanyName は"[0-9] [会社名]"形式のファイル名を解析します
// 会社名内にハイフンが含まれている場合は関連会社名または業務内容として処理します
// 戻り値: [0]業種, [1]会社名, error
func ParseCompanyName(name string) ([2]string, error) {
	var result [2]string

	// 最小長チェック（"0 X"の最小3文字）
	if len(name) < 3 {
		return result, errors.New("ファイル名が短すぎます")
	}

	// 最初の文字が数字かチェック
	if name[0] < '0' || name[0] > '9' {
		return result, errors.New("最初の文字が数字ではありません")
	}

	// 2番目の文字がスペースかチェック
	if name[1] != ' ' {
		return result, errors.New("2番目の文字がスペースではありません")
	}

	// スペース以降の文字列を取得
	companyPart := name[2:]

	// 空文字チェック
	if companyPart == "" {
		return result, errors.New("会社名が空です")
	}

	// 会社名内にハイフンが含まれているかチェック
	var shortName string
	if hyphenIndex := strings.Index(companyPart, "-"); hyphenIndex != -1 {
		// ハイフンが見つかった場合、ハイフン以前を会社名とする
		shortName = companyPart[:hyphenIndex]
	} else {
		// ハイフンがない場合は全体を会社名とする
		shortName = companyPart
	}

	companyCategory := ParseCompanyCategoryFromCode(name[:1])
	result[0] = companyCategory.String()
	result[1] = shortName

	return result, nil
}

// UpdateIdentity CategoryとShortNameからIDとFolderNameを更新します
// 戻り値: 変更前のフォルダー名, エラー
func (c *Company) UpdateIdentity() error {
	// 変更前のフォルダー名を取得
	prevFolderName := c.FolderName

	// フォルダー名の生成
	folderName := c.Category.Code() + " " + c.ShortName

	// 変更前のフォルダー名と変更後のフォルダー名が同じ場合はエラー
	if folderName == prevFolderName {
		return errors.New("ID及びFolderNameは変更されませんでした")
	}

	// フォルダー名とIDを更新
	c.ID = NewIDFromString(folderName).Len5()
	c.FolderName = folderName

	return nil
}

// CompanyCategory は業種を表すenum型
type CompanyCategory int

const (
	CompanyCategorySpecial  CompanyCategory = iota // 0: 特別(自社・組合・未定)
	CompanyCategoryAgency                          // 1: 応援会社
	CompanyCategoryPeer                            // 2: 同業者（築炉会社）
	CompanyCategoryPersonal                        // 3: 一人親方
	CompanyCategoryPrime                           // 4: 元請け
	CompanyCategoryLease                           // 5: リース会社
	CompanyCategorySales                           // 6: 販売会社
	CompanyCategorySales2                          // 7: 販売会社（重複）
	CompanyCategoryRecruit                         // 8: 求人会社
	CompanyCategoryOther                           // 9: その他
)

// String はCompanyCategoryの文字列表現を返します
func (bt CompanyCategory) String() string {
	switch bt {
	case CompanyCategorySpecial:
		return "特別"
	case CompanyCategoryAgency:
		return "下請会社"
	case CompanyCategoryPeer:
		return "築炉会社"
	case CompanyCategoryPersonal:
		return "一人親方"
	case CompanyCategoryPrime:
		return "元請け"
	case CompanyCategoryLease:
		return "リース会社"
	case CompanyCategorySales, CompanyCategorySales2:
		return "販売会社"
	case CompanyCategoryRecruit:
		return "求人会社"
	case CompanyCategoryOther:
		return "その他"
	default:
		return "特別"
	}
}

// 事前計算されたコード文字列
var companyCategoryCodes = [10]string{
	"0", // CompanyCategorySpecial
	"1", // CompanyCategoryAgency
	"2", // CompanyCategoryPeer
	"3", // CompanyCategoryPersonal
	"4", // CompanyCategoryPrime
	"5", // CompanyCategoryLease
	"6", // CompanyCategorySales
	"7", // CompanyCategorySales2
	"8", // CompanyCategoryRecruit
	"9", // CompanyCategoryOther
}

// Code はCompanyCategoryの数字コードを返します（極限高速化版）
//
// 戻り値: "0"〜"9"の文字列、無効な値の場合は"0"
func (bt CompanyCategory) Code() string {
	if bt >= 0 && bt < 10 {
		return companyCategoryCodes[bt]
	}
	return "0" // 無効値のデフォルト
}

// IsValid はCompanyCategoryが有効な値かチェックします
func (bt CompanyCategory) IsValid() bool {
	return bt >= CompanyCategorySpecial && bt <= CompanyCategoryOther
}

// AllCompanyCategories は全ての有効なCompanyCategoryを返します
func AllCompanyCategories() []CompanyCategory {
	return []CompanyCategory{
		CompanyCategorySpecial,
		CompanyCategoryAgency,
		CompanyCategoryPeer,
		CompanyCategoryPersonal,
		CompanyCategoryPrime,
		CompanyCategoryLease,
		CompanyCategorySales,
		CompanyCategorySales2,
		CompanyCategoryRecruit,
		CompanyCategoryOther,
	}
}

// ParseCompanyCategoryFromCode は数字コードからCompanyCategoryに変換します
func ParseCompanyCategoryFromCode(code string) CompanyCategory {
	switch code {
	case "0":
		return CompanyCategorySpecial
	case "1":
		return CompanyCategoryAgency
	case "2":
		return CompanyCategoryPeer
	case "3":
		return CompanyCategoryPersonal
	case "4":
		return CompanyCategoryPrime
	case "5":
		return CompanyCategoryLease
	case "6":
		return CompanyCategorySales
	case "7":
		return CompanyCategorySales2
	case "8":
		return CompanyCategoryRecruit
	case "9":
		return CompanyCategoryOther
	default:
		return CompanyCategorySpecial
	}
}

// ParseCompanyCategoryFromString は文字列からCompanyCategoryに変換します
func ParseCompanyCategoryFromString(name string) CompanyCategory {
	// slices.IndexFuncを使用して高速化
	index := slices.IndexFunc(AllCompanyCategories(), func(category CompanyCategory) bool {
		return category.String() == name
	})
	if index != -1 {
		return AllCompanyCategories()[index]
	}
	return CompanyCategorySpecial
}

// GetAllCompanyCategoriesMap は互換性のためのマップを返します
func GetAllCompanyCategoriesMap() map[string]string {
	result := make(map[string]string)
	for _, category := range AllCompanyCategories() {
		result[category.Code()] = category.String()
	}
	return result
}

// AddTag は会社のタグリストにタグを追加します
func (c *Company) AddTag(tag string) {
	if slices.Contains(c.Tags, tag) {
		return // タグは既に存在します
	}
	c.Tags = append(c.Tags, tag)
}
