## Exploration: PostgreSQL repositories and ULID ID generator

### Current State

Issue #3 is already satisfied: `internal/infrastructure/postgres/` contains `Conectar`, which creates and pings a `*pgxpool.Pool`, and `Migrar`, which executes the embedded canonical SQL from `migrations/001_esquema_inicial.sql` through `migrations.EsquemaInicial`. The migration defines the seven Contract v1.1 tables and the four required indexes; this change must not modify that schema or migration helper.

Issue #11 is already satisfied: `internal/usecase/puertos.go` owns the application ports and imports only the standard library plus `internal/domain`. The relevant contracts are `LeadRepository`, `PlanRepository`, `FichaRepository`, `CatalogoRepository`, and `GeneradorID`. The fake implementation and tests establish important behavioral semantics that production adapters must match: lead creation normalizes version 0 to 1, `Guardar` is compare-and-swap on `Lead.Version`, stale writes return `ErrOptimisticLock` without mutation, missing records return `NotFoundError`/`ErrNoEncontrado`, lead lists are filtered conjunctively and ordered by priority descending then ID ascending, and conversation results are chronological.

The domain entities map directly to the schema. `Lead` stores scalar fields plus `Perfil`, optional `Capacidad`, and optional `Intencion` as JSONB. `Mensaje` stores a JSONB attachment in the `mensajes` table. `PlanNutricion` and `Hito` span the `planes` and `hitos` tables. `Ficha` is represented by one JSONB `contenido` column with a unique `lead_id`. `domain.Comprador` matches the `compradores` seed table. There are currently no production repository adapters, no catalog cache, and no `GeneradorID` implementation.

The static catalog contract is file-backed rather than database-backed: `data/proyectos.json`, `data/compradores.json`, `data/afiliados_mock.json`, and `data/brochures/*.md` are the source material. The project pipeline already defines the normalized JSON/domain shapes, including project slugs, buyer categories, normalized age bands, and the contract-required brochure metadata. Contract v1.1 and NFR-E-04 require these read-only responses to be cached in memory so repeated recommendation or brochure requests do not reread disk or PostgreSQL.

Clean Architecture places all implementations in `internal/infrastructure`; adapters may import `internal/usecase` and `internal/domain`, while public usecase/domain contracts and the schema remain unchanged. `cmd/servidor/main.go` currently has a dedicated TODO seam for registering repositories and the motor, but wiring the full application is outside the established Issue #15 scope.

### Affected Areas
- `internal/infrastructure/postgres/` — add PostgreSQL implementations for `LeadRepository`, `PlanRepository`, `FichaRepository`, and the shared SQL/JSON/error mapping needed by those adapters; retain the existing connection and migration helpers.
- `internal/infrastructure/postgres/` — add the file-backed `CatalogoRepository` implementation, loading the four static catalog sources once and returning defensive copies of mutable collections.
- `internal/infrastructure/ids/` — add a ULID-backed implementation of `usecase.GeneradorID`; IDs must remain opaque strings, unique, and no longer than the contract limit of 40 characters.
- `internal/infrastructure/postgres/*_test.go` and `internal/infrastructure/ids/*_test.go` — add focused unit tests plus optional, environment-gated PostgreSQL integration tests. Tests must stay in infrastructure and must not alter public ports, domain types, or the migration.
- `go.mod` and generated module sums — add one pinned ULID dependency if a library is selected; preserve the existing pinned `pgx/v5` version unless a separate authority explicitly changes it.
- Read-only inputs `data/proyectos.json`, `data/compradores.json`, `data/afiliados_mock.json`, and `data/brochures/*.md` — consume their existing shapes only; no source-data or contract changes are needed.

### Coupling and Constraints

