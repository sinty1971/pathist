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

// ProjectService は工事情報取得サービスを提供する
type ProjectService struct {
	FileService      *FileService
	AttributeService *AttributeService[*models.Project]
	FolderName       string
}

// NewProjectService はProjectServiceを初期化する
// @Param folderName query string true "フォルダー名(FileService.BasePathからの相対パス)" default("2-工事")
func NewProjectService(businessFileService *FileService, folderName string) (*ProjectService, error) {
	ps := &ProjectService{}

	// フォルダー名を設定
	ps.FolderName = folderName

	// フォルダーのフルパスの取得
	folderPath, err := businessFileService.GetFullpath(folderName)
	if err != nil {
		return nil, err
	}

	// FileServiceを初期化
	ps.FileService, err = NewFileService(folderPath)
	if err != nil {
		return nil, err
	}

	// AttributeServiceを初期化
	ps.AttributeService = NewAttributeService[*models.Project](ps.FileService, ".detail.yaml")

	return ps, nil
}

// GetProject は指定されたパスから工事を取得する
// pathは工事フォルダーのファイル名
// 工事を返す
func (ps *ProjectService) GetProject(folderName string) (models.Project, error) {
	// 工事フォルダーのフルパスを取得
	folderPath, err := ps.FileService.GetFullpath(folderName)
	if err != nil {
		return models.Project{}, fmt.Errorf("工事フォルダーが見つかりません: %s", folderName)
	}

	// 工事フォルダーのFileInfoを取得
	folderInfo, err := models.NewFileInfo(folderPath)
	if err != nil {
		return models.Project{}, fmt.Errorf("工事フォルダー情報の取得に失敗しました: %v", err)
	}

	// 工事データモデルを作成
	project, err := models.NewProject(*folderInfo)
	if err != nil {
		return models.Project{}, fmt.Errorf("工事データモデルの作成に失敗しました: %v", err)
	}

	// 管理ファイルの設定
	err = ps.SetManagedFiles(&project)
	if err != nil {
		return models.Project{}, fmt.Errorf("管理ファイルの設定に失敗しました: %v", err)
	}

	// 属性ファイルと同期する
	err = ps.SyncAttributeFile(&project)
	if err != nil {
		return models.Project{}, fmt.Errorf("属性ファイルとの同期に失敗しました: %v", err)
	}

	return project, nil
}

// ProjectResult 並列処理用の結果構造体
type ProjectResult struct {
	Project models.Project
	Index   int
	Error   error
}

