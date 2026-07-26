# Proposal: ProcesarMensaje conversational turn (Issue #19)

## Intent

After #18 every lead stops at `PERFILANDO`: nothing consumes a message, so US-03 (value exchange), US-04 (voice notes), US-05 (colloquial input) and US-06 (intent) are unreachable and RF-M4-02/07 are unimplemented. Deliver one usecase-only synchronous turn — validate, **one** LLM call (NFR-R-01, p95 < 3 s), merge extraction, recalculate capacity with the motor, persist ordered messages, complete once — with no orchestration layer.

## Scope

### In Scope
- `internal/usecase/siguiente_pregunta.go`: `prioridadCampos`, `SiguienteMejorPregunta`, `PerfilEstaCompleto` — never returns a `VERIFICADO_BASE` field (RF-M4-02).
- `internal/usecase/procesar_mensaje.go`: `EntradaMensaje{LeadID,Tipo,Texto,Audio}`, `ProcesarMensaje{Leads,LLM,IDs,Bus,Reloj}`, `Ejecutar`, entry validation, `aplicarCampos`, `numerosDelMotor`, turn counting, bounded completion.
- `internal/usecase/procesar_mensaje_test.go`: table-driven tests with local counting/scripted LLM fake plus existing repo/clock/ID fakes; no DB, no real provider, no wall clock.

### Out of Scope
ADK graph and `internal/adapters/agentes`, HTTP route and frontend/polling, any edit to `LLMProvider`/`EntradaTurno`/`SalidaTurno`/`Audio`/repositories/domain/Contract, kNN, recommendations, routing, ficha, nutrition plan, migrations, usecase-level provider retry, and provider-returned audio transcription.

## Capabilities

### New Capabilities
- `procesar-mensaje`: validated turn, single-call dispatch, extraction merge, capacity recalculation, turn limits, message ordering, one-time completion.

### Modified Capabilities
- None. `capacidad`, `llm-guardrails` and `perfilar-lead` are consumed unchanged.

## Approach

Issue snippets are stale; merged code wins (doc 13 owns money, Contract v1.1 owns vocabulary):

