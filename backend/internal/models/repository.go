package models

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// Repository はファイルベースのデータ永続化を行うリポジトリです
// 現在はYAML形式をサポートしていますが、将来的に他の形式にも対応可能です
type Repository[T Persistable] struct {
	persistPath string
}

// Persistable はファイルベースで永続化可能なエンティティのインターフェースを定義します
type Persistable interface {
	// GetPersistPath はデータが保存されているパス名を取得します
	GetPersistPath() string
}

// NewRepository はFileRepositoryを初期化する
func NewRepository[T Persistable](filePath string) *Repository[T] {
	return &Repository[T]{
		persistPath: filePath,
	}
}

// Load はYAMLファイルからデータを読み込む
// @Param folderName query string true "フォルダー名(FileService.BasePathからの相対パス)"
func (r *Repository[T]) Load(ref T) (T, error) {

	// Initialize output with default data
	var output T

	// データファイルを読み込む
	yamlData, err := os.ReadFile(r.persistPath)
	if err != nil {
		return output, err
	}

	// データファイルをデコード
	if err := yaml.Unmarshal(yamlData, &output); err != nil {
		return output, err
	}

	return output, nil
}

// Save はデータをファイルに保存する
// @Param folderName query string true "フォルダー名(FileService.BasePathからの相対パス)"
func (r *Repository[T]) Save(input T) error {

	// データをエンコード
	yamlData, err := yaml.Marshal(input)
	if err != nil {
		return fmt.Errorf("データのエンコードに失敗しました: %v", err)
	}

	// ファイルが存在しない場合は作成
	if _, err := os.Stat(r.persistPath); os.IsNotExist(err) {
		os.Create(r.persistPath)
	}

	return os.WriteFile(r.persistPath, yamlData, 0644)
}
