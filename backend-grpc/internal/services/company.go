package services

import (
	"context"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	grpcv1 "backend-grpc/gen/grpc/v1"
	grpcv1connect "backend-grpc/gen/grpc/v1/grpcv1connect"
	"backend-grpc/internal/ext"
	"backend-grpc/internal/models"

	"connectrpc.com/connect"
	"github.com/fsnotify/fsnotify"
	"github.com/gohugoio/hugo/watcher/filenotify"
)

// CompanyService の実装
type CompanyService struct {
	// Embed the unimplemented handler for forward compatibility
	grpcv1connect.UnimplementedCompanyServiceHandler

	// services は任意のgrpcサービスハンドラーへの参照
	services *Services

	// managedFolder はこのサービスが管理する会社データのルートフォルダー
	managedFolder string

	// managedFolderWatcher は managedFolder のファイルシステム監視オブジェクト
	managedFolderWatcher filenotify.FileWatcher
	managedWatchCtx      context.Context
	managedWatchCancel   context.CancelFunc

	// cachedCompanyMap は会社データのキャッシュマップ
	cachedCompanyMap map[string]*models.Company
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
	srv.cachedCompanyMap = make(map[string]*models.Company, 1000)
	srv.managedWatchCtx, srv.managedWatchCancel = context.WithCancel(context.Background())

	// companiesの情報を取得
	if err = srv.UpdateCachedCompanyMap(); err != nil {
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
		srv.managedFolderWatcher.Close()
	}
}

// UpdateCachedCompanyMap 会社のキャッシュデータを更新します
func (srv *CompanyService) UpdateCachedCompanyMap() error {

	// ファイルシステムから会社フォルダー一覧を取得
	entries, err := os.ReadDir(srv.managedFolder)
	if err != nil {
		return err
	}

	// キャッシュデータのクリア
	srv.cachedCompanyMap = make(map[string]*models.Company, len(entries))

	// 会社データモデルを作成
	for _, entry := range entries {
		// Companyインスタンスの作成と初期化
		company := models.NewCompany()
		if err := company.ParseFromManagedFolder(srv.managedFolder, entry.Name()); err == nil {
			srv.cachedCompanyMap[company.GetId()] = company
		}
	}

	// 会社の内部情報の取得
	for _, company := range srv.cachedCompanyMap {
		if err := company.LoadPersistData(); err != nil {
			log.Printf("Failed to load persist info for company ID %s: %v", company.GetId(), err)
		}
	}
	return nil
}

func (srv *CompanyService) UpdateCachedCompany(newCompany *models.Company) (*models.Company, error) {

	// Idから更新前の会社情報を取得
	prevCompany, exist := srv.cachedCompanyMap[newCompany.GetId()]
	if !exist {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("指定された会社Idを見つけることができませんでした。"))
	}

	// 会社の短縮名・管理フォルダー・業種カテゴリーが変更された場合は
	// 管理フォルダーの情報を元に既存情報を削除し新たに作成する
	if (newCompany.GetShortName() != prevCompany.GetShortName()) ||
		(newCompany.GetManagedFolder() != prevCompany.GetManagedFolder()) ||
		(newCompany.GetCategoryIndex() != prevCompany.GetCategoryIndex()) {
		managedFolder := newCompany.GetManagedFolder()

		// Companyインスタンスの作成
		if err := newCompany.ParseFromManagedFolder(managedFolder); err != nil {
			return nil, err
		}
		// キャッシュから削除
		delete(srv.cachedCompanyMap, prevCompany.GetId())

		// フォルダー監視の登録解除及び登録
		srv.managedFolderWatcher.Remove(prevCompany.GetManagedFolder())
		srv.managedFolderWatcher.Add(managedFolder)
	}

	// キャッシュに登録
	srv.cachedCompanyMap[newCompany.GetId()] = newCompany

	// persist情報の書き込み
	if err := newCompany.SavePersistData(); err != nil {
		log.Printf("Failed to save persist info for company ID %s: %v", newCompany.GetId(), err)
	}

	return prevCompany, nil
}

