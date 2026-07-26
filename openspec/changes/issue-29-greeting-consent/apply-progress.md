# Apply Progress: Issue 29 — Greeting and Consent

## Mode

Standard mode (`strict_tdd: false`). The repository configuration disables strict TDD; focused behavioral tests and the required post-implementation checks were run.

## Completed Tasks

- [x] 1.1 Added the optional provider dependency, controlled `EntradaTurno`, one-call provider attempt, and deterministic fallback persistence.
- [x] 1.2 Added `URLPolitica`, motor-derived integer subsidy formatting, exact consent templates, and pure greeting validation.
- [x] 1.3 Added table-driven provider, fallback, validation, and persisted-copy tests.
- [x] 2.1 Added ordered denial persistence, existing lifecycle transitions, `DESPEDIDA`, CAS save, and fixed farewell.
- [x] 2.2 Added state, route, farewell, profile immutability, and write-failure ordering tests.
- [x] 2.3 Injected the greeting collaborator and branched on `CONSENTIMIENTO_NO` before normalization and motor mutation.
- [x] 2.4 Added malicious extracted-field denial coverage with zero completion events.
- [x] 3.1 Wired the same provider and `SaludarLead` instance through `cmd/servidor/main.go`; coordinator registration remains unchanged.
- [x] 3.2 Preserved and passed the observe-only `LeadNuevo` coordinator test suite.
- [x] 4.1 Passed focused use-case, coordinator, and command tests.
- [x] 4.2 Passed the full Go suite, build, and vet.

## Work Unit Evidence

| Work unit | Focused test command and result | Runtime harness | Rollback boundary |
|---|---|---|---|
| Safe selected greeting | `go test ./internal/usecase -run 'TestSaludarLead|TestValidarSaludo'` — PASS | N/A: event observer is proven with a fake provider and repository | Revert `saludar_lead.go` and `saludar_lead_test.go` greeting changes |
| Terminal consent denial | `go test ./internal/usecase -run 'TestProcesarMensajeConsent|TestRechazarConsentimiento'` — PASS | N/A: use-case repository/bus harness proves ordered writes and no publication | Revert denial branch, greeting collaborator wiring, and related tests |
| Regression boundary | `go test ./internal/usecase/... ./internal/adapters/agentes/... ./cmd/...` — PASS | N/A: no external runtime provider is contacted by tests | Revert the five source/test files in this change |

## Validation

- `go test ./internal/usecase/... ./internal/adapters/agentes/... ./cmd/...` — PASS
- `go test ./...` — PASS
- `go build ./...` — PASS
- `go vet ./...` — PASS

## Files Changed

- `internal/usecase/saludar_lead.go`
- `internal/usecase/saludar_lead_test.go`
- `internal/usecase/procesar_mensaje.go`
- `internal/usecase/procesar_mensaje_test.go`
- `cmd/servidor/main.go`
- `openspec/changes/issue-29-greeting-consent/tasks.md`
- `openspec/changes/issue-29-greeting-consent/state.yaml`
- `openspec/changes/issue-29-greeting-consent/apply-progress.md`

## Deviations and Risks

None from the approved design. Provider failures, nil providers, empty drafts, and invalid drafts intentionally select the deterministic fallback; no provider retry is performed. No external provider was contacted during local validation.

## Status

11/11 tasks complete. Ready for `sdd-verify`.


## Remediation: Verify Report #650 Warnings 1–3

- [x] Added persisted-greeting assertions that call the production `ValidarSaludo` validator for accepted drafts and deterministic fallback cases.
- [x] Strengthened consent-denial coverage with a structurally complete profile and an explicit assertion that no `PerfilCompleto` event is published.
- [x] Changed subsidy formatting to always emit one decimal COP-millions (`$52,0M` and `$52,5M` cases covered).

## Remediation Work Unit Evidence

| Work unit | Focused test command and result | Runtime harness | Rollback boundary |
|---|---|---|---|
| Report #650 warnings 1–3 | `go test ./internal/usecase -run 'TestSaludarLeadUsesOneValidatedProviderDraftOrFallback|TestFormatoSubsidioAlwaysUsesOneDecimalMillion|TestProcesarMensajeConsentDenialIsTerminalAndProfileSafe' -count=1` — PASS | N/A: deterministic validator, repository, and in-memory bus tests cover the runtime boundaries without external providers | Revert the formatter and the three focused test assertions in `internal/usecase/` |

## Remediation Validation

- `go test ./internal/usecase/... ./internal/adapters/agentes/... ./cmd/... -count=1` — PASS
- `go test ./... -count=1` — PASS
- `go build ./...` — PASS
- `go vet ./...` — PASS

## Remediation Status

Report #650 warnings 1–3 addressed. Original task progress remains 11/11 complete; ready for `sdd-verify`.
