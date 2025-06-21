package models

// FileEntry はファイルまたはディレクトリを表す
// @Description ファイルまたはディレクトリの情報
type FileEntry struct {
	ID uint64 `json:"id" yaml:"id" example:"123456"`
	// Name of the file or folder
	Name string `json:"name" yaml:"name" example:"documents"`
	// Full path to the file or folder
	Path string `json:"path" yaml:"path" example:"/home/user/documents"`
	// Whether this item is a directory
	IsDirectory bool `json:"is_directory" yaml:"is_directory" example:"true"`
	// Size of the file in bytes
	Size int64 `json:"size" yaml:"size" example:"4096"`
	// Last modification time
	ModifiedTime Timestamp `json:"modified_time" yaml:"modified_time"`
}

// GetFileEntriesResponse はファイルエントリ一覧のレスポンスを表す
// @Description ファイルエントリ一覧を含むレスポンス
type GetFileEntriesResponse struct {
	// File entries
	FileEntries []FileEntry `json:"file_entries" yaml:"file_entries"`
	// Folder number of file entries
	FolderCount int `json:"folder_count" example:"10"`
	// File number of file entries
	FileCount int `json:"file_count" example:"10"`
}
