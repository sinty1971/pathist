package services

import (
	"context"
	"errors"
	"log"
	"os"
	"path/filepath"

	grpcv1 "grpc-backend/gen/grpc/v1"
	grpcv1connect "grpc-backend/gen/grpc/v1/grpcv1connect"
	"grpc-backend/internal/models"

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

	// companiesById は管理されている会社データのインデックスがIdのキャッシュマップ
	companiesById map[string]*models.Company
}

// NewCompanyService CompanyService インスタンスを作成します
func NewCompanyService(services *Services, managedFolder string) (*CompanyService, error) {

	// インスタンス作成
	cs := &CompanyService{
		services:      services,
		managedFolder: managedFolder,
		companiesById: make(map[string]*models.Company, 1000),
	}

	// companiesの情報を取得
	if err := cs.UpdateCompanies(); err != nil {
		return nil, err
	}

	// managedFolderの監視を開始
	if err := watchManagedFolder(managedFolder); err != nil {
		return nil, err
	}

	return cs, nil
}

// UpdateCompanies ファイルシステムから会社データを再読み込みします
func (s *CompanyService) UpdateCompanies() error {

	// ファイルシステムから会社フォルダー一覧を取得
	entries, err := os.ReadDir(s.managedFolder)
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		return errors.New("会社フォルダーが1つも存在しません")
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
func watchManagedFolder(managedFolder string) error {
	absPath, err := filepath.Abs(managedFolder)
	if err != nil {
		return err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	// 監視終了時に閉じる
	go func() {
		<-watcher.Errors
		watcher.Close()
	}()

	// イベントループ
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Printf("[managed-folder] event=%s path=%s", event.Op, event.Name)

				if event.Op&(fsnotify.Create|fsnotify.Remove|fsnotify.Rename|fsnotify.Write) != 0 {
					// 必要に応じてサービスへ通知する
					// 例: reload metadata, update cache, etc.
				}

			case err := <-watcher.Errors:
				log.Printf("[managed-folder] watcher error: %v", err)
			}
		}
	}()

	// フォルダを監視対象に追加
	if err := watcher.Add(absPath); err != nil {
		return err
	}

	log.Printf("watching managed folder: %s", absPath)
	return nil
}

func (s *CompanyService) ListCompanies(
	ctx context.Context, _ *connect.Request[grpcv1.ListCompaniesRequest]) (
	response *connect.Response[grpcv1.ListCompaniesResponse], err error) {

	// 会社データモデルを作成
	grpcv1Companies := make([]*grpcv1.Company, len(s.companiesById))
	count := 0
	for _, v := range s.companiesById {
		grpcv1Companies[count] = v.Company
		count++
	}

	// Responseの更新とリターン
	response.Msg.SetCompanies(grpcv1Companies)
	return response, nil
}

// GetCompany は会社IDから会社情報を取得します
func (s *CompanyService) GetCompany(
	ctx context.Context, req *connect.Request[grpcv1.GetCompanyRequest]) (
	response *connect.Response[grpcv1.GetCompanyResponse], err error) {

	// Idの取得
	id := req.Msg.GetId()

	// 会社情報を取得
	company, exist := s.companiesById[id]
	if !exist {
		return response, connect.NewError(connect.CodeNotFound, errors.New("company not found"))
	}
	response.Msg.SetCompany(company.Company)

	return response, nil
}

func (h *CompanyService) UpdateCompany(ctx context.Context, req *connect.Request[grpcv1.UpdateCompanyRequest]) (*connect.Response[grpcv1.UpdateCompanyResponse], error) {
	if req.Msg.GetCompany() == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("company is required"))
	}

	company := convertProtoCompany(req.Msg.GetCompany())
	if company.GetManagedFolder() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("target_folder is required"))
	}

	if err := h.companyService.Update(company); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	reloaded, err := h.companyService.GetCompany(company.GetManagedFolder())
	if err == nil {
		company = reloaded
	}

	res := grpcv1.UpdateCompanyResponse_builder{
		Company: convertModelCompany(company),
	}.Build()

	return connect.NewResponse(res), nil
}

func (h *CompanyService) ListCompanyCategories(ctx context.Context, _ *connect.Request[grpcv1.ListCompanyCategoriesRequest]) (*connect.Response[grpcv1.ListCompanyCategoriesResponse], error) {
	categories := h.companyService.Categories()
	items := make([]*grpcv1.CompanyCategoryInfo, 0, len(categories))
	for _, category := range categories {
		items = append(items, convertModelCompanyCategory(category))
	}

	res := grpcv1.ListCompanyCategoriesResponse_builder{
		Categories: items,
	}.Build()

	return connect.NewResponse(res), nil
}

func convertModelCompanyCategory(src models.CompanyCategory) *grpcv1.CompanyCategoryInfo {
	return grpcv1.CompanyCategoryInfo_builder{
		Code:  string(src.Index),
		Label: src.Label,
	}.Build()
}
