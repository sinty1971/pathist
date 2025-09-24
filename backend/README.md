# Penguin Backend API

Go + Fiber製のフォルダー管理API

## 機能
- ファイルエントリー（フォルダー・ファイル）の一覧取得
- 工事プロジェクトの管理
- 日時文字列のパース機能
- CORS対応でフロントエンドとの連携
- Swagger APIドキュメント

## エンドポイント

### GET /api/file/entries
ファイルエントリー（フォルダー・ファイル）一覧を取得

**クエリパラメータ:**
- `path` (optional): 対象パス (デフォルト: `~/penguin`)

### GET /api/kouji/entries
工事一覧を取得

**クエリパラメータ:**
- `path` (optional): 対象パス (デフォルト: `~/penguin/豊田築炉/2-工事`)

### POST /api/kouji/entries/save
工事プロジェクト情報をYAMLファイルに保存

**レスポンス例:**
```json
{
  "folders": [
    {
      "id": "12345678",
      "name": "プロジェクトA",
      "path": "/home/user/penguin/プロジェクトA",
      "is_directory": true,
      "size": 4096,
      "modified_time": "2024-01-01T12:00:00Z"
    }
  ],
  "count": 1,
  "path": "/home/user/penguin"
}
```

## 実行方法

```bash
cd backend
go mod tidy
go run cmd/main.go
```

サーバーは http://localhost:8080 で起動します。

## API使用例

### 基本的な使用方法

1. **ファイルエントリー一覧の取得**:
   ```bash
   curl "http://localhost:8080/api/file/entries"
   ```

2. **工事一覧の取得**:
   ```bash
   curl "http://localhost:8080/api/kouji/entries"
   ```

3. **APIドキュメント (Huma)**:
   - `http://localhost:8080/docs` をブラウザで開く

### レスポンスの説明
- `folders`: ファイルエントリーの配列
- `count`: 取得された項目数
- `path`: 実際に読み取られたパス
- 各ファイルエントリー:
  - `id`: 一意のID（Unix inode）
  - `name`: ファイル/フォルダー名
  - `path`: フルパス
  - `is_directory`: ディレクトリかどうかのフラグ
  - `size`: ファイルサイズ（バイト）
  - `modified_time`: 最終更新時刻

## ディレクトリ構造

```
backend/
├── cmd/
│   └── main.go              # エントリーポイント
├── internal/
│   ├── handlers/            # HTTPハンドラー
│   ├── models/              # データモデル
│   └── services/            # ビジネスロジック
├── pkg/                     # 外部公開パッケージ
└── go.mod                   # Go modules
```
