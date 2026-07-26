# api-http Specification

## Requirements

### Requirement: JSON, errors, and privacy
API responses MUST use `application/json; charset=utf-8`; `snake_case`; UTC RFC 3339; integer COP; and UPPER_SNAKE enums. Failures MUST only be `{error:{codigo,mensaje,detalles}}`. Codes SHALL map: `VALIDACION` 400; `LEAD_NO_ENCONTRADO`/`FICHA_NO_DISPONIBLE` 404; `TRANSICION_INVALIDA` 409; `AUDIO_INVALIDO` 422; `LIMITE_TASA` 429; `PROVEEDOR_LLM_NO_DISPONIBLE` 503; `ERROR_INTERNO` 500. Unknown failures MUST be generic `ERROR_INTERNO`. Responses and logs MUST NOT expose cédulas, phones, message/audio content, providers, SQL, or stacks.

#### Scenario: Safe failure
- GIVEN invalid or internal failure
- WHEN endpoint handles it
- THEN it returns the mapped error envelope
- AND exposes no sensitive value.

### Requirement: Lead resources and turns
`POST /api/leads` MUST accept Contract fields or seed `ana|carlos|luisa` (overriding form fields), create `PERFILANDO`, initiate greeting, and return 201 `{lead_id,estado,afiliado_detectado}`. GET resources MUST return Contract shapes; ficha MUST distinguish absent lead/ficha. Message POST MUST accept text ≤2,000 characters or audio ≤60 seconds, ≤2 MiB decoded, MIME `audio/webm|audio/ogg|audio/mpeg`; audio MUST neither persist nor echo. It MUST return 202 `{mensaje_id,recibido_en,turno_en_proceso:true}` before background processing, keep the flag true during it, then clear it. Clients MAY poll every 1.5 seconds. MUST use plain use cases; no ADK.

#### Scenario: Turn acceptance
- GIVEN an existing lead and valid text
- WHEN its message is posted
- THEN it receives 202
- AND polling later observes completion.

#### Scenario: Rejected or absent resource
- GIVEN invalid audio or a lead without ficha
- WHEN its endpoint is requested
- THEN audio is 422 `AUDIO_INVALIDO` and ficha 404 `FICHA_NO_DISPONIBLE`.
