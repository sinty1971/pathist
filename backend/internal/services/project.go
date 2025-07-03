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
	FileService    *FileService
	YamlService    *YamlService[models.Project]
	CompanyService *CompanyService
	DetailFile     string
}

// NewProjectService はProjectServiceを初期化する
// @Param projectsPath query string true "プロジェクト一覧ファイルのパス" default("~/penguin/豊田築炉/2-工事")
func NewProjectService(projectsPath string) (*ProjectService, error) {
	// FileServiceを初期化
	fileService, err := NewFileService(projectsPath)
	if err != nil {
		return nil, err
	}

	// YamlServiceを初期化
	yamlService := NewYamlService[models.Project](fileService)
	
	// CompanyServiceを初期化
	companyService := NewCompanyService()

	return &ProjectService{
		FileService:    fileService,
		YamlService:    yamlService,
		CompanyService: companyService,
		DetailFile:     ".detail.yaml",
	}, nil
}

// GetProject は指定されたパスから工事を取得する
// pathは工事フォルダーのファイル名
// 工事を返す
func (s *ProjectService) GetProject(path string) (models.Project, error) {
	// 工事フォルダー自体の情報を取得
	fullPath, err := s.FileService.GetFullpath(path)
	if err != nil {
		return models.Project{}, fmt.Errorf("工事フォルダーが見つかりません: %s", path)
	}

	fileInfo, err := models.NewFileInfo(fullPath)
	if err != nil {
		return models.Project{}, fmt.Errorf("ファイル情報の取得に失敗しました: %v", err)
	}

	// 工事を作成
	project, err := models.NewProject(*fileInfo)
	if err != nil {
		return models.Project{}, fmt.Errorf("工事データの作成に失敗しました: %v", err)
	}

	// 管理ファイルの設定
	err = s.SetManagedFiles(&project)
	if err != nil {
		return models.Project{}, fmt.Errorf("管理ファイルの設定に失敗しました: %v", err)
	}

	// .detail.yamlと同期する
	err = s.SyncDetailFile(&project)
	if err != nil {
		return models.Project{}, fmt.Errorf(".detail.yamlの同期に失敗しました: %v", err)
	}

	// 企業情報を取得
	s.EnrichProjectWithCompanyInfo(&project)

	return project, nil
}

// ProjectResult 並列処理用の結果構造体
type ProjectResult struct {
	Project models.Project
	Index   int
	Error   error
}

