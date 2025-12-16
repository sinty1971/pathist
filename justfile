set shell := ["powershell", "-NoLogo", "-Command"]

# buf コマンドのパス設定
buf := `"$env:USERPROFILE/prj/bin/buf-windows-x86_64.exe"`

# gRPC Server & Frontend Commands

# Start the gRPC server (HTTP/2 over h2c)
server-grpc:
    cd ./server-grpc ; go run cmd/grpc/main.go

# Test file service client
fileclient *ARGS:
    cd ./server-grpc ; go run cmd/fileclient/main.go {{ARGS}}

# Test company service client
companyclient *ARGS:
    cd ./server-grpc ; go run cmd/companyclient/main.go {{ARGS}}

# Start the gRPC server with TLS enabled
server-tls:
    cd ./server-grpc ; go run cmd/grpc/main.go -enable-tls

# Start the frontend development server (React Router v7)
frontend:
    cd ./frontend ; npm run dev

# Start the fe-bun development server (React Router v7 on Bun)
fe-bun:
    cd ./fe-bun ; npm run dev

# Generate SSL certificate for HTTP/2
generate-cert:
    cd ./server-grpc ; ./generate-cert.sh

# Update all Go packages to latest versions
# ただしメジャーなバージョンは更新しない
server-grpc-update:
    cd ./server-grpc ; go get -u ./...
    cd ./server-grpc ; go mod tidy

# Generate gRPC stubs for all services
generate-grpc:
    @echo "Generating gRPC stubs..."
    @ {{buf}} generate
    
# Generate Connect-Web stubs for the frontend
frontend-generate-grpc: generate-grpc
    @echo "Connect-Web stubs generated at frontend/src/gen/"

# Generate Connect-Web stubs for fe-bun (React Router v7)
fe-bun-generate: generate-grpc
    @echo "Connect-Web stubs generated at fe-bun/app/gen/"

# Install frontend dependencies  
frontend-deps:
    cd ./frontend && npm install

# Update all npm packages to latest versions
frontend-update:
    cd ./frontend && npm update
    cd ./frontend && npm audit fix

# Generate TypeScript types from OpenAPI spec
generate-types: frontend-generate-grpc

# Generate React Router v7 route structure diagram
generate-routes:
    @echo "Next.js へ移行したため自動ルート図生成は未対応です"
    @echo "app ディレクトリ構成を直接参照してください"

# Build frontend for production (React Router v7)
frontend-build:
    cd ./frontend && npm run build

# Preview frontend production build
frontend-preview:
    cd ./frontend && npm run preview

# Run frontend linting
frontend-lint:
    cd ./frontend && npm run lint

# Start both server-grpc and frontend (requires tmux or run in separate terminals)
dev:
    @echo "Starting server-grpc and frontend..."
    @echo "Run 'just server-grpc' in one terminal and 'just frontend' in another"

# Generate both server-grpc and frontend gRPC stubs
generate-all: generate-grpc generate-types fe-bun-generate

# Install server-grpc dependencies
server-grpc-deps:
    cd ./server-grpc && go mod tidy
    cd ./server-grpc && go mod download

# Update all dependencies (Go and npm)
update-all: server-grpc-update frontend-update

# Clean and reinstall all dependencies
clean-install: server-grpc-deps frontend-deps

# Stop server-grpc server (Go application)
stop-server-grpc:
    #!/bin/bash
    set +e
    echo "Stopping server-grpc server..."
    pkill -f "go run cmd/grpc/main.go" 2>/dev/null
    pkill -f "cmd/grpc/main.go" 2>/dev/null
    pkill -f "server-grpc" 2>/dev/null
    for port in 9090 9443; do
        lsof -ti:$port | xargs -r kill -15 2>/dev/null
    done
    sleep 1
    for port in 9090 9443; do
        lsof -ti:$port | xargs -r kill -9 2>/dev/null
    done
    echo "gRPC server stopped"

# Stop frontend development server (React Router v7 / Vite)
stop-frontend:
    #!/bin/bash
    set +e
    echo "Stopping frontend development server..."
    pkill -f "npm run dev" 2>/dev/null
    pkill -f "react-router dev" 2>/dev/null
    pkill -f "vite" 2>/dev/null
    pkill -f "node.*vite" 2>/dev/null
    for port in 5173 5174 5175 5176 5177; do
        lsof -ti:$port | xargs -r kill -15 2>/dev/null
    done
    sleep 1
    for port in 5173 5174 5175 5176 5177; do
        lsof -ti:$port | xargs -r kill -9 2>/dev/null
    done
    echo "Frontend development server stopped"

# Stop both server-grpc and frontend servers
stop-all: stop-server-grpc stop-frontend
    @echo "All servers stopped"

# Restart gRPC server (HTTP/1.1)
restart-server-grpc: stop-server-grpc
    @echo "Restarting gRPC server..."
    @sleep 1
    cd ./server-grpc && go run cmd/grpc/main.go &
    @echo "gRPC server restarted"

# Restart frontend development server
restart-frontend: stop-frontend
    @echo "Restarting frontend development server..."
    @sleep 1
    cd ./frontend && npm run dev &
    @echo "Frontend development server restarted"

# Restart both server-grpc and frontend servers
restart-all: stop-all
    @echo "Restarting all servers..."
    @sleep 2
    cd ./server-grpc && go run cmd/grpc/main.go &
    @sleep 1
    cd ./frontend && npm run dev &
    @echo "All servers restarted"

# Kill process running on port 9090 (legacy command)
kill-port:
    @echo "Stopping process on port 9090..."
    @-lsof -ti:9090 | xargs -r kill -15 2>/dev/null || true
    @sleep 1
    @-lsof -ti:9090 | xargs -r kill -9 2>/dev/null || true
    @echo "Port 9090 cleanup completed"

# Update claude-code
claude-code:
    npm i -g @anthropic-ai/claude-code

# Show server-grpc architecture diagram in browser
architecture:
    @echo "Opening architecture diagram..."
    @xdg-open "https://mermaid.live/edit#$(cat doc/server-architecture.md | grep -A 100 '```mermaid' | grep -B 100 '```' | grep -v '```' | base64 -w 0)" 2>/dev/null || open "https://mermaid.live/" 2>/dev/null || echo "Please visit https://mermaid.live/ and paste the mermaid code from doc/server-architecture.md"

# Show API documentation in browser
docs:
    @echo "Opening API documentation..."
    @if command -v xdg-open > /dev/null; then xdg-open docs/proto/apis.md; elif command -v open > /dev/null; then open docs/proto/apis.md; else cat docs/proto/apis.md; fi

# Show available commands
help:
    @just --list
