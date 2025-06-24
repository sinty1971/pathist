package utils

import (
	"os"
	"strings"
)

// EntryIsExcelFile エクセルファイルかどうかをチェック
func EntryIsExcelFile(entry os.DirEntry) bool {
	// エクセルサフィックスの定義
	excelSuffix := []string{".xlsx", ".xls"}

	// エクセルサフィックスをチェック
	for _, suffix := range excelSuffix {
		if strings.HasSuffix(entry.Name(), suffix) {
			return true
		}
	}
	return false
}