// GetRecentProjects は指定されたパスから最近の工事一覧を取得する（高速並行処理版）
func (s *ProjectService) GetRecentProjects() []models.Project {
	// ファイルシステムから工事一覧を取得
	fileInfos, err := s.FileService.GetFileInfos()
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

				// 管理ファイルの設定（エラーは無視して継続）
				s.setManagedFilesOptimized(&project)

				// .detail.yamlと同期する（エラーは無視して継続）
				s.syncDetailFileOptimized(&project)

				// 企業情報の取得（エラーは無視して継続）
				s.EnrichProjectWithCompanyInfo(&project)

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

// SyncDetailFile は工事の詳細を.detail.yamlから取り込む
func (s *ProjectService) SyncDetailFile(current *models.Project) error {
	return s.syncDetailFileOptimized(current)
}

// syncDetailFileOptimized は工事の詳細を高速で.detail.yamlから取り込む（内部用）
func (s *ProjectService) syncDetailFileOptimized(current *models.Project) error {
	// .detail.yamlを読み込む
	detailPath := filepath.Join(current.FileInfo.Name, s.DetailFile)
	detail, err := s.YamlService.Load(detailPath)
	if err != nil {
		// ファイルが存在しない場合はデフォルト値を設定
		current.Status = models.DetermineProjectStatus(current.StartDate, current.EndDate)
		
		// 新規ファイルの場合は非同期で保存（必要最小限の書き込み）
		go func() {
			s.YamlService.Save(detailPath, *current)
		}()
		return nil
	}

	// detailにしか保持されない内容をprojectに反映
	current.Description = detail.Description
	current.EndDate = detail.EndDate
	current.Tags = detail.Tags

	// 計算が必要な項目の更新
	current.Status = models.DetermineProjectStatus(current.StartDate, current.EndDate)

	// フォルダー名から解析した情報をdetailに反映
	detail.StartDate = current.StartDate
	detail.CompanyName = current.CompanyName
	detail.LocationName = current.LocationName
	detail.Status = current.Status

	// データに変更がある場合のみ非同期で保存（パフォーマンス重視）
	if s.hasProjectChanged(detail, *current) {
		go func() {
			s.YamlService.Save(detailPath, detail)
		}()
	}

	return nil
}

// hasProjectChanged プロジェクトデータに変更があるかチェック
func (s *ProjectService) hasProjectChanged(detail, current models.Project) bool {
	return detail.StartDate != current.StartDate ||
		detail.CompanyName != current.CompanyName ||
		detail.LocationName != current.LocationName ||
		detail.Status != current.Status
}

// SetManagedFiles は工事の管理ファイルを設定します
func (s *ProjectService) SetManagedFiles(project *models.Project) error {
	return s.setManagedFilesOptimized(project)
}

// setManagedFilesOptimized は工事の管理ファイルを高速設定します（内部用）
func (s *ProjectService) setManagedFilesOptimized(project *models.Project) error {
	// 工事フォルダー内のファイルを取得
	fileInfos, err := s.FileService.GetFileInfos(project.FileInfo.Name)
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
func (s *ProjectService) UpdateProjectFileInfo(project *models.Project) error {

	// 工事開始日・会社名・現場名からフォルダー名を生成
	updatesStartDate, err := project.StartDate.Format("2006-0102")
	if err != nil {
		return err
	}
	generatedFolderName := fmt.Sprintf("%s %s %s", updatesStartDate, project.CompanyName, project.LocationName)

	// 生成されたフォルダー名とFileInfo.Nameを比較
	if generatedFolderName != project.FileInfo.Name {
		// フォルダー名の変更
		err = s.FileService.MoveFile(project.FileInfo.Name, generatedFolderName)
		if err != nil {
			return err
		}

		// FileInfoの作成
		generatedFullpath, err := s.FileService.GetFullpath(generatedFolderName)
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
	err = s.SetManagedFiles(project)
	if err != nil {
		return err
	}

	// 更新後の工事情報を.detail.yamlに反映
	detailPath := filepath.Join(project.FileInfo.Name, s.DetailFile)
	return s.YamlService.Save(detailPath, *project)
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
func (s *ProjectService) RenameManagedFile(project models.Project, currents []string) []string {
	renamedFiles := make([]string, len(currents))

	count := 0
	for _, current := range currents {
		for _, managedFile := range project.ManagedFiles {
			if managedFile.Current == current {
				currentPath := filepath.Join(project.FileInfo.Name, current)
				recommendedPath := filepath.Join(project.FileInfo.Name, managedFile.Recommended)
				fmt.Printf("current: %s, recommended: %s\n", currentPath, recommendedPath)
				err := s.FileService.MoveFile(currentPath, recommendedPath)
				if err == nil {
					renamedFiles[count] = recommendedPath
					count++
				}
			}
		}
	}
	return renamedFiles[:count]
}

// EnrichProjectWithCompanyInfo enriches a project with company information
func (s *ProjectService) EnrichProjectWithCompanyInfo(project *models.Project) error {
	if project.CompanyName == "" {
		return nil // No company name to lookup
	}
	
	// Try to get or create company by name
	company, err := s.CompanyService.GetOrCreateCompanyByName(project.CompanyName)
	if err != nil {
		return fmt.Errorf("failed to get/create company %s: %v", project.CompanyName, err)
	}
	
	// Update project with company information
	project.Company = company
	project.CompanyID = company.ID
	
	return nil
}

// EnrichProjectsWithCompanyInfo enriches multiple projects with company information
func (s *ProjectService) EnrichProjectsWithCompanyInfo(projects []models.Project) ([]models.Project, error) {
	enrichedProjects := make([]models.Project, len(projects))
	
	for i, project := range projects {
		enrichedProjects[i] = project
		if err := s.EnrichProjectWithCompanyInfo(&enrichedProjects[i]); err != nil {
			// Log error but continue with other projects
			fmt.Printf("Warning: failed to enrich project %s with company info: %v\n", project.ID, err)
		}
	}
	
	return enrichedProjects, nil
}

// UpdateProjectCompanyInfo updates company information for a project
func (s *ProjectService) UpdateProjectCompanyInfo(project *models.Project, companyInfo models.Company) error {
	// Update or create the company
	existing, err := s.CompanyService.GetCompanyByID(companyInfo.ID)
	if err != nil || existing == nil {
		// Create new company
		if err := s.CompanyService.CreateCompany(companyInfo); err != nil {
			return fmt.Errorf("failed to create company: %v", err)
		}
	} else {
		// Update existing company
		if err := s.CompanyService.UpdateCompany(companyInfo); err != nil {
			return fmt.Errorf("failed to update company: %v", err)
		}
	}
	
	// Update project reference
	project.CompanyID = companyInfo.ID
	project.CompanyName = companyInfo.Name
	project.Company = &companyInfo
	
	return nil
}
