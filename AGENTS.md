AGENTS.md

Build / Lint / Test
- Build: go build ./...
- Run all tests: go test ./... -v
- Run single package tests: go test ./internal/handlers -v
- Run single test: go test ./internal/auth -run TestAuthenticate -v
- Lint & format: go fmt ./... ; golangci-lint run (if installed)

Code style guidelines
- Formatting: use go fmt/gofmt for all files; run `go fmt ./...` before commits.
- Imports: group stdlib, third-party, then internal packages; use goimports to auto-fix ordering.
- Types & naming: use CamelCase for exported identifiers, mixedCaps for unexported. Prefer small, focused types; keep functions short and single-responsibility.
- Errors: return errors rather than logging in libraries; wrap with fmt.Errorf("context: %w", err) when adding context. Use sentinel errors only when necessary.
- Context: pass context.Context as first parameter to public functions that may be cancellable or I/O bound.
- Testing: table-driven tests encouraged; use testutils helpers where present. Keep tests deterministic, avoid network/file system reliance unless using test_database utilities.
- Mocks: use generated mocks in internal/mocks; prefer interfaces from service/ and repository/ to make units testable.
- Logging: handlers may use error_handler.go patterns; avoid global state.

Repository rules
- Cursor/Copilot: include project-specific assistant rules if present
  - .cursor/rules/ (if exists) — follow those cursor rules
  - .github/copilot-instructions.md (if exists) — follow Copilot instructions

Commit & PR
- Run `go test ./...` and `golangci-lint run` locally before committing. Keep commits small and focused.

Notes
- This file is for automated agents operating in this repo: always run tests after changes, prefer edits to existing files, and avoid creating new top-level docs unless requested.
