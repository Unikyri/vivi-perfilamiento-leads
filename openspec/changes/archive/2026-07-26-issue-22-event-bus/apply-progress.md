# Apply Progress: In-Memory Event Bus

**Change**: `issue-22-event-bus`
**Mode**: Standard (`strict_tdd: false`)
**Delivery**: Auto-chain, stacked-to-main; Slice 2 / PR 2

## Completed Tasks

- [x] 1.1 Created `internal/infrastructure/bus/memoria.go` with synchronous registration-ordered snapshot dispatch under `RWMutex`, unlocked before callbacks, nil-handler ignore, and unchanged context forwarding.
- [x] 1.2 Added producer/per-subscriber cycle-aware reflection cloning, per-handler panic recovery, and fixed privacy-safe `slog` metadata (`tipo`, opaque `lead_id`, handler index, outcome).
- [x] 1.3 Added table-driven tests for standalone delivery/order, unknown and zero subscribers, nil handler, nested publish/subscribe without deadlock, and context identity.
- [x] 1.4 Added tests for nested map and `[]domain.Recomendacion` isolation, panic continuation, and absence of sensitive payload/panic text from logs.
- [x] 1.5 Completed focused, package, race, full-suite, build, vet, formatting, and diff checks.

## Work Unit Evidence

| Evidence | Result |
|---|---|
| Focused test | `go test ./internal/infrastructure/bus -run TestEnMemoria -count=1` — PASS |
| Package test | `go test ./internal/infrastructure/bus` — PASS |
| Race test | `go test -race ./internal/infrastructure/bus` — PASS |
| Full test | `go test ./...` — PASS |
| Build | `go build ./...` — PASS |
| Vet | `go vet ./...` — PASS |
| Formatting/diff | `gofmt -d` empty and `git diff --check` — PASS |
| Runtime harness | N/A — this slice has no HTTP, ADK, process, or composition-root boundary; deterministic bus tests are the runtime evidence. |
| Rollback boundary | Revert exactly `internal/infrastructure/bus/memoria.go` and `internal/infrastructure/bus/memoria_test.go`; no unrelated behavior is removed. |

## Scope and Budget

- Runtime/test files changed: `internal/infrastructure/bus/memoria.go`, `internal/infrastructure/bus/memoria_test.go`.
- Authored runtime/test lines: 308 total (160 implementation + 148 tests), below the Slice 1 forecast of 340 and the 400-line hard budget.
- No ports, domain, usecase, adapters, migrations, configuration, or composition-root files changed.

## Slice 2 — Coordinator (PR 2)

- [x] 2.1 Created `internal/adapters/agentes/agentes.go` with only the narrow typed contracts and optional dependencies.
- [x] 2.2 Created `coordinadora.go` with a fixed, idempotent `sync.Once` registration table for the four supported event types.
- [x] 2.3 Implemented deterministic qualification, advisor-only ficha, tick parsing, reprofile skip, no route republish, no automatic nutrition plan, and safe malformed/error handling.
- [x] 2.4 Added table-driven tests for ten-run routing, normal qualification, duplicate-route prevention, reprofile skip, advisor/nutrition routing, route/tick parsing, nil safety, and metadata-only logs.
- [x] 2.5 Completed focused, race, full-suite, build, vet, formatting, diff, and budget checks.

## Work Unit Evidence — Slice 2

| Evidence | Result |
|---|---|
| Focused test | `go test ./internal/adapters/agentes -run TestCoordinadora -count=1` — PASS |
| Race test | `go test -race ./internal/adapters/agentes` — PASS |
| Full test | `go test ./...` — PASS |
| Build | `go build ./...` — PASS |
| Vet | `go vet ./...` — PASS |
| Formatting/diff | `gofmt -d` empty and `git diff --check` — PASS |
| Runtime harness | N/A — no HTTP, ADK, process, or composition-root boundary; deterministic in-memory tests are the runtime evidence. |
| Rollback boundary | Revert exactly `internal/adapters/agentes/agentes.go`, `coordinadora.go`, and `coordinadora_test.go`; retain Slice 1 bus files. |

## Scope and Budget

- Runtime/test files changed in Slice 2: `internal/adapters/agentes/agentes.go`, `coordinadora.go`, `coordinadora_test.go`.
- Authored runtime/test lines: 295 total, below the Slice 2 forecast of 340 and the 400-line hard budget.
- No ports, domain, existing use cases, migrations, HTTP, configuration, composition-root, or Slice 1 files changed.

## Deviations

None — implementation follows the design and assigned Slice 2 scope.

## Remaining Tasks

None. All tasks 1.1–1.5 and 2.1–2.5 are complete.

**Status**: Slice 2 complete; ready for `sdd-verify`.
