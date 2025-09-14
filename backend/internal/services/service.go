package services

// Service はRootServiceに登録可能なサービスのインターフェースを定義します
// インターフェイスメソッドは通常ポインタレシーバーで実装されます
type Service interface {
	GetServiceName() string
	GetService(serviceName string) *Service
	Cleanup() error
	Initialize(rs *RootService, opts ...ConfigFunc) error
}

// ConfigFunc はサービス登録時のオプション設定用の関数型です
type ConfigFunc func(*Config)

// Config はサービス登録時の設定を保持します
type Config struct {
	FileName     string
	PathName     string
	Priority     int
	LazyLoad     bool
	Dependencies []string
}

// ConfigFileName はファイル名を設定するコンフィグレーション関数です
func ConfigFileName(fileName string) ConfigFunc {
	return func(c *Config) {
		c.FileName = fileName
	}
}

// ConfigPathName はパス名を設定するコンフィグレーション関数です
func ConfigPathName(pathName string) ConfigFunc {
	return func(c *Config) {
		c.PathName = pathName
	}
}

// ConfigPriority はサービスの優先度を設定するコンフィグレーション関数です
func ConfigPriority(priority int) ConfigFunc {
	return func(c *Config) {
		c.Priority = priority
	}
}

// ConfigLazyLoad は遅延ロードを有効にするコンフィグレーション関数です
func ConfigLazyLoad(enabled bool) ConfigFunc {
	return func(c *Config) {
		c.LazyLoad = enabled
	}
}

// ConfigDependencies はサービスの依存関係を設定するコンフィグレーション関数です
func ConfigDependencies(deps ...string) ConfigFunc {
	return func(c *Config) {
		c.Dependencies = append(c.Dependencies, deps...)
	}
}

// cfs から Config を作成する
func NewConfig(cfs ...ConfigFunc) *Config {
	config := &Config{}
	for _, cf := range cfs {
		cf(config)
	}
	return config
}
