package services

import (
	"fmt"
	"path/filepath"
	"penguin-backend/internal/models"
	"regexp"
	"sort"
	"strings"
)

// ProjectService は工事情報取得サービスを提供する
type ProjectService struct {
	FileService *FileService
	YamlService *YamlService[models.Project]
	DetailFile  string
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

	return &ProjectService{
		FileService: fileService,
		YamlService: yamlService,
		DetailFile:  ".detail.yaml",
	}, nil
}

// GetProject は指定されたパスから工事を取得する
// pathは工事フォルダーのファイル名
// 工事を返す
func (s *ProjectService) GetProject(path string) (models.Project, error) {
	fmt.Printf("GetProject called with path: '%s'\n", path)

	// 工事フォルダー自体の情報を取得
	fullPath, err := s.FileService.GetFullpath(path)
	if err != nil {
		fmt.Printf("GetFullpath error for path '%s': %v\n", path, err)
		return models.Project{}, fmt.Errorf("工事フォルダーが見つかりません: %s", path)
	}
	fmt.Printf("Full path: %s\n", fullPath)

	fileInfo, err := models.NewFileInfo(fullPath)
	if err != nil {
		fmt.Printf("NewFileInfo error for fullPath '%s': %v\n", fullPath, err)
		return models.Project{}, fmt.Errorf("ファイル情報の取得に失敗しました: %v", err)
	}

	// 工事を作成
	project, err := models.NewProject(*fileInfo)
	if err != nil {
		fmt.Printf("NewProject error: %v\n", err)
		return models.Project{}, fmt.Errorf("工事データの作成に失敗しました: %v", err)
	}

	// 管理ファイルの設定
	err = s.SetManagedFiles(&project)
	if err != nil {
		fmt.Printf("SetManagedFiles error: %v\n", err)
		return models.Project{}, fmt.Errorf("管理ファイルの設定に失敗しました: %v", err)
	}

	// .detail.yamlと同期する
	err = s.SyncDetailFile(&project)
	if err != nil {
		fmt.Printf("SyncDetailFile error: %v\n", err)
		return models.Project{}, fmt.Errorf(".detail.yamlの同期に失敗しました: %v", err)
	}

	return project, nil
}

// GetRecentProjects は指定されたパスから最近の工事一覧を取得する
func (s *ProjectService) GetRecentProjects() []models.Project {
	// ファイルシステムから工事一覧を取得
	fileInfos, err := s.FileService.GetFileInfos()
	if err != nil {
		return []models.Project{}
	}

	// 工事一覧を作成
	projects := make([]models.Project, len(fileInfos))
	count := 0
	for _, fileInfo := range fileInfos {
		// fileInfoから工事を作成
		project, err := models.NewProject(fileInfo)
		if err != nil {
			continue
		}

		// 管理ファイルの設定
		s.SetManagedFiles(&project)

		// .detail.yamlと同期する
		s.SyncDetailFile(&project)

		// 工事を登録
		projects[count] = project
		count++

	}

	// 工事開始日で降順ソート（新しいものが最初）
	validProjects := projects[:count]
	sort.Slice(validProjects, func(i, j int) bool {
		// 両方の開始日が有効な場合
		if !validProjects[i].StartDate.IsZero() && !validProjects[j].StartDate.IsZero() {
			return validProjects[i].StartDate.After(validProjects[j].StartDate.Time)
		}
		// iの開始日が無効でjが有効な場合、jを先に
		if validProjects[i].StartDate.IsZero() && !validProjects[j].StartDate.IsZero() {
			return false
		}
		// iの開始日が有効でjが無効な場合、iを先に
		if !validProjects[i].StartDate.IsZero() && validProjects[j].StartDate.IsZero() {
			return true
		}
		// 両方の開始日が無効な場合、フォルダー名で降順
		return validProjects[i].FileInfo.Name > validProjects[j].FileInfo.Name
	})

	return validProjects
}

// SyncDetailFile は工事の詳細を.detail.yamlから取り込む
// 取り込む詳細情報は、工事完了日・工事説明・工事タグ
// Statusは取り込み後に計算する
// 更新後の工事情報を.detail.yamlに反映する。
// 更新後の工事情報を返す。
func (s *ProjectService) SyncDetailFile(current *models.Project) error {

	// .detail.yamlを読み込む
	detailPath := filepath.Join(current.FileInfo.Name, s.DetailFile)
	detail, err := s.YamlService.Load(detailPath)
	if err != nil {
		// ファイルが存在しない場合は新規作成
		err = s.YamlService.Save(detailPath, *current)
		if err != nil {
			return err
		}
		return nil
	}

	// detailにしか保持されない内容をprojectに反映
	current.Description = detail.Description
	current.EndDate = detail.EndDate
	current.Tags = detail.Tags

	// 計算が必要な項目の更新
	current.Status = models.DetermineProjectStatus(current.StartDate, current.EndDate)

	// 工事詳細を更新
	err = s.YamlService.Save(detailPath, detail)
	if err != nil {
		return err
	}

	return nil
}

// SetManagedFiles は工事の管理ファイルを設定します
func (s *ProjectService) SetManagedFiles(project *models.Project) error {

	// 工事フォルダーのフルパスを取得
	projectFullpath, err := s.FileService.GetFullpath(project.FileInfo.Name)
	if err != nil {
		return err
	}

	// 工事フォルダー内のファイルを取得
	fileInfos, err := s.FileService.GetFileInfos(project.FileInfo.Name)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 管理ファイルの設定
	managedFiles := make([]models.ManagedFile, len(fileInfos))

	// 工事ファイルの設定
	count := 0
	datePattern := regexp.MustCompile(`^20\d{2}-\d{4}`)
	for _, fileInfo := range fileInfos {
		// ファイルの拡張子
		fileExt := filepath.Ext(fileInfo.Name)

		// 拡張子を除いたファイル名
		filenameWithoutExt := strings.TrimSuffix(fileInfo.Name, fileExt)

		// エクセルファイル以外はスキップ
		if !fileInfo.IsExcel() {
			continue
		}

		// 拡張子以外のファイル名がyyyy-mmddで始まる場合または工事のファイル名以外はスキップ
		if !datePattern.MatchString(filenameWithoutExt) && filenameWithoutExt != "工事" {
			continue
		}

		recommendedName := strings.Join([]string{project.Name, fileExt}, "")

		// 推奨ファイルと現在のファイルが同じ時は推奨ファイルの表示を不必要と分かるようにする
		if recommendedName == fileInfo.Name {
			recommendedName = ""
		}

		// 管理ファイルの設定
		managedFiles[count] = models.ManagedFile{
			Recommended: recommendedName,
			Current:     fileInfo.Name,
		}
		count++
		break
	}

	// 工事ファイルが見つからない場合はひな形を作成
	if count == 0 {
		recommendedName := strings.Join([]string{project.Name, ".xlsx"}, "")

		managedFiles[0] = models.ManagedFile{
			Recommended: recommendedName,
			Current:     "工事.xlsx",
		}

		// テンプレートファイルのパスを取得
		templatePath, err := s.FileService.GetFullpath("@工事 新規テンプレート", "工事.xlsx")
		if err != nil {
			return err
		}

		// コピー先のパス（工事フォルダ内の "工事.xlsx"）
		dstPath := filepath.Join(projectFullpath, "工事.xlsx")

		// テンプレートファイルをコピー
		err = s.FileService.CopyFile(templatePath, dstPath)
		if err != nil {
			return err
		}

		count = 1
	}

	// 管理ファイルの設定
	project.ManagedFiles = managedFiles[:count]

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
