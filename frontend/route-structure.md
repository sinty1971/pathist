# React Router v7 ルート階層構造

## Mermaid図

```mermaid
graph TD
    A[/"/" - Root] --> B[_layout.tsx]
    B --> C[/"/" - Home<br/>_layout._index.tsx]
    B --> D["/files" - ファイル一覧<br/>_layout.files.tsx]
    B --> E["/projects" - 工程表<br/>_layout.projects.tsx]
    B --> F["/gantt" - ガントチャート<br/>_layout.gantt.tsx]
    
    style A fill:#e1f5fe
    style B fill:#f3e5f5
    style C fill:#e8f5e8
    style D fill:#e8f5e8
    style E fill:#e8f5e8
    style F fill:#e8f5e8
```

## ツリー構造（ASCII）

```
/
└── _layout.tsx (共通レイアウト)
    ├── / (ホーム) → _layout._index.tsx
    ├── /files (ファイル一覧) → _layout.files.tsx
    ├── /projects (工程表) → _layout.projects.tsx
    └── /gantt (ガントチャート) → _layout.gantt.tsx
```

## 詳細情報

| パス | ファイル | 説明 | コンポーネント |
|------|----------|------|----------------|
| `/` | `_layout._index.tsx` | ホームページ（機能紹介） | 機能カード表示 |
| `/files` | `_layout.files.tsx` | ファイル一覧 | TreeViewベースのファイルブラウザ |
| `/projects` | `_layout.projects.tsx` | 工程表 | 工事プロジェクト管理 |
| `/gantt` | `_layout.gantt.tsx` | ガントチャート | プロジェクトガントチャート |

## レイアウト構造

### 共通レイアウト (`_layout.tsx`)

- ナビゲーションバー
- ページタイトル表示
- FileInfoContext Provider
- 各ページの共通スタイル

### 各ページの特徴

1. **ホーム** (`/`)
   - 3つの機能カードを表示
   - Files, Projects, Gantt への導線

2. **ファイル一覧** (`/files`)
   - MUI TreeView使用
   - 階層ファイル表示
   - ファイル詳細モーダル

3. **工程表** (`/projects`)
   - 工事プロジェクト一覧
   - プロジェクト編集モーダル
   - YAML形式でのデータ保存

4. **ガントチャート** (`/gantt`)
   - プロジェクトの視覚的なタイムライン
   - 工事期間の可視化

## 自動生成について

このドキュメントは `routes.ts` の設定を基に作成されています。
ルート設定を変更した場合は、このドキュメントも更新する必要があります。
