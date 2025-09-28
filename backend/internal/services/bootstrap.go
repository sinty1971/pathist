package services

import (
	"os"
	"path/filepath"
)

type Container struct {
	FileService    *FileService
	CompanyService *CompanyService
	KojiService    *KojiService
}

type Options struct {
	FileServiceTargetFolder    string
	CompanyServiceTargetFolder string
	KojiServiceTargetFolder    string
	PersistFilename            string
}

var DefaultOptions = Options{
	PersistFilename:            ".detail.yaml",
	FileServiceTargetFolder:    "~/penguin",
	CompanyServiceTargetFolder: "~/penguin/豊田築炉/1 会社",
	KojiServiceTargetFolder:    "~/penguin/豊田築炉/2 工事",
}

// ResolveOptions は環境変数などを考慮したサービスオプションを生成します。
func ResolveOptions() Options {
	opts := DefaultOptions
	if dataRoot := os.Getenv("PENGUIN_DATA_ROOT"); dataRoot != "" {
		opts.FileServiceTargetFolder = dataRoot
		opts.CompanyServiceTargetFolder = filepath.Join(dataRoot, "豊田築炉", "1 会社")
		opts.KojiServiceTargetFolder = filepath.Join(dataRoot, "豊田築炉", "2 工事")
	}
	return opts
}

// CreateContainer は与えられたオプションでサービス群を初期化します。
func CreateContainer() (*Container, error) {
	// コンテナを作成
	container := &Container{}
	var err error

	// FileServiceを初期化
	fileService := &FileService{}
	container.FileService, err = fileService.Initialize(container, &DefaultOptions)
	if err != nil {
		return nil, err
	}

	// CompanyServiceを初期化
	companyService := &CompanyService{}
	container.CompanyService, err = companyService.Initialize(container, &DefaultOptions)
	if err != nil {
		return nil, err
	}

	// KojiServiceを初期化
	kojiService := &KojiService{}
	container.KojiService, err = kojiService.Initialize(container, &DefaultOptions)
	if err != nil {
		return nil, err
	}

	return container, nil
}

// Cleanup はサービスをクリーンアップする
func (c *Container) Cleanup() {
	if c.FileService != nil {
		c.FileService.Cleanup()
	}
	if c.CompanyService != nil {
		c.CompanyService.Cleanup()
	}
	if c.KojiService != nil {
		c.KojiService.Cleanup()
	}
}
