package services

import (
	"fmt"
	"path/filepath"
	"penguin-backend/internal/models"
	"penguin-backend/internal/utils"
	"regexp"
	"sort"
	"strings"
	"sync"
)

// KojiService は工事情報取得サービスを提供する
type KojiService struct {
	FileService     *FileService
	DatabaseService *DatabaseService[*models.Koji]
}

// NewKojiService はKojiServiceを初期化する
// businessFileService: ビジネスファイルサービス
// folderName: 工事一覧フォルダーのファイル名(FileService.BasePathからの相対パス)
func NewKojiService(businessFileService *FileService, folderName string) (*KojiService, error) {

	// 工事一覧フォルダーのフルパスの取得
	folderPath, err := businessFileService.GetFullpath(folderName)
	if err != nil {
		return nil, err
	}

	// FileServiceを初期化
	fileService, err := NewFileService(folderPath)
	if err != nil {
		return nil, err
	}

	// DatabaseServiceを初期化
	databaseService := NewDatabaseService[*models.Koji](fileService, ".detail.yaml")

	// KojiServiceのインスタンスを作成
	return &KojiService{
		FileService:     fileService,
		DatabaseService: databaseService,
	}, nil
}

// GetKoji は指定されたパスから工事を取得する
// pathは工事フォルダーのファイル名
// 工事を返す
func (ks *KojiService) GetKoji(folderName string) (models.Koji, error) {
	// 工事データモデルを作成
	koji, err := models.NewKoji(folderName)
	if err != nil {
		return models.Koji{}, fmt.Errorf("工事データモデルの作成に失敗しました: %v", err)
	}

	// 補助ファイルの設定
	err = ks.UpdateAssistFiles(&koji)
	if err != nil {
		return models.Koji{}, fmt.Errorf("補助ファイルの設定に失敗しました: %v", err)
	}

	// 属性ファイルと同期する
	err = ks.ImportDatabaseFile(&koji)
	if err != nil {
		return models.Koji{}, fmt.Errorf("属性ファイルとの同期に失敗しました: %v", err)
	}

	return koji, nil
}

// KojiResult 並列処理用の結果構造体
type KojiResult struct {
	Koji    models.Koji
	Success bool
}

// GetRecentKojies は指定されたパスから最近の工事データモデル一覧を取得する
func (ks *KojiService) GetRecentKojies() []models.Koji {
	// 工事フォルダー一覧をデータベースファイルと同期せずに取得
	kojies := ks.GetRecentKojiesNoSyncDatabaseFile()

	if len(kojies) == 0 {
		return kojies
	}

	// 並列処理用のワーカー数を決定
	numWorkers := utils.DecideNumWorkers(len(kojies),
		utils.WithMinWorkers(2),
		utils.WithMaxWorkers(8), // SyncDatabaseFileはI/O処理なので少なめ
		utils.WithCPUMultiplier(1),
	)

	// ワーカープールでSyncDatabaseFileを並列実行
	var wg sync.WaitGroup
	jobs := make(chan int, len(kojies))

	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			for idx := range jobs {
				// データベースファイルと同期（エラーは無視）
				_ = ks.ImportDatabaseFile(&kojies[idx])
			}
		}()
	}

	// ジョブを投入
	for i := range kojies {
		jobs <- i
	}
	close(jobs)

	// 全ワーカーの完了を待つ
	wg.Wait()

	return kojies
}

