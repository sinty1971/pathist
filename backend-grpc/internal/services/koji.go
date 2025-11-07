package services

import (
	"backend-grpc/internal/utils"
	"context"
	"errors"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"

	grpcv1 "backend-grpc/gen/grpc/v1"
	grpcv1connect "backend-grpc/gen/grpc/v1/grpcv1connect"
	"backend-grpc/internal/models"

	"connectrpc.com/connect"
	"github.com/fsnotify/fsnotify"
)

// KojiService bridges existing KojiService logic to Connect handlers.
type KojiService struct {
	// Embed the unimplemented handler for forward compatibility
	grpcv1connect.UnimplementedKojiServiceHandler

	// services は任意のgrpcサービスハンドラーへの参照
	services *Services

	// managedFolder はこのサービスが管理する工事データのルートフォルダー
	managedFolder string

	// managedFolderWatcher は managedFolder のファイルシステム監視オブジェクト
	managedFolderWatcher *fsnotify.Watcher

	// kojiMapById は管理されている工事データのインデックスがIdのキャッシュマップ
	kojiMapById map[string]*models.Koji
}

func NewKojiService(
	services *Services, options *ServiceOptions) (
	s *KojiService, err error) {
	// インスタンス作成
	s = &KojiService{
		services:      services,
		managedFolder: options.KojiServiceManagedFolder,
		kojiMapById:   make(map[string]*models.Koji, 1000),
	}

	// kojiesByIdの情報を取得
	if err = s.UpdateKojies(); err != nil {
		return
	}

	// managedFolderの監視を開始
	if err = s.watchManagedFolder(); err != nil {
		return
	}

	return
}

func (s *KojiService) Cleanup() {
	// 現在はクリーンアップ処理は不要
}

// watchManagedFolder starts watching the provided managedFolder for changes.
// Add callbacks or channels as needed to propagate events to your services.
func (s *KojiService) watchManagedFolder() error {
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

func (s *KojiService) UpdateKojies() error {
	// ファイルシステムから工事フォルダー一覧を取得
	entries, err := os.ReadDir(s.managedFolder)
	if err != nil {
		return err
	}

	// 工事フォルダー一覧の要素数を取得
	kojiesSize := len(entries)

	// 並列処理用のワーカー数を決定
	numWorkers := utils.DecideNumWorkers(kojiesSize,
		utils.WithMinWorkers(2),
		utils.WithMaxWorkers(16),
		utils.WithCPUMultiplier(2),
	)

	// バッファ付きチャンネルで効率化
	jobs := make(chan int, kojiesSize)
	results := make(chan *models.Koji, kojiesSize)

	// ワーカープールの起動
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			for idx := range jobs {
				kojiPath := path.Join(s.managedFolder, entries[idx].Name())
				if koji, err := models.NewKoji(kojiPath); err == nil {
					results <- koji
				} else {
					results <- nil // エラーの場合はnilを返す
				}
			}
		}()
	}

	// ジョブの投入
	go func() {
		for i := range entries {
			jobs <- i
		}
		close(jobs)
	}()

	// 結果収集用のゴルーチン
	go func() {
		wg.Wait()
		close(results)
	}()

	// 結果を収集（最大サイズで確保し、後でスライス）
	for result := range results {
		if result != nil {
			s.kojiMapById[result.GetId()] = result
		}
	}

	return nil
}

// GetKojies は管理されている工事データ一覧を返す
func (s *KojiService) GetKojiMapById(
	ctx context.Context, req *connect.Request[grpcv1.GetKojiMapByIdRequest]) (
	res *connect.Response[grpcv1.GetKojiMapByIdResponse], err error) {
	_ = req // 現状フィルター未対応

	grpcKojisById := make(map[string]*grpcv1.Koji, len(s.kojiMapById))
	for _, v := range s.kojiMapById {
		grpcKojisById[v.GetId()] = v.Koji
	}

	return res, nil
}

