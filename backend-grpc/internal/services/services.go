package services

// DefaultServicesOptionMap はサービスのデフォルト設定を定義します。
// これらの値は、環境変数が設定されていない場合に使用されます。
var DefaultServicesOptionMap = map[string]string{
	"FileServiceFolder":           "~/penguin",
	"CompanyServiceManagedFolder": "~/penguin/豊田築炉/1 会社",
	"KojiServiceManagedFolder":    "~/penguin/豊田築炉/2 工事",
	"PersistCompanyFilename":      "@inside.yaml",
	"PersistKojiFilename":         "@koji.yaml",
	"PersistMemberFilename":       "@member.yaml",
}

// Sevice は各サービスが実装すべきインターフェースを定義します。
type Sevice interface {
	Start(*Services, *map[string]string) error
	Cleanup()
}

// Services は各サービスのハンドラーをまとめた構造体です。
type Services struct {
	ServiceMap map[string]*Sevice
}

// NewServices は与えられたオプションでサービス群を初期化します。
func NewServices() *Services {

	// 変数宣言
	services := &Services{}
	services.ServiceMap = make(map[string]*Sevice)
	return services
}

// AddService はサービスを追加する
func (ss *Services) AddService(serviceName string, service Sevice) {
	ss.ServiceMap[serviceName] = &service
}

// StartAll はすべてのサービスを起動する
func (ss *Services) StartAll() error {
	for _, s := range ss.ServiceMap {
		if err := (*s).Start(ss, &DefaultServicesOptionMap); err != nil {
			return err
		}
	}
	return nil
}

// CleanupAll はサービスをクリーンアップする
func (ss *Services) CleanupAll() {
	for _, srv := range ss.ServiceMap {
		(*srv).Cleanup()
	}
}
