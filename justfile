# Grpc Backend & Frontend Commands

# Start the backend gRPC server (HTTP/2 over h2c)
backend:
    cd ./backend && go run cmd/grpc/main.go

# Start the backend gRPC server with TLS enabled
backend-tls:
    cd ./backend && go run cmd/grpc/main.go -enable-tls

# Start the frontend development server (React Router v7)
frontend:
    cd ./frontend && npm run dev

# Install backend dependencies
backend-deps:
    cd ./backend && go mod tidy

# Generate SSL certificate for HTTP/2
generate-cert:
    cd ./backend && ./generate-cert.sh

# Update all Go packages to latest versions
# ただしメジャーなバージョンは更新しない
backend-update:
    cd ./backend && go get -u ./...
    cd ./backend && go mod tidy

# Generate gRPC stubs for Koji service
generate-grpc:
    buf generate --path proto/v1/penguin.proto

# Generate Connect-Web stubs for the frontend
frontend-generate-grpc: generate-grpc
    @echo "Connect-Web stubs generated at frontend/src/gen/"

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

# Start both backend and frontend (requires tmux or run in separate terminals)
dev:
    @echo "Starting backend and frontend..."
    @echo "Run 'just backend' in one terminal and 'just frontend' in another"

# Generate both backend and frontend gRPC stubs
generate-all: generate-grpc generate-types

# Update all dependencies (Go and npm)
update-all: backend-update frontend-update

# Clean and reinstall all dependencies
clean-install: backend-deps frontend-deps

# Stop backend server (Go application)
stop-backend:
    #!/bin/bash
    set +e
    echo "Stopping backend server..."
    pkill -f "go run cmd/grpc/main.go" 2>/dev/null
    pkill -f "cmd/grpc/main.go" 2>/dev/null
    pkill -f "backend-grpc" 2>/dev/null
    for port in 9090 9443; do
        lsof -ti:$port | xargs -r kill -15 2>/dev/null
    done
    sleep 1
    for port in 9090 9443; do
        lsof -ti:$port | xargs -r kill -9 2>/dev/null
    done
    echo "Backend server stopped"

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

# Stop both backend and frontend servers
stop-all: stop-backend stop-frontend
    @echo "All servers stopped"

# Restart backend server (HTTP/1.1)
restart-backend: stop-backend
    @echo "Restarting backend server..."
    @sleep 1
    cd ./backend && go run cmd/grpc/main.go &
    @echo "Backend server restarted"

# Restart frontend development server
restart-frontend: stop-frontend
    @echo "Restarting frontend development server..."
    @sleep 1
    cd ./frontend && npm run dev &
    @echo "Frontend development server restarted"

# Restart both backend and frontend servers
restart-all: stop-all
    @echo "Restarting all servers..."
    @sleep 2
    cd ./backend && go run cmd/grpc/main.go &
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

# Show backend architecture diagram in browser
architecture:
    @echo "Opening architecture diagram..."
    @xdg-open "https://mermaid.live/edit#$(cat doc/backend-architecture.md | grep -A 100 '```mermaid' | grep -B 100 '```' | grep -v '```' | base64 -w 0)" 2>/dev/null || open "https://mermaid.live/" 2>/dev/null || echo "Please visit https://mermaid.live/ and paste the mermaid code from doc/backend-architecture.md"

# Show available commands
help:
    @just --list