func (s *KojiService) GetKojiById(
	ctx context.Context,
	req *connect.Request[grpcv1.GetKojiByIdRequest]) (
	res *connect.Response[grpcv1.GetKojiByIdResponse],
	err error) {

	// リクエスト情報の取得
	id := req.Msg.GetId()

	// 工事情報を取得
	koji, exist := s.kojiMapById[id]
	if !exist {
		err = connect.NewError(connect.CodeNotFound, errors.New("koji not found"))
		return
	}

	// Responseの更新
	res.Msg.SetKoji(koji.Koji)

	return
}

func (s *KojiService) UpdateKoji(
	ctx context.Context,
	req *connect.Request[grpcv1.UpdateKojiRequest]) (
	res *connect.Response[grpcv1.UpdateKojiResponse],
	err error) {

	// 既存の工事情報を取得
	currentKojiId := req.Msg.GetCurrentKojiId()
	currentKoji, exist := s.kojiMapById[currentKojiId]
	if !exist {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("koji not found"))
	}
	// 更新後の工事情報を取得
	grpcUpdatedKoji := req.Msg.GetUpdatedKoji()
	if grpcUpdatedKoji == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("updated koji is nil"))
	}
	updatedKoji := &models.Koji{
		Koji: grpcUpdatedKoji,
	}

	// 工事情報を更新
	updatedKoji, err = currentKoji.Update(updatedKoji)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 工事情報のインデックスを更新
	if _, exist := s.kojiMapById[currentKojiId]; exist {
		delete(s.kojiMapById, currentKojiId)
		// 新しいIDで再登録
		s.kojiMapById[updatedKoji.GetId()] = updatedKoji
	}

	// Responseの作成
	grpcv1KojiMapById := make(map[string]*grpcv1.Koji, len(s.kojiMapById))
	for _, v := range s.kojiMapById {
		grpcv1KojiMapById[v.GetId()] = v.Koji
	}
	res.Msg.SetKojiMapById(grpcv1KojiMapById)

	return res, nil
}

// RenameStandardFile は標準ファイルの名前を変更し、工事データも更新する
// TODO: StandardFile型が定義されていないため、一時的にコメントアウト
// func (ks *KojiService) RenameStandardFile(koji models.Koji, actuals []string) []string {
// 	// マップの作成
// 	actualToStandardMap := make(map[string]*models.StandardFile)
// 	for i := range koji.StandardFiles {
// 		sf := &koji.StandardFiles[i]
// 		actualToStandardMap[sf.ActualName] = sf
// 	}
//
// 	// 変更後の標準ファイル名を格納する配列
// 	renamedFiles := make([]string, len(actuals))
//
// 	// 変更前の標準ファイル名をループ
// 	count := 0
// 	for _, actual := range actuals {
// 		if sf, exists := actualToStandardMap[actual]; exists {
// 			actualFullpath, err := ks.BaseFolderService.GetFullpath(sf.GetPath())
// 			if err != nil {
// 				continue
// 			}
//
// 			standardFullpath, err := ks.BaseFolderService.GetFullpath(sf.Name)
// 			if err != nil {
// 				continue
// 			}
// 			count++
// 		}
//
// 		// ファイル名変更後、工事の必須ファイル情報を更新
// 		if count > 0 {
// 			// 必須ファイル情報を再設定
// 			err := ks.UpdateRequiredFiles(&koji)
// 			if err == nil {
// 				// 属性ファイルに反映
// 				ks.DatabaseService.Save(&koji)
// 			}
// 		}
// 	}
//
// 	// ファイル名変更後、工事の必須ファイル情報を更新
// 	if count > 0 {
// 		// 必須ファイル情報を再設定
// 		err := ks.UpdateRequiredFiles(&koji)
// 		if err == nil {
// 			// 属性ファイルに反映
// 			ks.DatabaseService.Save(&koji)
// 		}
// 	}
//
// 	return renamedFiles[:count]
// }