1. **Lead CAS is mandatory.** `Guardar` must execute an atomic update with `WHERE lead_id = $id AND version = $expected`. A zero-row update must distinguish a missing lead from a stale version, either through a targeted existence query or an equivalent `UPDATE ... RETURNING`/follow-up strategy. The adapter must increment the version only after a successful compare-and-swap and must not mutate storage on conflict.
2. **Stable queue semantics are SQL semantics.** `Listar` must use optional `afiliado` and `ruta` predicates and `ORDER BY prioridad DESC, lead_id ASC`; relying on database default order would diverge from the existing fake and Contract §3.5.
3. **JSONB preserves contract values.** JSON marshal/unmarshal must cover dynamic profile values, nullable capacity/intention, message attachments, and the complete ficha. JSON decoding normally produces `float64` for numbers, which is compatible with `Perfil.Entero`; adapter tests should still cover nulls and nested values. Nil JSONB values must not be decoded as an empty object accidentally.
4. **Plan persistence is relational plus aggregate reconstruction.** A plan write must persist the base row and its milestones consistently, preferably in one transaction. `PorLead` reconstructs `Hitos`; `HitosVencidos` joins active plans to pending milestones whose date is at or before the requested date; `MarcarHito` updates by milestone ID and should report absence rather than silently succeeding. The port has no plan version, so speculative plan CAS must not be introduced.
5. **Ficha uniqueness is authoritative.** `FichaRepository.Guardar` must support first insert and replacement for the unique `lead_id` record without creating duplicate fichas. `PorLead` must distinguish a missing ficha from a missing lead using the existing repository error vocabulary.
6. **Catalog data is not schema data.** The catalog adapter should be constructed with a data root (or an injected filesystem) and eagerly load or `sync.Once`-load all JSON and brochure content. It should cache the load error deterministically, validate IDs and JSON shapes at load time, and return defensive copies of maps/slices so consumers cannot mutate the cache. Affiliate and brochure misses should return an explicit not-found error.
7. **ULID generation must be concurrency-safe.** The implementation should use a cryptographically safe or library-provided ULID source, avoid global mutable entropy without synchronization, and expose only `Nuevo() string`. Tests should assert parseability, 26-character output, uniqueness over concurrent calls, and no leakage of generator internals.
8. **Authority precedence is fixed.** Contract v1.1 §§0, 3.5, 4, and 5; NFR-E-01, NFR-E-04, NFR-M-01, and NFR-M-03.4; and Software Architecture §3 override the issue sketch. No public port/domain/schema modification is justified by this exploration.

### Approaches

1. **Resource-specific PostgreSQL adapters with explicit SQL (recommended)** — implement lead, plan, and ficha repositories as focused types over the existing `*pgxpool.Pool`, with explicit queries, small JSON helpers, transaction boundaries for plan aggregates, and error translation into `usecase.NotFoundError`, `ErrNoEncontrado`, and `ErrOptimisticLock`.
   - Pros: preserves resource-specific behavior and CAS; makes SQL ordering and filters auditable; keeps Clean Architecture dependencies inward; maps directly to the frozen schema; easy to verify with focused integration scenarios.
   - Cons: more handwritten SQL and mapping code; plan aggregate persistence needs careful transaction design; repository tests need a PostgreSQL boundary for full confidence.
   - Effort: High

2. **Generic reflection/JSON repository over tables** — create one generic persistence helper that serializes arbitrary domain values and generates CRUD operations by convention.
   - Pros: fewer apparent repository files and less repeated CRUD scaffolding; potentially fast for a prototype.
   - Cons: obscures lead CAS and `Listar` ordering; cannot safely express plan/milestone joins and ficha uniqueness without resource-specific escape hatches; weakens error mapping and auditability; likely increases hidden coupling and test complexity.
   - Effort: High

3. **Persist all catalog data in PostgreSQL and read it through the repositories** — seed projects, buyers, affiliates, and brochures into database tables and use SQL for catalog lookups.
   - Pros: one runtime data access mechanism; centralized querying and possible future updates.
   - Cons: contradicts the existing Contract §4 file-backed B→A boundary; requires schema/migration changes explicitly excluded from Issue #15; adds startup seed and deployment coupling; does not satisfy the minimal no-repeat disk-read requirement without another cache.
   - Effort: Very high; rejected by authority

