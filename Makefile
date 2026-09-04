.PHONY: generate build test unit-test integration-test test-coverage unit-test-coverage integration-test-coverage test-all clean openapi-validate test-openapi docs coverage mocks lint

# Generate Go code from SQL queries
generate:
	go tool sqlc generate

# Generate mocks using mockery
mocks:
	go tool mockery --config=.mockery.yml

# Run linter
lint:
	@echo "🔍 Running linter..."
	@golangci-lint run ./...
	@echo "✅ Linting completed"

# Build the application
build:
	go build -o bin/inventory cmd/inventory/main.go

# Run unit tests
test:
	go tool mockery --config=.mockery.yml
	go tool gotestsum --format dots -- -v ./internal/...

# Alias kept for backwards compatibility
unit-test: test

# Run integration tests
integration-test:
	docker compose -f docker-compose.test.yml up --abort-on-container-exit --exit-code-from app

# Run unit tests with coverage
test-coverage:
	go tool mockery --config=.mockery.yml
	go test -coverprofile=coverage.out -covermode=count ./internal/...
	go tool cover -func=coverage.out | grep "total:" | awk '{print $$3}' | sed 's/%//' > coverage_percentage.txt
	@coverage=$$(cat coverage_percentage.txt | cut -d. -f1); \
	if [ $$coverage -lt 90 ]; then \
		echo "❌ Test coverage is below 90% (current: $$(cat coverage_percentage.txt)%)"; \
		exit 1; \
	else \
		echo "✅ Test coverage is $$(cat coverage_percentage.txt)% (meets 90% threshold)"; \
	fi
	go tool cover -html=coverage.out -o coverage.html

# Alias kept for backwards compatibility
unit-test-coverage: test-coverage

# Measure and display current coverage
coverage:
	@scripts/coverage.sh

# Run integration tests with coverage
integration-test-coverage:
	@echo "🧪 Running integration tests with coverage..."
	@if [ -f coverage_integration.out ]; then \
		rm coverage_integration.out; \
	fi
	docker compose -f docker-compose.test.yml up -d
	@echo "⏳ Waiting for tests to complete..."
	@until [ "$$(docker inspect -f '{{.State.Running}}' inventory-integration-test-app 2>/dev/null)" = "false" ] 2>/dev/null; do \
		sleep 1; \
	done
	@echo "📊 Extracting integration test coverage data..."
	docker cp inventory-integration-test-app:/app/coverage_integration.out . 2>/dev/null || true
	docker compose -f docker-compose.test.yml down
	@if [ -f coverage_integration.out ]; then \
		echo "✅ Integration test coverage data extracted"; \
		go tool cover -func=coverage_integration.out | grep total | awk '{print $3}' | sed 's/%//' > coverage_integration_percentage.txt; \
		echo "📈 Integration test coverage: $$(cat coverage_integration_percentage.txt)%"; \
		go tool cover -html=coverage_integration.out -o coverage_integration.html; \
		echo "📄 HTML coverage report generated: coverage_integration.html"; \
	else \
		echo "❌ No integration test coverage data found"; \
	fi

# Run all tests (unit + integration)
test-all:
	@echo "🧪 Running all tests (unit + integration)..."
	@echo "=========================================="
	@echo "📋 Running unit tests..."
	@go tool gotestsum --format dots -- -v ./internal/...
	@echo "✅ Unit tests completed"
	@echo ""
	@echo "🔧 Running integration tests..."
	@docker compose -f docker-compose.test.yml up --abort-on-container-exit --exit-code-from app
	@echo "✅ Integration tests completed"
	@echo "🎉 All tests completed successfully!"

# Clean generated files
clean:
	rm -rf internal/db
	rm -rf bin
	rm -rf coverage.out coverage.html coverage_percentage.txt
	rm -rf coverage_integration.out coverage_integration.html coverage_integration_percentage.txt

# Validate OpenAPI specification
openapi-validate:
	@echo "📋 Validating OpenAPI specification..."
	@go run scripts/validate_openapi.go
	@echo "✅ OpenAPI specification validation completed"

# Run OpenAPI compliance tests
test-openapi:
	@echo "🧪 Running OpenAPI compliance tests..."
	@go test -v ./internal/handlers
	@echo "✅ OpenAPI compliance tests completed"

# Generate API documentation
docs:
	@echo "📚 Generating API documentation..."
	@echo "📖 OpenAPI specification available at: api/openapi.yaml"
	@echo "🌐 You can view the API documentation using Swagger UI or Redoc"
	@echo "   - Online Swagger UI: https://editor.swagger.io/"
	@echo "   - Online Redoc: https://redocly.github.io/redoc/"
	@echo "   - Upload api/openapi.yaml to either platform to view interactive documentation"
