package models

import (
	"errors"
	"path"
	"slices"
	"strings"
)

// CompanyCategoryIndex は業種を表すenum型（string）
type CompanyCategoryIndex string

const (
	CompanyCategorySpecial  CompanyCategoryIndex = "0" // "0": 特別(自社・組合・未定)
	CompanyCategoryAgency   CompanyCategoryIndex = "1" // "1": 応援会社
	CompanyCategoryPeer     CompanyCategoryIndex = "2" // "2": 同業者（築炉会社）
	CompanyCategoryPersonal CompanyCategoryIndex = "3" // "3": 一人親方
	CompanyCategoryPrime    CompanyCategoryIndex = "4" // "4": 元請け
	CompanyCategoryLease    CompanyCategoryIndex = "5" // "5": リース会社
	CompanyCategorySales    CompanyCategoryIndex = "6" // "6": 販売会社
	CompanyCategorySales2   CompanyCategoryIndex = "7" // "7": 販売会社（重複）
	CompanyCategoryRecruit  CompanyCategoryIndex = "8" // "8": 求人会社
	CompanyCategoryOther    CompanyCategoryIndex = "9" // 9: その他
)

// CompanyCategories はアプリケーション全体で使用される業種カテゴリーの一覧です
// 将来的にはyamlファイルから読み込む予定
var CompanyCategoryMap = map[CompanyCategoryIndex]string{
	CompanyCategorySpecial:  "特別",
	CompanyCategoryAgency:   "下請会社",
	CompanyCategoryPeer:     "築炉会社",
	CompanyCategoryPersonal: "一人親方",
	CompanyCategoryPrime:    "元請け",
	CompanyCategoryLease:    "リース会社",
	CompanyCategorySales:    "販売会社",
	CompanyCategorySales2:   "販売会社２",
	CompanyCategoryRecruit:  "求人会社",
	CompanyCategoryOther:    "その他",
}

var CompanyCategoryReverseMap = map[string]CompanyCategoryIndex{}

// Company は工事会社の基本情報をファイル名から取得したモデルを表します
// @Description 工事会社の基本情報をファイル名から取得したモデル
type Company struct {
	// ID is generated from FolderName.
	ID string `json:"id" yaml:"-" example:"TC001"`

	// Identity fields. FolderName is linked ShortName and Category.
	FolderPath string               `json:"folderName" yaml:"-" example:"~/penguin/豊田築炉/1 会社/0 豊田築炉"`
	ShortName  string               `json:"shortName,omitempty" yaml:"-" example:"豊田築炉"`
	Category   CompanyCategoryIndex `json:"category" yaml:"-" example:"1"`

	// データベースファイルフィールド
	LegalName  string   `json:"legalName,omitempty" yaml:"legal_name" example:"有限会社 豊田築炉"`
	PostalCode string   `json:"postalCode,omitempty" yaml:"postal_code" example:"456-0001"`
	Address    string   `json:"address,omitempty" yaml:"address" example:"愛知県名古屋市熱田区三本松町1-1"`
	Phone      string   `json:"phone,omitempty" yaml:"phone" example:"052-681-8111"`
	Email      string   `json:"email,omitempty" yaml:"email" example:"info@toyotachikuro.jp"`
	Website    string   `json:"website,omitempty" yaml:"website" example:"https://www.toyotachikuro.jp"`
	Tags       []string `json:"tags,omitempty" yaml:"tags" example:"['元請け', '製造業']"`

	// 標準ファイルフィールド
	StandardFiles []FileInfo `json:"standardFiles" yaml:"standard_files"`
}

// GetFolderPath フォルダーパスを取得します
// DatabaseServiceで使用
func (c *Company) GetFolderPath() string {
	return c.FolderPath
}

