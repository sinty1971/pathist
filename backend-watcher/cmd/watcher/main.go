//go:build windows
// +build windows

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"backend-watcher/core"
)

func main() {
	var dirFlag string
	flag.StringVar(&dirFlag, "dir", "", "監視するディレクトリのパス")
	flag.Parse()

	dir := dirFlag
	if dir == "" && flag.NArg() > 0 {
		dir = flag.Arg(0)
	}
	if dir == "" {
		log.Fatal("監視するディレクトリを -dir もしくは引数で指定してください")
	}

	events, errs, cleanup, err := core.WatcherOnWindows(dir)
	if err != nil {
		log.Fatalf("WatcherOnWindows の起動に失敗しました: %v", err)
	}
	defer func() {
		if err := cleanup(); err != nil {
			log.Printf("クリーンアップ中にエラー: %v", err)
		}
	}()

	log.Printf("監視を開始します: %s", dir)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case ev, ok := <-events:
			if !ok {
				log.Println("イベントチャネルが閉じられました。終了します。")
				return
			}
			logDetailedEvent(ev)
		case err := <-errs:
			if err != nil {
				log.Printf("[ERROR] %v", err)
			}
		case <-sigCh:
			log.Println("割り込みを受信しました。終了します。")
			return
		}
	}
}

// logDetailedEvent はイベントの詳細情報をログ出力する
func logDetailedEvent(ev core.FileEvent) {
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	eventType := getEventDescription(ev.Op)

	// ファイル情報を取得
	fileInfo, err := os.Stat(ev.Name)
	var fileType, sizeInfo string

	if err == nil {
		if fileInfo.IsDir() {
			fileType = "ディレクトリ"
			sizeInfo = ""
		} else {
			fileType = "ファイル"
			sizeInfo = fmt.Sprintf(" (サイズ: %s)", formatFileSize(fileInfo.Size()))
		}
	} else if os.IsNotExist(err) {
		// 削除されたファイルの場合
		fileType = "削除済み"
		sizeInfo = ""
	} else {
		fileType = "不明"
		sizeInfo = ""
	}

	fileName := filepath.Base(ev.Name)
	dirPath := filepath.Dir(ev.Name)

	log.Printf("[%s] %s: %s (%s)%s\n  パス: %s",
		timestamp, eventType, fileName, fileType, sizeInfo, dirPath)
}

// getEventDescription はイベントタイプの日本語説明を返す
func getEventDescription(op core.FileOp) string {
	switch op {
	case core.Create:
		return "作成"
	case core.Write:
		return "更新"
	case core.Remove:
		return "削除"
	case core.Rename:
		return "リネーム"
	case core.Chmod:
		return "権限変更"
	default:
		return op.String()
	}
}

// formatFileSize はファイルサイズを人間が読みやすい形式に変換する
func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
