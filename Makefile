.PHONY: generate build test integration-test clean

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
	docker-compose -f docker-compose.test.yml up --build

# Clean generated files
clean:
	rm -rf internal/db
	rm -rf bin
