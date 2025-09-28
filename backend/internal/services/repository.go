package services

import (
	"penguin-backend/internal/models"
)

// RepositoryService は FileRepository のエイリアスで、後方互換性のために一時的に残します
// TODO: 将来的にはすべての参照を FileRepository に置き換える
type RepositoryService[T models.Persistable] struct {
	*models.Repository[T]
	databaseFilename string
}

// NewRepositoryService は FileRepository を作成します
// 後方互換性のために一時的に残します
func NewRepositoryService[T models.Persistable](databaseFilename string) *RepositoryService[T] {
	return &RepositoryService[T]{
		databaseFilename: databaseFilename,
	}
}

// DatabaseFilename は設定されたデータベースファイル名を返します。
func (rs *RepositoryService[T]) DatabaseFilename() string {
	return rs.databaseFilename
}

// Load はデータをロードします
func (rs *RepositoryService[T]) Load(entity T) (T, error) {
	// FileRepositoryを作成してロード
	repo := models.NewRepository[T](entity.GetPersistPath())
	return repo.Load(entity)
}

// Save はデータを保存します
func (rs *RepositoryService[T]) Save(entity T) error {
	// FileRepositoryを作成して保存
	repo := models.NewRepository[T](entity.GetPersistPath())
	return repo.Save(entity)
}
