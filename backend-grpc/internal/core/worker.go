package core

import (
	"runtime"
)

// DecideNumWorkers は要素数とシステムリソースに基づいて最適なワーカー数を決定する
// itemCount: 処理する要素数
func DecideNumWorkers(itemCount int) int {

	// 要素数が0の場合は0を返す
	if itemCount <= 0 {
		return 0
	}

	// ワーカーの環境変数を取得
	cpuMultiplier := WorkerConfigMap["CpuMultiplier"]
	minWorkers := WorkerConfigMap["MinumWorkers"]
	maxWorkers := WorkerConfigMap["MaximumWorkers"]

	// CPU数ベースでワーカー数を計算
	numCPU := runtime.NumCPU()
	idealWorkers := numCPU * cpuMultiplier

	// 最小値と最大値の範囲内に収める
	numWorkers := max(minWorkers, min(maxWorkers, idealWorkers))

	// 要素数が少ない場合はワーカー数を調整
	if itemCount < numWorkers {
		// 要素数の半分程度のワーカー数にする（最小1）
		numWorkers = max(1, itemCount/2)
	}

	return numWorkers
}
