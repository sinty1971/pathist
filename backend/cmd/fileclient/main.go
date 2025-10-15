package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	penguinv1 "penguin-backend/gen/penguin/v1"
	penguinv1connect "penguin-backend/gen/penguin/v1/penguinv1connect"

	"connectrpc.com/connect"
)

func main() {
	var (
		baseURL = flag.String("base-url", "http://localhost:9090", "FileService のベース URL (例: http://localhost:9090)")
		pathArg = flag.String("path", "", "一覧を取得する相対パス (未指定時はルート)")
		jsonOut = flag.Bool("json", false, "JSON 形式で出力します")
		timeout = flag.Duration("timeout", 5*time.Second, "RPC 呼び出しのタイムアウト")
	)
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	client := penguinv1connect.NewFileServiceClient(http.DefaultClient, *baseURL)

	baseResp, err := client.GetFileBasePath(ctx, connect.NewRequest(
		penguinv1.GetFileBasePathRequest_builder{}.Build(),
	))
	if err != nil {
		log.Fatalf("GetFileBasePath の呼び出しに失敗しました: %v", err)
	}

	req := connect.NewRequest(penguinv1.ListFilesRequest_builder{
		Path: *pathArg,
	}.Build())

	resp, err := client.ListFiles(ctx, req)
	if err != nil {
		log.Fatalf("ListFiles の呼び出しに失敗しました: %v", err)
	}

	if *jsonOut {
		output := struct {
			BasePath string                `json:"basePath"`
			Path     string                `json:"path"`
			Files    []*penguinv1.FileInfo `json:"files"`
		}{
			BasePath: baseResp.Msg.GetBasePath(),
			Path:     req.Msg.GetPath(),
			Files:    resp.Msg.GetFiles(),
		}

		data, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			log.Fatalf("レスポンスの JSON 変換に失敗しました: %v", err)
		}
		fmt.Println(string(data))
		return
	}

	// ターミナルで読みやすいように簡易フォーマットで出力する。
	fmt.Printf("BasePath: %s\n", baseResp.Msg.GetBasePath())
	fmt.Printf("Path: %s\n", req.Msg.GetPath())
	fmt.Println("IsDir\tSize\tModified\tTargetPath\tIdealPath")
	for _, file := range resp.Msg.GetFiles() {
		modified := ""
		if ts := file.GetModifiedTime(); ts != nil {
			modified = ts.AsTime().Format(time.RFC3339)
		}
		fmt.Printf("%t\t%d\t%s\t%s\t%s\n",
			file.GetIsDirectory(),
			file.GetSize(),
			modified,
			file.GetTargetPath(),
			file.GetIdealPath(),
		)
	}
}
