package app

import "penguin-backend/internal/services"

type ServiceOptions struct {
	FileFolderPath     string
	BusinessFolderPath string
	CompanyFolderPath  string
	KojiFolderPath     string
	DatabaseFilename   string
}

var DefaultServiceOptions = ServiceOptions{
	FileFolderPath:     "~/penguin",
	BusinessFolderPath: "~/penguin/豊田築炉",
	CompanyFolderPath:  "~/penguin/豊田築炉/1 会社",
	KojiFolderPath:     "~/penguin/豊田築炉/2 工事",
	DatabaseFilename:   ".detail.yaml",
}

type ServiceContainer struct {
	Root     *services.RootService
	File     *services.FileService
	Company  *services.CompanyService
	Koji     *services.KojiService
	Business *services.BusinessService
}

func SetupServices(opts ServiceOptions) ServiceContainer {
	root := services.CreateRootService()

	fileSvc := &services.FileService{}
	root.AddService(fileSvc, services.ConfigPathName(opts.FileFolderPath))

	companySvc := &services.CompanyService{}
	root.AddService(companySvc,
		services.ConfigPathName(opts.CompanyFolderPath),
		services.ConfigFileName(opts.DatabaseFilename),
	)

	kojiSvc := &services.KojiService{}
	root.AddService(kojiSvc,
		services.ConfigPathName(opts.KojiFolderPath),
		services.ConfigFileName(opts.DatabaseFilename),
	)

	businessSvc := &services.BusinessService{}
	root.AddService(businessSvc)

	return ServiceContainer{
		Root:     root,
		File:     fileSvc,
		Company:  companySvc,
		Koji:     kojiSvc,
		Business: businessSvc,
	}
}
