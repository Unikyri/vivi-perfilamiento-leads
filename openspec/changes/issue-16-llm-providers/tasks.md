# Tasks: LLM Provider Adapters and Fallback (Issue #16)

## Review Workload Forecast

| PR | Deliverable and forecast | Focused test | Runtime / rollback |
|---|---|---|---|
| 1 | Errors, transport, prompt, parser, tests; ≤300 | `go test ./internal/infrastructure/llm -run 'Test(Prompt|Parser|Transport)'` | N/A: fake transport; revert `llm/{errors,transport,prompt,parser}*.go` |
| 2 | Gemini/Qwen HTTP adapters, tests; ≤400 | `go test ./internal/infrastructure/llm -run 'Test(Gemini|Qwen)'` | N/A: live providers prohibited; revert `llm/{gemini,qwen}*.go` |
| 3 | Factory, fallback, breaker, wiring, tests; ≤400 | `go test ./internal/infrastructure/llm ./cmd/servidor` | Fake-provider health test; revert PR 3 files and `main.go` wiring |

Delivery strategy: chained feature-branch chain; PR 1 base `feat/issue-16-llm-providers`, PR 2 base PR 1, PR 3 base PR 2. Total ≤1,100 authored lines; no real APIs, secrets, or message/audio logging.

Decision needed before apply: No
Chained PRs recommended: Yes
Chain strategy: feature-branch-chain
400-line budget risk: High

### PR Boundary Handoffs
- **1 → 2:** sanitized errors, `HTTPDoer`, deterministic prompt, strict parser; keep `LLMProvider`, `Config`, `go.mod`, and CI unchanged.
- **2 → 3:** tested adapters and capability/error types; runtime remains unwired.
- **3 finish:** compose without port/config changes; each PR links approved Issue #16 with one `type:feature` label.

## PR Slice 1: Shared Contract Core (≤300)
- [x] 1.1 Create `llm/errors.go` and `transport.go`: sanitized typed errors, injected `HTTPDoer`, cancellation, 8-second cap. **Trace:** Transient failure; Safe validation.
- [x] 1.2 Create `llm/prompt.go`: fixed instruction separate from delimited user data; sorted maps, ordered history, no `VERIFICADO_BASE`. **Trace:** Prompt injection.
- [x] 1.3 Create `llm/parser.go`: strict Contract §7 JSON/EOF, fields/sources/confidence/level/actions, and motor-only money validation. **Trace:** Valid output; Unauthorized output.
- [x] 1.4 Add table-driven `prompt_test.go`, `parser_test.go`, `transport_test.go`: injection, invalid-output/no-fallback, timeout/cancel, redacted errors; run PR 1 test.

## PR Slice 2: Provider HTTP Adapters (≤400)
- [x] 2.1 Create `llm/gemini.go`: private structured text/inline-audio REST DTOs, shared prompt/parser, typed mapping. **Trace:** Gemini capability; One call/recovery.
- [x] 2.2 Create `llm/qwen.go`: configured OpenAI-compatible text request; unsupported-audio before transport and never reroute. **Trace:** Qwen primary audio.
- [x] 2.3 Add `gemini_test.go`/`qwen_test.go` with `httptest`/fake transport: body/auth, extraction, one normal call, one malformed retry, 429/5xx typing, audio zero calls, no sensitive errors; run PR 2 test.

## PR Slice 3: Composition and Resilience (≤400)
- [x] 3.1 Create `llm/factory.go` using only existing provider/fallback/key/base-URL config; typed invalid-name/credential errors. **Trace:** Qwen configuration switch; Scope boundary.
- [x] 3.2 Create `llm/fallback.go`/`breaker.go`: eligible routing, 10-second budget, rate limit, 3-failure/60-second half-open breaker, typed exhaustion. **Trace:** Gemini rate-limit fallback; Both unavailable.
- [x] 3.3 Add composable guardrail/metrics seams and `resilience_test.go`: guardrail zero calls, non-fallback validation/capability, breaker reset/bypass, no secret/content logs. **Trace:** Guardrail rejection; Unauthorized output.
- [x] 3.4 Wire only `cmd/servidor/main.go` to factory, `Nombre()`, and breaker health; add `factory_test.go`/wiring test while preserving `puertos.go`/`config.go`; run PR 3 test.

