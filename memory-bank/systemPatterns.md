# System Patterns

High-level architecture
- Layered design: CLI -> Service -> Repository -> DB (sqlc-generated) -> PostgreSQL.
- Each layer has a single responsibility and communicates with the adjacent layer via well-defined interfaces.

Common patterns
- Dependency Injection via constructor functions: services receive repository interfaces; repositories receive `*db.Queries` / DB connection.
- Repository pattern: repository implementations map between application models (`internal/models`) and generated DB models (`internal/db`).
- Thin CLI layer: Cobra commands parse arguments/flags and delegate to service layer; CLI contains minimal business logic.
- Transactional operations: multi-step stock operations (e.g., move-stock) execute within a DB transaction to ensure atomicity and record stock_movements.
- Error handling: services translate low-level DB errors into domain errors with clear messages for CLI output; use of wrapped errors for traceability.

Database & sqlc rules
- SQL in `queries/` is the source of truth for DB interactions.
- `sqlc generate` produces `internal/db/` code and models — never edit generated files.
- For schema changes: add/modify SQL in `queries/` and/or `migrations/`, then run `sqlc generate` and update mapping code in repositories.
- Keep queries explicit and tested; avoid embedding dynamic SQL in Go where possible.

Testing patterns
- Unit tests: mock repositories (interfaces) to isolate services; use `testify` assertions.
- Integration tests: run against a real Postgres instance (Docker Compose or dockertest) to validate sqlc-generated code, transactions, and migrations.
- Test utilities: `internal/testutils` contains helpers for creating test data and managing test database lifecycle.

Mapping & DTO strategy
- Application models (`internal/models`) are used across service and CLI layers.
- Repository implementations translate between `internal/models` and `internal/db` structs returned by sqlc.
- Keep mapping logic centralized in repository/mapper.go when applicable.

Operational notes
- Use Makefile targets for common tasks: `make generate` (sqlc), `make build`, `make test`, `make integration-test`.
- Database initialization in development is handled via Docker Compose mounting `migrations/` to the Postgres container.
