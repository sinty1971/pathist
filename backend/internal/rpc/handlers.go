package rpc

type Handlers struct {
	FileService    *FileServiceHandler
	CompanyService *CompanyServiceHandler
	KojiService    *KojiServiceHandler
}