## Final Cross-Slice Verification
- [x] 4.1 After PRs 1–3 and all bounded remediation slices, run `go test ./... && go build ./cmd/servidor`; verify every scenario, unchanged ports/config, and zero real API/secrets use. Do not create archive or review-receipt artifacts.

## Bounded Chained Remediation Slice 4 (≤400 physical added lines)
- [x] 5.1 Make exhausted malformed/non-JSON output fallback-eligible after exactly one same-provider retry while keeping valid Contract-invalid JSON non-fallbackable; add the primary-two/fallback-one runtime regression.
- [x] 5.2 Add an injected-clock, bounded per-provider token bucket before provider invocation; prove denial bypasses primary and uses configured fallback without live APIs.
- [x] 5.3 Retain the constructed provider composition in the server health path and report the live primary breaker state; preserve `unconfigured` without credential/database test hacks.
- [x] 5.4 Correct Slice 1 accounting and restrict numeric validation only for currency amounts; record all validation and exact physical-budget evidence.

### Remediation Slice Accounting Clarification
- Slice 1's original source/test snapshot was historically estimated at **295 authored physical lines**. That baseline excludes later additions to `errors.go`; current `errors.go` physical lines MUST NOT be presented as the original 295, and the uncommitted worktree prevents Git reconstruction.
- Historical estimate for the Slice 4 remediation additions: **247 authored physical lines** — `rate_limiter.go` 74 + `health.go` 23 + `cmd/servidor/main_test.go` 41 + 47 appended resilience-test lines + 11 appended parser-test lines + 29 fallback additions + 13 main additions + 8 parser additions + 1 errors addition. The worktree was uncommitted/untracked, so this historical estimate cannot be reconstructed from Git.
- This slice remains below the **400 physical added-line** bound. The numeric validation warning is resolved: only currency-marked amounts are restricted to motor values; count/date tokens remain unrestricted.

## Slice 4 Work Unit Evidence
- **Focused test:** `go test ./internal/infrastructure/llm ./cmd/servidor -count=1` — exit 0; malformed recovery/fallback, token-bucket deny/refill, parser currency scope, and live health-path tests passed.
- **Runtime harness:** fake Gemini HTTP transport, fake providers, injected fake clock, and `httptest` health handler; zero live provider calls, credentials, or message/audio logging.
- **Rollback boundary:** revert only `rate_limiter.go`, `health.go`, the focused test additions, and the fallback/parser/main changes listed above; leave Slice 1/2 implementation and unrelated modules intact.

- **Fallback construction:** `NewFromConfig` now ignores an unavailable or invalid optional fallback (`fallback, _ := providerFor(...)`) after primary validation; a valid Gemini primary remains usable without a Qwen key. The focused test exercises both `config.Cargar`'s default `LLM_FALLBACK=qwen` and explicit `LLM_FALLBACK=qwen` with only a Gemini key.
- **Historical budget estimate:** Slice 3-owned `factory.go` (48) + `fallback.go` (131) + `breaker.go` (87) + `resilience_test.go` (128) was recorded as **394 physical lines**; prior `errors.go` content and touched error kinds were separately estimated, with conservative review accounting recorded as **400 changed lines**. Because the files were uncommitted/untracked, these are historical estimates rather than Git-reconstructible deltas. No earlier Slice 1/2 implementation was altered.
- **Validation:** `gofmt`; `go test ./internal/infrastructure/llm ./cmd/servidor -count=1`; `go test ./... -count=1`; `go test -race ./...`; `go build ./...`; `go vet ./...`; `go mod verify`; `git diff --check`; and untracked Go `git diff --no-index --check` all passed.
- **Scope:** This remediation does not mark final verification task 4.1 complete and creates no verify/archive/review artifacts.

## Bounded Chained Remediation Slice 5 (≤400 physical added lines)
- [x] 6.1 Preserve a real `FallbackProvider` composition from `NewFromConfig` for valid single-provider and optional-invalid-fallback configurations, including per-provider limiter/breaker state and live health after three failures; add direct runtime regressions for breaker health and limiter denial.
- [x] 6.2 Parse Colombian monetary expressions such as `2 millones`, `$2 millones`, and `COP 2 millones` as 2,000,000; reject unauthorized monetary values while accepting ordinary counts/dates; add precise parser regressions.

