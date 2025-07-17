package services

import (
	"fmt"
	"os"
	"penguin-backend/internal/models"

	"gopkg.in/yaml.v3"
)

// DatabaseService はYAML形式のファイルでフォルダーパスからは得られない情報をファイルで保持します
type DatabaseService[T models.Database] struct {
	FileService *FileService
	Filename    string
}

// NewDatabaseService はDatabaseServiceを初期化する
// @Param fileService query string true "ファイルサービス"
// @Param filename query string true "データベースファイル名(.detail.yaml)"
func NewDatabaseService[T models.Database](fileService *FileService, filename string) *DatabaseService[T] {
	return &DatabaseService[T]{
		FileService: fileService,
		Filename:    filename,
	}
}

// Load はYAMLファイルからデータを読み込む
// @Param folderName query string true "フォルダー名(FileService.BasePathからの相対パス)"
func (as *DatabaseService[T]) Load(ref T) (T, error) {

	// Initialize output with default data
	var output T

	// データベースファイルのフルパスを取得
	databasePath, err := as.FileService.GetFullpath(ref.GetFolderName(), as.Filename)
	if err != nil {
		return output, fmt.Errorf("データベースファイルのフルパスの取得に失敗しました: %v", err)
	}

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
func (as *DatabaseService[T]) Save(input T) error {

	// データベースファイルのフルパスを取得
	databasePath, err := as.FileService.GetFullpath(input.GetFolderName(), as.Filename)
	if err != nil {
		return fmt.Errorf("データベースファイルのフルパスの取得に失敗しました: %v", err)
	}

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
