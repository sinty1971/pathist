package models

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// Persistable はファイルベースで永続化可能なエンティティのインターフェースを定義します
type Persistable interface {
	// GetFolderPath はデータが保存されているフォルダーのパスを取得します
	GetFolderPath() string
}

// FileRepository はファイルベースのデータ永続化を行うリポジトリです
// 現在はYAML形式をサポートしていますが、将来的に他の形式にも対応可能です
type FileRepository[T Persistable] struct {
	filePath string
}

// NewFileRepository はFileRepositoryを初期化する
func NewFileRepository[T Persistable](filePath string) *FileRepository[T] {
	return &FileRepository[T]{
		filePath: filePath,
	}
}

// Load はYAMLファイルからデータを読み込む
// @Param folderName query string true "フォルダー名(FileService.BasePathからの相対パス)"
func (r *FileRepository[T]) Load(ref T) (T, error) {

	// Initialize output with default data
	var output T

	// データファイルを読み込む
	yamlData, err := os.ReadFile(r.filePath)
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
func (r *FileRepository[T]) Save(input T) error {

	// データをエンコード
	yamlData, err := yaml.Marshal(input)
	if err != nil {
		return fmt.Errorf("データのエンコードに失敗しました: %v", err)
	}

	// ファイルが存在しない場合は作成
	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		os.Create(r.filePath)
	}

	return os.WriteFile(r.filePath, yamlData, 0644)
}
