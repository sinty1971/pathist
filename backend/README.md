# Penguin Backend (gRPC)

Connect ベースの gRPC サーバーです。ファイル・会社・工事ドメインの情報を公開し、フロントエンドからは Connect-Web 経由で利用します。

## 主な機能

- `FileService` : ファイル／フォルダの一覧取得、基準パスの問い合わせ
- `CompanyService` : 会社データの取得・更新、カテゴリー一覧
- `KojiService` : 工事データの取得・更新、標準ファイルの更新

API の定義は `proto/penguin/v1/penguin.proto` にまとまっており、`just generate-grpc` コマンドでサーバー側とフロントエンド側のスタブを再生成できます。

## 実行方法

```bash
cd backend
go run cmd/grpc/main.go
```

デフォルトでは HTTP/2 over h2c で `http://localhost:9090` を待ち受けます。TLS を有効化する場合は:

```bash
go run cmd/grpc/main.go -enable-tls -cert=cert.pem -key=key.pem
```

証明書が未作成の場合は `just generate-cert` で自己署名証明書を生成できます。

## 動作確認

CLI で簡易確認を行いたい場合は、同梱の `cmd/fileclient` を利用できます。

```bash
go run cmd/fileclient/main.go \
  -base-url http://localhost:9090 \
  -json
```

レスポンスには基準パスとファイル一覧が JSON で表示されます。

## ディレクトリ構成

```
backend/
├── cmd/
│   ├── fileclient/   # gRPC を叩く CLI
│   └── grpc/         # サーバーエントリポイント
├── gen/              # プロト生成コード（go generate で更新）
├── internal/
│   ├── models/       # ドメインモデル
│   ├── rpc/          # Connect ハンドラー
│   └── services/     # ビジネスロジック
├── tools/            # go generate スクリプト
└── go.mod
```

## 参考コマンド

- `just backend` : gRPC サーバーを h2c で起動
- `just backend-tls` : TLS 有効で gRPC サーバーを起動
- `just generate-grpc` : Go サーバー側のスタブを再生成

詳細はリポジトリ直下の `justfile` を参照してください。
