# Apply Progress: issue-23-http-api

## Slice S1 — HTTP foundation

**Mode:** Standard (strict_tdd=false)
**Delivery:** Feature-branch chain, auto-chain; S1 targets `feat/issue-23-http-api`.
**Scope:** HTTP foundation, error/JSON helpers, testable router seam, `POST /api/leads`, `GET /api/leads/{lead_id}/conversacion`, deterministic `LeadNuevo` greeting, and Bloque A wiring only.

### Completed tasks
- [x] 1.1 Added HTTP and greeting contract tests for seed creation, conversation, strict JSON, and privacy-safe failures.
- [x] 1.2 Added `errores.go`, `rutas.go`, `leads.go`, and `saludar_lead.go`; wired only the S1 dependencies and routes in the `BLOQUE A` section of `cmd/servidor/main.go`.

### Implementation
- `internal/adapters/http/errores.go`: JSON content type, complete Contract error code/status catalog, generic internal errors, and no sensitive error details.
- `internal/adapters/http/rutas.go`: dependency-validated `Controlador`, `Registrar(*http.ServeMux)` seam, API fallback, and temporary wall-clock adapter for S1 composition.
- `internal/adapters/http/leads.go`: strict single-value JSON decoding, seed override for `ana|carlos|luisa`, create response, and conversation DTO.
- `internal/usecase/saludar_lead.go`: deterministic plain-Go `LeadNuevo` subscriber that persists the initial Vivi message without ADK/provider calls.
- `cmd/servidor/main.go`: repository, catalog, IDs, bus, clock, profiler, greeting subscriber, coordinator, and router registration inside `BLOQUE A`.

### Work Unit Evidence
| Evidence | Result |
|---|---|
| Focused test | `go test ./internal/adapters/http ./internal/usecase ./cmd/servidor` — PASS |
| Runtime harness | `httptest` create + conversation contract tests — PASS; seed override, strict JSON, privacy-safe 404, and deterministic greeting covered |
| Full test | `go test ./...` — PASS |
| Race test | `go test -race ./...` — PASS |
| Build | `go build ./...` — PASS |
| Static checks | `go vet ./...`, `gofmt` check, `git diff --check` — PASS |
| Authored runtime/test budget | 379 new S1 file lines + 21 `main.go` additions = 400; cap respected |
| Rollback boundary | Revert `internal/adapters/http/{errores,rutas,leads,leads_test}.go`, `internal/usecase/{saludar_lead,saludar_lead_test}.go`, and the S1 `BLOQUE A` wiring in `cmd/servidor/main.go`. |

### Deviations and risks
- S1 uses a wall-clock adapter only to satisfy the existing `Reloj` port; S5 must replace it with the persisted simulated clock and update `/salud` as designed.
- The in-process greeting is synchronous through the existing in-memory bus; no ADK or new dependency was introduced.
- S2 asynchronous messages, queue/detail/ficha, buyer persona, demo time, and reset are intentionally not implemented.

### Remaining work
- [ ] 3.1–3.2 S3 advisor read model
- [ ] 4.1–4.2 S4 buyer persona
- [ ] 5.1–5.2 S5 simulated clock
- [ ] 6.1–6.2 S6 confirmed reset (blocked pending maintainer confirmation)

**Next:** `sdd-apply` for S2, then `sdd-verify` after all assigned implementation tasks are complete.

## Slice S2 — Asynchronous turns

**Mode:** Standard (strict_tdd=false)
**Delivery:** Feature-branch chain, auto-chain; S2 targets S1 on `feat/issue-23-http-api`.
**Scope:** `POST /api/leads/{lead_id}/mensajes` strict validation, 202 acceptance metadata, bounded per-lead background execution, cancellation, and polling-compatible `turno_en_proceso`.

### Completed tasks
- [x] 2.1 Added `turnos_test.go` contract coverage for 202 acceptance, active-to-clear polling, duplicate active turn rejection, shutdown cancellation, text/audio limits, decoded-size, MIME/duration errors, and raw-audio non-echo.
- [x] 2.2 Added `EjecutorTurnos`, strict message transport decoding, and additive `EntradaMensaje.MensajeID`/`RecibidoEn` metadata consumed by `ProcesarMensaje`.

### Implementation
- `internal/adapters/http/turnos.go`: process-root context, mutex-protected one-active-turn-per-lead map, `WaitGroup` shutdown, strict Contract §3.2 request mapping, 202 response, and no raw-audio response path.
- `internal/adapters/http/{rutas,leads,errores}.go`: registered the message route, expose tracker state in conversation polling, allow the larger base64 request body, and map duplicate active turns to `LIMITE_TASA` without sensitive details.
- `internal/usecase/procesar_mensaje.go`: exported the existing validation once for the adapter and uses accepted message ID/time when present; legacy direct callers retain generated ID/current clock behavior.
- `cmd/servidor/main.go`: wires `ProcesarMensaje` and tracker in Bloque A and closes the tracker during process shutdown.

