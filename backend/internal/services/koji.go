package services

import (
	"fmt"
	"path/filepath"
	"penguin-backend/internal/models"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"
)

// KojiService は工事情報取得サービスを提供する
type KojiService struct {
	FileService      *FileService
	AttributeService *AttributeService[*models.Koji]
	FolderName       string
}

// NewKojiService はKojiServiceを初期化する
// @Param folderName query string true "フォルダー名(FileService.BasePathからの相対パス)" default("2-工事")
func NewKojiService(businessFileService *FileService, folderName string) (*KojiService, error) {
	ks := &KojiService{}

	// フォルダー名を設定
	ks.FolderName = folderName

	// フォルダーのフルパスの取得
	folderPath, err := businessFileService.GetFullpath(folderName)
	if err != nil {
		return nil, err
	}

	// FileServiceを初期化
	ks.FileService, err = NewFileService(folderPath)
	if err != nil {
		return nil, err
	}

	// AttributeServiceを初期化
	ks.AttributeService = NewAttributeService[*models.Koji](ks.FileService, ".detail.yaml")

	return ks, nil
}

// GetKoji は指定されたパスから工事を取得する
// pathは工事フォルダーのファイル名
// 工事を返す
func (ks *KojiService) GetKoji(folderName string) (models.Koji, error) {
	// 工事フォルダーのフルパスを取得
	folderPath, err := ks.FileService.GetFullpath(folderName)
	if err != nil {
		return models.Koji{}, fmt.Errorf("工事フォルダーが見つかりません: %s", folderName)
	}

	// 工事フォルダーのFileInfoを取得
	folderInfo, err := models.NewFileInfo(folderPath)
	if err != nil {
		return models.Koji{}, fmt.Errorf("工事フォルダー情報の取得に失敗しました: %v", err)
	}

	// 工事データモデルを作成
	koji, err := models.NewKoji(*folderInfo)
	if err != nil {
		return models.Koji{}, fmt.Errorf("工事データモデルの作成に失敗しました: %v", err)
	}

	// 管理ファイルの設定
	err = ks.SetManagedFiles(&koji)
	if err != nil {
		return models.Koji{}, fmt.Errorf("管理ファイルの設定に失敗しました: %v", err)
	}

	// 属性ファイルと同期する
	err = ks.SyncAttributeFile(&koji)
	if err != nil {
		return models.Koji{}, fmt.Errorf("属性ファイルとの同期に失敗しました: %v", err)
	}

	return koji, nil
}

// KojiResult 並列処理用の結果構造体
type KojiResult struct {
	Koji  models.Koji
	Index int
	Error error
}

// GetRecentKojies は指定されたパスから最近の工事データモデル一覧を取得する
func (ks *KojiService) GetRecentKojies() []models.Koji {
	// ファイルシステムから工事フォルダー一覧を取得
	fileInfos, err := ks.FileService.GetFileInfos()
	if err != nil {
		return []models.Koji{}
	}

	if len(fileInfos) == 0 {
		return []models.Koji{}
	}

	// 並列処理用のワーカー数を決定（CPU数と同じ、最大8）
	numWorkers := min(runtime.NumCPU(), min(8, len(fileInfos)))

	// チャンネルとワーカーグループを設定
	jobs := make(chan int, len(fileInfos))
	results := make(chan KojiResult, len(fileInfos))
	var wg sync.WaitGroup

	// ワーカーを起動
	for range numWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range jobs {
				fileInfo := fileInfos[index]

				// fileInfoから工事を作成
				koji, err := models.NewKoji(fileInfo)
				if err != nil {
					results <- KojiResult{Index: index, Error: err}
					continue
				}

				// 管理ファイルの設定、エラーは無視
				ks.SetManagedFiles(&koji)

				// 属性ファイルと同期、エラーは無視
				ks.SyncAttributeFile(&koji)

				results <- KojiResult{Koji: koji, Index: index, Error: nil}
			}
		}()
	}

	// ジョブを送信
	for i := range fileInfos {
		jobs <- i
	}
	close(jobs)

	// ワーカーの完了を待つ
	go func() {
		wg.Wait()
		close(results)
	}()

	// 結果を収集（元の順序を保持）
	kojies := make([]models.Koji, 0, len(fileInfos))
	kojiMap := make(map[int]models.Koji, len(fileInfos))

	for result := range results {
		if result.Error == nil {
			kojiMap[result.Index] = result.Koji
		}
	}

	// 元の順序で工事一覧を構築
	for i := range len(fileInfos) {
		if koji, exists := kojiMap[i]; exists {
			kojies = append(kojies, koji)
		}
	}

	// 工事開始日で降順ソート（新しいものが最初）
	sort.Slice(kojies, func(i, j int) bool {
		// 両方の開始日が有効な場合
		if !kojies[i].StartDate.IsZero() && !kojies[j].StartDate.IsZero() {
			return kojies[i].StartDate.After(kojies[j].StartDate.Time)
		}
		// iの開始日が無効でjが有効な場合、jを先に
		if kojies[i].StartDate.IsZero() && !kojies[j].StartDate.IsZero() {
			return false
		}
		// iの開始日が有効でjが無効な場合、iを先に
		if !kojies[i].StartDate.IsZero() && kojies[j].StartDate.IsZero() {
			return true
		}
		// 両方の開始日が無効な場合、フォルダー名で降順
		return kojies[i].FileInfo.Name > kojies[j].FileInfo.Name
	})

	return kojies
}

