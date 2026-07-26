# Apply Progress: Nutrition Plan Slice 1

## Change
`issue-21-nutrition-plan` · Issue #21 · Slice 1 (planner and consented creation)

## Mode and Delivery
- Artifact store: hybrid (OpenSpec + Engram)
- Execution mode: standard (`strict_tdd: false`)
- Delivery strategy: auto-chain / feature-branch-chain
- Work unit: Slice 1 pre-review remediation only; Slice 2 remains unchecked

## Completed Tasks
- [x] 1.1 Added pure `motor.DisenarHitos` with UTC calendar expansion, stable date/source ordering, conversion-first `AFILIACION`, clamped gap allocation, four upcoming economic occurrences across annual rollover, and final `REEVALUACION`.
- [x] 1.2 Added table-driven planner tests for deterministic repeatability, calendar rollover/order, conversion-first behavior, small/odd/negative gaps, input immutability, and a positive gap requiring all four events across a year boundary.
- [x] 1.3 Added `GestionarPlan.CrearPlan` with explicit target/capacity validation, valid consent/frequency gate, no-consent reminder path, existing-plan retry reuse, and `Crear -> lead transition/save` ordering.
- [x] 1.4 Added local test doubles and discriminating tests for validation side effects, reminder success/failure, active timestamped plans, create failure ordering, lead-save failure, and retry without duplicate plan creation.

## Files Changed
| File | Action | Scope |
|---|---|---|
| `internal/domain/motor/plan.go` | Updated | Expand each valid calendar entry into recurring UTC occurrences before selecting the global four-event cap |
| `internal/domain/motor/plan_test.go` | Updated | Add annual-rollover coverage requiring four monetary milestones; simplify locally to preserve the line budget |
| `internal/usecase/gestionar_plan.go` | Created | `CrearPlan` only; pause deferred to Slice 2 |
| `internal/usecase/gestionar_plan_test.go` | Created | Local doubles and F2-F5 tests |
| `openspec/changes/issue-21-nutrition-plan/tasks.md` | Unchanged | Slice 1 tasks were already marked complete |
| `openspec/changes/issue-21-nutrition-plan/apply-progress.md` | Updated | Cumulative remediation evidence |
| `openspec/changes/issue-21-nutrition-plan/state.yaml` | Updated | Slice 1 authored line count and remediation state |

## Work Unit Evidence
| Evidence | Result |
|---|---|
| Focused tests | `go test ./internal/domain/motor ./internal/usecase -run 'TestDisenarHitos|TestGestionarPlanCrear' -count=1` — PASS |
| Package tests | `go test ./internal/domain/motor ./internal/usecase` — PASS |
| Full suite | Previously PASS before this focused remediation; not rerun because the requested scope is the two affected packages |
| Runtime harness | N/A — Slice 1 intentionally exposes callable services only; no HTTP, event wiring, or composition root is in scope |
| Authored line budget | 400 runtime/test lines total (`92 + 60 + 98 + 150`), at the hard Slice 1 limit |
| Rollback boundary | Revert the recurrence change and rollover test in the four runtime/test files plus this Slice 1 apply/docs update; no ports, existing use cases, adapters, migrations, or Slice 2 files were changed |

## Remediation Details
- `DisenarHitos` now emits up to four annual occurrences per valid source entry, sorts all generated occurrences by UTC date and original source order, then selects the earliest four monetary milestones globally.
- The discriminating test starts after the last event in a year and verifies events on `2027-01-15`, `2027-12-15`, `2028-01-15`, and `2028-12-15` for a positive gap.
- Caller calendar input remains unmodified; conversion-first and UTC normalization behavior remain unchanged.

## Native Attempt Evidence
Native runtime attempt attestation was unavailable. The transparent fallback command was:

```text
gentle-ai sdd-attempt status --cwd /tmp/vivi-issue-21 --change issue-21-nutrition-plan
exit status: 1
Error: unknown command "sdd-attempt" — run 'gentle-ai help' for available commands
```

No native attempt ordinal, budget attestation, or runtime evidence is claimed. Validation evidence above is direct command output.

## Deviations and Risks
- A pre-review remediation closes the annual-rollover gap: the planner accepts at most four monetary events strictly before `REEVALUACION`, and `TestDisenarHitosKeepsReevaluationLastByDate` proves chronological execution without inventing a later Contract date.
- The four-file authored line count is exactly 400, so later Slice 1 edits must simplify locally before adding lines.
- Cross-repository plan/lead persistence is intentionally non-atomic per the design; save-failure retry reuses the persisted plan through `PorLead`.

## Next
Slice 1 remediation is ready for review/verification. The overall change still has pending Slice 2 tasks (`2.1-2.3`), so the next apply phase is Slice 2 rather than final verification.
