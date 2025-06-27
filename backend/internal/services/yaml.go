package services

import (
	"os"

	"gopkg.in/yaml.v3"
)

// YamlService is a service for managing YAML files
type YamlService[T any] struct {
	FileService *FileService
}

// NewYamlService creates a new YamlService
func NewYamlService[T any](fileService *FileService) *YamlService[T] {
	return &YamlService[T]{
		FileService: fileService,
	}
}

// LoadFromYAML loads the database from a YAML file
// @Param yamlPath query string true "YAMLファイルのパス(FileSystem.Rootからの相対パス)"
func (s *YamlService[T]) Load(yamlPath string) (T, error) {

	// Initialize output with default data
	var output T

	// Read YAML file
	fullpath, err := s.FileService.GetFullpath(yamlPath)
	if err != nil {
		return output, err
	}
	yamlData, err := os.ReadFile(fullpath)
	if err != nil {
		return output, err
	}

	if err := yaml.Unmarshal(yamlData, &output); err != nil {
		return output, err
	}

	return output, nil
}

// SaveToYAML saves the database to a YAML file
func (s *YamlService[T]) Save(yamlPath string, input T) error {
	yamlData, err := yaml.Marshal(input)
	if err != nil {
		return err
	}

	fullpath, err := s.FileService.GetFullpath(yamlPath)
	if err != nil {
		return err
	}

	// ファイルが存在しない場合は作成
	if _, err := os.Stat(fullpath); os.IsNotExist(err) {
		os.Create(fullpath)
	}

	return os.WriteFile(fullpath, yamlData, 0644)
}
