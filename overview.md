# Developer Overview: CLI Inventory Management System

## 1. Introduction & Purpose

Welcome to the CLI Inventory Management System! This document is designed to get you, a new developer, up to speed quickly. It provides a comprehensive overview of the project's architecture, codebase structure, and development workflow.

**What this document is:**
*   A practical guide to understanding how the code is organized.
*   A roadmap for setting up your local development environment.
*   An explanation of the day-to-day development process.

**What this document is not:**
*   A detailed specification of the application's features. For that, please refer to `PROJECT_SPEC.md`, which outlines the "what" and "why" of the application's requirements.

The goal is to equip you with the knowledge needed to navigate the codebase confidently and start making meaningful contributions.

## 2. High-Level Architecture

The application is built with a clean, layered architecture to promote separation of concerns, testability, and maintainability. Understanding these layers is key to finding your way around the code.

Here's a top-down view of the architecture:

```
┌─────────────────────────────────────┐
│             CLI Layer               │
│         (internal/cli/)            │
│  • Handles user input & output     │
│  • Command parsing (Cobra)         │
└─────────────────┬───────────────────┘
                  │
┌─────────────────▼───────────────────┐
│           Service Layer             │
│        (internal/service/)          │
│  • Business logic & rules          │
│  • Orchestrates data flow          │
└─────────────────┬───────────────────┘
                  │
┌─────────────────▼───────────────────┐
│         Repository Layer            │
│       (internal/repository/)        │
│  • Data access logic               │
│  • Interacts with generated DB code│
└─────────────────┬───────────────────┘
                  │
┌─────────────────▼───────────────────┐
│      Database Access Layer          │
│          (internal/db/)             │
│  • Type-safe Go code (sqlc)         │
│  • Direct DB interaction            │
└─────────────────┬───────────────────┘
                  │
┌─────────────────▼───────────────────┐
│          PostgreSQL DB               │
│    (Managed via Docker Compose)     │
└─────────────────────────────────────┘
```

**Layer Responsibilities:**

