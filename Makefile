.PHONY: swagger
swagger:
	@echo "Generating Swagger documentation..."
	@go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/main.go -o ./docs
	@echo "Swagger documentation generated successfully!"

.PHONY: run
run:
	@go run main.go

.PHONY: build
build:
	@go build -o bin/ecommerce-be main.go

.PHONY: test
test:
	@echo "Running all tests..."
	@go test -v ./...

.PHONY: test-unit
test-unit:
	@echo "Running unit tests..."
	@go test -v -short ./repositories/... ./services/... ./controllers/...

.PHONY: test-integration
test-integration:
	@echo "Running integration tests..."
	@go test -v ./tests/integration/...

.PHONY: setup-test-db
setup-test-db:
	@echo "Setting up test database..."
	@if [ -f scripts/setup-test-db.sh ]; then \
		bash scripts/setup-test-db.sh; \
	else \
		echo "Creating test database manually..."; \
		createdb ecommerce_test 2>/dev/null || echo "Database might already exist or PostgreSQL not accessible"; \
	fi

.PHONY: test-all
test-all: setup-test-db
	@echo "Running all tests (unit + integration)..."
	@go test -v ./...

.PHONY: test-db-up
test-db-up:
	@echo "Starting PostgreSQL test database in Docker..."
	@docker-compose -f docker-compose.test.yml up -d postgres-test
	@echo "Waiting for PostgreSQL to be ready..."
	@sleep 3
	@echo "PostgreSQL test database is ready!"

.PHONY: test-db-down
test-db-down:
	@echo "Stopping PostgreSQL test database..."
	@docker-compose -f docker-compose.test.yml down

.PHONY: test-db-restart
test-db-restart: test-db-down test-db-up
	@echo "PostgreSQL test database restarted!"

.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: test-coverage-html
test-coverage-html: test-coverage
	@echo "Opening coverage report in browser..."
	@xdg-open coverage.html 2>/dev/null || open coverage.html 2>/dev/null || echo "Please open coverage.html manually"

.PHONY: test-race
test-race:
	@echo "Running tests with race detector..."
	@go test -v -race ./...

.PHONY: test-benchmark
test-benchmark:
	@echo "Running benchmark tests..."
	@go test -v -bench=. -benchmem ./...
