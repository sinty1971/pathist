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
// 工事一覧フォルダーを管理します
type KojiService struct {
	// ルートサービス
	Root *RootService

	// リポジトリーサービス
	Repository *RepositoryService[*models.Koji]

	// 工事一覧フォルダーのフルパス
	Folder string
}

func (ks *KojiService) Register(rs *RegistableService, opts ...Option) {
	ks.Root = rs.GetRootService()
	ks.Repository = NewRepositoryService[*models.Koji](repositoryFilename)
	ks.Folder = folder
}

// CreateKojiService は opt でKojiServiceを初期化する
// root: ルートサービス
// repositoryFilename: リポジトリーサービスのファイル名
// folder: 工事一覧フォルダーのフルパス
func (ks *KojiService) CreateKojiService(root *RootService, repositoryFilename, folder string) error {

	// 工事一覧フォルダーのフルパスの取得
	folder, err := utils.CleanAbsPath(folder)
	if err != nil {
		return err
	}

	// folderPath がアクセス可能かチェック
	fi, err := os.Stat(folder)
	if err != nil {
		return err
	}

	// フォルダーで無ければエラー
	if !fi.IsDir() {
		return fmt.Errorf("工事一覧対象パスがフォルダーではありません: %s", folder)
	}

	// ルートサービスに工事サービスを登録
	root.KojiService = ks

	// 工事一覧用のリポジトリーサービスを作成
	ks.Repository = NewRepositoryService[*models.Koji](repositoryFilename)

	return nil
}

// GetKoji は指定されたパスから工事を取得する
// folder は工事フォルダー名 (Ex. '2025-0618 豊田築炉 名和工場')
// 工事を返す
func (ks *KojiService) GetKoji(folderName string) (*models.Koji, error) {
	// フォルダー名からフルパスを作成
	folderPath := filepath.Join(ks.Folder, folderName)

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

	// 必須ファイルの設定
	err = ks.UpdateRequiredFiles(koji)
	if err != nil {
		return nil, fmt.Errorf("必須ファイルの設定に失敗しました: %v", err)
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
	entries, err := os.ReadDir(ks.Folder)
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
				kojiPath := path.Join(ks.Folder, entries[idx].Name())
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
	database, err := ks.Repository.Load(target)
	if err != nil {
		// ファイルが存在しない場合はデフォルト値を設定
		target.UpdateStatus()

		// 新規ファイルの場合は非同期で保存
		go func() {
			ks.Repository.Save(target)
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

// UpdateRequiredFiles は工事の必須ファイルを更新します
// データベースファイルの情報を優先するので、値が空の場合のみ更新する
func (ks *KojiService) UpdateRequiredFiles(koji *models.Koji) error {
	// 必須ファイルがすでに設定されている場合は更新しない
	if len(koji.RequiredFiles) > 0 {
		return errors.New("必須ファイルがすでに設定されています")
	}

	// 工事フォルダー内のファイルを取得
	entries, err := os.ReadDir(koji.FolderPath)
	if err != nil {
		// エラーの場合はデフォルトの必須ファイルを設定
		requiredFile, err := models.NewFileInfo(path.Join(koji.FolderPath, "工事.xlsx"))
		if err != nil {
			return err
		}
		koji.RequiredFiles = []models.FileInfo{*requiredFile}
		return nil
	}

	// 必須ファイルの設定（現在は1個のみ）
	var requiredFile *models.FileInfo
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
		requiredFile, err = models.NewFileInfo(actualPath, standardPath)
		if err != nil {
			continue
		}
		found = true
		break
	}

	// 工事ファイルが見つからない場合はひな形を設定
	if !found {
		standardPath := path.Join(koji.FolderPath, koji.GetFolderName()+".xlsx")
		requiredFile, err = models.NewFileInfo("工事.xlsx", standardPath)

		// テンプレートファイルのコピーは省略（高速化のため）
		// 必要に応じて別途実行する
	}

	// 管理ファイルの設定
	koji.RequiredFiles = []models.FileInfo{*requiredFile}

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
	if temp.UpdateFolderPath() {
		// 工事IDの更新（フォルダー名の変更に伴う）
		temp.UpdateID()

		// 必須ファイル情報の更新（フォルダー名の変更に伴う）
		temp.RequiredFiles = []models.FileInfo{}
		if err := ks.UpdateRequiredFiles(&temp); err != nil {
			return err
		}

		// ファイルの移動
		// TODO: ファイル移動の実装
		// if err := ks.RootService.MoveFile(temp.GetFolderName(), target.GetFolderName()); err != nil {
		// 	return err
		// }
	}

	// フォルダー名の変更完了後、情報を更新
	*target = temp

	// 計算が必要な項目の更新
	target.UpdateStatus()

	// 更新後の工事情報を属性ファイルに反映
	return ks.Repository.Save(target)
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