// SyncAttributeFile は工事の属性データを読み込み、工事データモデルに反映する
func (ks *KojiService) SyncAttributeFile(koji *models.Koji) error {
	// 工事の属性データを読み込む
	attribute, err := ks.AttributeService.Load(koji)
	if err != nil {
		// ファイルが存在しない場合はデフォルト値を設定
		koji.UpdateStatus()

		// 新規ファイルの場合は非同期で保存（必要最小限の書き込み）
		go func() {
			ks.AttributeService.Save(koji)
		}()
		return nil
	}

	// detailにしか保持されない内容をkojiに反映
	koji.Description = attribute.Description
	koji.EndDate = attribute.EndDate
	koji.Tags = attribute.Tags

	// 計算が必要な項目の更新
	koji.UpdateStatus()

	// フォルダー名から解析した情報をdetailに反映
	attribute.StartDate = koji.StartDate
	attribute.CompanyName = koji.CompanyName
	attribute.LocationName = koji.LocationName
	attribute.Status = koji.Status

	// データに変更がある場合のみ非同期で保存（パフォーマンス重視）
	if ks.hasKojiChanged(*attribute, *koji) {
		go func() {
			ks.AttributeService.Save(attribute)
		}()
	}

	return nil
}

// hasKojiChanged 工事データに変更があるかチェック
func (ks *KojiService) hasKojiChanged(detail, current models.Koji) bool {
	return detail.StartDate != current.StartDate ||
		detail.CompanyName != current.CompanyName ||
		detail.LocationName != current.LocationName ||
		detail.Status != current.Status
}

// SetManagedFiles は工事の管理ファイルを設定します
func (ks *KojiService) SetManagedFiles(koji *models.Koji) error {
	// 工事フォルダー内のファイルを取得
	fileInfos, err := ks.FileService.GetFileInfos(koji.FileInfo.Name)
	if err != nil {
		// エラーの場合はデフォルトの管理ファイルを設定
		koji.ManagedFiles = []models.ManagedFile{{
			Recommended: koji.Name + ".xlsx",
			Current:     "工事.xlsx",
		}}
		return nil
	}

	// 管理ファイルの設定（最大1個のみ）
	var managedFile models.ManagedFile
	var found bool

	// 日付パターンのコンパイル（1回のみ）
	datePattern := regexp.MustCompile(`^20\d{2}-\d{4}`)

	// 工事ファイルの検索
	for _, fileInfo := range fileInfos {
		// エクセルファイル以外はスキップ
		if !fileInfo.IsExcel() {
			continue
		}

		// ファイル名の解析
		fileExt := filepath.Ext(fileInfo.Name)
		filenameWithoutExt := strings.TrimSuffix(fileInfo.Name, fileExt)

		// 条件に合うファイルかチェック
		if !datePattern.MatchString(filenameWithoutExt) && filenameWithoutExt != "工事" {
			continue
		}

		recommendedName := koji.Name + fileExt

		// 推奨ファイルと現在のファイルが同じ時は推奨ファイルの表示を不要にする
		if recommendedName == fileInfo.Name {
			recommendedName = ""
		}

		managedFile = models.ManagedFile{
			Recommended: recommendedName,
			Current:     fileInfo.Name,
		}
		found = true
		break
	}

	// 工事ファイルが見つからない場合はひな形を設定
	if !found {
		managedFile = models.ManagedFile{
			Recommended: koji.Name + ".xlsx",
			Current:     "工事.xlsx",
		}

		// テンプレートファイルのコピーは省略（高速化のため）
		// 必要に応じて別途実行する
	}

	// 管理ファイルの設定
	koji.ManagedFiles = []models.ManagedFile{managedFile}

	return nil
}

