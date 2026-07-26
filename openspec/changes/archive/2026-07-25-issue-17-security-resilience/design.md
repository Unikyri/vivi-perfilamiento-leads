# Design: Issue #17 Security and Resilience Decorators

## Technical Approach

Implement security and observability as `usecase.LLMProvider` decorators, preserving Contract v1.1 §7 and Issue #16 provider behavior. Production composition is:

```text
Metricas → Guardarrailes → FallbackProvider
                            ├─ primary: limiter → breaker → provider
                            └─ fallback: limiter → breaker → provider
```

`Guardarrailes` performs deterministic local input classification and output validation. `Metricas` uses only Go standard-library JSON logging. Neither decorator makes an LLM call.

## Architecture Decisions

| Option | Tradeoff | Decision and rationale |
|---|---|---|
| Outer concrete guardrail vs existing `WithGuardrail` error hook | The hook runs before providers but can only return an error. | Use `ConGuardarrailes(LLMProvider)`. NFR-S-01.2 requires a safe reply, so blocked input returns a complete `FUERA_DE_DOMINIO` template with `nil` error and zero wrapped calls. Leave the hook compatible and unused. |
| Local safe replacement vs LLM regeneration | Regeneration conflicts with NFR-R-01’s single-call budget and can leak again. | On unsafe successful output, replace the entire output locally with a safe template; never retry or fallback. Reuse package-local `validResponse` for motor-money checks instead of duplicating its parser. |
| Existing per-provider breakers vs outer breaker | An outer breaker would double-count one logical turn, obscure which provider failed, duplicate timeout/state, and diverge from `/salud`. | Reuse `FallbackProvider` breakers: three consecutive eligible failures, 60-second open interval, half-open probe, and provider-specific fallback already satisfy NFR-E-03. `breaker.go` remains unchanged. |
| Concrete type health lookup vs capability interface | Wrappers break the current `*FallbackProvider` assertion. | Introduce package-local `interface{ CircuitBreakerState() BreakerState }`; both decorators delegate it, and `CircuitBreakerHealth` consumes it. |
| External telemetry vs standard library | A dependency expands risk and scope. | An `encoding/json` observer with an injected `io.Writer` records logical latency and per-provider attempt outcomes through the existing `Metrics` seam. It never logs messages, audio, responses, prompts, cédulas, credentials, or error text. |

## Data Flow and Semantics

```mermaid
sequenceDiagram
  participant C as Caller
  participant M as Metricas
  participant G as Guardarrailes
  participant F as FallbackProvider
  C->>M: GenerarTurno/ProcesarAudio
  M->>G: delegate once
  alt blocked input
    G-->>M: safe template, nil (0 provider calls)
  else accepted input
    G->>F: delegate once
    F->>F: existing primary/fallback policy
    alt provider error
      F-->>G: output,error unchanged
    else unsafe output
      G-->>M: safe template, nil (no retry)
    else safe output
      G-->>M: output,nil
    end
  end
  M-->>C: result; emit privacy-safe JSON
```

Both methods share pre/post checks. Input categories are jailbreak/prompt extraction, role change, third-party data, explicit/illegal content, and clear out-of-domain tasks. Third-party and generic blocks use distinct fixed templates; templates contain empty extracted fields, low intention, and `FUERA_DE_DOMINIO`. Output rejection covers prompt/skill/config/internal-file markers, any explicit `lead_id` serialization, and money rejected by `validResponse`. Arbitrary unlabeled cross-lead PII cannot be proven at this seam; Issue #19 must enforce session/tool identity.

A downstream error is returned unchanged and is not output-validated. New decorators add zero physical calls: accepted input invokes the wrapped provider once logically; Issue #16 may still perform its documented malformed retry or eligible fallback. Metrics log only typed error kind (or `unknown`), never `err.Error()`.

## File Changes

| File | Action | Description |
|---|---|---|
| `internal/infrastructure/llm/guardarrailes.go` | Create | Classifier, safe templates, post-validation, health delegation. |
| `internal/infrastructure/llm/metricas.go` | Create | `encoding/json` observer/decorator with injected writer/clock and health delegation. |
| `internal/infrastructure/llm/factory.go` | Modify | Build one observer, inject it into fallback attempts, then wrap guardrails and metrics in the stated order. |
| `internal/infrastructure/llm/health.go` | Modify | Read breaker state via capability interface. |
| `internal/infrastructure/llm/guardarrailes_test.go` | Create | Counting/scripted fake, fixture, input/output/audio/order/health semantics. |
| `internal/infrastructure/llm/metricas_test.go` | Create | Decode JSON lines and prove field/privacy/latency behavior. |
| `tests/adversarios.json` | Create | Exactly the 15 issue-defined adversarial cases. |
| `internal/infrastructure/llm/{fallback,breaker,parser}.go` | No change | Reused behavior and validation. |

## Testing Strategy

Fixture tests require 15/15 safe templates and zero provider calls; benign housing text calls once. Scripted fakes cover invalid/valid output, motor money, propagated errors, audio counts, factory order, and wrapped breaker health. Metrics tests parse buffered JSON with a fake clock and assert seeded secrets are absent. Run focused LLM tests, then `go test ./...`, `go build ./...`, `go vet ./...`, `go mod verify`, and `git diff --check`; CI makes no live LLM calls.

## Threat Matrix

N/A — no routing, shell, subprocess, VCS/PR automation, executable-file classification, or process-integration boundary.

## Rollout, Rollback, and Budget

No migration/configuration change; one PR targets `main`. Roll back by restoring `factory.go` to the undecorated fallback return and the old health assertion, or revert the PR; new files then become unreachable.

Forecast: **365–425 authored implementation/test/fixture changed lines** (center ≈395), excluding this design. **400-line budget risk: High at the upper bound.** Tasks must trim duplication and re-count before apply; if forecast exceeds 400, the single-PR strategy requires explicit `size:exception` rather than an implicit overrun.

## Open Questions

None blocking; the documented cross-lead limitation belongs to Issue #19.
