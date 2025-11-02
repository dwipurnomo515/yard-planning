.PHONY: help run build test clean migrate-up migrate-down db-create

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Run the application
	@echo "Starting server..."
	go run cmd/api/main.go

build: ## Build the application
	@echo "Building..."
	go build -o bin/api cmd/api/main.go
	@echo "Build complete! Binary: bin/api"

test: ## Run tests
	@echo "Running tests..."
	go test -v -cover ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean: ## Clean build files
	@echo "Cleaning..."
	rm -rf bin/
	@echo "Clean complete!"

db-create: ## Create database
	@echo "Creating database..."
	createdb yard_planning
	@echo "Database created!"

migrate-up: ## Run database migrations
	@echo "Running migrations..."
	psql -U postgres -d yard_planning -f migrations/001_init_schema.sql
	@echo "Migrations complete!"

migrate-down: ## Drop all tables
	@echo "Dropping all tables..."
	psql -U postgres -d yard_planning -c "DROP TABLE IF EXISTS containers, yard_plans, blocks, yards CASCADE;"
	@echo "Tables dropped!"

install: ## Install dependencies
	@echo "Installing dependencies..."
	go mod download
	@echo "Dependencies installed!"

lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run
	@echo "Lint complete!"

docker-up: ## Start Docker containers
	docker-compose up -d

docker-down: ## Stop Docker containers
	docker-compose down

api-test: ## Test API endpoints
	@chmod +x test_api.sh
	@./test_api.sh