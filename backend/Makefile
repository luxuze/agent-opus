.PHONY: proto proto-clean build run test clean help

# Proto related commands
proto: ## Generate Go code from proto files
	@echo "Generating Go code from proto files..."
	@mkdir -p api/proto/gen
	@protoc --go_out=api/proto/gen --go_opt=paths=source_relative \
		--go-grpc_out=api/proto/gen --go-grpc_opt=paths=source_relative \
		-I api/proto \
		api/proto/*.proto
	@echo "Proto generation complete!"

proto-clean: ## Clean generated proto files
	@echo "Cleaning generated proto files..."
	@rm -rf api/proto/gen
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

install-tools: ## Install required tools
	@echo "Installing protoc plugins..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Tools installed!"

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
