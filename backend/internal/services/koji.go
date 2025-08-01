package services

import (
	"errors"
	"fmt"
	"os"
	"path"
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
	// RootService はトップコンテナのインスタンス
	RootService *ContainerService

	// 工事一覧フォルダーのフルパス名
	FolderPath string

	// データベースサービス
	DatabaseService *DatabaseService[*models.Koji]
}

// BuildWithContext は opt でKojiServiceを初期化する
// businessFileService: ビジネスファイルサービス
// folderName: 工事一覧フォルダーのファイル名 (Ex. "2 工事")
func (ks *KojiService) BuildWithContext(opt ContainerOption, folderPath string, databaseFilename string) error {

	// folderPath がアクセス可能かチェック
	fi, err := os.Stat(folderPath)
	if err != nil {
		return err
	}

	// フォルダーで無ければエラー
	if !fi.IsDir() {
		return fmt.Errorf("工事一覧対象パスがフォルダーではありません: %s", folderPath)
	}

	// ルートサービスを設定
	ks.RootService = opt.RootService

	// 工事一覧フォルダーのフルパスの取得
	ks.FolderPath = folderPath

	// 工事一覧用のDatabaseServiceを作成
	ks.DatabaseService = NewDatabaseService[*models.Koji](databaseFilename)

	return nil
}

// GetKoji は指定されたパスから工事を取得する
// folder は工事フォルダー名 (Ex. '2025-0618 豊田築炉 名和工場')
// 工事を返す
func (ks *KojiService) GetKoji(folderName string) (*models.Koji, error) {
	// フォルダー名からフルパスを作成
	folderPath := filepath.Join(ks.FolderPath, folderName)

	// 工事データモデルを作成
	koji, err := models.NewKoji(folderPath)
	if err != nil {
		return nil, fmt.Errorf("工事データモデルの作成に失敗しました: %v", err)
	}

	// データベースファイルから読み込む
	err = ks.LoadDatabaseFile(koji)
	if err != nil {
		return nil, fmt.Errorf("データベースファイルからの読み込みに失敗しました: %v", err)
	}

	// 標準ファイルの設定
	err = ks.UpdateStandardFiles(koji)
	if err != nil {
		return nil, fmt.Errorf("標準ファイルの設定に失敗しました: %v", err)
	}

	return koji, nil
}

type getKojiesMode int

const (
	GetKoujiesModeRecent getKojiesMode = 1 << iota
	GetKoujiesModeOld
	GetKojiesModeSyncDatabase
)

func getKojiesModeFunc(modes ...getKojiesMode) getKojiesMode {
	opt := getKojiesMode(GetKoujiesModeRecent)
	for _, mode := range modes {
		opt |= mode
	}
	return opt
}

