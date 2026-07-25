# Apply Progress: Issue 16 — Final Bounded Verification Slice 8

## Current Slice 8 Completion

- [x] 7.1 Bare dot-grouped monetary values are recognized only with nearby monetary context, normalized through the existing safe parser, and required to match `NUMEROS_DEL_MOTOR`.
- [x] 7.2 Focused parser tests cover authorized contextual amounts, unauthorized contextual amounts, and ordinary grouped count/date text.
- [x] 4.1 Final cross-slice verification completed through the full post-implementation validation; native verify admission remains blocked because the local native binary is unavailable.

### Slice 6 Changes

- `internal/infrastructure/llm/parser.go`: added complete-boundary matching for dot-grouped values and a bounded context vocabulary (`mensual`, `mensuales`, `cuota`, `ingreso`, `salario`, `presupuesto`, `precio`, `valor`); reused `parseCurrencyAmount` for safe normalization and overflow handling.
- `internal/infrastructure/llm/parser_test.go`: added authorized, unauthorized, and non-monetary grouped count/date regressions, including trailing sentence punctuation.

### Slice 6 Work Unit Evidence

| Evidence | Result |
|---|---|
| Focused test | `go test ./internal/infrastructure/llm -run 'TestParser' -count=1` — exit 0; parser and all existing currency/million regressions passed. |
| Runtime harness | Direct `ParseSalida` parser execution through focused Go tests; no provider/network/API, credential, message, or audio calls. |
| Rollback boundary | Revert only the Slice 6 additions in `internal/infrastructure/llm/parser.go` and `parser_test.go`; retain Slice 5 factory/resilience/parser behavior. |

### Slice 6 Physical Budget

- Historical estimate of authored additions: **76 physical lines** — `parser.go` 42 + `parser_test.go` 34. The uncommitted/untracked worktree prevents independent Git reconstruction of this slice.
- Bound: **400**; remaining **324** lines.
- Historical cumulative authored implementation estimate: **1,517** (prior estimate 1,441 + estimated Slice 6 additions 76).

### Slice 6 Validation

- `gofmt` / `gofmt -l` — passed.
- `go test ./internal/infrastructure/llm -run 'TestParser' -count=1` — passed.
- `go test ./... -count=1` — passed.
- `go test -race ./...` — passed.
- `go build ./...` — passed.
- `go vet ./...` — passed.
- `go mod verify` — all modules verified.
- `git diff --check` and untracked changed-file whitespace checks — passed.
- No dependencies, environment/config changes, verify report, review receipt, commit, push, PR, archive, or issue action was performed.

## Work Unit

- **Mode:** Standard (strict TDD disabled in `openspec/config.yaml`)
- **Delivery:** Chained feature-branch remediation slice, ≤400 physical added lines
- **Scope:** Preserve resilient single-provider factory composition and live breaker health; parse Colombian million currency expressions without restricting ordinary counts/dates. No ports, config signatures, dependencies, live APIs, commits, PRs, archive, or verify admission.
- **Boundary:** Native verify admission remains blocked; no verify report, receipt, or archive artifact exists.

## Completed Tasks

