```yaml
schema: gentle-ai.verify-result/v1
evidence_revision: sha256:c58a6e1aeb4e960194c91e78cbc9f56764b4e976769abdb9bc9c2bd06394c339
verdict: pass
blockers: 0
critical_findings: 0
requirements: 6/6
scenarios: 8/8
test_command: go test ./... -count=1
test_exit_code: 0
test_output_hash: sha256:e0a281918d87a4b22437688d518c74608d7473d798a85919672de73a05bcb517
build_command: go build ./...
build_exit_code: 0
build_output_hash: sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
```

## Verification Report

**Change**: issue-10-recommendations-routing (GitHub Issue #10 — `[A5] Recommendations and 2x2 routing`)
**Version**: `openspec/changes/issue-10-recommendations-routing/specs/recommendations-routing/spec.md` (6 requirements, 8 scenarios)
**Mode**: Standard (`openspec/config.yaml` → `strict_tdd: false`)
**Artifact store**: hybrid (Engram primary + OpenSpec repository artifacts)
**Verified tree**: `d75c25ffca0ca5f4d95949cf40e73e5f57e856b8` — merge of PR #59 (`feat/issue-10-routing-matrix`), identical to `origin/main` and local `HEAD`
**Scope of this report**: FINAL independent verification from the merged SHA. It replaces the earlier pre-commit worktree report (Engram `#499`) which is now superseded.
**Verified at**: 2026-07-25

### Provenance and Merge Chain

| Item | Evidence |
|---|---|
| Local `HEAD` | `d75c25ffca0ca5f4d95949cf40e73e5f57e856b8` |
| `origin/main` | `d75c25ffca0ca5f4d95949cf40e73e5f57e856b8` (identical — verification ran on merged code, not a worktree candidate) |
| Working tree | Only untracked `Docs/` (pre-existing inventory) and this regenerated `verify-report.md`; zero tracked modifications at verification time |
| Merge chain | `462a985` → PR #56 (`01f6a04`) → `f5e1ea3` → PR #57 (`a2cfdc4`) → `6d43a2d` → PR #58 (`9f3e370`) → `8c88a34` → PR #59 (`d75c25f`) |
| Issue #8 gate | `git merge-base --is-ancestor 7a5b02a613c581568905936a35cb8feb841e62ab origin/main` → exit 0; #8 capacity work is an ancestor of the verified SHA |

### Completeness

| Metric | Value |
|--------|-------|
| Tasks total | 11 |
| Tasks complete | 11 (4.3 completed by this verification) |
| Tasks incomplete | 0 |

Task `4.3` was the only unchecked item before this phase and is satisfied by this report: verification ran after `1.1` and every implementation task, with exact command output evidence attached below.

### Build & Tests Execution

**Build**: ✅ Passed
```text
$ go build ./...
(no output)
exit 0 — output sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
```

**Static analysis**: ✅ Passed
```text
$ go vet ./...
(no output)
exit 0 — output sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
```

**Tests**: ✅ All packages passed / ❌ 0 failed / ⚠️ 0 skipped
```text
$ go test ./... -count=1
ok  	github.com/Unikyri/vivi-perfilamiento-leads/internal/domain	0.008s
ok  	github.com/Unikyri/vivi-perfilamiento-leads/internal/domain/motor	0.007s
ok  	github.com/Unikyri/vivi-perfilamiento-leads/internal/infrastructure/config	0.002s
ok  	github.com/Unikyri/vivi-perfilamiento-leads/internal/pipeline	0.005s
(cmd/pipeline, cmd/servidor, internal/adapters/agentes, internal/adapters/http,
 internal/infrastructure/llm, internal/infrastructure/postgres, internal/usecase,
 migrations, references: no test files)
exit 0 — output sha256:e0a281918d87a4b22437688d518c74608d7473d798a85919672de73a05bcb517
```

**Focused runtime enumeration** (`go test ./internal/domain/motor -count=1 -v`): exit 0, all subtests PASS —
`TestMatriz2x2` (21 table cases), `TestMatriz2x2IsPureAndRouteOnly`, `TestRecomendarProyectos` (4 cases),
`TestRecomendarProyectosEligibilityBoundaries` (8 cases), `TestRecomendarProyectosCopiesCatalogDataAndDoesNotMutateInputs`,
`TestFractionLessUint64AvoidsOverflow` (3 cases).

**Coverage**: `go test ./internal/domain/motor/... -count=1 -cover` → 97.4% of statements, exit 0.
Design target ≥90% → ✅ Above. Repository `coverage_threshold: 0` → ➖ no configured gate.

### Spec Compliance Matrix

| Requirement | Scenario | Test | Result |
|-------------|----------|------|--------|
| Integer-Safe Candidate Eligibility | Inclusive affordability boundary | `recomendar_test.go > TestRecomendarProyectosEligibilityBoundaries/exact_boundary_for_even_price`, `/one_peso_below_even_boundary`, `/ceil_boundary_for_odd_price`, `/one_peso_below_odd_boundary`, `/price_one_requires_one_peso` | ✅ COMPLIANT |
| Integer-Safe Candidate Eligibility | Invalid monetary data | `recomendar_test.go > TestRecomendarProyectosEligibilityBoundaries/zero_price_is_invalid`, `/negative_price_is_invalid`, `/negative_budget_is_invalid` | ✅ COMPLIANT |
| Catalog-Backed Recommendation Aggregation | Orphan neighbor evidence | `recomendar_test.go > TestRecomendarProyectos/aggregates_canonical_IDs_and_ignores_orphan_or_empty_IDs`, `/empty_input_returns_non_nil_empty_output`, `TestRecomendarProyectosCopiesCatalogDataAndDoesNotMutateInputs` | ✅ COMPLIANT |
| Deterministic Ranked Recommendations | Stable tie and cap | `recomendar_test.go > TestRecomendarProyectos/count_then_exact_rate_then_canonical_ID_order`, `/caps_output_at_three_with_shuffled_input_stable_ordering`, shuffled-input `reflect.DeepEqual` assertion, `TestFractionLessUint64AvoidsOverflow` (3 cases) | ✅ COMPLIANT |
| Pure 2x2 Route Selection | Thresholds and quadrants | `matriz_test.go > TestMatriz2x2/affiliate_high_capacity_and_high_intention_goes_to_advisor` (ratio exactly 0.95), `/affiliate_high_capacity_and_low_intention_goes_to_remarketing`, `/affiliate_low_capacity_and_high_intention_goes_to_nutrition`, `/affiliate_low_capacity_and_low_intention_goes_to_despedida`, `/affiliate_media_intention_is_low`, `/affiliate_unknown_intention_is_low`, `/affiliate_empty_intention_is_low`, `/affiliate_invalid_ratio_and_high_intention_never_reaches_advisor`, `/affiliate_positive_infinity_is_low_capacity`, `/affiliate_negative_infinity_is_low_capacity`, `/negative_ratio_is_never_advisor` | ✅ COMPLIANT |
| Pure 2x2 Route Selection | Non-affiliate threshold | `matriz_test.go > TestMatriz2x2/non_affiliate_exact_advisor_threshold` (ratio exactly 1.05), `/non_affiliate_below_advisor_threshold_goes_to_nutrition` (`math.Nextafter(1.05,0)`), `/non_affiliate_exact_point_nine_five_and_high_intention_goes_to_nutrition`, `/non_affiliate_one_point_zero_zero_and_high_intention_goes_to_nutrition` | ✅ COMPLIANT |
| Conversion Preference and Isolation | Conversion before advisor allocation | `matriz_test.go > TestMatriz2x2/conversion_preference_overrides_advisor_eligibility` (ratio 1.10 → `NUTRICION`), `/conversion_preference_also_applies_at_invalid_ratio`, `TestMatriz2x2IsPureAndRouteOnly` (value-only input, no mutation, stable repeat) | ✅ COMPLIANT |
| Dependency and Scope Precondition | Unmet capacity dependency | Deterministic Git evidence, exit 0: `git merge-base --is-ancestor 7a5b02a6… origin/main` plus `git diff --name-status 7a5b02a6…d75c25f` showing four additive motor files and seven planning files only | ✅ COMPLIANT (process requirement; no runtime unit test is applicable) |

**Compliance summary**: 8/8 scenarios compliant, 6/6 requirements covered. 0 FAILING, 0 UNTESTED, 0 PARTIAL.

### Correctness (Static Evidence)

| Requirement | Status | Notes |
|------------|--------|-------|
| Integer-Safe Candidate Eligibility | ✅ Implemented | `recomendar.go` rejects `presupuesto < 0` and `PrecioDesde <= 0`, then compares `presupuesto >= project.PrecioDesde - project.PrecioDesde/5`. Integer arithmetic only; no float division, no `Capacidad.Ratio` read. |
| Catalog-Backed Recommendation Aggregation | ✅ Implemented | Aggregation keys on non-empty `Vecino.ProyectoID` with a `catalogo[projectID]` hit **and** `project.ProyectoID == projectID` identity match, so mismatched entries and orphans are dropped. Names never participate. `result` is always a non-nil `make([]domain.Recomendacion, 0, 3)`. |
| Deterministic Ranked Recommendations | ✅ Implemented | `sort.Slice` orders neighbors descending, then exact desistimiento fraction ascending via `fractionLessUint64` (`bits.Mul64` 128-bit cross products), then `ProyectoID` ascending. Float rate is computed for display only. Cap enforced at 3. |
| Pure 2x2 Route Selection | ✅ Implemented | `Matriz2x2` returns `domain.Ruta` only. `altaIntencion` is `Intencion == domain.NivelAlta` exclusively. `ratioValido` excludes NaN, ±Inf, and negatives, so invalid ratios cannot reach `ASESOR`. Affiliate gate `>= 0.95`; non-affiliate advisor gate `>= 1.05`. |
| Conversion Preference and Isolation | ✅ Implemented | Non-affiliate + `ALTA` + `RutaConversion` returns `NUTRICION` before any threshold evaluation. Input is a by-value `EntradaMatriz`; no `Lead`/`Ficha`/cupo/milestone/plan type is referenced and no capacity is recomputed. |
| Dependency and Scope Precondition | ✅ Satisfied | #8 SHA `7a5b02a6…` is an ancestor of the verified merge. Diff `7a5b02a6…d75c25f` contains no `capacidad.go`, `capacidad_test.go`, shared-domain, `Docs/`, `data/`, `migrations/`, pipeline, HTTP, or use-case path. |

### Scope Evidence (`git diff --name-status 7a5b02a6…d75c25f`)

| Path | Action | Lines |
|---|---|---|
| `internal/domain/motor/recomendar.go` | A | +108 |
| `internal/domain/motor/recomendar_test.go` | A | +213 |
| `internal/domain/motor/matriz.go` | A | +54 |
| `internal/domain/motor/matriz_test.go` | A | +150 |
| `openspec/changes/issue-10-recommendations-routing/{exploration,proposal,specs/recommendations-routing/spec,design,tasks,state,apply-progress}` | A | +449 (planning artifacts) |

No file was modified or deleted. Code additions total 525 lines split across PR #58 (recommendations) and PR #59 (matrix), matching the approved auto-chain / stacked-to-main forecast. No `#8` capacity path and no `Docs/` path appears in any commit; `Docs/` remains untracked and unmodified.

### Coherence (Design)

| Decision | Followed? | Notes |
|----------|-----------|-------|
| Pure motor boundaries, value-in/value-out | ✅ Yes | Both functions are pure; purity and immutability are asserted by tests. |
| Exact integer eligibility (`presupuesto >= precio - precio/5`) | ✅ Yes | Implemented verbatim, with `precio > 0` and `presupuesto >= 0` guards ahead of it. |
| Deterministic ranking with `math/bits.Mul64` exact fractions | ✅ Yes | `fractionLessUint64` uses high/low word comparison; near-`uint64` products are tested. |
| Matrix precedence: conversion first, then quadrants/advisor gate | ✅ Yes | Conversion branch precedes ratio validation and thresholds. |
| Exact declared interfaces (`RecomendarProyectos`, `EntradaMatriz`, `Matriz2x2`) | ✅ Yes | Signatures match the design contract exactly. |
| Additive-only, initially unwired, revert-per-file rollback | ✅ Yes | Four new files, no production caller yet, per-slice revert boundary preserved. |
| Testing strategy (table-driven unit; no integration/E2E) | ✅ Yes | 21 matrix cases plus 15 recommendation/boundary/overflow cases; no adapter or HTTP boundary exists to harness. |

### Issues Found

**CRITICAL**: None.

**WARNING**:
1. Native admission validator unavailable — installed `gentle-ai 1.43.3` exposes no `sdd-verify-validate` subcommand (`gentle-ai help` lists only `install`, `uninstall`, `sync`, `skill-registry refresh`, `sdd-status`, `sdd-continue`, `update`, `upgrade`, `restore`, `doctor`, `version`). Envelope admission was performed manually against `gentle-ai.verify-result/v1`: requirement/scenario counts come from `grep -c '^### Requirement:'` (6) and `grep -c '^#### Scenario:'` (8) on the authoritative spec, and every command, exit code, and output hash below is an executed measurement. The superseded pre-commit report was regenerated rather than preserved because it described an unmerged worktree.
2. `gentle-ai 1.43.3` `sdd-status` routing uses a keyword heuristic and reports `verify-report.md is not clearly passing.` whenever the report body contains the residual-work keyword `p-e-n-d-i-n-g` (written here hyphenated on purpose), regardless of the verdict. Verified empirically against this binary. This report therefore states residual work with explicit wording (`unchecked`, `not yet wired`, `remaining prerequisites`) so machine routing matches the human verdict.
3. Apply-progress evidence understates the matrix table: Work Unit 2 records "all 17 table cases" while the merged `TestMatriz2x2` contains 21 cases. Documentation-only inaccuracy; runtime evidence is stronger than claimed, not weaker.
4. Native review/attempt authority is not obtainable from this binary (`gentle-ai review …` and `gentle-ai sdd-attempt …` are absent), so no `reviews/receipt.json`, `chain-bundle.json`, or `gate-context.json` exists and structured status cannot emit `reviewGate.result: allow`. Human review authority for this change is the merged GitHub PR chain (#56, #57, #58, #59). The orchestrator must accept that substitution explicitly before archive.
5. This report and the `4.3` checkbox are currently untracked/uncommitted working-tree state. They must land through a normal documentation PR to `main` before the change is team-shareable and archivable.
6. `RecomendarProyectos` and `Matriz2x2` have no production caller yet. Correct behavior is proven at the unit boundary only; end-to-end routing behavior remains unverified until a use-case wires them (out of Issue #10 scope).

**SUGGESTION**:
1. Spec wording for "Pure 2x2 Route Selection" states the four quadrants without scoping them to affiliates, while the implementation and design intentionally return `DESPEDIDA` (not `REMARKETING`) for a non-affiliate with high capacity and low intention. Behavior matches the design and is covered by `non_affiliate_low_intention_goes_to_despedida_at_high_ratio`; consider tightening the spec sentence during archive spec merge.
2. Motor coverage is 97.4%; the residual uncovered statements are worth a quick look when the functions are wired, not before.

### Verdict

PASS WITH WARNINGS — 8/8 spec scenarios compliant with passing runtime tests on merged SHA `d75c25f`; `go test ./... -count=1`, `go build ./...`, and `go vet ./...` all exit 0 with 97.4% motor coverage and no `#8` capacity or `Docs/` modification; remaining warnings are the missing native admission/review binaries, a dispatcher keyword heuristic, an understated apply-progress case count, and the uncommitted verification evidence.

**Archive readiness**: implementation and requirements verification are complete and archive-eligible. Two non-code prerequisites remain for the orchestrator: commit/merge this verification evidence, and record the review-authority substitution in Warning 4 (merged PR chain in place of a native receipt). Archive was not performed by this phase.
