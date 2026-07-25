# Tasks: Lead State Machine (Issue #7)

## Review Workload Forecast

| Field | Value |
|---|---|
| Estimated changed lines | 300–360 authored lines: `estado.go` 60–75, `estado_test.go` 165–205, tasks 75–80 |
| 400-line budget risk | Low |
| Chained PRs recommended | No |
| Suggested split | Single PR from `feat/issue-7-lead-state-machine` to `main` |
| Delivery strategy | ask-on-risk |
| Chain strategy | stacked-to-main |

Decision needed before apply: No
Chained PRs recommended: No
Chain strategy: stacked-to-main
400-line budget risk: Low

### Suggested Work Units

| Unit | Goal | Likely PR | Focused test command | Runtime harness | Rollback boundary |
|---|---|---|---|---|---|
| 1 | Add and prove the pure-domain lifecycle policy | Single PR | `go test ./internal/domain/... -run 'Test(Transicion|Estados)'` | N/A — deterministic in-memory domain behavior has no runtime boundary | Remove `internal/domain/estado.go` and `internal/domain/estado_test.go` |

## Phase 1: Domain Policy

- [x] 1.1 Create `internal/domain/estado.go` with one unexported, immutable-in-practice adjacency table and private Contract destination order for exactly the 11 allowed edges.
- [x] 1.2 Add `ErrTransicionInvalida` (including `TRANSICION_INVALIDA` error text), `PuedeTransicionar`, and defensive, ordered `EstadosPosibles`; return zero-length results for terminal/unknown sources.
- [x] 1.3 Add `(*Lead).Transicionar` so it validates before assignment and returns the typed value without changing `Lead.Estado` on failure.

## Phase 2: Domain Tests

- [x] 2.1 Create table-driven `estado_test.go` coverage for all 11 valid edges, successful assignment, and expected ordered query results.
- [x] 2.2 Test every non-listed known-state pair, self-transition, empty/`DESCONOCIDO` source or target: `errors.As` exposes exact `Desde`/`Hacia`, error code is present, and failure is atomic.
- [x] 2.3 Test `CERRADO` and `DESPEDIDO` as terminal and unknown sources as unreachable: `PuedeTransicionar` is false and `EstadosPosibles` is zero-length.
- [x] 2.4 Mutate a `EstadosPosibles(EstadoLeadCalificado)` result and assert a second call retains `Entregado`, `EnNutricion`, `Remarketing`, `Despedido` order.

## Phase 3: Verification and Delivery Guard

- [x] 3.1 Format both new Go files and run the focused transition/query test command; record its exact result with the work unit.
- [x] 3.2 Run `go test ./internal/domain/...`, then `go test ./...`, `go vet ./...`, and `go build ./...`.
- [x] 3.3 Guard the domain boundary with `go list -f '{{join .Imports " "}}' ./internal/domain` and diff review: only the two domain files change; no ADK, outer-layer import, enum, field, or JSON-tag change.
- [x] 3.4 Commit code and tests together as one reviewable work unit (for example, `feat(domain): add lead state machine`) after all checks pass.
