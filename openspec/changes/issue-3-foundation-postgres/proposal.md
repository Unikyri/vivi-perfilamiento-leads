# Proposal: PostgreSQL Foundation (Issue #3)

## Intent

Deliver the Contract §5 PostgreSQL schema, a legal embedded migration package, pgx v5.10.0 pool connection, Docker Compose local environment, and Make targets — exposing `Conectar` / `Migrar` for Issue #4's composition root.

## Scope

### In Scope
- `migrations/embed.go` + `001_esquema_inicial.sql` (7 tables, 4 indexes, idempotent DDL)
- `internal/infrastructure/postgres/conectar.go` — `Conectar(ctx, url) (*pgxpool.Pool, error)`
- `internal/infrastructure/postgres/migrar.go` — `Migrar(ctx, pool) error`
- `docker-compose.yml` (postgres:16-alpine with healthcheck)
- `.env.example` with `DATABASE_URL`
- Makefile targets: `db-up`, `db-down`, `db-reset`, `db-validate`
- `go.mod` pinned `github.com/jackc/pgx/v5 v5.10.0`

### Out of Scope
- Migration tracking table (IF NOT EXISTS idempotency suffices for Phase 0)
- Production pool tuning (NFR-E-01 max_conns deferred to deploy issue)
- Seed data / pipeline load (Issue #4+)
- `Docs/` and `README.md`

## Capabilities

### New Capabilities
- `postgres-foundation`: schema DDL, embedded migrations, pool connection, local Docker environment

### Modified Capabilities
None

## Approach

1. **Migrations embed package** at `migrations/` (mirrors `references/embed.go` pattern): `embed.go` exposes `EsquemaInicial string` from co-located SQL — legal Go embed, single canonical source.
2. **pgx v5.10.0** pinned in `go.mod`; `pgxpool.New` for connection, `.Ping` for validation.
3. **Idempotent DDL** — all `CREATE TABLE IF NOT EXISTS` / `CREATE INDEX IF NOT EXISTS`; single `pool.Exec(ctx, sql)` call; no transaction wrapper needed.
4. **Docker Compose** with healthcheck; Make targets wrap compose lifecycle and psql validation.
5. **Interface contract**: exact signatures `Conectar` and `Migrar` preserved for Issue #4.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `migrations/embed.go` | New | Go embed package exposing SQL |
| `migrations/001_esquema_inicial.sql` | New | Contract §5 schema (7 tables, 4 indexes) |
| `internal/infrastructure/postgres/conectar.go` | New | `Conectar` pool factory |
| `internal/infrastructure/postgres/migrar.go` | New | `Migrar` executor |
| `docker-compose.yml` | New | Local Postgres service |
| `.env.example` | New | Environment variable template |
| `Makefile` | Modified | Add db-up, db-down, db-reset, db-validate |
| `go.mod` / `go.sum` | Modified | Add pgx v5.10.0 |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| pgx v5.10.0 has breaking changes vs future upgrades | Low | Pinned; upgrade is a separate issue |
| Docker unavailable in CI (Heroku) | Low | `db-validate` is local-only; CI uses `go build` + `go vet` |
| Multi-migration ordering in future phases | Med | Design for `embed.FS` upgrade path; document in embed.go |

## Rollback Plan

1. `git revert` the merge commit on `main`.
2. `docker compose down -v` removes local volumes.
3. No production data exists at this stage — rollback is purely code removal.
4. `go.mod` returns to zero-dependency state via revert.

## Dependencies

- Docker + Docker Compose available locally (confirmed v29.6.1 / v5.3.0)
- Issue #2 completed (repo structure, references embed pattern established)
- Contract §5 schema definition (authoritative, no modifications)

## Success Criteria

- [ ] `go build ./...` passes with pgx v5.10.0 imported
- [ ] `go vet ./...` clean
- [ ] `make db-up` starts Postgres, healthcheck passes within 10s
- [ ] `make db-validate` applies migration via psql, verifies 7 tables exist
- [ ] `make db-reset` drops and recreates cleanly
- [ ] Migration is idempotent: second `Migrar` call produces no error
- [ ] `postgres.Conectar` and `postgres.Migrar` signatures match Issue #4 expectations exactly
