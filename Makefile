.PHONY: help dev build up down logs clean install

BINARY := goma-admin
BUILD_DIR := bin
UI_DIR := web

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

dev: ## Start development dependencies (postgres)
	@echo "Starting development environment..."
	docker compose up -d goma-postgres
	@echo "Waiting for services to be ready..."

install: ## Install all dependencies
	@echo "Installing backend dependencies..."
	go mod download
	@echo "Installing frontend dependencies..."
	cd $(UI_DIR) && npm install
	@echo "Done!"

build-ui: ## Build frontend only
	cd $(UI_DIR) && npm install && npm run build

build-api: ## Build backend only
	go build -o $(BUILD_DIR)/$(BINARY) cmd/main.go

build: build-ui build-api ## Build both API and UI

up: ## Start all services with Docker Compose
	docker compose up -d

down: ## Stop all services
	docker compose down

logs: ## Show logs from all services
	docker compose logs -f

clean: ## Clean build artifacts and dependencies
	@echo "Cleaning API..."
	rm -rf $(BUILD_DIR)/ tmp/
	@echo "Cleaning UI..."
	cd $(UI_DIR) && rm -rf dist/ node_modules/
	@echo "Cleaning Docker volumes..."
	docker compose down -v
	@echo "Done!"

test: ## Run Go tests
	go test -v ./...

run-api: ## Run backend locally
	go run cmd/main.go

run-ui: ## Run frontend dev server
	cd $(UI_DIR) && npm run dev

run: build ## Build and run the full application
	./$(BUILD_DIR)/$(BINARY)
fmt:
	gofmt -w .
tidy:
	go mod tidy

lint:
	golangci-lint run