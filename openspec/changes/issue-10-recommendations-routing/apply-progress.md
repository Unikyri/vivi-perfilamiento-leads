# Apply Progress: Recommendations and 2x2 Routing — Issue #10

## Change

`issue-10-recommendations-routing`

## Mode and Delivery

- Artifact store: hybrid
- Implementation mode: Standard (`strict_tdd: false`)
- Delivery: auto-chain, stacked-to-main
- Current work unit: Work Unit 2 — 2x2 routing matrix (final PR targets `main` after recommendation PR #58 merged at `9f3e3703f310276badce8bf1c2b8d04bcd087cd4`)
- Scope boundary: only `internal/domain/motor/matriz.go` and `internal/domain/motor/matriz_test.go` for this slice; cumulative change retains Work Unit 1 recommendation evidence.

## Completed Tasks

- [x] 1.1 Issue #8 gate recorded as merge SHA `7a5b02a613c581568905936a35cb8feb841e62ab`; maintainer evidence states it is an ancestor of `origin/main` and current `HEAD`.
- [x] 1.2 Planned diff inspected; no capacity, shared-domain, Docs, data, HTTP, use-case, or pipeline paths were edited.
- [x] 2.1 Added table-driven recommendation behavior tests for canonical-ID aggregation, orphan/empty IDs, non-nil empty output, cap-three behavior, and shuffled-input determinism.
- [x] 2.2 Added tests for inclusive `ceil(0.8 * PrecioDesde)` boundaries, one-peso failures, negative budget/non-positive price rejection, input immutability, copied catalog fields/rationale, and overflow-safe fraction equality/order.
- [x] 2.3 Implemented pure `RecomendarProyectos` with catalog-only aggregation, integer eligibility, exact `math/bits` fraction ranking, deterministic total ordering, copied card data, rationale, and top-three truncation.
- [x] 3.1 Added table-driven matrix tests for all affiliate quadrants, exact 0.95/1.05 boundaries, non-affiliate conversion precedence, `MEDIA` as low intention, and empty/unknown intention.
- [x] 3.2 Added matrix tests for NaN, positive/negative infinity, negative ratios never reaching `ASESOR`, and route-only purity/repeated deterministic calls without input mutation.
- [x] 3.3 Implemented `EntradaMatriz` and pure `Matriz2x2` with finite non-negative ratio validation, conversion-first precedence, affiliate 0.95 routing, and non-affiliate 1.05 advisor threshold.
- [x] 4.1 Ran focused matrix tests plus full motor tests and coverage; coverage exceeded the 90% target.
- [x] 4.2 Ran repository build and final diff/path inspection; protected capacity files and all paths outside the two matrix files remained untouched. Pre-existing untracked `Docs/` inventory was not modified.
- [x] 4.3 Performed final SDD verification from merged SHA `d75c25ffca0ca5f4d95949cf40e73e5f57e856b8`; 6/6 requirements and 8/8 scenarios passed.

## Files Changed

| File | Action | Description |
|---|---|---|
| `internal/domain/motor/recomendar.go` | Created in Work Unit 1 | Catalog-backed aggregation, integer affordability, exact ranking, card construction, and cap. |
| `internal/domain/motor/recomendar_test.go` | Created in Work Unit 1 | Table-driven recommendation behavior, boundary, purity, and fraction-order tests. |
| `internal/domain/motor/matriz.go` | Created in Work Unit 2 | Pure 2x2 route selection with exact interface and finite-ratio guard. |
| `internal/domain/motor/matriz_test.go` | Created in Work Unit 2 | Table-driven route matrix, invalid input, threshold, precedence, and purity tests. |

## Work Unit Evidence

### Work Unit 1 — Catalog recommendations

| Evidence | Result |
|---|---|
| Focused test command and exact result | `go test ./internal/domain/motor -run '^TestRecomendarProyectos$'` — PASS, exit 0. Re-run with `-count=1` also PASS; 1 focused test function completed. |
| Runtime harness command/scenario and exact result | N/A — pure in-memory domain function with no adapter, persistence, HTTP, LLM, or orchestration runtime boundary. |
| Broader validation | `go test ./internal/domain/motor/...` — PASS, exit 0. Prior coverage 97.2%. `go build ./...` — PASS, exit 0. |
| Rollback boundary | Revert/delete exactly `internal/domain/motor/recomendar.go` and `internal/domain/motor/recomendar_test.go`; no unrelated behavior or protected path is removed. |

### Work Unit 2 — 2x2 routing matrix

| Evidence | Result |
|---|---|
| Focused test command and exact result | `go test ./internal/domain/motor -run '^TestMatriz2x2$' -count=1` — PASS, exit 0; all 21 table cases completed. |
| Runtime harness command/scenario and exact result | N/A — `Matriz2x2` is a pure in-memory route-only function with no adapter, persistence, HTTP, LLM, or orchestration runtime boundary. |
| Broader validation | `go test ./internal/domain/motor/...` — PASS, exit 0. `go test ./internal/domain/motor/... -cover` — PASS, exit 0, coverage 97.4%. `go test ./...` — PASS, exit 0. `go build ./...` — PASS, exit 0. `go vet ./...` — PASS, exit 0. |
| Scope/diff check | PASS, exit 0; only `internal/domain/motor/matriz.go` and `internal/domain/motor/matriz_test.go` are new implementation paths, with no capacity-file change. Pre-existing untracked `Docs/` files were detected and left untouched. |
| Rollback boundary | Revert/delete exactly `internal/domain/motor/matriz.go` and `internal/domain/motor/matriz_test.go`; recommendation behavior and all protected paths remain independently rollback-safe. |

## Deviations

None — implementation follows the approved pure-function design and exact `EntradaMatriz`/`Matriz2x2` interface. Invalid ratios are treated as low capacity and therefore cannot route to `ASESOR`; non-affiliate high-intention conversion precedence is evaluated before ratio thresholds.

## Remaining Tasks

None — all 11 tasks are complete. Archive remains the next SDD phase after this verification evidence is merged.

## Status

Implementation and verification are complete. Ready for archive after the verification-evidence PR merges.