- [x] 1.1 Sanitized typed error taxonomy, injected `HTTPDoer`, cancellation, and 8-second request cap.
- [x] 1.2 Deterministic hardened prompt with delimited user data, sorted JSON maps, ordered history, retained verified-base fields and source, and an instruction not to request them.
- [x] 1.3 Strict Contract §7 parser with unknown-field rejection, EOF enforcement, field/source/confidence/intent/action validation, and motor-number response validation.
- [x] 1.4 Table-driven prompt, parser, and transport tests for injection boundaries, determinism, invalid output, redaction, timeout, and cancellation.
- [x] 2.1 Gemini REST adapter: default/overridable model endpoint, API-key header, structured JSON request, candidate-0 text extraction, inline audio, shared parser, and typed transport/parser mapping.
- [x] 2.2 Qwen REST adapter: configurable OpenAI-compatible endpoint, Bearer auth, ordered system/user messages, JSON response format, and pre-transport capability error for audio.
- [x] 2.3 Provider tests: exact endpoint/header/body assertions, response extraction, one normal call, exactly one malformed parser retry, 429/5xx typed errors, Gemini inline audio, Qwen audio zero-call behavior, and sanitized errors.
- [x] 3.1 Configuration-only factory preserves primary validation and treats an invalid or unavailable optional fallback as absent.
- [x] 3.2 Fallback routing, logical timeout, breaker, typed exhaustion, and resilience behavior remain intact.
- [x] 3.3 Guardrail and metrics seams, non-fallback validation/capability behavior, and no-sensitive-data tests remain intact.
- [x] 3.4 Server health wiring remains scoped to the safe factory identity and is included in the all-change budget.
- [x] 5.1 Exhausted malformed/non-JSON output is fallback-eligible after exactly one same-provider retry; `TestMalformedAfterRepairFallsBack` asserts primary 2 and fallback 1 calls.
- [x] 5.2 Bounded injected-clock `TokenBucket` gates each provider before invocation; `TestTokenBucketDenialUsesFallback` proves denial, fallback, and refill.
- [x] 5.3 The server retains one `NewFromConfig` composition and `/salud` reports its live breaker; `TestSaludReportsLiveBreakerState` proves `ABIERTO`.
- [x] 5.4 Currency-only numeric validation and Slice 1 accounting clarification are persisted.

## Remediation Changes

- `internal/infrastructure/llm/errors.go`: `KindMalformed` is eligible only as exhausted format-repair failure; valid JSON Contract violations remain `KindInvalidOutput` and non-fallbackable.
- `internal/infrastructure/llm/rate_limiter.go` and `fallback.go`: bounded injected-clock token buckets gate primary and fallback independently before provider invocation; local denial is sanitized and fallback-eligible.
- `internal/infrastructure/llm/health.go` and `cmd/servidor/main.go`: retain the constructed provider composition and expose its live primary breaker through `/salud`; nil/unconfigured remains safe.
- `internal/infrastructure/llm/parser.go`: only currency-marked amounts are checked against motor values; arbitrary count/date tokens are allowed.
- Focused tests cover malformed primary + valid fallback, limiter denial/refill, parser scope, and the actual health handler path.

## Work Unit Evidence

| Evidence | Result |
|---|---|
| Focused test | `go test ./internal/infrastructure/llm ./cmd/servidor -count=1` — exit 0; malformed fallback, limiter, parser, breaker, and server health checks passed. |
| Runtime harness | Fake Gemini transport, fake providers, injected fake clock, and `httptest` health handler; no real APIs, credentials, or message/audio logging. |
| Rollback boundary | Revert Slice 4 files and focused fallback/parser/main changes; retain Slice 1/2 and unrelated modules. |

## Physical Budget Accounting

- Slice 1 original source/test snapshot: **historical estimate of 295 authored physical lines**. Later `errors.go` additions are separate; the uncommitted/untracked worktree prevents independent Git reconstruction.
- Historical estimate for Slice 4 remediation additions: **247 authored physical lines** — `rate_limiter.go` 74 + `health.go` 23 + `cmd/servidor/main_test.go` 41 + resilience test 47 + parser test 11 + fallback 29 + main 13 + parser 8 + errors 1. The uncommitted/untracked worktree prevents independent Git reconstruction.
- Historical cumulative authored implementation estimate: 1,341 (prior estimate 1,094 + estimated Slice 4 additions 247).

## Validation

- `gofmt` — passed on all changed Go files.
- `go test ./internal/infrastructure/llm ./cmd/servidor -count=1` — passed.
- `go test ./... -count=1` — passed.
- `go test -race ./...` — passed.
- `go build ./...` — passed.
- `go vet ./...` — passed.
- `go mod verify` — all modules verified.
- `git diff --check` and untracked Go `git diff --no-index --check` — passed.
- No-live-network/secret audit — passed; tests use fake transports/providers and no credential-like provider endpoint literals or direct network calls.

