package models

import (
	"errors"
	"slices"
	"strings"
)

// Company は工事会社の基本情報をファイル名から取得したモデルを表します
// @Description 工事会社の基本情報をファイル名から取得したモデル
type Company struct {
	// 基本のFileInfo構造体を埋め込み
	FileInfo

	// 計算フィールド
	ID string `json:"id" yaml:"-" example:"TC001"`

	// パス名からの固有フィールド
	ShortName    string `json:"short_name,omitempty" yaml:"-" example:"豊田築炉"`
	BusinessType string `json:"business_type,omitempty" yaml:"-" example:"元請け"`

	// 属性ファイルフィールド
	FullName   string   `json:"full_name,omitempty" yaml:"full_name" example:"有限会社 豊田築炉"`
	PostalCode string   `json:"postal_code,omitempty" yaml:"postal_code" example:"456-0001"`
	Address    string   `json:"address,omitempty" yaml:"address" example:"愛知県名古屋市熱田区三本松町1-1"`
	Phone      string   `json:"phone,omitempty" yaml:"phone" example:"052-681-8111"`
	Email      string   `json:"email,omitempty" yaml:"email" example:"info@toyotachikuro.jp"`
	Website    string   `json:"website,omitempty" yaml:"website" example:"https://www.toyotachikuro.jp"`
	Tags       []string `json:"tags,omitempty" yaml:"tags" example:"['元請け', '製造業']"`
}

// GetFileInfo AttributeServiceで使用するためのメソッド
func (c *Company) GetFileInfo() FileInfo {
	return c.FileInfo
}

// NewCompany FileInfoからCompanyを作成します
func NewCompany(fileInfo FileInfo) (Company, error) {
	// ファイル名の[0-9] [会社名]の規則を解析
	parts, err := ParseCompanyName(fileInfo.Name)
	if err != nil {
		return Company{}, err
	}

	// 基本タグを作成
	tags := []string{"会社", parts[0], parts[1]}

	// 会社名内にハイフンが含まれているときは関連会社名または業務内容をタグに追加
	companyPart := fileInfo.Name[2:] // "0 " を除いた部分
	if hyphenIndex := strings.Index(companyPart, "-"); hyphenIndex != -1 {
		// ハイフン以降の文字列を取得
		relatedInfo := companyPart[hyphenIndex+1:]
		if relatedInfo != "" {
			tags = append(tags, relatedInfo)
		}
	}

	return Company{
		FileInfo:     fileInfo,
		ID:           NewIDFromString(parts[1]).Len5(),
		FullName:     parts[1],
		ShortName:    parts[1],
		BusinessType: parts[0],
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

	result[0] = DetermineBusinessType(name[:1])
	result[1] = shortName

	return result, nil
}

func DetermineBusinessType(name string) string {
	switch name {
	case "0":
		return "自社"
	case "1":
		return "下請会社"
	case "2":
		return "築炉会社"
	case "3":
		return "一人親方"
	case "4":
		return "元請け"
	case "5":
		return "リース会社"
	case "6":
		return "販売会社"
	case "7":
		return "販売会社"
	case "8":
		return "求人会社"
	case "9":
		return "その他"
	}
	return "不明"
}

// AddTag adds a tag to the company's tags list
func (c *Company) AddTag(tag string) {
	if slices.Contains(c.Tags, tag) {
		return // Tag already exists
	}
	c.Tags = append(c.Tags, tag)
}
