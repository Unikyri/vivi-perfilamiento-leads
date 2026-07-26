# Tasks: Issue 30 — Buyer Clock

## Review Workload Forecast

| Field | Value |
|---|---|
| Estimated changed lines | 120–180 authored lines |
| 400-line budget risk | Low |
| Chained PRs recommended | No |
| Suggested split | Single PR |
| Delivery strategy | single-pr |
| Chain strategy | pending (not needed) |

Decision needed before apply: No
Chained PRs recommended: No
Chain strategy: pending
400-line budget risk: Low

### Suggested Work Units

| Unit | Goal | Likely PR | Focused test command | Runtime harness | Rollback boundary |
|---|---|---|---|---|---|
| 1 | Separate adapter wall/demo clocks; prove HTTP contracts and retain buyers. | Single PR | `go test ./internal/infrastructure/reloj/... ./internal/adapters/http/...` | `go test ./internal/adapters/http -run 'TestDemoTiempo|TestGerenciaBuyerPersona'` uses `httptest`; no external service. | Revert adapter, tests, fallback comment, and decision record; no schema/data/API rollback. |

## Phase 1: Characterization Tests

- [x] 1.1 Update `internal/infrastructure/reloj/postgres_test.go` with table-driven RED coverage: loaded UTC date, zero-value fallback saved once, forward advance changes only `FechaSimulada()`, and `Ahora()` is bounded by observed real UTC wall time.
- [x] 1.2 Extend `internal/adapters/http/gerencia_test.go` to decode and assert deterministic filtered and catalog-wide buyer-persona response shapes, plus `VALIDACION` for empty or repeated `proyecto_id`; do not change aggregation behavior.
- [x] 1.3 Retain/verify `internal/adapters/http/demo_tiempo_test.go` coverage that demo advance persists once and exposes the simulated date; add a focused assertion only if needed to distinguish it from operational wall time.

## Phase 2: Adapter Clock Split

- [x] 2.1 Modify `internal/infrastructure/reloj/postgres.go`: keep mutex-protected persisted simulated UTC state for `FechaSimulada()` and `Avanzar`, while `Ahora()` returns `time.Now().UTC()` without changing `usecase.Reloj` or callers.
- [x] 2.2 Add a comment in `internal/adapters/http/rutas.go` documenting `relojSistema` as non-persistent real wall time whose `Avanzar` is intentionally a no-op; do not alter fallback behavior or routes.

## Phase 3: Preservation and Decision Record

- [x] 3.1 Create `docs/decisiones/0001-conservar-tabla-compradores.md` in Spanish: retain Contract-declared `compradores` for JSON loading, kNN, ficha, and buyer-persona; state reset preserves it and removal requires a separate Contract §9 PR with both-block approval.
- [x] 3.2 Confirm no changes to `migrations/001_esquema_inicial.sql`, `data/compradores.json`, `internal/usecase/buyer_persona.go`, buyer endpoints, or the `usecase.Reloj` interface.

## Phase 4: Validation

- [x] 4.1 Run `go test ./internal/infrastructure/reloj/... ./internal/adapters/http/...` and ensure split-clock, `/api/demo/tiempo`, and buyer-persona endpoint scenarios pass.
- [x] 4.2 Run `go test ./...`, `go vet ./...`, and `go build ./...`; record that no buyer data, migration, Contract, or API response shape changed.
