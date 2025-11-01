package services

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// Persistable はファイルに永続化されるエンティティに必要な振る舞いを定義します。
type Persistable interface {
	GetParsistPath() string
	Load(entity Persistable) (Persistable, error)
}

// PersistService はファイルベースのデータ永続化を担当します。
type PersistService[T Persistable] struct{}

func (s *Persistable) Load(entyt T) (T, error) {
	var output T
	return output, nil
}

// Load は永続化ファイルからデータを読み込みます。
func (s *PersistService[T]) Load(entity T) (T, error) {
	var output T

	// ファイルを読み込み
	data, err := os.ReadFile(entity.GetParsistPath())
	if err != nil {
		return output, err
	}

	if err := yaml.Unmarshal(data, &output); err != nil {
		return output, err
	}

	return output, nil
}

// Save はデータを永続化ファイルに保存します。
func (s *PersistService[T]) Save(entity T) error {
	// ディレクトリが存在しない場合はエラー
	targetFolder := entity.GetParsistPath()
	fi, err := os.Stat(targetFolder)
	if os.IsNotExist(err) {
		return fmt.Errorf("ターゲットフォルダーが存在しません: %s", targetFolder)
	}
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return fmt.Errorf("ターゲットフォルダーがディレクトリではありません: %s", targetFolder)
	}

	// データをYAMLにエンコード
	yamlText, err := yaml.Marshal(entity)
	if err != nil {
		return fmt.Errorf("データのエンコードに失敗しました: %w", err)
	}

	// ファイルに書き込み
	persistPath := s.PersistPath(entity)
	return os.WriteFile(persistPath, yamlText, 0644)
}
