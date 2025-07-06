# CLAUDE.md

このファイルは、このリポジトリでコードを扱う際のClaude Code (claude.ai/code) へのガイダンスを提供します。
レポート等は原則として日本語で行います。

## システム概要

フォルダー管理システムは、フォルダーの閲覧と管理、および工事情報の管理のためのWebインターフェースを提供します。

- **バックエンド**: Go (1.21) と Fiber v2 フレームワーク
- **フロントエンド**: React (19.1.0) と TypeScript、React Router v7 (Remix)

## ドキュメント構成

システムのドキュメントは以下のように分割されています：

### 📖 詳細ドキュメント
- **[backend/CLAUDE.md](./backend/CLAUDE.md)** - バックエンド（Go + Fiber）の詳細情報
- **[frontend/CLAUDE.md](./frontend/CLAUDE.md)** - フロントエンド（React + React Router v7）の詳細情報

### 🔧 開発用
- **[justfile](./justfile)** - 開発コマンドの一覧
- **[schemas/](./schemas/)** - OpenAPI仕様（バックエンド ↔ フロントエンド間のAPI定義）

## クイックスタート

### 開発サーバー起動
```bash
# バックエンド
just backend

# フロントエンド  
just frontend

# API仕様生成（バックエンド変更時）
just generate-all
```

### ディレクトリ構造
```
/home/shin/dev/claude/
├── backend/           # Go + Fiber v2
│   └── CLAUDE.md      # バックエンド詳細ドキュメント
├── frontend/          # React + React Router v7
│   └── CLAUDE.md      # フロントエンド詳細ドキュメント
├── schemas/           # 共有API仕様（OpenAPI）
└── justfile           # 開発コマンド
```

## 主要機能

1. **ファイル管理** - ~/penguinディレクトリのファイル・フォルダ閲覧
2. **工事管理** - 工事情報の一覧・編集・管理
3. **ガントチャート** - 工事のタイムライン表示

## データ保存場所

- **一般ファイル**: `~/penguin`
- **工事データ**: `~/penguin/豊田築炉/2-工事`
- **工事詳細**: `~/penguin/豊田築炉/2-工事/.detail.yaml`