# PostgreSQL Foundation Specification

## Purpose

Contract §5 schema as idempotent embedded DDL, pgx v5.10.0 pool connection, Docker Compose local environment, and Make targets — exposing `Conectar`/`Migrar` for downstream composition (Issue #4). Single-source migration SQL embedded via legal Go package at `migrations/`.

## Requirements

### Requirement: Schema DDL (7 Tables, 4 Indexes)

The migration SQL file `migrations/001_esquema_inicial.sql` MUST define exactly these 7 tables with `CREATE TABLE IF NOT EXISTS`:

| Table | Primary Key | Notable Columns |
|-------|-------------|-----------------|
| `leads` | `lead_id TEXT` | perfil/capacidad/intencion JSONB, version INT NOT NULL DEFAULT 1, creado_en/actualizado_en TIMESTAMPTZ |
| `mensajes` | `mensaje_id TEXT` | lead_id TEXT FK, autor TEXT, tipo_contenido TEXT, adjunto JSONB, creado_en TIMESTAMPTZ |
| `planes` | `plan_id TEXT` | lead_id TEXT FK, estado TEXT, meta_monto BIGINT, consentimiento_en TIMESTAMPTZ |
| `hitos` | `hito_id TEXT` | plan_id TEXT FK, tipo TEXT, fecha DATE, monto BIGINT, estado TEXT |
| `fichas` | `ficha_id TEXT` | lead_id TEXT FK UNIQUE, contenido JSONB, generada_en TIMESTAMPTZ |
| `compradores` | `id INTEGER` | All §4.1 columns; seed-loaded by pipeline |
| `demo` | `clave TEXT` | valor TEXT |

The file MUST define exactly 4 indexes with `CREATE INDEX IF NOT EXISTS`:

| Index | Table | Column(s) |
|-------|-------|-----------|
| `idx_mensajes_lead_id` | mensajes | lead_id |
| `idx_planes_lead_id` | planes | lead_id |
| `idx_hitos_plan_id` | hitos | plan_id |
| `idx_compradores_proyecto_id` | compradores | proyecto_id |

The DDL MUST NOT use any migration tracking table. Column definitions MUST match Contract §5 exactly.

#### Scenario: All 7 tables created on fresh database

- GIVEN a PostgreSQL 16 instance with no application tables
- WHEN the migration SQL is executed
- THEN exactly 7 tables MUST exist: leads, mensajes, planes, hitos, fichas, compradores, demo
- AND exactly 4 indexes MUST exist matching the specification above

#### Scenario: Idempotent re-execution produces no error

- GIVEN the migration has already been applied once successfully
- WHEN the same SQL is executed a second time
- THEN no error SHALL occur
- AND table/index definitions MUST remain unchanged

### Requirement: Legal Single-Source Embed Package

`migrations/embed.go` MUST declare `package migrations` and expose the SQL via `//go:embed 001_esquema_inicial.sql` as `var EsquemaInicial string`. The embed directive MUST NOT use `..` path traversal. No other copy of the schema SQL SHALL exist in the repository.

#### Scenario: Embed compiles without illegal paths

- GIVEN `embed.go` is co-located with `001_esquema_inicial.sql` in `migrations/`
- WHEN `go build ./...` is executed
- THEN compilation MUST succeed with no embed path errors

#### Scenario: Single canonical source

- GIVEN the full repository tree
- WHEN searched for `CREATE TABLE` statements defining the schema
- THEN only `migrations/001_esquema_inicial.sql` SHALL contain them

### Requirement: Pinned pgx Connection API

`internal/infrastructure/postgres/conectar.go` MUST export `func Conectar(ctx context.Context, databaseURL string) (*pgxpool.Pool, error)`. It MUST use `pgxpool.New` with the provided URL and validate connectivity via `pool.Ping(ctx)`. `go.mod` MUST pin `github.com/jackc/pgx/v5` at exactly `v5.10.0`.

#### Scenario: Successful connection to running Postgres

- GIVEN a running PostgreSQL instance at `DATABASE_URL`
- WHEN `Conectar(ctx, databaseURL)` is called
- THEN it MUST return a non-nil `*pgxpool.Pool` and nil error

#### Scenario: Invalid URL returns wrapped error

- GIVEN a malformed or unreachable `DATABASE_URL`
- WHEN `Conectar(ctx, databaseURL)` is called
- THEN it MUST return nil pool and a non-nil error wrapping the cause

### Requirement: Migration Executor

`internal/infrastructure/postgres/migrar.go` MUST export `func Migrar(ctx context.Context, pool *pgxpool.Pool) error`. It MUST import `migrations.EsquemaInicial` and execute it via `pool.Exec(ctx, sql)`. On success it MUST return nil; on failure it MUST return an error wrapping the cause with `fmt.Errorf`.

#### Scenario: Migrar applies schema successfully

- GIVEN a valid pool connected to a fresh database
- WHEN `Migrar(ctx, pool)` is called
- THEN it MUST return nil and all 7 tables MUST exist

#### Scenario: Migrar is idempotent

- GIVEN `Migrar` has already succeeded once
- WHEN `Migrar(ctx, pool)` is called again
- THEN it MUST return nil without error

### Requirement: Docker Local Environment

A `docker-compose.yml` at repo root MUST define a `postgres` service using image `postgres:16-alpine` with a healthcheck (`pg_isready`). Environment MUST set `POSTGRES_DB`, `POSTGRES_USER`, `POSTGRES_PASSWORD`. Port 5432 MUST be mapped to host.

#### Scenario: Container starts and passes healthcheck

- GIVEN Docker Compose is available
- WHEN `docker compose up -d` is executed
- THEN the postgres container MUST reach healthy state within 10 seconds

### Requirement: Environment Variable Template

`.env.example` MUST contain at minimum `DATABASE_URL=postgres://vivi:vivi@localhost:5432/vivi?sslmode=disable`. It SHALL serve as documentation for NFR-P-01 12-factor configuration.

#### Scenario: .env.example documents required vars

- GIVEN the file `.env.example` exists at repo root
- WHEN its content is read
- THEN it MUST contain a valid `DATABASE_URL` line

### Requirement: Makefile Targets

The root `Makefile` MUST add targets `db-up`, `db-down`, `db-reset`, and `db-validate` to the existing `.PHONY` list. Semantics:

| Target | Action |
|--------|--------|
| `db-up` | `docker compose up -d` and wait for healthy |
| `db-down` | `docker compose down -v` |
| `db-reset` | `db-down` then `db-up` |
| `db-validate` | Apply migration via psql, verify 7 tables exist |

#### Scenario: Full lifecycle via Make

- GIVEN Docker is available and no prior containers
- WHEN `make db-up` then `make db-validate` then `make db-reset` then `make db-down` are executed in sequence
- THEN each MUST exit with status 0
- AND after `db-validate`, exactly 7 application tables MUST exist

### Requirement: Excluded Paths

This change MUST NOT modify any file under `Docs/` or the root `README.md`.

#### Scenario: Protected paths unchanged

- GIVEN the repository state after this change
- WHEN `Docs/` and `README.md` are diffed against baseline
- THEN there SHALL be zero modifications
