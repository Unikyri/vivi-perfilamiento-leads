# Tasks: Usecase Ports and In-Memory Lead Fake — Issue #11

## Review Workload Forecast

| Field | Value |
|---|---|
| Estimated changed lines | Slice 1: 190–250; Slice 2: 280–360; total: 470–610 |
| 400-line budget risk | High overall; each enforced slice is below 400 |
| Chained PRs recommended | Yes |
| Suggested split | Ports + shape tests → fakes + behavioral tests |
| Delivery strategy | Resolved chained delivery |
| Chain strategy | sequential-to-main (PR 2 starts only after PR 1 merges) |

Decision needed before apply: No
Chained PRs recommended: Yes
Chain strategy: sequential-to-main
400-line budget risk: High

### Suggested Work Units

| Unit | Goal | Likely PR | Focused test command | Runtime harness | Rollback boundary |
|---|---|---|---|---|---|
| 1 | Ports and shape proof | PR 1: base `main` | `go test ./internal/usecase/... -run 'Test(Puertos|DTO|Errores)'` | N/A—contracts have no runtime boundary | Revert `puertos.go`, `puertos_test.go` |
| 2 | Lead fake and behavior proof | PR 2: base `main` after PR 1 merges | `go test -race ./internal/usecase/...` | N/A—in-memory doubles prohibit external calls | Revert `fakes_test.go`, `fakes_behavior_test.go` |

## Phase 1: Slice 1 — Ports and Shape Tests

- [ ] 1.1 Create `internal/usecase/puertos.go` with Spanish pointer ports (`Crear`/`PorID`/`Guardar`/`Listar`), Contract §§0/4.3/6/7 DTOs, `FiltroLeads`, errors, and event constants; allow only stdlib and `internal/domain` imports.
- [ ] 1.2 Declare `LeadRepository`, `PlanRepository`, `FichaRepository`, `CatalogoRepository`, `LLMProvider.Nombre`, `MensajeriaGateway`, `Reloj`, `BusEventos`, and `GeneradorID`; defer Plan/Ficha/Catalog fakes and plan CAS.
- [ ] 1.3 Create `internal/usecase/puertos_test.go` with compile-time method-shape, pointer-result, DTO JSON, and `errors.Is(ErrNoEncontrado/ErrOptimisticLock)` tests.
- [ ] 1.4 Prove Slice 1 with its focused `go test`, `go list -deps ./internal/usecase/...` import audit, and a path audit limited to `internal/usecase` plus Issue #11 artifacts; keep authored diff below 400 lines.

## Phase 2: Slice 2 — Fake and Behavioral Tests

- [ ] 2.1 Create `internal/usecase/fakes_test.go`: mutex-protected `FakeLeadRepository`, real CAS, duplicate/absence errors, recursive defensive clones, chronological message handling, and minimal deterministic Lead/LLM/clock/ID doubles.
- [ ] 2.2 Create `internal/usecase/fakes_behavior_test.go` table tests for stale/committed versions, unchanged stale storage, nested-map/slice isolation, conjunctive filters, priority/ID ordering, non-nil empty lists, and cloned conversations.
- [ ] 2.3 Run `go test -race ./internal/usecase/...`; retain only the approved fake set, with no outer-layer imports or plan/ficha/catalog fake.
- [ ] 2.4 Run `go test ./... && go build ./... && go vet ./...`; audit the complete chained diff and each slice’s authored additions/deletions (<400) before review.

## Phase 3: Review Handoff

- [ ] 3.1 Commit each verified slice as one work unit with its tests, Issue #11 linkage, a conventional message, and the stated rollback boundary; merge PR 1 to `main` before creating PR 2 from that updated `main`.
