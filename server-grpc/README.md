# Penguin gRPC server

Connect ベースの gRPC サーバーです。ファイル・会社・工事ドメインの情報を公開し、フロントエンドからは Connect-Web 経由で利用します。

## 主な機能

- `FileService` : ファイル／フォルダの一覧取得、基準パスの問い合わせ
- `CompanyService` : 会社データの取得・更新、カテゴリー一覧
- `KojiService` : 工事データの取得・更新、標準ファイルの更新

API の定義は `proto/grpc/v1/penguin.proto` にまとまっており、`buf generate --path proto/grpc/v1/penguin.proto` または `just generate-grpc` コマンドでサーバー側とフロントエンド側のスタブを再生成できます。

## 実行方法

```bash
cd server-grpc
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

```terminal
server-grpc/
├── cmd/
│   ├── companyclient/ # 会社 API を叩く CLI
│   ├── fileclient/    # ファイル API を叩く CLI
│   └── grpc/          # サーバーエントリポイント
├── gen/               # プロト生成コード（buf generate で更新）
├── internal/
│   ├── core/          # 共通ユーティリティ（ID生成、永続化など）
│   ├── models/        # ドメインモデル
│   └── services/      # ビジネスロジック（gRPC ハンドラー）
├── tools/             # go generate スクリプト
└── go.mod
```

## 参考コマンド

- `just server-grpc` : gRPC サーバーを h2c で起動
- `just server-grpc-tls` : TLS 有効で gRPC サーバーを起動
- `just generate-grpc` : Go サーバー側のスタブを再生成

詳細はリポジトリ直下の `justfile` を参照してください。

## buf.build

- [インストール手順](https://buf.build/docs/installation) に従って `buf` CLI を導入してください。
- リポジトリルートで `buf generate --path proto/grpc/v1/penguin.proto` を実行すると、Go/Connect-Go/TypeScript スタブがそれぞれ `server-grpc/gen/` と `frontend/src/gen/` に出力されます。
