## Exploration: issue-16-llm-providers

### Current State

The isolated worktree is `/home/daikyri/Workspace/Hackathon-Colsubsidio/vivi-issue-16`, on `feat/issue-16-llm-providers` at `43907ecd8aafa71f36344f5ee592191396e85ff`, the required base revision. No Issue #16 OpenSpec folder or prior Engram artifact exists.

The application-owned port is already present in `internal/usecase/puertos.go` and MUST remain unchanged: `LLMProvider` exposes `GenerarTurno(context.Context, EntradaTurno) (SalidaTurno, error)`, `ProcesarAudio(context.Context, Audio, EntradaTurno) (SalidaTurno, error)`, and `Nombre() string`. The provider-neutral DTOs and Contract §7 JSON field names already exist. The existing `LLMFake` in `internal/usecase/fakes_test.go` returns static text/audio outputs but does not count calls or model provider failures.

`internal/infrastructure/llm` contains only `doc.go`; there is no Gemini/Qwen implementation, transport abstraction, parser, decorator, factory, or provider SDK dependency. `go.mod` contains PostgreSQL and data-pipeline dependencies only. `internal/infrastructure/config/config.go` already exposes `LLMProvider`, `GeminiAPIKey`, `QwenAPIKey`, `QwenBaseURL`, and `LLMFallback` with defaults (`gemini`, empty keys/base URL, `qwen`) and MUST be preserved. It does not currently validate provider names/keys or expose a timeout field. `cmd/servidor/main.go` only reports the configured provider in `/salud`; its composition root has a placeholder explicitly naming Issue #16 and does not construct an LLM provider or conversation use case.

The governing Contract v1.1 §7 requires a single structured `SalidaTurno` containing `campos_extraidos`, `intencion`, `respuesta`, and `accion`; valid actions are the seven constants already represented in `usecase`. Extracted fields must be recognized Contract profile keys, use only `DECLARADO` or `INFERIDO`, and monetary figures in the response must come from the injected `NUMEROS_DEL_MOTOR` context. Invalid provider JSON receives one format retry before fallback. §8 requires configuration-only provider selection and the Gemini/Qwen keys/base URL. NFR-R-01/R-03 require one combined extraction/decision/redaction invocation per logical turn and an 8-second hard provider timeout. NFR-E-03 requires provider rate limiting, a circuit breaker after three consecutive failures/429s for 60 seconds, and Gemini → Qwen → courtesy-message cascading. NFR-S-01 requires input/output guardrails and no secret or personal-message logging. NFR-M-03 and ADR-5 require an application-owned Strategy/Adapter port, a configuration Factory, and composable guardrail/circuit-breaker/metrics decorators; ADR-8 requires the single structured turn call. US-19 requires `LLM_PROVIDER=qwen` to switch all calls after restart and Gemini 429 failure to complete through configured Qwen fallback.

Audio is accepted at the HTTP contract as base64 plus MIME/duration, with contract limits of 60 seconds and 2 MB decoded. Contract §2.5 stores a lead audio message as text transcription plus `audio_original: true`. Architecture documentation identifies Gemini as multimodal/audio-capable and Qwen as an OpenAI-compatible fallback. The Qwen adapter therefore needs an explicit unsupported-audio result with no accidental text-only request; a composite may use a configured Gemini fallback when Qwen is primary, but capability errors must be distinguished from transient provider failures.

A repository-level validation mismatch is present: `go.mod` declares Go 1.25.0, while `.github/workflows/ci.yml` configures Go 1.24.0. This is not an Issue #16 behavior change, but it can make provider PR CI fail before tests run and should be resolved or explicitly accepted during proposal/design.

### Affected Areas
- `internal/usecase/puertos.go` — read-only compatibility boundary; preserve method signatures and DTOs exactly.
- `internal/infrastructure/config/config.go` — read-only configuration shape; factory must consume existing fields rather than inventing new environment variables.
- `internal/infrastructure/llm/` — new Gemini/Qwen adapters, provider-neutral wire models, structured-response parser/validator, typed provider and audio errors, factory, fallback, and resilience decorators.
- `cmd/servidor/main.go` — composition-root wiring and accurate active-provider health identity may be needed; currently it only echoes `cfg.LLMProvider` and has an Issue #16 placeholder.
- `internal/usecase/fakes_test.go` and new `internal/infrastructure/llm/*_test.go` — extend test doubles/call counters and add `httptest`/transport tests; real LLM calls MUST NOT enter CI.
- `openspec/changes/issue-16-llm-providers/exploration.md` and `state.yaml` — SDD artifacts only; no proposal/spec/design/tasks are created in this phase.
- Wiki documents 09, 10, 11, and 07 — read-only authority for NFRs, Contract §7/§8, ADR-5/ADR-8, and US-19.

