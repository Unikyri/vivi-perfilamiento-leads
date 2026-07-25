# Apply Progress: Lead State Machine (Issue #7)

## Mode

Standard mode (`strict_tdd: false`), single-pr delivery. No prior apply-progress existed.

## Completed Tasks

### Phase 1: Domain Policy

- [x] 1.1 Added the private, package-owned adjacency map with exactly the 11 Contract lifecycle edges and a private deterministic destination order.
- [x] 1.2 Added `ErrTransicionInvalida`, `PuedeTransicionar`, and defensive ordered `EstadosPosibles`; terminal and unknown sources return zero-length results.
- [x] 1.3 Added guarded `(*Lead).Transicionar`, validating before assignment and preserving state on rejection.

### Phase 2: Domain Tests

- [x] 2.1 Added table-driven coverage for all 11 valid edges, successful assignment, and ordered query results.
- [x] 2.2 Added exhaustive invalid-pair coverage across known and unknown states, typed error assertions, error code checks, and atomicity checks.
- [x] 2.3 Covered `CERRADO`, `DESPEDIDO`, and unknown sources as terminal/unreachable.
- [x] 2.4 Covered defensive copying of `EstadosPosibles(EstadoLeadCalificado)` and preserved Contract order.

### Phase 3: Verification and Delivery Guard

- [x] 3.1 Formatted both new Go files and ran the focused transition/query suite successfully.
- [x] 3.2 Ran the full domain suite, full repository tests, `go vet ./...`, and `go build ./...` successfully.
- [x] 3.3 Confirmed the domain imports remain standard-library-only (`encoding/json fmt time`) and reviewed the scoped diff; existing `enums.go` and `lead.go` were not modified.
- [x] 3.4 Prepared the code and tests as one reviewable work unit; commit is created after all checks pass.

## Work Unit Evidence

| Evidence | Result |
|---|---|
| Focused test | `go test ./internal/domain/... -run 'Test(Transicion|Estados)'` — PASS (`internal/domain`; motor has no test files) |
| Full domain tests | `go test ./internal/domain/...` — PASS |
| Full repository tests | `go test ./...` — PASS |
| Vet | `go vet ./...` — PASS |
| Build | `go build ./...` — PASS |
| Domain boundary | `go list -f '{{join .Imports " "}}' ./internal/domain` — `encoding/json fmt time`; no forbidden imports |
| Runtime harness | N/A — pure in-memory domain behavior has no runtime or integration boundary |
| Rollback boundary | Remove/revert `internal/domain/estado.go` and `internal/domain/estado_test.go`; existing Issue #6 domain files remain untouched |

## Files Changed

- `internal/domain/estado.go` — lifecycle policy, typed error, query, and guarded mutation.
- `internal/domain/estado_test.go` — table-driven state-machine tests.
- `openspec/changes/issue-7-lead-state-machine/tasks.md` — all 11 tasks marked complete.
- `openspec/changes/issue-7-lead-state-machine/apply-progress.md` — this evidence.

## Deviations and Risks

None. Implementation follows the design and remains standard-library-only; no ADK, agent, outer-layer, Docs, enum, field, or JSON-tag changes were introduced.
