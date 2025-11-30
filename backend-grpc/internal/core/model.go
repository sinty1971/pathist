package core

import "slices"

type DefaultModel interface {
	GetInsideTags() []string
	SetInsideTags([]string)
}

// DefaultAddInsideTags はゲッター・セッターを使用してタグを追加します（重複・空文字はスキップ）
// 想定タグ数: 10個程度
func ModelAddInsideTags(obj DefaultModel, newTags ...string) {
	existingTags := obj.GetInsideTags()
	for _, tag := range newTags {
		if tag != "" && !slices.Contains(existingTags, tag) {
			existingTags = append(existingTags, tag)
		}
	}
	obj.SetInsideTags(existingTags)
}
