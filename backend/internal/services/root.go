package services

type RootHandler struct {
	FileService    *FileServiceHandler
	CompanyService *CompanyServiceHandler
	KojiService    *KojiServiceHandler
}
