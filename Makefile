# Makefile for Agent Platform

.PHONY: help dev build start stop clean logs test

# Default target
help:
	@echo "Agent Platform - Available commands:"
	@echo "  make dev        - Start development environment"
	@echo "  make build      - Build all services"
	@echo "  make start      - Start all services"
	@echo "  make stop       - Stop all services"
	@echo "  make clean      - Clean up containers and volumes"
	@echo "  make logs       - Show logs"
	@echo "  make test       - Run tests"
	@echo "  make backend    - Build backend only"
	@echo "  make frontend   - Build frontend only"

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
	@echo "Frontend: http://localhost:3000"
	@echo "Backend:  http://localhost:8080"

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
