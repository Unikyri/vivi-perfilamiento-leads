# Exploration: Personalized Greeting and Data Consent (Issue #29)

## Current state
`SaludarLead` is already a deterministic `LeadNuevo` subscriber, but its affiliate copy hard-codes the subsidy and does not include the policy link. `ProcesarMensaje` recognizes consent actions but can merge extracted fields before handling a denial. The current lifecycle lacks `PERFILANDO -> DESPEDIDO`.

## Decision
Implement exactly one deterministic greeting path through the existing coordinator wiring; do not add an LLM dependency. Derive the subsidy only from `lead.Capacidad.SubsidioAplicable`; if it is unavailable, do not claim an amount. Use a declarative consent request plus policy URL so the greeting contains only one discovery question. On `CONSENTIMIENTO_NO`, do not merge profile fields or calculate capacity, append one dignified farewell, transition to `DESPEDIDO`, and persist the lead. This is data-treatment consent, distinct from the existing nutrition-plan consent.

## Constraints
No provider calls, schema migration, or coordinator routing expansion. The documented provider-failure fallback in the issue is not applicable to the present deterministic greeting and must remain successful without a provider. Preserve one greeting per lead under duplicate event delivery.

## Evidence
Issue #29 body, current `saludar_lead.go`, `procesar_mensaje.go`, coordinator/composition wiring, domain transition tests, and persisted Engram exploration #632.