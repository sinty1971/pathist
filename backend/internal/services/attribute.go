package services

import (
	"fmt"
	"os"
	"penguin-backend/internal/models"

	"gopkg.in/yaml.v3"
)

// AttributeService はYAML形式のファイルでフォルダーパスからは得られない情報をファイルで保持します
type AttributeService[T models.Attributable] struct {
	FileService *FileService
	Filename    string
}

// NewAttributeService はAttributeServiceを初期化する
// @Param fileService query string true "ファイルサービス"
// @Param filename query string true "属性ファイル名(.detail.yaml)"
func NewAttributeService[T models.Attributable](fileService *FileService, filename string) *AttributeService[T] {
	return &AttributeService[T]{
		FileService: fileService,
		Filename:    filename,
	}
}

// Load はYAMLファイルからデータを読み込む
// @Param folderName query string true "フォルダー名(FileService.BasePathからの相対パス)"
func (as *AttributeService[T]) Load(ref T) (T, error) {

	// Initialize output with default data
	var output T

	// 属性ファイルのフルパスを取得
	attributePath, err := as.FileService.GetFullpath(ref.GetFileInfo().Name, as.Filename)
	if err != nil {
		return output, fmt.Errorf("属性ファイルのフルパスの取得に失敗しました: %v", err)
	}

	// 属性ファイルを読み込む
	yamlData, err := os.ReadFile(attributePath)
	if err != nil {
		return output, err
	}

	// 属性ファイルをデコード
	if err := yaml.Unmarshal(yamlData, &output); err != nil {
		return output, err
	}

	return output, nil
}

// Save はデータを属性ファイルに保存する
// @Param folderName query string true "フォルダー名(FileService.BasePathからの相対パス)"
func (as *AttributeService[T]) Save(input T) error {

	// 属性ファイルのフルパスを取得
	attributePath, err := as.FileService.GetFullpath(input.GetFileInfo().Name, as.Filename)
	if err != nil {
		return fmt.Errorf("属性ファイルのフルパスの取得に失敗しました: %v", err)
	}

	// データをエンコード
	yamlData, err := yaml.Marshal(input)
	if err != nil {
		return fmt.Errorf("データのエンコードに失敗しました: %v", err)
	}

	// ファイルが存在しない場合は作成
	if _, err := os.Stat(attributePath); os.IsNotExist(err) {
		os.Create(attributePath)
	}

	return os.WriteFile(attributePath, yamlData, 0644)
}
