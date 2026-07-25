# Tasks: Security and Resilience Decorators (Issue #17)

## Review Workload Forecast

| File | Kind | Lines |
|---|---|---:|
| `internal/infrastructure/llm/guardarrailes.go` | production | 96 |
| `internal/infrastructure/llm/metricas.go` | production | 59 |
| `internal/infrastructure/llm/factory.go` | production | 15 |
| `internal/infrastructure/llm/health.go` | production | 9 |
| `tests/adversarios.json` | fixture | 19 |
| `internal/infrastructure/llm/guardarrailes_test.go` | test | 119 |
| `internal/infrastructure/llm/metricas_test.go` | test | 49 |
| `internal/infrastructure/llm/resilience_test.go` | test | 23 |
| **Total** | **authored implementation** | **389** |

| Field | Value |
|---|---|
| 400-line budget risk | Medium; 11-line buffer |
| Chained PRs recommended | No |
| Suggested split | One PR → `main` |
| Delivery strategy | `single-pr` |
| Chain strategy | `pending` |

Decision needed before apply: No
Chained PRs recommended: No
Chain strategy: pending
400-line budget risk: Medium

### Suggested Work Units

| Unit | Goal | Likely PR | Focused test command | Runtime harness | Rollback boundary |
|---|---|---|---|---|---|
| 1 | Full decorator slice | 1 → `main` | `go test ./internal/infrastructure/llm -count=1` | N/A: deterministic fakes exercise both port methods | Revert factory/health; new files become unreachable |

## Phase 1: Contract-Faithful Decorators

- [x] 1.1 Create `internal/infrastructure/llm/guardarrailes.go`: classify all 15 fixture categories before delegation; return generic/privacy `FUERA_DE_DOMINIO` templates with zero text/audio calls.
- [x] 1.2 Add post-validation for prompt/skill/config/internal-file markers, foreign `lead_id`, and `validResponse` money denial; locally replace unsafe output without retry/fallback.
- [x] 1.3 Create `internal/infrastructure/llm/metricas.go`: injected clock/writer JSON logs only lead ID, event, latency, provider, and typed outcome; delegate both port methods and breaker state.

## Phase 2: Existing Resilience Composition

- [x] 2.1 Modify `internal/infrastructure/llm/factory.go` for `Metricas(Guardarrailes(FallbackProvider))` and the existing `WithMetrics` seam; add no breaker, limiter, timeout, retry, or fallback.
- [x] 2.2 Modify `internal/infrastructure/llm/health.go` to query `CircuitBreakerState()` by capability interface, preserving `/salud` through wrappers.

## Phase 3: Deterministic Evidence

- [x] 3.1 Replaced the altered fixture with exactly 15 categorized Issue #17 inputs: jailbreak 1/2/14, extraccion 3/4/13, rol 5/6/7, terceros 8/9/15, and fuera_dominio 10/11/12, preserving each literal Spanish `texto`.
- [x] 3.2 Create table-driven `guardarrailes_test.go`: 15/15 templates/zero calls; benign one-call pass-through; blocked audio; leak, foreign lead, unauthorized/allowed motor money, and no-extra-call cases.
- [x] 3.3 Create `metricas_test.go`: decode JSON for accepted, rejected, failed, and breaker-open outcomes; assert required fields/latency and absence of seeded message, audio, cédula, prompt, response, credential, and error text.
- [x] 3.4 Extend `resilience_test.go` (no breaker test) for factory order and decorated `ABIERTO` health after existing three-failure/60-second behavior; guardrails do not affect eligibility.

## Phase 4: Verification and Budget Gate

- [x] 4.1 Run `go test ./internal/infrastructure/llm -count=1`, `go test ./...`, `go build ./...`, `go vet ./...`, `go mod verify`, and `git diff --check`; record outcomes.
- [x] 4.2 Before commit, run `git diff --stat` and count additions plus deletions. The reproducible non-SDD count is exactly 400: tracked `+34/-8` plus 358 untracked additions; prune completed with no size exception.
- [x] 4.3 Commit the single work unit with tests/fixture; record focused result, runtime N/A rationale, and factory/health rollback boundary. (Reconciled at archive against reviewed commit `d7cc14c`.)
