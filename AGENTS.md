# Repository Guidelines

## Project Structure & Module Organization
- `cmd/notebook/` — application entry point (`main.go`).
- `internal/` — packages by concern: `config/`, `content/`, `http/`, `store/`, `feed/`.
- `content/` — Markdown posts (`YYYY-MM-DD-slug.md`).
- `migrations/` — SQL schema (initialized programmatically in `store/sqlite.go`).
- `themes/`, `docs/` — assets and documentation.

## Build, Test, and Development Commands
- `go mod tidy` — sync dependencies.
- `go build -o notebook ./cmd/notebook` — build binary.
- `ENV=dev go run ./cmd/notebook` — run locally (shows drafts).
- `go test ./...` — run all tests; add `-v -cover` for detail/coverage.
- `go fmt ./... && go vet ./...` — format and static checks.

## Coding Style & Naming Conventions
- Language: Go 1.22.x. Use `gofmt` (tabs, standard formatting).
- Packages: short, lower-case names (`content`, `store`).
- Identifiers: `CamelCase` for exported, `lowerCamel` for unexported.
- Files: use underscores and clear intent (`loader_test.go`, `sqlite.go`).
- Imports: standard → third-party → internal; keep unused imports out.

## Testing Guidelines
- Framework: standard `testing` package; table-driven tests preferred.
- Naming: files end with `_test.go`; tests `func TestXxx(t *testing.T)`.
- Scope: cover parsing (`internal/content`), HTTP handlers/middleware, feeds, and store.
- Run examples: `go test -v ./internal/http`, `go test -cover ./internal/content`.

## Commit & Pull Request Guidelines
- Conventional style preferred: `feat:`, `fix:`, `docs:`, `refactor:`, `test:`, `chore:`; optional scope: `feat(http): ...`.
- Subject in imperative, ≤72 chars; include a concise body when relevant.
- PRs: describe intent, link issues, list affected packages, include screenshots for UI/HTML changes, and note test coverage/impact.

## Security & Configuration Tips
- Config via env: `ENV`, `PORT`, `DB_PATH`, `CONTENT_DIR`, `SITE_BASEURL`, `SITE_TITLE` (see `internal/config`).
- Do not hardcode secrets; avoid committing local DBs or large binaries.
- `ENV=dev` exposes drafts; ensure `ENV=prod` for deployments.

## Agent-Specific Instructions
- Keep changes minimal and scoped; align with `README.md` and `ARCHITECTURE.md`.
- Don’t introduce new dependencies without discussion; prefer standard library.
- If you change routes, data models, or env vars, update tests and docs.
