# フォルダー管理システム - フロントエンド

React + TypeScript + React Router v7を使用したフォルダー管理システムのフロントエンドアプリケーションです。

## 技術スタック

- **React** 19.1.0
- **TypeScript** 5.6.3
- **React Router v7** (旧 Remix)
- **Vite** 6.0.7
- **@hey-api/openapi-ts** - OpenAPIからTypeScript型とAPIクライアントを自動生成

## プロジェクト構成

```
frontend/
├── app/
│   ├── root.tsx              # ルートコンポーネント
│   ├── routes.ts             # ルート設定
│   ├── routes/               # ページコンポーネント
│   │   ├── _layout.tsx       # メインレイアウト
│   │   ├── _layout._index.tsx # ホーム（ファイル一覧）
│   │   ├── _layout.kouji.tsx  # 工事一覧
│   │   └── _layout.gantt.tsx  # ガントチャート
│   ├── components/           # UIコンポーネント
│   │   ├── FileInfoGrid.tsx  # ファイル一覧グリッド
│   │   ├── FileInfoModal.tsx # ファイル詳細モーダル
│   │   ├── ProjectGrid.tsx   # 工事一覧
│   │   ├── ProjectPage.tsx   # 工事プロジェクトページ
│   │   ├── ProjectEditModal.tsx # 工事編集モーダル
│   │   ├── ProjectGanttChart.tsx # ガントチャート（詳細）
│   │   └── Navigation.tsx    # ナビゲーションバー
│   ├── api/                  # 自動生成されたAPIクライアント
│   │   ├── client/           # クライアント実装
│   │   ├── core/             # コア機能
│   │   ├── client.gen.ts     # クライアント設定
│   │   ├── sdk.gen.ts        # SDK関数
│   │   └── types.gen.ts      # TypeScript型定義
│   ├── services/             # APIサービス層
│   │   └── api.ts            # APIサービスラッパー
│   ├── styles/               # CSSファイル
│   │   ├── App.css           # アプリケーション全体のスタイル
│   │   ├── gantt.css         # ガントチャート用スタイル
│   │   └── modal.css         # モーダル用スタイル
│   ├── types/                # 型定義
│   │   └── css.d.ts          # CSS モジュール型定義
│   └── utils/                # ユーティリティ関数
│       ├── date.ts           # 日付処理
│       └── timestamp.ts      # タイムスタンプ変換
├── public/
│   └── favicon.ico           # アプリケーションアイコン
├── package.json
├── vite.config.ts
└── tsconfig.json
```

## 開発コマンド

### 依存関係のインストール
```bash
npm install
```

### 開発サーバーの起動
```bash
npm run dev
```
開発サーバーは http://localhost:5173 で起動します。

### 本番用ビルド
```bash
npm run build
```

### 本番サーバーの起動
```bash
npm run start
```

### その他のコマンド
```bash
npm run lint          # ESLintを実行
npm run preview       # 本番ビルドをプレビュー
npm run generate-api  # OpenAPIからTypeScript型を生成
```

## API連携

### 自動生成
バックエンドのOpenAPI仕様からTypeScript型とAPIクライアントを自動生成しています：

```bash
# Justfileを使用する場合（推奨）
just generate-types

# npm scriptを使用する場合
npm run generate-api
```

### APIクライアントの使用例
```typescript
import { getFileFileinfos, getProjectRecent } from '../api/sdk.gen';

// ファイル一覧の取得
const response = await getFileFileinfos({ 
  query: { path: '/path/to/folder' }
});

// 工事一覧の取得
const projects = await getProjectRecent();
```

## 主な機能

### ファイル管理
- ファイル・フォルダーの一覧表示
- ファイル種別に応じたアイコン表示
- ファイル詳細情報の表示

### 工事プロジェクト管理
- 工事プロジェクトの一覧表示
- 工事情報の編集（会社名、現場名、期間、タグ等）
- ガントチャートによる工期の視覚化
- YAMLファイルへの情報永続化

### タイムスタンプ処理
- バックエンドの`ModelsTimestamp`型（RFC3339Nano形式）の適切な処理
- JSTタイムゾーンでの日時表示・編集
- 柔軟な日時フォーマット解析

## 開発上の注意点

### TypeScript型
- API型は`app/api/types.gen.ts`から自動生成されます
- 手動で型を編集しないでください（再生成時に上書きされます）

### 日時処理
- バックエンドは`ModelsTimestamp`型（オブジェクト形式）を使用
- フロントエンドでは`utils/timestamp.ts`の関数で変換処理を行う
- 日付入力は常にJSTで処理される

### スタイリング
- CSSモジュールは使用せず、通常のCSSファイルをインポート
- `css.d.ts`でCSSインポートの型定義を提供

## トラブルシューティング

### favicon.icoエラー
`public/favicon.ico`が存在することを確認してください。

### TypeScriptエラー
API型が最新でない場合は、以下を実行してください：
```bash
just generate-all  # またはnpm run generate-api
```

### 開発サーバーが起動しない
ポート5173が使用されていないか確認してください。