// NewCompany 会社フォルダーパス名からCompanyを作成します
func NewCompany(folderPath string) (*Company, error) {
	// フォルダー名を取得
	var folderName string
	if lastSlash := strings.LastIndex(folderPath, "/"); lastSlash != -1 {
		folderName = folderPath[lastSlash+1:]
	}

	// ファイル名の[0-9] [会社名]の規則を解析
	result, err := ParseCompanyName(folderName)
	if err != nil {
		return nil, err
	}

	return &Company{
		ID:         NewIDFromString(result.ShortName).Len5(),
		FolderPath: folderPath, // フォルダー名のみを格納
		LegalName:  result.ShortName,
		ShortName:  result.ShortName,
		Category:   result.Category,
		Tags:       append([]string{"会社"}, result.Tags...),
	}, nil
}

type parseCompanyNameResult struct {
	Category  CompanyCategoryIndex
	ShortName string
	Tags      []string
}

// ParseCompanyName は"[0-9] [会社名]"形式のファイル名を解析します
// 会社名内にハイフンが含まれている場合は関連会社名または業務内容として処理します
// 戻り値: [0]業種, [1]会社名, error
func ParseCompanyName(name string) (parseCompanyNameResult, error) {
	result := parseCompanyNameResult{}

	// 最小長チェック（"0 X"の最小3文字）
	if len(name) < 3 {
		return result, errors.New("ファイル名が短すぎます")
	}

	// 2番目の文字がスペースかチェック
	if name[1] != ' ' {
		return result, errors.New("2番目の文字がスペースではありません")
	}

	// 最初の文字がCompanyCategoryCodeかチェック
	code := CompanyCategoryIndex(name[0])
	if !code.IsValid() {
		return result, errors.New("最初の文字が業種コードではありません")
	}

	// 会社名部分を取得
	companyPart := name[2:]
	if companyPart == "" {
		return result, errors.New("会社名が空です")
	}

	// 業種を取得
	result.Category = code

	// 会社名内にハイフンが含まれているかチェック
	var relatedInfo string
	if hyphenIndex := strings.Index(companyPart, "-"); hyphenIndex != -1 {
		// ハイフン以降の文字列を取得
		relatedInfo = companyPart[hyphenIndex+1:]
		result.ShortName = companyPart[:hyphenIndex]
	} else {
		// ハイフンがない場合は全体を会社名とする
		result.ShortName = companyPart
	}

	// タグを作成
	result.Tags = []string{CompanyCategoryMap[result.Category], result.ShortName}
	if relatedInfo != "" {
		result.Tags = append(result.Tags, relatedInfo)
	}

	return result, nil
}

// UpdateIdentity CategoryとShortNameからIDとFolderNameを更新します
// 戻り値: 変更前のフォルダー名, エラー
func (c *Company) UpdateIdentity() {
	// フォルダー名の生成
	folderName := string(c.Category) + " " + c.ShortName

	// フォルダー名とIDを更新
	c.ID = NewIDFromString(folderName).Len5()

	// フォルダーパスを更新
	dir := path.Dir(c.FolderPath)
	if dir == "." {
		return
	}
	c.FolderPath = path.Join(dir, folderName)
}

func (cc *CompanyCategoryIndex) IsValid() bool {
	return *cc >= CompanyCategorySpecial && *cc <= CompanyCategoryOther
}

// ParseCompanyCategoryCode は数字コードからカテゴリー文字列に変換します
func ParseCompanyCategoryCode(codeStr string) string {
	// 数字コードをCompanyCategoryCodeに変換
	// 変換に失敗してもCompanyCategoryCodeは文字列のため codeStr がそのまま返される
	code := CompanyCategoryIndex(codeStr)
	if str, exists := CompanyCategoryMap[code]; exists {
		return str
	}
	return CompanyCategoryMap[CompanyCategorySpecial]
}

// ParseCompanyCategoryName は文字列からCompanyCategoryCodeに変換します
func ParseCompanyCategoryName(name string) CompanyCategoryIndex {
	if code, exists := CompanyCategoryReverseMap[name]; exists {
		return code
	}
	return CompanyCategorySpecial
}

// AddTag は会社のタグリストにタグを追加します
func (c *Company) AddTag(tag string) {
	if slices.Contains(c.Tags, tag) {
		return // タグは既に存在します
	}
	c.Tags = append(c.Tags, tag)
}
