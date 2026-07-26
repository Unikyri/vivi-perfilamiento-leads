# Proposal: Security Guardrails and Resilience Composition (Issue #17)

## Intent

NFR-S-01 forbids taking Vivi out of her role. Layer 1 (hardened prompt) shipped with Issue #16, but no production input classifier or output leakage validator exists, so any adversarial message still reaches a paid provider call and any hallucinated peso figure still reaches the user. This change adds NFR-S-01 layers 2–3 and NFR-D-03 observability as stackable decorators over the same `usecase.LLMProvider` port (NFR-M-03.3), reusing Issue #16 resilience instead of duplicating it.

## Scope

### In Scope
- `Guardarrailes` decorator: deterministic input classifier (jailbreak, prompt extraction, role change, third-party data, out-of-domain/illicit) answering with a template `SalidaTurno{Accion: FUERA_DE_DOMINIO}` and zero provider calls.
- Output validation: suppress system-prompt/skill/internal-file leakage, foreign `lead_id` data, and peso amounts absent from `EntradaTurno.NumerosDelMotor` (Contract v1.1 §7).
- `Metricas` decorator: JSON structured logs with `lead_id`, event, latency, provider, outcome kind — never message text, audio, cédulas, prompts, or credentials (NFR-S-02, NFR-D-03).
- Production composition in `factory.go`: `Metricas(Guardarrailes(FallbackProvider(primary, fallback)))`, preserving live breaker health on `/salud`.
- `tests/adversarios.json` with exactly the 15 issue-defined cases, plus deterministic fake-provider tests (containment, zero-call proof, figure allow/deny, prompt-leak, breaker, benign housing text).

### Out of Scope (non-goals)
- No `ProcesarMensaje` usecase, HTTP turn endpoint, agent graph, or session persistence (Issue #19).
- No Contract/Wiki edits, new dependencies, provider SDKs, or config variables.
- No second breaker, limiter, timeout, or fallback path; no changes to `ParseSalida` strictness.
- No live LLM calls in tests or CI.

## Capabilities

### New Capabilities
- `llm-guardrails`: input containment, output leakage/figure validation, privacy-safe metrics, and decorator composition order.

### Modified Capabilities
- `llm-providers` (delta lands in this change folder; Issue #16 spec is not yet archived into `openspec/specs/`): factory composition now returns a decorated provider, and breaker health MUST stay observable through decorators.

## Approach

| Decision | Rationale |
|---|---|
| Decorators implement `LLMProvider` | Keeps security out of business logic; `Nombre()` and both turn methods delegate. |
| Input guardrail returns a template, not an error | Issue DoD and NFR-S-01.2 require a redirection reply; `KindGuardrail` (error-based `WithGuardrail` seam) would deny the user a response. Seam left untouched. |
| Reuse Issue #16 `Breaker` | `FallbackProvider` already enforces 3 eligible failures → open 60 s, per-provider token buckets, 8 s attempts, 10 s logical timeout, typed errors. The issue's literal outer `CircuitBreaker` snippet would duplicate state, double-count failures, and add a conflicting 8 s timeout. Documented deviation: the breaker layer stays inside the composite. |
| Health via small interface | `CircuitBreakerHealth` currently type-asserts `*FallbackProvider`; wrapping would silently report `CERRADO`. Assert `interface{ CircuitBreakerState() BreakerState }` and delegate through decorators. |
| Validation is local | Guardrail outcomes never trigger fallback, retry, or a second decision call. |

## Affected Areas

| Area | Impact | Description |
|---|---|---|
| `internal/infrastructure/llm/guardarrailes.go` | New | Input classifier + output validator decorator. |
| `internal/infrastructure/llm/metricas.go` | New | Latency/provider logger; may also satisfy the existing `Metrics` seam. |
| `internal/infrastructure/llm/factory.go` | Modified | Compose decorators around `NewFallbackProvider`. |
| `internal/infrastructure/llm/health.go` | Modified | Interface-based live breaker state through decorators. |
| `internal/infrastructure/llm/{guardarrailes,metricas,breaker}_test.go` | New | Deterministic fake-provider tests. |
| `tests/adversarios.json` | New | 15-case fixture. |
| `internal/infrastructure/llm/breaker.go`, `fallback.go`, `parser.go`, `prompt.go` | Unchanged | Reused as-is. |

## Risks

| Risk | Likelihood | Mitigation |
|---|---|---|
| Regex over-blocking legitimate housing Spanish | Med | Narrow anchored patterns; benign near-miss regression cases. |
| Decorator wrapping hides breaker state on `/salud` | High if unhandled | Interface-based health delegation with a direct runtime test. |
| Deviating from the issue's literal breaker snippet is read as missing scope | Med | Deviation stated here; DoD (3 failures / 60 s / visible state) proven against the existing breaker. |
| Foreign-lead isolation only partially provable at this seam | Med | Check explicit foreign `lead_id` tokens; full tool/session binding is Issue #19. |
| Review budget | Med | Forecast: ~230–270 authored implementation lines + ~150–180 test/fixture lines ≈ 380–450 total. Single PR to `main`; if the count crosses 400, request `size:exception` rather than splitting a security layer. |

## Rollback Plan

All behavior is additive and file-local. Revert by (1) restoring `factory.go` to `return NewFallbackProvider(...)` undecorated and `health.go` to the `*FallbackProvider` assertion — a two-hunk revert that restores exact Issue #16 behavior — or (2) `git revert` of the single PR merge commit. New files are unreferenced after step 1; no data, schema, config, or Contract migration is involved.

## Dependencies

- Issue #16 (`internal/infrastructure/llm` foundation) — merged.
- Blocks Issue #19 (`ProcesarMensaje`).

## Success Criteria

- [ ] 15/15 adversarial cases contained; none yields a normal model answer.
- [ ] Input-guardrail rejections invoke zero provider calls (proven by fake-provider counter).
- [ ] A figure absent from `NumerosDelMotor` is suppressed; a present one passes.
- [ ] Responses never leak system prompt, skill/internal file names, or foreign-lead data.
- [ ] Breaker still opens after 3 eligible failures for 60 s and `/salud` reports `ABIERTO` through the decorated provider.
- [ ] Logs contain no message content, audio, cédulas, prompts, or credentials.
- [ ] `go build ./...`, `go vet ./...`, `go test ./...`, `go mod verify`, and CI (backend, frontend, secrets) pass.

## Proposal question round

Non-interactive run; assumptions above need user review if any is wrong:
1. Business rule: a blocked message should still get a friendly redirection reply (not an error/silence) — assumed yes per issue DoD.
2. Out-of-domain detection scope: pattern/category-based only, accepting that novel out-of-domain phrasing may pass to layer 1 (hardened prompt) — assumed acceptable for Fase 0.
3. Composition: keeping the breaker inside `FallbackProvider` rather than adding the issue's outer `CircuitBreaker` — assumed acceptable since observable DoD is unchanged.
4. Tradeoff preference: prefer occasional false negatives over blocking real housing conversations — assumed yes (demo credibility).
