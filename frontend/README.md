# フォルダー管理システム - フロントエンド

Next.js (App Router) + TypeScript をベースにしたフォルダー管理システムのフロントエンドアプリケーションです。バックエンドの gRPC (Connect-Web) API を通じてファイル・工事・会社情報を扱います。

## 技術スタック

- **Next.js** 14.x
- **React** 18.x
- **TypeScript** 5.9
- **MUI / @mui/x-tree-view**
- **Tailwind CSS + DaisyUI**
- **@connectrpc/connect / connect-web** – gRPC over HTTP/1.1 をブラウザから利用

## ディレクトリ構成 (主要箇所)

```
frontend/
├── app/                     # Next.js App Router
│   ├── layout.tsx           # ルートレイアウト（Providers, Navigation）
│   ├── page.tsx             # ホーム
│   ├── files/page.tsx       # ファイル一覧
│   ├── kojies/page.tsx      # 工事一覧
│   └── kojies/gantt/page.tsx# ガントチャート
├── src/
│   ├── components/          # UI コンポーネント
│   ├── contexts/            # 状態管理コンテキスト
│   ├── gen/                 # `npm run generate-grpc` で生成される Connect-Web スタブ
│   ├── services/            # gRPC / Connect-Web クライアントラッパー
│   ├── types/               # フロントエンド共通モデル定義
│   ├── styles/              # 共通スタイル
│   └── utils/               # ユーティリティ関数
├── public/                  # 静的アセット
├── next.config.mjs
├── package.json
└── tsconfig.json
```

## 開発コマンド

開発前に依存関係を再インストールしてください（React Router 依存削除後の再ロック推奨）。

```bash
npm install
```

開発サーバー: http://localhost:3000

```bash
npm run dev
```

ビルド / 本番起動:

```bash
npm run build
npm run start
```

Lint と gRPC スタブ再生成:

```bash
npm run lint
npm run generate-grpc  # ../proto から ./src/gen に TypeScript を生成
```

## API 仕様との連携

バックエンドの `proto/penguin/v1/penguin.proto` から Connect-Web 用のクライアントスタブを生成しています。`just generate-grpc` または `npm run generate-grpc` を実行すると `src/gen` 配下が再生成されます。生成物を直接編集しないでください。

## 主な UI

- **ファイル一覧**: TreeView によるファイル構造の閲覧、詳細モーダル表示。
- **工事一覧**: 工事エントリーの編集、補助ファイルの状況可視化、ガントチャート遷移。
- **会社一覧**: カテゴリーフィルタ / 詳細モーダルによる編集。
- **ナビゲーション**: Next.js の Link と usePathname を利用したアクティブ表示。

## 注意事項

- ルーティングは Next.js のディレクトリ構造に従います。React Router 用スクリプト `scripts/generate-route-diagram.js` は非推奨となりました。
- グローバルスタイルは `app/globals.css` で読み込み、コンポーネント固有のスタイルは `src/styles/` からインポートしてください。
- Context Provider (`KojiProvider`, `FileInfoProvider`) は `app/providers.tsx` で一括ラップしています。
- 旧来の Vite / React Router 向け設定ファイルは削除済みです。不要なキャッシュは `rm -rf frontend/node_modules` 後に再インストールしてください。

## トラブルシューティング

- **依存関係エラー**: `package-lock.json` が古い場合は `rm package-lock.json` → `npm install` で再生成してください。
- **API 型の不整合**: `just generate-grpc` または `npm run generate-grpc` で最新のスタブを再生成し、必要に応じて `npm run lint` を実行してください。
- **スタイル崩れ**: Tailwind の監視対象パス (`tailwind.config.js`) に `app` / `src` が含まれていることを確認してください。

## ライセンス

社内利用を想定した非公開プロジェクトのため、ライセンスは設定していません。
