# Repository Guidelines

## プロジェクト構成とモジュール
- `backend/` は Go Fiber ベースの API で、`cmd/main.go` がエントリポイント、`internal/{endpoints,handlers,models,routes,services,utils}` にドメインロジックを集約しています。`docs/` は Swagger 生成物、`scripts/` はユーティリティスクリプト、証明書関連は `cert.pem` と `key.pem` を参照してください。
- `frontend/` は React Router v7 + TypeScript 構成です。画面ルートは `app/routes/`、UI コンポーネントは `app/components/`、API 型は `app/api/` に生成され、スタイル共通化は `app/styles/` で管理します。
- `schemas/` には OpenAPI 3.1 定義があり、バックエンド更新後はここも更新します。横断的な開発コマンドはリポジトリ直下の `justfile` にまとまっています。

## ビルド・テスト・開発コマンド
- `just backend` / `just backend-http2` で API サーバーを起動します。カスタムポートは `just backend-port <port>` を使用してください。
- `just frontend` で開発サーバー、`just frontend-build` で本番ビルド、`just frontend-preview` で生成物の確認ができます。
- バックエンドテストは `go test ./...`、ベンチマークは `go test -bench . ./...`。API 仕様と型定義は `just generate-api` と `just generate-types`、ルート図は `just generate-routes` で再生成します。

## コーディングスタイルと命名規則
- Go コードは `gofmt` 準拠 (タブインデント)。公開シンボルは UpperCamelCase、内部スコープは lowerCamelCase とし、ハンドラは `<Resource>Handler`、サービスは `<Resource>Service` 命名を守ります。
- フロントエンドは ESLint (`eslint.config.js`) と TypeScript に従い 2 スペースインデント。コンポーネント/フックは PascalCase、ユーティリティは camelCase、API クライアントは `services/` 配下にまとめます。
- スタイルは Tailwind と MUI テーマを併用します。トークンやカラーパレットは `tailwind.config.js` と `app/styles/*.css` を同期させ、カスタムクラスは接頭辞 `penguin-` を推奨します。

## テスト指針
- Go の単体テストは `_test.go` に配置し、テーブル駆動 (`cases := []struct{...}`) で失敗時ログを明瞭にします。ベンチマークは `benchmark_test.go` に `Benchmark<機能>` を追加してください。
- フロントエンドは現状自動テストが限定的なため、UI 変更時にはローカル実行結果のスクリーンショットまたはショートクリップを PR に添付し、必要に応じて `@testing-library/react` を用いたスモークテストを追加します。
- 主要ドメインの変更時は `go test -cover ./...` でカバレッジを確認し、低下が大きい場合はフォローアップのテスト追加を計画してください。

## コミットと PR ガイドライン
- コミットメッセージは日本語一文で主要変更点を具体的に述べ、複数トピックがある場合は本文で箇条書きを追加します。不要な生成物や鍵ファイルは含めません。
- PR では概要、テスト結果、影響範囲、関連 Issue を明記し、API 変更時は `schemas/` と Swagger の再生成結果をリンク、UI 変更時はビフォー/アフターを添付してください。
- `generate-cert.sh` で生成する証明書は開発用途のみです。鍵類はローカル管理とし、再現手順だけを共有すること。

## セキュリティと構成のヒント
- HTTP/2 を使う場合は `just generate-cert` でローカル証明書を再生成し、配布リポジトリにはコミットしないでください。運用環境では実証済みの証明書を `-cert` と `-key` フラグで指定します。
- バックエンドはデフォルトで `~/penguin` 以下のファイルツリーを参照します。別ディレクトリを使用する場合はサービス初期化ロジックを更新し、アクセス権限とバックアップ方針を事前に決めてください。
- Fiber の CORS 設定は全許可になっているため、公開環境では許可オリジンをホワイトリスト形式に更新し、ロギングや TLS 設定の値も環境変数やフラグで上書きしてください。
