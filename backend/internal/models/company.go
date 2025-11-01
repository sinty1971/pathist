package models

import (
	"encoding/json"
	"errors"
	"slices"
	"strings"

	grpcv1 "grpc-backend/gen/grpc/v1"
	"grpc-backend/internal/utils"
)

// CompanyEx は gRPC grpc.v1.Company メッセージの拡張版です。
type Company struct {
	*grpcv1.Company

	//insideFilename は管理フォルダー内の会社データファイル名を保持します
	insideFilename string
}

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

// CompanyCategory は業種コードとラベルのペアを表します。
type CompanyCategory struct {
	Index CompanyCategoryIndex `json:"code" example:"1"`
	Label string               `json:"label" example:"下請会社"`
}

// NewCompany 会社フォルダーパス名からCompanyを作成します
func NewCompany(managedFolder string) (*Company, error) {
	// フォルダー名を取得
	var folderName string
	if lastSlash := strings.LastIndex(managedFolder, "/"); lastSlash != -1 {
		folderName = managedFolder[lastSlash+1:]
	}

	// ファイル名の[0-9] [会社名]の規則を解析
	result, err := ParseCompanyName(folderName)
	if err != nil {
		return nil, err
	}

	company := &Company{
		Company: grpcv1.Company_builder{
			Id:            GenerateCompanyId(result.ShortName),
			ManagedFolder: managedFolder,
			ShortName:     result.ShortName,
			Category:      string(result.Category),

			InsideIdealPath:     "",
			InsideLegalName:     result.ShortName,
			InsideTags:          append([]string{"会社"}, result.Tags...),
			InsideRequiredFiles: []*grpcv1.FileInfo{},
			InsidePostalCode:    "",
			InsideAddress:       "",
			InsidePhone:         "",
			InsideEmail:         "",
			InsideWebsite:       "",
		}.Build(),
	}

	// IDを更新とリターン
	return company, nil
}

func GenerateCompanyId(shortName string) string {
	return utils.GenerateIdFromString(shortName)
}

// parseCompanyNameResult は ParseCompanyName の戻り値を表します
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
	if !(&code).IsValid() {
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

// IsValid は CompanyCategoryIndex が有効な範囲内かをチェックします
func (cc *CompanyCategoryIndex) IsValid() bool {
	if cc == nil {
		return false
	}
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
	if c == nil {
		return
	}
	if slices.Contains(c.GetInsideTags(), tag) {
		return // タグは既に存在します
	}
	c.SetInsideTags(append(c.GetInsideTags(), tag))
}

// MarshalYAML は YAML 用のシリアライズを行います。
func (c Company) MarshalYAML() (any, error) {
	return c.Company, nil
}

// UnmarshalYAML は YAML からの復元を行います。
func (c *Company) UnmarshalYAML(unmarshal func(any) error) error {
	enc := &grpcv1.Company{}
	if err := unmarshal(enc); err != nil {
		return err
	}
	c.Company = enc
	return nil
}

// MarshalJSON は JSON 用のシリアライズを行います。
func (c Company) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Company)
}

// UnmarshalJSON は JSON からの復元を行います。
func (c *Company) UnmarshalJSON(data []byte) error {
	enc := &grpcv1.Company{}
	if err := json.Unmarshal(data, enc); err != nil {
		return err
	}
	c.Company = enc
	return nil
}
