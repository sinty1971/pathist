package services

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Persistable はファイルに永続化されるエンティティに必要な振る舞いを定義します。
type Persistable interface {
	GetTargetFolder() string
}

// PersistService はファイルベースのデータ永続化を担当します。
type PersistService[T Persistable] struct {
	persistFilename string
}

// NewPersistService はデータ永続化サービスを初期化します。
func NewPersistService[T Persistable](persistFilename string) *PersistService[T] {
	return &PersistService[T]{persistFilename}
}

// PersistPath は永続化ファイルのフルパスを返します。
func (rs *PersistService[T]) PersistPath(entity T) string {
	return filepath.Join(entity.GetTargetFolder(), rs.persistFilename)
}

// Load は永続化ファイルからデータを読み込みます。
func (rs *PersistService[T]) Load(entity T) (T, error) {
	var output T

	data, err := os.ReadFile(
		rs.PersistPath(entity),
	)
	if err != nil {
		return output, err
	}

	if err := yaml.Unmarshal(data, &output); err != nil {
		return output, err
	}

	return output, nil
}

// Save はデータを永続化ファイルに保存します。
func (rs *PersistService[T]) Save(entity T) error {
	// ディレクトリが存在しない場合はエラー
	targetFolder := entity.GetTargetFolder()
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
	persistPath := rs.PersistPath(entity)
	return os.WriteFile(persistPath, yamlText, 0644)
}
