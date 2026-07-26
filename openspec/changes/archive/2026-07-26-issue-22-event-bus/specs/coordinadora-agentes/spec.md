# Agent Coordinator Specification

## Purpose
Define deterministic event handoff among existing use cases without duplicating their ownership.

## Requirements

### Requirement: F2 Canonical Route Decision Ownership
The coordinator MUST consume `RutaDecidida` and MUST NOT publish, synthesize, or republish it. `CalificarLead` remains its sole producer.

#### Scenario: Qualification handoff
- GIVEN a completed `PerfilCompleto` qualification
- WHEN the coordinator invokes the qualifier
- THEN it does not emit an additional route-decision event

### Requirement: F7 Deterministic Registration Table
`Registrar` MUST install a fixed table: `LeadNuevo` observe-only; `PerfilCompleto` qualifies by lead ID; `RutaDecidida` evaluates ficha eligibility; `TickReloj` dispatches milestones from `hasta`. It MUST NOT subscribe unlisted event types or fabricate greeting, pause, or plan work.

#### Scenario: Repeatable routing
- GIVEN equivalent registered coordinators
- WHEN each receives ten identical events of one supported type
- THEN each invokes the same designated handler ten times
- AND no LLM determines the route

### Requirement: F8 Reprofile Qualification Skip
A `PerfilCompleto` qualification returning `ErrLeadNoCalificable` MUST be treated as a side-effect-free skip. The coordinator MUST NOT generate a ficha or publish an event for that skip.

#### Scenario: Profiling reprofile
- GIVEN a `PERFILANDO` lead emits `PerfilCompleto`
- WHEN qualification returns `ErrLeadNoCalificable`
- THEN no ficha, event, or durable route change is produced

### Requirement: F9 Advisor-Only Ficha Handoff
For `RutaDecidida`, the coordinator MUST invoke ficha generation only for `ASESOR`. It MUST accept `ruta` as a route value or string; missing, invalid, or `NUTRICION` routes MUST generate neither ficha nor automatic plan. Handoffs MUST carry only lead ID.

#### Scenario: Nutrition route
- GIVEN a `RutaDecidida` event whose route is `NUTRICION`
- WHEN the coordinator handles it
- THEN it creates neither ficha nor plan

### Requirement: F10 Safe Partial Wiring and Payloads
Nil injected dependencies MUST omit only their dependent subscription and MUST NOT panic. A nil, missing, or malformed payload field required by a route MUST be a no-op and MUST NOT trigger an unintended handler.

#### Scenario: Partial coordinator
- GIVEN a coordinator without a ficha generator
- WHEN it registers and receives `RutaDecidida` with nil payload
- THEN registration and publication complete normally
- AND no ficha handler is invoked
