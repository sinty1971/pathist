package services

import (
	"context"
	"errors"
	"log"
	"os"
	"path/filepath"

	grpcv1 "backend-grpc/gen/grpc/v1"
	grpcv1connect "backend-grpc/gen/grpc/v1/grpcv1connect"
	"backend-grpc/internal/ext"
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
	managedWatchCtx      context.Context
	managedWatchCancel   context.CancelFunc

	// cacheById は管理されている会社データのインデックスがIdのキャッシュマップ
	cacheById map[string]*models.Company
}

// Start は CompanyService を初期化して開始します
func (srv *CompanyService) Start(services *Services, options *map[string]string) error {
	// パスをの取得と正規化
	optManagedFolder, exists := (*options)["CompanyServiceManagedFolder"]
	if !exists {
		return errors.New("CompanyServiceManagedFolder option is required")
	}
	managedFolder, err := ext.NormalizeAbsPath(optManagedFolder)
	if err != nil {
		return err
	}

	// 既存インスタンスに値をセット（再代入しないこと）
	srv.services = services
	srv.managedFolder = managedFolder
	srv.cacheById = make(map[string]*models.Company, 1000)
	srv.managedWatchCtx, srv.managedWatchCancel = context.WithCancel(context.Background())

	// companiesの情報を取得
	if err = srv.UpdateCompanies(); err != nil {
		return err
	}

	// managedFolderの監視を開始
	if err = srv.watchManagedFolder(srv.managedWatchCtx); err != nil {
		return err
	}

	return nil
}

func (srv *CompanyService) Cleanup() {
	if srv.managedWatchCancel != nil {
		srv.managedWatchCancel()
	}
	if srv.managedFolderWatcher != nil {
		_ = srv.managedFolderWatcher.Close()
	}
}

// UpdateCompanies ファイルシステムから会社データを再読み込みします
func (srv *CompanyService) UpdateCompanies() (err error) {

	// 変数定義
	var entries []os.DirEntry

	// ファイルシステムから会社フォルダー一覧を取得
	entries, err = os.ReadDir(srv.managedFolder)
	if err != nil {
		return
	}

	// キャッシュを作り直す（削除済み会社を残さない）
	srv.cacheById = make(map[string]*models.Company, len(entries))

	// 会社データモデルを作成
	for _, entry := range entries {
		// ディレクトリ以外はスキップ
		if !entry.IsDir() {
			continue
		}
		// 会社データモデルを作成、これはデータベースアクセスを行いません
		managedFolder := filepath.Join(srv.managedFolder, entry.Name())
		if company, err := models.NewCompany(managedFolder); err == nil {
			srv.cacheById[company.Company.GetId()] = company
		}
	}

	// 会社の内部情報の取得
	for _, company := range srv.cacheById {
		ps := ext.CreatePersistService(company)
		if err := ps.LoadPersistInfo(); err != nil {
			log.Printf("Failed to load persist info for company ID %s: %v", company.GetId(), err)
		}
	}
	return nil
}

// watchManagedFolder は指定された managedFolder の変更を監視します。
// 必要に応じてイベントをサービスに伝播するためのコールバックやチャネルを追加してください。
func (srv *CompanyService) watchManagedFolder(ctx context.Context) error {
	// 変数定義
	var err error

	// fsnotify ウォッチャーを作成
	srv.managedFolderWatcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	// イベントループ（イベント・エラーを単一ゴルーチンで処理）
	go func() {
		// 終了時にウォッチャーをクローズ
		defer srv.managedFolderWatcher.Close()

		for {
			log.Println("Waiting for managed folder events...")
			select {
			case event, ok := <-srv.managedFolderWatcher.Events:
				log.Printf("[managed-folder] event=%s path=%s", event.Op, event.Name)
				if !ok {
					return
				}

				// 新しいディレクトリが作成された場合、監視対象に追加
				if event.Op&fsnotify.Create != 0 {
					if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
						if err := srv.addWatchRecursively(event.Name); err != nil {
							log.Printf("Failed to add new directory recursively: %v", err)
						} else {
							log.Printf("Added new directory to watch: %s", event.Name)
						}
					}
				}

				if event.Op&(fsnotify.Create|fsnotify.Remove|fsnotify.Rename|fsnotify.Write) != 0 {
					// @company.yamlファイルの変更を検出
					if filepath.Base(event.Name) == "@company.yaml" {
						log.Printf("Company YAML file changed: %s", event.Name)
						if err := srv.UpdateCompanies(); err != nil {
							log.Printf("Error updating companies: %v", err)
						}
					}
				}

			case err, ok := <-srv.managedFolderWatcher.Errors:
				log.Printf("[managed-folder] watcher error: %v", err)
				if !ok {
					return
				}

			case <-ctx.Done():
				log.Printf("[managed-folder] stop watching (context canceled)")
				return
			}
		}
	}()

	// ルートフォルダと全サブディレクトリを監視対象に追加
	if err := srv.addWatchRecursively(srv.managedFolder); err != nil {
		return err
	}

	log.Printf("watching managed folder (recursively): %s", srv.managedFolder)
	return nil
}

