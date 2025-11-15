.PHONY: proto proto-clean build run test clean help

# Proto related commands
proto: ## Generate Go code from proto files
	@echo "Generating Go code from proto files..."
	@mkdir -p gen/go
	@protoc --go_out=gen/go --go_opt=paths=source_relative \
		--go-grpc_out=gen/go --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=gen/go --grpc-gateway_opt=paths=source_relative \
		--grpc-gateway_opt=generate_unbound_methods=true \
		-I ../proto \
		../proto/*.proto
	@echo "Proto generation complete!"

proto-clean: ## Clean generated proto files
	@echo "Cleaning generated proto files..."
	@rm -rf gen/go
	@echo "Clean complete!"

# Build commands
build: ## Build the application
	@echo "Building application..."
	@go build -o bin/agent-platform ./cmd/server
	@echo "Build complete!"

run: ## Run the application
	@echo "Running application..."
	@go run ./cmd/server/main.go

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

clean: proto-clean ## Clean all build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@echo "Clean complete!"

# Development commands
dev: proto build run ## Generate proto, build and run

install-tools: ## Install required Go tools (backend only)
	@echo "Installing Go protoc plugins..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "✓ Go tools installed!"
	@echo ""
	@echo "For full installation (including frontend tools), run from project root:"
	@echo "  cd .. && make install-proto-tools"

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
