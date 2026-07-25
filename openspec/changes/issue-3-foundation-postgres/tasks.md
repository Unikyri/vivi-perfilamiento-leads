# Tasks: PostgreSQL Foundation (Issue #3)

## Review Workload Forecast

| Field | Value |
|-------|-------|
| Estimated changed lines | 180–230 authored + ~50 go.sum generated |
| 400-line budget risk | Low |
| Chained PRs recommended | No |
| Suggested split | Single PR: feat/foundation-postgres → main |
| Delivery strategy | feature-branch-pr-to-main |
| Chain strategy | size-exception |

Decision needed before apply: No
Chained PRs recommended: No
Chain strategy: size-exception
400-line budget risk: Low

### Suggested Work Units

| Unit | Goal | Likely PR | Focused test command | Runtime harness | Rollback boundary |
|------|------|-----------|----------------------|-----------------|-------------------|
| 1 | Full postgres foundation | Single PR feat/foundation-postgres → main | `go build ./... && go vet ./... && make db-up && make db-validate && make db-down` | Docker compose lifecycle: `make db-up` → `make db-validate` (integration) | All new files + Makefile additions revertible via `git revert`; `docker compose down -v` cleans data |

## Phase 1: Canonical SQL & Embed Package

- [x] 1.1 Create `migrations/001_esquema_inicial.sql` — Contract §5 DDL: 7 tables (`leads`, `mensajes`, `planes`, `hitos`, `fichas`, `compradores`, `demo`), 4 indexes, all `IF NOT EXISTS`
- [x] 1.2 Create `migrations/embed.go` — `package migrations`; `//go:embed 001_esquema_inicial.sql` → `var EsquemaInicial string`
- [x] 1.3 Delete `migrations/.gitkeep` — superseded by real files

## Phase 2: Pinned pgx Pool API

- [x] 2.1 Run `go get github.com/jackc/pgx/v5@v5.10.0` to pin dependency in `go.mod`/`go.sum`
- [x] 2.2 Create `internal/infrastructure/postgres/conectar.go` — `func Conectar(ctx context.Context, databaseURL string) (*pgxpool.Pool, error)` using `pgxpool.New` + `pool.Ping`; wrap errors with `fmt.Errorf`
- [x] 2.3 Create `internal/infrastructure/postgres/migrar.go` — `func Migrar(ctx context.Context, pool *pgxpool.Pool) error` importing `migrations.EsquemaInicial`; execute via `pool.Exec`

## Phase 3: Docker Compose / Environment / Make

- [x] 3.1 Create `docker-compose.yml` at repo root — `postgres:16-alpine`, healthcheck (`pg_isready`), port 5432, env vars `POSTGRES_DB=vivi`, `POSTGRES_USER=vivi`, `POSTGRES_PASSWORD=vivi`
- [x] 3.2 Create `.env.example` — `DATABASE_URL=postgres://vivi:vivi@localhost:5432/vivi?sslmode=disable`
- [x] 3.3 Add Makefile targets `db-up`, `db-down`, `db-reset`, `db-validate` to `.PHONY` list and implement: `db-up` = compose up + wait healthy; `db-down` = compose down -v; `db-reset` = db-down + db-up; `db-validate` = psql migration apply + assert 7 tables

## Phase 4: Integration Validation (Docker Idempotency)

- [x] 4.1 Run `go build ./...` — verify pgx import compiles with embed
- [x] 4.2 Run `go vet ./...` — verify clean static analysis
- [x] 4.3 Run `make db-up` — container healthy within 10s
- [x] 4.4 Run `make db-validate` — migration applied, 7 tables confirmed
- [x] 4.5 Run `make db-validate` again — idempotency: second run exits 0, no error
- [x] 4.6 Run `make db-reset` then `make db-validate` — clean slate re-validates
- [x] 4.7 Run `make db-down` — teardown exits 0, no lingering volumes

## Phase 5: Cleanup & PR Delivery

- [x] 5.1 Run `go mod tidy` — ensure `go.sum` is complete and minimal
- [x] 5.2 Verify no modifications to `Docs/` or `README.md` (excluded paths per spec)
- [ ] 5.3 Commit all changes on `feat/foundation-postgres` branch with descriptive message
- [ ] 5.4 Push branch and open PR targeting `main`
