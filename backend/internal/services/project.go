package services

import (
	"fmt"
	"path/filepath"
	"penguin-backend/internal/models"
	"sort"
)

// ProjectService は工事情報取得サービスを提供する
type ProjectService struct {
	FileService *FileService
	YamlService *YamlService[models.Project]
}

// NewProjectService はProjectServiceを初期化する
// @param projectsPath: プロジェクト一覧ファイルのパス (ex. ~/penguin/豊田築炉/2-工事)
func NewProjectService(projectsPath string) (*ProjectService, error) {
	// FileServiceを初期化
	fileService, err := NewFileService(projectsPath)
	if err != nil {
		return nil, err
	}

	// YamlServiceを初期化
	yamlService := NewYamlService[models.Project](fileService)
	if err != nil {
		return nil, err
	}

	return &ProjectService{
		FileService: fileService,
		YamlService: yamlService,
	}, nil
}

// GetFullpath は工事一覧フォルダーを基準としてファイルのフルパスを取得する
func (s *ProjectService) GetFullpath(joinPath ...string) (string, error) {
	return s.FileService.GetFullpath(joinPath...)
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

		// 工事の詳細を.detail.yamlか取り込む
		s.ImportDetailFromYaml(&project)

		// 工事を登録
		projects[count] = project
		count++

	}

	return projects[:count]
}

// ImportDetailFromYaml は工事の詳細を.detail.yamlから取り込む
// また、工事詳細内の情報も更新する
func (s *ProjectService) ImportDetailFromYaml(target *models.Project) error {
	// 工事のフルパスを取得
	targetFullpath, err := s.GetFullpath(target.FileInfo.Name)
	if err != nil {
		return err
	}

	// .detail.yamlを読み込む
	dotDetailPath := filepath.Join(targetFullpath, ".detail.yaml")
	detail, err := s.YamlService.Load(dotDetailPath)
	if err != nil {
		return err
	}

	// detailにしか保持されない内容をprojectに反映
	target.Description = detail.Description
	target.EndDate = detail.EndDate
	target.Tags = detail.Tags

	// 計算が必要な項目の更新
	target.Status = models.DetermineProjectStatus(target.StartDate, target.EndDate)

	// 工事詳細を更新
	err = s.YamlService.Save(dotDetailPath, detail)
	if err != nil {
		return err
	}

	return nil
}

// 工事情報 from を工事情報 to に更新する
func (s *ProjectService) UpdateProject(from, to models.Project) error {

	// from フォルダーが存在するか確認
	if !from.FileInfo.IsExist() {
		return fmt.Errorf("from フォルダーが存在しません: %s", from.FileInfo.Path)
	}

	// to フォルダーが存在するか確認
	if !to.FileInfo.IsExist() {
		return fmt.Errorf("to フォルダーが存在しません: %s", to.FileInfo.Path)
	}
	// データベースの工事一覧をmapに変換する
	dbEntryMap := KoujiEntriesToMapByID(dbEntries)

	// ファイルシステムの工事一覧を更新する
	margedEntries := make([]models.Project, len(fsEntries))
	for i, fsEntry := range fsEntries {
		// データベースに情報が存在しているときの処理
		if dbEntry, exists := dbEntryMap[fsEntry.ID]; exists {
			// 工事開始日の更新
			fsEntry.StartDate = dbEntry.StartDate
			// 工事終了日の更新
			fsEntry.EndDate = dbEntry.EndDate
			// 工事説明の更新
			fsEntry.Description = dbEntry.Description
			// 工事ステータスの更新
			fsEntry.Status = models.DetermineProjectStatus(dbEntry.StartDate, dbEntry.EndDate)
			// 工事タグの更新
			fsEntry.Tags = dbEntry.Tags

			// データベースから削除（検索スピードを向上させるため）
			delete(dbEntryMap, fsEntry.ID)
		}
		// ファイルシステムの工事を更新
		margedEntries[i] = fsEntry
	}

	// 開始日の降順でソート（新しい順）
	sort.Slice(margedEntries, func(i, j int) bool {
		return margedEntries[i].StartDate.Time.After(margedEntries[j].StartDate.Time)
	})

	return margedEntries
}

// KoujiEntriesToMapByID は[]KoujiEntryをmap[string]KoujiEntryに変換する
func KoujiEntriesToMapByID(entries []models.Project) map[string]models.Project {
	entryMap := make(map[string]models.Project, len(entries))
	for _, entry := range entries {
		entryMap[entry.ID] = entry
	}
	return entryMap
}

// UpdateEntry は工事情報を更新する
func (s *ProjectService) UpdateEntry(updateEntry models.Project) (models.Project, error) {
	// データベースから工事一覧を取得
	dbEntries, err := s.Database.LoadFromYAML()
	if err != nil {
		return models.Project{}, err
	}

	// Find and update the project
	foundIndex := -1
	for i, entry := range dbEntries {
		if entry.ID == updateEntry.ID {
			dbEntries[i].StartDate = updateEntry.StartDate
			dbEntries[i].EndDate = updateEntry.EndDate
			dbEntries[i].Status = DetermineKoujiStatus(updateEntry.StartDate, updateEntry.EndDate)
			foundIndex = i
			break
		}
	}

	if foundIndex == -1 {
		return models.Project{}, fmt.Errorf("工事情報がデータベースにありません: %s", updateEntry.ID)
	}

	// Save updated kouji entries to YAML
	err = s.SaveToYAML(dbEntries)

	// 更新した工事を返す
	return dbEntries[foundIndex], err
}

// SaveToYAML は引数<entries>をデータベースに保存する
func (s *ProjectService) SaveToYAML(entries []models.Project) error {

	// データベースから工事一覧を取得
	fsEntries := s.GetKoujiEntriesFromFileSystem()

	//
	margedEntries := MargeProjects(fsEntries, entries)

	return s.Database.SaveToYAML(margedEntries)
}

// checkFile はファイル名が規則に沿っているかを返す
// "新規工事.xlsx"は規則に沿っていいるがテンプレートのため"<工事フォルダー名>.xlsx"への更新が必要
// 他のファイルは未対応
func checkFile(file string) bool {
	// ファイル名が規則に沿っているかを返す
	return true
}
