package services

import (
	"fmt"
	"penguin-backend/internal/models"
	"sort"
	"strings"
	"time"
)

// KoujiService は工事情報取得サービスを提供する
type KoujiService struct {
	FileSystem *FileService
	Database   *DatabaseService[[]models.KoujiEntry]
}

// NewKoujiService はKoujiServiceを初期化する
// @param listFolderPath 工事一覧フォルダーのパス
// @param databasePath 工事データベースのパス
// @return KoujiService 工事サービス
func NewKoujiService(fileSystemPath string, databasePath string) (*KoujiService, error) {
	// ファイルサービスを初期化
	fileSystemService, err := NewFileSystemService(fileSystemPath)
	if err != nil {
		return nil, err
	}

	// データベースサービスを初期化
	databaseService, err := NewDatabaseService[[]models.KoujiEntry](databasePath)
	if err != nil {
		return nil, err
	}

	return &KoujiService{
		FileSystem: fileSystemService,
		Database:   databaseService,
	}, nil
}

// NewKoujiEntry 工事詳細フォルダーから工事詳細エントリを作成します
func NewKoujiEntry(folderEntry models.FileEntry) (models.KoujiEntry, error) {

	if !folderEntry.IsDirectory {
		return models.KoujiEntry{}, fmt.Errorf("ファイルではありません")
	}

	// Parse date from folder name
	startDate := models.Timestamp{}
	restStr, err := models.ParseTimestampAndRest(folderEntry.Name, &startDate)
	if err != nil {
		return models.KoujiEntry{}, err
	}

	// restStr is like "豊田築炉 名和工場 詳細"
	// companyName は会社名(ex. 豊田築炉)
	// locationName は会社名以降の文字列(ex. "名和工場 詳細")
	parts := strings.Split(restStr, " ")
	companyName := parts[0]
	locationName := ""
	if len(parts) > 1 {
		locationName = strings.Join(parts[1:], " ")
	}

	// Generate unique project ID using folder creation date, company name, and location name
	// This ensures the same project always gets the same ID
	// Use project date instead of creation date for more stable ID generation
	idSource := fmt.Sprintf("%d%s%s", folderEntry.ID, companyName, locationName)
	id := models.NewIDFromString(idSource)

	koujiEntry := models.KoujiEntry{
		// Generate project metadata based on folder name
		ID:           id.Len5(),
		CompanyName:  companyName,
		LocationName: locationName,
		StartDate:    startDate,
		EndDate:      startDate,
		Status:       DetermineKoujiStatus(startDate, startDate),
		Description:  companyName + "の" + locationName + "における工事プロジェクト",
		Tags:         []string{"工事", companyName, locationName, startDate.Time.Format("2006")}, // Include year as tag
		// FileEntry: ファイルシステムから取得したフォルダー情報
		FileEntry: folderEntry,
	}

	return koujiEntry, nil
}

// GetKoujiEntries は指定されたパスから工事一覧を取得する（ファイルシステムとデータベースをマージ）
func (s *KoujiService) GetKoujiEntries() []models.KoujiEntry {
	// ファイルシステムから工事を取得
	fsEntries := s.GetKoujiEntriesFromFileSystem()

	// データベースから工事を取得
	dbEntries, err := s.Database.LoadFromYAML()
	if err != nil {
		dbEntries = []models.KoujiEntry{}
	}

	// データベースの工事一覧をmapに変換する
	dbEntryMap := KoujiEntriesToMapByID(dbEntries)

	// ファイルシステムの工事一覧を更新する
	margedEntries := make([]models.KoujiEntry, len(fsEntries))
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
			fsEntry.Status = DetermineKoujiStatus(dbEntry.StartDate, dbEntry.EndDate)
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

	// データベースに保存
	s.SaveToYAML(margedEntries)

	return margedEntries
}

// GetKoujiEntriesFromFileSystem はファイルシステムから工事一覧を取得する
func (s *KoujiService) GetKoujiEntriesFromFileSystem() []models.KoujiEntry {
	// ファイルシステムから工事一覧を取得
	response, err := s.FileSystem.GetFileEntries()
	if err != nil {
		return []models.KoujiEntry{}
	}

	// 工事一覧を作成
	entries := make([]models.KoujiEntry, 0, len(response.FileEntries))
	entryCount := 0
	for _, entry := range response.FileEntries {
		if koujiEntry, err := NewKoujiEntry(entry); err == nil {
			entries[entryCount] = koujiEntry
			entryCount++
		}
	}

	return entries[:entryCount]
}

// ConvertToMap は[]KoujiEntryをmap[string]KoujiEntryに変換する
func KoujiEntriesToMapByID(entries []models.KoujiEntry) map[string]models.KoujiEntry {
	entryMap := make(map[string]models.KoujiEntry, len(entries))
	for _, entry := range entries {
		entryMap[entry.ID] = entry
	}
	return entryMap
}

// DetermineKoujiStatus determines the project status based on the date
func DetermineKoujiStatus(startDate models.Timestamp, endDate models.Timestamp) string {
	if startDate.Time.IsZero() {
		return "不明"
	}

	now := time.Now()

	if now.Before(startDate.Time) {
		return "予定"
	} else if now.After(endDate.Time) {
		return "完了"
	} else {
		return "進行中"
	}
}

// UpdateEntry は工事情報を更新する
func (s *KoujiService) UpdateEntry(updateEntry models.KoujiEntry) (models.KoujiEntry, error) {
	// データベースから工事一覧を取得
	dbEntries, err := s.Database.LoadFromYAML()
	if err != nil {
		return models.KoujiEntry{}, err
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
		return models.KoujiEntry{}, fmt.Errorf("工事情報がデータベースにありません: %s", updateEntry.ID)
	}

	// Save updated kouji entries to YAML
	err = s.SaveToYAML(dbEntries)

	// 更新した工事を返す
	return dbEntries[foundIndex], err
}

// SaveToYAML は引数<entries>をデータベースに保存する
func (s *KoujiService) SaveToYAML(entries []models.KoujiEntry) error {
	// 引数の工事一覧をmapに変換する
	entryMap := KoujiEntriesToMapByID(entries)

	// データベースから工事一覧を取得
	fsEntries := s.GetKoujiEntriesFromFileSystem()

	// ファイルシステムの工事一覧を更新する
	updatedEntries := make([]models.KoujiEntry, len(fsEntries))
	entryCount := 0
	for _, fsEntry := range fsEntries {
		// データベースに情報が存在しているときの処理
		if dbEntry, exists := entryMap[fsEntry.ID]; exists {
			// 工事開始日の更新
			fsEntry.StartDate = dbEntry.StartDate
			// 工事終了日の更新
			fsEntry.EndDate = dbEntry.EndDate
			// 工事説明の更新
			fsEntry.Description = dbEntry.Description
			// 工事ステータスの更新
			fsEntry.Status = DetermineKoujiStatus(dbEntry.StartDate, dbEntry.EndDate)
			// 工事タグの更新
			fsEntry.Tags = dbEntry.Tags

			// Remove from map so we don't add it again
			delete(entryMap, fsEntry.ID)
		}
		// New project from file system - add it
		updatedEntries[entryCount] = fsEntry
		entryCount++
	}

	return s.Database.SaveToYAML(updatedEntries)
}

// checkFile はファイル名が規則に沿っているかを返す
// "新規工事.xlsx"は規則に沿っていいるがテンプレートのため"<工事フォルダー名>.xlsx"への更新が必要
// 他のファイルは未対応
func checkFile(file string) bool {
	// ファイル名が規則に沿っているかを返す
	return true
}
