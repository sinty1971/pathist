package services

// Services は各サービスのハンドラーをまとめた構造体です。
type Services struct {
	FileService    *FileService
	CompanyService *CompanyService
	KojiService    *KojiService
}

// Options はサービスの初期化に使用される設定を保持します。
// これらの設定は環境変数やデフォルト値に基づいて決定されます。
type Options struct {
	FileServiceFolder           string
	CompanyServiceManagedFolder string
	KojiServiceManagedFolder    string
	PersistFilename             string
}

// DefaultOptions はサービスのデフォルト設定を定義します。
// これらの値は、環境変数が設定されていない場合に使用されます。
var DefaultOptions = Options{
	PersistFilename:             "@inside.yaml",
	FileServiceFolder:           "~/penguin",
	CompanyServiceManagedFolder: "~/penguin/豊田築炉/1 会社",
	KojiServiceManagedFolder:    "~/penguin/豊田築炉/2 工事",
}

// NewServices は与えられたオプションでサービス群を初期化します。
func NewServices() (*Services, error) {

	// 変数宣言
	var (
		services = &Services{}
		err      error
	)

	// FileServiceのインスタンス作成
	services.FileService = NewFileService(services, &DefaultOptions)

	// CompanyServiceを初期化
	services.CompanyService, err = NewCompanyService(services, &DefaultOptions)
	if err != nil {
		return nil, err
	}

	// KojiServiceを初期化
	kojiService := &KojiService{}
	services.KojiService, err = kojiService.Initialize(services, &DefaultOptions)
	if err != nil {
		return nil, err
	}

	return services, nil
}

// Cleanup はサービスをクリーンアップする
func (s *Services) Cleanup() {
	if s.FileService != nil {
		s.FileService.Cleanup()
	}
	if s.CompanyService != nil {
		s.CompanyService.Cleanup()
	}
	if s.KojiService != nil {
		s.KojiService.Cleanup()
	}
}
