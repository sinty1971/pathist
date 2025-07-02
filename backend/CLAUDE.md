# BACKEND.md

バックエンド（Go + Fiber v2）の詳細情報とガイダンス

## 技術スタック

- **言語**: Go 1.21+
- **フレームワーク**: Fiber v2
- **API仕様**: Swagger 2.0 (OpenAPI 2.0)
- **CORS**: 全オリジン許可

## 開発コマンド

```bash
cd backend
go mod tidy          # 依存関係をインストール/更新
go run cmd/main.go   # サーバーを起動 (http://localhost:8080)
```

### Justfileコマンド（推奨）
```bash
just backend         # バックエンドサーバーを起動
just generate-api    # バックエンドからOpenAPI仕様を生成
just swagger         # Swagger UIをブラウザで開く
just stop-backend    # バックエンドサーバーを停止
```

## プロジェクト構造

```
backend/
├── cmd/
│   └── main.go                 # エントリーポイント、CORSを持つFiberサーバー
├── internal/
│   ├── handlers/               # HTTPリクエストハンドラー
│   │   ├── project.go          # プロジェクト関連API
│   │   └── time.go             # 時刻関連API
│   ├── services/               # ビジネスロジック
│   │   ├── file.go             # ファイル操作
│   │   ├── project.go          # プロジェクト管理
│   │   └── yaml.go             # YAML操作
│   └── models/                 # データモデル
│       ├── file.go             # ファイル関連モデル
│       ├── project.go          # プロジェクト関連モデル
│       ├── id.go               # ID生成
│       ├── time.go             # 時刻関連モデル
│       └── timestamp.go        # タイムスタンプ処理
└── docs/                       # Swagger生成ファイル（一時的）
```

## API エンドポイント

バックエンドは `http://localhost:8080/api` でREST APIを提供します。

### ファイル管理
- `GET /api/file/fileinfos?path=<オプション-パス>` - フォルダーの内容を返す

### プロジェクト管理
- `GET /api/project/recent` - 工事プロジェクトの一覧を返す
- `POST /api/project/update` - 工事プロジェクト情報をYAMLファイルに保存
- `POST /api/project/rename-managed-file` - 管理ファイルの名前変更

### Swagger UI
- URL: http://localhost:8080/swagger/index.html
- コマンド: `just swagger` でブラウザで開く

## データソース定義

### FileSystem (fs)
ファイルシステムから取得した情報
- ファイル/フォルダーの基本情報
- パス、サイズ、更新日時

### Database (db)
データベース（`.detail.yaml`ファイル）から取得した情報
- 工事プロジェクトデータベースの保存場所: `~/penguin/豊田築炉/2-工事/.detail.yaml`
- プロジェクト詳細情報（開始日、終了日、説明、タグ等）

### Merge (mg)
FileSystemとDatabaseのデータマージ処理
- この処理は工事プロジェクト管理において重要な役割を持つ
- 原則このデータをフロントエンドに提供する

## データ形式

### ファイル情報
- **ID**: UnixシステムのInoのuint64を使用
- **日時データ**: RFC3339Nano形式の文字列で保存
- **パス**: 相対パス形式で管理

### 工事プロジェクト
- **フォルダー命名規則**: `YYYY-MMDD 会社名 現場名` (例: `2025-0618 豊田築炉 名和工場`)
- **工事ID**: フォルダーのID+元請け会社名+現場名から一意ID生成
- **データ永続化**: `.detail.yaml`ファイルで工事情報を保存
- **タイムゾーン**: JST（ローカルタイム）で日時を保持

## API仕様生成

### Swagger アノテーション
```go
// @Summary プロジェクト一覧取得
// @Description 工事プロジェクトの一覧を取得
// @Tags project
// @Accept json
// @Produce json
// @Success 200 {array} models.Project
// @Router /project/recent [get]
```

### 生成ワークフロー
1. Goコード内のSwaggerアノテーションから生成
2. `just generate-api` でOpenAPI仕様を生成
3. `schemas/openapi.{yaml,json}` に出力
4. フロントエンドがこれを使ってTypeScript型を生成

## OpenAPI仕様について

- **現在**: Swagger 2.0 (OpenAPI 2.0)を使用
- **生成ツール**: swaggo/swag
- **変換**: OpenAPI 3.0への変換スクリプトあり (`just convert-openapi3`)
- **TODO**: go-swagger3がリリースされたら、OpenAPI 3.0ネイティブサポートへの移行を検討

## CORS設定

```go
app.Use(cors.New(cors.Config{
    AllowOrigins: "*",
    AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
    AllowHeaders: "*",
}))
```

## 開発時の注意事項

### ポート使用
- **バックエンド**: http://localhost:8080
- **フロントエンド**: http://localhost:5173 (開発時)

### データベース
- ファイルベースのYAML形式
- バックアップ機能なし（ファイルシステム依存）
- 同時書き込み制御なし

### セキュリティ
- 認証・認可機能なし
- ローカル開発環境想定
- ファイルシステムへの直接アクセス