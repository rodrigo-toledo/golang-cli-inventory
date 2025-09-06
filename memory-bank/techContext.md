# Technical Context

Language & Runtime
- Go 1.24 (module-based project).

Key libraries
- CLI: github.com/spf13/cobra
- Database driver: github.com/jackc/pgx/v5 (pgxpool for connection pooling)
- SQL codegen: sqlc (github.com/sqlc-dev/sqlc)
- Testing: github.com/stretchr/testify
- Integration test helper: github.com/ory/dockertest/v3

Database
- PostgreSQL 15 (development via Docker Compose).
- Migrations executed by Postgres container using files in `migrations/` (mounted to `/docker-entrypoint-initdb.d`).

Project layout (important files)
- cmd/inventory/main.go — application entrypoint.
- internal/cli/ — Cobra commands and CLI wiring.
- internal/service/ — business logic & validation.
- internal/repository/ — data access, mapping to internal models.
- internal/db/ — sqlc-generated code (do not edit).
- queries/ — SQL query source for sqlc.
- migrations/ — SQL migrations applied to DB initialization.
- Makefile — development commands (generate, build, test, integration-test).

Developer setup (common commands)
- Ensure Go and Docker are installed.
- Install sqlc:
  - go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
- Fetch dependencies:
  - go mod tidy
- Start local DB:
  - docker-compose up -d
- Generate db code from SQL:
  - sqlc generate
- Build binary:
  - make build
- Run CLI:
  - ./bin/inventory <command>
- Tests:
  - Unit tests: go test ./...
  - Integration tests (requires Docker): make integration-test
  - All tests: make test-all

Code generation & workflow rules
- SQL in `queries/` is authoritative; edit SQL and run `sqlc generate` to update `internal/db/`.
- Never edit files in `internal/db/` — they are generated.
- Repositories map between `internal/db` types and `internal/models` types.
- For schema changes that require migrations: add a new migration file in `migrations/` and reinitialize DB (or manage volumes) to apply.

Notes & constraints
- Keep CLI free of heavy business logic; validation and rules live in services.
- Use repository interfaces in services to make unit testing easier (mocks).
- Integration tests may rely on Docker Compose or dockertest. Use `docker-compose.test.yml` or the project's test utilities where available.
