# Exploration: Security Guardrails and Resilience Decorators

## Context
Issue #17 implements the NFR-S-01 defense-in-depth layers that remain after Issue #16: a deterministic input/output guardrail, adversarial regression battery, and privacy-safe metrics. The existing Issue #16 `FallbackProvider` already provides the required per-provider 8-second attempts, token bucket, three eligible failures/60-second breaker, fallback, typed errors, and live `/salud` state.

## Findings
- `usecase.LLMProvider` is stable and has `GenerarTurno`, `ProcesarAudio`, and `Nombre`; decorators can preserve Clean Architecture by implementing this port in `internal/infrastructure/llm`.
- `FallbackProvider` has extension points for `Guardrail` and `Metrics`, but no production guardrail is supplied and the metric interface captures no duration.
- The issue explicitly scopes the concrete LLM decorators, a 15-case fixture, and tests. It does not require the missing HTTP/`ProcesarMensaje` orchestration, so no Contract/API changes are needed.
- The input guardrail must reject jailbreak, prompt/config extraction, third-party data requests, explicit/illegal content, and out-of-domain tasks without invoking either provider. Its return is the fixed safe `SalidaTurno`, not an error, because the issue requires a template response.
- Output validation must suppress internal prompt/skill/config leakage and monetary claims not represented by `EntradaTurno.NumerosDelMotor`; no raw messages, cédulas, prompts, or responses may be logged.

## Decision
Implement self-contained `Guardarrailes` and `Metricas` decorators. Reuse (rather than duplicate) the existing circuit-breaker behavior in `FallbackProvider`; expose a `ConCircuitBreaker` compatibility decorator only if necessary to satisfy the intended stack without a second breaker state. Compose the factory as metrics over guardrails over the fallback resilience composite. Add deterministic fake-provider tests and `tests/adversarios.json` (15 cases). No live LLM calls, SDKs, config keys, Contract edits, or user-facing HTTP endpoints are introduced.

## Risks and mitigations
- Regex false positives: test known in-domain housing text and only classify clear patterns/categories.
- Cross-lead data cannot be inferred from arbitrary output without session data; reject explicit foreign `lead_id` values relative to `EntradaTurno.LeadID` and known leak tokens. The future conversation use case owns repository/tool session binding.
- Preserve Issue #16 fallback eligibility: guardrail outcomes must never cause fallback; transient breaker behavior remains unchanged.
- Keep one normal provider invocation per accepted turn; validations are local and must not retry/regenerate.

## Validation plan
Run focused LLM package tests (guardrails, metrics, resilience) and repository-wide `go test ./...`, `go build ./...`, `go vet ./...`, `go mod verify`, and `git diff --check`. CI must pass backend, frontend, and secrets checks before merge.
