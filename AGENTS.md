# Repository Guidelines

## Project Structure & Module Organization
- `cmd/` contains service entry points (`go install ./cmd/...` builds binaries).
- `core/` holds domain and shared business logic.
- `internal/` contains application-specific code and wiring.
- `vendor/` is committed; builds/tests must use vendored deps.
- `bin/` is the default output location for installed binaries.
- Tests live alongside code as `*_test.go` files.

## Build, Test, and Development Commands
- `make fmt` — format all Go code with `gofumpt`.
- `make vet` — run `go vet`, `staticcheck`, and `shadow` (includes formatting).
- `make lint` — run additional linters (`deadcode`, `golangci-lint`).
- `make build` — build/install all commands in `cmd/`.
- `make test` — run `go test` with race and coverage, writing `cover.out`.
- `make cover` — open coverage HTML from `cover.out`.

## Coding Style & Naming Conventions
- Use standard Go naming: exported `CamelCase`, unexported `camelCase`.
- Add doc comments for all exported identifiers.
- Follow `gofumpt` formatting and lint rules (`wsl_v5`, `nlreturn`).
  - Leave a blank line before `return` when there are multiple statements.
  - Separate consecutive `if` blocks with a blank line.
- Keep package names short and lowercase.

## Testing Guidelines
- Use Go’s testing package; prefer table-driven tests.
- Name test files `*_test.go` and test functions `TestXxx`.
- Run tests via `make test` (uses `-mod=vendor`, race, and coverage).

## Commit & Pull Request Guidelines
- Commit messages are short and imperative, often with a scope prefix.
  - Examples: `pinger: fix data race`, `sqlite: remove ID from endpoints table`, `chore: update example.png`.
- PRs should include a clear description, linked issues (if any), and test results.
- Include screenshots only when UI or visual output changes.

## Agent-Specific Notes
- This repo expects vendored builds: `GOFLAGS=-mod=vendor -race`.
- Do not add tests unless explicitly requested.
