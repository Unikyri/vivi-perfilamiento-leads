## Exploration: issue-21-nutrition-plan

### Current State
Issue #21 supplies the missing provider-free nutrition application layer after #18–#20. The existing Contract-shaped `domain.PlanNutricion` and `domain.Hito`, `PlanRepository`, `LeadRepository`, `Reloj`, `GeneradorID`, and `BusEventos` ports already cover the required boundary. `postgres.PlanRepository` atomically persists a plan with its milestones, retrieves only pending milestones from active plans, and already has schema support.

`CalificarLead` leaves the NUTRICION route in `EN_NUTRICION` but deliberately creates no plan. `ProcesarMensaje` only accepts `PERFILANDO` and currently answers `PAUSAR_CONTACTO` without changing a plan or lead. No nutrition use case, concrete bus, production simulated clock, or reusable in-memory plan/gateway fake exists. The canonical economic calendar is parsed by `domain.LoadCalendario`; its CESANTIAS and PRIMA dates must not be copied into code.

### Authority Reconciliation
The Issue #21 body and Wiki authority require explicit consent plus a valid contact frequency; no consent creates no plan and yields only one door-open reminder. Milestones are deterministic and calendar-anchored; a conversion-capable lead receives `AFILIACION` first. Tick execution must use the simulated clock, make paused plans silent, offer pause in every nutrition message without presuming financial distress, and hand off successful requalification without duplicating #20 route events.

### Affected Areas
- `internal/domain/motor/` — add a pure deterministic milestone planner with explicit inputs.
- `internal/usecase/` — add provider-free create/consent, pause, and due-milestone services with focused tests.
- `internal/usecase/fakes_test.go` or nutrition-local fakes — supply deep-copying plan/gateway failure doubles within the slice budget.
- Existing ports, migration, Postgres repository, calendar inputs, qualification, HTTP, ADK, and composition root remain dependencies and must not be changed by this issue.

### Risks and Decisions Required
1. `PlanNutricion` exposes only `MetaMonto`/`MetaDescripcion`; it has no target-project, reminder, version, or pause-reason field. The implementation MUST not invent Contract fields.
2. The target price is not persisted on `Lead`; define a deterministic target/brecha input and clamp negative gaps to zero. Do not silently infer a new persisted contract.
3. Plan and lead writes, and gateway/message/mark writes, cross repository boundaries without a transaction. Specify durable ordering, retry/idempotency, and mark-after-success tests.
4. `TickReloj` and `Reloj` are ports only. #21 MUST expose callable use cases and defer event subscription, bus, concrete clock, HTTP, and coordinator wiring to #22/#23/#30.
5. The existing `ProcesarMensaje` pause behavior is an integration handoff; broadening it is outside the safe first slice.
6. Requalification after a milestone is a handoff to the established profiling/qualification flow, not a direct duplicate invocation of `CalificarLead` on an `EN_NUTRICION` lead.

### Approaches
1. **Application orchestration plus a pure planner (recommended).** Add provider-free `GestionarPlan` and `EjecutarHitos`; isolate dates and amounts in a pure planner consuming explicit calendar, target, capacity, and conversion inputs.
   - Pros: deterministic, auditable, table-testable, reuses existing ports, preserves Clean Architecture.
   - Cons: requires explicit target and partial-failure semantics.
2. **Immediate event/coordinator integration.** Subscribe now to route/tick/message events.
   - Rejected: couples this change to the missing #22/#23 implementation, obscures failure ordering, and risks duplicate route publication.

### Recommendation
Use the first approach and deliver two chained, independently testable slices. Slice 1 contains the pure planner and explicit-consent plan creation. Slice 2 contains pause and due-milestone execution with gateway/message/mark ordering. No LLM, ADK, HTTP, frontend, migration, Contract change, new port, or repository rewrite is justified.

### Safe Delivery Forecast
- **Slice 1 — planner + consent/create:** forecast 140–165 runtime plus 190–220 tests, 330–385 authored lines.
- **Slice 2 — tick execution + pause:** forecast 145–175 runtime plus 180–215 tests, 325–390 authored lines.

Both remain below the hard 400-line review budget. Use a feature-branch chain and separate native review/verification evidence per slice.

### Ready for Proposal
Yes. Proposal must freeze the explicit target/brecha and failure-order rules from this exploration before implementation.