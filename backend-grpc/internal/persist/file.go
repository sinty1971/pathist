package persist

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// FilePersistable は永続化されるエンティティに必要な振る舞いを定義します。
type FilePersistable interface {
	// GetFilePersistPath はファイルの永続化フルパスを取得します。
	GetFilePersistPath() string

	// GetObject は永続化対象のオブジェクトを取得します。
	GetObject() any

	// SetObject は永続化対象のオブジェクトを設定します。
	SetObject(any)
}

// FilePersistService はファイルベースのデータ永続化のサービスを提供します。
type FilePersistService[T FilePersistable] struct {
	PersistFilename string
}

// Load は永続化ファイルからデータを読み込みます。
func (s *FilePersistService[T]) Load(entity T) (*T, error) {

	// 永続化ファイルのフルパスを取得
	persistPath := entity.GetFilePersistPath()

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
func (s *FilePersistService[T]) Save(entity T) error {
	// 永続化ファイルの有無等チェック
	persistPath := entity.GetFilePersistPath()
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
