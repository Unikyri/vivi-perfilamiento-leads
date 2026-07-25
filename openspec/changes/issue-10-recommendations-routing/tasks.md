# Tasks: Recommendations and 2x2 Routing — Issue #10

## Hard Dependency Gate

**Gate satisfied on 2026-07-25.** Issue #8 merged through PR #55 as `7a5b02a613c581568905936a35cb8feb841e62ab`. `git merge-base --is-ancestor 7a5b02a613c581568905936a35cb8feb841e62ab origin/main` and the same command against the rebased #10 `HEAD` both succeeded. Do not edit `internal/domain/motor/capacidad.go`, `internal/domain/motor/capacidad_test.go`, or any other #8 file.

## Review Workload Forecast

| Field | Value |
|---|---|
| Estimated changed lines | 520–650 (two production files, table-driven tests, SDD artifacts) |
| 400-line budget risk | High |
| Chained PRs recommended | Yes |
| Suggested split | PR 1 recommendations; PR 2 matrix |
| Delivery strategy | auto-chain: two sequential PRs to `main` |
| Chain strategy | stacked-to-main (PR 2 starts only after PR 1 merges) |

Decision needed before apply: No — resolved to two sequential slices
Chained PRs recommended: Yes
Chain strategy: stacked-to-main
400-line budget risk: High

### Suggested Work Units

| Unit | Goal | Likely PR | Focused test command | Runtime harness | Rollback boundary |
|---|---|---|---|---|---|
| 1 | Catalog recommendations | PR 1 | `go test ./internal/domain/motor -run '^TestRecomendarProyectos$'` | N/A: pure in-memory function | Revert `recomendar.go` + `recomendar_test.go` |
| 2 | 2x2 routing | PR 2 | `go test ./internal/domain/motor -run '^TestMatriz2x2$'` | N/A: pure in-memory function | Revert `matriz.go` + `matriz_test.go` |

The first slice targets `main`. The second slice starts only after the first one merges to `main`.

## Phase 1: Dependency and Scope Gate

- [x] 1.1 Record the #8 merge SHA and prove it is reachable from `origin/main` and `HEAD` descends from that `origin/main`; otherwise stop before PR/apply/verification.
- [x] 1.2 Before each slice, inspect the planned diff and reject any capacity, shared-domain, Docs, data, HTTP, use-case, or pipeline path.

## Phase 2: Recommendations (PR 1 after gate)

- [x] 2.1 Add RED table tests in `internal/domain/motor/recomendar_test.go` for REC-1..3: canonical-ID aggregation, orphan/empty IDs, non-nil empty output, cap three, and shuffled-input stable ordering.
- [x] 2.2 Extend those tests with exact `ceil(0.8*PrecioDesde)` pass/fail-by-one-peso cases, negative budget/non-positive price rejection, input immutability, and overflow-safe fraction tie/equality cases.
- [x] 2.3 Create `internal/domain/motor/recomendar.go` with `RecomendarProyectos`: catalog-only aggregation, integer eligibility, copied card data/reason, `math/bits` fraction comparison, total order, and top-three truncation.

## Phase 3: Matrix (PR 2 after gate)

- [x] 3.1 Add RED table tests in `internal/domain/motor/matriz_test.go` for MAT-1..7: all affiliate quadrants, exact 0.95/1.05 edges, non-affiliate conversion precedence, `MEDIA` as low, and empty/unknown intention.
- [x] 3.2 Extend matrix tests for NaN, infinities, and negative ratios never reaching `ASESOR`, plus route-only purity with no lead/ficha/cupo/milestone side effect.
- [x] 3.3 Create `internal/domain/motor/matriz.go` with `EntradaMatriz` and pure `Matriz2x2`: finite-ratio guard, conversion-first precedence, and the frozen route thresholds.

## Phase 4: Post-Gate Validation

- [x] 4.1 Run each focused slice test, then `go test ./internal/domain/motor/...` and `go test ./internal/domain/motor/... -cover` (target at least 90%).
- [x] 4.2 Run `go build ./...`; inspect the final diff for only the four planned motor files and no capacity-file change.
- [x] 4.3 Perform SDD verification only after 1.1 and all implementation tasks pass; attach command output to verification evidence.

Final verification completed on 2026-07-25 from merged SHA `d75c25ffca0ca5f4d95949cf40e73e5f57e856b8` (PR #59, identical to `origin/main`): `sdd/issue-10-recommendations-routing/verify-report` (`openspec/changes/issue-10-recommendations-routing/verify-report.md`) — PASS WITH WARNINGS, 6/6 requirements, 8/8 scenarios, `go test ./... -count=1` / `go build ./...` / `go vet ./...` all exit 0, motor coverage 97.4%, scope diff limited to the four additive motor files plus planning artifacts with no #8 capacity or `Docs/` change. This supersedes the earlier pre-commit worktree report (Engram `#499`). Warnings: native `sdd-verify-validate` and `gentle-ai review` authority are unavailable in `gentle-ai 1.43.3`, so envelope admission was manual and archive review authority rests on the merged PR chain (#56, #57, #58, #59); this evidence still needs a documentation commit to `main`. Archive not performed.