## Apply Completion State

- [x] 4.1 Final cross-slice verification completed through the full post-implementation validation; native verify admission remains blocked because the local native binary is unavailable.

## Prior Verification Attempt (Superseded)

The earlier Slice 3 verification evidence below remains historical context only. It does not complete current task 4.1 because this bounded remediation changes provider fallback, rate limiting, health reporting, and parser behavior.

## Historical Verification Evidence

- **Scope:** Final cross-slice checks only; no production code changed. No `sdd-verify`, archive, review-receipt, commit, push, or PR actions were performed.
- **Focused package test:** `go test ./internal/infrastructure/llm -count=1` — exit 0.
- **Formatting:** `gofmt -l` over all Issue 16 changed Go files — exit 0; no files reported.
- **Full validation:**
  - `go test ./... -count=1` — exit 0; all packages passed.
  - `go test -race ./...` — exit 0; all packages passed.
  - `go build ./...` — exit 0.
  - `go vet ./...` — exit 0.
  - `go mod verify` — exit 0; all modules verified.
  - `git diff --check` — exit 0.
  - Untracked Issue 16 Go-file whitespace audit with `git diff --no-index --check` — pass.
- **Compatibility audit:** `internal/usecase/puertos.go` and `internal/infrastructure/config/config.go` have no working-tree diff. `LLMProvider` remains `GenerarTurno(context.Context, EntradaTurno) (SalidaTurno, error)`, `ProcesarAudio(context.Context, Audio, EntradaTurno) (SalidaTurno, error)`, and `Nombre() string`; `Config` fields and `Cargar() (Config, error)` are unchanged.
- **Scenario coverage:** 9/9 specification scenarios have runtime test coverage: valid/unauthorized Contract output (`TestParserContractValidation`), prompt injection (`TestPromptInvariantsAndDeterminism`), malformed recovery (`TestGeminiMalformedRetry`, `TestQwenMalformedRetryAndStatuses`), Qwen primary audio (`TestQwenAudioHasNoTransportOrReroute`, `TestQwenAudioNeverReroutes`), Gemini rate-limit fallback and exhaustion (`TestFallbackEligibilityAndExhaustion`, `TestBreakerOpenResetAndTimeout`), Qwen configuration switch (`TestFactorySelectionAndSafeIdentity`), guardrail rejection (`TestGuardrailStopsCallsWithoutLeakingInput`), and both providers unavailable (`TestFallbackEligibilityAndExhaustion`).
- **Qwen audio audit:** capability failure occurs before transport; tests assert zero Qwen calls and zero Gemini fallback calls.
- **API and secret audit:** tests use injected fake transports/providers; no direct live-network test calls, real credentials, message/audio logging, or credential-like literals are present. Provider default endpoints remain configuration-only production defaults and are not invoked by tests.
- **Conservative slice budgets:** Slice 1 **295/300**, Slice 2 **399/400**, Slice 3 **400/400**; cumulative authored implementation budget **1094**.
- **Rollback boundary:** Revert only the Issue 16 implementation slices and their tests; this final verification adds no production behavior.

**Historical remaining tasks at that time:** None in the prior Slice 3 apply. The current task 4.1 is complete after full validation; native admission remains blocked because the local binary is unavailable.

## Current Slice 5 Completion

- [x] 6.1 `NewFromConfig` now retains a `FallbackProvider` for valid single-provider and optional-invalid-fallback compositions, preserving genuine limiter/breaker state; direct runtime tests prove three failures yield live `ABIERTO` health and the token bucket denies a second provider call.
- [x] 6.2 Monetary parsing now normalizes `2 millones`, `$2 millones`, and `COP 2 millones` to 2,000,000, rejects unauthorized monetary figures, and leaves ordinary counts/dates unrestricted; precise parser tests cover all cases.

## Slice 5 Changes

