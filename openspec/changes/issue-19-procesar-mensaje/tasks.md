# Tasks: ProcesarMensaje conversational turn (Issue #19)

## Review Workload Forecast

| File | Likely + | Likely - | PR 1 / PR 2 |
|---|---:|---:|---:|
| `internal/usecase/procesar_mensaje.go` | 265 | 0 | 145 / 120 |
| `internal/usecase/siguiente_pregunta.go` | 42 | 0 | 42 / 0 |
| `internal/usecase/procesar_mensaje_test.go` | 278 | 0 | 130 / 148 |
| **Total** | **585** | **0** | **317 / 268** |

Delivery strategy: feature-branch-chain; hard cap 400 per PR, no size exception. PR 1 base: `feat/issue-19-procesar-mensaje`; PR 2 base: PR 1 branch.

Decision needed before apply: No
Chained PRs recommended: Yes
Chain strategy: feature-branch-chain
400-line budget risk: High

### Suggested Work Units

| Unit | Goal and files | Focused test | Runtime / rollback |
|---|---|---|---|
| PR 1 core (317) | Helpers, validated text turn, one-call dispatch; all three files. | `go test ./internal/usecase -run 'TestProcesarMensajeCore|TestSiguientePregunta' -count=1` | N/A: provider-free usecase is intentionally unwired; fixed fakes prove the public boundary. Revert PR 1 removes all new files. |
| PR 2 edges (268) | Merge/motor, audio, caps, ordered failures; `procesar_mensaje.go` and `_test.go`. | `go test ./internal/usecase -run 'TestProcesarMensajeEdges|TestAplicarCampos' -count=1` | N/A: no HTTP/ADK/runtime path is in scope. Revert PR 2 removes only edge branches/tests, preserving PR 1 core. |

## Phase 1: Core boundary — PR 1

- [x] 1.1 In `internal/usecase/procesar_mensaje_test.go`, add fixed-clock/ID scripted-provider RED tables for Invalid entry, Provider error, Audio transcript, and Verified field scenarios (zero writes/call counts and caller text assertions).
- [x] 1.2 Create `internal/usecase/siguiente_pregunta.go` with exact priority and completion trio; prove `SiguienteMejorPregunta` skips `VERIFICADO_BASE` (Question helpers / Verified field).
- [x] 1.3 Create `internal/usecase/procesar_mensaje.go` DTO, sentinels, validation, bounded history, text dispatch, and injected dependency sequence; make 1.1 pass with one provider call and no external I/O (Input validation, One modality dispatch, Provider isolation).

## Phase 2: Integrity and completion — PR 2

- [ ] 2.1 Extend `internal/usecase/procesar_mensaje_test.go` with RED tables for protected merge/motor, unintelligible audio, 4/6 caps, completion de-duplication, and CAS/response failure prefixes.
- [ ] 2.2 Complete `internal/usecase/procesar_mensaje.go`: two-pass field validation/merge, one motor refresh, audio-safe handling, ordered LEAD→CAS→VIVI→event writes, cap completion and terminal-action precedence; make 2.1 pass (Protected fields, Unintelligible audio, Response failure, Duplicate completion).

## Phase 3: Evidence and budget gates

- [ ] 3.1 Per PR: run its focused command, `gofmt -w` then `gofmt -l`, `go vet ./...`, `git diff --check`, and `git diff --numstat <slice-base>...HEAD` summed additions+deletions ≤400; record exact outcomes.
- [ ] 3.2 Before final review: run `go test ./... -count=1`, `go test -race ./... -count=1`, `go build ./...`, `go vet ./...`, `go mod tidy && git diff --exit-code -- go.mod go.sum`, and final format/diff checks. Threat matrix runtime cases are N/A: no shell, routing, subprocess, VCS, or process integration.

Rollback boundary: no migration, configuration, port, DTO, Contract, or runtime wiring changes; revert the applicable chained PR. Do not edit outside the three named files.
