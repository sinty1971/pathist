package models

import "errors"

// 業種カテゴリーの定義
const (
	CompanyCategoryMin      int = 0
	CompanyCategoryUnion    int = 0
	CompanyCategoryAgency   int = 1
	CompanyCategoryPeer     int = 2
	CompanyCategoryPersonal int = 3
	CompanyCategoryPrime    int = 4
	CompanyCategoryLease    int = 5
	CompanyCategorySales    int = 6
	CompanyCategorySales2   int = 7
	CompanyCategoryRecruit  int = 8
	CompanyCategoryOther    int = 9
	CompanyCategoryMax      int = 9
)

// CompanyCategoryMap 業種カテゴリーのラベルマップです
// 将来的にはyamlファイルから読み込む予定
var CompanyCategoryMap = map[int]string{
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
var CompanyCategoryReverseMap = map[string]int{}

// init はパッケージ初期化時に呼び出され、逆引きマップを初期化します
func init() {
	// CompanyCategoryReverseMapを初期化
	for code, label := range CompanyCategoryMap {
		CompanyCategoryReverseMap[label] = code
	}
}

// ErrorCompanyCategoryIndex 引数: idx が有効な範囲内かをチェックします
func ErrorCompanyCategoryIndex(idx int) error {
	if CompanyCategoryMin <= idx && idx <= CompanyCategoryMax {
		return nil
	}
	return errors.New("invalid CompanyCategoryIndex")
}
