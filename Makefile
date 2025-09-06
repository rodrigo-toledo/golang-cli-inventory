.PHONY: generate build test integration-test test-coverage integration-test-coverage test-all clean openapi-validate test-openapi docs coverage

# Generate Go code from SQL queries
generate:
	sqlc generate

# Build the application
build:
	go build -o bin/inventory cmd/inventory/main.go

# Run unit tests
test:
	go test ./...

# Run integration tests
integration-test:
	docker-compose -f docker-compose.test.yml run --rm app go test ./...

# Run unit tests with coverage
test-coverage:
	go test -coverprofile=coverage.out -covermode=count ./...
	go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//' > coverage_percentage.txt
	@if [ $(cat coverage_percentage.txt) -lt 90 ]; then \
		echo "❌ Test coverage is below 90% (current: $(cat coverage_percentage.txt)%)"; \
		exit 1; \
	else \
		echo "✅ Test coverage is $(cat coverage_percentage.txt)% (meets 90% threshold)"; \
	fi
	go tool cover -html=coverage.out -o coverage.html

# Measure and display current coverage
coverage:
	@echo "📊 Measuring test coverage..."
	@go test -coverprofile=coverage.out -covermode=count ./... >/dev/null
	@echo ""
	@echo "📈 Coverage by package:"
	@go tool cover -func=coverage.out | grep -v "total:" | awk '{print $1 ": " $3}' | sort
	@echo ""
	@echo "📈 Overall coverage:"
	@go tool cover -func=coverage.out | grep "total:" | awk '{print $3}'
	@echo ""
	@echo "📄 HTML coverage report generated: coverage.html"
	@go tool cover -html=coverage.out -o coverage.html

# Run integration tests with coverage
integration-test-coverage:
	docker-compose -f docker-compose.test.yml run --rm app go test -coverprofile=coverage.out ./...
	docker-compose -f docker-compose.test.yml run --rm app go tool cover -html=coverage.out -o coverage.html

# Run all tests (unit + integration) with comprehensive output
test-all:
	@echo "🧪 Running all tests (unit + integration)..."
	@echo "=========================================="
	@echo "📋 Running unit tests..."
	@go test -v ./...
	@echo "✅ Unit tests completed"
	@echo ""
	@echo "🔧 Running integration tests..."
	@docker-compose -f docker-compose.test.yml run --rm app go test -v ./...
	@echo "✅ Integration tests completed"
	@echo ""
	@echo "🎉 All tests completed successfully!"

# Clean generated files
clean:
	rm -rf internal/db
	rm -rf bin
	rm -rf coverage.out coverage.html coverage_percentage.txt

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