package services

// RootService は全てのサービスを管理するルートサービス
// このサービスに登録されているサービスはこのサービスを経由して他のサービスを利用出来るようになります。
type RootService struct {
	// 登録されているサービスのマップ
	services map[string]Service
}

func (rs *RootService) GetServiceName() string {
	return "RootService"
}

func (rs *RootService) Initialize(_rs *RootService, _opts ...ConfigFunc) error {
	return nil
}

// GetService はサービスを返す
func (rs *RootService) GetService(serviceName string) Service {
	if serviceName == rs.GetServiceName() {
		return rs
	}
	if service, ok := rs.services[serviceName]; ok {
		return service
	}
	return nil
}

// Cleanup はリソースのクリーンアップを行う
func (rs *RootService) Cleanup() error {
	// 必要に応じて各サービスのクリーンアップ処理を実行
	// 例: データベース接続のクローズ、ファイルハンドルのクローズなど
	for _, service := range rs.services {
		service.Cleanup()
	}
	return nil
}

// CreateRootService はRootServiceを初期化する
func CreateRootService() *RootService {
	rs := &RootService{
		services: make(map[string]Service),
	}

	return rs
}

// AddService はサービスを追加する
// serviceName: サービス名
// service: サービスのインスタンス
// opts: サービス登録時のオプション
func (rs *RootService) AddService(service Service, opts ...ConfigFunc) error {
	if err := service.Initialize(rs, opts...); err != nil {
		return err
	}
	rs.services[service.GetServiceName()] = service
	return nil
}
