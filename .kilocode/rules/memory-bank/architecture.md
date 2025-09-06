# System Architecture

## Layered Architecture

The application follows a clean, layered architecture to ensure separation of concerns and testability.

```
┌─────────────────────────────────────┐
│             CLI Layer               │
│         (internal/cli/)             │
└─────────────────┬───────────────────┘
                  │
┌─────────────────▼───────────────────┐
│           Service Layer             │
│        (internal/service/)          │
└─────────────────┬───────────────────┘
                  │
┌─────────────────▼───────────────────┐
│         Repository Layer            │
│       (internal/repository/)        │
└─────────────────┬───────────────────┘
                  │
┌─────────────────▼───────────────────┐
│      Database Access Layer (sqlc)   │
│          (internal/db/)             │
└─────────────────┬───────────────────┘
                  │
┌─────────────────▼───────────────────┐
│          PostgreSQL DB              │
└─────────────────────────────────────┘
```

- **CLI Layer (`internal/cli/`)**: Handles user input and output using the Cobra library.
- **Service Layer (`internal/service/`)**: Contains the core business logic and orchestrates data flow.
- **Repository Layer (`internal/repository/`)**: Abstracts data access and interacts with the generated database code.
- **Database Access Layer (`internal/db/`)**: Contains auto-generated, type-safe Go code from `sqlc`. **This layer should not be edited manually.**

## Key Technical Decisions

- **Database Access**: `sqlc` is used to generate type-safe Go code from raw SQL queries located in the `queries/` directory. This approach minimizes boilerplate and prevents SQL injection vulnerabilities.
- **Dependency Management**: Go Modules (`go.mod`) are used for managing project dependencies.
- **Containerization**: Docker and Docker Compose are used to create a reproducible development environment, including the PostgreSQL database.

## Design Patterns

- **Repository Pattern**: The repository layer isolates the data access logic, making it easier to manage and test.
- **Service Layer Pattern**: The service layer encapsulates the business logic, promoting a clean separation from the CLI and data access layers.
- **Dependency Injection**: Dependencies (like services and repositories) are initialized in `cmd/inventory/main.go` and passed down to the components that need them, which facilitates testing and modularity.

## Critical Implementation Paths

- **Adding a new command**:
  1.  Define the command in `internal/cli/`.
  2.  Add business logic to the appropriate service in `internal/service/`.
  3.  Implement data access in the corresponding repository in `internal/repository/`.
  4.  If necessary, add or modify SQL queries in `queries/` and regenerate the database code with `sqlc generate`.

- **Modifying the database schema**:
  1.  Create a new SQL migration file in the `migrations/` directory.
  2.  Update the relevant SQL queries in the `queries/` directory.
  3.  Regenerate the database code using `sqlc generate`.
  4.  Update the repository and service layers to reflect the changes.