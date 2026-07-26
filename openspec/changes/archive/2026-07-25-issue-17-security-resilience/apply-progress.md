# Apply Progress: Issue #17 Security and Resilience

## Status

- Change: `issue-17-security-resilience`
- Phase: corrective apply retry complete; ready for independent verification
- Delivery: single PR to `main`
- Mode: Standard (strict TDD disabled by project configuration)
- Next: `sdd-verify`

## Completed Tasks

- [x] 1.1 Added `Guardarrailes` input classification with generic/privacy-safe `FUERA_DE_DOMINIO` templates and zero provider calls.
- [x] 1.2 Added local output validation for prompt/skill/config/internal-file markers, foreign `lead_id`, and unauthorized monetary values via existing `validResponse`; unsafe output is replaced without retry or fallback.
- [x] 1.3 Added `Metricas` decorator with injected `Clock`/`io.Writer`, JSON fields limited to lead ID, event, latency, provider, and typed outcome; both provider methods and breaker state delegate.
- [x] 2.1 Composed production providers as `Metricas(Guardarrailes(FallbackProvider))` and injected the same metrics observer through the existing `WithMetrics` seam. No outer breaker, limiter, timeout, retry, or fallback was added.
- [x] 2.2 Changed health lookup to the package-local `CircuitBreakerState()` capability so `/salud` observes wrapped live breaker state.
- [x] 3.1 Replaced the altered fixture with the exact 15 `id`/`categoria`/`texto` rows from Issue #17: jailbreak 1/2/14, extraccion 3/4/13, rol 5/6/7, terceros 8/9/15, and fuera_dominio 10/11/12.
- [x] 3.2 Extended deterministic input matching for every canonical fixture prompt, including exact cédula, out-of-domain, SKILL.md, and other-lead forms; focused tests prove 15/15 containment and zero provider calls while benign housing text still passes through.
- [x] 3.3 Added JSON metrics tests for accepted, rejected, failed, and breaker-open outcomes, required operational fields/latency, and absence of seeded message/audio/cédula/prompt/response/credential/error text.
- [x] 3.4 Extended resilience tests for factory order, decorated live `ABIERTO` health after three eligible failures, and preserved limiter behavior; no new breaker test or breaker implementation was added.
- [x] 4.1 Re-ran focused LLM tests, full repository tests and race tests, build, vet, module verification, fixture validation, gofmt, and diff checks; all passed.
- [x] 4.2 Narrowed false-positive patterns, added seven benign one-call regressions, compacted duplicate/dead test code, and reproduced the exact 400-line non-SDD budget count without a size exception.
- [ ] 4.3 Commit the work unit — intentionally deferred because this remediation retry must not commit.

## Focused Remediation

- Scope included narrow deterministic classifier patterns, seven required benign regressions, and safe test/helper compaction to meet the hard budget.
- No predecessor resilience code, provider composition, breaker, timeout, limiter, fallback, public contract, or dependency changed.
- Canonical fixture verification: JSON parses to exactly 15 rows and retains all Issue #17 Spanish prompt text verbatim.

## Files Changed

- `internal/infrastructure/llm/guardarrailes.go` — new input/output guardrail decorator and breaker capability delegation.
- `internal/infrastructure/llm/metricas.go` — new privacy-safe JSON metrics decorator and existing metrics seam implementation.
- `internal/infrastructure/llm/factory.go` — decorated production composition.
- `internal/infrastructure/llm/health.go` — capability-based breaker health.
- `internal/infrastructure/llm/guardarrailes_test.go` — deterministic guardrail fixture, audio, output, and call-count evidence.
- `internal/infrastructure/llm/metricas_test.go` — structured metrics/privacy evidence.
- `internal/infrastructure/llm/resilience_test.go` — decorated factory/health/limiter assertions.
- `tests/adversarios.json` — exact 15-case adversarial fixture.
- `openspec/changes/issue-17-security-resilience/tasks.md` — tasks 1.1–3.4 and 4.1–4.2 checked; 4.3 remains pending by explicit no-commit instruction.
- `openspec/changes/issue-17-security-resilience/state.yaml` — remediation complete and verify next.

## Work Unit Evidence

| Evidence | Result |
|---|---|
| Focused fixture/regression | `go test ./internal/infrastructure/llm -run 'TestGuardrailsContainFixtureAndMakeZeroCalls\\|TestGuardrailsPermittedTextAndBlockedAudio' -count=1 -v` — exit 0; all 15 fixture rows blocked with zero calls and all 7 benign phrases delegated once. |
| Full LLM tests | `go test ./internal/infrastructure/llm -count=1` — exit 0. |
| Full repository tests | `go test ./... -count=1` — exit 0. |
| Race tests | `go test ./... -count=1 -race` — exit 0. |
| Build/static | `go build ./...`, `go vet ./...`, `go mod verify` — exit 0; all modules verified. |
| Formatting/whitespace | `gofmt` changed-file check and `git diff --check` — exit 0. |
| Fixture check | Python JSON parse — exit 0; exactly 15 rows with IDs 1..15. |
| Runtime harness | N/A — no runtime boundary exists in this unit; deterministic fake providers exercise both port methods and decorated health. |
| Rollback boundary | Revert `factory.go` to return the undecorated `FallbackProvider` and restore the concrete health assertion; new decorator/fixture/test files then become unreachable/removable without schema, config, Contract, or Issue #16 breaker changes. |

## Budget Evidence

The reproducible complete non-SDD authored snapshot is exactly **400 lines** (OpenSpec artifacts excluded):

```text
git diff --numstat
9  4  internal/infrastructure/llm/factory.go
1  1  internal/infrastructure/llm/health.go
24 3  internal/infrastructure/llm/resilience_test.go

# git ls-files --others --exclude-standard | awk '!/^openspec\\// {print}' | while read f; do wc -l < "$f"; done
89 internal/infrastructure/llm/guardarrailes.go
94 internal/infrastructure/llm/guardarrailes_test.go
98 internal/infrastructure/llm/metricas.go
62 internal/infrastructure/llm/metricas_test.go
15 tests/adversarios.json

tracked additions: 34
tracked deletions: 8
untracked non-OpenSpec additions: 358
exact total: 34 + 8 + 358 = 400
```

No size exception was used; the hard budget is met exactly.

## Deviations and Risks

- No behavioral deviation from the design. Existing Issue #16 `FallbackProvider` remains the sole breaker/limiter/timeout/fallback owner.
- Classifier matching is intentionally narrower than the previous implementation: exact adversarial phrases and explicit foreign identifiers remain blocked, while normal housing/affiliate language passes through.
- Cross-lead isolation is limited to explicit `lead_id` output markers at this provider seam; session/tool identity remains Issue #19 scope.
- Task 4.3 remains unchecked because the user explicitly prohibited commit/PR actions in this retry; independent verification is the next action.
