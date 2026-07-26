# Proposal: Complete Contract v1.1 §3 HTTP API + demo endpoints (Issue #23)

## Intent

Block B (frontend, dashboard) is fully blocked: `web/src/models/api.ts` already calls nine `/api` routes that do not exist. The server exposes only `GET /salud`, while Issues #18–#22 delivered callable use cases (`PerfilarLead`, `ProcesarMensaje`, `CalificarLead`, `GenerarFicha`, `GestionarPlan`, `EjecutarHitos`), the in-memory bus and the deterministic coordinator. This change adds the exact Contract §3 surface as a thin adapter over those use cases so A/B integration happens only over HTTP + data files.

## Scope

### In Scope
- Single error envelope `{error:{codigo,mensaje,detalles}}` and the Contract §3 code→status catalog, including `FICHA_NO_DISPONIBLE` 404 (§3.6).
- `POST /api/leads` (201, `precargado_id` seed override), `POST /api/leads/{lead_id}/mensajes` (**202** + background turn, ADR-3), `GET /api/leads/{lead_id}/conversacion` (with `turno_en_proceso`).
- `GET /api/leads` (queue: `cupo_10`, `prioridad` desc, `semaforo`, `resumen`, `?afiliado`/`?ruta` filters), `GET /api/leads/{lead_id}`, `GET /api/leads/{lead_id}/ficha` (read-only).
- `GET /api/gerencia/buyer-persona` (per project and catalog summary), `POST /api/demo/tiempo`, `POST /api/demo/reiniciar`.
- Persisted simulated clock (`demo` table, Contract §5) replacing `time.Now()` in `/salud.fecha_simulada`.
- Composition wiring inside the `// === BLOQUE A ===` section of `cmd/servidor/main.go` only.
- `httptest`-level contract tests per endpoint with deterministic fakes.

