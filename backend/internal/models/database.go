package models

// DatabaseFile はデータベースのインターフェースを定義します
type DatabaseFile interface {
	GetFilePath() string
}
