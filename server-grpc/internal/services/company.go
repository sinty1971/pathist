package services

import (
	"context"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	grpcv1 "server-grpc/gen/grpc/v1"
	grpcv1connect "server-grpc/gen/grpc/v1/grpcv1connect"
	"server-grpc/internal/core"
	"server-grpc/internal/models"

	"connectrpc.com/connect"
	"github.com/fsnotify/fsnotify"
)

// CompanyService の実装
type CompanyService struct {
	// services は任意のgrpcサービスハンドラーへの参照
	services *Services

	// Embed the unimplemented handler for forward compatibility
	grpcv1connect.UnimplementedCompanyServiceHandler

	// companies は会社データのキャッシュマップ
	companies map[string]*models.Company

	// target はこのサービスが管理する会社データのルートフォルダー
	target string

	// targetWatcher は target のファイルシステム監視オブジェクト
	targetWatcher *fsnotify.Watcher

	// watchedDirs は監視登録済みディレクトリの集合
	watchedDirs map[string]struct{}
}

const companyWatcherMaxDepth = 2

// Start は CompanyService を初期化して開始します
func (srv *CompanyService) Start(services *Services, options *map[string]string) error {
	// パスをの取得と正規化
	optTarget, exists := (*options)["CompanyServiceTarget"]
	if !exists {
		return errors.New("CompanyServiceTarget option is required")
	}
	target, err := core.NormalizeAbsPath(optTarget)
	if err != nil {
		return err
	}

	// 既存インスタンスに値をセット（再代入しないこと）
	srv.services = services
	srv.target = target
	srv.companies = map[string]*models.Company{}
	srv.watchedDirs = make(map[string]struct{})

	// companiesの情報を取得
	if err = srv.UpdateCompanies(); err != nil {
		return err
	}

	// watcherの開始
	if err = srv.watchTarget(); err != nil {
		return err
	}

	return nil
}

func (srv *CompanyService) Cleanup() {
	// 現状クリーンアップは不要（監視を廃止したため）
}

// watchTarget は target のファイルシステム監視を開始します
// 階層は2階層まで監視します
func (srv *CompanyService) watchTarget() error {
	// fsnotifyウォッチャーの作成
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	srv.targetWatcher = watcher

	// target配下2階層までを監視対象に追加
	if err := srv.addWatchersRecursively(srv.target, 0); err != nil {
		return err
	}

	// ゴルーチンで監視イベントを処理
	go srv.consumeWatcherEvents()

	return nil
}

// consumeWatcherEvents はファイルシステム監視イベントを処理します
func (srv *CompanyService) consumeWatcherEvents() {
	for {
		select {
		case event, ok := <-srv.targetWatcher.Events:
			if !ok {
				return
			}
			log.Printf("CompanyService: File system event: %s", event)
			srv.handleWatcherEvent(event)

			// 会社キャッシュの更新
			if err := srv.UpdateCompanies(); err != nil {
				log.Printf("CompanyService: Failed to update company cache map: %v", err)
			}

		case err, ok := <-srv.targetWatcher.Errors:
			if !ok {
				return
			}
			log.Printf("CompanyService: File system watcher error: %v", err)
		}
	}
}

// handleWatcherEvent はファイルシステム監視イベントを処理します
func (srv *CompanyService) handleWatcherEvent(event fsnotify.Event) {
	// 無効なパスの場合は無視
	if event.Name == "" {
		return
	}

	// ディレクトリ作成イベントの場合は監視対象を拡張
	if event.Op&fsnotify.Create != 0 {
		if err := srv.addWatchIfDirectory(event.Name); err != nil {
			log.Printf("CompanyService: Failed to expand watcher for %s: %v", event.Name, err)
		}
	}

	// ディレクトリ削除・移動イベントの場合は監視対象を解除
	if event.Op&(fsnotify.Remove|fsnotify.Rename) != 0 {
		srv.unregisterWatcherTree(event.Name)
	}
}

// addWatchIfDirectory は指定パスがディレクトリの場合に監視対象に追加します
func (srv *CompanyService) addWatchIfDirectory(path string) error {
	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		return nil
	}

	depth, ok := srv.relativeDepth(path)
	if !ok || depth > companyWatcherMaxDepth {
		return nil
	}

	return srv.addWatchersRecursively(path, depth)
}

