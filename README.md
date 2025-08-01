# CLI Inventory Management

This is a command-line interface (CLI) application for managing inventory. It allows users to add products, manage stock levels, move inventory between locations, and generate reports.

## Features

- Add new products to the inventory
- Add stock for existing products at specific locations
- Find products by SKU
- Move stock between locations with atomic transactions
- Generate low-stock reports

## Technical Stack

- Language: Go
- Database: PostgreSQL
- Docker and Docker Compose for containerization
- SQLC for type-safe database queries
- Go Playground Validator for input validation

## Getting Started

### Prerequisites

- Go 1.21 or later
- Docker and Docker Compose
- PostgreSQL client (for direct database access if needed)
- SQLC (for code generation from SQL queries)

### Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd cli-inventory
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Install SQLC:
   ```bash
   go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
   ```

4. Generate Go code from SQL queries:
   ```bash
   sqlc generate
   ```

5. Start the database using Docker Compose:
   ```bash
   docker-compose up -d
   ```

6. Run database migrations:
   ```bash
   migrate -path migrations -database "postgres://inventory_user:inventory_password@localhost:5432/inventory_db?sslmode=disable" up
   ```

### Building the Application

Using Makefile:
```bash
make build
```

Or manually:
```bash
go build -o bin/inventory cmd/inventory/main.go
```

### Running the Application

```bash
./bin/inventory [command] [arguments]
```

## Usage

### Add a Product

```bash
./bin/inventory add-product <sku> <name> <description> <price>
```

Example:
```bash
./bin/inventory add-product PROD001 "Laptop" "High-performance laptop" 1299.99
```

### Add Stock

```bash
./bin/inventory add-stock <product-id> <location-id> <quantity>
```

Example:
```bash
./bin/inventory add-stock 1 1 50
```

### Find a Product

```bash
./bin/inventory find-product <sku>
```

Example:
```bash
./bin/inventory find-product PROD001
```

### Move Stock

```bash
./bin/inventory move-stock <product-id> <from-location-id> <to-location-id> <quantity>
```

Example:
```bash
./bin/inventory move-stock 1 1 2 10
```

### Generate Report

```bash
./bin/inventory generate-report low-stock [threshold]
```

Example:
```bash
./bin/inventory generate-report low-stock 20
```

## Database Schema

The application uses a PostgreSQL database with the following tables:

- `products`: Stores product definitions
- `locations`: Stores location information
- `stock`: Stores stock levels for each product at each location
- `stock_movements`: Tracks all stock movements for audit purposes

## Testing

### Unit Tests

To run unit tests:
```bash
make test
```

Or manually:
```bash
go test ./...
```

### Integration Tests

To run integration tests:
```bash
make integration-test
```

Or manually:
```bash
docker-compose -f docker-compose.test.yml up --build
```

## Project Structure

```
.
├── cmd/inventory/main.go          # Main application entry point
├── go.mod                         # Go module definition
├── go.sum                         # Go module checksums
├── Makefile                       # Build and test commands
├── sqlc.yaml                      # SQLC configuration
├── Dockerfile                     # Docker configuration for the application
├── docker-compose.yml             # Docker Compose configuration
├── docker-compose.test.yml       # Docker Compose configuration for testing
├── migrations/                   # Database migration files
│   ├── 000001_create_tables.up.sql
│   └── 000001_create_tables.down.sql
├── internal/
│   ├── database/                 # Database connection and utilities
│   ├── db/                       # Generated SQLC code
│   ├── models/                   # Data models
│   ├── repository/               # Data access layer
│   └── service/                  # Business logic layer
└── queries/                      # SQL queries
    ├── products.sql
    ├── stock.sql
    ├── locations.sql
    └── stock_movements.sql
```

## License

This project is licensed under the MIT License.