### Work Unit Evidence
| Evidence | Result |
|---|---|
| Focused test | `go test ./internal/adapters/http ./internal/usecase` — PASS |
| Runtime harness | `httptest` POST 202 → GET conversation active → release → GET clear; duplicate active request 429; shutdown cancellation — PASS |
| Full test | `go test ./...` — PASS |
| Race test | `go test -race ./...` — PASS |
| Build | `go build ./...` — PASS |
| Static checks | `go vet ./...`, changed-file `gofmt`, `git diff --check` — PASS; two pre-existing S1 pipeline tests remain reported by repository-wide `gofmt -l .` |
| Authored runtime/test budget | 373 estimated S2 additions/changes; within the 400-line slice cap |
| Rollback boundary | Revert `internal/adapters/http/turnos.go`, `turnos_test.go`, S2 route/polling/error edits, the `EntradaMensaje` metadata changes/tests, and S2 Bloque A tracker wiring in `cmd/servidor/main.go`; retain S1 foundation. |

### Deviations and risks
- The in-process tracker is intentionally single-instance and bounded to one active turn per lead; multi-dyno deployment still requires a shared queue/tracker as documented.
- Audio bytes are decoded only for validation/provider dispatch and are not included in the 202 response or any HTTP error. Existing Contract-compatible `audio_original` metadata remains the only persisted audio marker.
- No ADK, queue, S3+ read model, demo clock, or reset functionality was added.

### Remaining work
- [ ] 3.1–3.2 S3 advisor read model
- [ ] 4.1–4.2 S4 buyer persona
- [ ] 5.1–5.2 S5 simulated clock
- [ ] 6.1–6.2 S6 confirmed reset (blocked pending maintainer confirmation)

**Next:** `sdd-apply` for S3, then `sdd-verify` after all assigned implementation tasks are complete.

## Slice S3 — Advisor read model

**Mode:** Standard (strict_tdd=false)
**Delivery:** Feature-branch chain, auto-chain; S3 targets S2 on `feat/issue-23-http-api`.
**Scope:** Deterministic advisor queue, lead detail, read-only ficha, filter validation, and distinct missing-lead/missing-ficha errors. No priority/state mutation, ADK, buyer persona, simulated time, or reset work.

### Completed tasks
- [x] 3.1 Added `consultar_leads_test.go` and `cola_test.go` for persisted-priority ordering, cupo counting, filters, detail privacy, and ficha absence distinction.
- [x] 3.2 Added `consultar_leads.go` and `cola.go`; the read model preserves persisted priority, derives `semaforo`/`resumen`, and performs no writes.

### Implementation
- `internal/usecase/consultar_leads.go`: `ConsultarLeads` queue/detail/ficha boundary over `LeadRepository` and `FichaRepository`; deterministic local priority ordering, `cupo_10` count for non-affiliated ASESOR leads, derived semáforo/resumen, safe detail DTO, and explicit ficha not-found normalization.
- `internal/usecase/consultar_leads_test.go`: queue ranking/cupo/summary/read-only and filter/ficha distinction coverage.
- `internal/adapters/http/cola.go`: strict `afiliado=true|false` and Contract route filters, queue/detail/ficha handlers.
- `internal/adapters/http/cola_test.go`: `httptest` queue, invalid filters, detail, ficha success, missing ficha, and missing lead contract cases.
- `internal/adapters/http/rutas.go`: additive GET queue/detail/ficha registrations and ficha repository dependency.
- `internal/adapters/http/errores.go`: maps `NotFoundError{Resource: "ficha"}` to `FICHA_NO_DISPONIBLE`.
- `cmd/servidor/main.go`: wires the existing Postgres `FichaRepository` into the controller in Bloque A.

### Work Unit Evidence
| Evidence | Result |
|---|---|
| Focused test | `go test ./internal/adapters/http ./internal/usecase` — PASS |
| Runtime harness | `httptest` GET queue/detail/ficha — PASS; priority order, cupo, filters, safe detail, and distinct 404 codes covered |
| Full test | `go test ./...` — PASS |
| Race test | `go test -race ./...` — PASS |
| Build | `go build ./...` — PASS |
| Static checks | `go vet ./...` and `git diff --check` — PASS; repository-wide `gofmt -l .` reports only pre-existing `internal/pipeline/{compradores,proyectos}_test.go` |
| Authored runtime/test budget | 398 nonblank S3 runtime/test lines; within the 400-line slice cap |
| Rollback boundary | Revert `internal/usecase/{consultar_leads,consultar_leads_test}.go`, `internal/adapters/http/{cola,cola_test}.go`, S3 route/error edits, and the S3 ficha-repository wiring in `cmd/servidor/main.go`; retain S1/S2. |

