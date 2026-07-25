# Tasks: Domain Contract Types (Issue #6)

## Review Workload Forecast

| Field | Value |
|---|---|
| Estimated changed lines | 961–1,078: SDD planning 431 + domain/test 530–647 |
| SDD artifact accounting | 431: exploration 116, proposal 82, spec 61, design 81, state 43, tasks 48 |
| 400-line budget risk | High |
| Chained PRs recommended | Yes |
| Suggested split | PRs 1–3 artifacts; PRs 4–7 code; sequential to `main` |
| Delivery strategy | auto-chain |
| Chain strategy | stacked-to-main |

Decision needed before apply: No
Chained PRs recommended: Yes
Chain strategy: stacked-to-main
400-line budget risk: High

### Suggested Work Units

| Unit | ≤ lines | Scope | Focused test command | Runtime / rollback |
|---|---:|---|---|---|
| 1 | 198 | PR 1: exploration + proposal | `git diff --check` | N/A: artifact-only; revert PR 1 |
| 2 | 142 | PR 2: delta spec + design | `git diff --check` | N/A: artifact-only; revert PR 2 |
| 3 | 91 | PR 3: state + tasks | `git diff --check` | N/A: artifact-only; revert PR 3 |
| 4 | 210 | PR 4: 2.1, 2.2, 1.1, enum/profile 1.2 | `go test ./internal/domain/... -run 'Test(Enum|Perfil)' -v` | N/A: pure Go; revert enum/profile/tests |
| 5 | 165 | PR 5: 2.3, capacity 1.2 | `go test ./internal/domain/... -run 'Test(Capacidad|Recomendacion|Comprador|Proyecto)' -v` | N/A: pure Go; restore legacy files |
| 6 | 125 | PR 6: lead/plan 2.4 and 1.2 | `go test ./internal/domain/... -run 'Test(Lead|Mensaje|PlanNutricion|Hito)' -v` | N/A: pure Go; revert lead/plan/tests |
| 7 | 147 | PR 7: ficha 2.4/1.2, 3.1–3.3 | `go test ./internal/domain/... -v && go test ./internal/pipeline/... -v` | N/A: pure Go; revert ficha/final tests |

## Phase 1: Contract Tests

- [x] 1.1 Create `internal/domain/perfil_test.go` table tests for issue cases: `int64`, `int`, JSON `float64`, bool, string, absent, and incompatible accessors; verify only `VERIFICADO_BASE` and exact 18/4 key sets.
- [ ] 1.2 Add enum JSON-wire assertions and reflection/JSON cases for every specified field type, tag, omitted `-`, pointer nullability, and exported schema.

## Phase 2: Domain Declarations

- [x] 2.1 Create `internal/domain/enums.go` with the exact 11 typed-string enum groups and literals.
- [x] 2.2 Create `internal/domain/perfil.go` with `CampoPerfil`, `Perfil`, the two exact sets, and only the four specified accessors.
- [ ] 2.3 Create `internal/domain/capacidad.go` with capacity/recommendation types and field-identical `Comprador`/`Proyecto`; delete `comprador.go` and `proyecto.go` after consumer search.
- [ ] 2.4 Create `internal/domain/lead.go`, `plan.go`, and `ficha.go` with only the Contract layouts/tags and stdlib (`time`) dependency.

## Phase 3: Compatibility and Evidence

- [ ] 3.1 Run `go test ./internal/domain/... -v` and `go test ./internal/pipeline/... -v`; repair only contract-test failures.
- [ ] 3.2 Run `go build ./internal/domain/...`, `go build ./...`, `go vet ./internal/domain/...`, and `go list -deps ./internal/domain/... | grep -E 'internal/(usecase|adapters|infrastructure)'` (expect no output).
- [ ] 3.3 Confirm no pipeline or Docs changes and record each slice’s focused evidence before its PR.
