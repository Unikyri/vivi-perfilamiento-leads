# Tasks: Buyer-Twin kNN — Issue #9

## Review Workload Forecast

| Field | Value |
|---|---|
| Estimated changed lines | 954–1,144 total; each sequential work unit is capped at 400 authored lines |
| 400-line budget risk | High |
| Chained PRs recommended | Yes |
| Suggested split | Five sequential, independently revertible PRs to `main` |
| Delivery strategy | auto-chain |
| Chain strategy | stacked-to-main |

Decision needed before apply: No
Chained PRs recommended: Yes
Chain strategy: stacked-to-main
400-line budget risk: High

Each PR targets `main`, merges before the next begins, and contains only its named unit.

### Suggested Work Units

| Unit | Goal / likely PR (max lines) | Focused test command | Runtime harness | Rollback boundary |
|---|---|---|---|---|
| 1 | SDD exploration, proposal, state (160) | `git diff --check` | N/A: planning-only | Those three artifacts |
| 2 | SDD spec, design, tasks (≤260) | `git diff --check` | N/A: planning-only | Those three artifacts |
| 3 | API/projection code and RED/GREEN tests (≤340) | `go test ./internal/domain/motor/... -run TestGemeloKNNProjection` | N/A: pure, unwired service | Projection portions of `knn.go`/`knn_test.go` |
| 4 | Statistics/Gower code and RED/GREEN tests (≤380) | `go test ./internal/domain/motor/... -run 'TestGemeloKNN(Gower|Imputation)'` | N/A: pure, unwired service | Statistics/Gower portions of both files |
| 5 | Selection/safety tests, evidence (≤330) | `go test ./internal/domain/motor/... -run TestGemeloKNN` | N/A: pure, unwired service | Selection portions of both files |

## Phase 1: SDD Review Slices

- [ ] 1.1 PR 1: review `exploration.md`, `proposal.md`, and `state.yaml`; preserve the #8 exclusion and scope boundary.
- [ ] 1.2 PR 2: review `specs/buyer-twin-knn/spec.md`, `design.md`, and this plan; retain all 13 scenarios.

## Phase 2: Contract-Safe Projection (PR 3)

- [ ] 2.1 RED: in `knn_test.go`, cover keyed/absent catalog zones, income A/B/C boundaries, all age brackets including `55+=60`, normalization, and missing optional values.
- [ ] 2.2 GREEN: in `knn.go`, add `EntradaGemelo`/`Vecino` and pure projection from records plus exact-ID catalog zones; forbid parallel inputs, name, slug, and price reads.

## Phase 3: Local Statistics and Gower (PR 4)

- [ ] 3.1 RED: test weighted renormalization, symmetry/identity, `[0,1]`, zero dependents, five-affiliate mode/odd-even median, below-threshold omission, and cross-project rejection.
- [ ] 3.2 GREEN: add immutable per-`ProyectoID` affiliate statistics, independent category/age imputation, and capped Gower weights `.35/.20/.15/.15/.15`.

## Phase 4: Neighbors, Safety, and Evidence (PR 5)

- [ ] 4.1 RED: test 30 ordered results, `K>n`, ID ties, `K<=0`/empty non-nil output, repeated/pure calls, and name/price inertness under fixed catalog zones.
- [ ] 4.2 GREEN: add non-mutating distance ordering and value-only `Vecino`; keep `capacidad.go`, `capacidad_test.go`, finance, pipeline, data, and Docs untouched.
- [ ] 4.3 Run focused tests per PR, then `go test ./internal/domain/motor/...`, `go test ./...`, `go build ./...`, ≥90% motor coverage, and the `origin/main...HEAD` path audit.
