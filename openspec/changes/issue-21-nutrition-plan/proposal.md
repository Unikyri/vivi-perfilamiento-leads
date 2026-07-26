# Proposal: Nutrition Plan with Explicit Consent and Time Advancement (Issue #21)

## Intent

A lead who cannot buy today must not be discarded. Issue #21 (UC-07, US-11/12/13) adds the missing
provider-free application layer that turns a `NUTRICION` route into a consented, calendar-anchored
plan, advances simulated time to fire due milestones, and honors "stop writing me" on the first ask.

## Scope

### In Scope
- `internal/domain/motor/plan.go`: pure `DisenarHitos` (deterministic, calendar-anchored, `AFILIACION` first when convertible).
- `internal/usecase/gestionar_plan.go`: `CrearPlan` (consent-gated) and `PausarPlan`.
- `internal/usecase/ejecutar_hitos.go`: `Ejecutar(ctx, hasta)` over ACTIVE plans only.
- Milestone microcopy rules: always offer pause in the last line; never presume financial distress.
- Focused Go tests per Issue #21 test list, including failure-ordering tests.

### Out of Scope (frozen non-goals)
- LLM/ADK/prompting, HTTP endpoints (#24), concrete `BusEventos`/`Reloj` implementations, event
  subscription and coordinator wiring (#22/#23/#30), frontend, migrations, Postgres repository changes.
- Contract v1.1 changes, new domain fields, new ports or repository methods.
- Editing `CalificarLead` or `ProcesarMensaje` (their `NUTRICION`/`PAUSAR_CONTACTO` handoffs stay as-is).

## Frozen Decisions

| # | Decision |
|---|----------|
| F1 | Provider-free: only callable use cases + pure planner. No integration, no wiring, no new ports. |
| F2 | Target/brecha: caller passes `PrecioObjetivo`; `brecha = max(0, PrecioObjetivo − lead.Capacidad.PresupuestoMax)`. Missing `Capacidad` or `PrecioObjetivo <= 0` → validation error, no plan, no message. Nothing new is persisted. |
| F3 | Consent+frequency: plan requires `consintio == true` AND `Frecuencia ∈ {QUINCENAL, MENSUAL, TRIMESTRAL}` (Contract §2.6). Invalid frequency → validation error, no plan, no message, no state change. `ConsentimientoEn = Reloj.Ahora()`. |
| F4 | No consent: exactly one door-open reminder appended, return `(nil, nil)`, no plan and no lead transition. A failed append returns the error (never swallowed); at most one reminder per call. |
| F5 | Create ordering: `Planes.Crear` → `lead.Transicionar(EN_NUTRICION)` → `Leads.Guardar`. `Crear` failure leaves the lead untouched. If the lead save fails, the plan stays persisted and a retry MUST reuse the existing plan via `Planes.PorLead` instead of creating a second one. |
| F6 | Tick ordering per due milestone: `Gateway.Enviar` → `Leads.AgregarMensaje` → `Planes.MarcarHito(NOTIFICADO)`. Mark only after send and append succeed; on failure skip that milestone, keep it `PENDIENTE`, continue the remaining ones, count only successes, return the aggregated error. No in-use-case retry loop; the next tick retries. |
| F7 | Pause mapping: plan → `PAUSADO`, then lead `EN_NUTRICION → PAUSADO`, then one farewell message. No plan → silent success. Already paused → idempotent success. Illegal lead transition → error without a farewell. Never asks why, never insists. A paused plan yields no messages on any later tick. |
| F8 | Requalification handoff is deterministic: when cumulative `NOTIFICADO` milestone `Monto` ≥ `plan.MetaMonto`, or the `REEVALUACION` milestone becomes `NOTIFICADO`, transition `EN_NUTRICION → PERFILANDO`, save, and publish exactly one `PerfilCompleto` per lead per tick. Never call `CalificarLead` directly and never publish `RutaDecidida` (#20 owns it). Nil bus is skipped silently. |
| F9 | Delivery: two chained PRs, each under 400 authored changed lines, each independently testable and revertible. |

## Capabilities

### New Capabilities
- `plan-nutricion`: consent-gated plan creation, deterministic milestone design, pause, tick execution, requalification handoff.

### Modified Capabilities
- None.

## Approach

Application orchestration plus a pure planner. `DisenarHitos` receives explicit `brecha`,
convertibility, `desde`, and the parsed economic calendar (never hardcoded dates) and returns
`[]domain.Hito`. `GestionarPlan` and `EjecutarHitos` compose existing `LeadRepository`,
`PlanRepository`, `MensajeriaGateway`, `Reloj`, `GeneradorID`, and `BusEventos` ports. Convertibility
reuses the existing package-private `esConversion(lead)`.

## Delivery Slices

| Slice | Content | Forecast |
|-------|---------|----------|
| 1 | `motor/plan.go` + `CrearPlan` (consent, frequency, brecha, F5 ordering) + tests | 330–385 lines |
| 2 | `PausarPlan` + `EjecutarHitos` (F6 ordering, microcopy, F8 handoff) + tests | 325–390 lines |

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `internal/domain/motor/plan.go` | New | Pure deterministic milestone planner |
| `internal/usecase/gestionar_plan.go` | New | Consent-gated create + pause |
| `internal/usecase/ejecutar_hitos.go` | New | Due-milestone execution |
| `internal/usecase/*_test.go` | New | Focused and failure-ordering tests, local fakes |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Inventing Contract fields (target, reminder flag, pause reason) | Med | F2/F3 keep all state inside existing Contract fields |
| Non-atomic cross-repository writes | High | F5/F6 freeze ordering, idempotent retry, mark-after-success tests |
| Duplicate route events vs #20 | Med | F8 forbids `RutaDecidida` and direct `CalificarLead` calls |
| Slice growth past 400 lines | Med | Chained slices, local test fakes only, no shared-fake refactor |

## Rollback Plan

Each slice is additive and self-contained: revert its commit/PR (new `plan.go`, `gestionar_plan.go`,
`ejecutar_hitos.go`, and their tests). No migration, schema, Contract, port, or existing use case is
touched, so reverting either slice leaves `main` deployable and slice 1 stands alone without slice 2.

## Dependencies

- #2 economic calendar, #8 motor, #15 repositories (all already merged). Blocks #24.

## Success Criteria

- [ ] No consent → no plan, exactly one door-open reminder.
- [ ] Consent → `ACTIVO` plan with `consentimiento_en`, valid frequency, and motor-designed milestones.
- [ ] Convertible lead → milestone 1 is `AFILIACION`; monetary milestones anchor to calendar events.
- [ ] Pause honored immediately; paused plans never message again.
- [ ] Milestone text offers pause and contains no distress presumption.
- [ ] Failure ordering tests prove no milestone is marked without a successful send and append.
- [ ] Gap closure hands off to `PERFILANDO` with one `PerfilCompleto` event.
- [ ] `go test ./internal/usecase/... ./internal/domain/motor/...` passes per slice.

## Proposal question round

Non-interactive execution; these assumptions need confirmation before spec:
1. Is `PrecioObjetivo` an explicit caller input (assumed yes) rather than derived from the top kNN recommendation?
2. Is cumulative `NOTIFICADO` milestone amount the accepted deterministic gap-closure signal (assumed yes)?
3. Should a no-consent reminder failure surface as an error (assumed yes) instead of being swallowed as in the issue sketch?