// GetRecentProjects は指定されたパスから最近の工事データモデル一覧を取得する
func (ps *ProjectService) GetRecentProjects() []models.Project {
	// ファイルシステムから工事フォルダー一覧を取得
	fileInfos, err := ps.FileService.GetFileInfos()
	if err != nil {
		return []models.Project{}
	}

	if len(fileInfos) == 0 {
		return []models.Project{}
	}

	// 並列処理用のワーカー数を決定（CPU数と同じ、最大8）
	numWorkers := min(runtime.NumCPU(), min(8, len(fileInfos)))

	// チャンネルとワーカーグループを設定
	jobs := make(chan int, len(fileInfos))
	results := make(chan ProjectResult, len(fileInfos))
	var wg sync.WaitGroup

	// ワーカーを起動
	for range numWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range jobs {
				fileInfo := fileInfos[index]

				// fileInfoから工事を作成
				project, err := models.NewProject(fileInfo)
				if err != nil {
					results <- ProjectResult{Index: index, Error: err}
					continue
				}

				// 管理ファイルの設定、エラーは無視
				ps.SetManagedFiles(&project)

				// 属性ファイルと同期、エラーは無視
				ps.SyncAttributeFile(&project)

				results <- ProjectResult{Project: project, Index: index, Error: nil}
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
	projects := make([]models.Project, 0, len(fileInfos))
	projectMap := make(map[int]models.Project, len(fileInfos))

	for result := range results {
		if result.Error == nil {
			projectMap[result.Index] = result.Project
		}
	}

	// 元の順序でプロジェクト一覧を構築
	for i := range len(fileInfos) {
		if project, exists := projectMap[i]; exists {
			projects = append(projects, project)
		}
	}

	// 工事開始日で降順ソート（新しいものが最初）
	sort.Slice(projects, func(i, j int) bool {
		// 両方の開始日が有効な場合
		if !projects[i].StartDate.IsZero() && !projects[j].StartDate.IsZero() {
			return projects[i].StartDate.After(projects[j].StartDate.Time)
		}
		// iの開始日が無効でjが有効な場合、jを先に
		if projects[i].StartDate.IsZero() && !projects[j].StartDate.IsZero() {
			return false
		}
		// iの開始日が有効でjが無効な場合、iを先に
		if !projects[i].StartDate.IsZero() && projects[j].StartDate.IsZero() {
			return true
		}
		// 両方の開始日が無効な場合、フォルダー名で降順
		return projects[i].FileInfo.Name > projects[j].FileInfo.Name
	})

	return projects
}

// SyncAttributeFile は工事の属性データを読み込み、工事データモデルに反映する
func (ps *ProjectService) SyncAttributeFile(project *models.Project) error {
	// 工事の属性データを読み込む
	attribute, err := ps.AttributeService.Load(project)
	if err != nil {
		// ファイルが存在しない場合はデフォルト値を設定
		project.UpdateStatus()

		// 新規ファイルの場合は非同期で保存（必要最小限の書き込み）
		go func() {
			ps.AttributeService.Save(project)
		}()
		return nil
	}

	// detailにしか保持されない内容をprojectに反映
	project.Description = attribute.Description
	project.EndDate = attribute.EndDate
	project.Tags = attribute.Tags

	// 計算が必要な項目の更新
	project.UpdateStatus()

	// フォルダー名から解析した情報をdetailに反映
	attribute.StartDate = project.StartDate
	attribute.CompanyName = project.CompanyName
	attribute.LocationName = project.LocationName
	attribute.Status = project.Status

	// データに変更がある場合のみ非同期で保存（パフォーマンス重視）
	if ps.hasProjectChanged(*attribute, *project) {
		go func() {
			ps.AttributeService.Save(attribute)
		}()
	}

	return nil
}

// hasProjectChanged プロジェクトデータに変更があるかチェック
func (ps *ProjectService) hasProjectChanged(detail, current models.Project) bool {
	return detail.StartDate != current.StartDate ||
		detail.CompanyName != current.CompanyName ||
		detail.LocationName != current.LocationName ||
		detail.Status != current.Status
}

// SetManagedFiles は工事の管理ファイルを設定します
func (ps *ProjectService) SetManagedFiles(project *models.Project) error {
	// 工事フォルダー内のファイルを取得
	fileInfos, err := ps.FileService.GetFileInfos(project.FileInfo.Name)
	if err != nil {
		// エラーの場合はデフォルトの管理ファイルを設定
		project.ManagedFiles = []models.ManagedFile{{
			Recommended: project.Name + ".xlsx",
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

		recommendedName := project.Name + fileExt

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
			Recommended: project.Name + ".xlsx",
			Current:     "工事.xlsx",
		}

		// テンプレートファイルのコピーは省略（高速化のため）
		// 必要に応じて別途実行する
	}

	// 管理ファイルの設定
	project.ManagedFiles = []models.ManagedFile{managedFile}

	return nil
}

// UpdateProjectFileInfo は工事情報 project.FileInfoの情報を更新します
// projectの詳細情報で.detail.yamlを更新します
func (ps *ProjectService) UpdateProjectFileInfo(project *models.Project) error {

	// 工事開始日・会社名・現場名からフォルダー名を生成
	updatesStartDate, err := project.StartDate.Format("2006-0102")
	if err != nil {
		return err
	}
	generatedFolderName := fmt.Sprintf("%s %s %s", updatesStartDate, project.CompanyName, project.LocationName)

	// 生成されたフォルダー名とFileInfo.Nameを比較
	if generatedFolderName != project.FileInfo.Name {
		// フォルダー名の変更
		err = ps.FileService.MoveFile(project.FileInfo.Name, generatedFolderName)
		if err != nil {
			return err
		}

		// FileInfoの作成
		generatedFullpath, err := ps.FileService.GetFullpath(generatedFolderName)
		if err != nil {
			return err
		}
		generatedFileInfo, err := models.NewFileInfo(generatedFullpath)
		if err != nil {
			return err
		}

		// 詳細情報以外の更新
		project.FileInfo = *generatedFileInfo
	}

	// 計算が必要な項目の更新
	project.ID = models.GenerateProjectID(project.StartDate, project.CompanyName, project.LocationName)
	project.Status = models.DetermineProjectStatus(project.StartDate, project.EndDate)

	// 管理ファイルを更新（推奨ファイル名が変更される可能性があるため）
	err = ps.SetManagedFiles(project)
	if err != nil {
		return err
	}

	// 更新後の工事情報を属性ファイルに反映
	return ps.AttributeService.Save(project)
}

// ProjectsToMapByID は[]Projectをmap[string]Projectに変換する
func ProjectsToMapByID(projects []models.Project) map[string]models.Project {
	projectMap := make(map[string]models.Project, len(projects))
	for _, project := range projects {
		projectMap[project.ID] = project
	}
	return projectMap
}

// RenameManagedFile は管理ファイルの名前を変更する
// projectは工事データ
// currentsは変更前の管理ファイル名
// 変更後の管理ファイル名を返す
func (ps *ProjectService) RenameManagedFile(project models.Project, currents []string) []string {
	renamedFiles := make([]string, len(currents))

	count := 0
	for _, current := range currents {
		for _, managedFile := range project.ManagedFiles {
			if managedFile.Current == current {
				currentPath := filepath.Join(project.FileInfo.Name, current)
				recommendedPath := filepath.Join(project.FileInfo.Name, managedFile.Recommended)
				fmt.Printf("current: %s, recommended: %s\n", currentPath, recommendedPath)
				err := ps.FileService.MoveFile(currentPath, recommendedPath)
				if err == nil {
					renamedFiles[count] = recommendedPath
					count++
				}
			}
		}
	}
	return renamedFiles[:count]
}
