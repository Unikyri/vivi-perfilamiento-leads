# Tasks: Contract v1.1 HTTP API and Demo Controls

## Review Workload Forecast

| Field | Value |
|---|---|
| Estimated authored runtime/test lines | 2,070 across six slices |
| 400-line budget risk | High overall; each slice <=400 |
| Chained PRs recommended | Yes |
| Delivery / chain | auto-chain / feature-branch-chain |
| Bases | S1→`feat/issue-23-http-api`; S2→S1; …; S6→S5 |

Decision needed before apply: No
Chained PRs recommended: Yes
Chain strategy: feature-branch-chain
400-line budget risk: High

### Suggested Work Units

| Unit | Status / cap | Focused test | Runtime harness | Rollback |
|---|---|---|---|---|
| S1 | READY / 390 | `go test ./internal/adapters/http ./internal/usecase` | `httptest` create + conversation | S1 files |
| S2 | READY / 320 | `go test ./internal/adapters/http ./internal/usecase` | `httptest` 202 then poll | S2 files |
| S3 | READY / 360 | `go test ./internal/adapters/http ./internal/usecase` | `httptest` queue/detail/ficha | S3 files |
| S4 | READY / 330 | `go test ./internal/adapters/http ./internal/usecase` | `httptest` persona filter/catalog | S4 files |
| S5 | READY / 370 | `go test ./internal/... ./cmd/servidor` | `httptest` time tick + health | S5 files |
| S6 | **BLOCKED** / 300 | `go test ./internal/...` | reset twice; preserve buyers | S6 files |

## Phase 1: S1 HTTP Foundation (READY, <=390)
- [x] 1.1 Add `internal/adapters/http/leads_test.go` and `internal/usecase/saludar_lead_test.go` contract tables for create, greeting, conversation, strict JSON, methods, and private errors.
- [x] 1.2 Create `errores.go`, `rutas.go`, `leads.go`, and `saludar_lead.go`; wire only Bloque A in `cmd/servidor/main.go`.

## Phase 2: S2 Asynchronous Turns (READY, <=320)
- [x] 2.1 Add `internal/adapters/http/turnos_test.go` cases for valid 202, active-to-clear polling, cancellation, text limit, MIME, duration, decoded-size, and no audio echo.
- [x] 2.2 Create `turnos.go`; add the bounded acceptance metadata seam to `internal/usecase/procesar_mensaje.go` and its tests.

## Phase 3: S3 Advisor Read Model (READY, <=360)
- [x] 3.1 Add `consultar_leads_test.go` and `cola_test.go` for priority order, cupo counting, filters, detail, and absent-ficha distinction.
- [x] 3.2 Create `consultar_leads.go` and `cola.go`; preserve stored priority and derive queue fields without mutation.

## Phase 4: S4 Buyer Persona (READY, <=330)
- [x] 4.1 Add `buyer_persona_test.go` and `gerencia_test.go` for deterministic project/catalog aggregates and immutable source reads.
- [x] 4.2 Create `buyer_persona.go` and `gerencia.go` over the catalog snapshot; handlers only present DTOs.

## Phase 5: S5 Simulated Clock (READY, <=370)
- [x] 5.1 Add `avanzar_demo_test.go`, `reloj/postgres_test.go`, `demo_repository_test.go`, `demo_tiempo_test.go`, and coordinator tests for persistence, one tick/count, and `/salud`.
- [x] 5.2 Create `avanzar_demo.go`, `reloj/postgres.go`, and `postgres/demo_repository.go`; extend `puertos.go`, `demo.go`, `salud.go`, and `coordinadora.go` without duplicate milestones.

## Phase 6: S6 Confirmed Reset (DONE, <=300)
- [x] 6.1 Add reset idempotency, disabled-gate, buyer-preservation, no-DDL, and <3s tests in `reiniciar_demo_test.go` and `demo_reset_test.go`.
- [x] 6.2 Create `reiniciar_demo.go`; extend `postgres/demo_repository.go` and tests with one transaction deleting only `fichas→hitos→planes→mensajes→leads`, then approved seed/date.
