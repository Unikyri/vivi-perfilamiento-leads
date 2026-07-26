```yaml
schema: gentle-ai.verify-result/v1
evidence_revision: sha256:450bd7b9e0073d0f915e0621ebfcd5d11d49440926b03ca5c90821f1c709d935
verdict: pass
blockers: 0
critical_findings: 0
requirements: 4/4
scenarios: 7/7
test_command: go test ./... -count=1
test_exit_code: 0
test_output_hash: sha256:aa6a223cd48268eb2e35ecd079287d988c0a5863cd019601100db50f37f78793
build_command: go build ./...
build_exit_code: 0
build_output_hash: sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
```

## Verification Report

**Change**: issue-29-greeting-consent
**Version**: N/A (delta specs: `saludar-lead`, `procesar-mensaje`)
**Mode**: Standard (`strict_tdd: false`, `openspec/config.yaml`)
**Round**: final verification after remediation of report #650 warnings 1-3
**Workspace**: `/tmp/vivi-issue-29`, branch `feat/issue-29-greeting-consent`, base `889b893`, candidate = uncommitted working tree (5 modified Go files, 0 new product files)
**Artifacts read**: OpenSpec `proposal.md`, `exploration.md`, `specs/{saludar-lead,procesar-mensaje}/spec.md`, `design.md`, `tasks.md`, `apply-progress.md`, `state.yaml`, prior `verify-report.md`; Engram `spec` (#641), `tasks` (#646), `design` (#640), `apply-progress` (#649, revision 2 with remediation), prior `verify-report` (#650)
**Authoritative spec counts**: 4 requirements, 7 scenarios (saludar-lead 3 req / 5 scen, procesar-mensaje 1 req / 2 scen)
**Toolchain**: go1.25.0 linux/amd64

### Completeness

| Metric | Value |
|--------|-------|
| Tasks total | 11 |
| Tasks complete | 11 |
| Tasks incomplete | 0 |
| Authored changed lines | 357 (321 additions + 36 deletions), budget 400 |

Native dispatcher agrees: `taskProgress 11/11`, `applyState: all_done`, `artifacts.applyProgress: done`.

### Prior Warning Closure (report #650)

| # | Prior warning | Remediation evidence | Status |
|---|---|---|---|
| 1 | Deterministic-template copy compliance not runtime-proven | `saludar_lead_test.go > TestSaludarLeadUsesOneValidatedProviderDraftOrFallback` now asserts `ValidarSaludo(persistedText, lead.Nombre, lead.Afiliado, wantAmount)` for all 4 cases, including the two deterministic-fallback paths. A template edit that broke name / amount / one-`?` / audience-marker / forbidden-prompt rules now fails the test. | ✅ CLOSED |
| 2 | `Denial does not complete a profile` lacked its GIVEN precondition | `procesar_mensaje_test.go > TestProcesarMensajeConsentDenialIsTerminalAndProfileSafe` seeds `recursos_propios` + `zona_deseada`, hard-asserts `PerfilEstaCompleto(before.Perfil)` before processing, and adds an explicit loop rejecting any `EvPerfilCompleto` event. | ✅ CLOSED |
| 3 | `formatoSubsidio` whole-million branch untested and deviated from `$X,YM` | `formatoSubsidio` now emits `fmt.Sprintf("$%d,%dM", millions, remainder/100_000)`; `TestFormatoSubsidioAlwaysUsesOneDecimalMillion` proves `52_000_000 -> $52,0M` and `52_500_000 -> $52,5M`. Integer arithmetic only, no floating point. | ✅ CLOSED |
| 4 | Native admission (`gentle-ai sdd-verify-validate`) unavailable | Still unavailable (see Verification Deviations). Process-only, not substantive. | ⚠️ OPEN (process) |

### Build & Tests Execution

```text
$ go test ./... -count=1   # exit 0 — 14 packages ok, 4 without test files, 0 failures, 0 skips
$ go build ./...           # exit 0, empty output
$ go vet ./...             # exit 0, empty output
$ gofmt -l <5 changed files>   # no output
```

Focused runtime evidence (`go test ./internal/usecase -run '<5 named tests>' -v -count=1`, exit 0):

```text
--- PASS: TestSaludarLeadUsesOneValidatedProviderDraftOrFallback (4/4 subtests)
--- PASS: TestFormatoSubsidioAlwaysUsesOneDecimalMillion (2/2 subtests: $52,0M, $52,5M)
--- PASS: TestValidarSaludoRejectsUnsafeDrafts (6/6 subtests)
--- PASS: TestProcesarMensajeConsentDenialIsTerminalAndProfileSafe
--- PASS: TestRechazarConsentimientoStopsAfterOrderedWriteFailures (3/3 subtests)
```

**Coverage**: ➖ Not available (threshold 0; no coverage command configured).

### Spec Compliance Matrix

| # | Requirement | Scenario | Test | Result |
|---|---|---|---|---|
| 1 | saludar-lead / Validated deterministic greeting | Provider draft is compliant | `saludar_lead_test.go > TestSaludarLeadUsesOneValidatedProviderDraftOrFallback/compliant_affiliate_draft` (persisted text == draft, `provider.calls == 1`) | ✅ COMPLIANT |
| 2 | saludar-lead / Validated deterministic greeting | Provider is unavailable or invalid | same test > `provider_error_uses_affiliate_fallback`, `invalid_draft_uses_non_affiliate_fallback`, `nil_provider_uses_deterministic_fallback` (1 message persisted, no retry) | ✅ COMPLIANT |
| 3 | saludar-lead / Audience-specific, single-question copy | Affiliate copy | same test > `provider_error_uses_affiliate_fallback` asserts the **persisted** affiliate template through `ValidarSaludo(..., true, "$52,5M")`; `TestFormatoSubsidioAlwaysUsesOneDecimalMillion` + `TestValidarSaludoRejectsUnsafeDrafts/{affiliate_accepted,wrong_amount_rejected,income_prompt_rejected}` | ✅ COMPLIANT |
| 4 | saludar-lead / Audience-specific, single-question copy | Non-affiliate copy | same test > `invalid_draft_uses_non_affiliate_fallback`, `nil_provider_uses_deterministic_fallback` assert the persisted job template through `ValidarSaludo(..., false, "")`; `TestValidarSaludoRejectsUnsafeDrafts/{job_copy_accepted,non_affiliate_amount_rejected,two_questions_rejected}` | ✅ COMPLIANT |
| 5 | saludar-lead / Consent refusal is terminal and profile-safe | Refusal farewell | `procesar_mensaje_test.go > TestProcesarMensajeConsentDenialIsTerminalAndProfileSafe` (`DESPEDIDO`/`DESPEDIDA`, 2 messages, 1 farewell, 0 events) + `TestRechazarConsentimientoStopsAfterOrderedWriteFailures` | ✅ COMPLIANT |
| 6 | procesar-mensaje / Consent denial precedes normal turn mutation | Denial with extracted fields | same denial test: malicious `ingreso_hogar=99999999` extracted field; `reflect.DeepEqual` over `Perfil`, `Capacidad`, `Intencion` on a cloning repository fake | ✅ COMPLIANT |
| 7 | procesar-mensaje / Consent denial precedes normal turn mutation | Denial does not complete a profile | same denial test with `PerfilEstaCompleto(before.Perfil)` precondition asserted and an explicit `EvPerfilCompleto` rejection loop | ✅ COMPLIANT |

**Compliance summary**: 7/7 scenarios compliant, 0 partial, 0 failing, 0 untested. All 4 requirements have complete passing runtime evidence.

### Correctness (Static + Runtime Evidence)

| Requirement | Status | Notes |
|---|---|---|
| At most one `GenerarTurno`, no retry | ✅ Implemented | Single call guarded by `uc.LLM != nil`; error/invalid draft falls through to the deterministic text. Runtime `provider.calls == 1`. |
| Validator contract | ✅ Implemented | Non-blank text, exact `URLPolitica`, exact clause `Al continuar autorizas el tratamiento de tus datos`, `strings.Count(text,"?") == 1`, lead name; rejects income/salary/household prompts; affiliate branch requires exact motor amount + `sueñ` and no second monetary expression; otherwise requires a job marker and no monetary expression. |
| Motor-only subsidy source, `$X,YM` shape | ✅ Implemented | Affiliate copy only when `Afiliado && Capacidad != nil && SubsidioAplicable > 0`; one-decimal COP-millions via integer division and modulo. |
| `URLPolitica` literal | ✅ Implemented | `https://www.colsubsidio.com/politica-tratamiento-datos`, single constant, no placeholder. |
| Denial ordering before mutation | ✅ Implemented | Branch immediately after provider output, before `normalizarCampos`, `aplicarCampos`, `motor.CalcularCapacidad`, intention assignment, `responder`, and `Bus.Publicar`. |
| No profile/capacity/intention mutation on denial | ✅ Implemented | `RechazarConsentimiento` touches only `Estado`, `Ruta`, `ActualizadoEn`; runtime deep-equality proof on a cloning fake. |
| Existing lifecycle edges only | ✅ Implemented | `PERFILANDO -> CALIFICADO -> DESPEDIDO`, one CAS `Guardar` after both transitions. |
| Ordered writes stop on failure | ✅ Implemented | Inbound -> CAS save -> farewell, each failure wrapped, no compensating write; 3 injected-failure subtests. |
| Composition wiring | ✅ Implemented | `cmd/servidor/main.go` shares one `SaludarLead` instance between the `LeadNuevo` observer and `ProcesarMensaje.Saludo`; coordinator registration unchanged. |
| Contract v1.1 surface | ✅ Unchanged | `git diff` on `internal/usecase/puertos.go`, `internal/adapters`, `internal/domain`, `migrations`, `data`, `openspec/specs` is empty. Exactly 5 files touched, none created or deleted. |

### Coherence (Design)

| Decision | Followed? | Notes |
|---|---|---|
| Greeting reuses `LLMProvider.GenerarTurno` + existing `EntradaTurno` | ✅ Yes | Controlled brief plus name, affiliation, capacity copy, and `NumerosDelMotor`; returned fields/action/intention ignored. |
| Acceptance boundary: fallback first, then validate draft with the same rules | ✅ Yes | `saludoDeterminista` runs before the provider call; invalid drafts never reach persistence. |
| Subsidy source and `$X,YM` formatting | ✅ Yes | Previously ⚠️ partial; now one decimal always and covered by a dedicated test. |
| Denial lifecycle | ✅ Yes | Implemented and runtime-proven. |
| Composition without coordinator-table change | ✅ Yes | `internal/adapters/agentes` suite green with observe-only registration. |
| Table-driven fakes only, no real provider | ✅ Yes | In-package fakes only; no network access in tests. |
| File changes limited to the 5 planned files | ✅ Yes | Confirmed by `git status` / `git diff --stat`. |

### Issues Found

**CRITICAL (substantive)**: None.

**PROCESS GATE (outside verification's authority)**

1. No bounded review authority exists for this change. `gentle-ai sdd-status issue-29-greeting-consent --cwd /tmp/vivi-issue-29 --json` reports `dependencies.verify: blocked`, `dependencies.archive: blocked`, `nextRecommended: resolve-review`, `blockedReasons: ["verify evidence cannot enter remediation: requirements are incomplete; bounded review transaction is missing"]` (the "requirements incomplete" clause referred to the superseded `requirements: 2/4` envelope of report #650 and is resolved by this report's `4/4`). `reviewState`, `reviewLedger`, and `reviewReceipt` are all `missing`; no lineage under `.git/gentle-ai/review-transactions/v2` references `issue-29`. Archive stays blocked until `gentle-ai review start` -> `review finalize` runs for this candidate. Verification cannot mint this authority.

**WARNING**

1. The zero-subsidy **affiliate** disjunct of "Non-affiliate copy" has no runtime case. `saludoDeterminista` routes `Afiliado && SubsidioAplicable == 0` to the amount-free job template and `ValidarSaludo` falls to the non-affiliate branch when `amount == ""`, so behaviour is correct by construction, but every test case uses `Afiliado: false`. One extra table row (`greetingLead(true, 0)`) would close it.
2. `gentle-ai sdd-verify-validate` does not exist in either installed binary (`.tools/gentle-ai` 2.1.11, `~/.local/bin/gentle-ai` 1.43.3), so strict-envelope admission was self-checked against `references/report-format.md` instead of machine-validated. Process-only; unchanged from report #650.
3. `test_output_hash` covers `go test` output that embeds per-package durations, so it is not byte-stable across runs. It proves this run, not reproducibility.

**SUGGESTION**

1. `ProcesarMensaje` returns `ErrValidacion` when `Saludo` is nil; the HTTP layer maps that to a 4xx, so a composition-root wiring bug would surface as a client error. A distinct internal error would diagnose better.
2. The greeting performs a synchronous provider call inside the `LeadNuevo` observer and `bus.EnMemoria.Publicar` dispatches synchronously, so lead creation can wait on the provider (bounded by `RequestTimeout` 8s / `LogicalTimeout` 10s, failures degrade to the deterministic text).
3. Task 2.2 says the refusal tests live in `saludar_lead_test.go`; they are in `procesar_mensaje_test.go`. Coverage exists — only the task text is stale.
4. The denial test asserts inbound author and text but not `MensajeID == "denial-1"`, so caller-supplied inbound metadata retention is only partially asserted.

### Verdict

**PASS WITH WARNINGS** — 11/11 tasks complete; 4/4 requirements and 7/7 spec scenarios fully compliant with passing runtime evidence; `go test ./... -count=1`, `go build ./...`, `go vet ./...`, and `gofmt` on all changed files are clean; report #650 warnings 1-3 are closed with committed tests and no behavioural regression; scope stays inside the 5 planned files and 357/400 authored lines with no Contract v1.1, motor, schema, data, adapter, or API change. Remaining items are one test-depth gap, one non-reproducible-hash note, and the unavailable native validator. Archive remains blocked by the missing bounded review transaction, which is a process gate this phase cannot satisfy.

**Next**: `gentle-ai review start` -> `review finalize` for this candidate, then `sdd-archive`.

### Verification Deviations

- Native admission could not run: `gentle-ai sdd-verify-validate` is absent from both installed binaries (2.1.11 and 1.43.3). The skill's hard rule prescribes zero persistence when the validator is unavailable; the orchestrator explicitly instructed a hybrid persist for this final verification, so this report was persisted with the deviation recorded here. It supersedes report #650 / Engram observation #650, whose four warnings are enumerated above so no prior finding is lost.
- The authority-only-failure recovery shape (`authority_only_failure`, exit `125`) was deliberately **not** emitted. Tests and build were executed under explicit orchestrator instruction, so truthful exit `0` values are recorded, and `gentle-ai` 2.1.11 rejects `authority_only_failure` as an unknown verify-result field (observed in report #651's `blockedReasons`). The review-authority gate is reported in prose and remains independently enforced by the native archive gate.
- Read-only guarantee: no product code, no commit, no push, no deploy, and no provider/network call. The only write is this report at `openspec/changes/issue-29-greeting-consent/verify-report.md` plus its Engram mirror.

### Post-Persist Routing Recheck

`gentle-ai sdd-status issue-29-greeting-consent --cwd /tmp/vivi-issue-29 --json` after persisting this report:

- `dependencies.verify: all_done` — the strict envelope above is now accepted by the native dispatcher (previously `blocked` on stale envelope fields).
- `dependencies.archive: blocked`, `nextRecommended: resolve-review`, `blockedReasons: ["multiple terminal native review receipts found; restore the change-local reviews/receipt.json mirror or remove stale terminal authority"]`.

The remaining gate is purely process/authority: 58 terminal review lineages exist under the shared Git common dir, none of them targets this change, and receipt discovery therefore requires explicit maintainer target selection. Substantive verification is complete; no code change can clear this gate.
