package services

// BusinessService は複数のドメインサービスをまとめ、ビジネス系ハンドラーへ提供する集約サービスです。
type BusinessService struct {
	rootService    *RootService
	FileService    *FileService
	CompanyService *CompanyService
	KojiService    *KojiService

	FolderPath       string
	DatabaseFilename string
}

func (bs *BusinessService) GetServiceName() string {
	return "BusinessService"
}

func (bs *BusinessService) GetService(name string) *Service {
	if name == bs.GetServiceName() {
		var svc Service = bs
		return &svc
	}
	if bs.rootService == nil {
		return nil
	}
	return bs.rootService.GetService(name)
}

func (bs *BusinessService) Cleanup() error {
	return nil
}

func (bs *BusinessService) Initialize(rs *RootService, _ ...ConfigFunc) error {
	bs.rootService = rs

	if svc := rs.GetService("FileService"); svc != nil {
		if typed, ok := (*svc).(*FileService); ok {
			bs.FileService = typed
			bs.FolderPath = typed.TargetFolder
		}
	}
	if svc := rs.GetService("CompanyService"); svc != nil {
		if typed, ok := (*svc).(*CompanyService); ok {
			bs.CompanyService = typed
			if typed.DatabaseService != nil {
				bs.DatabaseFilename = typed.DatabaseService.DatabaseFilename()
			}
		}
	}
	if svc := rs.GetService("KojiService"); svc != nil {
		if typed, ok := (*svc).(*KojiService); ok {
			bs.KojiService = typed
		}
	}

	return nil
}
