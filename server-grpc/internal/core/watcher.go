package core

import "github.com/fsnotify/fsnotify"

type Watcher struct {
	// pathistServiceFolder はこのサービスが管理する会社データのルートフォルダー
	pathistServiceFolder string

	// watcher は target のファイルシステム監視オブジェクト
	watcher *fsnotify.Watcher

	// watchedDirs は監視登録済みディレクトリの集合
	watchedDirs map[string]struct{}

	// maxDepth は監視するディレクトリの最大深度
	maxDepth int
}

type Watcherable interface {
	// watchTarget はファイルシステム監視を開始します
	watchTarget() error
}

func NewWatcher(target string, maxDepth int) (*Watcher, error) {
	// fsnotifyウォッチャーの作成
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &Watcher{
		pathistServiceFolder: target,
		watcher:              watcher,
		watchedDirs:          make(map[string]struct{}),
		maxDepth:             maxDepth,
	}, nil
}
