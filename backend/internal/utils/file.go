package utils

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// CleanAbsPath は絶対パスを最短パスに変換します。
// absPath は絶対パスです。
// '~' はホームディレクトリを展開して絶対パスに変換します。
func CleanAbsPath(absPath string) (string, error) {
	// ホームディレクトリに展開
	if strings.HasPrefix(absPath, "~/") {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		absPath = filepath.Join(usr.HomeDir, absPath[2:])
	}
	cleanPath := filepath.Clean(absPath)

	// 絶対パスチェック
	if filepath.IsAbs(cleanPath) {
		return cleanPath, nil
	}

	return "", errors.New("絶対パスではありません")
}

// IsExcel エクセルファイルかどうかをチェック
func FilenameIsExcel(filename string) bool {
	excelSuffix := []string{".xlsx", ".xls"}
	nameLower := strings.ToLower(filename)

	for _, suffix := range excelSuffix {
		if strings.HasSuffix(nameLower, suffix) {
			return true
		}
	}
	return false
}

// EntryIsExcelFile エクセルファイルかどうかをチェック
func EntryIsExcelFile(entry os.DirEntry) bool {
	return FilenameIsExcel(entry.Name())
}
