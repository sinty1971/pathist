package services

// RootService は全てのサービスを管理するルートサービス
// このサービスに登録されているサービスはこのサービスを経由して他のサービスを利用出来るようになります。
type RootService struct {
	// ファイルサービス
	FileService *FileService

	// 工事サービス
	KojiService *KojiService

	// 会社管理サービス
	CompanyService *CompanyService
}

// GetRootService はRootServiceを返す
func (rs *RootService) GetRootService() *RootService {
	return rs
}

func (rs *RootService) Register(root *RootService, opts ...Option) {
	rs.FileService = &FileService{}
	rs.KojiService = &KojiService{}
	rs.CompanyService = &CompanyService{}
}

// CreateRootService はRootServiceを初期化する
func CreateRootService() *RootService {
	rs := &RootService{
		FileService:    &FileService{},
		KojiService:    &KojiService{},
		CompanyService: &CompanyService{},
	}

	return rs
}

// Cleanup はリソースのクリーンアップを行う
func (rs *RootService) Cleanup() error {
	// 必要に応じて各サービスのクリーンアップ処理を実行
	// 例: データベース接続のクローズ、ファイルハンドルのクローズなど
	return nil
}
