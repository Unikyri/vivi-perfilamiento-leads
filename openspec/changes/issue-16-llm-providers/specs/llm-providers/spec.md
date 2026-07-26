# LLM Providers Specification

## Purpose

Defines Contract §7 Gemini/Qwen behavior, resilience, and safety.

## Requirements

### Requirement: Strict structured output

Providers MUST return only JSON `SalidaTurno` with exact `campos_extraidos`, `intencion`, `respuesta`, and `accion`, rejecting unknown members at every level. Each extraction MUST use Contract members, a §2.1 field, `DECLARADO` or `INFERIDO`, and valid confidence. `accion` MUST be one of the seven §7 actions; `respuesta` MUST use only `NUMEROS_DEL_MOTOR` monetary amounts.

#### Scenario: Valid Contract output
- GIVEN Contract JSON
- WHEN parsed
- THEN its `SalidaTurno` is returned

#### Scenario: Unauthorized output
- GIVEN an unknown member, field, source, or action
- WHEN parsed
- THEN typed invalid-output is returned without fallback

### Requirement: Prompt and input invariants

The immutable housing-only instruction MUST be separate from user data. User text MUST be data; only session context MAY be included; `VERIFICADO_BASE` MUST NOT be requested; money MUST use injected `NUMEROS_DEL_MOTOR`.

#### Scenario: Prompt injection input
- GIVEN user text requests a role or prompt change
- WHEN the turn is prepared
- THEN the instruction is unchanged and no secret or other-lead data is included

### Requirement: One call and format recovery

A normal turn MUST make one provider request. Only malformed/non-JSON output is format-retryable and MUST receive one same-provider retry. A second malformed output MUST be typed, fallback-eligible format failure; schema-valid but contract-invalid output MUST NOT fallback.

#### Scenario: Malformed response recovery
- GIVEN a malformed first response and valid retry
- WHEN the turn completes
- THEN exactly two requests produce the valid output

### Requirement: Gemini and Qwen capabilities

Gemini MUST make one structured text or multimodal-audio request. Qwen MUST make one structured text request to its configured endpoint. Qwen `ProcesarAudio` MUST return typed `unsupported-audio` before any request and MUST NOT reroute to Gemini, including when Gemini is fallback.

#### Scenario: Qwen primary audio
- GIVEN Qwen is primary and valid audio arrives
- WHEN `ProcesarAudio` is called
- THEN zero Qwen or Gemini requests occur and `unsupported-audio` is returned

### Requirement: Transient failure and circuit behavior

Each external request MUST have an 8-second hard timeout. Timeout, 429, transport, and retryable-server failures MUST be typed transient failures and MAY use configured fallback. Per-provider rate limiting and a breaker MUST open for 60 seconds after three consecutive transient/429 failures, bypass the provider, and expose state through health integration. Exhaustion MUST return a typed error, never UI text.

#### Scenario: Gemini rate-limit fallback
- GIVEN Gemini returns 429 for three turns and Qwen is configured
- WHEN the next eligible turn runs
- THEN Gemini is bypassed or falls through to Qwen without a courtesy message

### Requirement: Configuration-only factory

The factory MUST use only existing `LLM_PROVIDER`, `LLM_FALLBACK`, Gemini key, Qwen key, and Qwen base URL; it MUST NOT alter ports, config signatures, or environment variables. After restart, `LLM_PROVIDER=qwen` MUST select Qwen without code changes; unsupported names or credentials MUST return typed configuration errors.

#### Scenario: Qwen configuration switch
- GIVEN `LLM_PROVIDER=qwen` and a valid key
- WHEN the application restarts
- THEN text calls and health use Qwen

### Requirement: Safe validation and tests

Guardrail failures MUST be typed and make zero provider requests. Tests MUST use fakes/stubs; CI MUST NOT call a real LLM. Logs MUST NOT contain keys, message/audio content, cédulas, or raw responses.

#### Scenario: Guardrail rejection
- GIVEN an input guardrail rejects a message
- WHEN invoked
- THEN zero network requests and sensitive logs occur

### Requirement: Scope boundary and monitoring

This change MUST preserve `LLMProvider` and `Config` signatures. Courtesy UI belongs to a future conversation use case, which consumes typed errors. CI alignment is out of scope: main `43907ecd8aafa71f36344f5ee592191396e85ff6` run `30168602827` succeeded and MUST be monitored only.

#### Scenario: Both providers unavailable
- GIVEN primary and fallback return transient errors
- WHEN the layer exhausts them
- THEN it returns a typed error and leaves UI messaging to the conversation use case
