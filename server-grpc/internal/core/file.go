package core

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// NormalizeAbsPath は絶対パスを最短パスに変換します。
// absPath は絶対パスです。
// '~' はホームディレクトリを展開して絶対パスに変換します。
func NormalizeAbsPath(absPath string) (string, error) {
	// ホームディレクトリに展開
	if strings.HasPrefix(absPath, "~/") {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		absPath = filepath.Join(usr.HomeDir, absPath[2:])
	}
	cleanPath := filepath.Clean(absPath)

	// シンボリックリンクを解決して絶対パスチェック
	resolvedPath, err := filepath.EvalSymlinks(cleanPath)
	if err == nil {
		if filepath.IsAbs(cleanPath) {
			// 絶対パスチェック
			return resolvedPath, nil
		}
		err = errors.New("絶対パスの条件を満たしていません。")
	}
	return "", err
}

// パスからファイル名またはフォルダー名を取得します
func GetBaseName(pathname string) string {
	basename := filepath.Base(pathname)
	if basename == "." || basename == "/" || basename == "\\" {
		return ""
	}
	return basename
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
