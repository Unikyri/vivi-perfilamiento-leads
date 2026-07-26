# Apply Progress: ProcesarMensaje Core (PR 1)

## Status

- Change: `issue-19-procesar-mensaje`
- Mode: Standard (strict TDD disabled by project configuration)
- Delivery: feature-branch-chain, slice `core`, hard authored budget 400 lines
- Snapshot authored lines: 384 (`228` + `36` + `120`), within the 400-line cap
- Completed tasks: 1.1, 1.2, 1.3
- Remaining tasks: 2.1, 2.2, 3.1, 3.2

## Implementation

Created only the requested implementation files:

- `internal/usecase/procesar_mensaje.go`: `EntradaMensaje`, sentinels, input validation, lead-state/loading/history/cap checks, bounded turn context, one text/audio provider dispatch, basic recognized-field safety, and inbound/CAS/response persistence order.
- `internal/usecase/siguiente_pregunta.go`: exact six-field priority and conversational critical trio helpers.
- `internal/usecase/procesar_mensaje_test.go`: deterministic fixed-clock/ID/provider tests for invalid input, provider errors, helper behavior, one-call audio dispatch, caller transcript persistence, verified-field protection, unknown-field ignore, and ordered messages.

## Work Unit Evidence

| Evidence | Result |
|---|---|
| Focused test | `go test ./internal/usecase -run 'TestProcesarMensajeCore|TestSiguientePreguntaCore' -count=1` — PASS |
| Runtime harness | N/A — provider-free use case is intentionally unwired; deterministic injected fakes prove the public boundary |
| Rollback boundary | Revert/delete the three new `internal/usecase/procesar_mensaje.go`, `siguiente_pregunta.go`, and `procesar_mensaje_test.go` files; no existing ports, domain, contracts, or wiring changed |

## Validation Evidence

- `gofmt -l` on all three files — PASS (no output)
- `go test ./... -count=1` — PASS
- `go test -race ./... -count=1` — PASS
- `go build ./...` — PASS
- `go vet ./...` — PASS
- `go mod tidy` followed by `git diff --exit-code -- go.mod go.sum` — PASS; no module changes
- `git diff --no-index --check` for each new file — PASS

## Scope Notes

PR2-only complex field type/merge edge cases, motor recalculation, unintelligible-audio safety, completion transitions/events, and persistence-failure matrix remain intentionally pending for tasks 2.1 and 2.2. Phase 3 evidence tasks remain unchecked even though this slice's checks were run.
