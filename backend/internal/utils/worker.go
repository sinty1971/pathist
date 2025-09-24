package utils

import (
	"runtime"
)

// DecideNumWorkers は要素数とシステムリソースに基づいて最適なワーカー数を決定する
// itemCount: 処理する要素数
// minWorkers: 最小ワーカー数（デフォルト: 1）
// maxWorkers: 最大ワーカー数（デフォルト: 16）
// cpuMultiplier: CPUコア数に対する倍率（デフォルト: 2）
func DecideNumWorkers(itemCount int, opts ...WorkerOption) int {
	// デフォルト設定
	config := &workerConfig{
		minWorkers:    1,
		maxWorkers:    16,
		cpuMultiplier: 2,
	}

	// オプションを適用
	for _, opt := range opts {
		opt(config)
	}

	// 要素数が0の場合は0を返す
	if itemCount <= 0 {
		return 0
	}

	// CPU数ベースでワーカー数を計算
	numCPU := runtime.NumCPU()
	idealWorkers := numCPU * config.cpuMultiplier

	// 最小値と最大値の範囲内に収める
	numWorkers := max(config.minWorkers, min(config.maxWorkers, idealWorkers))

	// 要素数が少ない場合はワーカー数を調整
	if itemCount < numWorkers {
		// 要素数の半分程度のワーカー数にする（最小1）
		numWorkers = max(1, itemCount/2)
	}

	return numWorkers
}

// workerConfig はワーカー数決定のための設定
type workerConfig struct {
	minWorkers    int
	maxWorkers    int
	cpuMultiplier int
}

// WorkerOption はワーカー数決定のオプション
type WorkerOption func(*workerConfig)

// WithMinWorkers は最小ワーカー数を設定する
func WithMinWorkers(min int) WorkerOption {
	return func(c *workerConfig) {
		if min > 0 {
			c.minWorkers = min
		}
	}
}

// WithMaxWorkers は最大ワーカー数を設定する
func WithMaxWorkers(max int) WorkerOption {
	return func(c *workerConfig) {
		if max > 0 {
			c.maxWorkers = max
		}
	}
}

// WithCPUMultiplier はCPUコア数に対する倍率を設定する
func WithCPUMultiplier(multiplier int) WorkerOption {
	return func(c *workerConfig) {
		if multiplier > 0 {
			c.cpuMultiplier = multiplier
		}
	}
}
