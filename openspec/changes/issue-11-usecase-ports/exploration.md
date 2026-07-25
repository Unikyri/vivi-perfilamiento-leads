# Exploration: issue-11-usecase-ports

## Current State

`internal/usecase` currently contains only `doc.go`; the merged domain contract is the source for port signatures. The domain exposes `Lead` (with `Version`, `Perfil`, optional `Capacidad`, and optional `Intencion`), `PlanNutricion` (with mutable `Hitos`), `Ficha`, `Mensaje`, `Proyecto`, and `Comprador`. `Lead.Transicionar` owns lifecycle validation; optimistic locking remains repository responsibility.

## Authority Reconciliation

The governing order is Contract v1.1, then NFR-M-01/M-04, then Software Architecture §3, then merged-code compatibility. The Issue #11 sketch is supplementary.

- Contract §3.5 supports only optional `afiliado` and `ruta` lead-list filters and requires priority-descending output.
- Contract §5 defines lead `version` optimistic locking; stale writes must not overwrite stored data.
- Contract §§4.3, 6, and 7 define affiliated-record, event, and structured-LLM shapes.
- NFR-M-01 requires usecase to import domain only and owns application ports; NFR-M-04 requires in-memory repositories and LLM stubs for tests.
- Architecture §3 additionally establishes `FichaRepository`, `Reloj`, and an event bus as application boundaries.

The sketch's unconditional `Guardar` version increment is rejected: it violates Contract §5. Its map-backed listing is also insufficient because map iteration leaks nondeterministic order.

## Recommended Port Surface

Create context-aware, domain-value contracts in `internal/usecase/puertos.go`. The port layer must not import HTTP, SQL, ADK, or provider SDK types.

- `LeadRepository`: create, get by ID, compare-and-swap update, and deterministic list.
- `PlanRepository`: create, look up by lead, simple update, overdue milestones, and mark milestone. Do not invent plan compare-and-swap because `PlanNutricion` has no version field.
- `FichaRepository`: save and get by lead.
- `CatalogoRepository`: projects, buyers, affiliate lookup, and brochure content.
- `LLMProvider`: text and audio turns using `EntradaTurno`, `SalidaTurno`, and `Audio` application DTOs matching Contract §7.
- `Reloj`, `BusEventos`, and `GeneradorID` as infrastructure-free application boundaries.

Use inspectable repository errors (`ErrNotFound`, `ErrOptimisticLock`) rather than nil-success or zero values. A missing ficha remains distinguishable from a missing lead at the usecase boundary.

## Fake Semantics

The initial fake set is intentionally narrow: a real `LeadRepository` fake, plus minimal LLM, clock, and deterministic-ID fakes. Plan/ficha fakes are deferred until a focused usecase scenario consumes them; declaring every port does not justify speculative test doubles.

The lead fake must:

1. Reject duplicate creates and return `ErrNotFound` for absent single-record lookups.
2. Implement real compare-and-swap: compare incoming `Lead.Version` with stored version, return `ErrOptimisticLock` without mutation on mismatch, and return a committed copy whose version incremented exactly once on success.
3. Deep-copy every mutable value at the fake boundary. A shallow copy leaks maps, slices, pointers, profile values, capacity breakdown, intention signals, and message attachments.
4. Apply non-nil optional filters conjunctively and return a non-nil empty list when unmatched.
5. Sort list results by `Prioridad` descending, then `LeadID` ascending. The Contract specifies priority but not ties; the secondary key is a frozen deterministic implementation rule.

## Approaches

The selected approach is a semantic port layer plus behaviorally correct test fakes. Thin shallow fakes are rejected because they hide aliases and stale-write defects. Generic repositories are rejected because they erase lead CAS and ficha absence semantics.

## Delivery Forecast

Expected authored work is 380–520 lines if all fakes are included. To keep reviewable scopes below 400 lines, use two sequential slices:

1. ports, DTOs, error/filter contracts, and focused compile/shape tests;
2. lead fake and tests for CAS, cloning, filtering, ordering, and LLM/clock/ID test doubles.

A third slice is only warranted if plan/ficha fake behavior becomes justified by approved scenarios.

## Risks

- Clone helpers must preserve concrete integer values in `CampoPerfil.Valor`; avoid JSON round-tripping.
- Plan CAS would be invented behavior and is excluded.
- Port names may differ from the sketch, but semantics and strict inward dependencies are authoritative.
- `FichaRepository` remains minimal because Contract §3.6 and Architecture §3 require it even though the NFR's short list omits it.

## Recommendation

Proceed to proposal with all authority-backed interfaces but only the LeadRepository fake and minimal LLM/clock/ID fakes. Preserve the two-slice delivery plan and the strict `internal/usecase`-only scope.