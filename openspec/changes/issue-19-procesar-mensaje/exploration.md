# Exploration: ProcessMessage conversation turn (Issue #19)

## Current State

Issue #19 delivers the usecase-only synchronous conversational turn after Issue #18 created `PERFILANDO` leads. Existing ports already provide `LLMProvider`, `LeadRepository`, `BusEventos`, `Reloj`, and IDs. `internal/adapters/agentes` is only a package marker; no ADK graph is available or required for this usecase.

The authoritative APIs differ from stale issue snippets: profile keys are Contract strings in `domain.CamposReconocidos`, provenance constants are `FuenteCampoVerificadoBase`, `FuenteCampoDeclarado`, and `FuenteCampoInferido`, and capacity is calculated strictly with `motor.CalcularCapacidad(perfil, esAfiliado, precioCandidato)`.

## Minimum Boundary

- Add a provider-free `ProcesarMensaje` usecase and deterministic tests only.
- Validate text/audio entry before mutation; one `LLMProvider` call per accepted turn (`GenerarTurno` for text, `ProcesarAudio` for audio).
- Preserve verified profile fields; merge only recognized `DECLARADO`/`INFERIDO` extraction, recalculate capacity using the three-argument motor, persist intention and ordered messages, and emit one post-durability `PerfilCompleto` at action/completion cap.
- Count persisted inbound `LEAD` messages: cap 4 for affiliates and 6 otherwise; prevent duplicate completion events.
- Treat `MensajeEntrante` as the trigger and do not republish it, avoiding an event loop.

## Explicit Exclusions

No ADK graph, HTTP route, frontend, Contract/port/domain changes, routing, kNN, recommendations, ficha, plan, migration, external service, or usecase-level provider retry belongs in this issue. Existing providers/guardrails own retries, fallback, and security.

## Audio Constraint

Contract §2.5 calls for a persisted audio transcription, but the frozen `SalidaTurno` DTO and strict provider parser expose only extracted fields, intention, response, and action. This change must not fabricate a transcript or store Vivi's response as the lead message. The usecase accepts/persists only the caller-provided text transcription for audio (when available) with `audio_original=true`; raw audio is processed in memory and discarded. End-to-end transcription DTO/adapter evolution requires a separately approved Contract-compatible change.

## Failure Boundary

The repository exposes independent message and CAS operations, not a transaction. Validation/provider failure produces no writes. After a provider result, the safe ordering is inbound message → lead CAS save → Vivi response → completion event; no event precedes final message persistence. Partial persistence after a repository error is documented rather than hidden with compensating writes.

## Delivery Risk

The initial estimate was 440–580 authored lines, above the 400 review budget. Tasks must measure a compact implementation/test plan before apply; if it remains over 400, use chained delivery rather than a size exception.

## Recommendation

Proceed with proposal/spec/design/tasks for the isolated usecase boundary, preserving every existing port and Contract type. The final plan must define input validation, audio caller-transcript behavior, message ordering, event de-duplication, and measured review workload.
