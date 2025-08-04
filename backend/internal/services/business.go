package services

// BusinessService は統合ビジネスデータ管理サービス
type BusinessService struct {
	// RootService はトップコンテナのインスタンス
	RootService *ContainerService `json:"-" yaml:"-"`

	// ビジネスフォルダーフルパス名（Ex. "~/penguin/豊田築炉"）
	FolderPath string `json:"folderPath" yaml:"folder-path"`

	// データベースサービス
	DatabaseService *DatabaseFileService[*BusinessService] `json:"-" yaml:"-"`

	// 企業データ管理サービス
	CompanyService *CompanyService `json:"-" yaml:"-"`

	// 工事データ管理サービス
	KojiService *KojiService `json:"-" yaml:"-"`

	// 将来追加予定
	// メンバーデータ管理サービス
	// MemberService  *MemberService

	// 注文データ管理サービス
	// OrderService   *OrderService

	// 請求データ管理サービス
	// InvoiceService *InvoiceService

	// 見積りデータ管理サービス
	// EstimateService *EstimateService

	// 設定
	DatabaseFilename string `json:"databaseFilename" yaml:"database-filename"`
}

// GetFolderPath Databaseインターフェースの実装
func (bs *BusinessService) GetFolderPath() string {
	return bs.FolderPath
}

// BuildWithOption はBusinessServiceを初期化する
// folderPath は会社フォルダーパス(Ex. "~/penguin/豊田築炉")
// databaseFilename はデータベースファイルのファイル名
func (bs *BusinessService) BuildWithOption(opt ContainerOption, folderPath string) {
	// コンテナを設定
	bs.RootService = opt.RootService

	// ビジネスフォルダーパスを設定
	bs.FolderPath = folderPath
}

// 将来追加予定のメソッド
/*
// GetMembers は社員一覧を取得する
func (s *BusinessService) GetMembers() ([]models.Member, error) {
	return s.MemberService.GetAllMembers()
}

// GetOrders は注文一覧を取得する
func (s *BusinessService) GetOrders() ([]models.Order, error) {
	return s.OrderService.GetAllOrders()
}

// GetInvoices は請求一覧を取得する
func (s *BusinessService) GetInvoices() ([]models.Invoice, error) {
	return s.InvoiceService.GetAllInvoices()
}

// GetEstimates は見積一覧を取得する
func (s *BusinessService) GetEstimates() ([]models.Estimate, error) {
	return s.EstimateService.GetAllEstimates()
}
*/
