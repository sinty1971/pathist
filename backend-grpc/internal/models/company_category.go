package models

import "errors"

// CompanyCategoryIndex は業種を表すenum型（int）
type CompanyCategoryIndex int

// 業種カテゴリーの定義
const (
	CompanyCategoryMin      CompanyCategoryIndex = 0
	CompanyCategoryUnion    CompanyCategoryIndex = 0
	CompanyCategoryAgency   CompanyCategoryIndex = 1
	CompanyCategoryPeer     CompanyCategoryIndex = 2
	CompanyCategoryPersonal CompanyCategoryIndex = 3
	CompanyCategoryPrime    CompanyCategoryIndex = 4
	CompanyCategoryLease    CompanyCategoryIndex = 5
	CompanyCategorySales    CompanyCategoryIndex = 6
	CompanyCategorySales2   CompanyCategoryIndex = 7
	CompanyCategoryRecruit  CompanyCategoryIndex = 8
	CompanyCategoryOther    CompanyCategoryIndex = 9
	CompanyCategoryMax      CompanyCategoryIndex = 9
)

// CompanyCategoryMap 業種カテゴリーのラベルマップです
// 将来的にはyamlファイルから読み込む予定
var CompanyCategoryMap = map[CompanyCategoryIndex]string{
	CompanyCategoryUnion:    "自社組合",
	CompanyCategoryAgency:   "下請会社",
	CompanyCategoryPeer:     "築炉会社",
	CompanyCategoryPersonal: "一人親方",
	CompanyCategoryPrime:    "元請会社",
	CompanyCategoryLease:    "リース会社",
	CompanyCategorySales:    "販売会社",
	CompanyCategorySales2:   "販売会社２",
	CompanyCategoryRecruit:  "求人会社",
	CompanyCategoryOther:    "一般会社",
}

// CompanyCategoryReverseMap は業種カテゴリーの逆引きマップです
var CompanyCategoryReverseMap = map[string]CompanyCategoryIndex{}

// init はパッケージ初期化時に呼び出され、逆引きマップを初期化します
func init() {
	// CompanyCategoryReverseMapを初期化
	for code, label := range CompanyCategoryMap {
		CompanyCategoryReverseMap[label] = code
	}
}

func (h CompanyCategoryIndex) ToInt32() int32 {
	return int32(h)
}

// IsValid は CompanyCategoryIndex が有効な範囲内かをチェックします
func (h CompanyCategoryIndex) Error() error {
	if CompanyCategoryMin <= h && h <= CompanyCategoryMax {
		return nil
	}
	return errors.New("invalid CompanyCategoryIndex")
}
