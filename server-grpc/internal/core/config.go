package core

import (
	"os"
	"strings"
)

// DefaultStringMap はサービスのデフォルト設定を定義します。
// これらの値は、環境変数が設定されていない場合に使用されます。
var ConfigMap = map[string]string{
	"FileServiceTarget":          "{ROOT}",
	"CompanyServiceFolder":       "{ROOT}/1 会社",
	"CompanyPersistFilename":     "@company.yaml",
	"CompanyPollIntervalMillSec": "3000",
	"KojiServiceFolder":          "{ROOT}/2 工事",
	"KojiPersistFilename":        "@koji.yaml",
	"MemberPersistFilename":      "@member.yaml",
}

var WorkerConfigMap = map[string]int{
	"MinumWorkers":   2,
	"MaximumWorkers": 16,
	"CpuMultiplier":  2,
}

func init() {
	// PCのホスト名の取得
	hostname, err := os.Hostname()
	if err != nil {
		return
	}
	// ホスト名に基づく設定の上書き
	root := ""
	switch hostname {
	case "DESKTOP-HHR7FT6":
		root = "C:/SyncFolder/SynologyDrive/豊田築炉"
	case "SINTY-OMEN":
		root = "O:/"
	}
	for key, value := range ConfigMap {
		ConfigMap[key] = strings.ReplaceAll(value, "{ROOT}", root)
	}
}