// addWatchRecursively は指定ディレクトリとそのサブディレクトリ・ファイルを再帰的に監視対象に追加します
func (srv *CompanyService) addWatchRecursively(path string) error {
	// 現在のディレクトリを監視対象に追加
	if err := srv.managedFolderWatcher.Add(path); err != nil {
		return err
	}

	// サブディレクトリとファイルを取得
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	// 各エントリを処理
	for _, entry := range entries {
		entryPath := filepath.Join(path, entry.Name())

		if entry.IsDir() {
			// サブディレクトリを再帰的に監視対象に追加
			if err := srv.addWatchRecursively(entryPath); err != nil {
				log.Printf("Failed to watch directory %s: %v", entryPath, err)
				// エラーがあっても続行
			}
		}
	}

	return nil
}

// GetCompanies は管理されている会社情報の一覧を取得します
// gRPCサービスの実装です
func (srv *CompanyService) GetCompanyMapById(
	ctx context.Context,
	_ *grpcv1.GetCompanyMapByIdRequest) (
	res *grpcv1.GetCompanyMapByIdResponse,
	err error) {

	// レスポンスを初期化
	res = &grpcv1.GetCompanyMapByIdResponse{}

	// 会社データモデルを作成
	grpcv1CompanyMapById := make(map[string]*grpcv1.Company, len(srv.cacheById))
	for _, v := range srv.cacheById {
		grpcv1CompanyMapById[v.Company.GetId()] = v.Company
	}

	// Responseの更新とリターン
	res.SetCompanyMapById(grpcv1CompanyMapById)
	return res, nil
}

// GetCompany は会社IDから会社情報を取得します
// gRPCサービスの実装です
func (srv *CompanyService) GetCompanyById(
	ctx context.Context,
	req *grpcv1.GetCompanyByIdRequest) (
	res *grpcv1.GetCompanyByIdResponse,
	err error) {

	// レスポンスを初期化
	res = &grpcv1.GetCompanyByIdResponse{}

	// Idの取得
	id := req.GetId()

	// 会社情報を取得
	company, exist := srv.cacheById[id]
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
	// 引数
	ctx context.Context,
	req *grpcv1.UpdateCompanyRequest) (
	// 戻り値
	res *grpcv1.UpdateCompanyResponse,
	err error) {

	// 既存の会社情報を取得
	currentCompanyId := req.GetCurrentCompanyId()
	currentCompany, exist := srv.cacheById[currentCompanyId]
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
	if _, exist := srv.cacheById[currentCompanyId]; exist {
		delete(srv.cacheById, currentCompanyId)
		// 新しいIDで再登録
		srv.cacheById[updatedCompany.GetId()] = updatedCompany
	}

	// レスポンスを初期化
	res = &grpcv1.UpdateCompanyResponse{}

	// Responseの作成
	grpcv1CompanyMapById := make(map[string]*grpcv1.Company, len(srv.cacheById))
	for _, v := range srv.cacheById {
		grpcv1CompanyMapById[v.Company.GetId()] = v.Company
	}
	res.SetCompanyMapById(grpcv1CompanyMapById)

	return res, nil
}

// GetCompanyCategories は業種カテゴリーの一覧を取得します
func (srv *CompanyService) GetCompanyCategories(
	ctx context.Context,
	_ *grpcv1.GetCompanyCategoriesRequest) (
	res *grpcv1.GetCompanyCategoriesResponse,
	err error) {

	// レスポンスを初期化
	res = &grpcv1.GetCompanyCategoriesResponse{}

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
