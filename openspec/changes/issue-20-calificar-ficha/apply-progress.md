# Apply Progress: issue-20-calificar-ficha

## Batch
- Change: `issue-20-calificar-ficha`
- Slice: qualification-edge-hardening after independent verification
- Branch: `test/issue-20-qualification-edges`
- Worktree: `/tmp/vivi-issue-20-verify`
- Mode: Standard (strict TDD disabled by project configuration)
- Delivery: feature-branch-chain, focused test-only slice; no size exception

## Completed Tasks
- [x] 1.1–1.4 Slice 1 qualification, shared decision, routing, CAS, and event ordering (inherited from Slice 1).
- [x] 2.1 Eligibility/read/cancellation tests prove invalid leads and failures produce no ficha or lead writes.
- [x] 2.2 `GenerarFicha` reconstructs the shared deterministic decision and emits ordered Contract content, warning, benefits, rent argument, and strict withdrawal alert.
- [x] 2.3 Fixed-clock/ID tests cover deterministic content, recommendation parity, threshold `.20` versus `>.20`, no-rent behavior, and read-only output isolation.
- [x] 2.4 Ficha-first persistence and repair retry preserve stable ficha ID/time, leave failed lead CAS unchanged, and complete `ENTREGADO` on retry without events; the public API is `Ejecutar(ctx context.Context, leadID string) (*domain.Ficha, error)` and successful handoff stamps `ActualizadoEn` from `Reloj`.
- [x] 3.1 Slice validation and physical source budget evidence recorded below.
- [x] Qualification edge hardening: added three deterministic tests for ratio `>1.2` with confidence `<1`, exact catalog-key zone behavior, and complete copied `RutaDecidida` payload fields.
- [ ] 3.2 Commit/work-unit closure remains pending because this executor was explicitly instructed not to commit.

## Implementation
- `internal/usecase/generar_ficha.go`: unchanged in this batch; previous Slice 2 implementation remains cumulative progress.
- `internal/usecase/generar_ficha_test.go`: unchanged in this batch; previous Slice 2 tests remain cumulative progress.
- `internal/usecase/calificar_lead.go`: unchanged; no production code was modified.
- `internal/usecase/calificar_lead_edges_test.go`: added 63 formatted test lines covering the three verifier warnings. The priority case uses ratio `2.0` and confidence `0.85`, proving priority is `1.2 * confidence`; the zone case compares canonical and noncanonical keys and asserts the normalized missing-zone distance; the event case asserts `ruta`, `semaforo`, `consume_cupo_10`, `prioridad`, and copied `recomendaciones`.

## Work Unit Evidence
| Evidence | Result |
|---|---|
| Focused test | `go test ./internal/usecase -run 'TestCalificarLead(PriorityCapsRatioAndUsesPartialConfidence|NonCanonicalCatalogKeyOmitsKNNZone|RutaDecididaPayloadContainsCopiedDecisionData)$' -count=1 -v` — PASS; all 3 edge tests passed. |
| Runtime harness | N/A — synchronous provider-free usecase has no HTTP, process, shell, or external runtime boundary; the focused tests exercise the deterministic usecase, motor KNN, event payload, and fake repository/bus boundary. |
| Rollback boundary | Revert only `internal/usecase/calificar_lead_edges_test.go`; no production, domain, motor, ports, infrastructure, adapter, HTTP, frontend, migration, config, Contract, or Wiki path changed. |

## Full Validation
- `go test ./internal/usecase -run 'TestCalificarLead(PriorityCapsRatioAndUsesPartialConfidence|NonCanonicalCatalogKeyOmitsKNNZone|RutaDecididaPayloadContainsCopiedDecisionData)$' -count=1 -v` — PASS.
- `go test ./... -count=1` — PASS.
- `go test -race ./... -count=1` — PASS.
- `go build ./...` — PASS.
- `go vet ./...` — PASS.
- `go mod verify` — PASS (`all modules verified`).
- `test -z "$(gofmt -l internal/usecase/calificar_lead_edges_test.go)"` — PASS; no output.
- `git diff --check` — PASS.
- Test-only batch source budget — `63` authored lines, `0` runtime lines, `63` test lines; well below the 400-line hard ceiling.

## Deviations and Risks
- No production or design deviation; this batch only adds readable deterministic coverage for the verifier's three concrete warnings.
- Commit/work-unit closure is intentionally pending under the explicit no-commit constraint; independent verification should rerun after the parent workflow establishes review authority and work-unit closure.
- The pre-existing authority-only verification blockers remain unchanged: bound review receipt and commit closure are outside this executor's scope.

## Next
Run independent SDD verification after the parent workflow handles the existing authority/commit gates. Do not archive in this executor.
