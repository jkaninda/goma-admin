.PHONY: help dev build up down logs clean install

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

dev: ## Start development environment
	@echo "Starting development environment..."
	docker compose up -d goma-postgres goma-redis
	@echo "Waiting for services to be ready..."

install: ## Install all dependencies
	@echo "Installing backend dependencies..."
	go mod download
	@echo "Installing frontend dependencies..."
	cd ui && npm install
	@echo "Done!"

build: ## Build both api and UI
	@echo "Building API..."
	cd go build -o bin/goma-admin cmd/main.go
	@echo "Building Ui..."
	cd ui && npm run build
	@echo "Done!"

up: ## Start all services with Docker Compose
	docker compose up -d

down: ## Stop all services
	docker compose down

logs: ## Show logs from all services
	docker compose logs -f

clean: ## Clean build artifacts and dependencies
	@echo "Cleaning API..."
	cd rm -rf bin/ tmp/
	@echo "Cleaning UI..."
	cd ui && rm -rf dist/ node_modules/
	@echo "Cleaning Docker volumes..."
	docker compose down -v
	@echo "Done!"


test-api: ## Run api tests
	go test -v ./...

build-api: ## Build backend only
	 go build -o bin/goma-admin cmd/main.go

ui-build: ## Build ui only
	cd ui && npm run build

run-api: ## Run api only
	go run cmd/main.go

run-ui: ## Run ui only
	cd ui && npm run dev
