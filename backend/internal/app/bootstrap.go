package app

import "penguin-backend/internal/services"

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
