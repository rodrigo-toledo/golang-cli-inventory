# CLI Inventory Management System - Project Specification

## 1. Core Features / Commands

### add-product
**Purpose**: Add a new product definition to the inventory system.

**Arguments (Positional)**:
1.  `sku` (required): Stock Keeping Unit - unique product identifier.
2.  `name` (required): Product name.
3.  `description` (optional): Product description.
4.  `price` (required): Product price (must be a valid number).

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

**Arguments (Positional)**:
1.  `product-id` (required): The ID of the product to add stock for.
2.  `location-id` (required): The ID of the location where stock will be added.
3.  `quantity` (required): The quantity of stock to add (must be a positive integer).

**Expected Success Output**:
```
Successfully added 50 units of product (ID: 1) to location (ID: 1)
```

**Key Error Scenarios**:
- Product with given SKU does not exist
- Invalid quantity (negative or zero)
- Location identifier is invalid
- Insufficient database permissions

### find-product
**Purpose**: Find a product by its unique SKU.

**Arguments (Positional)**:
1.  `sku` (required): The exact SKU of the product to find.

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

**Arguments (Positional)**:
1.  `product-id` (required): The ID of the product to move.
2.  `from-location-id` (required): The ID of the source location.
3.  `to-location-id` (required): The ID of the destination location.
4.  `quantity` (required): The quantity of stock to move (must be a positive integer).

**Expected Success Output**:
```
Successfully moved 10 units of product (ID: 1) from location (ID: 1) to location (ID: 2)
```

**Key Error Scenarios**:
- Insufficient stock at source location
- Product not found
- Source and destination locations are the same
- Database transaction failure

### generate-report
**Purpose**: Generate various inventory reports.

**Arguments (Positional)**:
1.  `report-type` (required): The type of report to generate. Currently supports "low-stock".
2.  `threshold` (optional): The minimum stock threshold for the "low-stock" report. Defaults to 10 if not provided.

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

The following schema represents the current state of the database as defined in the initial migration file (`migrations/000001_create_tables.up.sql`).

```sql
CREATE TABLE locations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    sku VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE stock (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    location_id INTEGER NOT NULL REFERENCES locations(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(product_id, location_id)
);

CREATE TABLE stock_movements (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    from_location_id INTEGER REFERENCES locations(id) ON DELETE SET NULL,
    to_location_id INTEGER REFERENCES locations(id) ON DELETE SET NULL,
    quantity INTEGER NOT NULL,
    movement_type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```
**Note:** The indexes mentioned in previous versions of this spec (e.g., `idx_products_sku`) are not present in the initial migration. They may be added in future migrations if performance analysis indicates a need.

## 3. Technical Stack & Local Environment

### Language
Go (Golang)

### Database
PostgreSQL

### Key Go Libraries & Tools

1.  **CLI Framework**:
    *   `github.com/spf13/cobra` - A powerful library for creating modern command-line applications.

2.  **Database Access**:
    *   `github.com/jackc/pgx/v5` - A robust PostgreSQL driver and toolkit. The `pgxpool` is used for connection pooling.

3.  **SQL Code Generation**:
    *   `sqlc` - A tool for generating type-safe Go code from SQL queries. This eliminates boilerplate and reduces the risk of SQL injection.

4.  **Testing**:
    *   `github.com/stretchr/testify` - A popular toolkit with helpers for writing cleaner tests (e.g., assertions, mock suites).
    *   `github.com/ory/dockertest/v3` - A library for managing Docker containers during integration tests.

**Note:** While `github.com/go-playground/validator/v10` was previously considered for input validation, it is not currently a direct dependency or actively used in the project. Validation is currently handled manually within the service layer.

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

Database schema changes are managed using SQL migration files located in the `migrations/` directory. These files are executed automatically when the PostgreSQL container is first created using Docker Compose.

The `docker-compose.yml` file mounts the `migrations/` directory to `/docker-entrypoint-initdb.d` inside the PostgreSQL container. Any `.sql` files found in this directory will be executed in alphabetical order when the container initializes its database.

**Migration File Naming Convention:**
```
migrations/
  └── 000001_create_tables.up.sql
```
*   `000001`: A version number to ensure ordering.
*   `create_tables`: A descriptive name for the migration.
*   `up.sql`: Indicates this file applies the migration (creates tables, adds columns, etc.). Corresponding `.down.sql` files for rollbacks are planned but not yet implemented in the current workflow.

**Process for Adding a New Migration:**
1.  Create a new `.sql` file in the `migrations/` directory following the naming convention (e.g., `000002_add_new_column.up.sql`).
2.  Write the SQL schema changes (e.g., `ALTER TABLE ...`) in this new file.
3.  To apply the migration to a fresh local database:
    *   Stop and remove the existing Docker volume: `docker-compose down -v`.
    *   Restart the database: `docker-compose up -d`. The new migration will be executed automatically.

**Note:** This current setup is ideal for development and initial setup. For production environments or more complex schema versioning and rollback scenarios, a dedicated migration tool like `golang-migrate` or `Flyway` would be recommended.

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
