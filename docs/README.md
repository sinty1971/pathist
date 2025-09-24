# API ドキュメント

`backend/cmd/main.go` で `DocsPath = "/docs"` を設定しているため、バックエンド起動後は `http://localhost:8080/docs` から Huma のドキュメント UI にアクセスできます。

OpenAPI 定義はルートディレクトリの `schemas/` 配下に保存され、`just docs` コマンドでブラウザを開くことができます。
