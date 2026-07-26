# Exploration: PerfilarLead (Issue #18)

## Context
Issue #18 implements deterministic UC-01 steps 1–3 and UC-03 family-affiliate re-consultation. It depends on the merged state machine, capacity motor, application ports, and repositories. The implementation belongs entirely in `internal/usecase`; HTTP, LLM, messaging, ADK, greetings, recommendations, and persistence adapters remain out of scope.

## Confirmed contracts
- `LeadRepository.Crear` persists a new `NUEVO` lead and initializes version; `Guardar` is CAS for existing leads.
- `CatalogoRepository.AfiliadoPorCedula` exposes synthetic affiliate data; a missing/inactive affiliate selects the non-affiliate path rather than failing creation.
- `motor.CalcularCapacidad` is pure, takes `(Perfil, esAfiliado, precioCandidato)`, and owns all monetary rules.
- `Lead.Transicionar` is the only valid lifecycle path; profiling ends in `PERFILANDO`.
- Contract events include `LeadNuevo`; publish only after successful persistence.
- The project contract requires the ratio denominator to be the lowest affordable positive catalog price, falling back to the motor median when none exists.

## Policy decision
The exact Issue #18 acceptance narrative defines the demo baseline that recognized affiliates have no home and no prior subsidy for initial pre-profiling. This use case records `tiene_vivienda=false` and `recibio_subsidio=false` as `VERIFICADO_BASE` **only for this synthetic demo affiliate baseline**, enabling the mandated 30-SMMLV Ana result. No inferred values or user-provided values are reclassified as verified; later production data integration must supply explicit eligibility fields before retaining this behavior.

## Design boundary
Create one provider-free application service with `EntradaPerfilar`, `SalidaPerfilar`, `Ejecutar`, and `ReconsultarPorFamiliar`. It creates and persists a new lead, optionally maps verified affiliate fields, computes capacity using the existing motor, transitions to `PERFILANDO`, and emits a minimal `LeadNuevo` event only after the create succeeds. Re-consultation loads the existing lead, marks household affiliation, adds family income, recalculates/persists capacity, or records the unknown family cedula as declared/confirmation-required.

## Risks and mitigations
- Never directly set a state transition; use `Lead.Transicionar`.
- Preserve repository errors and never publish successful events following a failed create/save.
- Tests use only deterministic test fakes: no database, LLM, HTTP, or wall clock.
- Keep implementation/tests below the 400-line review budget; no new port/config/Contract/API changes.

## Validation
Focused usecase tests must prove affiliate hit, missing affiliate, inactive affiliate, family hit/miss, verified provenance, 52,527,150 subsidy, persistence/event ordering, and no real dependencies. Repository-wide Go test/build/vet/module/diff checks precede delivery.
