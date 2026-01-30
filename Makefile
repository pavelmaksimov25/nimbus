.PHONY: help install migrate migrate-up migrate-down migrate-create migrate-status run build clean docker-build docker-up docker-down docker-logs test lint docs

# Default target
.DEFAULT_GOAL := help

# Load environment variables from .env file
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Variables
APP_NAME=nimbus
BINARY_NAME=nimbus
DOCKER_IMAGE=$(APP_NAME):latest
MIGRATIONS_DIR=./migrations
GO_FILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*")

# Colors for output
COLOR_RESET=\033[0m
COLOR_BOLD=\033[1m
COLOR_GREEN=\033[32m
COLOR_YELLOW=\033[33m
COLOR_BLUE=\033[34m

##@ General

help: ## Display this help message
	@echo "$(COLOR_BOLD)$(APP_NAME) - Makefile Commands$(COLOR_RESET)"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make $(COLOR_BLUE)<target>$(COLOR_RESET)\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  $(COLOR_BLUE)%-20s$(COLOR_RESET) %s\n", $$1, $$2 } /^##@/ { printf "\n$(COLOR_BOLD)%s$(COLOR_RESET)\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

install: ## Install dependencies and setup environment
	@echo "$(COLOR_GREEN)Installing dependencies...$(COLOR_RESET)"
	@if [ ! -f .env ]; then \
		echo "$(COLOR_YELLOW)Creating .env file from .env.example...$(COLOR_RESET)"; \
		cp .env.example .env; \
		echo "$(COLOR_YELLOW)Please update .env file with your configuration$(COLOR_RESET)"; \
	fi
	@echo "$(COLOR_GREEN)Downloading Go modules...$(COLOR_RESET)"
	@go mod download
	@go mod verify
	@echo "$(COLOR_GREEN)Installing development tools...$(COLOR_RESET)"
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@echo "$(COLOR_GREEN)✓ Installation complete!$(COLOR_RESET)"

run: ## Run the application locally
	@echo "$(COLOR_GREEN)Starting $(APP_NAME)...$(COLOR_RESET)"
	@go run cmd/app/main.go

build: ## Build the application binary
	@echo "$(COLOR_GREEN)Building $(APP_NAME)...$(COLOR_RESET)"
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o bin/$(BINARY_NAME) ./cmd/app
	@echo "$(COLOR_GREEN)✓ Binary built: bin/$(BINARY_NAME)$(COLOR_RESET)"

clean: ## Clean build artifacts
	@echo "$(COLOR_YELLOW)Cleaning build artifacts...$(COLOR_RESET)"
	@rm -rf bin/
	@rm -rf tmp/
	@go clean
	@echo "$(COLOR_GREEN)✓ Clean complete!$(COLOR_RESET)"

##@ Database & Migrations

migrate: migrate-up ## Run all pending migrations (alias for migrate-up)

migrate-up: ## Run all pending database migrations
	@echo "$(COLOR_GREEN)Running migrations...$(COLOR_RESET)"
	@goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" up
	@echo "$(COLOR_GREEN)✓ Migrations applied successfully!$(COLOR_RESET)"

migrate-down: ## Rollback the last migration
	@echo "$(COLOR_YELLOW)Rolling back last migration...$(COLOR_RESET)"
	@goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" down
	@echo "$(COLOR_GREEN)✓ Migration rolled back!$(COLOR_RESET)"

migrate-create: ## Create a new migration file (usage: make migrate-create name=migration_name)
	@if [ -z "$(name)" ]; then \
		echo "$(COLOR_YELLOW)Usage: make migrate-create name=migration_name$(COLOR_RESET)"; \
		exit 1; \
	fi
	@echo "$(COLOR_GREEN)Creating migration: $(name)...$(COLOR_RESET)"
	@goose -dir $(MIGRATIONS_DIR) -s create $(name) sql
	@echo "$(COLOR_GREEN)✓ Migration file created!$(COLOR_RESET)"

migrate-status: ## Show migration status
	@echo "$(COLOR_BLUE)Migration status:$(COLOR_RESET)"
	@goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" status

migrate-redo: ## Re-run the latest migration
	@echo "$(COLOR_YELLOW)Re-running latest migration...$(COLOR_RESET)"
	@goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" redo
	@echo "$(COLOR_GREEN)✓ Migration re-applied!$(COLOR_RESET)"

##@ Docker

docker-build: ## Build Docker image
	@echo "$(COLOR_GREEN)Building Docker image...$(COLOR_RESET)"
	@docker build -t $(DOCKER_IMAGE) .
	@echo "$(COLOR_GREEN)✓ Docker image built: $(DOCKER_IMAGE)$(COLOR_RESET)"

docker-up: ## Start all services with docker-compose
	@echo "$(COLOR_GREEN)Starting services with docker-compose...$(COLOR_RESET)"
	@docker-compose up -d
	@echo "$(COLOR_GREEN)✓ Services started!$(COLOR_RESET)"
	@echo "$(COLOR_BLUE)View logs with: make docker-logs$(COLOR_RESET)"

docker-down: ## Stop all services
	@echo "$(COLOR_YELLOW)Stopping services...$(COLOR_RESET)"
	@docker-compose down
	@echo "$(COLOR_GREEN)✓ Services stopped!$(COLOR_RESET)"

docker-logs: ## View docker-compose logs
	@docker-compose logs -f

docker-restart: docker-down docker-up ## Restart all services

docker-clean: ## Remove all containers, volumes, and images
	@echo "$(COLOR_YELLOW)Cleaning Docker resources...$(COLOR_RESET)"
	@docker-compose down -v
	@docker rmi $(DOCKER_IMAGE) 2>/dev/null || true
	@echo "$(COLOR_GREEN)✓ Docker cleanup complete!$(COLOR_RESET)"

##@ Testing

test: ## Run tests
	@echo "$(COLOR_GREEN)Running tests...$(COLOR_RESET)"
	@go test -v -race -coverprofile=coverage.out ./...
	@echo "$(COLOR_GREEN)✓ Tests complete!$(COLOR_RESET)"

test-coverage: test ## Run tests with coverage report
	@echo "$(COLOR_BLUE)Generating coverage report...$(COLOR_RESET)"
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(COLOR_GREEN)✓ Coverage report: coverage.html$(COLOR_RESET)"

lint: ## Run linter (requires golangci-lint)
	@echo "$(COLOR_GREEN)Running linter...$(COLOR_RESET)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
		echo "$(COLOR_GREEN)✓ Linting complete!$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_YELLOW)golangci-lint not installed. Install with:$(COLOR_RESET)"; \
		echo "  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$$(go env GOPATH)/bin"; \
	fi

##@ Documentation

docs: ## Generate and serve API documentation
	@echo "$(COLOR_GREEN)API Documentation:$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_BLUE)📚 Available Documentation:$(COLOR_RESET)"
	@echo "  • README.md           - Project overview and setup"
	@echo "  • docs/migration.md   - Database migration guide"
	@echo "  • docs/apidoc/        - API documentation"
	@echo ""
	@echo "$(COLOR_YELLOW)To view API docs, check: docs/apidoc/apidoc.yaml$(COLOR_RESET)"
	@if command -v swagger >/dev/null 2>&1; then \
		echo "$(COLOR_GREEN)Starting Swagger UI...$(COLOR_RESET)"; \
		swagger serve docs/apidoc/apidoc.yaml; \
	else \
		echo "$(COLOR_YELLOW)Swagger CLI not installed. View the YAML file directly or install with:$(COLOR_RESET)"; \
		echo "  go install github.com/go-swagger/go-swagger/cmd/swagger@latest"; \
	fi

docs-migration: ## Open migration documentation
	@if command -v open >/dev/null 2>&1; then \
		open docs/migration.md; \
	elif command -v xdg-open >/dev/null 2>&1; then \
		xdg-open docs/migration.md; \
	else \
		cat docs/migration.md; \
	fi

##@ Database Management

db-start: ## Start only the database service
	@echo "$(COLOR_GREEN)Starting database...$(COLOR_RESET)"
	@docker-compose up -d database
	@echo "$(COLOR_GREEN)✓ Database started!$(COLOR_RESET)"

db-stop: ## Stop the database service
	@echo "$(COLOR_YELLOW)Stopping database...$(COLOR_RESET)"
	@docker-compose stop database
	@echo "$(COLOR_GREEN)✓ Database stopped!$(COLOR_RESET)"

db-shell: ## Open PostgreSQL shell
	@echo "$(COLOR_BLUE)Opening database shell...$(COLOR_RESET)"
	@docker-compose exec database psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

db-reset: ## Reset database (WARNING: destroys all data)
	@echo "$(COLOR_YELLOW)⚠️  WARNING: This will destroy all data!$(COLOR_RESET)"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		echo "$(COLOR_YELLOW)Resetting database...$(COLOR_RESET)"; \
		docker-compose down -v; \
		docker-compose up -d database; \
		sleep 5; \
		$(MAKE) migrate-up; \
		echo "$(COLOR_GREEN)✓ Database reset complete!$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_BLUE)Cancelled.$(COLOR_RESET)"; \
	fi

##@ Quick Start

setup: install docker-up ## Complete setup: install dependencies and start services
	@echo ""
	@echo "$(COLOR_GREEN)✓ Setup complete!$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_BOLD)Next steps:$(COLOR_RESET)"
	@echo "  1. Update .env file with your configuration"
	@echo "  2. Run: $(COLOR_BLUE)make run$(COLOR_RESET) to start the application"
	@echo "  3. Visit: http://localhost:8080"
	@echo ""

dev: db-start migrate run ## Start development environment (db + migrations + app)

all: clean install build test ## Run all: clean, install, build, and test
