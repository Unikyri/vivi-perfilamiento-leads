# Exploration: PostgreSQL Foundation (`issue-3-foundation-postgres`)

## Current State

The project has:
- `migrations/` directory at the repo root (currently only `.gitkeep`).
- `internal/infrastructure/postgres/doc.go` — package stub, no implementation yet.
- `references/embed.go` — established pattern from Issue #2: embed canonical files from a package in the same directory (no `..` traversal).
- `go.mod` declares Go 1.24 with module `github.com/Unikyri/vivi-perfilamiento-leads`; no `go.sum` yet (zero external dependencies).
- Docker v29.6.1 and Docker Compose v5.3.0 are available in the local environment.
- Issue #4 (composition root) expects `postgres.Conectar(ctx, databaseURL) (*pgxpool.Pool, error)` and `postgres.Migrar(ctx, pool) error` — these form the downstream contract.

## The Embedding Problem

Issue #3's sample code proposes:
```go
//go:embed all:../../../migrations/001_esquema_inicial.sql
var esquemaInicial string
```

**This is illegal in Go.** The `//go:embed` directive specification explicitly forbids patterns containing `..` elements, regardless of whether the resolved path stays within the module tree. The Go compiler rejects it at build time with: `pattern may not contain '..'`.

The issue itself acknowledges this: *"Si go:embed con ../../../ falla, copiar el .sql a internal/infrastructure/postgres/migrations/ y ajustar la ruta."*

The core constraint: **exactly one migration source file** must exist in the repository. We must not duplicate the SQL content.

## Affected Areas

- `migrations/001_esquema_inicial.sql` — canonical schema SQL (new file, Contract §5 authority)
- `internal/infrastructure/postgres/` — `migrar.go` with `Conectar` and `Migrar` functions
- `docker-compose.yml` — local Postgres environment (new file at repo root)
- `.env.example` — environment variable documentation (new file)
- `Makefile` — new `db-up`, `db-down`, `db-reset` targets
- `go.mod` / `go.sum` — new dependency `github.com/jackc/pgx/v5`

## Approaches

### 1. **Migrations Embed Package at Repo Root** ✅ RECOMMENDED

- **Description:** Turn `migrations/` into a Go package (add `embed.go`) that embeds the co-located `.sql` files. The `postgres` package imports `migrations` for the raw bytes.
- **Layout:**
  ```
  migrations/
  ├── 001_esquema_inicial.sql   (canonical schema — Contract §5)
  └── embed.go                  (Go package: exposes embedded SQL as []byte or string)
  ```
  ```go
  // migrations/embed.go
  package migrations

  import _ "embed"

  //go:embed 001_esquema_inicial.sql
  var EsquemaInicial string
  ```
  ```go
  // internal/infrastructure/postgres/migrar.go
  package postgres

  import (
      "context"
      "fmt"
      "github.com/jackc/pgx/v5/pgxpool"
      "github.com/Unikyri/vivi-perfilamiento-leads/migrations"
  )

  func Migrar(ctx context.Context, pool *pgxpool.Pool) error {
      _, err := pool.Exec(ctx, migrations.EsquemaInicial)
      if err != nil {
          return fmt.Errorf("aplicando migraciones: %w", err)
      }
      return nil
  }
  ```
- **Pros:**
  - Single source of truth — SQL lives in exactly one place (`migrations/001_esquema_inicial.sql`).
  - Legal Go — embed is in the same directory as the file; no `..` traversal.
  - Follows the pattern established by Issue #2 (`references/embed.go` embeds co-located JSON).
  - Self-contained binary — embed bakes the SQL into the binary (Heroku deploy).
  - Future migrations (002, 003…) add files to the same directory and extend `embed.go`.
  - Clear separation: `migrations/` is the source of truth; `internal/infrastructure/postgres/` is the execution logic.
