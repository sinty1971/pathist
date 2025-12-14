//go:build windows
// +build windows

package core

import (
	"io/fs"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// FileOp はファイル操作の種類を表す
type FileOp int

const (
	Create FileOp = iota
	Write
	Remove
	Rename
	Chmod
)

func (op FileOp) String() string {
	switch op {
	case Create:
		return "CREATE"
	case Write:
		return "WRITE"
	case Remove:
		return "REMOVE"
	case Rename:
		return "RENAME"
	case Chmod:
		return "CHMOD"
	default:
		return "UNKNOWN"
	}
}

// FileEvent はファイルシステムイベントの詳細情報を持つ
type FileEvent struct {
	Name string
	Op   FileOp
}

func convertEvent(ev fsnotify.Event) FileEvent {
	var op FileOp
	switch {
	case ev.Op&fsnotify.Create == fsnotify.Create:
		op = Create
	case ev.Op&fsnotify.Write == fsnotify.Write:
		op = Write
	case ev.Op&fsnotify.Remove == fsnotify.Remove:
		op = Remove
	case ev.Op&fsnotify.Rename == fsnotify.Rename:
		op = Rename
	case ev.Op&fsnotify.Chmod == fsnotify.Chmod:
		op = Chmod
	default:
		op = Write
	}
	return FileEvent{Name: ev.Name, Op: op}
}

// WatcherOnWindows は Windows 上のディレクトリを監視するシンプルなヘルパー。
// 呼び出し側は返却されるイベントチャネルを受信し、終了時に cleanup を必ず実行してください。
func WatcherOnWindows(dir string) (<-chan FileEvent, <-chan error, func() error, error) {
	watchDir := filepath.Clean(dir)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, nil, nil, err
	}

	eventCh := make(chan FileEvent, 8)
	errCh := make(chan error, 1)

	// イベント転送用の goroutine。軽量なため caller 側で待つ必要はない。
	go func() {
		defer close(eventCh)
		defer close(errCh)
		debouncer := newEventDebouncer(eventCh, time.Second)
		defer debouncer.Close()
		for {
			select {
			case ev, ok := <-watcher.Events:
				if !ok {
					debouncer.FlushAll()
					return
				}
				if strings.Contains(ev.Name, ".SynologyWorkingDirectory") {
					continue
				}
				debouncer.Push(convertEvent(ev))
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				errCh <- err
			}
		}
	}()

	if err := addDirWithDepth(watcher, watchDir, 2); err != nil {
		_ = watcher.Close()
		return nil, nil, nil, err
	}

	cleanup := func() error {
		return watcher.Close()
	}

	return eventCh, errCh, cleanup, nil
}

// addDirWithDepth は root から depth 階層分までのディレクトリを watcher に登録する。
func addDirWithDepth(w *fsnotify.Watcher, root string, depth int) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}

		if strings.Contains(path, ".SynologyWorkingDirectory") {
			return filepath.SkipDir
		}

		level := 0
		if path != root {
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			rel = filepath.ToSlash(rel)
			level = strings.Count(rel, "/") + 1
		}

		if level > depth {
			return filepath.SkipDir
		}

		if err := w.Add(path); err != nil {
			return err
		}

		return nil
	})
}

type eventDebouncer struct {
	window  time.Duration
	out     chan<- FileEvent
	mu      sync.Mutex
	pending map[string]debounceEntry
	ticker  *time.Ticker
	done    chan struct{}
	wg      sync.WaitGroup
}

type debounceEntry struct {
	event     FileEvent
	updatedAt time.Time
}

func newEventDebouncer(out chan<- FileEvent, window time.Duration) *eventDebouncer {
	d := &eventDebouncer{
		window:  window,
		out:     out,
		pending: make(map[string]debounceEntry),
		ticker:  time.NewTicker(window / 2),
		done:    make(chan struct{}),
	}
	d.wg.Add(1)
	go d.run()
	return d
}

func (d *eventDebouncer) run() {
	defer d.wg.Done()
	for {
		select {
		case <-d.ticker.C:
			d.flush(false)
		case <-d.done:
			d.ticker.Stop()
			d.flush(true)
			return
		}
	}
}

func (d *eventDebouncer) Push(ev FileEvent) {
	name := ev.Name
	d.mu.Lock()
	d.pending[name] = debounceEntry{event: ev, updatedAt: time.Now()}
	d.mu.Unlock()
}

func (d *eventDebouncer) FlushAll() {
	d.flush(true)
}

func (d *eventDebouncer) Close() {
	close(d.done)
	d.wg.Wait()
}

func (d *eventDebouncer) flush(all bool) {
	now := time.Now()
	var ready []FileEvent

	d.mu.Lock()
	for key, entry := range d.pending {
		if all || now.Sub(entry.updatedAt) >= d.window {
			ready = append(ready, entry.event)
			delete(d.pending, key)
		}
	}
	d.mu.Unlock()

	for _, ev := range ready {
		d.out <- ev
	}
}