### Deviations and risks
- Queue ordering is enforced again on the local read-model projection for deterministic behavior even if a repository implementation returns an unsorted slice; persisted `prioridad` is never recalculated or saved.
- `FichaRepository` remains optional in test controllers and is required only when `/ficha` is requested; production Bloque A now supplies the existing Postgres implementation.
- No ADK, buyer persona, time/reset, priority mutation, or lead state transition was added.

### Remaining work
- [ ] 4.1–4.2 S4 buyer persona
- [ ] 5.1–5.2 S5 simulated clock
- [ ] 6.1–6.2 S6 confirmed reset (blocked pending maintainer confirmation)

**Next:** `sdd-apply` for S4.


## Slice S4 — Buyer persona

**Mode:** Standard (strict_tdd=false)
**Delivery:** Feature-branch chain, auto-chain; S4 targets S3 on `feat/issue-23-http-api`.
**Scope:** Deterministic buyer-persona project/catalog aggregates from the immutable catalog snapshot and `GET /api/gerencia/buyer-persona`. No ADK, lead mutation, simulated time, reset, migrations, or frontend work.

### Completed tasks
- [x] 4.1 Added `buyer_persona_test.go` and `gerencia_test.go` for project/catalog proportions, deterministic dates/order, filter validation, and immutable source reads.
- [x] 4.2 Added `buyer_persona.go` and `gerencia.go`; wired the catalog dependency and route while keeping aggregation out of HTTP handlers.

### Implementation
- `internal/usecase/buyer_persona.go`: aggregates affiliation, normalized category/age-band proportions, abandonment rate, sample count, and newest valid `fecha_opcion` as UTC RFC 3339; includes sorted catalog summaries and zero-sample projects without mutating catalog results.
- `internal/adapters/http/gerencia.go`: validates optional `proyecto_id`, delegates project/catalog reads, and presents Contract DTOs through the existing JSON/error helpers.
- `internal/adapters/http/rutas.go`: adds the optional catalog dependency and registers the Gerencia route.
- `cmd/servidor/main.go`: supplies the existing eager `CatalogoRepository` to the HTTP controller.
- `internal/usecase/buyer_persona_test.go` and `internal/adapters/http/gerencia_test.go`: focused aggregation, immutability, deterministic repeat, catalog ordering, HTTP project/catalog response, and invalid duplicate/empty filters.

### Work Unit Evidence
| Evidence | Result |
|---|---|
| Focused test | `go test ./internal/usecase ./internal/adapters/http` — PASS |
| Runtime harness | `httptest` GET project and catalog buyer-persona routes — PASS; proportions, sorted catalog, and invalid filters covered |
| Full test | `go test ./...` — PASS |
| Race test | `go test -race ./...` — PASS |
| Build | `go build ./...` — PASS |
| Static checks | `go vet ./...`, changed-file `gofmt`, `git diff --check` — PASS |
| Authored runtime/test budget | 330 S4 lines after trimming; within the 330-line slice cap and 400-line review budget |
| Rollback boundary | Revert `internal/usecase/buyer_persona{,_test}.go`, `internal/adapters/http/gerencia{,_test}.go`, S4 route/dependency edits in `rutas.go`, and the catalog dependency wiring in `cmd/servidor/main.go`; retain S1–S3. |

### Deviations and risks
- Catalog rows with no valid source date return the deterministic zero UTC RFC 3339 timestamp `0001-01-01T00:00:00Z`; no wall-clock value is introduced.
- All Contract-normalized category and age-band keys are emitted with zero proportions when absent, preserving a stable response shape.
- No ADK, lead writes, time/reset functionality, migrations, or frontend files were added.

### Remaining work
- [ ] 5.1–5.2 S5 simulated clock
- [ ] 6.1–6.2 S6 confirmed reset (blocked pending maintainer confirmation)

**Next:** `sdd-apply` for S5, then `sdd-verify` after all assigned implementation tasks are complete.


## Slice S5 — Simulated clock

**Mode:** Standard (strict_tdd=false)
**Delivery:** Feature-branch chain, auto-chain; S5 targets S4 on `feat/issue-23-http-api`.
**Scope:** Persisted simulated clock, `POST /api/demo/tiempo`, health date, one coordinator TickReloj dispatch, and focused tests. No ADK, reset, migrations, or frontend work.