// watchManagedFolder は指定された managedFolder の変更を監視します。
// 必要に応じてイベントをサービスに伝播するためのコールバックやチャネルを追加してください。
func (srv *CompanyService) watchManagedFolder(ctx context.Context) error {
	// ポーリング間隔の取得
	pollIntervalMillSec, err := strconv.ParseInt(ext.ConfigMap["CompanyPollIntervalMillSec"], 10, 64)
	if err != nil {
		pollIntervalMillSec = 3000
	}

	// ポーリングベースのウォッチャーを作成（fsnotify が使えない環境向け）
	srv.managedFolderWatcher = filenotify.NewPollingWatcher(time.Duration(pollIntervalMillSec) * time.Millisecond)

	// イベントループ（イベント・エラーを単一ゴルーチンで処理）
	go func() {
		defer srv.managedFolderWatcher.Close()

		for {
			select {
			case event, ok := <-srv.managedFolderWatcher.Events():
				if !ok {
					return
				}

				// 新規ディレクトリを検出したら即座にウォッチに追加（再帰監視の補完）
				if event.Op&fsnotify.Create == fsnotify.Create {
					if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
						if err := srv.addWatchRecursive(event.Name); err != nil {
							log.Printf("[managed-folder] failed to add new dir watch: %v", err)
						}
					}
				}

				// 会社フォルダーを取得
				companyFolder := strings.Replace(event.Name, srv.managedFolder, "", 1)
				companyFolder = strings.Trim(companyFolder, string(os.PathSeparator))
				log.Printf("[managed-folder] detected change: %s %s", event.Op.String(), companyFolder)
				companyFolders := strings.Split(companyFolder, string(os.PathSeparator))
				if len(companyFolders) < 1 {
					return
				}
				managedFolder := filepath.Join(srv.managedFolder, companyFolders[0])

				// 会社のキャッシュを更新
				newCompany := models.NewCompany()
				newCompany.SetManagedFolder(managedFolder)
				_, err := srv.UpdateCachedCompany(newCompany)
				if err != nil {
					log.Println("Failed to update. ManagedFolder =", managedFolder)
					return
				}
				log.Println("Update Company.ShortName=", newCompany.GetShortName())

			case err := <-srv.managedFolderWatcher.Errors():
				log.Printf("[managed-folder] error: %v", err)

			case <-ctx.Done():
				log.Printf("[managed-folder] stop watching (context canceled)")
				return
			}
		}
	}()

	// 監視対象を登録（全ディレクトリを再帰的に追加）
	if err := srv.addWatchRecursive(srv.managedFolder); err != nil {
		return err
	}

	log.Printf("watching managed folder recursively (polling): %s", srv.managedFolder)
	return nil
}

// addWatchRecursive は指定ディレクトリ以下のすべてのディレクトリをウォッチャーに登録します。
// filenotify の Add は重複時にエラーを返すため、その場合は無視します。
func (srv *CompanyService) addWatchRecursive(root string) error {
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		if err := srv.managedFolderWatcher.Add(path); err != nil {
			// "watch exists" を許容
			if !strings.Contains(err.Error(), "exists") {
				return err
			}
		}
		return nil
	})
}

// GetCompanies は管理されている会社情報の一覧を取得します
// gRPCサービスの実装です
func (srv *CompanyService) GetCompanyMap(
	ctx context.Context,
	_ *grpcv1.GetCompanyMapRequest) (
	res *grpcv1.GetCompanyMapResponse,
	err error) {

	// レスポンスを初期化
	res = grpcv1.GetCompanyMapResponse_builder{}.Build()

	// 会社データモデルを作成
	grpcv1CompanyMap := make(map[string]*grpcv1.Company, len(srv.cachedCompanyMap))
	for _, v := range srv.cachedCompanyMap {
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
	company, exist := srv.cachedCompanyMap[id]
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

	// models.Companyインスタンスの作成
	newCompany := &models.Company{Company: req.GetNewCompany()}

	prevCompany, err := srv.UpdateCachedCompany(newCompany)
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
