.PHONY: help build run test clean docker-up docker-down docker-logs migrate-up migrate-down

help:
	@echo "Available commands:"
	@echo "  make build         - Build Go application"
	@echo "  make run           - Run application locally"
	@echo "  make test          - Run tests"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make docker-up     - Start Docker containers"
	@echo "  make docker-down   - Stop Docker containers"
	@echo "  make docker-logs   - View Docker logs"
	@echo "  make migrate-up    - Run database migrations"
	@echo "  make migrate-down  - Rollback database migrations"

build:
	cd cmd/api && go build -o ../../bin/api main.go

run:
	go run cmd/api/main.go

test:
	go test -v ./... -cover

clean:
	rm -rf bin/

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f api

docker-build:
	docker-compose build

migrate-up:
	go run cmd/migrate/main.go up

migrate-down:
	go run cmd/migrate/main.go down