*   **CLI Layer (`internal/cli/`):** This is the entry point. It uses the [Cobra](https://github.com/spf13/cobra) library to define commands, parse arguments, and display results to the user. It delegates the actual work to the Service Layer.
*   **Service Layer (`internal/service/`):** This is where the core business logic lives. For example, `ProductService` handles creating products and validating SKUs, while `StockService` manages stock movements and ensures atomicity. It calls upon the Repository Layer for data persistence and retrieval.
*   **Repository Layer (`internal/repository/`):** This layer is responsible for abstracting the data source. It implements the interfaces defined by the Service Layer and uses the generated code from the Database Access Layer to interact with PostgreSQL.
*   **Database Access Layer (`internal/db/`):** This directory contains **auto-generated Go code**. The [sqlc](https://sqlc.dev/) tool reads the raw SQL queries in the `queries/` directory and generates type-safe Go functions for executing them. This layer should **never be edited manually**.
*   **PostgreSQL Database:** The application's data is stored in a PostgreSQL database, which is easily managed locally using Docker and Docker Compose.

## 3. Key Technologies & Dependencies

Here's a look at the key technologies and how they fit into your development workflow:

*   **Go 1.24:** The programming language for the entire project.
*   **Cobra (`github.com/spf13/cobra`):** A powerful library for creating modern CLI applications. We use it to define commands, arguments, flags, and help text. You'll be working with it in `internal/cli/`.
*   **pgx/v5 (`github.com/jackc/pgx/v5`):** A robust PostgreSQL driver and toolkit for Go. `sqlc` generates code that uses `pgx/v5` under the hood for efficient and type-safe database interactions.
*   **sqlc:** This is a crucial part of our workflow. It allows us to write SQL queries in `.sql` files (in the `queries/` directory) and automatically generates idiomatic, type-safe Go code (in `internal/db/`). This eliminates boilerplate and reduces the risk of SQL injection errors.
    *   **Workflow:** Edit a `.sql` file in `queries/` -> Run `sqlc generate` -> Use the new/updated functions in `internal/repository/`.
*   **PostgreSQL 15:** Our relational database of choice.
*   **Docker & Docker Compose:** These tools containerize our application and its PostgreSQL database. This ensures a consistent and isolated development environment for everyone. You'll use `docker-compose up` to start the database.
*   **Testify (`github.com/stretchr/testify`):** A popular toolkit with helpful utilities for writing cleaner and more expressive tests in Go.

## 4. Project Structure Deep Dive

Let's take a closer look at the most important directories in the project:

```
.
├── cmd/
│   └── inventory/
│       └── main.go             # Application entry point. Very simple, just calls cli.Execute().
├── internal/
│   ├── cli/                    # Command-Line Interface definitions
│   │   ├── root.go            # Root command setup, service initialization, and DB connection logic.
│   │   ├── product_commands.go # Definitions for product-related commands (add-product, find-product, etc.).
│   │   └── stock_commands.go   # Definitions for stock-related commands (add-stock, move-stock, etc.).
│   ├── service/                # Business Logic Layer
│   │   ├── product.go         # ProductService: handles product creation, retrieval, validation.
│   │   ├── stock.go           # StockService: handles stock logic, movements, transactions.
│   │   └── location.go        # LocationService: handles location management.
│   ├── repository/             # Data Access Layer
│   │   ├── products.go        # ProductRepository: implements data access for products using generated SQLC code.
│   │   ├── stock.go           # StockRepository: implements data access for stock.
│   │   ├── locations.go       # LocationRepository: implements data access for locations.
│   │   └── stock_movements.go # StockMovementRepository: implements data access for audit logs.
│   ├── db/                     # **GENERATED** Database Access Code
│   │   ├── db.go              # Main DB struct and connection setup.
│   │   ├── models.go          # Go structs mirroring database tables (generated by sqlc).
│   │   ├── querier.go         # Interface for all generated queries.
│   │   └── *.sql.go           # Files containing generated Go functions for each .sql file in queries/.
│   ├── models/                 # Core Application Data Models
│   │   ├── product.go         # Defines the `Product` and `CreateProductRequest` structs used throughout the app.
│   │   ├── stock.go           # Defines stock-related models.
│   │   └── location.go        # Defines location-related models.
│   ├── database/               # Database Connection Management
│   │   └── database.go        # Handles the singleton initialization of the pgx connection pool.
│   └── testutils/              # Utilities for Testing
│       ├── test_data.go       # Helper functions to create test data.
│       └── test_database.go   # Helper functions to set up and tear down a test database.
├── queries/                    # Raw SQL Queries for sqlc
│   ├── products.sql           # SQL queries for product operations (CREATE, READ, etc.).
│   ├── stock.sql              # SQL queries for stock operations.
│   ├── locations.sql          # SQL queries for location operations.
│   └── stock_movements.sql    # SQL queries for stock movement logging.
├── migrations/                 # Database Schema Migrations
│   └── 000001_create_tables.up.sql # The initial database schema.
├── Makefile                    # Centralized commands for building, testing, and generating code.
├── sqlc.yaml                   # Configuration file for the sqlc code generator.
├── go.mod                      # Go module definition and dependencies.
└── docker-compose.yml          # Docker configuration for local development (app + db).
```

**Key Takeaways:**
*   When you need to change a command's behavior, start in `internal/cli/`.
*   When you need to change business rules, go to `internal/service/`.
*   When you need to change how data is fetched or saved, look in `internal/repository/`.
*   When you need to change a SQL query, edit the appropriate file in `queries/` and then run `sqlc generate`.
*   **Do not edit files in `internal/db/` directly.**

## 5. Development Workflow

Here is the typical workflow for making a change to the application:

**1. Initial Setup (One-time or when pulling major changes):**
   ```bash
   # 1. Ensure dependencies are installed
   go mod tidy

   # 2. Start the PostgreSQL database in the background
   docker-compose up -d

   # 3. Generate the initial database access code from SQL queries
   sqlc generate
   ```

**2. The Development Loop (Making Changes):**
   Let's say you want to add a new field to a product.

   *   **Step A: Modify the Database Schema (If needed):**
       *   Edit the `queries/products.sql` file to include your new field in the relevant queries (e.g., `CreateProduct`, `GetProductBySKU`).
       *   **Regenerate the DB code:** `sqlc generate`. This will update the `internal/db/` directory with the new/modified functions and models.

   *   **Step B: Update the Application Models:**
       *   Go to `internal/models/product.go`.
       *   Add the new field to the `Product` struct and the `CreateProductRequest` struct.

   *   **Step C: Implement Business Logic (If needed):**
       *   Go to `internal/service/product.go`.
       *   Update the service functions to handle the new field (e.g., validation, passing it to the repository).

   *   **Step D: Update the Data Access Layer:**
       *   Go to `internal/repository/products.go`.
       *   Modify the repository functions to map the new field between your application model (`models.Product`) and the generated database model (`db.Product`).

   *   **Step E: Update the CLI Command:**
       *   Go to `internal/cli/product_commands.go`.
       *   Modify the command (e.g., `addProductCmd`) to accept the new field as an argument and pass it to the service.

   *   **Step F: Test Your Changes:**
       *   Run all tests to ensure you haven't broken anything: `make test-all`.
       *   Write new unit or integration tests for your new functionality.

   *   **Step G: Build and Run:**
       *   Build the application: `make build`.
       *   Test your new feature: `./bin/inventory add-product ...`.

**Useful `Makefile` Commands:**
*   `make generate`: Runs `sqlc generate`.
*   `make build`: Builds the application into `bin/inventory`.
*   `make test`: Runs unit tests.
*   `make integration-test`: Runs integration tests (requires Docker).
*   `make test-all`: Runs both unit and integration tests.
*   `make clean`: Removes generated files (`internal/db/`, `bin/`, coverage files).

## 6. Database Schema Overview

The application uses four main tables to manage inventory. For a detailed, column-by-column breakdown, please refer to the `PROJECT_SPEC.md`. Here's a high-level summary:

*   **`products`**: Stores information about each product, such as its unique SKU, name, description, and price.
*   **`locations`**: Stores the different physical locations or warehouses where inventory can be kept.
*   **`stock`**: This is a junction table that tracks the quantity of a specific product at a specific location. It ensures you know how much of each item is where.
*   **`stock_movements`**: An audit log that records every change in stock levels. It tracks movements *from* a location, *to* a location, the quantity, and the type of movement (e.g., 'add', 'transfer'). This is crucial for traceability.

The relationships are:
*   A `Product` can be in `Stock` at multiple `Locations`.
*   A `Location` can hold `Stock` for multiple `Products`.
*   Each entry in `Stock` is linked to one `Product` and one `Location`.
*   Each `StockMovement` is linked to a `Product` and optionally a `from_location` and/or `to_location`.

Database schema changes are managed via SQL files in the `migrations/` directory.

## 7. How to Get Started: A Quick Checklist

1.  **Prerequisites:** Ensure you have Go 1.24, Docker, and `sqlc` installed.
    ```bash
    # To install sqlc (if you haven't already)
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
    ```
2.  **Clone the Repository:**
    ```bash
    git clone <repository-url>
    cd cli-inventory
    ```
3.  **Install Go Dependencies:**
    ```bash
    go mod tidy
    ```
4.  **Start the Database:**
    ```bash
    docker-compose up -d
    ```
5.  **Generate Database Code:**
    ```bash
    sqlc generate
    ```
6.  **Build the Application:**
    ```bash
    make build
    ```
7.  **Verify it Works:**
    ```bash
    ./bin/inventory list-products
    # You should see a message like "No products found in inventory." or a list of products.
    ```

Congratulations, you're ready to start developing!

## 8. Next Steps for a New Developer

Once you're set up, here are some suggestions to get familiar with the codebase:

1.  **Trace a Command:** Pick a simple command, like `find-product`. Start in `internal/cli/product_commands.go`, follow the call to `productService.GetProductBySKU` in `internal/service/product.go`, and then see how it calls `productRepo.GetBySKU` in `internal/repository/products.go`, which finally uses a function from the generated code in `internal/db/products.sql.go`.
2.  **Add a Simple Feature:** Try adding a new, non-intrusive feature. For example, you could try to add a `--verbose` flag to `list-products` that shows the product description. This will touch the CLI, service, and repository layers.
3.  **Read the Tests:** Look at the existing tests in `internal/service/` and `internal/repository/` (e.g., `product_test.go`). They are great examples of how the different layers are meant to be used and tested.
4.  **Review `PROJECT_SPEC.md`:** Now that you have a handle on the "how," dive into `PROJECT_SPEC.md` to understand the detailed "what" and "why" behind the application's features and intended behavior.

Welcome aboard, and happy coding!
