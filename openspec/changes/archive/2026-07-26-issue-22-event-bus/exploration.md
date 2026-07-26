## Exploration: issue-22-event-bus

### Current State
`BusEventos` already exists as an application port, but has no production implementation. `PerfilarLead` publishes durable `LeadNuevo`; `ProcesarMensaje` publishes durable `PerfilCompleto`; `CalificarLead` already publishes durable `RutaDecidida`; and #21 supplies callable nutrition use cases without subscriptions. No greeting use case exists, so the coordinator must not invent one.

### Scope Reconciliation
Issue #22 owns a synchronous, in-memory Observer implementation and a plain-Go deterministic Mediator adapter. It must not edit ports, use cases, domain, repositories, HTTP, ADK, composition root, migrations, Contract, or frontend. The coordinator MUST NOT publish a second `RutaDecidida`, create an automatic nutrition plan, or directly qualify `EN_NUTRICION` reprofile events.

### Required Behavior
The bus must copy event payloads per handler, preserve registration order, release its lock before callbacks (nested publication safe), recover each handler panic independently, and log only event type, opaque lead ID, and handler outcome. The coordinator uses explicit injected callbacks/use cases and deterministic event-type registration; handoffs carry `lead_id`, not conversation history. No LLM judges routing.

### Approaches
1. **Pure sync bus plus callback-driven coordinator (recommended).** Add a panic-safe `bus.EnMemoria`, then a plain-Go coordinator that subscribes fixed handlers and delegates only valid normal-flow events.
2. **Direct ADK graph or automatic workflow integration.** Rejected: ADK is unvalidated, greeting is absent, and automatic plan/reprofile handling would duplicate ownership from #19–#21.

### Delivery Forecast
- Slice 1: new in-memory bus and tests, 320–385 runtime/test lines.
- Slice 2: new coordinator and tests, 310–375 runtime/test lines.

### Risks
- `Evento.Payload` is untyped; handlers must validate only fields they consume.
- The in-memory bus is intentionally non-durable and has no acknowledgment semantics.
- Normal profile completion and nutrition reprofile handoff must remain distinct to avoid invalid `CalificarLead` state use.

### Ready for Proposal
Yes. Use provider-free plain Go, two bounded slices, and no runtime integration outside this issue.