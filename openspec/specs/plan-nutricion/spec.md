# Plan Nutricion Specification

## Purpose
Define provider-free, consented nutrition plans and deterministic simulated-time advancement using the existing Contract boundary.

## Requirements

### Requirement: F1 — Provider-Free Deterministic Planning
The system MUST expose only callable use cases and a pure planner, using existing fields and ports. It MUST NOT add integrations, wiring, ports, or domain fields. Milestones MUST derive deterministically from explicit gap, conversion, start date, and parsed economic calendar; monetary dates SHALL come from that calendar, and a conversion-capable lead MUST receive `AFILIACION` first.

#### Scenario: Deterministic, isolated plan
- GIVEN identical explicit planner inputs and calendar data
- WHEN milestones are designed repeatedly
- THEN their ordered values are identical and calendar-anchored
- AND no external provider, HTTP, ADK, or persisted field is required

### Requirement: F2 — Explicit Target and Gap Validation
The caller MUST provide `PrecioObjetivo`. The system MUST calculate `brecha = max(0, PrecioObjetivo - lead.Capacidad.PresupuestoMax)` and MUST reject a missing capacity or non-positive target with no plan, message, or state change.

#### Scenario: Target validation and clamped gap
- GIVEN a valid capacity and an objective above or below its budget, or an invalid target input
- WHEN plan creation is requested
- THEN the gap is the positive difference or zero, respectively
- AND invalid input returns validation failure with zero side effects

### Requirement: F3 — Consent and Frequency Gate
A plan MUST require `consintio == true` and frequency `QUINCENAL`, `MENSUAL`, or `TRIMESTRAL`. On valid creation, it MUST be `ACTIVO` and set `ConsentimientoEn` to `Reloj.Ahora()`; invalid frequency MUST have no side effects.

#### Scenario: Valid consent gate
- GIVEN consent with a permitted or unpermitted frequency
- WHEN creation is requested
- THEN only the permitted case creates an active, timestamped plan
- AND the other returns validation failure without message or state change

### Requirement: F4 — No-Consent Reminder
Without consent, the system MUST append exactly one door-open reminder, return `(nil, nil)`, and MUST NOT create a plan or transition the lead. A failed append MUST return its error and MUST NOT be swallowed.

#### Scenario: Consent declined or absent
- GIVEN `consintio == false` and either successful or failed reminder persistence
- WHEN creation is requested
- THEN exactly one append is attempted and no plan or transition occurs
- AND success returns `(nil, nil)` while failure returns the append error

### Requirement: F5 — Durable, Idempotent Creation Ordering
Creation MUST order `Planes.Crear`, lead transition to `EN_NUTRICION`, then `Leads.Guardar`. A create failure MUST leave the lead untouched. If lead save fails, the persisted plan MUST remain; a retry MUST reuse it through `Planes.PorLead`, never create a duplicate.

#### Scenario: Cross-repository creation failure
- GIVEN failures at plan creation or lead save followed by a retry
- WHEN creation executes
- THEN create failure causes no lead mutation
- AND save failure retains one plan that the retry reuses

### Requirement: F6 — Due-Milestone Delivery Ordering
For each due milestone of an active plan, the system MUST order `Gateway.Enviar`, `Leads.AgregarMensaje`, then `Planes.MarcarHito(NOTIFICADO)`. Nutrition text MUST offer pausing in its final line and MUST NOT presume financial distress. On either failure it MUST leave that milestone `PENDIENTE`, continue others, count only successes, return aggregated error, and defer retry to the next tick.

#### Scenario: Partial tick failure
- GIVEN due milestones with send or append failures
- WHEN the tick executes
- THEN only milestones with successful send and append are marked and counted
- AND failed milestones remain pending while independent milestones continue

### Requirement: F7 — Immediate, Respectful Pause
Pausing MUST set the plan to `PAUSADO`, transition the lead `EN_NUTRICION` to `PAUSADO`, then append one farewell. No plan MUST succeed silently; an already paused plan MUST be idempotent. An illegal lead transition MUST return an error without farewell, and paused plans MUST produce no later tick messages.

#### Scenario: Pause variants
- GIVEN absent, active, already paused, or illegally transitionable plans
- WHEN pause is requested and later ticks run
- THEN only a valid active plan receives one farewell after both state updates
- AND all other cases are silent or error as specified, with no paused-plan delivery

### Requirement: F8 — Requalification Handoff
When cumulative `NOTIFICADO` milestone amounts reach `plan.MetaMonto`, or `REEVALUACION` becomes `NOTIFICADO`, the system MUST transition `EN_NUTRICION` to `PERFILANDO`, save, and publish exactly one `PerfilCompleto` per lead per tick. It MUST NOT call `CalificarLead` or publish `RutaDecidida`; a nil bus MUST be silently skipped.

#### Scenario: One deterministic handoff
- GIVEN one or more successful milestones meeting either requalification condition
- WHEN a tick completes
- THEN the lead is durably `PERFILANDO` before at most one `PerfilCompleto`
- AND no qualification call or route event occurs, including with a nil bus
