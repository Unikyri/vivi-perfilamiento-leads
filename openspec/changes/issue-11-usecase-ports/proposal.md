# Proposal: Usecase Ports and In-Memory Fakes — Issue #11

## Intent

`internal/usecase` holds only `doc.go`, so no downstream work has a contract to code against: Postgres/LLM adapters, HTTP handlers, and `CalificarLead` (#17) all wait on it. NFR-M-01 assigns port ownership to `usecase` (domain-only imports); NFR-M-04 requires in-memory repositories and an LLM stub so CI never calls a real provider. Contract §5 optimistic locking and §3.5 list filtering/ordering must be provable by tests *before* a real adapter can be trusted to honor them.

## Scope

### In Scope
- `internal/usecase/puertos.go`: `LeadRepository`, `PlanRepository`, `FichaRepository`, `CatalogoRepository`, `LLMProvider`, `MensajeriaGateway`, `Reloj`, `BusEventos`, `GeneradorID`; `LeadFilter`; `EntradaTurno`/`SalidaTurno`/`Audio` DTOs; `ErrNotFound`, `ErrOptimisticLock`, wrapping `NotFoundError{Resource, ID}`.
- `internal/usecase/fakes_test.go`: `FakeLeadRepository` with real compare-and-swap and deep clones; minimal LLM, clock, and deterministic-ID doubles.
- Focused usecase tests: CAS, isolation, missing-record, filtering, ordering, interface assertions.

### Out of Scope
- Plan, ficha, catalog, messaging, and event-bus fakes (deferred until a scenario consumes them).
- Any usecase orchestration behavior (`CalificarLead`, conversation turns, nutrition flows).
- Plan versioning/CAS, HTTP status mapping, SQL/ADK/provider imports.
- Edits to `internal/domain/`, `internal/adapters/`, `internal/infrastructure/`, `data/`, `Docs/`.

## Capabilities

### New Capabilities
- `usecase-ports`: application-owned port contracts, application errors, and test-double semantics for the usecase layer.

### Modified Capabilities
- None. `Lead`, `PlanNutricion`, `Ficha`, `Mensaje`, `Proyecto`, `Comprador` are consumed unchanged.

## Approach — Frozen Decisions

| # | Decision | Rule |
|---|---|---|
| 1 | Port ownership | Interfaces live in `usecase`, take `context.Context`, exchange domain values. Imports limited to `context`, `errors`, `time`, `internal/domain`. |
| 2 | Compatibility, not redesign | Signatures must compile against merged domain types as-is. The Issue #11 sketch is supplementary; where it conflicts with Contract v1.1 it loses. |
| 3 | Real lead CAS | `Update` compares `Lead.Version` with stored version, returns `ErrOptimisticLock` **without mutating storage** on mismatch, increments exactly once on success, and returns the committed copy. The sketch's unconditional increment is rejected (Contract §5). |
| 4 | Explicit absence | Single-record lookups return a zero value plus a sentinel/wrapped not-found error. Never nil-success. A missing ficha stays distinguishable from a missing lead so HTTP can later split `FICHA_NO_DISPONIBLE` from `LEAD_NO_ENCONTRADO`. |
| 5 | Deep defensive copies | Every mutable value crossing the fake boundary is cloned: `Perfil`/`CampoPerfil.Valor` nested maps and slices, `Capacidad`+`Desglose`, `Intencion.Senales`, plan `Hitos`, `ConsentimientoEn`, message `Adjunto`. No JSON round-trip (integer type preservation). |
| 6 | Deterministic listing | `LeadFilter{Afiliado *bool, Ruta *domain.Ruta}` applied conjunctively; sort a copy by `Prioridad` desc then `LeadID` asc; return non-nil empty slice when unmatched. Map iteration never reaches the result. |
| 7 | Narrow fake set | Declaring a port does not justify a fake. Only the lead fake plus minimal LLM/clock/ID doubles ship now; plan/ficha fakes are deferred rather than inventing plan CAS. |

## Affected Areas

| Area | Impact | Description |
|---|---|---|
| `internal/usecase/puertos.go` | New | Ports, DTOs, filter, application errors |
| `internal/usecase/fakes_test.go` | New | Lead fake, clone helpers, minimal doubles |
| `internal/usecase/*_test.go` | New | Table-driven contract tests |
| `internal/usecase/doc.go` | Read-only | Package boundary comment already correct |
| `internal/domain/*` | Read-only | Compatibility source only |
| Wiki docs 09/10/11 | Read-only | Authority; no edits in scope |

## Risks

| Risk | Likelihood | Mitigation |
|---|---|---|
| Ad-hoc clone type switches miss nested `any` values | High | Explicit recursive clone for `map[string]any` / `[]any`; tests mutate returned values and re-read storage |
| Port names drift from the issue sketch | Med | Preserve semantics, pick one vocabulary, record the mapping in the delta spec |
| Speculative ports without adapters (catalog, event bus) rot | Med | Declare minimal method sets only; no fakes, no callers, easy to delete |
| Contract leaves priority ties unspecified | Med | `LeadID` ascending frozen as a documented implementation rule (decision 6) |
| 380–520 line forecast exceeds the 400-line review budget | High | Two sequential slices (below); each has its own tests and rollback |
| Declared-but-unimplemented ports block compilation of later slices | Low | Slice 1 adds compile-time `var _ Port = (*fake)(nil)` assertions only for implemented fakes |

## Delivery

Two sequential review slices (`ask-on-risk` resolved to chained delivery; no `size:exception` needed):

1. **Ports slice** — `puertos.go`, DTOs, errors, filter, shape/assertion tests.
2. **Fake slice** — lead fake, clone helpers, CAS/isolation/filter/order tests, minimal LLM/clock/ID doubles.

Slice 2 targets slice 1's branch. A third slice is only warranted if approved scenarios require plan/ficha fakes.

## Rollback Plan

Both files are new and imported by no production code. Rollback is `git revert` of the slice commit (or deleting `puertos.go` and `fakes_test.go`) on the Issue #11 branch. Nothing in `domain`, adapters, infrastructure, migrations, or data is touched, so there is no schema, data, or API rollback path to manage. Reverting slice 2 alone leaves slice 1 compiling.

## Dependencies

- Issue #6 domain contract types — satisfied (`a3654ce` reachable from `origin/main`).
- Issue #7 lead state machine — satisfied (merge `e78f9be…` is an ancestor of `origin/main`).
- Go 1.24+ toolchain for `go build` / `go test`.

## Success Criteria

- [ ] `go build ./...` and `go test ./internal/usecase/...` pass; `go vet` clean.
- [ ] `internal/usecase` imports only stdlib plus `internal/domain` (verified by inspection or `go list -deps`).
- [ ] Stale-version `Update` returns `ErrOptimisticLock`, storage unchanged; successful `Update` increments version exactly once.
- [ ] Duplicate `Create` rejected; absent-ID lookups return `ErrNotFound` (matched via `errors.Is`) with a zero value.
- [ ] Mutating any returned map/slice/pointer does not alter fake storage, and vice versa.
- [ ] `List` with shuffled insertion order yields identical `Prioridad` desc / `LeadID` asc output; unmatched filters return a non-nil empty slice.
- [ ] No file outside `internal/usecase/` and the Issue #11 OpenSpec folder appears in the diff.
- [ ] Each slice stays under the 400 authored-line review budget.

## Proposal question round

Autonomous run — these assumptions need confirmation or correction:

1. Declaring `CatalogoRepository`, `MensajeriaGateway`, `BusEventos`, and `GeneradorID` now (no implementations) is acceptable, versus deferring them to the issue that first consumes them.
2. New leads start at `Version = 1` to match `leads.version INT NOT NULL DEFAULT 1`, and `Create` rejects a supplied version other than 0 or 1.
3. `Prioridad` ties break on `LeadID` ascending — a deliberate addition beyond Contract §3.5.
4. Port vocabulary is `Create`/`GetByID`/`Update`/`List` (not the sketch's `Guardar`/`Buscar`), with Spanish reserved for domain and DTO field names.
5. `EntradaTurno`/`SalidaTurno` stay minimal Contract §7 shapes in this issue; provider-specific fields wait for the LLM adapter issue.

## Consumer Compatibility Reconciliation

Inspection of dependent Issues #15, #16, and #18 found a consistent, compatible vocabulary: `Crear`, `PorID`, `Guardar`, `Listar`, pointer repository entities, `FiltroLeads`, `LLMProvider.Nombre`, and the `Reloj`/`BusEventos` shapes proposed by Issue #11. Those consumer sketches remain supplementary, but no Contract, NFR, architecture, or merged-code rule conflicts with retaining this vocabulary. It is therefore preferred over the earlier English value-return alternative.

This reconciliation **does not** accept the sketch's behavioral defects: `Guardar` MUST compare the incoming `Lead.Version` with stored version, fail with `ErrOptimisticLock` without storage mutation on mismatch, and increment the caller's version exactly once only after success. Lookup/list results remain defensive copies; filters remain conjunctive and results deterministic by `Prioridad` descending then `LeadID` ascending. `LLMProvider` includes `Nombre() string` because the architecture's health/metrics flow needs provider identity. `MensajeriaGateway` is declared to satisfy NFR-M-01 but has no fake in this issue.
