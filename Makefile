.PHONY: help build run test clean docker-build docker-up docker-down migrate-up migrate-down migrate-create swagger

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	go build -o bin/server ./cmd/server

run: ## Run the application
	go run ./cmd/server/main.go

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

test-integration: ## Run integration tests
	go test -v -tags=integration ./...

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html

deps: ## Download dependencies
	go mod download
	go mod tidy

fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

lint: ## Run linters
	golangci-lint run

# Database migration targets
migrate-up: ## Run all database migrations
	@echo "Running migrations..."
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/subspace?sslmode=disable" up

migrate-down: ## Rollback last database migration
	@echo "Rolling back migration..."
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/subspace?sslmode=disable" down 1

migrate-create: ## Create a new migration (usage: make migrate-create name=migration_name)
	@if [ -z "$(name)" ]; then echo "Usage: make migrate-create name=migration_name"; exit 1; fi
	migrate create -ext sql -dir migrations -seq $(name)

migrate-version: ## Show current migration version
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/subspace?sslmode=disable" version

migrate-force: ## Force set migration version (usage: make migrate-force version=1)
	@if [ -z "$(version)" ]; then echo "Usage: make migrate-force version=1"; exit 1; fi
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/subspace?sslmode=disable" force $(version)

# Swagger/OpenAPI targets
swagger: ## Generate Swagger documentation
	swag init -g cmd/server/main.go -o docs

swagger-install: ## Install swag CLI tool
	go install github.com/swaggo/swag/cmd/swag@latest

# Docker targets
docker-build: ## Build Docker image
	docker build -t subspace-backend:latest .

docker-up: ## Start Docker containers
	docker compose up -d

docker-down: ## Stop Docker containers
	docker compose down

docker-logs: ## View Docker logs
	docker compose logs -f

docker-migrate: ## Run migrations inside Docker
	docker compose exec api migrate -path /app/migrations -database "postgres://postgres:postgres@postgres:5432/subspace?sslmode=disable" up

# Development targets
dev: ## Run in development mode with hot reload (requires air)
	air

dev-migrate: ## Run application with auto-migration enabled
	AUTO_MIGRATE=true go run ./cmd/server/main.go
