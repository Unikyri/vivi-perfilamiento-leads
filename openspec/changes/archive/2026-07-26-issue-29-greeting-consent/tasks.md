# Tasks: Personalized Greeting and Data Consent

## Review Workload Forecast

| Field | Value |
|---|---|
| Estimated changed lines | 260–330 authored lines |
| 400-line budget risk | Low |
| Chained PRs recommended | No |
| Suggested split | Single Block A PR |
| Delivery strategy | single-pr |
| Chain strategy | pending |

Decision needed before apply: No
Chained PRs recommended: No
Chain strategy: pending
400-line budget risk: Low

### Suggested Work Units

| Unit | Goal | Likely PR | Focused test command | Runtime harness | Rollback boundary |
|---|---|---|---|---|---|
| 1 | Safe selected greeting | Single PR | `go test ./internal/usecase -run TestSaludarLead` | N/A—event observer verified with a fake provider | greeting source/tests |
| 2 | Terminal consent denial | Single PR | `go test ./internal/usecase -run TestProcesarMensaje.*Consent` | N/A—use-case test proves no event/writes after failure | denial branch, tests, composition wiring |

## Phase 1: Validated Greeting

- [x] 1.1 In `internal/usecase/saludar_lead.go`, add the provider dependency and construct one controlled `EntradaTurno`; call `GenerarTurno` at most once and persist the accepted response or fallback.
- [x] 1.2 In `internal/usecase/saludar_lead.go`, define `URLPolitica = "https://www.colsubsidio.com/politica-tratamiento-datos"`, exact templates, integer `$X,YM` formatting, and a pure validator for consent, one `?`, audience markers, and forbidden fields/amounts.
- [x] 1.3 In `internal/usecase/saludar_lead_test.go`, table-test compliant affiliate/non-affiliate drafts, real/zero subsidy, and fallback for nil/error/empty/invalid drafts; assert one provider call and persisted-copy parity.

## Phase 2: Profile-Safe Consent Refusal

- [x] 2.1 In `internal/usecase/saludar_lead.go`, implement `RechazarConsentimiento`: retain inbound denial, traverse `PERFILANDO → CALIFICADO → DESPEDIDO`, set `DESPEDIDA`, CAS-save once, append one fixed farewell, and stop on each write failure.
- [x] 2.2 In `internal/usecase/saludar_lead_test.go`, test refusal state/route/farewell plus add/save failure ordering, with profile, capacity, and intention unchanged.
- [x] 2.3 In `internal/usecase/procesar_mensaje.go`, inject the greeting collaborator and branch on `CONSENTIMIENTO_NO` immediately after provider output, before normalization, field merge, motor calculation, response, or publication.
- [x] 2.4 In `internal/usecase/procesar_mensaje_test.go`, use malicious extracted fields to prove retained inbound metadata, immutable lead data, `DESPEDIDO`/`DESPEDIDA`, exactly one farewell, and zero `PerfilCompleto` events.

## Phase 3: Composition and Regression Boundary

- [x] 3.1 In `cmd/servidor/main.go`, inject the existing provider into `SaludarLead` and the same greeting instance into `ProcesarMensaje`; leave coordinator registration, ports, adapters, DTOs, lifecycle, and migrations unchanged.
- [x] 3.2 Run `go test ./internal/adapters/agentes/...` to preserve the existing `LeadNuevo` observe-only registration contract.

## Phase 4: Verification

- [x] 4.1 Run `go test ./internal/usecase/... ./internal/adapters/agentes/... ./cmd/...` after the focused tests pass.
- [x] 4.2 Run `go test ./...`, `go build ./...`, and `go vet ./...`; record failure output before any corrective work.