package services

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type DatabaseService[T any] struct {
	DatabasePath string `json:"database_path" yaml:"database_path" example:"/home/<user>/penguin/database.yaml"`
}

func NewDatabaseService[T any](databasePath string) (*DatabaseService[T], error) {
	// Expand ~ to home directory
	if strings.HasPrefix(databasePath, "~/") {
		usr, err := user.Current()
		if err != nil {
			return nil, err
		}
		databasePath = filepath.Join(usr.HomeDir, databasePath[2:])
	}

	// Get absolute path
	databasePath, err := filepath.Abs(databasePath)
	if err != nil {
		return nil, err
	}

	return &DatabaseService[T]{
		DatabasePath: databasePath,
	}, nil
}

func (s *DatabaseService[T]) LoadFromYAML() (T, error) {

	// Initialize output with default data
	var output T

	// Check if file exists
	if _, err := os.Stat(s.DatabasePath); os.IsNotExist(err) {
		return output, errors.New("file does not exist")
	}

	// Read YAML file
	yamlData, err := os.ReadFile(s.DatabasePath)
	if err != nil {
		return output, err
	}

	if err := yaml.Unmarshal(yamlData, &output); err != nil {
		return output, err
	}

	return output, nil
}

func (s *DatabaseService[T]) SaveToYAML(input T) error {
	yamlData, err := yaml.Marshal(input)
	if err != nil {
		return err
	}

	return os.WriteFile(s.DatabasePath, yamlData, 0644)

}
