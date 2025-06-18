# CLAUDE.md

このファイルは、このリポジトリでコードを扱う際のClaude Code (claude.ai/code) へのガイダンスを提供します。
レポート等は原則として日本語で行います。

## プロジェクト概要

フォルダー管理システムは、フォルダーの閲覧と管理、および工事プロジェクトの管理のためのWebインターフェースを提供します。構成は以下の通りです：
- **バックエンド**: Go (1.21) と Fiber v2 フレームワーク
- **フロントエンド**: React (19.1.0) と TypeScript、React Router v7 (Remix)

## 開発コマンド

### フロントエンド開発 (React Router v7)
```bash
cd frontend
npm install          # 依存関係をインストール
npm run dev          # React Router v7開発サーバーを起動 (http://localhost:5173)
npm run build        # React Router v7本番用にビルド
npm run start        # 本番サーバーを起動
npm run lint         # ESLintを実行
npm run preview      # 本番ビルドをプレビュー
npm run generate-api # OpenAPIからTypeScript型を生成
```

### バックエンド開発
```bash
cd backend
go mod tidy          # 依存関係をインストール/更新
go run cmd/main.go   # サーバーを起動 (http://localhost:8080)
```

### Justfileコマンド（推奨）
```bash
just backend         # バックエンドサーバーを起動
just frontend        # フロントエンドサーバーを起動
just generate-api    # バックエンドからOpenAPI仕様を生成
just generate-types  # フロントエンド用TypeScript型を生成
just generate-all    # API仕様と型を一括生成
just swagger         # Swagger UIをブラウザで開く
just help            # 利用可能なコマンド一覧
```

## プロジェクト構造

```
/home/shin/dev/claude/
├── schemas/                    # 共有API仕様（OpenAPI）
│   ├── openapi.yaml
│   └── openapi.json
├── backend/                    # Go + Fiber v2
│   ├── cmd/main.go
│   ├── internal/
│   │   ├── handlers/
│   │   ├── services/
│   │   └── models/
│   └── docs/                   # Swagger生成ファイル（一時的）
└── frontend/                   # React + React Router v7
    ├── app/
    │   ├── root.tsx
    │   ├── routes.ts
    │   ├── routes/
    │   ├── components/
    │   ├── api/                # 自動生成APIクライアント
    │   ├── types/
    │   └── styles/
    ├── package.json
    └── vite.config.ts
```

## アーキテクチャ

### バックエンド構造
- `cmd/main.go`: エントリーポイント、CORSを持つFiberサーバーをセットアップ
- `internal/handlers/`: HTTPリクエストハンドラー (filesystem.go, kouji.go, time.go)
- `internal/services/`: ビジネスロジック (filesystem.go, kouji.go)
- `internal/models/`: データモデル (fileentry.go, kouji.go, id.go, time.go, timestamp.go)

バックエンドは `http://localhost:8080/api` でREST APIを提供し、主要なエンドポイントは以下です：
- `GET /api/file-entries?path=<オプション-パス>` - フォルダーの内容を返す
- `GET /api/kouji-entries?path=<オプション-パス>` - 工事プロジェクトの一覧を返す
- `POST /api/kouji-entries/save` - 工事プロジェクト情報をYAMLファイルに保存

### フロントエンド構造 (React Router v7)
- `app/root.tsx`: ルートレイアウトコンポーネント
- `app/routes.ts`: ルート設定ファイル
- `app/routes/`: ファイルベースルーティング
  - `_layout.tsx`: メインレイアウト
  - `_layout._index.tsx`: ホームページ（フォルダー一覧）
  - `_layout.kouji.tsx`: 工事一覧ページ
- `app/components/`: UIコンポーネント (FileEntryGrid, FileEntryModal, KoujiProjectGrid, KoujiProjectPage)
- `app/api/`: 自動生成されたAPIクライアント
- `app/types/`: TypeScript型定義 (kouji.ts)
- `app/styles/`: CSS ファイル

### 主要な実装詳細

1. **デフォルトパス**: 
   - 一般フォルダー: `~/penguin` ディレクトリを標準で参照
   - 工事プロジェクト: `~/penguin/豊田築炉/2-工事` ディレクトリを標準で参照

2. **CORS**: バックエンドは `AllowOrigins: "*"` で全てのオリジンを許可します

3. **ファイル種別検出**: フロントエンドはファイル拡張子に基づいて異なるアイコンを表示します：
   - フォルダー: 📁
   - PDF: 📄
   - 画像 (jpg, jpeg, png, gif): 🖼️
   - 動画 (mp4, avi, mov): 🎬
   - 音声 (mp3, wav): 🎵
   - その他: 📎

4. **フォルダーリストの管理
   - Id: UnixシステムのInoのuint64を使用
   - 日時データ: 更新日
       - バックエンドのデータベース保存はフォーマット形式RFC3339Nanoの文字列保存

4. **工事一覧の管理**: 
   - フォルダー命名規則: `YYYY-MMDD 会社名 現場名` (例: `2025-0618 豊田築炉 名和工場`)
   - 工事ID: フォルダーのID+元請け会社名+現場名から一意ID生成
   - 工事データー永続化: `/home/<user>/penguin/豊田築炉/2-工事/.inside.yaml` ファイルで工事情報を保存
   - 日時データ: 工事開始日、工事完了日、フォルダー更新日
       - バックエンドのデータベース保存はフォーマット形式RFC3339Nanoの文字列保存
   - タイムゾーン: JST（ローカルタイム）で日時を保持

5. **APIレスポンス形式**: 
   - 一般フォルダー: name、path、size、isDirectory などのプロパティを持つ配列
   - 工事プロジェクト: id、company_name、location_name などの拡張プロパティを含む配列

6. **バックエンド内のデータソース定義**:
   - **FileSystem (fs)**: ファイルシステムから取得した情報
   - **Database (db)**: データベース（`.inside.yaml`ファイル）から取得した情報
     - 工事プロジェクトデータベースの保存場所: `~/penguin/豊田築炉/2-工事/.inside.yaml`
   - **Merge (mg)**: FileSystemとDatabaseのデータマージ処理
     - この処理は工事プロジェクト管理において重要な役割を持つ
     - 原則このデータをフロントエンドに提供する

## API仕様とコード生成

### OpenAPI仕様の管理
- **生成元**: バックエンドのGoコード内のSwaggerアノテーション
- **共有場所**: `schemas/openapi.{yaml,json}`
- **フロントエンド**: `app/api/`に自動生成されるTypeScriptクライアント

### コード生成ワークフロー
1. `just generate-api` - バックエンドからOpenAPI仕様を生成
2. `just generate-types` - フロントエンド用TypeScript型とクライアントを生成
3. `just generate-all` - 上記を一括実行

### Swagger UI
- URL: http://localhost:8080/swagger/index.html
- コマンド: `just swagger` でブラウザで開く