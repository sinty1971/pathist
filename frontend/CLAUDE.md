# FRONTEND.md

フロントエンド（React + React Router v7）の詳細情報とガイダンス

## 技術スタック

- **言語**: TypeScript
- **フレームワーク**: React 19.1.0
- **ルーティング**: React Router v7 (Remix)
- **UIライブラリ**: 
  - Material-UI (MUI) - TreeView, DatePicker, Icons
  - daisyUI (Tailwind CSS ベース)
- **スタイリング**: Tailwind CSS + カスタムCSS
- **ビルドツール**: Vite
- **API**: 自動生成されたTypeScriptクライアント

## 開発コマンド

```bash
cd frontend
npm install          # 依存関係をインストール
npm run dev          # React Router v7開発サーバーを起動 (http://localhost:5173)
npm run build        # React Router v7本番用にビルド
npm run start        # 本番サーバーを起動
npm run lint         # ESLintを実行
npm run preview      # 本番ビルドをプレビュー
npm run generate-api # OpenAPIからTypeScript型を生成
npm run generate-routes # ルート構造図を生成
```

### Justfileコマンド（推奨）
```bash
just frontend        # フロントエンドサーバーを起動
just generate-types  # フロントエンド用TypeScript型を生成
just generate-routes # React Router v7ルート構造図を生成
just frontend-build  # 本番用ビルド
just frontend-lint   # ESLintを実行
just stop-frontend   # フロントエンドサーバーを停止
```

## プロジェクト構造

```
frontend/
├── app/
│   ├── root.tsx                    # ルートレイアウトコンポーネント
│   ├── routes.ts                   # ルート設定ファイル
│   ├── routes/                     # ファイルベースルーティング
│   │   ├── _layout.tsx             # メインレイアウト
│   │   ├── _layout._index.tsx      # ホームページ（機能紹介）
│   │   ├── _layout.files.tsx       # ファイル一覧ページ
│   │   ├── _layout.projects.tsx    # 工程表ページ
│   │   └── _layout.gantt.tsx       # ガントチャートページ（/projects/gantt）
│   ├── components/                 # UIコンポーネント
│   │   ├── Files.tsx               # TreeViewベースのファイルブラウザ
│   │   ├── FileDetailModal.tsx     # ファイル詳細モーダル
│   │   ├── Projects.tsx            # 工程表一覧
│   │   ├── ProjectDetailModal.tsx  # プロジェクト詳細編集モーダル
│   │   ├── ProjectGanttChart.tsx   # ガントチャート表示
│   │   ├── CalendarPicker.tsx      # MUI DatePickerベースのカレンダー
│   │   └── Navigation.tsx          # ナビゲーション
│   ├── contexts/                   # React Context
│   │   ├── FileInfoContext.tsx     # ファイル情報共有
│   │   └── ProjectContext.tsx      # プロジェクト情報共有
│   ├── api/                        # 自動生成APIクライアント
│   │   ├── client/                 # APIクライアント実装
│   │   ├── client.gen.ts           # 自動生成されたクライアント設定
│   │   ├── sdk.gen.ts              # 自動生成されたSDK関数
│   │   └── types.gen.ts            # 自動生成されたTypeScript型定義
│   ├── styles/                     # CSS ファイル
│   │   ├── App.css                 # メインスタイル
│   │   ├── gantt.css               # ガントチャート専用
│   │   └── modal.css               # モーダル専用
│   ├── utils/                      # ユーティリティ関数
│   │   ├── date.ts                 # 日付処理
│   │   └── timestamp.ts            # タイムスタンプ処理
│   └── services/                   # APIサービス
│       └── api.ts                  # API呼び出しサービス
├── scripts/                        # 開発スクリプト
│   └── generate-route-diagram.js   # ルート構造図生成
├── package.json
├── vite.config.ts
├── openapi-ts.config.ts            # API生成設定
└── public/
    └── favicon.ico                 # ファイル管理アプリ用アイコン
```

## ルーティング構成

### 現在のルート構造
```
/ (ホーム)
├── /files (ファイル一覧)
├── /projects (工程表)
└── /projects/gantt (ガントチャート)
```

### ルーティング変更時の注意事項

**重要**: ルーティング設定を変更した場合は、以下のファイルも同時に更新する必要があります：

1. **`app/routes.ts`**: メインのルート設定
2. **`app/routes/_layout._index.tsx`**: ホームページのリンク先を更新
3. **`app/components/Navigation.tsx`**: ナビゲーションのリンクとページタイトルを更新
4. **ルート構造図**: `just generate-routes` を実行して最新の図を生成

## コンポーネント設計