### Completed tasks
- [x] 5.1 Added usecase, persisted-clock, Postgres repository, HTTP, coordinator, and health tests for persistence, exact-one input validation, tick count, and simulated-date reporting.
- [x] 5.2 Added the demo application boundary, mutex-safe persisted clock, existing `demo` table repository, demo route, health wiring, and context-preserving coordinator result bridge.

### Implementation
- `internal/usecase/avanzar_demo.go`: validates exactly one target/days operation, rejects backwards time, persists before updating the cache, emits one synchronous `TickReloj`, and returns the coordinator's milestone count.
- `internal/usecase/puertos.go`: adds the persisted demo repository port and private context result sink used without mutating cloned bus payloads.
- `internal/infrastructure/reloj/postgres.go`: loads the persisted date once, initializes a missing value, and provides mutex-safe simulated `Ahora`/`FechaSimulada`/`Avanzar` behavior.
- `internal/infrastructure/postgres/demo_repository.go`: reads/writes `demo(clave, valor)` with an upsert; no migration was added.
- `internal/adapters/http/{demo,rutas,errores}.go`: registers strict `POST /api/demo/tiempo`, supports RFC3339/date targets, and maps backward time to `VALIDACION`.
- `internal/adapters/agentes/coordinadora.go`: accepts canonical `fecha_simulada` and writes the sole milestone result into the context sink; `sync.Once` prevents duplicate registration.
- `cmd/servidor/main.go`: wires one persisted clock, one coordinator milestone executor, demo repository, and demo use case; `/salud` reads the same clock.

### Work Unit Evidence
| Evidence | Result |
|---|---|
| Focused test | `go test ./internal/usecase ./internal/adapters/http ./internal/adapters/agentes ./internal/infrastructure/reloj ./internal/infrastructure/postgres ./cmd/servidor` — PASS |
| Runtime harness | `httptest` POST `/api/demo/tiempo` with valid/ambiguous inputs and synchronous TickReloj result bridge — PASS; coordinator canonical-date/count test passes |
| Full test | `go test ./...` — PASS |
| Race test | `go test -race ./...` — PASS |
| Build | `go build ./...` — PASS |
| Static checks | `go vet ./...`, changed-file `gofmt`, `git diff --check` — PASS |
| Authored runtime/test budget | 370-line S5 estimate; within the <=400-line slice cap |
| Rollback boundary | Revert S5 demo clock/usecase/HTTP/coordinator/main wiring and S5 tests; retain S1–S4 and existing migration/table. |

### Deviations and risks
- The existing `demo` table is reused exactly; no migration or reset behavior was added.
- A missing persisted clock is initialized from UTC wall time once at startup, then all health and application reads use the cache.
- S6 reset remains blocked pending explicit maintainer confirmation and was not changed.

## Remaining work
- [ ] 6.1–6.2 S6 confirmed reset (blocked pending explicit maintainer confirmation)

**Next:** `sdd-apply` for S6 only after confirmation; otherwise verification must wait until the blocked slice is resolved.


## Slice S6 — Confirmed reset

**Mode:** Standard (strict_tdd=false)
**Delivery:** Feature-branch chain, auto-chain; S6 targets S5 on `feat/issue-23-http-api`.
**Scope:** Config-gated `POST /api/demo/reiniciar`, one-transaction demo reset, approved date restoration, buyer/schema preservation, and idempotency tests. No ADK, DDL, migrations, or frontend work.

### Completed tasks
- [x] 6.1 Added use-case and HTTP tests for disabled-gate read-only behavior, repeated reset identity, and the under-three-second bound; repository tests assert ordered deletes and absence of catalog/DDL operations.
- [x] 6.2 Added `ReiniciarDemo`, the reset route, `DEMO_SEED` wiring, and `DemoRepository.Reiniciar`, which deletes only `fichas`, `hitos`, `planes`, `mensajes`, and `leads` in one transaction before restoring `demo.fecha_simulada` to `2026-07-26T00:00:00Z`.

### Implementation

- `internal/usecase/reiniciar_demo.go`: checks the explicit config gate before mutation, delegates the reset transaction, and updates the in-memory simulated clock only after success.
- `internal/infrastructure/postgres/demo_repository.go`: begins one transaction, deletes child-to-parent application rows in order, upserts the approved demo date, commits, and never issues DDL or touches `compradores`.
- `internal/adapters/http/demo.go` and `rutas.go`: register `POST /api/demo/reiniciar` and return the Contract response/error envelope.
- `internal/infrastructure/config/config.go`: defaults `DEMO_SEED` to disabled; only the exact value `true` enables destructive reset.
- `cmd/servidor/main.go`: passes the configured gate and reset repository into the Bloque A controller.

