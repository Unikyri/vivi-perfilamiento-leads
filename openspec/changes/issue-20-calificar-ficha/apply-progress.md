# Apply Progress: issue-20-calificar-ficha

## Batch
- Change: `issue-20-calificar-ficha`
- Slice: 1 — qualification
- Branch: `feat/issue-20-calificar-ficha`
- Mode: Standard (strict TDD disabled by project configuration)
- Delivery: feature-branch-chain, PR #1 boundary; no size exception

## Completed Tasks
- [x] 1.1 Deterministic guards/read-failure tests with no lead save or event.
- [x] 1.2 Provider-free qualification service and package-private decision/candidate helpers.
- [x] 1.3 Conversion, route, priority-cap, semáforo, cupo, candidate-zero fallback, KNN, and deterministic recommendation tests.
- [x] 1.4 Route/state mapping, required event payload, and durable lead-save-before-single-event ordering.
- [ ] 2.1–2.4 Slice 2 ficha generation (reserved for the next apply batch).
- [ ] 3.2 Final commit/work-unit closure (commit explicitly forbidden in this executor session).

## Implementation
- `internal/usecase/calificar_lead.go`: added `CalificarLead`, qualification result/input types, capacity candidate selection, exact catalog-zone KNN input, conversion predicate, matrix routing, priority cap, semáforo/cupo mapping, state transitions, CAS persistence, and post-save `RutaDecidida` publication.
- `internal/usecase/calificar_lead_test.go`: rewrote the previously semicolon-compressed test into gofmt-formatted table-driven subtests and readable helpers, preserving guards/no-write, candidate/fallback, KNN K=30, deterministic recommendations, conversions, routes, priority cap, semáforo/cupo, event payload, save-before-event, and save-failure coverage.

## Work Unit Evidence
| Evidence | Result |
|---|---|
| Focused test | `go test ./internal/usecase -run TestCalificarLead -count=1` — PASS; all three qualification test functions and subcases passed. |
| Runtime harness | N/A — slice 1 is a synchronous provider-free application service with no HTTP, process, shell, or external runtime boundary. Deterministic fake repository/bus ordering scenario exercised the durability boundary. |
| Rollback boundary | Revert/delete only `internal/usecase/calificar_lead.go` and `internal/usecase/calificar_lead_test.go`; slice 2 remains untouched. |

## Full Validation
- `go test ./...` — PASS.
- `go vet ./...` — PASS.
- `go build ./...` — PASS.

## Scope and Workload
- Final formatted source count: `238` lines in `calificar_lead.go` + `162` lines in `calificar_lead_test.go` = `400` total, exactly at the hard slice ceiling. Production behavior was not changed by this remediation; the test source was rewritten for readability.
- Only the two assigned usecase files and active SDD artifacts were changed; no domain, motor, port, adapter, infrastructure, HTTP, frontend, migration, config, Contract, or Wiki paths were modified.
- Deviation: none from the slice-1 design. The existing `ASESOR`/`CALIFICADO` result is read-only on repeat invocation to avoid duplicate CAS/event delivery.

## Next
Slice 1 is ready for independent verification. Slice 2 should implement `GenerarFicha` and its tests on this branch without modifying the qualification scope.
