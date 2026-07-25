# Apply Progress: Recommendations and 2x2 Routing — Issue #10

## Change

`issue-10-recommendations-routing`

## Mode and Delivery

- Artifact store: hybrid
- Implementation mode: Standard (`strict_tdd: false`)
- Delivery: auto-chain, stacked-to-main
- Current work unit: Work Unit 1 — catalog recommendations (PR 1 targets `main`)
- Scope boundary: only `internal/domain/motor/recomendar.go` and `internal/domain/motor/recomendar_test.go`

## Completed Tasks

- [x] 1.1 Issue #8 gate recorded as merge SHA `7a5b02a613c581568905936a35cb8feb841e62ab`; maintainer evidence states it is an ancestor of `origin/main` and current `HEAD`.
- [x] 1.2 Planned diff inspected; no capacity, shared-domain, Docs, data, HTTP, use-case, or pipeline paths were edited.
- [x] 2.1 Added table-driven recommendation behavior tests for canonical-ID aggregation, orphan/empty IDs, non-nil empty output, cap-three behavior, and shuffled-input determinism.
- [x] 2.2 Added tests for inclusive `ceil(0.8 * PrecioDesde)` boundaries, one-peso failures, negative budget/non-positive price rejection, input immutability, copied catalog fields/rationale, and overflow-safe fraction equality/order.
- [x] 2.3 Implemented pure `RecomendarProyectos` with catalog-only aggregation, integer eligibility, exact `math/bits` fraction ranking, deterministic total ordering, copied card data, rationale, and top-three truncation.

## Files Changed

| File | Action | Description |
|---|---|---|
| `internal/domain/motor/recomendar.go` | Created | Catalog-backed aggregation, integer affordability, exact ranking, card construction, and cap. |
| `internal/domain/motor/recomendar_test.go` | Created | Table-driven Work Unit 1 behavior, boundary, purity, and fraction-order tests. |

## Work Unit Evidence

| Evidence | Result |
|---|---|
| Focused test command and exact result | `go test ./internal/domain/motor -run '^TestRecomendarProyectos$'` — PASS, exit 0. Re-run with `-count=1` also PASS; 1 focused test function completed. |
| Runtime harness command/scenario and exact result | N/A — pure in-memory domain function with no adapter, persistence, HTTP, LLM, or orchestration runtime boundary. |
| Broader validation | `go test ./internal/domain/motor/...` — PASS, exit 0. `go test ./internal/domain/motor/... -cover` — PASS, exit 0, coverage 97.2%. `go build ./...` — PASS, exit 0. |
| Rollback boundary | Revert/delete exactly `internal/domain/motor/recomendar.go` and `internal/domain/motor/recomendar_test.go`; no unrelated behavior or protected path is removed. |

## Deviations

None — implementation follows the approved pure-function design. Ranking compares original integer counts with `math/bits.Mul64`; floating-point rate is used only for display in the output card.

## Remaining Tasks

- [ ] 3.1–3.3 Matrix work unit; intentionally deferred until this slice merges to `main`.
- [ ] 4.1–4.3 Final combined validation and SDD verification after all implementation work units.

## Status

Work Unit 1 complete. Ready for `sdd-verify` for this slice, subject to the orchestrator's chained-PR workflow; matrix implementation remains pending for the subsequent slice.
