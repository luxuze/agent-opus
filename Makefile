# Makefile for Agent Platform

.PHONY: help dev build start stop clean logs test proto proto-backend proto-frontend proto-clean install-proto-tools

# Default target
help:
	@echo "Agent Platform - Available commands:"
	@echo "  make dev              - Start development environment"
	@echo "  make build            - Build all services"
	@echo "  make start            - Start all services"
	@echo "  make stop             - Stop all services"
	@echo "  make clean            - Clean up containers and volumes"
	@echo "  make logs             - Show logs"
	@echo "  make test             - Run tests"
	@echo "  make backend          - Build backend only"
	@echo "  make frontend         - Build frontend only"
	@echo ""
	@echo "Proto generation:"
	@echo "  make proto            - Generate proto for both backend and frontend"
	@echo "  make proto-backend    - Generate Go code from proto files"
	@echo "  make proto-frontend   - Generate TypeScript code from proto files"
	@echo "  make proto-clean      - Clean all generated proto files"
	@echo "  make install-proto-tools - Install protoc plugins"

# Development
dev:
	docker-compose up

# Build all services
build:
	docker-compose build

# Start services in background
start:
	docker-compose up -d
	@echo "Services started successfully!"

# Stop services
stop:
	docker-compose stop

# Clean up
clean:
	docker-compose down -v
	@echo "Cleaned up containers and volumes"

# Show logs
logs:
	docker-compose logs -f

# Run tests
test:
	cd backend && go test ./...
	cd frontend && npm run test

# Backend only
backend:
	cd backend && go build -o agent-platform ./cmd/server

# Frontend only
frontend:
	cd frontend && npm run build

# Database migration
migrate:
	cd backend && go run cmd/migrate/main.go

# Install dependencies
install:
	cd backend && go mod download
	cd frontend && npm install

# Format code
fmt:
	cd backend && go fmt ./...
	cd frontend && npm run lint

# Proto generation
proto: proto-backend proto-frontend ## Generate proto for both backend and frontend
	@echo "All proto files generated successfully!"

proto-backend: ## Generate Go code from proto files
	@echo "Generating Go code from proto files..."
	@mkdir -p backend/gen/go
	@protoc --go_out=backend/gen/go --go_opt=paths=source_relative \
		--go-grpc_out=backend/gen/go --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=backend/gen/go --grpc-gateway_opt=paths=source_relative \
		--grpc-gateway_opt=generate_unbound_methods=true \
		-I proto \
		proto/*.proto
	@echo "Backend proto generation complete!"

proto-frontend: ## Generate TypeScript code from proto files
	@echo "Generating TypeScript code from proto files..."
	@if ! command -v protoc-gen-grpc-web >/dev/null 2>&1; then \
		echo ""; \
		echo "⚠️  protoc-gen-grpc-web not found!"; \
		echo ""; \
		echo "Please install the gRPC-Web protoc plugin:"; \
		echo "  macOS:   brew install grpc-web"; \
		echo "  or download from: https://github.com/grpc/grpc-web/releases"; \
		echo ""; \
		echo "Note: Frontend proto generation is optional for backend development."; \
		echo "      You can skip this step if you're only working on the backend."; \
		echo ""; \
		exit 0; \
	fi
	@mkdir -p frontend/src/proto
	@protoc -I proto \
		--js_out=import_style=commonjs:frontend/src/proto \
		--grpc-web_out=import_style=typescript,mode=grpcwebtext:frontend/src/proto \
		proto/*.proto
	@echo "Frontend proto generation complete!"

proto-clean: ## Clean all generated proto files
	@echo "Cleaning generated proto files..."
	@rm -rf backend/gen/go
	@rm -rf backend/api/proto/gen
	@rm -rf frontend/src/proto
	@echo "Proto clean complete!"

install-proto-tools: ## Install protoc plugins
	@echo "Installing protoc plugins..."
	@echo ""
	@echo "Installing Go plugins..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	@echo "✓ Go plugins installed (including gRPC-Gateway)"
	@echo ""
	@echo "Installing frontend plugins..."
	@if command -v brew >/dev/null 2>&1; then \
		echo "Using Homebrew to install protoc-gen-grpc-web and protoc-gen-js..."; \
		brew install protoc-gen-grpc-web protoc-gen-js || true; \
		echo "✓ Frontend plugins installed"; \
	else \
		echo "Homebrew not found. Please install manually:"; \
		echo "  Download protoc-gen-grpc-web from:"; \
		echo "  https://github.com/grpc/grpc-web/releases"; \
		echo "  Download protoc-gen-js from:"; \
		echo "  https://github.com/protocolbuffers/protobuf-javascript/releases"; \
	fi
	@echo ""
	@echo "All proto tools installation complete!"
