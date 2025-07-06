#!/bin/bash

# HTTP/2対応の自己署名証明書を生成するスクリプト

echo "🔐 HTTP/2対応の自己署名証明書を生成中..."

# 既存の証明書があるかチェック
if [ -f "cert.pem" ] && [ -f "key.pem" ]; then
    echo "⚠️  既存の証明書が見つかりました。"
    read -p "上書きしますか？ (y/N): " overwrite
    if [ "$overwrite" != "y" ] && [ "$overwrite" != "Y" ]; then
        echo "❌ 証明書生成をキャンセルしました。"
        exit 1
    fi
fi

# 自己署名証明書を生成
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes \
    -subj "/C=JP/ST=Tokyo/L=Tokyo/O=Development/OU=IT/CN=localhost" \
    -extensions v3_req \
    -config <(cat <<EOF
[req]
distinguished_name = req_distinguished_name
req_extensions = v3_req
prompt = no

[req_distinguished_name]
C = JP
ST = Tokyo
L = Tokyo
O = Development
OU = IT Department
CN = localhost

[v3_req]
keyUsage = keyEncipherment, dataEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
DNS.2 = 127.0.0.1
IP.1 = 127.0.0.1
IP.2 = ::1
EOF
)

if [ $? -eq 0 ]; then
    echo "✅ 証明書が正常に生成されました！"
    echo ""
    echo "📁 生成されたファイル:"
    echo "  - cert.pem (証明書)"
    echo "  - key.pem  (秘密鍵)"
    echo ""
    echo "🚀 HTTP/2サーバーを起動するには:"
    echo "  go run cmd/main_http2.go"
    echo ""
    echo "🌐 ブラウザでアクセス:"
    echo "  https://localhost:8443"
    echo "  https://localhost:8443/swagger/index.html"
    echo ""
    echo "⚠️  ブラウザで「安全でない」警告が表示されますが、"
    echo "   「詳細設定」→「localhost に進む」で続行してください。"
else
    echo "❌ 証明書の生成に失敗しました。"
    echo "OpenSSLがインストールされているか確認してください。"
    exit 1
fi