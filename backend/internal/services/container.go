package services

// ContainerOption はコンテナのオプションを表す
type ContainerOption struct {
	RootService *ContainerService
}

// ContainerService は全てのサービスを管理するコンテナ
type ContainerService struct {
	// ファイルサービス
	FileService *FileService

	// ビジネスデータ管理サービス
	BusinessService *BusinessService

	// マルチメディア管理サービス
	MultiMediaService *MultiMediaService

	// その他のサービス（将来追加予定）
	// AuthService     *AuthService
	// NotificationService *NotificationService
}

// NewContainerService はServiceContainerを初期化する
func NewContainerService() *ContainerService {
	return &ContainerService{
		FileService:       &FileService{},
		BusinessService:   &BusinessService{},
		MultiMediaService: &MultiMediaService{},
	}
}

// Cleanup はリソースのクリーンアップを行う
func (cs *ContainerService) Cleanup() error {
	// 必要に応じて各サービスのクリーンアップ処理を実行
	// 例: データベース接続のクローズ、ファイルハンドルのクローズなど
	return nil
}
