# Project Brief

Project: CLI Inventory Management System

Purpose:
- Provide a small, maintainable command-line application to manage products, locations, stock, and stock movements with a clear audit log.
- Serve as both a real tool and a learning project that demonstrates layered architecture, testing patterns, and use of sqlc for type-safe DB access.

Primary Goals:
- Reliable CLI commands for common inventory operations: add-product, add-stock, find-product, move-stock, generate-report.
- Strong separation of concerns (CLI, Service, Repository, DB layers).
- High test coverage with unit and integration tests.
- Simple developer experience: Docker for local DB, sqlc for DB code generation, Makefile shortcuts.

Scope / Non-goals:
- Scope: Local developer-focused CLI and integration tests; schema migrations handled via SQL files in `migrations/`.
- Non-goals: Production deployment orchestration and advanced migration tooling (e.g., golang-migrate) are out of scope for the initial project.

Success Criteria:
- All CLI commands implemented and covered by tests.
- DB interactions implemented via sqlc-generated code in `internal/db/` (never edited manually).
- Developer setup documented and reproducible using docker-compose and `make` targets.
