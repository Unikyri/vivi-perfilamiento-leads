# Tasks: Nutrition Plan (Issue #21)

## Scope Guard
Implementation is limited to the six runtime/test files below; `gestionar_plan.*` is created in Slice 1 and extended in Slice 2. Do not modify ports, domain contracts, existing use cases, fakes, repositories, migrations, wiring, or frontend. SDD bookkeeping files may be updated for orchestration evidence.

## Review Workload Forecast

| Field | Value |
|---|---|
| Estimated changed lines | 645–755 total; each delivery slice MUST remain at or below 400 authored runtime/test lines |
| Delivery strategy / chain | auto-chain / feature-branch-chain |

Decision needed before apply: No
Chained PRs recommended: Yes
Chain strategy: feature-branch-chain
400-line budget risk: Medium

## Chained Slices and Frozen Budget Evidence

| Slice | Frozen review budget | Current authored runtime/test evidence | Proof and rollback |
|---|---:|---:|---|
| 1 — Plan creation | 400 lines maximum | **400 lines** in the existing committed Slice 1 baseline: `92 + 60 + 98 + 150`; no Slice 1 runtime/test files changed during remediation. | Focused planner/creation tests; runtime harness N/A because scope is callable-only. Revert Slice-1 commit; PR #1 targets tracker `feat/issue-21-nutrition-plan`. |
| 2 — Pause and tick | 400 lines maximum | **391 lines** after remediation across `gestionar_plan.go`, `gestionar_plan_test.go`, `ejecutar_hitos.go`, and `ejecutar_hitos_test.go`; within budget. | Focused pause/tick tests plus full/race/build/vet/format/diff checks; runtime harness N/A because wiring is intentionally deferred. Revert Slice-2 files only; PR #2 targets the Slice 1 branch. |

Both slice budgets are explicit hard limits. Planning artifacts and generated output are not included in authored runtime/test counts; the implementation remains on the feature-branch chain.

## Phase 1: Slice 1 — Planner and Consented Creation

- [x] 1.1 Create `internal/domain/motor/plan.go` with pure `DisenarHitos`: UTC calendar expansion/order, conversion-first `AFILIACION`, clamped gap allocation, and final `REEVALUACION` (F1).
- [x] 1.2 Add table-driven `plan_test.go` coverage for repeated deterministic output, rollover/order, conversion-first, zero/odd/small gaps, and non-mutated inputs (F1).
- [x] 1.3 Create `internal/usecase/gestionar_plan.go` `CrearPlan` with target/capacity, consent/frequency, existing-plan retry, and `Crear → state/save` ordering only (F2–F5).
- [x] 1.4 Add `gestionar_plan_test.go` cases for zero-side-effect validation, one no-consent reminder/error, timestamped active plan, create/save failure order, and retry without duplicate plan (F2–F5).

## Phase 2: Slice 2 — Pause and Due-Milestone Execution

- [x] 2.1 Extend `gestionar_plan.go` and its test with plan-save → lead-save → farewell pause ordering, absent/already-paused idempotency, illegal-transition error, and no farewell on failure (F7).
- [x] 2.2 Create `internal/usecase/ejecutar_hitos.go` to reject backward time, advance the clock, send → append → mark, aggregate independent failures, and retain pending failures (F6).
- [x] 2.3 Add `ejecutar_hitos_test.go` cases for paused silence, final-line pause copy without distress language, each send/append/mark failure, continued processing, retry count, threshold/reevaluation handoff, one event, and nil bus (F6–F8).
