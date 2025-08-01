# Penguin Backend & Frontend Commands

# Start the backend server (HTTP/1.1)
backend:
    cd ./backend && go run cmd/main.go

# Start the backend server with HTTP/2 + HTTPS
backend-http2:
    cd ./backend && go run cmd/main.go -http2

# Start the backend server with custom port
backend-port PORT:
    cd ./backend && go run cmd/main.go -port={{PORT}}

# Start the backend server with HTTP/2 on custom port
backend-http2-port PORT:
    cd ./backend && go run cmd/main.go -http2 -port={{PORT}}

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

# Generate OpenAPI documentation from Go code (using swag v2.0.0-rc4)
generate-api:
    cd ./backend && go run github.com/swaggo/swag/v2/cmd/swag@v2.0.0-rc4 init -g cmd/main.go --v3.1
    cp ./backend/docs/swagger.yaml ./schemas/openapi.yaml
    cp ./backend/docs/swagger.json ./schemas/openapi.json
    cp ./backend/docs/swagger.yaml ./schemas/openapi-v3.yaml
    cp ./backend/docs/swagger.json ./schemas/openapi-v3.json
    @echo "✅ OpenAPI 3.1 documentation generated and copied to schemas/"

# Convert Swagger 2.0 to OpenAPI 3.0 (legacy - no longer needed with swag v2)
convert-openapi3:
    cd ./backend/scripts && node convert-to-openapi3.js
    @echo "OpenAPI 3.0 documentation generated at schemas/openapi-v3.{yaml,json}"

# Install frontend dependencies  
frontend-deps:
    cd ./frontend && npm install

# Update all npm packages to latest versions
frontend-update:
    cd ./frontend && npm update
    cd ./frontend && npm audit fix

# Generate TypeScript types from OpenAPI spec
generate-types:
    cd ./frontend && npm run generate-api
    @echo "TypeScript types generated at frontend/app/api/"

# Generate React Router v7 route structure diagram
generate-routes:
    cd ./frontend && npm run generate-routes
    @echo "Route structure diagram generated at frontend/route-structure-generated.md"

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

# Generate both API docs and TypeScript types
generate-all: generate-api generate-types

# Generate all documentation (API, types, and routes)
generate-docs: generate-api generate-types generate-routes

# Update all dependencies (Go and npm)
update-all: backend-update frontend-update

# Clean and reinstall all dependencies
clean-install: backend-deps frontend-deps

# Stop backend server (Go application)
stop-backend:
    #!/bin/bash
    set +e
    echo "Stopping backend server..."
    pkill -f "go run cmd/main.go" 2>/dev/null
    pkill -f "main.go" 2>/dev/null
    pkill -f "penguin-backend" 2>/dev/null
    lsof -ti:8080 | xargs -r kill -15 2>/dev/null
    sleep 1
    lsof -ti:8080 | xargs -r kill -9 2>/dev/null
    echo "Backend server stopped"

# Stop HTTP/2 backend server
stop-backend-http2:
    #!/bin/bash
    set +e
    echo "Stopping HTTP/2 backend server..."
    pkill -f "go run cmd/main.go -http2" 2>/dev/null
    pkill -f "main.go.*http2" 2>/dev/null
    lsof -ti:8443 | xargs -r kill -15 2>/dev/null
    sleep 1
    lsof -ti:8443 | xargs -r kill -9 2>/dev/null
    echo "HTTP/2 backend server stopped"

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
stop-all: stop-backend stop-backend-http2 stop-frontend
    @echo "All servers stopped"

# Restart backend server (HTTP/1.1)
restart-backend: stop-backend
    @echo "Restarting backend server..."
    @sleep 1
    cd ./backend && go run cmd/main.go &
    @echo "Backend server restarted"

# Restart backend server (HTTP/2 + HTTPS)
restart-backend-http2: stop-backend-http2
    @echo "Restarting HTTP/2 backend server..."
    @sleep 1
    cd ./backend && go run cmd/main.go -http2 &
    @echo "HTTP/2 backend server restarted"

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
    cd ./backend && go run cmd/main.go &
    @sleep 1
    cd ./frontend && npm run dev &
    @echo "All servers restarted"

# Kill process running on port 8080 (legacy command)
kill-port:
    @echo "Stopping process on port 8080..."
    @-lsof -ti:8080 | xargs -r kill -15 2>/dev/null || true
    @sleep 1
    @-lsof -ti:8080 | xargs -r kill -9 2>/dev/null || true
    @echo "Port 8080 cleanup completed"

# Update claude-code
claude-code:
    npm i -g @anthropic-ai/claude-code

# Show backend architecture diagram in browser
architecture:
    @echo "Opening architecture diagram..."
    @xdg-open "https://mermaid.live/edit#$(cat doc/backend-architecture.md | grep -A 100 '```mermaid' | grep -B 100 '```' | grep -v '```' | base64 -w 0)" 2>/dev/null || open "https://mermaid.live/" 2>/dev/null || echo "Please visit https://mermaid.live/ and paste the mermaid code from doc/backend-architecture.md"

# Open Swagger UI documentation in browser (HTTP/1.1)
swagger:
    @echo "Opening Swagger UI documentation..."
    @xdg-open "http://localhost:8080/swagger/index.html" 2>/dev/null || open "http://localhost:8080/swagger/index.html" 2>/dev/null || echo "Please visit http://localhost:8080/swagger/index.html"

# Open Swagger UI documentation in browser (HTTP/2 + HTTPS)
swagger-http2:
    @echo "Opening Swagger UI documentation (HTTP/2)..."
    @xdg-open "https://localhost:8080/swagger/index.html" 2>/dev/null || open "https://localhost:8080/swagger/index.html" 2>/dev/null || echo "Please visit https://localhost:8080/swagger/index.html"

# Show available commands
help:
    @just --list