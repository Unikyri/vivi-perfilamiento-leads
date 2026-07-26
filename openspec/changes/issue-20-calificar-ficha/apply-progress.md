# Apply Progress: issue-20-calificar-ficha

## Batch
- Change: `issue-20-calificar-ficha`
- Slice: 2 — deterministic ficha generation remediation
- Branch: `feat/issue-20-generar-ficha`
- Worktree: `/tmp/vivi-issue-20-ficha`
- Mode: Standard (strict TDD disabled by project configuration)
- Delivery: feature-branch-chain, Slice 2 boundary; no size exception

## Completed Tasks
- [x] 1.1–1.4 Slice 1 qualification, shared decision, routing, CAS, and event ordering (inherited from Slice 1).
- [x] 2.1 Eligibility/read/cancellation tests prove invalid leads and failures produce no ficha or lead writes.
- [x] 2.2 `GenerarFicha` reconstructs the shared deterministic decision and emits ordered Contract content, warning, benefits, rent argument, and strict withdrawal alert.
- [x] 2.3 Fixed-clock/ID tests cover deterministic content, recommendation parity, threshold `.20` versus `>.20`, no-rent behavior, and read-only output isolation.
- [x] 2.4 Ficha-first persistence and repair retry preserve stable ficha ID/time, leave failed lead CAS unchanged, and complete `ENTREGADO` on retry without events; the public API is `Ejecutar(ctx context.Context, leadID string) (*domain.Ficha, error)` and successful handoff stamps `ActualizadoEn` from `Reloj`.
- [x] 3.1 Slice validation and physical source budget evidence recorded below.
- [ ] 3.2 Commit/work-unit closure remains pending because this executor was explicitly instructed not to commit.

## Implementation
- `internal/usecase/generar_ficha.go`: removed `EntradaGenerarFicha`, `EntradaFicha`, and `SalidaGenerarFicha`; changed `Ejecutar` to the explicit leadID/pointer-result API; preserved all guards, read failures, ficha-first save, transition, CAS, retry, no-event, and no-LLM behavior; set `toSave.ActualizadoEn = uc.Reloj.Ahora()` immediately before the successful lead save.
- `internal/usecase/generar_ficha_test.go`: updated existing helper and retry calls to the explicit API and asserted the fixed Reloj timestamp on the persisted delivered lead.
- `internal/usecase/calificar_lead.go`: unchanged.

## Work Unit Evidence
| Evidence | Result |
|---|---|
| Focused test | `go test ./internal/usecase -run TestGenerarFicha -count=1` — PASS; all GenerarFicha tests and subcases passed. |
| Runtime harness | N/A — synchronous provider-free usecase has no HTTP, process, shell, or external runtime boundary; deterministic fake catalog, ficha repository, fixed clock/ID, and retry-after-CAS-failure harness exercised the persistence boundary. |
| Rollback boundary | Revert only `internal/usecase/generar_ficha.go` and `internal/usecase/generar_ficha_test.go`; Slice 1 qualification files remain untouched. |

## Full Validation
- `go test ./internal/usecase -run TestGenerarFicha -count=1` — PASS.
- `go test ./...` — PASS.
- `go test -race ./...` — PASS.
- `go build ./...` — PASS.
- `go vet ./...` — PASS.
- `go mod verify` — PASS (`all modules verified`).
- `go list ./...` — PASS; 14 packages listed.
- `gofmt -l internal/usecase/generar_ficha.go internal/usecase/generar_ficha_test.go` — no output.
- `git diff --check` — PASS.
- Legacy API grep (`EntradaGenerarFicha|EntradaFicha|SalidaGenerarFicha`) — no references.
- Scope/status — only the two Slice 2 Go files and active SDD artifacts are changed; no forbidden or Slice 1 implementation path changed.
- Physical Slice 2 source — `199 + 189 = 388` formatted lines, within the 400-line hard ceiling.

## Deviations and Risks
- No implementation deviation from the Slice 2 design; remediation only aligns the public API and lead timestamp behavior with Issue #20.
- Commit closure is intentionally pending under the explicit no-commit constraint; implementation is ready for independent verification after work-unit commit handling by the parent workflow.

## Next
Independent SDD verification should run after commit/work-unit closure. No events or LLM calls are introduced by this slice.
