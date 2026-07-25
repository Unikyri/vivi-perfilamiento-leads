# Proposal: LLM Provider Adapters and Fallback (Issue #16)

## Intent

`internal/infrastructure/llm` holds only `doc.go` and the composition root carries an Issue #16 placeholder, so no turn can reach a real model. Vivi needs configurable Gemini/Qwen providers returning one strictly validated Contract §7 `SalidaTurno` per turn, with resilient fallback and no secret or message-content leakage.

## Scope

### In Scope
- Provider-neutral wire DTOs, strict Contract §7 `SalidaTurno` parser/validator, typed error taxonomy.
- Gemini adapter (text + inline audio); Qwen adapter (text, typed unsupported-audio error).
- Config-only factory for primary/fallback selection using existing `Config` fields.
- Typed seams: 8s timeout, 429/rate limit, fallback cascade, circuit breaker (3 failures / 60s), guardrail and metrics decorators.
- Focused `httptest`/`RoundTripper` tests plus composition-root wiring and provider health identity.

### Out of Scope
- `usecase.LLMProvider` or `Config` signature/env-var changes.
- Conversation use case, HTTP message handler, ADK agent graph, prompt tuning.
- Motor/domain behavior, frontend, data pipeline, Contract/Wiki edits, provider SDK dependencies.

## Capabilities

### New Capabilities
- `llm-providers`: provider port implementations, structured-output validation, configured selection, resilience and guardrail behavior.

### Modified Capabilities
- None.

## Approach

`net/http` adapters behind an injectable transport (no SDK) keep wire types inside infrastructure, permit exact request-count and timeout assertions, and avoid unpinned dependencies. The parser rejects unknown fields, non-Contract profile keys, sources other than `DECLARADO`/`INFERIDO`, and invalid actions. Exactly one normal provider call per turn; a single invalid-JSON format retry is exceptional recovery only, then fallback. Validation and capability errors never silently trigger fallback. CI uses fake transports only — no real LLM calls, secrets, or logged content.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `internal/infrastructure/llm/` | New | Wire models, parser, adapters, factory, decorators, tests |
| `cmd/servidor/main.go` | Modified | Replace placeholder wiring; accurate active-provider health |
| `internal/usecase/puertos.go`, `config/config.go` | Unchanged | Read-only compatibility boundary |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| go.mod Go 1.25 vs CI Go 1.24 blocks PRs | High | Decide alignment before slice 1; treat as prerequisite, not provider work |
| Provider REST schema drift | Med | Pin documented endpoint/shape in adapter tests |
| Fallback masks validation/capability errors | Med | Typed error classes; fallback only for transient provider failures |
| Guardrail/secret leakage in logs | Med | Redaction seam; assert no key/content logging |
| No end-to-end turn exists yet | High | Verify at provider/seam level; defer full-turn proof |

## Rollback Plan

Each of the three chained PR slices is additive and independently revertible by reverting its merge commit. Runtime behavior is unchanged until slice 3 wires the factory; reverting slice 3 restores the placeholder composition root. No schema, data, or configuration migration is involved.

## Dependencies

- Issue #11 port satisfied; no signature changes required.
- Go toolchain alignment decision (go.mod 1.25 vs CI 1.24).

## Success Criteria

- [ ] `LLMProvider` and `Config` signatures unchanged; build and tests pass.
- [ ] Parser rejects malformed or unauthorized output with typed errors.
- [ ] `LLM_PROVIDER` switch and Gemini-429 → Qwen fallback proven by tests.
- [ ] Qwen `ProcesarAudio` returns typed unsupported-audio error with zero HTTP calls.
- [ ] One provider call per valid turn; timeout, circuit-breaker, guardrail seams tested.
- [ ] Three chained PRs, each ≤400 authored changed lines, no real LLM/secrets in CI.

## Proposal question round

Blocked from asking interactively; these need user confirmation before spec/design:

1. When both providers fail, is the courtesy fallback message owned here or by the future conversation use case?
2. Should the Go 1.25/1.24 alignment ship as a separate prerequisite PR outside this change?
3. Is audio-to-Gemini fallback acceptable when Qwen is the configured primary, or should audio simply fail?
