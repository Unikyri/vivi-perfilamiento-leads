# Tasks: In-Memory Event Bus and Deterministic Coordinator

## Review Workload Forecast

| Field | Value |
|---|---|
| Estimated changed lines | Slice 1: ~340; Slice 2: ~340; total ~680 authored runtime/test lines |
| Delivery strategy | auto-chain |
| Chain strategy | stacked-to-main: PR 1 bus, then PR 2 coordinator after PR 1 merges |
| 400-line budget risk | High — the aggregate crosses 400, but each independently reviewable slice is capped below 400 |

Decision needed before apply: No
Chained PRs recommended: Yes
Chain strategy: stacked-to-main
400-line budget risk: High

### Allowed implementation files (only)
- `internal/infrastructure/bus/memoria.go`
- `internal/infrastructure/bus/memoria_test.go`
- `internal/adapters/agentes/agentes.go`
- `internal/adapters/agentes/coordinadora.go`
- `internal/adapters/agentes/coordinadora_test.go`

## Slice 1 — Bus (PR 1, ≤340 authored runtime/test lines)

- [x] 1.1 Create `bus/memoria.go`: `EnMemoria`, synchronous registration-ordered snapshot dispatch under `RWMutex`, unlock before callbacks, nil-handler ignore, and unchanged context.
- [x] 1.2 Add producer and per-subscriber cycle-aware event cloning plus per-handler `recover`; emit only fixed safe `slog` metadata (`tipo`, opaque `lead_id`, handler, outcome).
- [x] 1.3 Create table-driven `memoria_test.go`: standalone/order, unknown/zero subscribers, nil payload/handler, nested publish/subscribe without deadlock, and context identity.
- [x] 1.4 Test producer/per-handler isolation with nested maps and `[]domain.Recomendacion`, panic continuation, and that sensitive payload/panic text never reaches logs.
- [x] 1.5 Focused validation: `go test ./internal/infrastructure/bus -run TestEnMemoria -count=1 && go test -race ./internal/infrastructure/bus`; full: `go test ./... && go build ./...`.

## Slice 2 — Coordinator (PR 2, ≤340 authored runtime/test lines)

- [x] 2.1 Create `agentes.go` with only the design’s narrow `Calificador`, `Documentadora`, `Nutricionista`, `Dependencias`, and event-handler contracts.
- [x] 2.2 Create `coordinadora.go`: nil-safe, idempotent `Registrar` table for optional `LeadNuevo`, `PerfilCompleto`, `RutaDecidida`, and `TickReloj`; subscribe no other events.
- [x] 2.3 Implement typed handlers: `ErrLeadNoCalificable` skip; `ASESOR`-only ficha from `domain.Ruta`/string; no republished route or plan; `hasta` as `time.Time`/RFC3339; malformed input and nonfatal errors safe-log then no-op.
- [x] 2.4 Create table-driven `coordinadora_test.go`: ten-run deterministic routing, normal qualification/no duplicate route, reprofile skip, ASESOR vs NUTRICION, route/tick parsing, nil dependencies/payload, and metadata-only logs.
- [x] 2.5 Focused validation: `go test ./internal/adapters/agentes -run TestCoordinadora -count=1 && go test -race ./internal/adapters/agentes`; full: `go test ./... && go build ./...`.

## Evidence, Review, and Delivery

- Runtime harness: N/A for both slices—scope has no HTTP, ADK, process, or composition-root boundary; focused deterministic in-memory tests are the runtime evidence.
- Native review gate per slice: after validation run `gentle-ai review start`, execute/capture its selected lenses, `gentle-ai review finalize`, then `gentle-ai review validate --gate post-apply --cwd /tmp/vivi-issue-22`; revalidate the same receipt at pre-commit, pre-push, and pre-PR without creating a new budget.
- Commit approach: include code and tests together as `feat(bus): add synchronous in-memory event bus` and `feat(agentes): add deterministic event coordinator`; record focused/full results before each PR.
- Rollback: revert Slice 1’s two bus files; revert Slice 2’s three coordinator files while retaining the unused bus. Do not modify ports, domain, existing use cases, composition, state, or runtime artifacts.
