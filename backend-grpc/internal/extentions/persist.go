package exts

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// ObjectPersistable は永続化されるエンティティに必要な振る舞いを定義します。
type ObjectPersistable interface {
	// GetPersistPath はファイルの永続化フルパスを取得します。
	GetPersistPath() string

	// GetObject は永続化対象のオブジェクトを取得します。
	GetObject() any

	// SetObject は永続化対象のオブジェクトを設定します。
	SetObject(any)
}

// ObjectPersistService はファイルベースのデータ永続化のサービスを提供します。
type ObjectPersistService[T ObjectPersistable] struct {
	PersistFilename string
}

// Load は永続化ファイルからデータを読み込みます。
func (s *ObjectPersistService[T]) Load(entity T) (*T, error) {

	// 永続化ファイルのフルパスを取得
	persistPath := entity.GetPersistPath()

	// ファイルを読み込み
	data, err := os.ReadFile(persistPath)
	if err != nil {
		return nil, err
	}

	// YAMLをデコード
	object := entity.GetObject()
	if err := yaml.Unmarshal(data, object); err != nil {
		return nil, err
	}

	// エンティティにデコード結果を設定
	entity.SetObject(object)

	return &entity, nil
}

// Save はデータを永続化ファイルに保存します。
func (s *ObjectPersistService[T]) Save(entity T) error {
	// 永続化ファイルの有無等チェック
	persistPath := entity.GetPersistPath()
	if _, err := os.Stat(persistPath); os.IsNotExist(err) {
		return fmt.Errorf("永続化ファイルが存在しません: %s", persistPath)
	} else if err != nil {
		return err
	}

	// データをYAMLにエンコード
	obj := entity.GetObject()
	yamlText, err := yaml.Marshal(obj)
	if err != nil {
		return fmt.Errorf("データのエンコードに失敗しました: %w", err)
	}

	// ファイルに書き込み
	return os.WriteFile(persistPath, yamlText, 0644)
}
