# Design: ProcesarMensaje Conversational Turn (Issue #19)

## Technical Approach

Add a synchronous, provider-free application service and question helpers in new `internal/usecase` files only. Existing `LLMProvider`, repositories, domain vocabulary, motor, Contract v1.1, and runtime wiring remain unchanged. The issue sample is stale where it names nonexistent symbols, adds `Catalogo`/`TasaEA`, saves before the provider, or emits `PausarContacto`; current code and the approved spec govern.

## Architecture Decisions

| Topic | Choice and rationale | Rejected |
|---|---|---|
| Boundary | `ProcesarMensaje` depends only on `LeadRepository`, `LLMProvider`, `GeneradorID`, `BusEventos`, and `Reloj`; all public additions live in the three new files. | Port/domain/HTTP/ADK edits. |
| Validation | Export sentinels `ErrValidacion` (`VALIDACION`), `ErrAudioInvalido` (`AUDIO_INVALIDO`), `ErrLeadNoPerfilando` (`TRANSICION_INVALIDA`), `ErrLimiteTurnos` (wrapping `ErrValidacion`), and `ErrSalidaTurnoInvalida`. Validate context/ID/type first; text is valid UTF-8, nonblank, ≤2,000 runes; audio is nonempty strict base64, ≤2 MiB decoded, 1–60 s, and allowed MIME. | Mutating or calling the provider before validation. |
| History/cap | Load the full chronological conversation; count every persisted `AutorMensajeLead`. Current turn is `count+1`; reject above 4/6 before dispatch. Pass only the last six existing messages plus current `MensajeUsuario`, matching the issue’s bounded context without counting Vivi messages as turns. | Counting message pairs or truncated history. |
| LLM input | Supply current lead/profile/capacity, affiliation, history, and `NumerosDelMotor={presupuesto_max,credito_max,subsidio_aplicable,recursos_propios}`. Call exactly one modality method. New motor values become next-turn context; the same call cannot safely quote post-merge values. | A second extraction/response call or usecase retry. |
| Merge | Two-pass `aplicarCampos`: reject duplicate recognized keys, invalid source/confidence/type, fractional integer fields, and invalid Contract enums before mutation; ignore unknown keys; preserve existing `VERIFICADO_BASE`; downgrade provider-claimed `VERIFICADO_BASE` to `DECLARADO`; normalize integral JSON numbers to `int64`. | Partial merge or invented profile keys. |
| Completion | Recalculate once with `motor.CalcularCapacidad(perfil, lead.Afiliado, 0)`. Complete on `PERFIL_COMPLETO`, `PerfilEstaCompleto`, or cap+`CONTINUAR`; safety/consent/out-of-domain actions prevail. Transition only `PERFILANDO→CALIFICADO`; publish `EvPerfilCompleto` after both messages are durable. Rejecting non-`PERFILANDO` leads prevents duplicate state/event emission. | Event-only completion or ad-hoc pause event. |
| Audio | Validate/decode only in memory, call `ProcesarAudio` once, persist caller `Texto` with `audio_original=true`, never fabricate transcription. `AUDIO_ININTELIGIBLE` stores only `salida.Respuesta`; profile, capacity, intention, state, inbound history, and event remain unchanged. | Storing raw audio or Vivi text as transcript. |

## Exact Sequence

```mermaid
sequenceDiagram
  Caller->>UC: Ejecutar(input)
  UC->>Repo: PorID; Conversacion
  UC->>LLM: GenerarTurno OR ProcesarAudio (once)
  UC->>UC: validate output; merge; motor; optional transition
  UC->>Repo: AgregarMensaje(LEAD)
  UC->>Repo: Guardar(lead CAS)
  UC->>Repo: AgregarMensaje(VIVI)
  UC->>Bus: PerfilCompleto (only if transitioned)
```

There is no cross-write transaction. Provider/validation failure writes nothing; persistence failure returns immediately, publishes no event, and may leave the documented prefix durable—no compensating writes.

## Interfaces and Helpers

```go
type EntradaMensaje struct { LeadID string; Tipo domain.TipoMensaje; Texto string; Audio *Audio }
type ProcesarMensaje struct { Leads LeadRepository; LLM LLMProvider; IDs GeneradorID; Bus BusEventos; Reloj Reloj }
func (uc *ProcesarMensaje) Ejecutar(context.Context, EntradaMensaje) error
```

`prioridadCampos` is exactly `ingreso_hogar, recursos_propios, zona_deseada, plazo_compra_meses, arriendo_actual, personas_hogar`; `camposCriticosPerfil` is the first three. `SiguienteMejorPregunta` skips present/verified fields; `PerfilEstaCompleto` requires all three. Private helpers: `validarEntrada`, `contarTurnosLead`, `historialReciente`, `numerosDelMotor`, `validarCampos`, `aplicarCampos`, and `responder`.

## File Changes

| File | Action | Responsibility |
|---|---|---|
| `internal/usecase/procesar_mensaje.go` | Create | DTO/service, validation, sequence, merge, persistence. |
| `internal/usecase/siguiente_pregunta.go` | Create | Exact priority/completion vocabulary. |
| `internal/usecase/procesar_mensaje_test.go` | Create | Deterministic tests and local counting/failure fakes. |

## Testing, Rollback, and Non-Goals

Use fixed clock/IDs, existing in-memory repo, local scripted LLM and fail-at-step repository/bus. Matrix: each invalid input/state/cap; text/audio one-call and provider errors; six-message history; unknown/duplicate/type/source fields and verified preservation; one motor refresh; unintelligible audio; 4/6 caps; complete/terminal actions; CAS/response failure prefixes; event post-durability and no duplicate. Run focused tests, then `go test ./...`, build, vet, and formatting.

Threat matrix: N/A—no shell, routing, subprocess, VCS, or process integration. No migration/flag. Roll back by deleting/reverting the three new files. Tasks must measure before apply; >400 authored lines requires chained slices, with no size exception. Non-goals: HTTP/frontend, ADK, Contract/port/domain/repository changes, provider-returned transcription, retries, routing, recommendations, ficha, nutrition, migrations, or external services.