func (srv *CompanyService) addWatchersRecursively(dir string, depth int) error {
	if depth > companyWatcherMaxDepth {
		return nil
	}

	cleanDir := filepath.Clean(dir)
	if err := srv.registerWatcher(cleanDir); err != nil {
		return err
	}

	if depth == companyWatcherMaxDepth {
		return nil
	}

	entries, err := os.ReadDir(cleanDir)
	if err != nil {
		log.Printf("CompanyService: Failed to read directory %s: %v", cleanDir, err)
		return nil
	}

	for _, entry := range entries {
		if entry.IsDir() {
			if err := srv.addWatchersRecursively(filepath.Join(cleanDir, entry.Name()), depth+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func (srv *CompanyService) registerWatcher(target string) error {
	if _, watched := srv.watchedDirs[target]; watched {
		return nil
	}

	if err := srv.targetWatcher.Add(target); err != nil {
		return err
	}
	srv.watchedDirs[target] = struct{}{}

	return nil
}

func (srv *CompanyService) unregisterWatcherTree(target string) {
	cleanPath := filepath.Clean(target)
	if _, ok := srv.relativeDepth(cleanPath); !ok {
		return
	}

	targets := make([]string, 0, len(srv.watchedDirs))
	prefix := cleanPath + string(os.PathSeparator)
	for dir := range srv.watchedDirs {
		if dir == cleanPath || strings.HasPrefix(dir, prefix) {
			targets = append(targets, dir)
		}
	}

	for _, dir := range targets {
		if err := srv.targetWatcher.Remove(dir); err != nil {
			log.Printf("CompanyService: Failed to remove watcher for %s: %v", dir, err)
		}
		delete(srv.watchedDirs, dir)
	}
}

// relativeDepth は指定パスの target からの相対深度を返します
func (srv *CompanyService) relativeDepth(target string) (int, bool) {
	rel, err := filepath.Rel(srv.target, target)
	if err != nil {
		return -1, false
	}

	slashed := filepath.ToSlash(rel)
	if slashed == "." || slashed == "" {
		return 0, true
	}

	if slashed == ".." || strings.HasPrefix(slashed, "../") {
		return -1, false
	}

	return strings.Count(slashed, "/") + 1, true
}

// UpdateCompanies 会社のキャッシュデータを更新します
func (srv *CompanyService) UpdateCompanies() error {
	// ファイルシステムから会社フォルダー一覧を取得
	entries, err := os.ReadDir(srv.target)
	if err != nil {
		return err
	}

	// キャッシュデータの初期化
	srv.companies = make(map[string]*models.Company, len(entries))

	// 全てのCompanyインスタンスを作成
	for _, entry := range entries {
		// Companyインスタンスの作成と初期化
		company := models.NewCompany()
		if err := company.ParseFrom(srv.target, entry.Name()); err == nil {
			srv.companies[company.GetId()] = company
		}
	}

	// 会社の内部情報の取得
	for _, company := range srv.companies {
		// persist情報の読み込み
		if err := company.Pathist.LoadPersists(); err != nil {
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
	prevCompany, exist := srv.companies[prevId]

	// キャッシュから削除
	if exist {
		delete(srv.companies, prevId)
	}

	// 新しい会社情報の管理フォルダー名を生成
	newTarget := models.GenerateCompanyPathistFolder(
		filepath.Dir(newCompany.GetPathistFolder()),
		newCompany.GetCategoryIndex(),
		newCompany.GetShortName())

	// 管理フォルダーの変更がある場合はフォルダー移動を実施
	if exist && prevCompany.GetPathistFolder() != newCompany.GetPathistFolder() {
		// 管理フォルダーの変更がある場合はフォルダー移動を実施
		prevTarget := prevCompany.GetPathistFolder()
		if err := os.Rename(prevTarget, newTarget); err != nil {
			return nil, err
		}
	}

	// 既存の情報を更新
	srv.companies[newCompany.GetId()] = newCompany

	// persist情報の書き込み
	if err := newCompany.Pathist.SavePersists(); err != nil {
		log.Printf("Failed to save persist info for company ShortName %s: %v", newCompany.GetShortName(), err)
	}

	return prevCompany, nil
}

// GetCompanies は管理されている会社情報の一覧を取得します
// gRPCサービスの実装です
func (srv *CompanyService) GetCompanies(
	ctx context.Context, req *grpcv1.GetCompaniesRequest) (
	*grpcv1.GetCompaniesResponse, error) {

	// レスポンスを初期化
	res := grpcv1.GetCompaniesResponse_builder{}.Build()

	// 必要に応じてキャッシュを更新
	if req.GetRefresh() {
		if err := srv.UpdateCompanies(); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	// 会社データモデルを作成
	grpcv1Companies := make(map[string]*grpcv1.Company, len(srv.companies))
	for _, v := range srv.companies {
		grpcv1Companies[v.Company.GetId()] = v.Company
	}

	// Responseの更新とリターン
	res.SetCompanies(grpcv1Companies)
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
	company, exist := srv.companies[id]
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
