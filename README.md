# CLI Inventory Management

This is a command-line interface (CLI) application for managing inventory. It allows users to add products, manage stock levels, move inventory between locations, and generate reports.

## Features

- Add new products to the inventory
- List all products in the inventory
- Find products by SKU
- Add stock for existing products at specific locations
- Move stock between locations with atomic transactions
- Generate low-stock reports

## Technical Stack

- Language: Go 1.24
- Database: PostgreSQL 15
- Docker and Docker Compose for containerization
- SQLC for type-safe database queries
- Cobra CLI framework for command-line interface
- Testify for testing utilities

## Getting Started

### Prerequisites

- Go 1.24 or later
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

   The database will be automatically initialized with migrations from the `migrations/` directory.

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

### List All Products

```bash
./bin/inventory list-products
```

Example:
```bash
./bin/inventory list-products
```

### Find a Product

```bash
./bin/inventory find-product <sku>
```

Example:
```bash
./bin/inventory find-product PROD001
```

### Add Stock

```bash
./bin/inventory add-stock <product-id> <location-id> <quantity>
```

Example:
```bash
./bin/inventory add-stock 1 1 50
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
./bin/inventory generate-report <report-type> [options]
```

Low-stock report example:
```bash
./bin/inventory generate-report low-stock 20
```

Available report types:
- `low-stock [threshold]` - Show products with stock below specified threshold

## Database Schema

The application uses a PostgreSQL database with the following tables:

### `products`
Stores product definitions:
- `id` (SERIAL PRIMARY KEY)
- `sku` (VARCHAR(50) UNIQUE NOT NULL)
- `name` (VARCHAR(255) NOT NULL)
- `description` (TEXT)
- `price` (DECIMAL(10, 2))
- `created_at` (TIMESTAMP WITH TIME ZONE DEFAULT NOW())

### `locations`
Stores location information:
- `id` (SERIAL PRIMARY KEY)
- `name` (VARCHAR(255) UNIQUE NOT NULL)
- `created_at` (TIMESTAMP WITH TIME ZONE DEFAULT NOW())

### `stock`
Stores stock levels for each product at each location:
- `id` (SERIAL PRIMARY KEY)
- `product_id` (INTEGER REFERENCES products(id) ON DELETE CASCADE)
- `location_id` (INTEGER REFERENCES locations(id) ON DELETE CASCADE)
- `quantity` (INTEGER NOT NULL DEFAULT 0)
- `created_at` (TIMESTAMP WITH TIME ZONE DEFAULT NOW())
- `updated_at` (TIMESTAMP WITH TIME ZONE DEFAULT NOW())
- UNIQUE constraint on (product_id, location_id)

### `stock_movements`
Tracks all stock movements for audit purposes:
- `id` (SERIAL PRIMARY KEY)
- `product_id` (INTEGER REFERENCES products(id) ON DELETE CASCADE)
- `from_location_id` (INTEGER REFERENCES locations(id) ON DELETE SET NULL)
- `to_location_id` (INTEGER REFERENCES locations(id) ON DELETE SET NULL)
- `quantity` (INTEGER NOT NULL)
- `movement_type` (VARCHAR(50) NOT NULL)
- `created_at` (TIMESTAMP WITH TIME ZONE DEFAULT NOW())

## Configuration

### Database Connection

The application uses the following environment variables for database configuration:

- `DATABASE_URL`: PostgreSQL connection string
  - Default: `postgres://inventory_user:inventory_password@db:5432/inventory_db?sslmode=disable`

### Docker Configuration

The `docker-compose.yml` file sets up:
- PostgreSQL 15 database
- Automatic migration execution on startup
- Data persistence via Docker volumes

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
docker-compose -f docker-compose.test.yml run --rm app go test ./...
```

### Test Coverage

To run unit tests with coverage report:
```bash
make test-coverage
```

To run integration tests with coverage report:
```bash
make integration-test-coverage
```

To run all tests (unit + integration):
```bash
make test-all
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
│   ├── cli/                      # Command-line interface
│   │   ├── root.go               # Root command and initialization
│   │   ├── product_commands.go   # Product-related commands
│   │   └── stock_commands.go     # Stock-related commands
│   ├── config/                   # Configuration management
│   ├── database/                 # Database connection and utilities
│   │   └── database.go
│   ├── db/                       # Generated SQLC code
│   │   ├── db.go
│   │   ├── models.go
│   │   ├── querier.go
│   │   ├── products.sql.go
│   │   ├── stock.sql.go
│   │   ├── locations.sql.go
│   │   └── stock_movements.sql.go
│   ├── models/                   # Data models
│   │   ├── product.go
│   │   ├── location.go
│   │   └── stock.go
│   ├── repository/               # Data access layer
│   │   ├── products.go
│   │   ├── locations.go
│   │   ├── stock.go
│   │   └── stock_movements.go
│   ├── service/                  # Business logic layer
│   │   ├── product.go
│   │   ├── location.go
│   │   └── stock.go
│   ├── testutils/                # Test utilities
│   │   ├── test_data.go
│   │   └── test_database.go
│   └── validation/               # Input validation
└── queries/                      # SQL queries
    ├── products.sql
    ├── stock.sql
    ├── locations.sql
    └── stock_movements.sql
```

## Development Workflow

1. **Make changes to SQL queries**: Edit files in the `queries/` directory
2. **Generate Go code**: Run `sqlc generate` to update `internal/db/`
3. **Update business logic**: Modify files in `internal/service/`
4. **Update CLI commands**: Modify files in `internal/cli/`
5. **Test changes**: Run `make test-all`
6. **Build application**: Run `make build`

## Cleaning Generated Files

To clean all generated files and start fresh:
```bash
make clean
```

This will remove:
- `internal/db/` directory (generated SQLC code)
- `bin/` directory (compiled binaries)
- Coverage files (`coverage.out`, `coverage.html`)

## License

This project is licensed under the MIT License.
