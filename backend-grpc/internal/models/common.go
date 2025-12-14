package models

import (
	"backend-grpc/internal/core"
	"path/filepath"
)

// CommonModel は共通フィールドを持つモデルのインターフェースを定義します。
type CommonModel interface {
	// GetTarget は永続化ファイルを保存するフルパスを取得します。
	// protobuf メッセージ内の target フィールドを返す実装が一般的です。
	GetTarget() string
	SetId(id string)
}

// ParamFunc はゲッター・セッター関数のペアを保持します。
type Common struct {
	target     CommonModel
	modeleName string
}

// NewCommon は Common インスタンスを作成します。
func NewCommon(target CommonModel, modeleName string) *Common {
	return &Common{
		target:     target,
		modeleName: modeleName,
	}
}

func (p *Common) SetDefaultId() {
	idText := p.modeleName + filepath.Base(p.target.GetTarget())
	id := core.GenerateIdFromString(idText)
	p.target.SetId(id)
}