// Update は工事情報 FileInfoの情報を更新
// kojiの詳細情報で.detail.yamlを更新
// 移動元フォルダー名：FileInfo.Name
// 移動先フォルダー名：StartDate, CompanyName, LocationNameから生成されたフォルダー名
func (ks *KojiService) Update(koji *models.Koji) error {

	// 工事開始日・会社名・現場名からフォルダー名を生成
	generatedFolderName, err := models.GenerateKojiFolderName(koji.StartDate, koji.CompanyName, koji.LocationName)
	if err != nil {
		return err
	}

	// 生成されたフォルダー名とFileInfo.Nameを比較
	if generatedFolderName != koji.FileInfo.Name {
		// フォルダー名の変更
		err = ks.FileService.MoveFile(koji.FileInfo.Name, generatedFolderName)
		if err != nil {
			return err
		}

		// FileInfoの作成
		generatedFullpath, err := ks.FileService.GetFullpath(generatedFolderName)
		if err != nil {
			return err
		}
		generatedFileInfo, err := models.NewFileInfo(generatedFullpath)
		if err != nil {
			return err
		}

		// 詳細情報以外の更新
		koji.FileInfo = *generatedFileInfo
	}

	// 計算が必要な項目の更新
	koji.ID = models.GenerateKojiID(koji.StartDate, koji.CompanyName, koji.LocationName)
	koji.Status = models.DetermineKojiStatus(koji.StartDate, koji.EndDate)

	// 管理ファイルを更新（推奨ファイル名が変更される可能性があるため）
	err = ks.SetManagedFiles(koji)
	if err != nil {
		return err
	}

	// 更新後の工事情報を属性ファイルに反映
	return ks.AttributeService.Save(koji)
}

// KojiesToMapByID は[]KojiをMap[string]Kojiに変換する
func KojiesToMapByID(kojies []models.Koji) map[string]models.Koji {
	kojiMap := make(map[string]models.Koji, len(kojies))
	for _, koji := range kojies {
		kojiMap[koji.ID] = koji
	}
	return kojiMap
}

// RenameManagedFile は管理ファイルの名前を変更し、工事データも更新する
// kojiは工事データ
// currentsは変更前の管理ファイル名
// 変更後の管理ファイル名を返す
func (ks *KojiService) RenameManagedFile(koji models.Koji, currents []string) []string {
	renamedFiles := make([]string, len(currents))

	count := 0
	for _, current := range currents {
		for _, managedFile := range koji.ManagedFiles {
			if managedFile.Current == current {
				currentPath := filepath.Join(koji.FileInfo.Name, current)
				recommendedPath := filepath.Join(koji.FileInfo.Name, managedFile.Recommended)
				err := ks.FileService.MoveFile(currentPath, recommendedPath)
				if err == nil {
					renamedFiles[count] = recommendedPath
					count++
				}
			}
		}
	}

	// ファイル名変更後、工事の管理ファイル情報を更新
	if count > 0 {
		// 管理ファイル情報を再設定
		err := ks.SetManagedFiles(&koji)
		if err == nil {
			// 属性ファイルに反映
			ks.AttributeService.Save(&koji)
		}
	}

	return renamedFiles[:count]
}
