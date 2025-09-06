# CRUSH.md

## Commands
- Build: `make build` (JSON v2 enabled)
- Lint: `golangci-lint run`
- Test all: `make test`
- Test single: `GOEXPERIMENT=jsonv2 go test ./internal/<package> -run ^TestName$`
- Coverage: `make test-coverage`
- Generate mocks: `make mocks`
- Generate SQL code: `make generate`

## Code Style
- **Imports**: Grouped (stdlib, project, third-party), sorted via `goimports`
- **Formatting**: `gofmt -s -w`, no manual line breaks
- **Naming**:
  - Functions/Types: `CamelCase` (exported), `camelCase` (unexported)
  - Variables: `snake_case` for constants, `camelCase` elsewhere
- **Errors**: Immediate check + `if err != nil`, wrap with `%w`
- **Tests**:
  - Table-driven with subtests (`t.Run`)
  - Generated mocks via mockery in `internal/mocks/`
  - Use `testify/assert`

## Notes
- Always use `GOEXPERIMENT=jsonv2` in test commands
- `.crush/` already in `.gitignore`
- Mocks are generated using mockery v3 and stored in `internal/mocks/`