// GetRecentKojiesNoSyncDatabaseFile は最近の工事データモデル一覧を取得する
// データベースへのアクセスは行わない
// goルーチンを使用して並列処理を行っています
func (ks *KojiService) GetRecentKojiesNoSyncDatabaseFile() []models.Koji {
	// ファイルシステムから工事フォルダー一覧を取得
	fileInfos, err := ks.FileService.GetFileInfos()
	if err != nil || len(fileInfos) == 0 {
		return []models.Koji{}
	}

	// 並列処理用のワーカー数を決定
	numWorkers := utils.DecideNumWorkers(len(fileInfos),
		utils.WithMinWorkers(2),
		utils.WithMaxWorkers(16),
		utils.WithCPUMultiplier(2),
	)

	// バッファ付きチャンネルで効率化
	jobs := make(chan int, len(fileInfos))
	results := make(chan KojiResult, len(fileInfos))

	// ワーカープールの起動
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			for idx := range jobs {
				if koji, err := models.NewKoji(fileInfos[idx].Name); err == nil {
					results <- KojiResult{Koji: koji, Success: true}
				} else {
					results <- KojiResult{Koji: models.Koji{}, Success: false}
				}
			}
		}()
	}

	// ジョブの投入
	go func() {
		for i := range fileInfos {
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
	kojies := make([]models.Koji, len(fileInfos))
	count := 0
	for result := range results {
		if result.Success {
			kojies[count] = result.Koji
			count++
		}
	}
	// 実際に成功した要素数にスライス
	kojies = kojies[:count]

	// 工事開始日で降順ソート（新しいものが最初）
	sort.Slice(kojies, func(i, j int) bool {
		// Compare関数を使用して比較（降順なので結果を反転）
		cmp := kojies[i].StartDate.Compare(kojies[j].StartDate)
		if cmp != 0 {
			return cmp > 0
		}
		// 両方の開始日が同じ場合、フォルダー名で降順
		return kojies[i].FolderName > kojies[j].FolderName
	})

	return kojies
}

// ImportDatabaseFile は工事データベースファイルからデータを取り込みます
func (ks *KojiService) ImportDatabaseFile(target *models.Koji) error {
	// 工事データベースファイルからデータを読み込む
	database, err := ks.DatabaseService.Load(target)
	if err != nil {
		// ファイルが存在しない場合はデフォルト値を設定
		target.UpdateStatus()

		// 新規ファイルの場合は非同期で保存（必要最小限の書き込み）
		go func() {
			ks.DatabaseService.Save(target)
		}()
		return nil
	}

	// データベースファイルにしか保持されない内容をtargetに反映
	target.Description = database.Description
	target.EndDate = database.EndDate
	target.Tags = database.Tags

	// 計算が必要な項目の更新
	target.UpdateStatus()

	return nil
}

// UpdateAssistFiles は工事の補助ファイルを更新します
func (ks *KojiService) UpdateAssistFiles(koji *models.Koji) error {
	// 工事フォルダー内のファイルを取得
	kojiFileInfos, err := ks.FileService.GetFileInfos(koji.FolderName)
	if err != nil {
		// エラーの場合はデフォルトの補助ファイルを設定
		koji.Assists = []models.AssistFile{{
			Desired: koji.FolderName + ".xlsx",
			Current: "工事.xlsx",
		}}
		return nil
	}

	// 管理ファイルの設定（最大1個のみ）
	var assistFile models.AssistFile
	var found bool

	// 日付パターンのコンパイル（1回のみ）
	datePattern := regexp.MustCompile(`^20\d{2}-\d{4}`)

	// 工事ファイルの検索
	for _, file := range kojiFileInfos {
		// エクセルファイル以外はスキップ
		if !file.IsExcel() {
			continue
		}

		// ファイル名の解析
		extension := filepath.Ext(file.Name)
		filenameWithoutExt := strings.TrimSuffix(file.Name, extension)

		// 条件に合うファイルかチェック
		if !datePattern.MatchString(filenameWithoutExt) && filenameWithoutExt != "工事" {
			continue
		}

		recommendedName := koji.FolderName + extension

		assistFile = models.AssistFile{
			Desired: recommendedName,
			Current: file.Name,
		}
		found = true
		break
	}

	// 工事ファイルが見つからない場合はひな形を設定
	if !found {
		assistFile = models.AssistFile{
			Desired: koji.FolderName + ".xlsx",
			Current: "工事.xlsx",
		}

		// テンプレートファイルのコピーは省略（高速化のため）
		// 必要に応じて別途実行する
	}

	// 管理ファイルの設定
	koji.Assists = []models.AssistFile{assistFile}

	return nil
}

// Update は工事情報 FileInfoの情報を更新
// kojiの詳細情報で.detail.yamlを更新
// 移動元フォルダー名：FileInfo.Name
// 移動先フォルダー名：StartDate, CompanyName, LocationNameから生成されたフォルダー名
func (ks *KojiService) Update(target *models.Koji) error {
	// 更新前の工事情報を保存
	temp := *target

	// フォルダー名の変更
	if temp.UpdateFolderName() {
		// 工事IDの更新（フォルダー名の変更に伴う）
		temp.UpdateID()

		// 補助ファイル情報の更新（フォルダー名の変更に伴う）
		if err := ks.UpdateAssistFiles(&temp); err != nil {
			return err
		}

		// ファイルの移動
		if err := ks.FileService.MoveFile(temp.FolderName, target.FolderName); err != nil {
			return err
		}
	}

	// フォルダー名の変更完了後、情報を更新
	*target = temp

	// 計算が必要な項目の更新
	target.UpdateStatus()

	// 更新後の工事情報を属性ファイルに反映
	return ks.DatabaseService.Save(target)
}

// RenameAssistFile は補助ファイルの名前を変更し、工事データも更新する
// kojiは工事データ
// currentsは変更前の補助ファイル名
// 変更後の補助ファイル名を返す
func (ks *KojiService) RenameAssistFile(koji models.Koji, currents []string) []string {
	renamedFiles := make([]string, len(currents))

	count := 0
	for _, current := range currents {
		for _, assistFile := range koji.Assists {
			if assistFile.Current == current {
				currentPath := filepath.Join(koji.FolderName, current)
				recommendedPath := filepath.Join(koji.FolderName, assistFile.Desired)
				err := ks.FileService.MoveFile(currentPath, recommendedPath)
				if err == nil {
					renamedFiles[count] = recommendedPath
					count++
				}
			}
		}
	}

	// ファイル名変更後、工事の補助ファイル情報を更新
	if count > 0 {
		// 補助ファイル情報を再設定
		err := ks.UpdateAssistFiles(&koji)
		if err == nil {
			// 属性ファイルに反映
			ks.DatabaseService.Save(&koji)
		}
	}

	return renamedFiles[:count]
}