- `internal/infrastructure/llm/factory.go`: always returns the resilient composition after primary validation, with a nil optional fallback when absent or invalid.
- `internal/infrastructure/llm/parser.go`: applies the million multiplier and safe separator normalization only when a currency unit is present.
- `internal/infrastructure/llm/parser_test.go`: tests bare, `$`, and `COP` million expressions, unauthorized amounts, and ordinary count/date text.
- `internal/infrastructure/llm/resilience_test.go`: tests factory-backed no-fallback live breaker health and genuine limiter denial.

## Slice 5 Work Unit Evidence

| Evidence | Result |
|---|---|
| Focused test | `go test ./internal/infrastructure/llm ./cmd/servidor -count=1` — exit 0. |
| Runtime harness | `NewFromConfig` with fake provider substitution, injected token-bucket clock, and direct parser execution; no live APIs, credentials, or message/audio logging. |
| Rollback boundary | Revert only the four Slice 5 files above; retain Slice 4 and earlier implementation. |

## Slice 5 Physical Budget

- Historical estimate for Slice 5 additions: **100 physical lines** — factory 1, parser 43, parser tests 17, resilience tests 39. The uncommitted/untracked worktree prevents independent Git reconstruction.
- Historical cumulative authored implementation estimate: **1,441** (prior estimate 1,341 + estimated Slice 5 additions 100).

## Slice 5 Validation

- `gofmt` / `gofmt -l` — passed on all Slice 5 Go files.
- `go test ./internal/infrastructure/llm ./cmd/servidor -count=1` — passed.
- `go test ./... -count=1` — passed.
- `go test -race ./...` — passed.
- `go build ./...` — passed.
- `go vet ./...` — passed.
- `go mod verify` — all modules verified.
- `git diff --check` and changed-file snapshot checks — passed.
- No live API, credential, message, or audio calls are used by the runtime tests. A repository-wide `/dev/null` probe reports an unrelated pre-existing EOF warning in `internal/pipeline/proyectos_test.go`; that out-of-scope file was not modified.

## Current Slice 7 Completion

- [x] 8.1 Contextual Colombian bare amounts now use lexical adjacency and sentence-aware boundaries; grouped zero-decimal and ungrouped whole amounts match motor values, while non-integral fractional pesos are rejected exactly without float rounding. Existing `$`, `COP`, `pesos`, and `millones` paths remain supported.
- [x] 8.2 Default provider rate limiting now uses named bounded chat-adapter constants: burst 3 and refill 0.5 tokens/second, with no new configuration or environment variables.
- [x] 8.3 Ineligible breaker failures reset the consecutive eligible-failure sequence; success continues to reset it, and opening requires exactly three consecutive eligible failures.
- [x] 4.1 Final cross-slice verification completed through the full post-implementation validation; native verify admission remains blocked because the local native binary is unavailable.

### Slice 7 Changes

- `internal/infrastructure/llm/parser.go`: replaced the 64-character context window with lexical adjacency (maximum three intervening words), no numeric token in the gap, and sentence/clause boundaries; added exact grouped/decimal validation and integer-only normalization with overflow protection.
- `internal/infrastructure/llm/parser_test.go`: added direct grouped-decimal, ungrouped-whole, unauthorized, date/count, sentence-boundary, and distant-linkage regressions.
- `internal/infrastructure/llm/rate_limiter.go` and `fallback.go`: named and applied intentional bounded default token-bucket values.
- `internal/infrastructure/llm/breaker.go` and `resilience_test.go`: reset ineligible sequences and prove the exact three-consecutive-eligible-failure rule.
- `openspec/changes/issue-16-llm-providers/design.md`: documented the Slice 7 parsing, limiter, and breaker decisions.

### Slice 7 Work Unit Evidence

| Evidence | Result |
|---|---|
| Focused test | `go test ./internal/infrastructure/llm ./cmd/servidor -count=1` — exit 0; parser, limiter, fallback, breaker, and health tests passed. |
| Runtime harness | Direct parser execution, fake providers, injected fake clock, and fake HTTP transports; no live provider/API, credentials, message, or audio calls. |
| Rollback boundary | Revert Slice 7 changes in `parser.go`, `parser_test.go`, `rate_limiter.go`, `fallback.go`, `breaker.go`, and `resilience_test.go`, plus the matching Slice 7 documentation additions; retain Slices 1–6 and unrelated modules. |

