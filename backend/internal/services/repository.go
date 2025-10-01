package services

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Persistable はファイルに永続化されるエンティティに必要な振る舞いを定義します。
type Persistable interface {
	PersistFolder() string
}

// RepositoryService はファイルベースのデータ永続化を担当します。
type RepositoryService[T Persistable] struct {
	databaseFilename string
}

// NewRepositoryService はデータ永続化サービスを初期化します。
func NewRepositoryService[T Persistable](databaseFilename string) *RepositoryService[T] {
	return &RepositoryService[T]{
		databaseFilename: databaseFilename,
	}
}

// DatabaseFilename は設定されているデータベースファイル名を返します。
func (rs *RepositoryService[T]) DatabaseFilename() string {
	return rs.databaseFilename
}

// Load は永続化ファイルからデータを読み込みます。
func (rs *RepositoryService[T]) Load(entity T) (T, error) {
	var output T

	persistPath := rs.persistPath(entity)
	data, err := os.ReadFile(persistPath)
	if err != nil {
		return output, err
	}

	if err := yaml.Unmarshal(data, &output); err != nil {
		return output, err
	}

	return output, nil
}

// Save はデータを永続化ファイルに保存します。
func (rs *RepositoryService[T]) Save(entity T) error {
	data, err := yaml.Marshal(entity)
	if err != nil {
		return fmt.Errorf("データのエンコードに失敗しました: %w", err)
	}

	persistPath := rs.persistPath(entity)
	if _, err := os.Stat(persistPath); os.IsNotExist(err) {
		file, err := os.Create(persistPath)
		if err != nil {
			return fmt.Errorf("永続化ファイルの作成に失敗しました: %w", err)
		}
		_ = file.Close()
	}

	return os.WriteFile(persistPath, data, 0644)
}

// persistPath は対象エンティティの永続化ファイルパスを組み立てます。
func (rs *RepositoryService[T]) persistPath(entity T) string {
	folder := entity.PersistFolder()
	if folder == "" {
		return rs.databaseFilename
	}
	return filepath.Join(folder, rs.databaseFilename)
}
