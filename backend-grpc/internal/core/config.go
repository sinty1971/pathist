package core

// DefaultStringMap はサービスのデフォルト設定を定義します。
// これらの値は、環境変数が設定されていない場合に使用されます。
var ConfigMap = map[string]string{
	"FileServiceTarget":          "C:/SyncFolder/SynologyDrive/豊田築炉",
	"CompanyServiceTarget":       "C:/SyncFolder/SynologyDrive/豊田築炉/1 会社",
	"CompanyPersistFilename":     "@company.yaml",
	"CompanyPollIntervalMillSec": "3000",
	"KojiServiceTarget":          "C:/SyncFolder/SynologyDrive/豊田築炉/2 工事",
	"KojiPersistFilename":        "@koji.yaml",
	"MemberPersistFilename":      "@member.yaml",
}

var WorkerConfigMap = map[string]int{
	"MinumWorkers":   2,
	"MaximumWorkers": 16,
	"CpuMultiplier":  2,
}
