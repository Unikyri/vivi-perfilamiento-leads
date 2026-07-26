# Tasks: Nutrition Plan (Issue #21)

## Scope Guard
Implementation is limited to the six files below; `gestionar_plan.*` is created in Slice 1 and extended in Slice 2. Do not modify `state.yaml`, ports, domain contracts, existing use cases, fakes, repositories, migrations, wiring, or frontend.

## Review Workload Forecast

| Field | Value |
|---|---|
| Estimated changed lines | 645–755 total; each slice 310–390 |
| Delivery strategy / chain | auto-chain / feature-branch-chain |

Decision needed before apply: No
Chained PRs recommended: Yes
Chain strategy: feature-branch-chain
400-line budget risk: Medium

## Chained Slices

| Slice | Files and count | Proof and native gate | Rollback and delivery |
|---|---|---|---|
| 1 — Plan creation | Create `internal/domain/motor/plan.go` (70–80 runtime), `plan_test.go` (95–110 tests), `internal/usecase/gestionar_plan.go` (75–90 runtime), `gestionar_plan_test.go` (95–110 tests): 145–170 runtime, 190–220 tests. | `go test ./internal/domain/motor ./internal/usecase -run 'TestDisenarHitos|TestGestionarPlanCrear'`; runtime harness N/A: callable-only scope. Run `gentle-ai review start`, `review finalize`, and `review validate --gate post-apply --cwd /tmp/vivi-issue-21`. | Revert Slice-1 commit. Commit `feat(nutrition): create consented deterministic plan`; PR #1 branch `feat/issue-21-nutrition-plan-s1` targets tracker `feat/issue-21-nutrition-plan`. |
| 2 — Pause and tick | Modify `gestionar_plan.go` (45–55 runtime) and its test (50–60 tests); create `ejecutar_hitos.go` (95–110 runtime) and test (120–140 tests): 140–165 runtime, 170–200 tests. | `go test ./internal/usecase -run 'TestGestionarPlanPausar|TestEjecutarHitos'`; runtime harness N/A: no wiring exists by design. Run the same native start/finalize/post-apply validation gate. | Revert Slice-2 commit only. Commit `feat(nutrition): pause and deliver due milestones`; PR #2 branch `feat/issue-21-nutrition-plan-s2` targets the PR #1 branch; retarget/rebase if its diff includes Slice 1. |

## Phase 1: Slice 1 — Planner and Consented Creation

- [x] 1.1 Create `internal/domain/motor/plan.go` with pure `DisenarHitos`: UTC calendar expansion/order, conversion-first `AFILIACION`, clamped gap allocation, and final `REEVALUACION` (F1).
- [x] 1.2 Add table-driven `plan_test.go` coverage for repeated deterministic output, rollover/order, conversion-first, zero/odd/small gaps, and non-mutated inputs (F1).
- [x] 1.3 Create `internal/usecase/gestionar_plan.go` `CrearPlan` with target/capacity, consent/frequency, existing-plan retry, and `Crear → state/save` ordering only (F2–F5).
- [x] 1.4 Add `gestionar_plan_test.go` cases for zero-side-effect validation, one no-consent reminder/error, timestamped active plan, create/save failure order, and retry without duplicate plan (F2–F5).

## Phase 2: Slice 2 — Pause and Due-Milestone Execution

- [ ] 2.1 Extend `gestionar_plan.go` and its test with plan-save → lead-save → farewell pause ordering, absent/already-paused idempotency, illegal-transition error, and no farewell on failure (F7).
- [ ] 2.2 Create `internal/usecase/ejecutar_hitos.go` to reject backward time, advance the clock, send → append → mark, aggregate independent failures, and retain pending failures (F6).
- [ ] 2.3 Add `ejecutar_hitos_test.go` cases for paused silence, final-line pause copy without distress language, each send/append/mark failure, continued processing, retry count, threshold/reevaluation handoff, one event, and nil bus (F6–F8).
