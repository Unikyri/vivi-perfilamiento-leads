# Design: Buyer Persona, Production Clock, and Buyers Table Decision

## Technical Approach

Change only the behavior of `reloj.Postgres`: retain its mutex-protected persisted simulated value, make `Ahora()` call `time.Now().UTC()`, leave `FechaSimulada()` on stored simulated time, and let `Avanzar` update only that stored value. `usecase.Reloj`, all callers, and existing fake clocks remain unchanged.

## Architecture Decisions

| Decision | Choice | Rationale |
|---|---|---|
| Clock split | Adapter-only change | Every use case already depends on the port, so audit timestamps are corrected without call-site churn. |
| Demo source | Keep `demo.fecha_simulada` | Retains restart, reset, non-regression, milestones, and health behavior. |
| Buyers table | Retain and document | Contract-declared JSON/port/kNN/buyer-persona boundary; removal needs separate Contract §9 PR and both-block approval. |

## Data Flow

```text
operational use case ──Reloj.Ahora()────────→ real UTC record timestamp
/demo/tiempo ──AvanzarDemo──→ DemoRepository + Reloj.Avanzar → simulated date
/salud, milestones ──FechaSimulada()────────→ stored simulated date
```

The initialization fallback remains: a zero repository value is assigned real UTC once and persisted. `AvanzarDemo` remains the authority for persistence and backward rejection; `Postgres.Avanzar` does not write the repository. The `rutas.go` `relojSistema` fallback receives a comment stating it is non-persistent wall time with no-op advance.

## File Changes

| File | Action | Description |
|---|---|---|
| `internal/infrastructure/reloj/postgres.go` | Modify | Split wall `Ahora` from stored `FechaSimulada`. |
| `internal/infrastructure/reloj/postgres_test.go` | Modify | Table-driven load/fallback/advance/split tests. |
| `internal/adapters/http/rutas.go` | Modify | Comment-only non-persistent fallback clarification. |
| `internal/adapters/http/gerencia_test.go` | Modify | Preserve filtered/catalog/invalid buyer-persona contract checks. |
| `docs/decisiones/0001-conservar-tabla-compradores.md` | Create | Spanish retention decision and Contract-PR removal rule. |

## Interfaces / Contracts

`usecase.Reloj` is unchanged: `Ahora`, `FechaSimulada`, `Avanzar`. No HTTP shape, data file, migration, or buyer-persona calculation changes.

## Testing Strategy

| Layer | Proof |
|---|---|
| Adapter unit | Loaded UTC date, one fallback save, advance changes only simulation, `Ahora` remains within wall-clock observation bounds. |
| Use-case/HTTP | Backward `AvanzarDemo` rejection; health exposes simulation; existing fake continues to work. |
| HTTP | Buyer-persona filtered/catalog shapes and invalid repeated/empty `proyecto_id`. |
| Commands | `go test ./internal/infrastructure/reloj/... ./internal/adapters/http/...`, then `go test ./...`, `go vet ./...`, `go build ./...`. |

## Threat Matrix

N/A — no routing, shell, subprocess, VCS/PR automation, executable-file classification, or process-integration boundary.

## Migration / Rollout

No migration, data conversion, Contract change, or call-site change. Revert the adapter/documentation change to roll back.
