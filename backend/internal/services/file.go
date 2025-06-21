package services

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"
	"penguin-backend/internal/models"
	"strings"
	"syscall"
)

// FileService is a service for managing the file system
type FileService struct {
	// Root is the root directory of the file system
	Root string `json:"root" yaml:"root" example:"/home/<user>/penguin"`
}

// NewFileSystemService creates a new FileSystemService
func NewFileSystemService(root string) (*FileService, error) {
	// Expand ~ to home directory
	if strings.HasPrefix(root, "~/") {
		usr, err := user.Current()
		if err != nil {
			return nil, err
		}
		root = filepath.Join(usr.HomeDir, root[2:])
	}

	// Get absolute path
	absPath, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	return &FileService{
		Root: absPath,
	}, nil
}

// GetAbsPath 相対パスをRootパスを追加したパスを返す
// 絶対パスの場合はエラーを返す
func (s *FileService) GetFullpath(joinPath ...string) (string, error) {
	if len(joinPath) == 0 {
		return s.Root, nil
	}
	var fullpath string

	// Expand ~ to home directory
	if strings.HasPrefix(joinPath[0], "~/") {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		joinPath[0] = filepath.Join(usr.HomeDir, joinPath[0][2:])
		fullpath = filepath.Join(joinPath...)
		return fullpath, nil
	} else if filepath.IsAbs(joinPath[0]) {
		// 絶対パスはエラーを返す
		return "", errors.New("absolute path is not allowed")
	}

	// 相対パスはRootパスを追加したパスを返す
	joinedPath := filepath.Join(joinPath...)
	fullpath = filepath.Join(s.Root, joinedPath)

	return fullpath, nil
}

// GetFileEntries gets the file entries from the file system
func (s *FileService) GetFileEntries(joinPath ...string) (*models.GetFileEntriesResponse, error) {
	fullpath, err := s.GetFullpath(joinPath...)
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(fullpath)
	if err != nil {
		return nil, err
	}

	fileEntries := make([]models.FileEntry, len(entries))
	var folderCount, fileCount int
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			// ファイル情報の取得に失敗した場合はスキップ
			continue
		}
		stat := info.Sys().(*syscall.Stat_t)

		// For symlinks, check the target's type
		entryFullpath := filepath.Join(fullpath, entry.Name())

		// Check if entry is a directory
		var isDirectory bool

		// If it's a symlink, check what it points to
		if info.Mode()&os.ModeSymlink != 0 {
			fileinfo, err := os.Stat(entryFullpath) // Follow the symlink
			if err == nil {
				isDirectory = fileinfo.IsDir()
			}
		} else {
			isDirectory = entry.IsDir()
		}

		// Add the file entry to the list
		fileEntries[folderCount+fileCount] = models.FileEntry{
			ID:           stat.Ino,
			Name:         entry.Name(),
			Path:         entryFullpath,
			IsDirectory:  isDirectory,
			Size:         info.Size(),
			ModifiedTime: models.NewTimestamp(info.ModTime()),
		}

		// Count the number of folders and files
		if isDirectory {
			folderCount++
		} else {
			fileCount++
		}

	}

	return &models.GetFileEntriesResponse{
		FileEntries: fileEntries[:folderCount+fileCount],
		FolderCount: folderCount,
		FileCount:   fileCount,
	}, nil
}
