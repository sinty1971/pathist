# BACKEND.md

バックエンド（Go + Fiber v3）の詳細情報とガイダンス

## 技術スタック

- **言語**: Go 1.21+
- **フレームワーク**: Fiber v3
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

## ディレクトリ構造

```
backend/
├── cmd/
│   └── main.go                 # エントリーポイント、CORSを持つFiberサーバー
├── internal/
│   ├── handlers/               # HTTPリクエストハンドラー
│   │   ├── koji.go             # 工事関連API
│   │   └── company.go          # 会社関連API
│   ├── services/               # ビジネスロジック
│   │   ├── file.go             # ファイル操作
│   │   ├── koji.go             # 工事管理
│   │   └── business.go         # ビジネスデータ管理
│   └── models/                 # データモデル
│       ├── file.go             # ファイル関連モデル
│       ├── koji.go             # 工事関連モデル
│       ├── company.go          # 会社関連モデル
│       ├── id.go               # ID生成
│       └── timestamp.go        # タイムスタンプ処理
└── docs/                       # Swagger生成ファイル（一時的）
```

## API エンドポイント

バックエンドは `http://localhost:8080/api` でRESTful APIを提供します。

### ファイル管理
- `GET /api/files?path=<オプション-パス>` - ファイル・フォルダーの一覧を返す

### 工事管理
- `GET /api/kojies?filter=recent` - 工事一覧を返す（filter=recentで最近の工事のみ）
- `GET /api/kojies/{path}` - 指定パスの工事を取得
- `PUT /api/kojies` - 工事情報をYAMLファイルに保存
- `PUT /api/kojies/managed-files` - 管理ファイルの名前変更

### 会社管理
- `GET /api/companies` - 会社一覧を返す
- `GET /api/companies/{id}` - 指定IDの会社詳細を取得
- `PUT /api/companies` - 会社情報を更新

### Swagger UI
- URL: http://localhost:8080/swagger/index.html
- コマンド: `just swagger` でブラウザで開く
- **注意**: Fiber v3対応のため、カスタムHTML実装を使用（CDN経由でSwagger UI読み込み）
- OpenAPI仕様: `/swagger/swagger.json`, `/swagger/swagger.yaml`で直接アクセス可能

## データソース定義

### FileSystem (fs)
ファイルシステムから取得した情報
- ファイル/フォルダーの基本情報
- パス、サイズ、更新日時

### Database (db)
データベース（`.detail.yaml`ファイル）から取得した情報
- 工事データベースの保存場所: `~/penguin/豊田築炉/2-工事/.detail.yaml`
- 工事詳細情報（開始日、終了日、説明、タグ等）

### Merge (mg)
FileSystemとDatabaseのデータマージ処理
- この処理は工事管理において重要な役割を持つ
- 原則このデータをフロントエンドに提供する

## データ形式

### ファイル情報
- **ID**: UnixシステムのInoのuint64を使用
- **日時データ**: RFC3339Nano形式の文字列で保存
- **パス**: 相対パス形式で管理

### 工事情報
- **フォルダー命名規則**: `YYYY-MMDD 会社名 現場名` (例: `2025-0618 豊田築炉 名和工場`)
- **工事ID**: フォルダーのID+元請け会社名+現場名から一意ID生成
- **データ永続化**: `.detail.yaml`ファイルで工事情報を保存
- **タイムゾーン**: JST（ローカルタイム）で日時を保持

### 会社管理
- **フォルダー命名規則**: `[数字] 会社名` (例: `4 豊田築炉`)
- **業種マップ**: 数字から業種名への変換マップを使用
  - 0: 自社, 1: 下請会社, 2: 築炉会社, 3: 一人親方, 4: 元請け
  - 5: リース会社, 6: 販売会社, 7: 販売会社, 8: 求人会社, 9: その他
- **業種変換関数**:
  - `DetermineBusinessType(number)`: 数字 → 業種名
  - `GetBusinessTypeNumber(typeName)`: 業種名 → 数字
  - `GetAllBusinessTypes()`: 全業種マップを取得
- **会社ID**: 会社名から一意ID生成
- **データ永続化**: `.detail.yaml`ファイルで会社詳細情報を保存

## API仕様生成

### Swagger アノテーション記述ルール

**重要**: Swaggerアノテーションは`internal/handlers/`内のファイルにのみ記述し、`internal/routes/`には記述しない。

#### 記述場所
- ✅ **handlers/*.go**: ハンドラー関数の直前に記述
- ❌ **routes/*.go**: ルート定義のみ、Swaggerアノテーションは記述しない

#### 標準テンプレート
```go
// [ハンドラー名] godoc
// @Summary      [機能の簡潔な説明]
// @Description  [機能の詳細説明]
// @Tags         [APIグループ名]
// @Accept       json
// @Produce      json
// @Param        [パラメータ名] [in] [型] [必須] "[説明]"
// @Success      200 {object/array} [レスポンス型] "[成功時の説明]"
// @Failure      [ステータス] {object} map[string]string "[エラー時の説明]"
// @Router       /[パス] [メソッド]
func (h *Handler) FunctionName(c *fiber.Ctx) error {
    // 実装
}
```

#### タグの統一ルール
- **ファイル管理**: `ファイル管理`
- **工事管理**: `工事管理`
- **会社管理**: `会社管理`

#### 実例
```go
// GetCompanies godoc
// @Summary      会社一覧の取得
// @Description  会社フォルダーの一覧を取得します
// @Tags         会社管理
// @Accept       json
// @Produce      json
// @Success      200 {array} models.Company "正常なレスポンス"
// @Failure      500 {object} map[string]string "サーバーエラー"
// @Router       /company/list [get]
func (h *CompanyHandler) GetCompanies(c *fiber.Ctx) error {
    companies := h.CompanyService.GetCompanies()
    return c.JSON(companies)
}
```

### 生成ワークフロー
1. `internal/handlers/`内のSwaggerアノテーションから生成
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


## Prefork設定について

Preforkは、Fiberフレームワークの高性能機能の一つで、**マルチプロセスモード**を指します。

### Preforkの仕組み

通常、Goアプリケーションは1つのプロセスで動作しますが、Preforkを有効にすると：

1. **親プロセス**が複数の**子プロセス**を事前に起動
2. 各子プロセスが同じポートをリッスン（SO_REUSEPORTを使用）
3. OSが自動的にリクエストを各プロセスに分散

### 利点と欠点

#### 利点 ✅
- **高負荷対応**: CPUコアを最大限活用
- **パフォーマンス向上**: リクエストを並列処理
- **障害耐性**: 1つのプロセスがクラッシュしても他は継続

#### 欠点 ❌
- **メモリ使用量増加**: プロセス数分のメモリが必要
- **ステート共有不可**: プロセス間でメモリ共有できない
- **開発環境では不要**: デバッグが複雑になる

### 設定方法

```go
app := fiber.New(fiber.Config{
    Prefork: true,  // Preforkを有効化
})
```

### 現在の設定

現在の設定では`Prefork: Disabled`なので、シングルプロセスで動作しています。開発環境ではこれが推奨設定です。

サーバー起動時の出力例：
```
INFO Server started on:     http://127.0.0.1:8080 (bound on host 0.0.0.0 and port 8080)
INFO Total handlers count:  10
INFO Prefork:              Disabled
INFO PID:                  373347
INFO Total process count:  1
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