### ファイル管理 (Files.tsx)
- **MUI TreeView**: 階層ファイル表示
- **動的読み込み**: ディレクトリ展開時に子要素を読み込み
- **展開状態保持**: React.useCallbackで最適化
- **アイコン表示**: ファイル拡張子に基づく自動判定
- **検索・ナビゲーション**: パス入力、ブレッドクラム、ホーム/リフレッシュボタン

### 工程表管理 (Projects.tsx)
- **プロジェクト一覧**: 開始日順ソート（新しい順）
- **編集機能**: 行クリックでProjectDetailModalを表示
- **ガントチャートリンク**: ヘッダーにアクセスボタン
- **ステータス表示**: 完了/進行中/予定の色分け
- **ファイル名変更**: 管理ファイルのリネーム機能

### プロジェクト詳細 (ProjectDetailModal.tsx)
- **MUI Dialog**: フルスクリーンモーダル
- **CalendarPicker**: 日本語対応の日付選択
- **水平レイアウト**: 全フィールドでタイトルと入力欄を水平配置
- **ファイル名関連フィールド**: 青い背景と枠線で強調表示
- **推奨ファイル名生成**: 開始日・会社名・現場名変更時に自動生成
- **段階的更新**: ファイル名影響フィールド変更時は専用ボタンで更新

### ガントチャート (ProjectGanttChart.tsx)
- **カスタム実装**: Canvas/SVGベースの時系列表示
- **プロジェクト表示**: 期間バーとマイルストーン
- **インタラクション**: クリックで詳細編集
- **レスポンシブ**: 画面サイズに応じた表示調整

## データ管理

### API通信
- **自動生成クライアント**: OpenAPI仕様からTypeScript型とSDK関数を生成
- **エラーハンドリング**: try-catchとユーザーフレンドリーなエラー表示
- **ローディング状態**: CircularProgressでローディング表示

### 状態管理
- **React Context**: FileInfoContext, ProjectContext
- **useState**: コンポーネントローカル状態
- **useEffect**: 副作用とライフサイクル管理
- **useCallback**: パフォーマンス最適化

### データ形式
- **日時**: ModelsTimestamp型（オブジェクト形式）で扱う
- **ファイル情報**: name, path, size, isDirectory等のプロパティ
- **プロジェクト**: id, company_name, location_name等の拡張プロパティ

## UI/UX 設計

### デザインシステム
- **カラーパレット**: 白ベース + ブランドカラー（#667eea）
- **タイポグラフィ**: Inter フォント
- **アイコン**: 絵文字 + Material Icons
- **レイアウト**: Flexbox, CSS Grid

### ユーザビリティ
- **ホバー効果**: ボタンや要素の視覚的フィードバック
- **ローディング**: 非同期操作中の状態表示
- **エラー表示**: 分かりやすいエラーメッセージ
- **レスポンシブ**: モバイル対応

### アクセシビリティ
- **キーボード操作**: TreeView, Modal, フォーム
- **ARIA属性**: 適切なラベルと説明
- **カラーコントラスト**: 読みやすい色の組み合わせ

## API仕様とコード生成

### TypeScript型生成
- **設定**: `openapi-ts.config.ts`
- **入力**: `../schemas/openapi.yaml`
- **出力**: `./app/api/`
- **コマンド**: `npm run generate-api`

### 生成されるファイル
- `types.gen.ts`: TypeScript型定義
- `sdk.gen.ts`: API呼び出し関数
- `client.gen.ts`: クライアント設定

## パフォーマンス最適化

### React最適化
- **React.memo**: 不要な再レンダリング防止
- **useCallback**: 関数の安定化
- **useMemo**: 計算結果のキャッシュ
- **lazy loading**: 必要時のコンポーネント読み込み

### バンドル最適化
- **Vite**: 高速なビルドとHMR
- **Tree shaking**: 未使用コードの除去
- **コード分割**: ルートベースの分割

## 開発時の注意事項

### ファイル命名規則
- **コンポーネント**: PascalCase (例: ProjectDetailModal.tsx)
- **ユーティリティ**: camelCase (例: timestamp.ts)
- **CSS**: kebab-case (例: gantt.css)

### インポート規則
- **React**: 名前付きインポートを優先
- **MUI**: 必要なコンポーネントのみインポート
- **相対パス**: `../` を使用して明確に

### TypeScript
- **型安全性**: any型の使用を最小限に
- **インターフェース**: Props型の明確な定義
- **null check**: optional chainingとnullish coalescingを活用

## メモリーズ

- CalendarPickerはMUI DatePickerベースのプロフェッショナルなカレンダー
- ファイル名変更時は推奨ファイル名を自動生成する
- 工程表の表示は開始日順でソートされている
- UIの一貫性を保つため水平レイアウトを採用
- ガントチャートは`/projects/gantt`に配置（projectsの下位ルート）
- TreeViewの展開状態はReact.useCallbackで最適化済み