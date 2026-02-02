# =============================================================================
# aLMS Backend Makefile
# =============================================================================

.PHONY: help build run test clean \
        docker-up docker-down docker-logs docker-migrate docker-reset docker-build \
        migrate-up migrate-down migrate-create

# Default target
.DEFAULT_GOAL := help

# Variables
APP_NAME := alms-api
GO := go
GOFLAGS := -v
DOCKER_COMPOSE := docker compose

# =============================================================================
# Development
# =============================================================================

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

## build: Build the API binary
build:
	CGO_ENABLED=0 $(GO) build $(GOFLAGS) -o bin/$(APP_NAME) ./cmd/api

## run: Run the API locally
run:
	$(GO) run ./cmd/api

## test: Run all tests
test:
	$(GO) test -v -race -cover ./...

## clean: Clean build artifacts
clean:
	rm -rf bin/
	$(GO) clean -cache

## lint: Run linter
lint:
	golangci-lint run ./...

## fmt: Format code
fmt:
	$(GO) fmt ./...
	goimports -w .

# =============================================================================
# Docker
# =============================================================================

## docker-up: Start all services with docker-compose (builds if needed)
docker-up:
	$(DOCKER_COMPOSE) up --build -d

## docker-down: Stop all docker-compose services
docker-down:
	$(DOCKER_COMPOSE) down

## docker-logs: Tail logs from all services
docker-logs:
	$(DOCKER_COMPOSE) logs -f

## docker-logs-api: Tail logs from API service only
docker-logs-api:
	$(DOCKER_COMPOSE) logs -f api

## docker-migrate: Re-run database migrations
docker-migrate:
	$(DOCKER_COMPOSE) up --build migrate

## docker-reset: Stop all services and remove volumes (full reset)
docker-reset:
	$(DOCKER_COMPOSE) down -v
	@echo "All containers stopped and volumes removed."

## docker-build: Build docker images without starting
docker-build:
	$(DOCKER_COMPOSE) build

## docker-ps: Show running containers
docker-ps:
	$(DOCKER_COMPOSE) ps

## docker-shell-api: Open shell in API container
docker-shell-api:
	$(DOCKER_COMPOSE) exec api /bin/sh

## docker-shell-postgres: Open psql in postgres container
docker-shell-postgres:
	$(DOCKER_COMPOSE) exec postgres psql -U alms -d alms_db

# =============================================================================
# Local Migrations (requires golang-migrate CLI)
# =============================================================================

## migrate-up: Run all pending migrations (local)
migrate-up:
	migrate -path ./migrations -database "$(DATABASE_URL)" up

## migrate-down: Rollback last migration (local)
migrate-down:
	migrate -path ./migrations -database "$(DATABASE_URL)" down 1

## migrate-create: Create a new migration (usage: make migrate-create NAME=migration_name)
migrate-create:
	migrate create -ext sql -dir ./migrations -seq $(NAME)

# =============================================================================
# Utilities
# =============================================================================

## deps: Download dependencies
deps:
	$(GO) mod download
	$(GO) mod tidy

## swagger: Generate swagger docs (requires swag CLI)
swagger:
	swag init -g cmd/api/main.go -o docs/

## seed: Run API data seeder (requires running API)
seed:
	$(GO) run ./scripts/seeder/main.go

## seed-url: Run seeder with custom API URL (usage: make seed-url URL=http://localhost:8080/api/v1)
seed-url:
	API_URL=$(URL) $(GO) run ./scripts/seeder/main.go
# =============================================================================
# Testing with Docker
# =============================================================================

## test-docker-up: Start test database containers
test-docker-up:
	$(DOCKER_COMPOSE) -f docker-compose.test.yml up -d
	@echo "Waiting for test services to be healthy..."
	@sleep 3

## test-docker-down: Stop and remove test containers
test-docker-down:
	$(DOCKER_COMPOSE) -f docker-compose.test.yml down -v

## test-migrate: Run migrations against test database
test-migrate:
	migrate -path ./migrations -database "postgres://alms_test:alms_test_secret@localhost:5433/alms_test_db?sslmode=disable" up

## test-migrate-down: Rollback all migrations in test database
test-migrate-down:
	migrate -path ./migrations -database "postgres://alms_test:alms_test_secret@localhost:5433/alms_test_db?sslmode=disable" down -all

## test-unit: Run unit tests only (no build tags)
test-unit:
	$(GO) test -v -race -cover ./internal/...

## test-integration: Run integration tests (requires test db)
test-integration: test-docker-up test-migrate
	$(GO) test -v -race -tags=integration ./tests/integration/...

## test-e2e: Run e2e tests (requires running API and test db)
test-e2e: test-docker-up test-migrate
	TEST_API_URL=http://localhost:8080 $(GO) test -v -tags=e2e ./tests/e2e/...

## test-all: Run all tests (unit, integration, e2e)
test-all: test-unit test-integration

## test-coverage: Run tests with coverage report
test-coverage:
	$(GO) test -v -race -coverprofile=coverage.out ./internal/...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## test-clean: Clean up test artifacts and stop test containers
test-clean: test-docker-down
	rm -f coverage.out coverage.html