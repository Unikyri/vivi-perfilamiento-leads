# ProcesarMensaje Specification

## Purpose

Define deterministic turn for a `PERFILANDO` lead through existing ports.

## Requirements

### Requirement: Input validation has no side effects
The lead MUST be `PERFILANDO`. Text MUST be nonblank and ≤2,000 characters. Audio MUST base64-decode to ≤2 MB, last 1–60 seconds, and use `audio/webm`, `audio/ogg`, or `audio/mpeg`. Invalid input MUST return `VALIDACION` or `AUDIO_INVALIDO` with no provider call, message, lead write, or event.

#### Scenario: Invalid entry
- GIVEN blank text, invalid audio, or a non-`PERFILANDO` lead
- WHEN a turn is requested
- THEN a typed error returns with zero side effects

### Requirement: One modality dispatch
Text MUST call `GenerarTurno` once; audio MUST call `ProcesarAudio` once. Audio bytes MUST be discarded after the call. Caller `Texto`, including empty text, MUST be the sole audio transcript; Vivi's response MUST NOT replace it.

#### Scenario: Audio transcript
- GIVEN valid audio and caller text
- WHEN its turn succeeds
- THEN one audio call occurs and the inbound record uses that text with `audio_original=true`

### Requirement: Field integrity and motor refresh
Only recognized keys MUST merge. Existing `VERIFICADO_BASE` MUST remain; provider-claimed `VERIFICADO_BASE` MUST become `DECLARADO`; unknown keys MUST be ignored. Duplicate recognized keys or invalid provenance MUST fail before mutation. Turn input MUST include current motor values; a normal merge MUST recalculate capacity once with `CalcularCapacidad(updatedPerfil, lead.Afiliado, 0)`.

#### Scenario: Protected fields
- GIVEN a verified field, unknown key, and declared recognized value
- WHEN output is applied
- THEN the verified value remains, the unknown key is absent, and capacity refreshes

### Requirement: Unintelligible audio is profile-safe
For `AUDIO_ININTELIGIBLE`, the system MUST leave profile, capacity, and state unchanged, persist only the repeat-or-write Vivi response, and MUST NOT apply fields, invent a transcript, or emit completion.

#### Scenario: Unintelligible audio
- GIVEN valid audio with an unintelligible outcome
- WHEN the turn completes
- THEN only the repeat-or-write response is stored without profile mutation

### Requirement: Ordered persistence
A normal result MUST persist inbound `LEAD` message, CAS-save the lead, persist `VIVI` response, then publish completion. Validation/provider failure MUST write nothing; persistence failure MUST publish no event. Writes are non-atomic: partial persistence MAY remain and MUST NOT cause compensating writes.

#### Scenario: Response failure
- GIVEN inbound message and CAS save succeeded but response persistence fails
- WHEN the repository returns an error
- THEN no completion event is published and partial persistence is returned

### Requirement: Bounded completion
Inbound `LEAD` count MUST cap accepted turns at 4 for affiliates and 6 otherwise; post-cap input MUST be rejected before dispatch. At the cap, `CONTINUAR` MUST complete available fields; terminal safety actions prevail. Completion MUST transition only `PERFILANDO` to `CALIFICADO` and publish one post-durability `PerfilCompleto` event.

#### Scenario: Duplicate completion
- GIVEN a `CALIFICADO` lead or completed profile
- WHEN completion is requested again
- THEN no transition or `PerfilCompleto` event is repeated

### Requirement: Question helpers
`SiguienteMejorPregunta` MUST choose the highest-priority missing conversational field and MUST NOT return `VERIFICADO_BASE`. `PerfilEstaCompleto` MUST require `ingreso_hogar`, `recursos_propios`, and `zona_deseada`, without redefining motor `CamposCriticos`.

#### Scenario: Verified field
- GIVEN the first-priority field is verified and the next is missing
- WHEN the next question is requested
- THEN the helper selects the missing non-verified field

### Requirement: Provider isolation
The use case MUST surface provider errors, including Qwen's typed unsupported-audio error, without retry or fallback. It MUST use injected repository, clock, ID generator, provider, and bus only; it MUST NOT use a database, real provider, wall clock, ADK graph, or foreign-lead data.

#### Scenario: Provider error
- GIVEN a typed audio-capability error
- WHEN valid audio is processed
- THEN the error returns with zero writes and exactly one provider call

### Requirement: Consent denial precedes normal turn mutation

When a valid provider result has `Accion=CONSENTIMIENTO_NO`, `ProcesarMensaje` MUST persist the inbound denial evidence and invoke `SaludarLead.RechazarConsentimiento` before field normalization, `aplicarCampos`, capacity calculation, intention assignment, or normal response/completion handling. It MUST preserve the pre-denial profile, capacity, and intention, append exactly the refusal farewell, and publish no event. It SHALL retain existing validation and provider-dispatch behavior before this classified action is available.

#### Scenario: Denial with extracted fields
- GIVEN a `PERFILANDO` lead and a result with `CONSENTIMIENTO_NO` plus extracted fields
- WHEN the text turn is processed
- THEN no extracted field, capacity, or intention is saved
- AND only the inbound denial, farewell, final route, and final state are persisted

#### Scenario: Denial does not complete a profile
- GIVEN a lead whose profile would otherwise be complete
- WHEN the result is `CONSENTIMIENTO_NO`
- THEN no `PerfilCompleto` event is published
- AND the final state is `DESPEDIDO`

## Non-Goals

This capability SHALL NOT add ADK/HTTP wiring, frontend, migrations, recommendations, routing, ficha, nutrition, provider retries, or changes to ports, DTOs, repositories, domain vocabulary, motor criteria, or Contract v1.1. Provider-returned transcription requires a separate Contract-labelled change.