- **Cons:**
  - `migrations/` becomes a Go package. Minor aesthetic change (same trade-off accepted in `references/`).
  - Infrastructure layer imports `migrations/` package (but `migrations/` has zero logic — it's pure data, not a dependency rule violation).
- **Effort:** Low

### 2. **Copy SQL into `internal/infrastructure/postgres/migrations/`** (Issue's fallback suggestion)

- **Description:** Place the SQL file at `internal/infrastructure/postgres/migrations/001_esquema_inicial.sql` and embed it directly in the `postgres` package.
- **Pros:** Zero import gymnastics; embed is trivially co-located.
- **Cons:**
  - **Violates single-source constraint** — the "canonical" location per the issue and Contract §5 is `migrations/` at the repo root. Having the SQL inside `internal/` makes discovery harder for Block B.
  - Makes the Git diff confusing for PR review (migration lives inside Go infrastructure instead of top-level).
  - Future `psql < migrations/…` verification commands (issue Step 6) would break or require a different path.
- **Effort:** Low
- **Verdict: INFERIOR** — violates discoverability and the issue's own file plan.

### 3. **Read SQL at Runtime from Filesystem (no embed)**

- **Description:** `Migrar` calls `os.ReadFile("migrations/001_esquema_inicial.sql")` at startup.
- **Pros:** No embed complexity. Simple.
- **Cons:**
  - **Binary is not self-contained** — the SQL file must be deployed alongside the binary.
  - Breaks Heroku single-binary deploy model.
  - Fragile path resolution depending on working directory.
  - Violates NFR-P-01 (12-factor: app must run identically in local and Heroku).
- **Effort:** Low
- **Verdict: REJECTED** — fails deployment requirement.

### 4. **Embed from `cmd/servidor/` and Inject**

- **Description:** Embed the SQL in the composition root and pass it down to `Migrar(ctx, pool, sql string)`.
- **Pros:** Keeps `migrations/` non-Go.
- **Cons:**
  - From `cmd/servidor/` the path to `migrations/` is `../../migrations/…` — **still illegal** (`..` forbidden).
  - Even if it worked, it adds coupling: every binary that needs migrations must embed them individually.
- **Effort:** N/A
- **Verdict: REJECTED** — still illegal embed path.

## Recommendation

**Option 1: Migrations Embed Package at `migrations/`.**

This mirrors the exact pattern already established by Issue #2 for `references/`. The team has accepted this trade-off: a directory that contains data files gains an `embed.go` to make them available to Go code. Benefits:

1. ✅ **Single canonical SQL source** — `migrations/001_esquema_inicial.sql`.
2. ✅ **Legal Go embed** — pattern is a local file, no `..` needed.
3. ✅ **Self-contained binary** — SQL baked in at compile time (Heroku requirement).
4. ✅ **Consistent with prior art** — same pattern as `references/embed.go`.
5. ✅ **Future-proof** — additional migrations add files + extend the embed package.
6. ✅ **Issue #4 interface preserved** — `postgres.Migrar(ctx, pool)` remains the exact signature.

## Additional Design Decisions

### pgx/v5 Version Pinning

- Latest stable: **v5.10.0** (released, stable).
- Recommendation: Pin to `github.com/jackc/pgx/v5 v5.10.0` in `go.mod`. Go modules lock the exact version in `go.sum` automatically. Avoid open ranges.
- The project uses `pgxpool` (connection pooling) per the issue and Issue #4's composition root code.

### Schema Exactness (Contract §5 Compliance)

The SQL in the issue body (`CREATE TABLE IF NOT EXISTS …`) matches Contract §5's column definitions. Verification points:
- `leads` table: `version INTEGER NOT NULL DEFAULT 1` supports optimistic locking (NFR-E-01).
- All `TEXT PRIMARY KEY` for IDs (ULID strings ≤ 40 chars, per Contract §0).
- `JSONB` for `perfil`, `capacidad`, `intencion`, `contenido`, `adjunto`.
- `compradores` table: `id INTEGER PRIMARY KEY` (seed data from pipeline, not ULID).
- Indexes: 4 indexes matching the issue's specification exactly.
- `CREATE TABLE IF NOT EXISTS` + `CREATE INDEX IF NOT EXISTS` = **idempotent** (can run N times safely).

### Migration Execution Security / Idempotency

- All DDL uses `IF NOT EXISTS` — running the migration multiple times is safe.
- The migration executes as a single `pool.Exec(ctx, sql)` call — for DDL, PostgreSQL auto-commits each statement.
- Consideration: wrapping in a transaction (`BEGIN; ... COMMIT;`) is unnecessary for DDL with `IF NOT EXISTS` and could cause issues (DDL in transactions is supported in Postgres but adds complexity without benefit for idempotent CREATE).
- Security: the SQL does not accept parameters — no injection risk. It's a compile-time constant embedded string.
- Startup-time execution aligns with RF-M8-02 (Heroku: migrate on deploy, no separate migration step).

### Docker Availability and Local Validation

- ✅ Docker v29.6.1 and Docker Compose v5.3.0 confirmed available.
- The `docker-compose.yml` from the issue uses `postgres:16-alpine` with healthcheck.
- Validation strategy: `docker compose up -d` → wait for healthcheck → apply SQL via psql → verify 7 tables → apply again → verify no errors (idempotency).
- The `Makefile` targets (`db-up`, `db-down`, `db-reset`) provide developer ergonomics.

### Issue #4 Usage Alignment

Issue #4's composition root expects:
```go
pool, err := postgres.Conectar(ctx, cfg.DatabaseURL)
// ...
if err := postgres.Migrar(ctx, pool); err != nil { ... }
```

Our implementation MUST export exactly:
- `func Conectar(ctx context.Context, databaseURL string) (*pgxpool.Pool, error)`
- `func Migrar(ctx context.Context, pool *pgxpool.Pool) error`

The `Conectar` function opens a `pgxpool.Pool` and pings. The `Migrar` function applies the embedded SQL. This interface is fixed by Issue #4's contract.

## Risks

- **Package naming**: `migrations` as a Go package name is clear and unambiguous. No stdlib collision.
- **Embed variable type**: Using `string` (not `[]byte`) for SQL is idiomatic since `pool.Exec` accepts a string query. The issue's sample also uses `string`.
- **Future migration ordering**: If multiple `.sql` files are added, the embed package must expose them in order. For Phase 0 (single file), this is trivial. For future phases, consider `embed.FS` with sorted file listing or numbered variables.
- **No migration tracking table**: The approach relies on `IF NOT EXISTS` idempotency rather than a `schema_migrations` table. This is acceptable for Phase 0 (single migration, hackathon scope). If additive migrations appear later, a tracking mechanism may be needed.
- **Pool sizing**: NFR-E-01 mentions `max_conns = min(dynos×10, plan limit)`. The default `pgxpool` config is adequate for local dev; production tuning is out of scope for this issue.

## Ready for Proposal

**Yes.** All technical questions are resolved:
1. Legal embedding pattern identified (mirrors `references/embed.go`).
2. Interface matches Issue #4's expected API exactly.
3. Schema SQL matches Contract §5.
4. Docker is available for validation.
5. pgx/v5 v5.10.0 is the pinning target.

The proposal phase should formalize the exact file list, embed package structure, Makefile additions, and acceptance criteria.
