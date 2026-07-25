# Apply Progress: Buyer-Twin kNN — PR 4/5

## Scope
- `auto-chain` / `stacked-to-main`; standard mode (`strict_tdd: false`).
- Completed 2.1, 2.2, 3.1, and 3.2 only: safe projection plus project-local statistics and weighted Gower.
- Remaining: PR 5 selection/safety/evidence (4.1–4.3).

## Delivered
- `knn.go`: immutable per-project affiliated statistics with a five-affiliate threshold, deterministic category mode (`A < B < C` ties), odd/even numeric median, independent non-affiliate category/age fallback, and capped weighted Gower distance.
- `knn_test.go`: focused tests for five-affiliate imputation, odd/even medians, deterministic mode, below-threshold omission, independent fallback dimensions, cross-project rejection, weighted renormalization, symmetry, identity, `[0,1]`, and zero dependents.
- Only tasks 3.1/3.2 were marked complete; no selection or `GemeloKNN` behavior was added.

## Decisions and Evidence
- Statistics count affiliated buyers by exact non-empty `ProyectoID`; category and age samples are collected independently, and invalid/missing samples do not create a dimension.
- Imputation applies only to non-affiliates and only within the exact project when at least five affiliates exist. A missing category never blocks an available age fallback, and vice versa.
- Gower weights are category `.35`, zone `.20`, age `.15`, affiliation `.15`, and dependents `.15`; only mutually present dimensions contribute and the denominator is renormalized. Numeric terms use age range `32.5`, dependents range `10`, and cap at `1`.
- Zero dependents remain present; inputs and source buyer slices are not mutated. No global, cross-project, price, or financial data is used.
- `go test ./internal/domain/motor/... -run 'TestGemeloKNN(Gower|Imputation)' -count=1` — PASS.
- `go test ./internal/domain/motor/... -count=1` — PASS.
- `git diff --check` — PASS; authored implementation/test delta is 276 lines before artifact updates.
- Runtime: N/A — pure, currently unwired domain service has no runtime/integration boundary in this work unit.
- Rollback: revert the statistics/Gower additions and focused tests in `internal/domain/motor/knn.go` and `knn_test.go`, plus the 3.1/3.2 task checkboxes and this progress artifact; no capacity, pipeline, data, Docs, or other production files are involved.

## Remaining Tasks
- [ ] 4.1 RED: selection and safety tests (PR 5).
- [ ] 4.2 GREEN: non-mutating ordering and value-only neighbor results (PR 5).
- [ ] 4.3 Focused/full verification and path audit (PR 5).

## Deviations and Risks
- None from the approved statistics/Gower design. Neighbor selection and `GemeloKNN` orchestration remain intentionally deferred to PR 5.
- The repository contains unrelated pre-existing untracked `Docs/` content; it was not modified or included in this work unit.
