.PHONY: help install build start stop restart logs clean test lint

help:
	@echo "3X-UI Modern - Makefile Commands"
	@echo "================================="
	@echo "install     - Install dependencies and setup"
	@echo "build       - Build Docker images"
	@echo "start       - Start all services"
	@echo "stop        - Stop all services"
	@echo "restart     - Restart all services"
	@echo "logs        - View logs"
	@echo "clean       - Clean up containers and volumes"
	@echo "test        - Run all tests"
	@echo "lint        - Run linters"
	@echo "backup      - Create backup"

install:
	@chmod +x scripts/*.sh
	@./scripts/install.sh

build:
	docker-compose build

start:
	docker-compose up -d
	@echo "Services started. Access at http://localhost:8080"

stop:
	docker-compose down

restart:
	docker-compose restart

logs:
	docker-compose logs -f

clean:
	docker-compose down -v
	@echo "Cleaned up containers and volumes"

test:
	@echo "Running backend tests..."
	cd backend && go test ./... -v -race -coverprofile=coverage.out
	@echo "Running frontend tests..."
	cd frontend && npm test

lint:
	@echo "Linting backend..."
	cd backend && golangci-lint run
	@echo "Linting frontend..."
	cd frontend && npm run lint

backup:
	@chmod +x scripts/backup.sh
	@./scripts/backup.sh
