package services

import (
	"context"
	"errors"
	"log"
	"os"
	"path/filepath"

	grpcv1 "backend-grpc/gen/grpc/v1"
	grpcv1connect "backend-grpc/gen/grpc/v1/grpcv1connect"
	"backend-grpc/internal/models"

	"connectrpc.com/connect"
	"github.com/fsnotify/fsnotify"
)

// CompanyService bridges CompanyService logic to Connect handlers.
type CompanyService struct {
	// Embed the unimplemented handler for forward compatibility
	grpcv1connect.UnimplementedCompanyServiceHandler

	// services は任意のgrpcサービスハンドラーへの参照
	services *Services

	// managedFolder はこのサービスが管理する会社データのルートフォルダー
	managedFolder string

	// managedFolderWatcher は managedFolder のファイルシステム監視オブジェクト
	managedFolderWatcher *fsnotify.Watcher

	// companiesById は管理されている会社データのインデックスがIdのキャッシュマップ
	companiesById map[string]*models.Company
}

// NewCompanyService CompanyService インスタンスを作成します
func NewCompanyService(
	services *Services,
	options *ServiceOptions) (
	s *CompanyService,
	err error) {

	// インスタンス作成
	s = &CompanyService{
		services:      services,
		managedFolder: options.CompanyServiceManagedFolder,
		companiesById: make(map[string]*models.Company, 1000),
	}

	// companiesの情報を取得
	if err = s.UpdateCompanies(); err != nil {
		return
	}

	// managedFolderの監視を開始
	if err = s.watchManagedFolder(); err != nil {
		return
	}

	return
}

func (s *CompanyService) Cleanup() {
	// 現在はクリーンアップ処理は不要
}

// UpdateCompanies ファイルシステムから会社データを再読み込みします
func (s *CompanyService) UpdateCompanies() (err error) {

	// 変数定義
	var entries []os.DirEntry

	// ファイルシステムから会社フォルダー一覧を取得
	entries, err = os.ReadDir(s.managedFolder)
	if err != nil {
		return
	}

	// 会社データモデルを作成
	for _, entry := range entries {
		// 会社データモデルを作成、これはデータベースアクセスを行いません
		companyFolder := filepath.Join(s.managedFolder, entry.Name())
		company, err := models.NewCompany(companyFolder)
		if err != nil {
			continue
		}

		s.companiesById[company.Company.GetId()] = company
	}
	return nil
}

// watchManagedFolder starts watching the provided managedFolder for changes.
// Add callbacks or channels as needed to propagate events to your services.
func (s *CompanyService) watchManagedFolder() error {
	absPath, err := filepath.Abs(s.managedFolder)
	if err != nil {
		return err
	}

	s.managedFolderWatcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	// 監視終了時に閉じる
	go func() {
		<-s.managedFolderWatcher.Errors
		s.managedFolderWatcher.Close()
	}()

	// イベントループ
	go func() {
		for {
			select {
			case event, ok := <-s.managedFolderWatcher.Events:
				if !ok {
					return
				}
				log.Printf("[managed-folder] event=%s path=%s", event.Op, event.Name)

				if event.Op&(fsnotify.Create|fsnotify.Remove|fsnotify.Rename|fsnotify.Write) != 0 {
					// 必要に応じてサービスへ通知する
					// 例: reload metadata, update cache, etc.
				}

			case err := <-s.managedFolderWatcher.Errors:
				log.Printf("[managed-folder] watcher error: %v", err)
			}
		}
	}()

	// フォルダを監視対象に追加
	if err := s.managedFolderWatcher.Add(absPath); err != nil {
		return err
	}

	log.Printf("watching managed folder: %s", absPath)
	return nil
}

// GetCompanies は管理されている会社情報の一覧を取得します
// gRPCサービスの実装です
func (s *CompanyService) GetCompanyMapById(
	ctx context.Context,
	_ *connect.Request[grpcv1.GetCompanyMapByIdRequest]) (
	response connect.Response[grpcv1.GetCompanyMapByIdResponse],
	err error) {

	// 会社データモデルを作成
	grpcv1CompanyMapById := make(map[string]*grpcv1.Company, len(s.companiesById))
	for _, v := range s.companiesById {
		grpcv1CompanyMapById[v.Company.GetId()] = v.Company
	}

	// Responseの更新とリターン
	response.Msg.SetCompanyMapById(grpcv1CompanyMapById)
	return response, nil
}

// GetCompany は会社IDから会社情報を取得します
// gRPCサービスの実装です
func (s *CompanyService) GetCompanyById(
	ctx context.Context,
	req *grpcv1.GetCompanyByIdRequest) (
	res *grpcv1.GetCompanyByIdResponse,
	err error) {

	// Idの取得
	id := req.GetId()

	// 会社情報を取得
	company, exist := s.companiesById[id]
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
func (s *CompanyService) UpdateCompany(
	// 引数
	ctx context.Context,
	req *grpcv1.UpdateCompanyRequest) (
	// 戻り値
	res *grpcv1.UpdateCompanyResponse,
	err error) {

	// 既存の会社情報を取得
	currentCompanyId := req.GetCurrentCompanyId()
	currentCompany, exist := s.companiesById[currentCompanyId]
	if !exist {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("company not found"))
	}
	// 更新後の会社情報を取得
	grpcUpdatedCompany := req.GetUpdatedCompany()
	if grpcUpdatedCompany == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("updated company is nil"))
	}
	updatedCompany := &models.Company{
		Company: grpcUpdatedCompany,
	}

	// 会社情報を更新
	updatedCompany, err = currentCompany.Update(updatedCompany)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 会社情報のインデックスを更新
	if _, exist := s.companiesById[currentCompanyId]; exist {
		delete(s.companiesById, currentCompanyId)
		// 新しいIDで再登録
		s.companiesById[updatedCompany.GetId()] = updatedCompany
	}

	// Responseの作成
	grpcv1CompanyMapById := make(map[string]*grpcv1.Company, len(s.companiesById))
	for _, v := range s.companiesById {
		grpcv1CompanyMapById[v.Company.GetId()] = v.Company
	}
	res.SetCompanyMapById(grpcv1CompanyMapById)

	return res, nil
}

// GetCompanyCategories は業種カテゴリーの一覧を取得します
func (s *CompanyService) GetCompanyCategories(
	ctx context.Context,
	_ *grpcv1.GetCompanyCategoriesRequest) (
	res *grpcv1.GetCompanyCategoriesResponse,
	err error) {

	//
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
