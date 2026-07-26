# Tasks: Calificar lead and generate deterministic advisor ficha (Issue #20)

## Review Workload Forecast

| Slice | Scope / PR base | Runtime | Tests | Total |
|---|---|---:|---:|---:|
| 1 | Qualification; base `feature/bloque-a` | 130 | 210 | 340 |
| 2 | Ficha; base PR #1 branch | 135 (includes helper <=20) | 215 | 350 |

Decision needed before apply: No
Chained PRs recommended: Yes
Chain strategy: feature-branch-chain
400-line budget risk: High

Both authored totals are <=400. PR #1 targets `feature/bloque-a`; PR #2 targets the PR #1 branch and must be retargeted/rebased if its diff includes slice 1.

| Unit | Focused test | Runtime harness | Rollback |
|---|---|---|---|
| 1 | `go test ./internal/usecase -run TestCalificarLead` | deterministic fake repos/bus ordering scenario | revert new `calificar_lead.go` and test |
| 2 | `go test ./internal/usecase -run TestGenerarFicha` | fixed clock/ID, retry-after-CAS-failure scenario | revert new ficha files and any <=20-line helper delta |

## Traceability

| Issue #20 capability | Tasks |
|---|---|
| Qualification: candidate, KNN, route, cupo, durable event | 1.1–1.4 |
| Deterministic ficha: ordered content, strict alert, handoff/retry | 2.1–2.4 |

## Phase 1: Slice 1 — Qualify lead

- [x] 1.1 In `internal/usecase/calificar_lead_test.go`, add deterministic RED cases for non-`CALIFICADO`, blank ID/intention/capacity, catalog/read failures, and assert no lead save or event.
- [x] 1.2 Create `internal/usecase/calificar_lead.go` with package-private decision/candidate helpers: capacity(0), lowest positive affordable price, final capacity, exact zones/dependents, K=30, and immutable ordered recommendations.
- [x] 1.3 Complete RED/GREEN matrix tests in `calificar_lead_test.go`: each conversion signal, blank external caja, four routes, 1.2 ratio cap, semáforo, cupo, candidate-zero median fallback, and recommendation determinism.
- [x] 1.4 Implement qualification route/state mapping and CAS save-before-one-`RutaDecidida` event; test `ASESOR` remains `CALIFICADO`, other mappings, save failure/no event, and event payload order.

## Phase 2: Slice 2 — Generate ficha

- [x] 2.1 In `internal/usecase/generar_ficha_test.go`, add RED eligibility/read/cancellation tests proving non-`CALIFICADO`/non-`ASESOR` and failures write neither ficha nor lead.
- [x] 2.2 Create `internal/usecase/generar_ficha.go` using the shared decision; build Contract fields in fixed order, exact low-confidence warning, ordered benefits/rent argument, and alert active only for rate `>0.20`.
- [x] 2.3 Add fixed-clock/ID tests for byte-stable content, recommendation parity, non-aliasing, threshold `0.20` inactive versus above active, and no LLM use.
- [x] 2.4 Implement and test ficha-upsert before `ENTREGADO` CAS: ficha-save leaves lead unchanged; lead-save failure remains `CALIFICADO`/`ASESOR`; retry reuses ID/time, upserts, then saves. Change `calificar_lead.go` helper only if <=20 authored lines.

## Phase 3: Slice evidence and guard

- [x] 3.1 Run each unit's focused command, record exact result and fail-at-step harness evidence; keep each PR's authored `git diff --stat` total <=400.
- [x] 3.2 Validation commands (`go test ./...`, `go vet ./...`, `go build ./...`, race/module/diff checks) pass, including the 63-line qualification edge suite; committed as `23f4a2e` (`test(usecase): cover qualification edges`).

## Scope and Definition of Done

Forbidden: `internal/domain/**`, `internal/domain/motor/**`, `internal/usecase/puertos.go`, adapters/infrastructure, HTTP/frontend, migrations/config, Contract/Wiki, ports, LLM/narrative. Done when all 10 tasks pass, both slices are independently revertible, all error paths prove no premature write/event, state/event ordering is asserted, and both slice totals remain <=400 authored lines.
