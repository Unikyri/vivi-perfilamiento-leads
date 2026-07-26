```yaml
schema: gentle-ai.verify-result/v1
evidence_revision: sha256:349469185fa6b87505dda1a3043cbfcf85762e6be6b39ed4406dd3949893b2f9
verdict: fail
blockers: 1
critical_findings: 1
requirements: 8/8
scenarios: 8/8
test_command: go test ./... -count=1
test_exit_code: 0
test_output_hash: sha256:97f3e75ebcda4bc4ea7eea531fa6bdf611d8246636b57755fb6aaaa7033f8bc6
build_command: go build ./...
build_exit_code: 0
build_output_hash: sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
```

## Verification Report

**Change**: issue-21-nutrition-plan (Issue #21)
**Version**: `plan-nutricion` delta spec — 8 requirements (F1–F8) / 8 scenarios, counted from the retrieved spec
**Mode**: Standard (`strict_tdd: false` in `openspec/config.yaml`; no TDD module loaded)
**Worktree**: `/tmp/vivi-issue-21` on `feat/issue-21-nutrition-plan` (HEAD `74c157a`, Slice 2 uncommitted)
**Artifact store**: hybrid — read in full from Engram (`#582` explore, `#584` proposal, `#585` spec, `#586` design, `#587` tasks, `#588` apply-progress, `#583` state, `#589` prior verify report) and the matching OpenSpec files; Engram and OpenSpec copies of spec, tasks, design, apply-progress and state agree
**Supersedes**: report `#589` / evidence revision `sha256:4bbf7fba…`, invalidated by the bounded F2/F7/F8 remediation

### Blocker: review authority, not implementation

`gentle-ai sdd-status` (native 2.1.11, the newer binary in `.tools/`) reports task completion but denies verify/archive dependency readiness:

```text
$ ./.tools/gentle-ai sdd-status issue-21-nutrition-plan --cwd /tmp/vivi-issue-21 --json
applyState: all_done            taskProgress: 7/7 (allComplete: true)
dependencies.verify: blocked    dependencies.archive: blocked
nextRecommended: resolve-review
blockedReasons: ["verify evidence cannot enter remediation: unknown verify result field
  validator_available; bounded review transaction is missing"]

$ ./.tools/gentle-ai review validate --gate post-apply --cwd /tmp/vivi-issue-21
result: invalidated   allowed: false   action: explicit-maintainer-action
reason: multiple terminal review receipts require explicit target selection
denial.stage: receipt-discovery
```

Two distinct causes are visible. The first — `unknown verify result field validator_available` — was caused by the prior report `#589` adding non-schema fields (`validator_available`, `validator_attested`, `admission_mode`) to the strict envelope; this report removes them and keeps validator transparency in prose only. The second — no bounded review transaction bound to this change plus an `invalidated` post-apply gate — is real and cannot be remedied by verification. `review status` shows the repository authority store (`…/.git/gentle-ai/review-transactions/v2/`) holding several approved lineages from earlier issues, so receipt discovery is ambiguous and no receipt matches the post-remediation tree.

This report therefore carries `verdict: fail` with `blockers: 1` even though every requirement, scenario and executed check is green: archive readiness is denied by review authority.

**Deviation from the authority-only recovery shape (stated explicitly, not hidden)**: the skill's authority-only preflight template requires `test_exit_code: 125` / `build_exit_code: 125` with empty-output hashes and unexecuted commands. Both commands *were* executed here and passed, so emitting `125` would misstate the evidence. The envelope records the real commands, real exit codes and real output digests, and omits `authority_only_failure`, `missing_review_authority`, `substantive_failure`, `command_failed` and `observed_authority_revision` rather than asserting a shape that does not describe what happened.

### Admission and validator transparency

No native admission attestation is claimed. `sdd-verify-validate` does not exist in either installed binary:

```text
$ ./.tools/gentle-ai --version                → gentle-ai 2.1.11
$ ./.tools/gentle-ai sdd-verify-validate --input … --requirements 8 --scenarios 8
Error: unknown command "sdd-verify-validate" — run 'gentle-ai help' for available commands
$ gentle-ai --version                         → gentle-ai 1.43.3   (PATH binary)
$ gentle-ai sdd-verify-validate --input … --requirements 8 --scenarios 8
Error: unknown command "sdd-verify-validate" — run 'gentle-ai help' for available commands
```

Admission is self-checked fallback only: a single leading envelope, every schema field present exactly once, no unknown fields, counts taken from the retrieved spec (8/8), and every command actually executed with recorded exit code and output digest. `evidence_revision` is the SHA-256 of a manifest containing the change name, `HEAD`, the 8/8 counts, both command hashes, and the SHA-256 of each of the six scoped runtime/test files. No validator attestation, no runtime-attempt attestation (`sdd-attempt` is absent from both binaries), and no review-gate attestation is claimed.

### Completeness
| Metric | Value |
|--------|-------|
| Tasks total | 7 |
| Tasks complete | 7 (1.1–1.4, 2.1–2.3 `[x]` in OpenSpec `tasks.md` and Engram `#587` rev 4) |
| Tasks incomplete | 0 |

### Build & Tests Execution

Go toolchain: `go1.25.0 linux/amd64`.

```text
$ go test ./... -count=1                                   # exit 0
ok  …/cmd/servidor  …/internal/domain  …/internal/domain/motor  …/internal/infrastructure/config
ok  …/internal/infrastructure/ids  …/internal/infrastructure/llm  …/internal/infrastructure/postgres
ok  …/internal/pipeline  …/internal/usecase        (4 packages report "no test files")
$ go test -race ./... -count=1                             # exit 0, all packages ok
$ go build ./...                                           # exit 0, empty output
$ go vet ./...                                             # exit 0
$ go mod verify                                            # exit 0, "all modules verified"
$ go mod tidy -diff                                        # exit 0, no go.mod/go.sum drift
$ gofmt -l internal/domain/motor internal/usecase          # empty
$ git diff --check                                         # exit 0
$ go test ./internal/domain/motor ./internal/usecase \
    -run 'TestDisenarHitos|TestGestionarPlan|TestEjecutarHitos' -v -count=1   # exit 0
10 test functions (4 planner + 6 use case) + 9 subtests: PASS, 0 FAIL, 0 SKIP
```

`gofmt -l internal cmd` also lists `internal/pipeline/compradores_test.go` and `internal/pipeline/proyectos_test.go`. Both are pre-existing drift from commit `4533b46`, untouched by this change (`git diff --stat -- internal/pipeline` is empty) and outside its scope guard.

**Coverage**: ➖ Not available — `coverage_threshold: 0` and no coverage command configured.

### Remediation Verification (independently confirmed)
| Remedy | Runtime evidence | Result |
|---|---|---|
| Below-budget plan yields `MetaMonto == 0` at the `CrearPlan` boundary | `gestionar_plan_test.go > TestGestionarPlanCrearConsentAndTimestamp`: target `50` vs `PresupuestoMax 100` → `MetaMonto == 0`, `creates == 1`, `len(messages) == 0`, lead `EN_NUTRICION`; runtime clamps `gap = PrecioObjetivo - PresupuestoMax`, `if gap < 0 { gap = 0 }` | ✅ Confirmed |
| Zero meta produces no threshold handoff, but `REEVALUACION` still hands off | `ejecutar_hitos_test.go > TestEjecutarHitosHandoffOnceAndNilBus`: `MetaMonto = 0` with a notified `Monto 1` milestone → count 1, lead stays `EN_NUTRICION`, `len(published) == 0`; the `REEVALUACION` case still reaches `PERFILANDO`. `requiereReperfilado` checks `TipoHitoReevaluacion` first, then returns false for `MetaMonto <= 0` | ✅ Confirmed |
| Active plan + already-`PAUSADO` lead becomes `PAUSADO` without a farewell | `gestionar_plan_test.go > TestGestionarPlanPausarOrderingAndIdempotency` third case: plan `ACTIVO`, lead already `PAUSADO` → plan persisted `PAUSADO`, `order == ["plan-save"]`, `len(messages) == 0`, nil error. Required because `domain.transiciones` has no `PAUSADO → PAUSADO` edge (`internal/domain/estado.go`) | ✅ Confirmed |

### Spec Compliance Matrix
| Requirement | Scenario | Test | Result |
|-------------|----------|------|--------|
| F1 Provider-free deterministic planning | Deterministic, isolated plan | `motor/plan_test.go > TestDisenarHitosDeterministicCalendarPlan` (repeat equality, calendar dates `2026-02-14`/`2026-06-30`, input slice unmutated), `TestDisenarHitosConversionFirstAndRollover` (`AFILIACION` at `desde+8d`), `TestDisenarHitosKeepsReevaluationLastByDate`, `TestDisenarHitosSmallOddAndNegativeGaps`; isolation from `go list -deps ./internal/domain/motor` and `plan.go` importing only `sort/strconv/strings/time` + `internal/domain` | ✅ COMPLIANT |
| F2 Explicit target and gap validation | Target validation and clamped gap | `TestGestionarPlanCrearValidationHasNoSideEffects/{missing capacity,non-positive target}` (ErrValidacion, zero creates/messages, state unchanged) **and** the remediated below-budget case in `TestGestionarPlanCrearConsentAndTimestamp` (`MetaMonto == 0`) plus above-budget `MetaMonto == 100` | ✅ COMPLIANT |
| F3 Consent and frequency gate | Valid consent gate | `TestGestionarPlanCrearConsentAndTimestamp` (`ACTIVO`, `ConsentimientoEn == Reloj.Ahora()`), `TestGestionarPlanCrearValidationHasNoSideEffects/invalid_frequency` (no create, no message, state unchanged) | ✅ COMPLIANT |
| F4 No-consent reminder | Consent declined or absent | `TestGestionarPlanCrearNoConsentReminderAndFailure` (exactly one repository interaction, `(nil, nil)` on success, append error propagated, zero creates) | ✅ COMPLIANT |
| F5 Durable, idempotent creation ordering | Cross-repository creation failure | `TestGestionarPlanCrearFailureOrderAndRetryReuse` (order `create`→`save`; create failure leaves lead `CALIFICADO` with order `["create"]`; save failure then retry keeps `creates == 2`, reuses the plan via `PorLead`, final order `["create","create","save","save"]`) | ✅ COMPLIANT |
| F6 Due-milestone delivery ordering | Partial tick failure | `TestEjecutarHitosOrdersAndUsesDignifiedPauseCopy` (backward-time rejection, events `["send","append","mark"]`, count 1, final line offers PAUSAR, no distress substrings), `TestEjecutarHitosContinuesFailuresAndRetriesPending` (send/append/mark failure each leaves 0 marked, both due milestones still attempted, aggregated error) | ✅ COMPLIANT |
| F7 Immediate, respectful pause | Pause variants | `TestGestionarPlanPausarOrderingAndIdempotency` (order `plan-save`→`save`→`message`, one farewell, repeat is a no-op, already-paused lead still normalizes the plan without farewell), `TestGestionarPlanPausarMissingAndFailuresHaveNoFarewell` (missing plan silent; illegal transition errors after `plan-save`; lead-save failure without farewell); paused silence via the tick's `lead.Estado == PAUSADO` skip plus `postgres.HitosVencidos` filtering `EstadoPlanActivo` + `EstadoHitoPendiente` | ✅ COMPLIANT |
| F8 Requalification handoff | One deterministic handoff | `TestEjecutarHitosHandoffOnceAndNilBus` (60+40 ≥ `MetaMonto` 100 → lead saved `PERFILANDO` before exactly one `PerfilCompleto`; zero-meta case does not hand off; `REEVALUACION` with nil bus still transitions and publishes nothing); no `CalificarLead`/`RutaDecidida` reference in the three runtime files (grep) | ✅ COMPLIANT |

**Compliance summary**: 8/8 scenarios compliant, 0 partial, 0 untested, 0 failing.

### Correctness (Static Evidence)
| Requirement | Status | Notes |
|------------|--------|-------|
| F1 | ✅ Implemented | Pure `motor.DisenarHitos`: UTC normalization, canonical `--MM-DD` expansion with annual rollover, stable date-then-source sort, ≤4 monetary milestones strictly before reevaluation, conversion-first `AFILIACION`, halving allocation with a 500 000 floor. |
| F2 | ✅ Implemented | `PrecioObjetivo <= 0 \|\| Capacidad == nil` → `ErrValidacion` before any write; gap clamped at 0 and persisted as `MetaMonto`. |
| F3 | ✅ Implemented | Frequency allowlist `QUINCENAL/MENSUAL/TRIMESTRAL`; plan persisted `ACTIVO` with `ConsentimientoEn = Reloj.Ahora()`. |
| F4 | ✅ Implemented | Single `Leads.AgregarMensaje` door-open reminder, `return nil, nil`; append error wrapped and returned. |
| F5 | ✅ Implemented | `PorLead` retry lookup → `Planes.Crear` → conditional lead transition/save; existing plan reused instead of duplicated. |
| F6 | ✅ Implemented | `Gateway.Enviar` → `Leads.AgregarMensaje` → `Planes.MarcarHito(NOTIFICADO)` with `continue` on failure, `errors.Join` aggregation, only fully marked deliveries counted. |
| F7 | ✅ Implemented | `Planes.Guardar(PAUSADO)` → optional lead transition/save → one farewell; already-paused plan and missing plan return nil; already-paused lead normalizes the plan and stops before the farewell. |
| F8 | ✅ Implemented | `requiereReperfilado` matches `REEVALUACION` first, ignores `MetaMonto <= 0`, otherwise sums notified amounts (including the current milestone when absent from the stored plan); `handoff` map caps one `PerfilCompleto` per lead per tick; nil bus skipped. |

### Coherence (Design)
| Decision | Followed? | Notes |
|----------|-----------|-------|
| Pure planner in `motor`, orchestration in `usecase` | ✅ Yes | `plan.go` has no repository, clock or port dependency. |
| Contract v1.1 names/types (`TipoHito*`, `EstadoHito*`, `EventoCalendario.Fecha string`) | ✅ Yes | `git diff --stat -- internal/domain internal/usecase/puertos.go` is empty. |
| No new ports or repository methods | ✅ Yes | Only pre-existing members used (`PorID`, `PorLead`, `Crear`, `Guardar`, `AgregarMensaje`, `HitosVencidos`, `MarcarHito`, `Enviar`, `Publicar`, `Ahora`, `Avanzar`, `FechaSimulada`). |
| Callable-only, provider-free integration | ✅ Yes | No `net/http`, ADK, `genai`, prompt, `Suscribir` or `TickReloj` reference in the three runtime files; no `cmd/`, `internal/adapters`, `internal/infrastructure` or `web/` reference to `GestionarPlan`/`EjecutarHitos`. |
| Lead transition only from legacy `CALIFICADO` | ✅ Yes (documented deviation) | Design records that `CalificarLead` already persists `EN_NUTRICION`, so F5's transition is a reconciliation path; `["create","save"]` order preserved when it applies. |
| Pause order plan-save → lead-save → farewell | ✅ Yes | Proven by asserted `order` slices. |
| Nutrition-local test doubles, no shared fake refactor | ✅ Yes | `planTestState`/`tickState` doubles are file-local; existing `fakes_test.go` untouched. |
| Runtime harness N/A | ✅ Yes | Wiring is intentionally deferred to #22/#23/#24; no runtime-attempt attestation claimed. |

### Boundary, Non-Goal and Budget Verification
| Check | Result |
|---|---|
| Frozen non-goals (LLM/ADK/prompting, HTTP endpoints, concrete `BusEventos`/`Reloj`, event subscription/coordinator wiring, frontend, migrations, Postgres repository changes) | ✅ None touched — `git diff --numstat origin/main...HEAD` plus working tree limited to the six runtime/test files and the change folder |
| Contract/domain/ports unchanged; `CalificarLead`/`ProcesarMensaje` unedited | ✅ Confirmed by empty diffs for `internal/domain`, `internal/usecase/puertos.go`, `calificar_lead.go`, `procesar_mensaje.go` |
| Scope guard (six files only) | ✅ Modified: `gestionar_plan.go`, `gestionar_plan_test.go`; new: `ejecutar_hitos.go`, `ejecutar_hitos_test.go`; Slice 1 `motor/plan*.go` unchanged during remediation |
| Slice 1 authored runtime/test budget = 400 (hard cap 400) | ✅ Exactly 400 — `git diff --numstat origin/main...HEAD` gives `92 + 60 + 98 + 150` in commit `74c157a`; at the cap, not over |
| Slice 2 authored runtime/test budget ≤ 400 | ✅ 391 — `gestionar_plan.go` 49+0, `gestionar_plan_test.go` 63+3, `ejecutar_hitos.go` 120, `ejecutar_hitos_test.go` 156 |
| Hybrid artifact parity | ✅ Engram `#585/#586/#587/#588/#583` match their OpenSpec files; the prior `tasks.md` ↔ `#587` drift reported in `#589` is resolved |

### Issues Found
**CRITICAL**:
1. **Archive gate blocked by review authority.** No bounded review transaction is bound to this change and `review validate --gate post-apply` returns `result: invalidated` (`multiple terminal review receipts require explicit target selection`, stage `receipt-discovery`). Native `sdd-status` routes to `resolve-review`. The remediation also moved the tree past the receipt identity cited in `#589`, so the prior lineage cannot be reused. This is the single blocker and requires explicit maintainer/orchestrator action; verification cannot resolve it and did not attempt to.

**WARNING**:
1. Slice 1 sits at exactly 400 authored runtime/test lines. `tasks.md` frames the budget as "400 lines maximum" (satisfied), while proposal decision F9 says "each under 400 authored changed lines" (not satisfied at equality). Reconcile the wording or accept the equality explicitly. Counting the planning docs committed in `74c157a`, PR #1 is 875 changed lines under the shared authored-text rule.
2. Slice 2 is still uncommitted (2 modified + 2 untracked runtime/test files, plus 3 modified and 1 untracked SDD bookkeeping files). Any commit/PR gate must stage exactly the reviewed paths without content change, and the frozen intended-untracked set must move together.
3. Two installed `gentle-ai` binaries disagree (`.tools/` 2.1.11 vs PATH 1.43.3). Only 2.1.11 exposes `review`/`sdd-status` review awareness, and neither exposes `sdd-verify-validate`. Orchestration should pin one binary before archive gating.
4. Cross-repository plan/lead writes remain intentionally non-atomic, and a gateway send followed by an append/mark failure can duplicate outbound delivery on the next tick (no gateway idempotency contract). Accepted design risk, bounded by ordering tests.
5. Pre-existing `gofmt` drift in `internal/pipeline/compradores_test.go` and `proyectos_test.go` is unrelated to this change but will keep showing in repository-wide format checks.

**SUGGESTION**:
1. `ejecutar_hitos.go`, `ejecutar_hitos_test.go` and `motor/plan_test.go` place the module import inside the stdlib import group; `gofmt` accepts it but the rest of the package separates them.
2. The "no financial distress" assertion rejects only the substrings `deuda` and `distres`; a positive assertion on the approved copy would be stronger.
3. `REEVALUACION` handoff is proven with `MetaMonto = 100` and a nil bus; a case combining `MetaMonto = 0` with `REEVALUACION` would pin the precedence order in `requiereReperfilado` directly.
4. No test composes `PausarPlan` with `EjecutarHitos`; paused silence is proven through the resulting lead state plus the repository's active-plan filter.

### Verdict
FAIL
All 7 tasks complete and all 8 requirements / 8 scenarios are compliant with green tests, race, build, vet, module and format evidence; the sole blocker is missing/ambiguous bounded review authority, which denies archive readiness and must be resolved by the orchestrator (`resolve-review`), not by re-implementation.
