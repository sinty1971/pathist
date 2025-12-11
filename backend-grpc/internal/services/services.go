package services

import "backend-grpc/internal/core"

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
		if err := (*s).Start(ss, &core.ConfigMap); err != nil {
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
