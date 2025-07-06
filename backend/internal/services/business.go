package services

// BusinessDataService は統合ビジネスデータ管理サービス
// ファイルサーバー内のファイル名から様々なビジネスデータを解析・提供する
type BusinessDataService struct {
	// 基準ファイルサービス（ファイルアクセスはこれより下位のみ）
	FileService *FileService

	// 工事データ管理サービス
	KojiService *KojiService

	// 企業データ管理サービス
	CompanyService *CompanyService

	// メンバーデータ管理サービス（将来追加予定）
	// 将来追加予定
	// MemberService  *MemberService
	// OrderService   *OrderService
	// InvoiceService *InvoiceService
	// EstimateService *EstimateService

	// 設定
	AttributeFilename string
}

// NewBusinessDataService はBusinessDataServiceを初期化する
// @Param businessFilePath query string true "基準ファイルサービスのパス" default("~/penguin/豊田築炉")
// @Param attributeFilename query string true "属性ファイルのファイル名" default(".detail.yaml")
func NewBusinessDataService(businessFilePath, attributeFilename string) (*BusinessDataService, error) {
	// fileServiceを初期化
	fileService, err := NewFileService(businessFilePath)
	if err != nil {
		return nil, err
	}

	// KojiServiceを初期化（CompanyService依存なしで）
	kojiService, err := NewKojiService(fileService, "2 工事")
	if err != nil {
		return nil, err
	}

	// CompanyServiceを初期化
	companyService, err := NewCompanyService(fileService, "1 会社")
	if err != nil {
		return nil, err
	}

	return &BusinessDataService{
		FileService:       fileService,
		KojiService:       kojiService,
		CompanyService:    companyService,
		AttributeFilename: attributeFilename,
	}, nil
}

// 将来追加予定のメソッド
/*
// GetMembers は社員一覧を取得する
func (s *BusinessDataService) GetMembers() ([]models.Member, error) {
	return s.MemberService.GetAllMembers()
}

// GetOrders は注文一覧を取得する
func (s *BusinessDataService) GetOrders() ([]models.Order, error) {
	return s.OrderService.GetAllOrders()
}

// GetInvoices は請求一覧を取得する
func (s *BusinessDataService) GetInvoices() ([]models.Invoice, error) {
	return s.InvoiceService.GetAllInvoices()
}

// GetEstimates は見積一覧を取得する
func (s *BusinessDataService) GetEstimates() ([]models.Estimate, error) {
	return s.EstimateService.GetAllEstimates()
}
*/
