# Design: Personalized Greeting and Data Consent

## Technical Approach

Extend the existing `SaludarLead` `LeadNuevo` observer rather than coordinator routing. It builds segment rules and an exact deterministic fallback, makes at most one application-level `LLMProvider.GenerarTurno` call, accepts only `SalidaTurno.Respuesta` when a pure validator proves every rule, then persists the accepted draft or fallback. Provider errors, nil provider, empty output, and invalid output are fallback conditions, not use-case failures.

`ProcesarMensaje` detects `CONSENTIMIENTO_NO` immediately after provider output and delegates to `SaludarLead.RechazarConsentimiento` before output-field normalization, profile merge, capacity calculation, intention assignment, or completion publication.

## Architecture Decisions

| Decision | Choice and rationale | Rejected alternative |
|---|---|---|
| Greeting composition | Reuse `LLMProvider.GenerarTurno` and existing `EntradaTurno`; pass a controlled greeting brief plus name, affiliation, capacity, and `NumerosDelMotor`. Existing adapters supply the system instruction. Ignore returned fields/action/intention. This adds no provider-specific API. | Provider adapter or port changes broaden scope; deterministic-only composition contradicts the approved proposal. |
| Acceptance boundary | Construct fallback first; validate the LLM draft with the same segment rule set. Invalid drafts never reach persistence. | Prompt-only compliance is nondeterministic. |
| Subsidy source | Affiliate treatment applies only when `Afiliado` and `Capacidad.SubsidioAplicable > 0`; format that integer as Colombian millions (`52_500_000` → `$52,5M`) without floating point. Zero/nil uses non-affiliate amount-free copy. | Hard-coded or inferred subsidy violates motor authority. |
| Denial lifecycle | Persist the inbound denial, transition the unmodified loaded lead through `CALIFICADO` then `DESPEDIDO`, set `RutaDespedida`, CAS-save once, append one fixed farewell, and return without bus publication. | A direct lifecycle edge or consent field changes Contract v1.1. |

## Composition and Data Flow

```mermaid
sequenceDiagram
  participant B as LeadNuevo bus/coordinator
  participant S as SaludarLead
  participant L as LLMProvider
  participant R as LeadRepository
  B->>S: Ejecutar(event)
  S->>R: PorID
  S->>S: rules + deterministic fallback
  S->>L: GenerarTurno once
  L-->>S: draft or error
  S->>S: validate; otherwise fallback
  S->>R: AgregarMensaje(VIVI)

  participant P as ProcesarMensaje
  P->>L: classify turn
  alt action == CONSENTIMIENTO_NO
    P->>S: RechazarConsentimiento(lead, entrada)
    S->>R: AgregarMensaje(LEAD denial)
    S->>S: PERFILANDO→CALIFICADO→DESPEDIDO; RutaDespedida
    S->>R: Guardar(lead)
    S->>R: AgregarMensaje(VIVI farewell)
  end
```

The validator MUST require nonblank text, the exact `URLPolitica`, the declarative clause `Al continuar autorizas el tratamiento de tus datos`, and `strings.Count(text, "?") == 1`. Affiliate output MUST contain the lead name, exact formatter output, a dream-question marker (`sueñ` case-insensitive), no other monetary expression, and no income/salary/household prompt. Non-affiliate/zero-subsidy output MUST contain no monetary expression and MUST contain a job-situation marker (`empleo`, `trabajo`, or `laboral`). Templates use those same exact clauses. `URLPolitica` is the single issue/spec-owned constant; implementation MUST NOT invent or substitute a placeholder URL.

## Error Boundaries

Greeting read/persistence errors are returned; provider and validation failures persist the fallback. Denial normalizes no returned fields, calls no motor function, and emits no event. Writes remain intentionally non-transactional: on denial, each failure stops later writes, no compensating write is attempted, and an error is returned. Success guarantees one retained inbound denial, final `DESPEDIDO`/`DESPEDIDA`, one farewell, and unchanged `Perfil`, `Capacidad`, and `Intencion`.

## File Changes

| File | Action | Description |
|---|---|---|
| `internal/usecase/saludar_lead.go` | Modify | LLM dependency, templates/formatter/validator, denial method and ordered writes. |
| `internal/usecase/saludar_lead_test.go` | Modify | Table-driven composition, validation, fallback, and persistence tests. |
| `internal/usecase/procesar_mensaje.go` | Modify | Inject greeting collaborator and early terminal denial branch. |
| `internal/usecase/procesar_mensaje_test.go` | Modify | Denial immutability, ordering, state/route, farewell, and zero-event assertions. |
| `cmd/servidor/main.go` | Modify | Inject the existing provider into `SaludarLead` and the same instance into `ProcesarMensaje`; no coordinator-table change. |

No files are created or deleted; `puertos.go`, provider adapters, domain lifecycle, repositories, DTOs, migrations, and coordinator registration remain unchanged.

## Testing Strategy

Use table-driven fakes only—no Gemini/Qwen/API calls. Greeting cases cover valid affiliate/non-affiliate drafts, actual motor amount, zero subsidy, extra question, forbidden fields/amounts, missing policy/consent, empty draft, and provider error; assert one provider call maximum and persisted fallback parity. Denial tests return malicious extracted fields with `CONSENTIMIENTO_NO` and deep-compare profile/capacity/intention, assert inbound metadata, farewell, `DESPEDIDO`, `RutaDespedida`, zero motor-observable change, and zero events; inject write/CAS failures to prove stop ordering. Preserve existing coordinator `LeadNuevo` observer tests. Run `go test ./internal/usecase/... ./internal/adapters/agentes/... ./cmd/...`, then `go test ./...`, `go build ./...`, and `go vet ./...`.

## Threat Matrix

N/A — no routing, shell, subprocess, VCS/PR automation, executable classification, or process-integration boundary.

## Migration / Rollout

No migration or feature flag. Single Block A PR; rollback is revert-only.

## Open Questions

None architecturally. Apply remains blocked until the spec supplies the literal `URLPolitica`; no placeholder is permitted.
