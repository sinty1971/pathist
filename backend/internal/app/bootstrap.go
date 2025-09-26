package app

import (
	"os"
	"path/filepath"
	"penguin-backend/internal/services"
)

type ServiceOptions struct {
	FileFolderPath    string
	CompanyFolderPath string
	KojiFolderPath    string
	DatabaseFilename  string
}

var DefaultServiceOptions = ServiceOptions{
	FileFolderPath:    "~/penguin",
	CompanyFolderPath: "~/penguin/豊田築炉/1 会社",
	KojiFolderPath:    "~/penguin/豊田築炉/2 工事",
	DatabaseFilename:  ".detail.yaml",
}

type ServiceContainer struct {
	Root    *services.RootService
	File    *services.FileService
	Company *services.CompanyService
	Koji    *services.KojiService
}

// ResolveServiceOptions は環境変数などを考慮したサービスオプションを生成します。
func ResolveServiceOptions() ServiceOptions {
	opts := DefaultServiceOptions
	if dataRoot := os.Getenv("PENGUIN_DATA_ROOT"); dataRoot != "" {
		opts.FileFolderPath = dataRoot
		opts.CompanyFolderPath = filepath.Join(dataRoot, "豊田築炉", "1 会社")
		opts.KojiFolderPath = filepath.Join(dataRoot, "豊田築炉", "2 工事")
	}
	return opts
}

// InitializeServicesWithOptions は与えられたオプションでサービス群を初期化します。
func InitializeServicesWithOptions(opts ServiceOptions) (*ServiceContainer, error) {
	rootService := services.CreateRootService()

	fileService := &services.FileService{}
	if err := rootService.AddService(fileService, services.WithPath(opts.FileFolderPath)); err != nil {
		return nil, err
	}

	companyService := &services.CompanyService{}
	if err := rootService.AddService(
		companyService,
		services.WithPath(opts.CompanyFolderPath),
		services.WithFileName(opts.DatabaseFilename),
	); err != nil {
		return nil, err
	}

	kojiService := &services.KojiService{}
	if err := rootService.AddService(
		kojiService,
		services.WithPath(opts.KojiFolderPath),
		services.WithFileName(opts.DatabaseFilename),
	); err != nil {
		return nil, err
	}

	return &ServiceContainer{
		Root:    rootService,
		File:    fileService,
		Company: companyService,
		Koji:    kojiService,
	}, nil
}

// InitializeServices は環境変数を考慮したオプション生成とサービス登録をまとめて行います。
func InitializeServices() (*ServiceContainer, error) {
	return InitializeServicesWithOptions(ResolveServiceOptions())
}
