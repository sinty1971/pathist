package services

import (
	"path/filepath"
	"penguin-backend/internal/models"
)

// DatabaseFileService は FileRepository のエイリアスで、後方互換性のために一時的に残します
// TODO: 将来的にはすべての参照を FileRepository に置き換える
type DatabaseFileService[T models.Persistable] struct {
	*models.FileRepository[T]
	databaseFilename string
}

// NewDatabaseFileService は FileRepository を作成します
// 後方互換性のために一時的に残します
func NewDatabaseFileService[T models.Persistable](databaseFilename string) *DatabaseFileService[T] {
	return &DatabaseFileService[T]{
		databaseFilename: databaseFilename,
	}
}

// Load はデータをロードします
func (dfs *DatabaseFileService[T]) Load(entity T) (T, error) {
	// エンティティのフォルダーパスとデータベースファイル名を結合
	filePath := filepath.Join(entity.GetFolderPath(), dfs.databaseFilename)
	
	// FileRepositoryを作成してロード
	repo := models.NewFileRepository[T](filePath)
	return repo.Load(entity)
}

// Save はデータを保存します
func (dfs *DatabaseFileService[T]) Save(entity T) error {
	// エンティティのフォルダーパスとデータベースファイル名を結合
	filePath := filepath.Join(entity.GetFolderPath(), dfs.databaseFilename)
	
	// FileRepositoryを作成して保存
	repo := models.NewFileRepository[T](filePath)
	return repo.Save(entity)
}