### Work Unit Evidence
| Evidence | Result |
|---|---|
| Focused test | `go test ./internal/usecase ./internal/adapters/http ./internal/infrastructure/postgres ./internal/infrastructure/config ./cmd/servidor -count=1` — PASS |
| Runtime harness | `httptest` POST `/api/demo/reiniciar` twice and disabled-gate request — PASS; identical date, generic 500, zero disabled mutations |
| Full test | `go test ./... -count=1` — PASS |
| Race test | `go test -race ./... -count=1` — PASS |
| Build | `go build ./...` — PASS |
| Static checks | `go vet ./...`, `gofmt` on changed files, and `git diff --check` — PASS |
| Authored runtime/test budget | 251 new S6 runtime/test lines; within the <=400-line slice cap |
| Rollback boundary | Revert `reiniciar_demo.go` and its tests, reset route/dependency wiring, config default/test, and `DemoRepository.Reiniciar`; retain S1–S5 and the existing demo clock methods/table. |

### Deviations and risks

- The approved reset date is a fixed UTC value (`2026-07-26T00:00:00Z`) so repeated resets are byte-for-byte deterministic; no seed lead rows are inserted because the approved `ana|carlos|luisa` seed identities remain request-level constants and reset is limited to the specified application tables.
- `DEMO_SEED` now defaults to false for safe-by-default destructive behavior; deployments must explicitly set `DEMO_SEED=true`.
- No migration was added; the existing `demo` table is reused.

## Remaining work

- None for apply; all 12 implementation tasks are complete.

**Next:** `sdd-verify`.

## Final verifier remediation — CRITICAL-1 / CRITICAL-2

**Mode:** Standard (`strict_tdd=false`). **Delivery:** Feature-branch chain, auto-chain; bounded follow-up on S3/S5 read and demo wire boundaries. Scope was limited to Contract §3.4/§3.8 verifier findings, with the trivial filtered `cupo_10` and `fuente` enum corrections included.

### Implementation

- `internal/usecase/consultar_leads.go`: lead detail now derives `semaforo` and reads the optional nutrition `plan` through the existing `PlanRepository`; absent plans serialize as `null`. Queue `cupo_10.usados` is calculated from the unfiltered lead population even when the visible lead list has filters.
- `internal/adapters/http/rutas.go` and `cmd/servidor/main.go`: pass the existing Postgres `PlanRepository` into the detail read model.
- `internal/adapters/http/demo.go`: format `POST /api/demo/tiempo` and reset `fecha_simulada` as date-only `YYYY-MM-DD`; internal use-case timestamps remain unchanged.
- `internal/adapters/http/leads.go`: reject non-Contract `fuente` values and use `DEMO` when omitted.
- Direct discriminating tests cover detail semáforo/plan plus PII exclusion, global filtered cupo, exact date-only demo/reset JSON, and invalid fuente.

### Work Unit Evidence
| Evidence | Result |
|---|---|
| Focused test | `go test ./internal/usecase ./internal/adapters/http ./cmd/servidor -count=1` — PASS |
| Runtime harness | `httptest` detail/queue/demo/reset routes — PASS; plan/null shape, PII absence, filtered global cupo, exact date-only wire values, and source enum rejection covered |
| Full test | `go test ./... -count=1` — PASS |
| Race test | `go test -race ./... -count=1` — PASS |
| Build | `go build ./...` — PASS |
| Static checks | `go vet ./...`, `gofmt` on all remediation-touched Go files, and `git diff --check` — PASS; repository-wide `gofmt -l .` reports only pre-existing `internal/pipeline/compradores_test.go` and `internal/pipeline/proyectos_test.go` |
| Verify admission | `gentle-ai sdd-verify-validate` unavailable in this environment; no verify report was created or overwritten |
| Rollback boundary | Revert only the detail plan/semaforo projection and PlanRepository wiring, filtered cupo calculation, demo response formatting, fuente validation, their direct tests, and this SDD evidence. |

### Exact test counts

- Focused packages: 3 packages, all PASS.
- Full suite: 13 tested packages PASS; 3 packages report `[no test files]`; 1 PostgreSQL integration test is skipped because `VIVI_TEST_DATABASE_URL` is unset.
- Race suite: same package result as full suite, PASS.

## Remaining work

- None for apply remediation. Next: `sdd-verify`; verify-report remains pending because the admission validator is unavailable.
