package services

// Services は各サービスのハンドラーをまとめた構造体です。
type Services struct {
	FileService    *FileService
	CompanyService *CompanyServiceOld
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
	companyService := &CompanyServiceOld{}
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
