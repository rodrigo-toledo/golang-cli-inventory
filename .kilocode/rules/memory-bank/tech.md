# Technical Stack and Tooling

## Core Technologies

- **Language**: Go (version 1.25), with the `GOEXPERIMENT=jsonv2` flag enabled.
- **Database**: PostgreSQL (version 17), managed via Docker.
- **CLI Framework**: `cobra` (`github.com/spf13/cobra`) for building the command-line interface.
- **Database Driver**: `pgx/v5` (`github.com/jackc/pgx/v5`) for PostgreSQL interaction.

## Development and Build Tooling

- **Database Code Generation**: `sqlc` (`github.com/sqlc-dev/sqlc`) is used to generate type-safe Go code from raw SQL queries located in the `queries/` directory.
- **Containerization**: Docker and Docker Compose are used to create a reproducible development environment, including the PostgreSQL database.
- **Build & Task Management**: A `Makefile` provides convenient targets for common development tasks like building, testing, and code generation.
- **Dependency Management**: Go Modules (`go.mod`) are used to manage all project dependencies.

## Testing

- **Unit & Integration Testing**: The `testify` library (`github.com/stretchr/testify`) is used for assertions and mocking.
- **Integration Test Environment**: `dockertest` (`github.com/ory/dockertest/v3`) is used to spin up a real PostgreSQL container for integration tests, ensuring that database interactions work as expected.

## Common Developer Commands

The following commands are available via the `Makefile`:

- `make generate`: Regenerates the Go database code from SQL queries using `sqlc`.
- `make build`: Compiles the application into the `bin/` directory.
- `make test-all`: Runs both unit and integration tests.
- `make unit-test`: Runs only the fast unit tests.
- `make integration-test`: Runs the slower integration tests that require a Docker environment.
- `make clean`: Removes all generated files (database code, binaries, coverage reports).