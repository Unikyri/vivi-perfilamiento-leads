# Proposal: In-Memory Event Bus and Deterministic Coordinator (Issue #22)

## Intent

Issues #18–#21 delivered callable use cases that publish `LeadNuevo`, `PerfilCompleto`, and
`RutaDecidida` into `usecase.BusEventos`, but no production bus and no subscriber exist: every
event is published into nothing. Issue #22 adds the missing Observer bus (NFR-M-03.7) and the
Mediator coordinator (NFR-M-03.6, RF-M3-01..03) whose routing is a deterministic table, never LLM
judgement, so `ASESOR` leads reach `GenerarFicha` without wiring agents to each other.

## Scope

### In Scope
- `internal/infrastructure/bus/memoria.go`: `EnMemoria` implementing `usecase.BusEventos` —
  synchronous, registration-ordered, deep-copying, panic-isolated, privacy-safe logging.
- `internal/adapters/agentes/agentes.go`: narrow handler interfaces satisfied directly by the
  existing use-case pointers (ADR-1: agents are orchestration mechanism, tools delegate to use cases).
- `internal/adapters/agentes/coordinadora.go`: `Registrar()` + deterministic dispatch handlers.
- Focused tests: `bus/memoria_test.go`, `agentes/coordinadora_test.go`.

