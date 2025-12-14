package services

import (
	"backend-grpc/internal/core"
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

	// target はこのサービスが管理する工事データのルートフォルダー
	target string

	// targetWatcher は target のファイルシステム監視オブジェクト
	targetWatcher *fsnotify.Watcher

	// kojies は管理されている工事データのインデックスがIdのキャッシュマップ
	kojies map[string]*models.Koji
}

func (s *KojiService) Start(services *Services, options *map[string]string) error {
	// オプションの取得
	optTarget, exists := (*options)["KojiServiceTarget"]
	if !exists {
		return errors.New("KojiServiceTarget option is required")
	}
	// パスを正規化
	target, err := core.NormalizeAbsPath(optTarget)
	if err != nil {
		return err
	}

	// 情報の初期化
	s.services = services
	s.target = target
	s.kojies = make(map[string]*models.Koji, 1000)

	// kojiesByIdの情報を取得
	if err = s.UpdateKojies(); err != nil {
		return err
	}

	// targetの監視を開始
	if err = s.watchTarget(); err != nil {
		return err
	}

	return nil
}

func (s *KojiService) Cleanup() {
	// 現在はクリーンアップ処理は不要
}

// watchTarget starts watching the provided target for changes.
// Add callbacks or channels as needed to propagate events to your services.
func (s *KojiService) watchTarget() error {
	absPath, err := filepath.Abs(s.target)
	if err != nil {
		return err
	}

	s.targetWatcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	// 監視終了時に閉じる
	go func() {
		<-s.targetWatcher.Errors
		s.targetWatcher.Close()
	}()

	// イベントループ
	go func() {
		for {
			select {
			case event, ok := <-s.targetWatcher.Events:
				if !ok {
					return
				}
				log.Printf("[target] event=%s path=%s", event.Op, event.Name)
				if event.Op&(fsnotify.Create|fsnotify.Remove|fsnotify.Rename|fsnotify.Write) != 0 {
					// 必要に応じてサービスへ通知する
					// 例: reload metadata, update cache, etc.
				}

			case err := <-s.targetWatcher.Errors:
				log.Printf("[target] watcher error: %v", err)
			}
		}
	}()

	// フォルダを監視対象に追加
	if err := s.targetWatcher.Add(absPath); err != nil {
		return err
	}

	log.Printf("watching target: %s", absPath)
	return nil
}

func (s *KojiService) UpdateKojies() error {
	// ファイルシステムから工事フォルダー一覧を取得
	entries, err := os.ReadDir(s.target)
	if err != nil {
		return err
	}

	// 工事フォルダー一覧の要素数を取得
	kojiesSize := len(entries)

	// 並列処理用のワーカー数を決定
	numWorkers := core.DecideNumWorkers(kojiesSize)

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
				kojiPath := path.Join(s.target, entries[idx].Name())
				koji := models.NewKoji()
				if err := koji.ParseKojiTarget(kojiPath); err == nil {
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
			s.kojies[result.GetId()] = result
		}
	}

	return nil
}

// GetKojies は管理されている工事データ一覧を返す
func (s *KojiService) GetKojies(
	ctx context.Context,
	req *grpcv1.GetKojiesRequest) (
	res *grpcv1.GetKojiesResponse,
	err error) {
	_ = req // 現状フィルター未対応

	grpcKojies := make(map[string]*grpcv1.Koji, len(s.kojies))
	for _, v := range s.kojies {
		grpcKojies[v.GetId()] = v.Koji
	}

	res.SetKojies(grpcKojies)

	return
}

// GetKojiById は指定されたIDの工事データを返す
func (s *KojiService) GetKojiById(
	ctx context.Context,
	req *grpcv1.GetKojiRequest) (
	res *grpcv1.GetKojiResponse,
	err error) {

	// リクエスト情報の取得
	id := req.GetId()

	// 工事情報を取得
	koji, exist := s.kojies[id]
	if !exist {
		err = connect.NewError(connect.CodeNotFound, errors.New("koji not found"))
		return
	}

	// Responseの更新
	res.SetKoji(koji.Koji)

	return
}

func (s *KojiService) UpdateKoji(
	_ context.Context, req *grpcv1.UpdateKojiRequest) (
	*grpcv1.UpdateKojiResponse, error) {

	// 既存の工事情報を取得
	grpcNewKoji := req.GetNewKoji()
	prevKoji, exist := s.kojies[grpcNewKoji.GetId()]
	if !exist {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("koji not found"))
	}

	newKoji := &models.Koji{Koji: grpcNewKoji}

	// 工事情報を更新
	newKoji, err := prevKoji.Update(newKoji)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 工事情報のインデックスを更新
	if _, exist := s.kojies[prevKoji.GetId()]; exist {
		delete(s.kojies, prevKoji.GetId())
		// 新しいIDで再登録
		s.kojies[newKoji.GetId()] = newKoji
	}

	// Responseの作成
	res := grpcv1.UpdateKojiResponse_builder{}.Build()

	grpcv1KojiMapById := make(map[string]*grpcv1.Koji, len(s.kojies))
	for _, v := range s.kojies {
		grpcv1KojiMapById[v.GetId()] = v.Koji
	}
	res.SetPrevKoji(prevKoji.Koji)

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
