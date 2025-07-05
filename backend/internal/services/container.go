package services

// ServiceContainer は全てのサービスを管理するコンテナ
type ServiceContainer struct {
	// 基盤サービス
	FileService *FileService

	// ビジネスデータ管理サービス
	BusinessService *BusinessDataService

	// メディアデータ管理サービス（将来実装）
	// MediaService *MediaDataService

	// その他のサービス（将来追加予定）
	// AuthService     *AuthService
	// NotificationService *NotificationService
}

// GetBusinessDataService はBusinessDataServiceを取得する
func (sc *ServiceContainer) GetBusinessDataService() *BusinessDataService {
	return sc.BusinessService
}

// GetMediaDataService はMediaDataServiceを取得する（将来実装）
// func (c *ServiceContainer) GetMediaDataService() *MediaDataService {
//     return c.MediaData
// }

// Cleanup はリソースのクリーンアップを行う
func (sc *ServiceContainer) Cleanup() error {
	// 必要に応じて各サービスのクリーンアップ処理を実行
	// 例: データベース接続のクローズ、ファイルハンドルのクローズなど
	return nil
}