// GetKojies は指定されたパスから工事データリストを取得する
// modes: 取得モード
// 工事データリストを返す
func (ks *KojiService) GetKojies(modes ...getKojiesMode) []models.Koji {
	// モードを決定
	mode := getKojiesModeFunc(modes...)

	// 未実装のモードはエラーを返す
	if mode&GetKoujiesModeOld != 0 {
		return []models.Koji{}
	}

	// ファイルシステムから工事フォルダー一覧を取得
	entries, err := os.ReadDir(ks.FolderPath)
	if err != nil || len(entries) == 0 {
		return []models.Koji{}
	}

	// 工事フォルダー一覧の要素数を取得
	entryLength := len(entries)

	// 並列処理用のワーカー数を決定
	numWorkers := utils.DecideNumWorkers(entryLength,
		utils.WithMinWorkers(2),
		utils.WithMaxWorkers(16),
		utils.WithCPUMultiplier(2),
	)

	// バッファ付きチャンネルで効率化
	jobs := make(chan int, entryLength)
	results := make(chan *models.Koji, entryLength)

	// ワーカープールの起動
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			for idx := range jobs {
				kojiPath := path.Join(ks.FolderPath, entries[idx].Name())
				if koji, err := models.NewKoji(kojiPath); err == nil {
					if mode&GetKojiesModeSyncDatabase != 0 {
						_ = ks.LoadDatabaseFile(koji)
					}
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
	kojies := make([]models.Koji, entryLength)
	count := 0
	for result := range results {
		if result != nil {
			kojies[count] = *result
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
		return kojies[i].FolderPath > kojies[j].FolderPath
	})

	return kojies
}

// LoadDatabaseFile は工事データベースファイルからデータを取り込みます
func (ks *KojiService) LoadDatabaseFile(target *models.Koji) error {
	// 工事データベースファイルからデータを読み込む
	database, err := ks.DatabaseService.Load(target)
	if err != nil {
		// ファイルが存在しない場合はデフォルト値を設定
		target.UpdateStatus()

		// 新規ファイルの場合は非同期で保存
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
// データベースファイルの情報を優先するので、値が空の場合のみ更新する
func (ks *KojiService) UpdateStandardFiles(koji *models.Koji) error {
	// 補助ファイルがすでに設定されている場合は更新しない
	if len(koji.StandardFiles) > 0 {
		return errors.New("補助ファイルがすでに設定されています")
	}

	// 工事フォルダー内のファイルを取得
	entries, err := os.ReadDir(koji.FolderPath)
	if err != nil {
		// エラーの場合はデフォルトの標準ファイルを設定
		standardFile, err := models.NewFileInfo(path.Join(koji.FolderPath, "工事.xlsx"))
		if err != nil {
			return err
		}
		koji.StandardFiles = []models.FileInfo{*standardFile}
		return nil
	}

	// 標準ファイルの設定（現在は1個のみ）
	var standardFile *models.FileInfo
	var found bool

	// 日付パターンのコンパイル（1回のみ）
	datePattern := regexp.MustCompile(`^20\d{2}-\d{4}`)

	// 工事ファイルの検索
	for _, file := range entries {
		filename := file.Name()
		// エクセルファイル以外はスキップ
		if !utils.FilenameIsExcel(filename) {
			continue
		}

		// ファイル名の解析
		extension := filepath.Ext(filename)
		filenameWithoutExt := strings.TrimSuffix(filename, extension)

		// 条件に合うファイルかチェック
		if !datePattern.MatchString(filenameWithoutExt) && filenameWithoutExt != "工事" {
			continue
		}

		actualPath := path.Join(koji.FolderPath, file.Name())
		standardName := koji.GetFolderName() + extension
		standardPath := path.Join(koji.FolderPath, standardName)
		standardFile, err = models.NewFileInfo(actualPath, standardPath)
		if err != nil {
			continue
		}
		found = true
		break
	}

	// 工事ファイルが見つからない場合はひな形を設定
	if !found {
		standardFile, err = models.NewFileInfo("工事.xlsx", standardPath)

		// テンプレートファイルのコピーは省略（高速化のため）
		// 必要に応じて別途実行する
	}

	// 管理ファイルの設定
	koji.StandardFiles = []models.StandardFile{standardFile}

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
		temp.StandardFiles = []models.StandardFile{}
		if err := ks.UpdateStandardFiles(&temp); err != nil {
			return err
		}

		// ファイルの移動
		if err := ks.BaseFolderService.MoveFile(temp.FolderName, target.FolderName); err != nil {
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

// RenameStandardFile は標準ファイルの名前を変更し、工事データも更新する
// kojiは工事データ
// currentsは変更前の標準ファイル名
// 変更後の標準ファイル名を返す
func (ks *KojiService) RenameStandardFile(koji models.Koji, actuals []string) []string {
	// マップの作成
	actualToStandardMap := make(map[string]*models.StandardFile)
	for i := range koji.StandardFiles {
		sf := &koji.StandardFiles[i]
		actualToStandardMap[sf.ActualName] = sf
	}

	// 変更後の標準ファイル名を格納する配列
	renamedFiles := make([]string, len(actuals))

	// 変更前の標準ファイル名をループ
	count := 0
	for _, actual := range actuals {
		if sf, exists := actualToStandardMap[actual]; exists {
			actualFullpath, err := ks.BaseFolderService.GetFullpath(sf.GetPath())
			if err != nil {
				continue
			}

			standardFullpath, err := ks.BaseFolderService.GetFullpath(sf.Name)
			if err != nil {
				continue
			}
			count++
		}

		// ファイル名変更後、工事の補助ファイル情報を更新
		if count > 0 {
			// 補助ファイル情報を再設定
			err := ks.UpdateStandardFiles(&koji)
			if err == nil {
				// 属性ファイルに反映
				ks.DatabaseService.Save(&koji)
			}
		}
	}

	// ファイル名変更後、工事の補助ファイル情報を更新
	if count > 0 {
		// 補助ファイル情報を再設定
		err := ks.UpdateStandardFiles(&koji)
		if err == nil {
			// 属性ファイルに反映
			ks.DatabaseService.Save(&koji)
		}
	}

	return renamedFiles[:count]
}
