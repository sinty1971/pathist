package services

// Option はサービス登録時のオプション設定用の関数型です
type Option func(*Config)

// Config はサービス登録時の設定を保持します
type Config struct {
	Name         string
	FolderPath   string
	Priority     int
	LazyLoad     bool
	Dependencies []string
}

// WithName はサービス名を設定するオプションです
func WithName(name string) Option {
	return func(c *Config) {
		c.Name = name
	}
}

func WithFolderPath(folderPath string) Option {
	return func(c *Config) {
		c.FolderPath = folderPath
	}
}

// WithPriority はサービスの優先度を設定するオプションです
func WithPriority(priority int) Option {
	return func(c *Config) {
		c.Priority = priority
	}
}

// WithLazyLoad は遅延ロードを有効にするオプションです
func WithLazyLoad(enabled bool) Option {
	return func(c *Config) {
		c.LazyLoad = enabled
	}
}

// WithDependencies はサービスの依存関係を設定するオプションです
func WithDependencies(deps ...string) Option {
	return func(c *Config) {
		c.Dependencies = append(c.Dependencies, deps...)
	}
}

// RegistableService はRootServiceに登録可能なサービスのインターフェースを定義します
type RegistableService interface {
	GetRootService() *RootService
	Register(rs *RegistableService, opts ...Option)
}