### Out of Scope (frozen non-goals)
- ADK graph/LLM/prompting (spike unvalidated; plain Go must work alone), HTTP endpoints and
  composition root (`cmd/servidor/main.go`), clock adapter, WhatsApp gateway wiring (#23/#24/#30).
- Any edit to `puertos.go`, `domain`, repositories, migrations, Contract v1.1, or #18–#21 use cases.
- Greeting/first message, plan creation, `PausarContacto` routing, outbox/queue/WebSocket (ADR-9).

## Frozen Decisions

| # | Decision |
|---|----------|
| F1 | Plain Go only. No LLM, ADK, HTTP, Postgres, or composition-root wiring; the demo must not depend on the framework. |
| F2 | Canonical event ownership: `CalificarLead` (#20) is the sole `RutaDecidida` producer. The coordinator **consumes** it and MUST NOT republish it. This overrides the issue's `alPerfilCompleto` snippet, which would emit a duplicate, lower-fidelity event. |
| F3 | Bus semantics: `Publicar` snapshots handlers for `e.Tipo` under lock, releases the lock, then calls them **synchronously in registration order** — no goroutine, queue, or retry. Nested publish/subscribe inside a handler MUST NOT deadlock. Unknown event type and zero subscribers are silent successes. |
| F4 | Isolation: deep-copy the event once at publish and again per subscriber, so a producer's later mutation and one handler's writes cannot be observed by the next. A shallow map copy is insufficient because `RutaDecidida.Payload` carries `[]domain.Recomendacion`. Nil handlers are ignored at `Suscribir`. |
| F5 | Panic isolation: each handler runs under its own `recover()`; a panic is logged as a fixed safe category and later subscribers still run. The port returns no error, so handler failures never propagate to the publisher. |
| F6 | Privacy (NFR-S-02): logs carry only `tipo`, opaque `lead_id`, handler identity, and outcome. Never log `Payload`, message text, names, cédulas, or raw panic/error values (repository errors can embed identifiers). |
| F7 | Registration table (exact, deterministic): `LeadNuevo` → observe only (optional injected handler); never calls `PerfilarLead`, which already published the event post-creation, and never fabricates a greeting because no greeting use case exists. `PerfilCompleto` → `CalificarLead.Ejecutar(ctx, EntradaCalificar{LeadID})`. `RutaDecidida` → `GenerarFicha` only when route is `ASESOR`. `TickReloj` → `EjecutarHitos` after parsing `hasta`. `MensajeEntrante`, `HitoVencido`, `ResultadoReal`, `PausarContacto` are **not** subscribed: they have no producer today, and `ProcesarMensaje` already handles `PAUSAR_CONTACTO` in-turn. |
| F8 | Reprofile guard: `PerfilCompleto` is emitted both by `ProcesarMensaje` (lead `CALIFICADO`) and by #21's requalification handoff (lead `PERFILANDO`). `ErrLeadNoCalificable` is a deterministic **skip**, not a failure — no ficha, no event, no log of payload. `CalificarLead` validates state before any write, so the skip is side-effect free. |
| F9 | Route reading: `RutaDecidida.Payload["ruta"]` is stored as `domain.Ruta`; the handler accepts `domain.Ruta` or `string`. `NUTRICION`, missing, or unrecognized route → no ficha and **no automatic plan** (#21 requires explicit consent). Handoffs pass `lead_id` only, never conversation history (ADR-4). |
| F10 | Nil-safe construction: a nil injected dependency skips its subscription instead of panicking, so partial wiring stays usable. |

## Capabilities

### New Capabilities
- `bus-eventos`: in-memory synchronous pub/sub delivery, ordering, isolation, panic containment, safe logging.
- `coordinadora-agentes`: deterministic event→handler routing table, qualification skip, ASESOR ficha, tick dispatch.

### Modified Capabilities
- None. No existing spec requirement changes; #18–#21 behavior is consumed as-is.

## Approach

Two additive files plus tests per slice. The coordinator depends on narrow interfaces declared in
`agentes.go` (`Calificador`, `Documentadora`, `Nutricionista`) that `*usecase.CalificarLead`,
`*usecase.GenerarFicha`, and `*usecase.EjecutarHitos` satisfy structurally. This keeps the Mediator
testable with in-package stubs — no Postgres, LLM, or HTTP — while `Registrar()` remains the single
literal routing table required by RF-M3-01.

## Delivery Slices

| Slice | Content | Forecast | Rollback |
|-------|---------|---------:|----------|
| 1 | `bus/memoria.go` + `memoria_test.go`: order, sync execution, nested publish, deep copy at both boundaries, nil handler, panic continuation, no-subscriber, context propagation, `-race` registration | 150 runtime + 190 tests ≈ 340 | Revert both files; no existing behavior changes |
| 2 | `agentes/agentes.go` + `coordinadora.go` + `coordinadora_test.go`: full table, 10× determinism, ASESOR-ficha vs NUTRICION-no-ficha, `ErrLeadNoCalificable` skip, no duplicate `RutaDecidida`, tick parse, nil-safe, log metadata | 140 runtime + 200 tests ≈ 340 | Revert the three files; slice 1 stands alone |

Both slices stay under the 400 authored-line budget and are independently reviewable.
Runtime harness is `N/A`: no HTTP, ADK, external process, or composition-root boundary is in scope.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `internal/infrastructure/bus/` | New | Production `BusEventos` implementation |
| `internal/adapters/agentes/` | New | Handler interfaces + Coordinadora (replaces `doc.go` marker only by addition) |
| `internal/usecase/`, `internal/domain/`, `cmd/` | Untouched | Consumed via existing ports |

Note: paths follow Issue #22 and ADR-1 (bus in `infrastructure`, agents in `adapters`), superseding
the exploration's tentative `internal/usecase/` placement.

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Duplicate `RutaDecidida` from the issue's snippet | High | F2 forbids republication; test asserts exactly one route event per qualification |
| Unguarded qualification breaks #21 reprofile | High | F8 treats `ErrLeadNoCalificable` as a skip; test covers a `PERFILANDO` lead |
| Untyped `Payload` mis-read | Med | F9 fixes accepted types; unknown values are no-ops, never guesses |
| Shared mutable payload across handlers | Med | F4 deep copy at publish and per subscriber |
| PII leaking into logs | Med | F6 metadata-only logging asserted by test |

## Rollback Plan

Every file is new and additive. Reverting slice 2 leaves the bus available and unused; reverting
slice 1 restores the pre-#22 state where publications are no-ops. No migration, schema, port,
Contract, or existing use case is modified, so `main` stays deployable after either revert.

## Dependencies

- #18, #19, #20, #21 (all merged). Blocks #24.

## Success Criteria

- [ ] Bus delivers to subscribers in registration order; no subscribers and unknown types do not break.
- [ ] A panicking handler neither kills the process nor prevents later subscribers.
- [ ] Handlers cannot observe another handler's payload mutation.
- [ ] Routing is a deterministic table: 10 identical publications invoke the same handler, no LLM.
- [ ] `ASESOR` route generates a ficha; `NUTRICION` generates neither ficha nor plan.
- [ ] Exactly one `RutaDecidida` per qualification; coordinator never publishes it.
- [ ] A `PERFILANDO` reprofile `PerfilCompleto` is skipped without error or side effects.
- [ ] Logs contain no message content, names, or cédulas.
- [ ] `go test ./internal/infrastructure/bus/... ./internal/adapters/agentes/... && go build ./...` passes per slice.

## Proposal question round

Non-interactive execution; these product assumptions need confirmation before spec:
1. Is consuming the canonical `RutaDecidida` (instead of the issue's republish) accepted as the ASESOR→ficha trigger? (assumed yes — it preserves route, priority, semaphore, and recommendations)
2. Should `LeadNuevo` stay observe-only until #24 owns the greeting? (assumed yes — no greeting use case exists)
3. Is leaving `PausarContacto` unrouted acceptable, given `ProcesarMensaje` already handles the pause in-turn? (assumed yes)