### Slice 7 Physical Budget

- Historical estimate for Slice 7 net physical implementation additions: **247 lines** — `parser.go` +141, `parser_test.go` +44, `rate_limiter.go` +8, `fallback.go` +0, `breaker.go` +4, `resilience_test.go` +50. The uncommitted/untracked worktree prevents independent Git reconstruction.
- Historical cumulative authored implementation estimate: **1,764** (prior estimate 1,517 + estimated Slice 7 additions 247).

### Slice 7 Validation

- `gofmt` / `gofmt -l` — passed.
- `go test ./internal/infrastructure/llm -run 'TestParser|TestBreaker|TestDefaultProviderRateLimiter|TestFallback|TestTokenBucket' -count=1` — passed.
- `go test ./internal/infrastructure/llm ./cmd/servidor -count=1` — passed.
- `go test ./... -count=1` — passed.
- `go test -race ./...` — passed.
- `go build ./...` — passed.
- `go vet ./...` — passed.
- `go mod verify` — all modules verified.
- `gofmt -l`, `git diff --check`, and no-live-network/credential audits — passed.
- No dependencies, environment/config changes, verify report, review artifacts, commit, push, PR, archive, or issue action was performed.


## Current Slice 8 Completion

- [x] 4.1 Final cross-slice verification completed through the full post-implementation validation; native verify admission remains blocked because the local native binary is unavailable.

### Slice 8 Changes

- `internal/infrastructure/llm/parser_test.go`: materialized the deterministic signed `int64` overflow probe; an overflowing monetary amount returns invalid output and cannot wrap to an authorized `MinInt64` motor value.
- `internal/infrastructure/llm/resilience_test.go`: materialized the deterministic composite-breaker probe; after three eligible primary failures, the fourth invocation bypasses primary and succeeds through the fake fallback.
- No externally exposed health string or production behavior was changed.

### Slice 8 Work Unit Evidence

| Evidence | Result |
|---|---|
| Focused test | `go test ./internal/infrastructure/llm -run 'TestParserInt64OverflowIsInvalidWithoutWrap|TestCompositeBreakerBypassesPrimaryAfterThreeEligibleFailures' -count=1` — exit 0. |
| Runtime harness | In-process parser plus fake providers and fake rate limiters; no network, credentials, provider, message, or audio calls. |
| Rollback boundary | Revert only the two appended test blocks in `parser_test.go` and `resilience_test.go`; retain all prior implementation slices. |

### Slice 8 Accounting Reconciliation

- Current-tree physical line counts after this slice: `parser_test.go` **141**, `resilience_test.go` **289**, `parser.go` **317**, `breaker.go` **91**, `fallback.go` **158**, `rate_limiter.go` **81**, `cmd/servidor/main.go` **102**.
- Direct Slice 8 edit content is **32 physical lines** (9 + 23), below the **400-line** bound. This is an edit measurement, not a reconstructed Git diff.
- Historical Slice 1–7 figures in this artifact are estimates/snapshots only. The Issue 16 implementation and test files remain uncommitted/untracked, so Git cannot reconstruct exact per-slice additions or deletions from the current tree. Current-tree counts and historical estimates remain distinct.

### Slice 8 Validation

- `go test ./internal/infrastructure/llm ./cmd/servidor -count=1` — passed.
- `go test ./... -count=1` — passed.
- `go test -race ./...` — passed.
- `go build ./...` — passed.
- `go vet ./...` — passed.
- `go mod verify` — all modules verified.
- Scoped `gofmt -l` — clean; tracked and untracked `git diff --check` audits — clean.
- Native verify admission was attempted but blocked because `./.tools/gentle-ai` is unavailable in this worktree. No verify report, receipt, review artifact, or archive artifact was created.