### Out of Scope
- **ADK in any form** (explicit user constraint; consistent with #19/#22 exclusions). Handlers call plain-Go use cases only.
- WebSocket/SSE, pagination, authentication/CORS policy (Contract defines none for Fase 0), IP rate limiting (NFR-S-04 *Should*; `LIMITE_TASA` stays reserved).
- Static frontend serving and dashboard routes (Bloque B, `=== BLOQUE B ===`), frontend changes, new migrations, load testing (NFR-E-02).
- Changing motor/use-case behavior, spec-level requirements of #18–#22, or the Contract itself.

## Capabilities

### New Capabilities
- `api-http`: Contract §3 request/response DTOs, single error envelope, status mapping, validation limits, 202 + turn-state semantics.
- `cola-leads`: advisor queue read model — ordering by persisted `prioridad`, `cupo_10` counting, derived `semaforo`/`resumen`, filters.
- `buyer-persona`: deterministic aggregation of `data/compradores.json` per project (proportions, samples, desistimiento).
- `demo-control`: simulated clock advance with `TickReloj` + milestone dispatch, and idempotent seed reset.

### Modified Capabilities
- None. No existing spec-level requirement changes; existing ports and use cases are consumed as-is.

## Approach

Thin adapter, business logic stays inside the application boundary:
1. `errores.go` presenters + code catalog; `rutas.go` `Controlador` and `Registrar(mux)` seam so tests build the router without `main`.
2. Handlers decode/validate wire limits (texto ≤ 2000 → `VALIDACION`; audio ≤ 60 s, ≤ 2 MB decoded, MIME ∈ {webm, ogg, mpeg} → `AUDIO_INVALIDO`), then delegate.
3. `POST /mensajes` returns 202 immediately, marks the lead's turn in progress in a mutex-guarded tracker, and runs `ProcesarMensaje` on a detached context; the tracker is cleared in `defer`. Frontend polls conversation (1.5 s).
4. Queue, lead detail and buyer-persona aggregation live in read-only use cases; handlers never recompute `prioridad` (Contract §3.5 formula already implemented in `CalificarLead`) and never generate a ficha on read.
5. `POST /demo/tiempo` accepts exactly one of `avanzar_hasta`/`avanzar_dias`, advances the persisted clock, publishes `TickReloj`, and returns `EjecutarHitos` count.
6. `POST /demo/reiniciar` deletes only app-owned rows (`fichas`, `hitos`, `planes`, `mensajes`, `leads`) and restores seed + `fecha_simulada`; idempotent, < 3 s (NFR-D-02), never DDL.

### Contract conflicts resolved in favor of the Wiki
| Conflict | Resolution |
|---|---|
| Exploration assumed `turno_en_proceso=false` (synchronous) | Contract §3.2 + ADR-3 win: 202 + background turn, `true` while processing |
| Issue sketch returns `err.Error()` as `ERROR_INTERNO` message | NFR-S-02 wins: generic message, detail only in structured logs |
| Issue queue text omits `porcentaje_ventana` | Contract §3.5 shape wins: `cupo_10 {usados, porcentaje_ventana}` |
| `semaforo`/`resumen` not persisted | Derived deterministically from persisted `ruta`, `afiliado`, `categoria`, `capacidad`, `intencion` in the queue read model |
| §3 catalog line omits `FICHA_NO_DISPONIBLE` | §3.6 wins: `FICHA_NO_DISPONIBLE` 404 when the lead exists without ficha; `LEAD_NO_ENCONTRADO` otherwise |

## Invariants

**Privacy (NFR-S-02).** No `cedula`, phone, message text or audio bytes in logs or in `detalles`; logs carry `lead_id`, event, latency only. Audio stays in memory, is never persisted or echoed. No provider names, SQL text or stack traces in responses.

**Errors.** Exactly one envelope for every endpoint and every failure path; status derives only from the code catalog; unknown code → 500 `ERROR_INTERNO`; `usecase.ErrNoEncontrado`/`NotFoundError` → 404 with the correct code; invalid state transition → 409 `TRANSICION_INVALIDA`; provider outage surfaced only as 503 `PROVEEDOR_LLM_NO_DISPONIBLE`. All responses `application/json; charset=utf-8`, money as integer COP, timestamps RFC 3339 UTC, enums UPPER_SNAKE.

## Affected Areas

| Area | Impact | Description |
|---|---|---|
| `internal/adapters/http/errores.go`, `leads.go`, `demo.go`, `rutas.go` | New | Presenters, handlers, router seam |
| `internal/adapters/http/*_test.go` | New | `httptest` contract tests |
| `internal/usecase/` (queue, detail, buyer-persona, reset ports) | New | Read models + reset port; no change to existing use cases |
| `internal/infrastructure/postgres/` | New | Reset/seed + demo-clock persistence implementations |
| `internal/infrastructure/reloj/` | New | Persisted simulated clock (`demo` table) |
| `cmd/servidor/main.go` | Modified | Wiring inside `=== BLOQUE A ===` only |
| `internal/adapters/http/salud.go` | Modified | `fecha_simulada` from the clock |

## Dependencies

- Issue #4 composition root and Issues #18–#22 (done, archived).
- Contract v1.1 §0–§3, §5, §6 (authority); doc 11 ADR-3/§4.1/§4.3; doc 09 NFR-R-02, NFR-D-02, NFR-S-02; doc 07 US-14.
- Postgres schema `001_esquema_inicial.sql` incl. `demo(clave,valor)`; `data/*.json` seeds. Blocks #24–#26.

## Delivery (chained PRs, ≤ 400 authored runtime/test lines per slice)

| Slice | Content | Est. lines |
|---|---|---|
| S1 | `errores.go`, router/controller seam, `POST /api/leads` + seeds, `GET /conversacion`, wiring + tests | ~300 |
| S2 | `POST /mensajes` 202, turn tracker, text/audio limits + tests | ~320 |
| S3 | `GET /api/leads` queue read model, `GET /api/leads/{id}`, `GET /ficha` + tests | ~360 |
| S4 | buyer-persona aggregation use case + endpoint + tests | ~330 |
| S5 | persisted simulated clock, `POST /demo/tiempo`, `/salud` clock + tests | ~370 |
| S6 | `POST /demo/reiniciar` reset port + Postgres impl + idempotency tests | ~300 |

Feature Branch Chain on `feat/issue-23-http-api`: S1 targets the tracker branch, each later slice targets the previous one. Every slice is independently runnable, tested and revertible.

## Risks

| Risk | Likelihood | Mitigation |
|---|---|---|
| Background turn goroutine outlives request context / leaks | Med | Detached `context.Background()` with timeout, `defer` clearing the tracker, logged failures, no fire-and-forget writes outside `ProcesarMensaje` |
| Turn tracker is in-process; multi-dyno deploy would desync `turno_en_proceso` | Med | Fase 0 runs a single dyno; documented limitation, tracker interface allows a `demo`-table implementation later |
| `POST /demo/reiniciar` deletes lead data | High impact | Restricted to app-owned tables, no DDL, config-gated (`DEMO_SEED`), separate final slice, explicit maintainer confirmation before apply |
| API is unauthenticated (Contract defines no auth) | Med | Contract-compliant for the demo; explicitly flagged as out of scope and recorded as an accepted Fase 0 risk |
| `resumen`/`semaforo` derivation drifts from dashboard expectations | Med | Derive only from persisted Contract fields; freeze wording in specs and golden-style tests |
| Wiring changes break startup (`DATABASE_URL` required) | Low | Additive registration inside `=== BLOQUE A ===`, router built by a constructor testable without Postgres |
| p95 < 300 ms for LLM-free endpoints (NFR-R-02) | Low | Queue/ficha use single indexed reads; buyer-persona aggregation cached per catalog load |

## Rollback Plan

Routes are additive: reverting a slice's commit/PR removes only its endpoints and leaves `/salud` and prior slices working. No migrations are added, so no schema rollback is needed. `POST /demo/reiniciar` is disabled by config without redeploy; if reset misbehaves, disable the flag first, then revert S6. If the async turn model proves unstable, revert S2 only — the frontend polling path degrades to "no reply" rather than corrupt state.

## Success Criteria

- [ ] All nine endpoints return the exact Contract §3 shapes and status codes; `POST /mensajes` returns 202 and processes in background.
- [ ] Every error path uses the single envelope with the catalog status; no PII or internal detail in bodies or logs.
- [ ] Queue ordered by `prioridad` desc with `cupo_10.usados` = non-affiliate ASESOR leads; ficha read returns 404 `FICHA_NO_DISPONIBLE` when absent.
- [ ] `POST /demo/tiempo` requires exactly one field and reports `hitos_disparados`; `POST /demo/reiniciar` is idempotent and < 3 s.
- [ ] `main.go` modified only inside `=== BLOQUE A ===`; no ADK dependency added (`go.mod` unchanged for ADK).
- [ ] ≥ 12 `httptest` tests from the issue DoD pass; `go build ./...` and `go test ./...` green per slice; each slice ≤ 400 authored runtime/test lines.

## Proposal question round

No interactive channel was available in this phase. Assumptions needing user review:

1. **Preloaded seeds** — `precargado_id ∈ {ana, carlos, luisa}` is assumed to be a server-side seed catalog (Carlos intentionally has no affiliate match; his wife's cédula `1015789456` exists), not a lookup in `afiliados_mock.json`. Confirm the three seeds' exact identity data.
2. **`resumen` wording** — assumed template `"Afiliada cat. A · presupuesto $166.8M · intención alta"` derived from persisted fields. Confirm whether the dashboard requires exact phrasing.
3. **Buyer-persona without `proyecto_id`** — assumed `{proyectos:[{summary per project}]}` with the same proportion semantics. Confirm the per-project summary fields.
4. **Reset blast radius** — assumed deletion of all `leads`/`mensajes`/`planes`/`hitos`/`fichas` plus clock reset, keeping `compradores`. Confirm nothing else must survive a reset.
5. **Multi-instance** — assumed single dyno for Fase 0 so the in-process turn tracker is acceptable. Confirm, otherwise S2 must persist turn state in `demo`.
