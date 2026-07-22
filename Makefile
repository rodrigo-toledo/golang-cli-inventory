.PHONY: generate build test unit-test integration-test test-coverage integration-test-coverage test-all clean openapi-validate test-openapi docs coverage mocks

# Generate Go code from SQL queries
generate:
	go tool sqlc generate

# Generate mocks using mockery
mocks:
	go tool mockery --config=.mockery.yml

# Build the application with JSON v2 experiment enabled
build:
	GOEXPERIMENT=jsonv2 go build -o bin/inventory cmd/inventory/main.go

# Run unit tests with JSON v2 experiment enabled
unit-test:
	go tool mockery --config=.mockery.yml
	GOEXPERIMENT=jsonv2 go tool gotestsum --format dots -- -v ./internal/... -tags=unit

# Run unit tests with JSON v2 experiment enabled
test:
	go tool mockery --config=.mockery.yml
	GOEXPERIMENT=jsonv2 go tool gotestsum --format dots -- -v ./internal/... -tags=unit

# Run integration tests with JSON v2 experiment enabled
integration-test:
	docker-compose -f docker-compose.test.yml up --abort-on-container-exit --exit-code-from app

# Run unit tests with coverage and JSON v2 experiment enabled
unit-test-coverage:
	go tool mockery --config=.mockery.yml
	GOEXPERIMENT=jsonv2 go test -coverprofile=coverage.out -covermode=count ./internal/... -tags=unit
	GOEXPERIMENT=jsonv2 go tool cover -func=coverage.out | grep "total:" | awk '{print $$3}' | sed 's/%//' > coverage_percentage.txt
	@coverage=$$(cat coverage_percentage.txt | cut -d. -f1); \
	if [ $$coverage -lt 90 ]; then \
		echo "❌ Test coverage is below 90% (current: $$(cat coverage_percentage.txt)%)"; \
		exit 1; \
	else \
		echo "✅ Test coverage is $$(cat coverage_percentage.txt)% (meets 90% threshold)"; \
	fi
	GOEXPERIMENT=jsonv2 go tool cover -html=coverage.out -o coverage.html

# Run unit tests with coverage and JSON v2 experiment enabled
test-coverage:
	go tool mockery --config=.mockery.yml
	GOEXPERIMENT=jsonv2 go test -coverprofile=coverage.out -covermode=count ./internal/... -tags=unit
	GOEXPERIMENT=jsonv2 go tool cover -func=coverage.out | grep "total:" | awk '{print $$3}' | sed 's/%//' > coverage_percentage.txt
	@coverage=$$(cat coverage_percentage.txt | cut -d. -f1); \
	if [ $$coverage -lt 90 ]; then \
		echo "❌ Test coverage is below 90% (current: $$(cat coverage_percentage.txt)%)"; \
		exit 1; \
	else \
		echo "✅ Test coverage is $$(cat coverage_percentage.txt)% (meets 90% threshold)"; \
	fi
	GOEXPERIMENT=jsonv2 go tool cover -html=coverage.out -o coverage.html

# Measure and display current coverage
coverage:
	@scripts/coverage.sh

# Run integration tests with coverage and JSON v2 experiment enabled
integration-test-coverage:
	@echo "🧪 Running integration tests with coverage..."
	@if [ -f coverage_integration.out ]; then \
		rm coverage_integration.out; \
	fi
	docker-compose -f docker-compose.test.yml up -d
	@echo "⏳ Waiting for tests to complete..."
	@until [ "$$(docker inspect -f '{{.State.Running}}' inventory-integration-test-app 2>/dev/null)" = "false" ] 2>/dev/null; do \
		sleep 1; \
	done
	@echo "📊 Extracting integration test coverage data..."
	docker cp inventory-integration-test-app:/app/coverage_integration.out . 2>/dev/null || true
	docker-compose -f docker-compose.test.yml down
	@if [ -f coverage_integration.out ]; then \
		echo "✅ Integration test coverage data extracted"; \
		GOEXPERIMENT=jsonv2 go tool cover -func=coverage_integration.out | grep total | awk '{print $3}' | sed 's/%//' > coverage_integration_percentage.txt; \
		echo "📈 Integration test coverage: $$(cat coverage_integration_percentage.txt)%"; \
		GOEXPERIMENT=jsonv2 go tool cover -html=coverage_integration.out -o coverage_integration.html; \
		echo "📄 HTML coverage report generated: coverage_integration.html"; \
	else \
		echo "❌ No integration test coverage data found"; \
	fi

# Run all tests (unit + integration) with comprehensive output and JSON v2 experiment enabled
test-all:
	@echo "🧪 Running all tests (unit + integration)..."
	@echo "=========================================="
	@echo "📋 Running unit tests..."
	@GOEXPERIMENT=jsonv2 go tool gotestsum --format dots -- -v ./internal/... -tags=unit
	@echo "✅ Unit tests completed"
	@echo ""
	@echo "🔧 Running integration tests..."
	@docker-compose -f docker-compose.test.yml up --abort-on-container-exit --exit-code-from app
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
	@GOEXPERIMENT=jsonv2 go run scripts/validate_openapi.go
	@echo "✅ OpenAPI specification validation completed"

# Run OpenAPI compliance tests
test-openapi:
	@echo "🧪 Running OpenAPI compliance tests..."
	@GOEXPERIMENT=jsonv2 go test -v ./internal/handlers
	@echo "✅ OpenAPI compliance tests completed"

# Generate API documentation
docs:
	@echo "📚 Generating API documentation..."
	@echo "📖 OpenAPI specification available at: api/openapi.yaml"
	@echo "🌐 You can view the API documentation using Swagger UI or Redoc"
	@echo "   - Online Swagger UI: https://editor.swagger.io/"
	@echo "   - Online Redoc: https://redocly.github.io/redoc/"
	@echo "   - Upload api/openapi.yaml to either platform to view interactive documentation"