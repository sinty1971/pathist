package models

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// DatabaseInterface はデータベースのインターフェースを定義します
type DatabaseInterface interface {
	// GetFolderPath はデータが保存されているフォルダーのパスを取得します
	GetFolderPath() string
}

// Database はデータベースのインターフェースを定義します
// Filename データサービスのファイル名
type Database[T DatabaseInterface] struct {
	Filename string
}

// NewDatabaseFileService はDatabaseFileServiceを初期化する
func NewDatabaseFileService[T DatabaseInterface](filename string) *Database[T] {
	return &Database[T]{
		Filename: filename,
	}
}

// Load はYAMLファイルからデータを読み込む
// @Param folderName query string true "フォルダー名(FileService.BasePathからの相対パス)"
func (ds *Database[T]) Load(ref T) (T, error) {

	// Initialize output with default data
	var output T

	// データベースファイルのフルパスを取得
	databasePath := filepath.Join(ref.GetFolderPath(), ds.Filename)

	// データベースファイルを読み込む
	yamlData, err := os.ReadFile(databasePath)
	if err != nil {
		return output, err
	}

	// データベースファイルをデコード
	if err := yaml.Unmarshal(yamlData, &output); err != nil {
		return output, err
	}

	return output, nil
}

// Save はデータをデータベースファイルに保存する
// @Param folderName query string true "フォルダー名(FileService.BasePathからの相対パス)"
func (ds *Database[T]) Save(input T) error {

	// データベースファイルのフルパスを取得
	databasePath := filepath.Join(input.GetFolderPath(), ds.Filename)

	// データをエンコード
	yamlData, err := yaml.Marshal(input)
	if err != nil {
		return fmt.Errorf("データのエンコードに失敗しました: %v", err)
	}

	// ファイルが存在しない場合は作成
	if _, err := os.Stat(databasePath); os.IsNotExist(err) {
		os.Create(databasePath)
	}

	return os.WriteFile(databasePath, yamlData, 0644)
}
