package models

import (
	"errors"
	"slices"
	"strings"
)

// Company は工事会社の基本情報をファイル名から取得したモデルを表します
// @Description 工事会社の基本情報をファイル名から取得したモデル
type Company struct {
	// 基本フィールド
	ID         string `json:"id" yaml:"-" example:"TC001"`
	FolderName string `json:"folderName" example:"0 豊田築炉"`

	// パス名からの固有フィールド
	ShortName    string      `json:"short_name,omitempty" yaml:"-" example:"豊田築炉"`
	BusinessType CompanyType `json:"business_type,omitempty" yaml:"-" example:"4"`

	// 属性ファイルフィールド
	LongName   string   `json:"long_name,omitempty" yaml:"long_name" example:"有限会社 豊田築炉"`
	PostalCode string   `json:"postal_code,omitempty" yaml:"postal_code" example:"456-0001"`
	Address    string   `json:"address,omitempty" yaml:"address" example:"愛知県名古屋市熱田区三本松町1-1"`
	Phone      string   `json:"phone,omitempty" yaml:"phone" example:"052-681-8111"`
	Email      string   `json:"email,omitempty" yaml:"email" example:"info@toyotachikuro.jp"`
	Website    string   `json:"website,omitempty" yaml:"website" example:"https://www.toyotachikuro.jp"`
	Tags       []string `json:"tags,omitempty" yaml:"tags" example:"['元請け', '製造業']"`
}

// GetFolderName フォルダー名を取得します
// AttributeServiceで使用
func (c *Company) GetFolderName() string {
	return c.FolderName
}

// UpdateFolderName ShortNameとBusinessTypeからFolderNameを更新します
// input: 引数が2つの場合は、引数のフォルダー名と会社名を使用して更新します
func (c *Company) UpdateFolderName(input ...string) error {
	if len(input) == 2 {
		c.FolderName = input[0] + " " + input[1]
		c.BusinessType = ParseBusinessTypeFromCode(input[0])
		c.ShortName = input[1]
	} else {
		c.FolderName = c.BusinessType.Code() + " " + c.ShortName
	}
	return nil
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

	businessType := ParseBusinessTypeFromCode(folderName[:1])

	return Company{
		ID:           NewIDFromString(parts[1]).Len5(),
		FolderName:   folderName, // フォルダー名のみを格納
		LongName:     parts[1],
		ShortName:    parts[1],
		BusinessType: businessType,
		Tags:         tags,
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

	businessType := ParseBusinessTypeFromCode(name[:1])
	result[0] = businessType.String()
	result[1] = shortName

	return result, nil
}

// CompanyType は業種を表すenum型
type CompanyType int

const (
	CompanyTypeSpecial  CompanyType = iota // 0: 特別(自社・組合・未定)
	CompanyTypeAgency                      // 1: 応援会社
	CompanyTypePeer                        // 2: 同業者（築炉会社）
	CompanyTypePersonal                    // 3: 一人親方
	CompanyTypePrime                       // 4: 元請け
	CompanyTypeLease                       // 5: リース会社
	CompanyTypeSales                       // 6: 販売会社
	CompanyTypeSales2                      // 7: 販売会社（重複）
	CompanyTypeRecruit                     // 8: 求人会社
	CompanyTypeOther                       // 9: その他
)

// String はBusinessTypeの文字列表現を返します
func (bt CompanyType) String() string {
	switch bt {
	case CompanyTypeSpecial:
		return "特別"
	case CompanyTypeAgency:
		return "下請会社"
	case CompanyTypePeer:
		return "築炉会社"
	case CompanyTypePersonal:
		return "一人親方"
	case CompanyTypePrime:
		return "元請け"
	case CompanyTypeLease:
		return "リース会社"
	case CompanyTypeSales, CompanyTypeSales2:
		return "販売会社"
	case CompanyTypeRecruit:
		return "求人会社"
	case CompanyTypeOther:
		return "その他"
	default:
		return "特別"
	}
}

// 事前計算されたコード文字列
var companyTypeCodes = [10]string{
	"0", // CompanyTypeSpecial
	"1", // CompanyTypeAgency
	"2", // CompanyTypePeer
	"3", // CompanyTypePersonal
	"4", // CompanyTypePrime
	"5", // CompanyTypeLease
	"6", // CompanyTypeSales
	"7", // CompanyTypeSales2
	"8", // CompanyTypeRecruit
	"9", // CompanyTypeOther
}

// Code はBusinessTypeの数字コードを返します（極限高速化版）
//
// 戻り値: "0"〜"9"の文字列、無効な値の場合は"0"
func (bt CompanyType) Code() string {
	if bt >= 0 && bt < 10 {
		return companyTypeCodes[bt]
	}
	return "0" // 無効値のデフォルト
}

// IsValid はBusinessTypeが有効な値かチェックします
func (bt CompanyType) IsValid() bool {
	return bt >= CompanyTypeSpecial && bt <= CompanyTypeOther
}

// AllBusinessTypes は全ての有効なBusinessTypeを返します
func AllBusinessTypes() []CompanyType {
	return []CompanyType{
		CompanyTypeSpecial,
		CompanyTypeAgency,
		CompanyTypePeer,
		CompanyTypePersonal,
		CompanyTypePrime,
		CompanyTypeLease,
		CompanyTypeSales,
		CompanyTypeSales2,
		CompanyTypeRecruit,
		CompanyTypeOther,
	}
}

// ParseBusinessTypeFromCode は数字コードからBusinessTypeに変換します
func ParseBusinessTypeFromCode(code string) CompanyType {
	switch code {
	case "0":
		return CompanyTypeSpecial
	case "1":
		return CompanyTypeAgency
	case "2":
		return CompanyTypePeer
	case "3":
		return CompanyTypePersonal
	case "4":
		return CompanyTypePrime
	case "5":
		return CompanyTypeLease
	case "6":
		return CompanyTypeSales
	case "7":
		return CompanyTypeSales2
	case "8":
		return CompanyTypeRecruit
	case "9":
		return CompanyTypeOther
	default:
		return CompanyTypeSpecial
	}
}

// ParseBusinessTypeFromString は文字列からBusinessTypeに変換します
func ParseBusinessTypeFromString(name string) CompanyType {
	for _, bt := range AllBusinessTypes() {
		if bt.String() == name {
			return bt
		}
	}
	return CompanyTypeOther
}

// GetAllBusinessTypesMap は互換性のためのマップを返します
func GetAllBusinessTypesMap() map[string]string {
	result := make(map[string]string)
	for _, bt := range AllBusinessTypes() {
		result[bt.Code()] = bt.String()
	}
	return result
}

// AddTag adds a tag to the company's tags list
func (c *Company) AddTag(tag string) {
	if slices.Contains(c.Tags, tag) {
		return // Tag already exists
	}
	c.Tags = append(c.Tags, tag)
}