### Slice 5 Work Unit Evidence
- **Focused:** `go test ./internal/infrastructure/llm ./cmd/servidor -count=1` — exit 0; factory-backed no-fallback breaker/health and limiter tests plus parser expression tests passed.
- **Runtime:** `NewFromConfig` composition with fake provider substitution, injected token bucket clock, and direct parser execution; no live APIs, credentials, or message/audio logging.
- **Rollback:** revert only `internal/infrastructure/llm/factory.go`, `parser.go`, `parser_test.go`, and `resilience_test.go` changes from this slice; retain Slice 4 and earlier implementation.

### Slice 5 Physical Budget Accounting
- Historical estimate for Slice 5 additions: **100 authored physical lines** — `factory.go` 1 + `parser.go` 43 + `parser_test.go` 17 + `resilience_test.go` 39. Because the implementation remained uncommitted/untracked, this is not a Git-reconstructible per-slice delta.
- Slice 5 is historically estimated at **300 lines below** the ≤400 bound; historical cumulative authored implementation estimate is **1,441** (prior estimate 1,341 + estimated Slice 5 additions 100).
- No verify report, review receipt, commit, push, PR, archive, or issue closure was created.

## Bounded Chained Remediation Slice 6 (≤400 physical added lines)
- [x] 7.1 Recognize bare dot-grouped monetary amounts only when a bounded nearby context term (`mensual`, `mensuales`, `cuota`, `ingreso`, `salario`, `presupuesto`, `precio`, or `valor`) is present; normalize through the existing overflow-safe parser and require the result in `NUMEROS_DEL_MOTOR`.
- [x] 7.2 Add focused parser regressions for authorized contextual amounts, unauthorized contextual amounts, and ordinary grouped counts/date text without monetary context.

### Slice 6 Work Unit Evidence
- **Focused:** `go test ./internal/infrastructure/llm -run 'TestParser' -count=1` — exit 0; existing currency/million cases and new contextual bare-group cases passed.
- **Runtime:** Direct `ParseSalida` parser harness through focused Go tests; no provider/network/API, credential, message, or audio calls.
- **Rollback:** Revert only `internal/infrastructure/llm/parser.go` and `internal/infrastructure/llm/parser_test.go` Slice 6 additions; retain Slice 5 factory, resilience, and existing parser behavior.

### Slice 6 Physical Budget Accounting
- Historical estimate for Slice 6 additions: **76 authored physical lines** — `parser.go` 42 + `parser_test.go` 34. The uncommitted/untracked worktree prevents independent Git reconstruction of this slice.
- Slice 6 is **324 lines below** the ≤400 physical added-line bound.
- Historical cumulative authored implementation estimate is **1,517** (prior estimate 1,441 + estimated Slice 6 additions 76); the uncommitted/untracked worktree prevents Git reconstruction.

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

## Bounded Chained Remediation Slice 7 (≤400 physical added lines)
- [x] 8.1 Validate contextual Colombian bare amounts using lexical adjacency and sentence-aware boundaries; accept grouped zero-decimal and ungrouped whole values only when they match `NUMEROS_DEL_MOTOR`, reject non-integral grouped decimals without float rounding, and preserve `$`, `COP`, `pesos`, and `millones` rules.
- [x] 8.2 Replace implicit 1 req/s defaults with named, documented bounded token-bucket defaults (burst 3, refill 0.5 tokens/second) and regression coverage without new configuration.
- [x] 8.3 Reset breaker consecutive-failure state after ineligible failures while retaining success reset; require exactly three consecutive eligible failures and add regressions.

### Slice 7 Work Unit Evidence
- **Focused:** `go test ./internal/infrastructure/llm ./cmd/servidor -count=1` — exit 0; parser, limiter, fallback, breaker, and health tests passed.
- **Runtime harness:** direct `ParseSalida`, fake providers, injected fake clock, and fake HTTP transports; no live provider/API, credential, message, or audio calls.
- **Rollback:** revert Slice 7 changes in `internal/infrastructure/llm/parser.go`, `parser_test.go`, `rate_limiter.go`, `fallback.go`, `breaker.go`, and `resilience_test.go`; revert the matching Slice 7 design/progress/task/state documentation only.

