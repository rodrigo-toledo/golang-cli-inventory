# CLI Inventory Management System - Project Specification

## 1. Core Features / Commands

### add-product
**Purpose**: Add a new product definition to the inventory system.

**Arguments/Flags**:
- `--name` (required): Product name
- `--sku` (required): Stock Keeping Unit - unique product identifier
- `--description` (optional): Product description
- `--price` (required): Product price
- `--category` (optional): Product category

**Expected Success Output**:
```
Product "Product Name" (SKU: ABC123) successfully added with ID: 12345
```

**Key Error Scenarios**:
- SKU already exists in the database
- Missing required arguments
- Invalid price format
- Database connection failure

### add-stock
**Purpose**: Add inventory of an existing product to a specific location.

**Arguments/Flags**:
- `--sku` (required): Product SKU to add stock for
- `--quantity` (required): Quantity to add (must be positive integer)
- `--location` (required): Warehouse/location identifier
- `--lot-number` (optional): Lot or batch number for tracking

**Expected Success Output**:
```
Successfully added 50 units of product (SKU: ABC123) to location WH-01
```

**Key Error Scenarios**:
- Product with given SKU does not exist
- Invalid quantity (negative or zero)
- Location identifier is invalid
- Insufficient database permissions

### find-product
**Purpose**: Search for products by name or SKU.

**Arguments/Flags**:
- `--name` (optional): Product name or partial name to search for
- `--sku` (optional): Exact SKU to search for
- `--category` (optional): Filter by product category

**Expected Success Output**:
```
Found 2 product(s):
1. Product Name A (SKU: ABC123) - $29.99 - In stock: 150 units
2. Product Name B (SKU: DEF456) - $45.50 - In stock: 75 units
```

**Key Error Scenarios**:
- No products match the search criteria
- Both name and SKU provided but no match
- Database query error

### move-stock
**Purpose**: Move a specific quantity of a product from one location to another in an atomic transaction.

**Arguments/Flags**:
- `--sku` (required): Product SKU to move
- `--quantity` (required): Quantity to move
- `--from-location` (required): Source location
- `--to-location` (required): Destination location

**Expected Success Output**:
```
Successfully moved 25 units of product (SKU: ABC123) from WH-01 to WH-02
Transaction ID: 789bea1b-2c3d-4e5f-6a7b-8c9d0e1f2a3b
```

**Key Error Scenarios**:
- Insufficient stock at source location
- Product not found
- Source and destination locations are the same
- Database transaction failure

### generate-report
**Purpose**: Generate various inventory reports.

**Arguments/Flags**:
- `--type` (required): Report type (e.g., "low-stock", "inventory-summary")
- `--threshold` (optional): Minimum stock threshold (used with low-stock report)
- `--location` (optional): Filter report by specific location
- `--output` (optional): Output file path (default: stdout)

**Expected Success Output**:
```
Low Stock Report (Threshold: 10 units)
Generated: 2025-08-01 23:56:00

1. Product A (SKU: ABC123) - Location: WH-01 - Quantity: 5
2. Product C (SKU: GHI789) - Location: WH-02 - Quantity: 3
```

**Key Error Scenarios**:
- Invalid report type
- Database query error
- File write permission error when using --output flag

## 2. Database Schema

```sql
-- Products table for product definitions
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    sku VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    category VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Locations table for warehouse/storage locations
CREATE TABLE locations (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    address TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Stock levels table
CREATE TABLE stock (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    location_id INTEGER NOT NULL REFERENCES locations(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL DEFAULT 0,
    lot_number VARCHAR(100),
    expiration_date DATE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(product_id, location_id, lot_number)
);

-- Stock movement transaction log
CREATE TABLE stock_movements (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES products(id),
    from_location_id INTEGER REFERENCES locations(id),
    to_location_id INTEGER REFERENCES locations(id),
    quantity INTEGER NOT NULL,
    transaction_type VARCHAR(50) NOT NULL, -- 'add', 'remove', 'transfer'
    reference_id VARCHAR(100), -- External reference (e.g., order number)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_products_sku ON products(sku);
CREATE INDEX idx_stock_product_location ON stock(product_id, location_id);
CREATE INDEX idx_stock_movements_product ON stock_movements(product_id);
```

## 3. Technical Stack & Local Environment

### Language
Go (Golang)

### Database
PostgreSQL

### Key Go Libraries & Tools

1. **Database Access**:
   - `database/sql` package (standard library)
   - `github.com/jackc/pgx/v5/stdlib` - PostgreSQL driver

2. **SQL Code Generation**:
   - `sqlc` - Generate type-safe Go code from SQL queries
   - SQL queries stored in `.sql` files

3. **Input Validation**:
   - `github.com/go-playground/validator/v10` - For validating user input

4. **Database Migrations**:
   - `github.com/golang-migrate/migrate/v4` - For versioned database schema migrations

### Local Development Environment

#### Docker & Docker Compose Setup

**docker-compose.yml**:
```yaml
version: '3.8'

services:
  app:
    build: .
    container_name: inventory-cli
    depends_on:
      - db
    environment:
      - DATABASE_URL=postgres://user:password@db:5432/inventory?sslmode=disable
    volumes:
      - .:/app

  db:
    image: postgres:15
    container_name: inventory-db
    environment:
      POSTGRES_DB: inventory
      POSTGRES_USER: user
      POSTGRES_PASSWORD: password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./migrations:/docker-entrypoint-initdb.d

volumes:
  postgres_data:
```

#### Database Migrations

Migrations are managed using `golang-migrate/migrate` with versioned files:
```
migrations/
  ├── 000001_create_tables.up.sql
  ├── 000001_create_tables.down.sql
  ├── 000002_add_indexes.up.sql
  └── 000002_add_indexes.down.sql
```

Migration commands:
```bash
# Apply migrations
migrate -path ./migrations -database "postgres://user:password@localhost:5432/inventory?sslmode=disable" up

# Rollback migrations
migrate -path ./migrations -database "postgres://user:password@localhost:5432/inventory?sslmode=disable" down
```

### Testing Strategy

#### Unit Tests
- Business logic testing (input validation, data transformation)
- Mock database dependencies using interfaces
- Test coverage target: >90%
- Run with: `go test ./...`

#### Integration Tests
- Test database interactions with real PostgreSQL instance
- Run against containerized database via Docker Compose
- Validate SQL queries, migrations, and transactions
- Separate test suite: `go test ./integration/...`

**docker-compose.test.yml** for integration tests:
```yaml
version: '3.8'

services:
  test-db:
    image: postgres:15
    environment:
      POSTGRES_DB: inventory_test
      POSTGRES_USER: test_user
      POSTGRES_PASSWORD: test_password
    ports:
      - "5433:5432"