4. **Catalog cache policy: lazy `sync.Once` versus eager constructor load** — both policies retain file-backed sources and one-time reads. Lazy loading keeps construction simple but defers missing-file errors until first use and permanently caches the first failure. Eager loading validates all sources at startup and makes repository methods trivial, but the constructor performs I/O and requires callers to handle initialization errors. The preferred implementation is an explicit constructor with an injected root/filesystem and eager validation; a `sync.Once` internal loader is acceptable if the composition root requires deferred startup.
   - Pros of eager loading: fail-fast deployment, one clear initialization error, no request-time disk I/O, straightforward defensive-copy methods.
   - Cons of eager loading: startup depends on all static assets, and tests must provide a complete fixture set.
   - Effort: Medium

### Recommendation

Proceed with resource-specific PostgreSQL adapters over the existing pgx pool, plus a separate ULID generator package and a file-backed catalog repository with a one-time in-memory cache. Keep the schema and all application-facing interfaces exactly as they are. Use explicit SQL for every port: atomic lead CAS, deterministic queue ordering, chronological conversations, transactional plan/milestone persistence, and upsert-by-lead fichas. Use small shared helpers for JSONB encoding/decoding and PostgreSQL error translation, but do not hide resource-specific queries behind a generic repository.

Construct the catalog adapter from a data root or injected `io/fs` so unit tests do not depend on the working directory. Load and validate projects, buyers, affiliates, and brochures once; return copies of maps/slices and copied affiliate records. Use a pinned, maintained ULID library or a small standard-library implementation with synchronized entropy, while preserving the existing `GeneradorID` interface. Do not change the current pgx pin as incidental cleanup.

The expected implementation is larger than the 400-authored-line review budget. Forecast approximately 650–950 authored lines including SQL mappings, constructors, error handling, ULID code, focused tests, and optional integration coverage. Therefore use sequential main PR slices rather than one oversized PR:

1. **Slice 1 — persistence primitives and leads:** shared JSON/error helpers, ULID generator, `LeadRepository` including CAS/list/conversation behavior, and its focused tests. Clear finish: lead lifecycle and ID tests pass; rollback is limited to the new lead/ID files and dependency.
2. **Slice 2 — nutrition and fichas:** `PlanRepository` and `FichaRepository`, including transactional milestone handling, not-found behavior, and focused PostgreSQL tests. Clear finish: plan/ficha round trips and milestone transitions pass; rollback removes only these adapters/tests.
3. **Slice 3 — catalog cache and optional integration harness:** one-time catalog loading, defensive copies, affiliate/brochure lookup, and skippable real-Postgres integration coverage if the environment supports it. Clear finish: repeated calls perform no repeated file reads and all source shapes are validated; rollback is limited to catalog adapter/tests.

Each slice should include its own focused tests, runtime/database harness result or an explicit N/A rationale, and a rollback boundary. A subsequent task phase should refine the line forecast after concrete file/task sizing.

Decision needed before apply: Yes
Chained PRs recommended: Yes
400-line budget risk: High

### Risks

- PostgreSQL integration tests may be unavailable in CI or local environments; keep them opt-in or skip under `testing.Short()` and retain deterministic unit coverage for JSON, catalog, and ID behavior.
- A zero-row CAS update can be misclassified as not-found or optimistic-lock unless the adapter performs an explicit distinction; this is the highest correctness risk.
- Plan writes that replace milestones incorrectly could lose externally observed milestone IDs or statuses; define update semantics before implementation and keep the operation transactional.
- `sync.Once` permanently caches an initial catalog load error; eager construction or explicit error state should be chosen deliberately.
- Returning internal catalog maps/slices directly would allow callers or the motor to corrupt the cache; defensive copies are required.
- A new ULID dependency must be pinned exactly and kept separate from unrelated pgx version changes.
- The existing server composition root has only a registration seam; full use-case wiring should remain a later issue unless the proposal explicitly expands scope.
- No schema, public port, domain, source-data, or unrelated workspace paths are required for this change.

### Ready for Proposal

Yes. The proposal should formalize the three sequential review slices, exact constructors and error semantics, the eager catalog-cache boundary, optional PostgreSQL test prerequisites, and the explicit no-change constraints for public contracts and schema.