### Approaches
1. **Official provider SDKs at the infrastructure boundary** — use the official Google GenAI/ADK-compatible client for Gemini structured output and audio, plus an OpenAI-compatible client or Qwen SDK for Qwen.
   - Pros: provider wire details and multimodal/schema handling are delegated to maintained clients; closest to ADR-5's Gemini `genai` wording.
   - Cons: no SDK is currently in `go.mod`; versions/API stability must be pinned; Qwen compatibility is still endpoint-dependent; SDK types must not cross `internal/usecase`; mocking and exact request-count tests require an additional seam.
   - Effort: High.

2. **Provider adapters over `net/http` with injectable transport** — keep a provider-neutral request/parser core, implement Gemini REST `generateContent` and Qwen OpenAI-compatible chat-completions wire calls behind small HTTP interfaces, and test with `httptest.Server` or a fake `RoundTripper`.
   - Pros: no new dependency, exact control over timeout, request count, schema validation, base URL, redaction, and Qwen compatibility; easy deterministic tests; preserves Clean Architecture and keeps SDK/wire types inside infrastructure.
   - Cons: manually maintains Gemini candidate/part and inline-audio extraction, structured-output request fields, and provider error mapping; must pin/document endpoint assumptions.
   - Effort: Medium-High.

3. **ADK-first graph integration** — add the ADK model registry and provider switching directly while building the agent graph and use cases.
   - Pros: follows the planned runtime architecture and could reuse ADK model abstractions.
   - Cons: `internal/adapters/agentes` is currently only a package marker and the use cases/HTTP conversation flow are not implemented; it couples Issue #16 to unrelated orchestration work, makes one-call/fallback tests harder, and is not independently reviewable.
   - Effort: High; reject for this change.

### Recommendation

Proceed with a provider-neutral core plus isolated provider adapters. Prefer `net/http` transport seams for this issue unless the proposal explicitly pins an official Gemini SDK version and demonstrates an injectable test boundary; the current module has no provider SDK and the REST approach best controls the exact one-call and fallback behavior. Keep `LLMProvider`, `Config`, and Contract JSON names unchanged. The parser should decode and validate `SalidaTurno`, reject unknown fields/actions/profile sources, preserve `domain.Intencion`, and map malformed output to a typed error. A Qwen `ProcesarAudio` call should return a typed unsupported-audio error without sending a text-only request; fallback policy should only try an alternate provider when its capability supports the input. Circuit-breaker/429/timeout failures should be classified separately from invalid output and unsupported audio, and provider names/keys must never be logged.

Use three chained, independently verifiable PR slices:
1. **Shared contract/transport slice** — provider wire DTOs, strict `SalidaTurno` parser/validator, typed errors, timeout/request context helpers, and unit tests for schema/source/action validation. Target ≤250 authored lines.
2. **Provider slice** — Gemini text+multimodal audio and Qwen text adapters, configured base URL/key handling, provider error mapping, audio capability behavior, and HTTP contract tests. Target ≤350 authored lines.
3. **Factory/resilience/wiring slice** — configured provider factory, fallback composition, circuit-breaker/metrics/guardrail seams required by the governing NFRs, call-count/429/timeout tests, and composition-root/health integration. Target ≤350 authored lines.

This is a chained plan rather than a single PR: aggregate work is likely above 400 authored lines and crosses external integrations, security, resilience, and composition-root boundaries. Each slice has a clear rollback boundary; slice 2 targets slice 1, and slice 3 targets slice 2. The exact issue checklist should be reconciled during proposal because it is not present in the local checkout.

### Risks

- The local repository does not contain the GitHub Issue #16 checklist; the proposal phase must verify any checklist items not represented in the user-provided scope.
- Provider SDK/API choice is unresolved. Adding a dependency requires a pinned version and CI compatibility review; direct REST requires stable endpoint/schema tests.
- `go.mod` requires Go 1.25 while CI selects Go 1.24; this may block all future provider PRs independently of implementation correctness.
- There is no conversation use case or HTTP message handler yet, so this phase can establish provider behavior and composition seams but cannot prove a complete end-to-end turn.
- Contract language allows one invalid-JSON format retry, while NFR-R-01 forbids multi-stage extraction/decision/redaction calls. Design must treat the retry as exceptional output recovery, not as a normal multi-call turn pipeline; the circuit's three-failure threshold should be tracked across requests rather than blindly issuing three calls in one turn.
- Fallback must not hide deterministic validation or unsupported-audio errors, and audio capability must be explicit when the configured primary is Qwen.
- Guardrails must prevent prompt/system leakage and untrusted monetary figures while not logging API keys, message content, cédulas, or raw provider responses.

### Ready for Proposal

Yes. The next phase should freeze the provider wire/API choice, typed error and retry/fallback classification, exact slice boundaries, and the Go-version CI discrepancy before implementation tasks are written. No production code has been changed.