### Slice 7 Physical Budget Accounting
- Historical estimate for Slice 7 net physical implementation additions: **247 lines** — `parser.go` +141, `parser_test.go` +44, `rate_limiter.go` +8, `fallback.go` +0, `breaker.go` +4, `resilience_test.go` +50. This cannot be independently reconstructed from Git because the Issue 16 files are uncommitted/untracked.
- Slice 7 is historically estimated at **153 lines below** the ≤400 physical-line bound; historical cumulative authored implementation estimate is **1,764** (prior estimate 1,517 + estimated Slice 7 additions 247).

### Slice 7 Validation
- `gofmt` / `gofmt -l` — passed.
- `go test ./internal/infrastructure/llm -run 'TestParser|TestBreaker|TestDefaultProviderRateLimiter|TestFallback|TestTokenBucket' -count=1` — passed.
- `go test ./internal/infrastructure/llm ./cmd/servidor -count=1` — passed.
- `go test ./... -count=1` — passed.
- `go test -race ./...` — passed.
- `go build ./...` — passed.
- `go vet ./...` — passed.
- `go mod verify` — all modules verified.
- `git diff --check` and no-live-network/credential audits — passed.
- No dependency, environment-variable, config-signature, verify-report, review, commit, push, PR, archive, or issue action was performed.


## Bounded Final Verification Slice 8 (≤400 physical added lines)

- Materialized the prior deterministic verifier probe for signed `int64` overflow in `internal/infrastructure/llm/parser_test.go`; an overflowing monetary value is invalid and cannot wrap to an authorized `MinInt64` value.
- Materialized the prior deterministic composite-breaker probe in `internal/infrastructure/llm/resilience_test.go`; three eligible primary failures open the breaker, and the fourth invocation bypasses primary and succeeds through fallback.
- No production code, health vocabulary, ports, configuration signatures, dependencies, verify report, receipt, archive, commit, push, PR, or issue action was changed or created.

### Slice 8 Work Unit Evidence

| Evidence | Result |
|---|---|
| Focused test | `go test ./internal/infrastructure/llm -run 'TestParserInt64OverflowIsInvalidWithoutWrap|TestCompositeBreakerBypassesPrimaryAfterThreeEligibleFailures' -count=1` — exit 0. |
| Runtime harness | Deterministic in-process parser execution and fake providers/rate limiters; no network, credentials, provider, message, or audio calls. |
| Rollback boundary | Revert only the two appended test blocks in `internal/infrastructure/llm/parser_test.go` and `internal/infrastructure/llm/resilience_test.go`; retain all prior implementation slices. |

### Slice 8 Accounting Reconciliation

- Current-tree physical line counts after this slice: `parser_test.go` **141**, `resilience_test.go` **289**, `parser.go` **317**, `breaker.go` **91**, `fallback.go` **158**, `rate_limiter.go` **81**, `cmd/servidor/main.go` **102**.
- Direct Slice 8 edit content is **32 physical lines** (9 in `parser_test.go` and 23 in `resilience_test.go`), below the **400-line** bound. This is an edit measurement, not a reconstructed Git diff.
- Historical Slice 1–7 figures in this artifact are estimates/snapshots only. Issue 16 implementation and test files remain uncommitted/untracked, so Git cannot reconstruct exact per-slice additions or deletions from the current tree. Current-tree counts and historical estimates MUST remain distinct.

### Slice 8 Validation

- `go test ./internal/infrastructure/llm ./cmd/servidor -count=1` — passed.
- `go test ./... -count=1` — passed.
- `go test -race ./...` — passed.
- `go build ./...` — passed.
- `go vet ./...` — passed.
- `go mod verify` — all modules verified.
- Scoped `gofmt -l` — clean; tracked and untracked `git diff --check` audits — clean.
- Native verify admission was attempted but blocked because `./.tools/gentle-ai` is unavailable in this worktree; no `verify-report.md`, receipt, or archive artifact was created.
