package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	grpcv1 "backend-grpc/gen/grpc/v1"
	"backend-grpc/gen/grpc/v1/grpcv1connect"
)

func main() {
	var (
		baseURL = flag.String("base-url", "http://localhost:9090", "FileService のベース URL (例: http://localhost:9090)")
		target  = flag.String("target", "", "一覧を取得する相対パス (未指定時はルート)")
		jsonOut = flag.Bool("json", false, "JSON 形式で出力します")
		timeout = flag.Duration("timeout", 5*time.Second, "RPC 呼び出しのタイムアウト")
	)
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	client := grpcv1connect.NewFileServiceClient(http.DefaultClient, *baseURL)

	resFileBasePath, err := client.GetFileTarget(
		ctx,
		grpcv1.GetFileTargetRequest_builder{}.Build(),
	)
	if err != nil {
		log.Fatalf("GetFileTarget の呼び出しに失敗しました: %v", err)
	}

	req := grpcv1.GetFilesRequest_builder{
		Target: *target,
	}.Build()

	resFiles, err := client.GetFiles(ctx, req)
	if err != nil {
		log.Fatalf("ListFiles の呼び出しに失敗しました: %v", err)
	}

	if *jsonOut {
		output := struct {
			ManagedFolder string         `json:"managedFolder"`
			Target        string         `json:"target"`
			Files         []*grpcv1.File `json:"files"`
		}{
			ManagedFolder: resFileBasePath.GetTarget(),
			Target:        req.GetTarget(),
			Files:         resFiles.GetFiles(),
		}

		data, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			log.Fatalf("レスポンスの JSON 変換に失敗しました: %v", err)
		}
		fmt.Println(string(data))
		return
	}

	// ターミナルで読みやすいように簡易フォーマットで出力する。
	fmt.Printf("ManagedFolder: %s\n", resFileBasePath.GetTarget())
	fmt.Printf("Target: %s\n", req.GetTarget())
	fmt.Println("IsDir\tSize\tModified\tTargetPath\tIdealPath")
	for _, file := range resFiles.GetFiles() {
		fmt.Printf("%d\t%s\t%s\t\n",
			file.GetSize(),
			file.GetModifiedTime().AsTime().Format(time.RFC3339Nano),
			file.GetTarget(),
		)
	}
}
