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

- Language: Go 1.25 (with experimental JSON v2 package)
- Database: PostgreSQL 15
- Docker and Docker Compose for containerization
- SQLC for type-safe database queries
- Cobra CLI framework for command-line interface
- Testify for testing utilities

## Getting Started

### Prerequisites

- Go 1.25 or later (with JSON v2 experimental package)
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

### HTTP API Server

The application can also be started as an HTTP server to expose a RESTful API.

#### Starting the Server

Ensure the database is running (`docker-compose up -d`). Then, start the API server:

```bash
./bin/inventory serve
```

The server will start on `http://localhost:8080`.

#### API Endpoints

The API provides the following endpoints. All requests and responses use JSON.

**Base URL:** `http://localhost:8080/api/v1`

---

**Products**

*   **List all products**
    *   `GET /products`
    *   **Response:** `200 OK` with an array of product objects.
    *   **Example `curl`:**
        ```bash
        curl http://localhost:8080/api/v1/products
        ```

*   **Get a single product by SKU**
    *   `GET /products/{sku}`
    *   **Response:** `200 OK` with a single product object.
    *   **Example `curl`:**
        ```bash
        curl http://localhost:8080/api/v1/products/PROD001
        ```

*   **Create a new product**
    *   `POST /products`
    *   **Request Body:** `CreateProductRequest` object.
        ```json
        {
          "sku": "PROD003",
          "name": "Wireless Mouse",
          "description": "Ergonomic wireless mouse",
          "price": 25.50
        }
        ```
    *   **Response:** `201 Created` with the created product object.
    *   **Example `curl`:**
        ```bash
        curl -X POST http://localhost:8080/api/v1/products \
        -H "Content-Type: application/json" \
        -d '{"sku":"PROD003","name":"Wireless Mouse","description":"Ergonomic wireless mouse","price":25.50}'
        ```

---

**Locations**

*   **List all locations**
    *   `GET /locations`
    *   **Response:** `200 OK` with an array of location objects.
    *   **Example `curl`:**
        ```bash
        curl http://localhost:8080/api/v1/locations
        ```

*   **Get a single location by name**
    *   `GET /locations/{name}`
    *   **Response:** `200 OK` with a single location object.
    *   **Example `curl`:**
        ```bash
        curl http://localhost:8080/api/v1locations/Main%20Warehouse
        ```

*   **Create a new location**
    *   `POST /locations`
    *   **Request Body:** `CreateLocationRequest` object.
        ```json
        {
          "name": "Secondary Warehouse"
        }
        ```
    *   **Response:** `201 Created` with the created location object.
    *   **Example `curl`:**
        ```bash
        curl -X POST http://localhost:8080/api/v1/locations \
        -H "Content-Type: application/json" \
        -d '{"name":"Secondary Warehouse"}'
        ```

---

**Stock**

*   **Add stock to a product at a location**
    *   `POST /stock/add`
    *   **Request Body:** `AddStockRequest` object.
        ```json
        {
          "product_id": 1,
          "location_id": 1,
          "quantity": 100
        }
        ```
    *   **Response:** `200 OK` with the updated stock object for that product/location.
    *   **Example `curl`:**
        ```bash
        curl -X POST http://localhost:8080/api/v1/stock/add \
        -H "Content-Type: application/json" \
        -d '{"product_id":1,"location_id":1,"quantity":100}'
        ```

*   **Move stock between locations**
    *   `POST /stock/move`
    *   **Request Body:** `MoveStockRequest` object.
        ```json
        {
          "product_id": 1,
          "from_location_id": 1,
          "to_location_id": 2,
          "quantity": 10
        }
        ```
    *   **Response:** `200 OK` with the stock object at the destination location after the move.
    *   **Example `curl`:**
        ```bash
        curl -X POST http://localhost:8080/api/v1/stock/move \
        -H "Content-Type: application/json" \
        -d '{"product_id":1,"from_location_id":1,"to_location_id":2,"quantity":10}'
        ```

*   **Get low stock report**
    *   `GET /stock/low-stock?threshold={threshold}`
    *   **Query Parameter:** `threshold` (optional, integer, defaults to 10).
    *   **Response:** `200 OK` with an array of stock objects where quantity is below the threshold.
    *   **Example `curl`:**
        ```bash
        # Get stock below 5 units
        curl http://localhost:8080/api/v1/stock/low-stock?threshold=5

        # Get stock below default threshold (10)
        curl http://localhost:8080/api/v1/stock/low-stock
        ```

#### Error Responses

*   **`400 Bad Request`**: Invalid JSON payload, missing required fields, or invalid input values (e.g., negative quantity).
*   **`404 Not Found`**: Resource not found (e.g., product with a given SKU does not exist). *Note: Currently, most "not found" scenarios return `500 Internal Server Error`, but this is planned to be improved to `404`.*
*   **`500 Internal Server Error`**: Unexpected server-side errors (e.g., database connection issues, service layer errors not specifically handled).

### Add a Product (CLI)

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

## JSON v2 Migration

This project uses the experimental JSON v2 package introduced in Go 1.25. To build and run the project with the new JSON implementation, you need to enable the `jsonv2` experiment:

```bash
GOEXPERIMENT=jsonv2 go build -o bin/inventory cmd/inventory/main.go
```

Or to run directly:

```bash
GOEXPERIMENT=jsonv2 go run cmd/inventory/main.go [command] [arguments]
```

The JSON v2 package provides performance improvements and new features while maintaining compatibility with the existing `encoding/json` package.

### Benefits of JSON v2

- Improved performance for JSON encoding and decoding
- Better error messages
- Enhanced streaming capabilities
- More efficient memory usage

For more information about the JSON v2 package, see the [Go 1.25 release notes](https://go.dev/doc/go1.25#json_v2).

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
