# Proposal: PostgreSQL Repositories and ULID ID Generator (Issue #15)

## Intent

The application ports delivered by Issue #11 have no production implementation, so no use case can persist a lead, plan, or ficha, no static catalog reaches the engine, and no component can mint an identifier. This change supplies the infrastructure adapters that make the existing contracts runnable, without altering any public contract, domain type, or database schema.

## Scope

### In Scope
- pgx adapters over the existing `*pgxpool.Pool`: `LeadRepository` (CAS write, filtered/ordered queue, chronological conversation), `PlanRepository` (transactional plan + milestones), `FichaRepository` (upsert on unique `lead_id`).
- Shared JSONB encode/decode helpers and PostgreSQL error translation into `usecase.NotFoundError`, `ErrNoEncontrado`, `ErrOptimisticLock`.
- File-backed `CatalogoRepository`: one-time load and validation of `data/proyectos.json`, `data/compradores.json`, `data/afiliados_mock.json`, `data/brochures/*.md`; defensive copies of maps/slices; explicit not-found for affiliate and brochure misses.
- Concurrency-safe ULID `GeneradorID` in `internal/infrastructure/ids/`: opaque string, unique, ≤ 40 chars.
- Focused unit tests plus opt-in, environment-gated PostgreSQL integration tests; one exactly pinned ULID dependency.

### Out of Scope
- Public `internal/usecase` port, `internal/domain`, `migrations/`, `data/`, and `Docs/` changes.
- Composition-root wiring in `cmd/servidor/main.go`; the registration seam stays a TODO.
- pgx version bump, generic/reflection repository, catalog persisted in PostgreSQL.

## Capabilities

### New Capabilities
- `postgres-repositories`: pgx-backed lead/plan/ficha persistence semantics — compare-and-swap, deterministic ordering, transactional aggregates, error mapping.
- `catalog-cache`: file-backed read-only catalog with one-time load and immutable responses.
- `id-generation`: opaque, unique, concurrency-safe identifier generation.

### Modified Capabilities
- None. `usecase-ports` and `postgres-foundation` requirements stay unchanged.

## Approach

Resource-specific adapters with explicit SQL (exploration approach 1).

- `Guardar` runs `UPDATE ... WHERE lead_id = $1 AND version = $2`; a zero-row result is disambiguated by an existence check so stale writes return `ErrOptimisticLock` and absent leads return not-found, with no storage mutation on conflict. Version increments only after a successful swap.
- `Listar` applies optional `afiliado`/`ruta` predicates conjunctively with `ORDER BY prioridad DESC, lead_id ASC`; conversation reads order by timestamp.
- Plan writes persist base row and milestones in one transaction; `HitosVencidos` joins active plans to pending milestones due on or before the requested date; `MarcarHito` reports absence instead of silently succeeding. No plan CAS is introduced.
- The catalog adapter is constructed from an injected data root / `io/fs` and loads eagerly, so tests are working-directory independent and asset failures surface at startup rather than per request.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `internal/infrastructure/postgres/` | New | Lead, plan, ficha adapters; JSON/error helpers; existing `Conectar`/`Migrar` untouched |
| `internal/infrastructure/postgres/` | New | File-backed catalog adapter with one-time cache |
| `internal/infrastructure/ids/` | New | ULID `GeneradorID` implementation |
| `internal/infrastructure/**/*_test.go` | New | Focused unit tests + gated integration tests |
| `go.mod` / `go.sum` | Modified | One pinned ULID dependency only |
| `data/`, `migrations/`, `internal/usecase`, `internal/domain`, `Docs/` | Unchanged | Read-only or frozen by authority |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Zero-row CAS misclassified as not-found vs stale | High | Explicit existence disambiguation plus dedicated tests for both paths |
| Milestone replacement loses observed IDs/status | Med | Transactional write with defined update semantics decided in design |
| PostgreSQL absent in CI/local | High | Integration tests opt-in and skipped under `-short`; deterministic unit coverage for JSON, catalog, IDs |
| Cached catalog mutated by callers | Med | Return defensive copies; assert immutability in tests |
| Cached first load error becomes permanent | Med | Eager constructor load with explicit initialization error |
| Incidental dependency drift | Low | Pin the ULID module exactly; do not touch the pgx pin |

## Delivery

- Decision needed before apply: Yes
- Chained PRs recommended: Yes
- 400-line budget risk: High (~650–950 authored lines forecast)

Three sequential PRs to `main`, no `size:exception` requested. Each slice is independently deployable because nothing is wired into the composition root yet, ships its own focused tests, records a database/runtime harness result or an explicit `N/A` rationale, and has a rollback boundary limited to its own new files.

| PR | Work unit | Finish state | Rollback boundary |
|----|-----------|--------------|-------------------|
| 1 | JSON/error helpers, ULID generator, `LeadRepository` | Lead CAS, list ordering, conversation, and ID tests pass | New lead/helper/ID files + ULID dependency |
| 2 | `PlanRepository`, `FichaRepository` | Plan/ficha round trips and milestone transitions pass | Plan/ficha adapters and their tests |
| 3 | Catalog cache + optional integration harness | No repeated disk reads; all catalog shapes validated | Catalog adapter, fixtures, and integration tests |

`sdd-tasks` refines this forecast after concrete file sizing.

## Rollback Plan

Revert the offending PR only; earlier slices remain valid because no adapter is referenced by the server or by any use case. No migration runs, no data file changes, and no public contract changes, so revert requires no data repair or coordinated release. If the ULID library proves unsuitable, replace only `internal/infrastructure/ids/` and its `go.mod` entry.

## Dependencies

- Issue #3 (schema and PostgreSQL helpers) — closed, on `origin/main` at `4f30652`.
- Issue #11 (application ports) — closed, on `origin/main` at `ffaf9f1`.
- Optional: a reachable PostgreSQL instance for gated integration tests.

## Success Criteria

- [ ] Every method of `LeadRepository`, `PlanRepository`, `FichaRepository`, `CatalogoRepository`, and `GeneradorID` has a production implementation matching the Issue #11 fake's observable semantics.
- [ ] Stale writes return `ErrOptimisticLock` without mutation; missing records return the repository not-found vocabulary.
- [ ] Queue results are filtered conjunctively and ordered `prioridad DESC, lead_id ASC`.
- [ ] Repeated catalog reads perform no additional disk I/O and cannot mutate the cache (NFR-E-04).
- [ ] Generated IDs are opaque, unique under concurrency, and within the 40-character contract limit.
- [ ] `go build ./...` and `go test ./internal/infrastructure/...` pass; public ports, domain, schema, and `data/` are byte-identical.
- [ ] Delivered as three PRs, each under the 400 authored-line budget.

## Proposal question round

No interactive channel was available to this executor, so these questions and assumptions need user review before `sdd-spec`.

1. Plan write semantics: should `Guardar` fully replace the milestone set (dropping externally created milestones) or merge by milestone ID? Assumption: replace within one transaction, preserving supplied IDs and statuses.
2. Catalog initialization policy: is fail-fast at startup acceptable if any static asset is missing or malformed, or must the service boot degraded? Assumption: fail-fast, eager validation.
3. Affiliate lookup misses: is an explicit not-found error the correct product behavior for an unaffiliated lead, or should it be an empty/neutral result? Assumption: explicit not-found, resolved by the caller.
4. Integration test prerequisites: is a reachable PostgreSQL expected in CI for this hackathon, or should database coverage stay locally opt-in? Assumption: opt-in, skipped under `-short`.
5. ULID source: is an external pinned library acceptable, or is a small standard-library implementation preferred to keep the dependency set minimal? Assumption: pinned maintained library.