| Snippet | Resolution |
|---|---|
| `motor.CalcularCapacidad(..., ParametrosCredito{TasaEA})` | Real 3-arg `motor.CalcularCapacidad(perfil, lead.Afiliado, 0)`; drop `TasaEA` field |
| `domain.Campo*`, `FuenteVerificadoBase`, `AutorLead`, `ContenidoTexto` | Contract string keys as unexported usecase constants; `FuenteCampo*`, `AutorMensajeLead`, `TipoContenidoTexto` |
| `Catalogo CatalogoRepository` field | Dropped — unused at this boundary |
| `Bus.Publicar(Evento{Tipo:"PausarContacto"})` | Not a Contract §6 event; `PAUSAR_CONTACTO` only stores the response, no state change, no ad-hoc event (consent/route lifecycle is #20) |
| `camposCriticosPerfil` vs `domain.CamposCriticos` | Keep the issue trio (`ingreso_hogar`, `recursos_propios`, `zona_deseada`) as the conversational completion gate; `domain.CamposCriticos` stays the motor-side set, unchanged |

Frozen decisions:

| Topic | Decision |
|---|---|
| Ordering | Validation → `PorID` → `Conversacion`/turn count → **one** provider call → inbound `LEAD` message → lead CAS `Guardar` → `VIVI` response → event. Deviates from issue step 2 (message before LLM) because ports expose no transaction: a provider failure must not consume a turn or leave a reply-less message, keeping retries safe. No DoD item constrains this order. |
| State/event de-dup | Completion only via `Lead.Transicionar(PERFILANDO→CALIFICADO)`; `PerfilCompleto` published exactly once, after both messages and the lead are durable. A lead already `CALIFICADO`+ never republishes. |
| Validation | `TEXTO`: non-blank, ≤ 2 000 chars. `AUDIO`: base64 decodes, ≤ 2 MB decoded, `duracion_s` 1–60, MIME ∈ {webm, ogg, mpeg} → else typed `VALIDACION`/`AUDIO_INVALIDO`. Lead must be `PERFILANDO`. Invalid input makes no provider call and no write. |
| Turn limit (RF-M4-07) | Count persisted `AutorMensajeLead` messages; accepted input is `existing+1`; max 4 affiliate / 6 non-affiliate. Over the cap → rejected before the provider call. At the cap with `CONTINUAR` → complete with the fields available; terminal safety actions stay authoritative. |
| Audio (§2.5 / RF-M4-03) | `EntradaMensaje.Texto` carries the **caller-supplied** transcription; `Audio` is passed once to `ProcesarAudio`, kept in memory and discarded. Stored as `LEAD/TEXTO` + `adjunto{audio_original:true}`; never `SalidaTurno.Respuesta` as a transcript, never fabricated. `AUDIO_ININTELIGIBLE` leaves profile, capacity and state untouched and only stores the repeat-or-write reply. Qwen's typed unsupported-audio error surfaces as-is. |
| Transcription gap | Full provider-returned transcription needs a **separate** `contrato`-labelled change: add `transcripcion` to `SalidaTurno`, update parser, guardrails, Gemini/Qwen adapters and tests. Explicitly not in #19. |
| Extraction integrity | Merge only `CamposReconocidos` keys; skip `EsVerificado` fields; provider-claimed `VERIFICADO_BASE` downgrades to `DECLARADO`; unknown keys ignored (Contract §2.1); duplicate keys rejected as invalid output. |

## Affected Areas

| Area | Impact | Description |
|---|---|---|
| `internal/usecase/procesar_mensaje.go` | New | Turn usecase, validation, merge, completion |
| `internal/usecase/siguiente_pregunta.go` | New | Next-best-question priority and completion gate |
| `internal/usecase/procesar_mensaje_test.go` | New | US-03..06 scenarios, local LLM fake |
| `internal/usecase/puertos.go`, `internal/domain/**`, `internal/infrastructure/llm/**` | Read-only | Consumed unchanged |

## Risks

| Risk | Likelihood | Mitigation |
|---|---|---|
| Non-atomic lead CAS + two message inserts leaves partial state | Med | Event only after full durability; failure path tested and documented; transaction/idempotency key is a follow-up |
| Contract §2.5 transcription unmet without caller transcript | High | Empty `texto` + `audio_original` flag, no fabrication; separate DTO/Contract change named above |
| Measured size exceeds the 400 authored-line cap | High | `sdd-tasks` must measure before apply; chained slices, never a size exception |
| Rejecting turns outside `PERFILANDO` blocks post-qualification chat | Med | Deliberate minimum boundary; re-entry belongs to #20 routing/nutrition |
| Forced completion at the cap with thin profile lowers `Confianza` | Low | RF-M4-07 accepts it; confidence comes from the motor |

## Rollback Plan

All three files are new and nothing imports them (`cmd/servidor` keeps its unwired Block A seam). Revert the single PR commit or delete `feat/issue-19-procesar-mensaje`; there is no migration, config, data, port, DTO or Contract change to undo. #18 stays functional because `PerfilarLead` is untouched.

## Dependencies

Merged #7 (state machine), #8 (motor), #11 (ports), #15 (repositories), #16/#17 (providers, guardrails), #18 (`PERFILANDO` leads). Blocks #20 and #24.

## Success Criteria

- [ ] Exactly one provider call per accepted turn, asserted by fake invocation count.
- [ ] `VERIFICADO_BASE` never overwritten and never re-asked; unknown keys ignored without error.
- [ ] Capacity recalculated every turn via the 3-arg motor; `NumerosDelMotor` always injected.
- [ ] `AUDIO_ININTELIGIBLE` leaves profile, capacity and state unchanged.
- [ ] Caps 4/6 enforced; `PerfilCompleto` published once, post-durability, never twice.
- [ ] Invalid text/audio and non-`PERFILANDO` leads: no call, no write, typed error.
- [ ] `go test ./internal/usecase/... -run 'TestTurno|TestAplicarCampos|TestNoRepregunta' -v`, then `go test ./...`, `go build ./...`, `go vet ./...`, `gofmt -l` clean.
- [ ] Measured authored lines ≤ 400 before apply, or chained delivery.

## Proposal question round

1. Confirm audio without a caller transcription may persist an empty `texto` (flagged `audio_original`) instead of blocking the turn.
2. Should the completion gate be the issue trio (`ingreso_hogar`, `recursos_propios`, `zona_deseada`) or `domain.CamposCriticos`?
3. Reject turns on a `CALIFICADO`+ lead, or accept them without further profiling?
4. Is losing the lead's message on provider failure (persist-after-call) preferable to consuming a turn (persist-before-call)?
5. Should `PAUSAR_CONTACTO` stay event-free until #20?

Assumed if unanswered: 1 yes, 2 issue trio, 3 reject, 4 persist-after-call, 5 event-free.
