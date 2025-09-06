.PHONY: generate build test integration-test test-coverage integration-test-coverage test-all clean

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
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	open coverage.html

# Run integration tests with coverage
integration-test-coverage:
	docker-compose -f docker-compose.test.yml run --rm app go test -coverprofile=coverage.out ./...
	docker-compose -f docker-compose.test.yml run --rm app go tool cover -html=coverage.out -o coverage.html

# Run all tests (unit + integration)
test-all: test integration-test

# Clean generated files
clean:
	rm -rf internal/db
	rm -rf bin
	rm -rf coverage.out coverage.html
