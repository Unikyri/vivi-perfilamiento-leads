# Apply Progress: Nutrition Plan Remediation

## Change
`issue-21-nutrition-plan` · Issue #21 · focused remediation after verify report #589

## Mode and Delivery
- Artifact store: hybrid (OpenSpec + Engram)
- Execution mode: standard (`strict_tdd: false`)
- Delivery strategy: auto-chain / feature-branch-chain
- Work unit: focused remediation of F2/F8 and F7 warnings
- Scope: `internal/usecase/gestionar_plan.go`, `gestionar_plan_test.go`, `ejecutar_hitos.go`, `ejecutar_hitos_test.go`, and SDD bookkeeping only

## Completed Tasks
- [x] 1.1 Added pure `motor.DisenarHitos` with UTC calendar expansion, stable date/source ordering, conversion-first `AFILIACION`, clamped gap allocation, annual rollover, and final `REEVALUACION`.
- [x] 1.2 Added table-driven planner tests for deterministic repeatability, calendar rollover/order, conversion-first behavior, small/odd/negative gaps, input immutability, and annual rollover.
- [x] 1.3 Added `GestionarPlan.CrearPlan` with explicit target/capacity validation, consent/frequency gate, no-consent reminder, existing-plan retry reuse, and `Crear → lead transition/save` ordering.
- [x] 1.4 Added creation tests for validation side effects, no-consent reminder/error, timestamped active plans, create/save failure order, retry reuse, and a valid below-budget target producing `MetaMonto == 0` without a reminder or invalid state change.
- [x] 2.1 Added `GestionarPlan.PausarPlan` with plan-save → lead transition/save → farewell ordering, missing/already-paused idempotency, illegal-transition errors, no farewell after persistence failure, and active-plan normalization when the lead is already `PAUSADO` without a farewell.
- [x] 2.2 Added `EjecutarHitos.Ejecutar` with monotonic simulated time, due active/pending milestone delivery, ordered gateway/message/mark operations, independent error aggregation, pending retry semantics, and deterministic profile handoff.
- [x] 2.3 Added execution tests for paused silence, dignified copy, send/append/mark failures, continued processing, retries, positive-target threshold handoff, zero-meta threshold non-handoff, reevaluation handoff, one event, and nil bus.

## Remediation Evidence
- `CrearPlan` still clamps `gap = max(0, PrecioObjetivo - PresupuestoMax)`; the below-budget test uses target `50` against budget `100`, asserts `MetaMonto == 0`, exactly one plan create, no reminder, and normal lead transition.
- `requiereReperfilado` checks `REEVALUACION` first, then returns false for `MetaMonto <= 0`; a positive monetary milestone with zero meta leaves the lead `EN_NUTRICION`, while `REEVALUACION` still transitions it to `PERFILANDO`.
- `PausarPlan` persists an active plan as `PAUSADO` before the already-paused lead returns silently; the normal active-lead path retains plan-save → lead-save → farewell ordering and repeat idempotency.

## Work Unit Evidence
| Evidence | Result |
|---|---|
| Focused tests | `go test ./internal/domain/motor ./internal/usecase -run 'TestDisenarHitos|TestGestionarPlan|TestEjecutarHitos' -v -count=1` — PASS |
| Full suite | `go test ./... -count=1` — PASS |
| Race suite | `go test -race ./... -count=1` — PASS |
| Build | `go build ./...` — PASS |
| Vet | `go vet ./...` — PASS |
| Modules | `go mod verify` — PASS (`all modules verified`) |
| Formatting/diff | `gofmt -l internal/domain/motor internal/usecase` empty; `git diff --check` — PASS |
| Runtime harness | N/A — callable-only scope; HTTP, event wiring, and composition root remain intentionally out of scope |
| Slice 1 authored runtime/test budget | **400 lines**, existing committed baseline `92 + 60 + 98 + 150`; unchanged during this remediation and at the 400-line hard cap |
| Slice 2 authored runtime/test budget | **391 lines**, measured across the four scoped runtime/test files; below the 400-line hard cap |
| Rollback boundary | Revert only the four scoped usecase runtime/test files and the related SDD bookkeeping; no ports, domain contract, repositories, integrations, migrations, HTTP, LLM/ADK, wiring, or frontend changed |

## Native Attempt Evidence
Native runtime attempt attestation remains unavailable in the installed toolchain. The prior transparent fallback was:

```text
gentle-ai sdd-attempt status --cwd /tmp/vivi-issue-21 --change issue-21-nutrition-plan
exit status: 1
Error: unknown command "sdd-attempt" — run 'gentle-ai help' for available commands
```

No native attempt ordinal, budget attestation, or runtime evidence is claimed. Validation evidence above is direct command output.

## Deviations and Risks
- None from the scoped design; the implementation remains callable-only and provider-free.
- Cross-repository plan/lead persistence remains intentionally non-atomic.
- Gateway send success followed by append or mark failure can duplicate outbound delivery on retry because the existing gateway port has no idempotency contract.
- The prior verify report #589 is historical evidence and must be refreshed by `sdd-verify` after this remediation before archive.

## Next
All implementation tasks (`1.1–1.4`, `2.1–2.3`) remain complete. Ready for `sdd-verify` to refresh the strict evidence and archive gate.
