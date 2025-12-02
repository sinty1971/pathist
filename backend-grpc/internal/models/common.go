package models

import "slices"

// CommonModel は共通フィールドを持つモデルのインターフェースを定義します。
type CommonModel interface {
	GetInsideTags() []string
	SetInsideTags([]string)
}

// AddInsideTags はゲッター・セッターを使用してタグを追加します（重複・空文字はスキップ）
// 想定タグ数: 10個程度
func AddInsideTags(obj CommonModel, newTags ...string) {
	existingTags := obj.GetInsideTags()
	for _, tag := range newTags {
		if tag != "" && !slices.Contains(existingTags, tag) {
			existingTags = append(existingTags, tag)
		}
	}
	obj.SetInsideTags(existingTags)
}
