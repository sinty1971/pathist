package services

import (
	"context"
	"errors"
	"log"
	"os"
	"path/filepath"

	grpcv1 "backend-grpc/gen/grpc/v1"
	grpcv1connect "backend-grpc/gen/grpc/v1/grpcv1connect"
	"backend-grpc/internal/core"
	"backend-grpc/internal/models"

	"connectrpc.com/connect"
)

// CompanyService の実装
type CompanyService struct {
	// Embed the unimplemented handler for forward compatibility
	grpcv1connect.UnimplementedCompanyServiceHandler

	// services は任意のgrpcサービスハンドラーへの参照
	services *Services

	// managedFolder はこのサービスが管理する会社データのルートフォルダー
	managedFolder string

	// companyMap は会社データのキャッシュマップ
	companyMap map[string]*models.Company
}

// Start は CompanyService を初期化して開始します
func (srv *CompanyService) Start(services *Services, options *map[string]string) error {
	// パスをの取得と正規化
	optManagedFolder, exists := (*options)["CompanyServiceManagedFolder"]
	if !exists {
		return errors.New("CompanyServiceManagedFolder option is required")
	}
	managedFolder, err := core.NormalizeAbsPath(optManagedFolder)
	if err != nil {
		return err
	}

	// 既存インスタンスに値をセット（再代入しないこと）
	srv.services = services
	srv.managedFolder = managedFolder
	srv.companyMap = map[string]*models.Company{}

	// companiesの情報を取得
	if err = srv.UpdateCompanyCacheMap(); err != nil {
		return err
	}

	return nil
}

func (srv *CompanyService) Cleanup() {
	// 現状クリーンアップは不要（監視を廃止したため）
}

// UpdateCompanyCacheMap 会社のキャッシュデータを更新します
func (srv *CompanyService) UpdateCompanyCacheMap() error {
	// ファイルシステムから会社フォルダー一覧を取得
	entries, err := os.ReadDir(srv.managedFolder)
	if err != nil {
		return err
	}

	// キャッシュデータの初期化
	srv.companyMap = make(map[string]*models.Company, len(entries))

	// 全てのCompanyインスタンスを作成
	for _, entry := range entries {
		// Companyインスタンスの作成と初期化
		company := models.NewCompany()
		if err := company.ParseFromManagedFolder(srv.managedFolder, entry.Name()); err == nil {
			srv.companyMap[company.GetId()] = company
		}
	}

	// 会社の内部情報の取得
	for _, company := range srv.companyMap {

		if err := company.Load(); err != nil {
			log.Printf("Failed to load persist info for company ShortName %s: %v", company.GetShortName(), err)
		}
	}
	return nil
}

// UpdateCompanyCache は指定 id のキャッシュ情報を新しい会社情報で更新します
// prevId: 更新対象の会社ID、存在しない場合は追加
// newCompany: 更新後の会社情報
func (srv *CompanyService) UpdateNewCompany(prevId string, newCompany *models.Company) (*models.Company, error) {
	// Idから更新前の会社情報を取得
	prevCompany, exist := srv.companyMap[prevId]

	// キャッシュから削除
	if exist {
		delete(srv.companyMap, prevId)
	}

	// 新しい会社情報の管理フォルダー名を生成
	newManageBaseFolder := filepath.Dir(newCompany.GetManagedFolder())
	newManagedFolder := newCompany.CreateNewManagedFolder(newManageBaseFolder, newCompany.GetCategoryIndex(), newCompany.GetShortName())

	// 管理フォルダーの変更がある場合はフォルダー移動を実施
	if exist && prevCompany.GetManagedFolder() != newCompany.GetManagedFolder() {
		// 管理フォルダーの変更がある場合はフォルダー移動を実施
		prevManagedFolder := prevCompany.GetManagedFolder()
		if err := os.Rename(prevManagedFolder, newManagedFolder); err != nil {
			return nil, err
		}
	}

	// 既存の情報を更新
	srv.companyMap[newCompany.GetId()] = newCompany

	// persist情報の書き込み
	if err := newCompany.Save(); err != nil {
		log.Printf("Failed to save persist info for company ShortName %s: %v", newCompany.GetShortName(), err)
	}

	return prevCompany, nil
}

// GetCompanies は管理されている会社情報の一覧を取得します
// gRPCサービスの実装です
func (srv *CompanyService) GetCompanyMap(
	ctx context.Context, req *grpcv1.GetCompanyMapRequest) (
	*grpcv1.GetCompanyMapResponse, error) {

	// レスポンスを初期化
	res := grpcv1.GetCompanyMapResponse_builder{}.Build()

	// 必要に応じてキャッシュを更新
	if req.GetRefresh() {
		if err := srv.UpdateCompanyCacheMap(); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	// 会社データモデルを作成
	grpcv1CompanyMap := make(map[string]*grpcv1.Company, len(srv.companyMap))
	for _, v := range srv.companyMap {
		grpcv1CompanyMap[v.Company.GetId()] = v.Company
	}

	// Responseの更新とリターン
	res.SetCompanyMap(grpcv1CompanyMap)
	return res, nil
}

// GetCompany は会社IDから会社情報を取得します
// gRPCサービスの実装です
func (srv *CompanyService) GetCompany(
	ctx context.Context,
	req *grpcv1.GetCompanyRequest) (
	res *grpcv1.GetCompanyResponse,
	err error) {

	// レスポンスを初期化
	res = grpcv1.GetCompanyResponse_builder{}.Build()

	// Idの取得
	id := req.GetId()

	// 会社情報を取得
	company, exist := srv.companyMap[id]
	if !exist {
		err = connect.NewError(connect.CodeNotFound, errors.New("company not found"))
		return
	}

	// Responseの更新
	res.SetCompany(company.Company)

	return
}

// UpdateCompany は会社情報を更新します
// gRPCサービスの実装です
// 既存の Id の会社情報を更新します。そのため Id の変更の可能性があります。
// また、フォルダーの移動も発生する可能性があります。
// Company.Id 更新対象の会社Id
func (srv *CompanyService) UpdateCompany(
	_ context.Context, req *grpcv1.UpdateCompanyRequest) (
	*grpcv1.UpdateCompanyResponse, error) {

	// リクエスト情報の取得
	prevId := req.GetPrevId()
	newCompany := &models.Company{Company: req.GetNewCompany()}

	prevCompany, err := srv.UpdateNewCompany(prevId, newCompany)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("updated company is nil"))
	}

	// Responseの作成
	res := grpcv1.UpdateCompanyResponse_builder{}.Build()
	res.SetPrevCompany(prevCompany.Company)

	return res, nil
}

// GetCompanyCategories は業種カテゴリーの一覧を取得します
func (srv *CompanyService) GetCompanyCategories(
	_ context.Context, _ *grpcv1.GetCompanyCategoriesRequest) (
	*grpcv1.GetCompanyCategoriesResponse, error) {

	// レスポンスを初期化
	res := grpcv1.GetCompanyCategoriesResponse_builder{}.Build()

	categories := make([]*grpcv1.CompanyCategory, 0, len(models.CompanyCategoryMap))
	for idx, label := range models.CompanyCategoryMap {
		categories = append(categories, grpcv1.CompanyCategory_builder{
			Index: int32(idx),
			Label: label,
		}.Build())
	}

	res.SetCategories(categories)

	return res, nil
}
