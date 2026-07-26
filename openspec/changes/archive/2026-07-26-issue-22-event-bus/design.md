# Design: In-Memory Event Bus and Deterministic Coordinator

## Technical Approach
Add a synchronous Observer in `internal/infrastructure/bus` implementing the existing `usecase.BusEventos`, plus a plain-Go Mediator in `internal/adapters/agentes`. The coordinator delegates through narrow interfaces already satisfied by `*usecase.CalificarLead`, `*usecase.GenerarFicha`, and `*usecase.EjecutarHitos`; no existing use case, port, domain type, or composition root changes.

## Architecture Decisions
| Concern | Choice and rationale | Rejected alternative |
|---|---|---|
| Placement | Infrastructure owns delivery; adapter owns orchestration, preserving Clean Architecture. | Bus/coordinator in `usecase`, which mixes mechanisms with application policy. |
| Dispatch | Snapshot matching handlers under `sync.RWMutex`, unlock, then invoke synchronously in registration order. This permits nested publish/subscribe and guarantees order per publication, not across concurrent publishers. | Lock during callbacks (deadlock risk), goroutines/queue (different ordering and durability contract). |
| Isolation | Clone at `Publicar` entry and clone again per handler with a cycle-aware reflection helper preserving concrete scalar, map, slice, array, pointer, interface, and struct types. | Shallow copy (nested aliases) or JSON round-trip (type loss). |
| Failure/privacy | Recover each callback independently; continue later handlers. Inject `*slog.Logger`; emit only `tipo`, opaque `lead_id`, handler index/name, and fixed outcome (`ok`, `skipped`, `error`, `panic`). Never emit payload or raw error/panic text. | Propagation is impossible through the void bus port; raw diagnostics can leak PII. |
| Policy | One literal registration table, no LLM/ADK judgment. `Registrar` uses `sync.Once`; nil bus/dependencies skip registration safely. | Dynamic graph or direct agent-to-agent calls. |

## Interfaces / Contracts
```go
// package bus
func Nuevo(logger *slog.Logger) *EnMemoria
func (*EnMemoria) Suscribir(string, func(context.Context, usecase.Evento))
func (*EnMemoria) Publicar(context.Context, usecase.Evento)

// package agentes
type ManejadorEvento func(context.Context, usecase.Evento) error
type Calificador interface {
    Ejecutar(context.Context, usecase.EntradaCalificar) (usecase.SalidaCalificar, error)
}
type Documentadora interface {
    Ejecutar(context.Context, string) (*domain.Ficha, error)
}
type Nutricionista interface {
    Ejecutar(context.Context, time.Time) (int, error)
}
type Dependencias struct {
    LeadNuevo ManejadorEvento
    Calificador Calificador
    Documentadora Documentadora
    Nutricionista Nutricionista
    Logger *slog.Logger
}
func Nueva(bus usecase.BusEventos, deps Dependencias) *Coordinadora
func (*Coordinadora) Registrar()
```
`Suscribir` ignores nil handlers. A nil logger means silent operation. Context is forwarded unchanged. Handler errors are classified and swallowed because `BusEventos.Publicar` has no error result.

## Coordinator Dispatch Policy
| Event | Action |
|---|---|
| `LeadNuevo` | Invoke only the optional callback; never recreate the lead or fabricate a greeting. |
| `PerfilCompleto` | Call `Calificador.Ejecutar(ctx, EntradaCalificar{LeadID: e.LeadID})`; `errors.Is(err, usecase.ErrLeadNoCalificable)` is a side-effect-free skip for nutrition reprofile. Other errors log fixed `error`. |
| `RutaDecidida` | Accept `payload["ruta"]` as `domain.Ruta` or string. Only `ASESOR` calls `Documentadora.Ejecutar(ctx, e.LeadID)`; missing/unknown/non-ASESOR routes skip. Never republish or create a nutrition plan. |
| `TickReloj` | Accept `payload["hasta"]` as `time.Time` or RFC3339 string, then call `Nutricionista.Ejecutar`; malformed values skip. |
| Other event constants | No subscription: no owned consumer exists. |

```mermaid
sequenceDiagram
  Producer->>EnMemoria: Publicar(ctx, Evento)
  EnMemoria->>EnMemoria: clone; snapshot under RLock; unlock
  loop registration order
    EnMemoria->>Coordinadora: callback(ctx, per-handler clone)
    Coordinadora->>UseCase: deterministic typed call with lead_id
    UseCase-->>Coordinadora: result/error
  end
```

## File Changes
| File | Action | Description |
|---|---|---|
| `internal/infrastructure/bus/memoria.go` | Create | Bus, locking, cloning, recovery, safe logs. |
| `internal/infrastructure/bus/memoria_test.go` | Create | Focused bus tests. |
| `internal/adapters/agentes/agentes.go` | Create | Narrow interfaces and dependencies. |
| `internal/adapters/agentes/coordinadora.go` | Create | Registration table and handlers. |
| `internal/adapters/agentes/coordinadora_test.go` | Create | Routing/error/privacy tests. |

## Testing Strategy
Use table-driven Go tests. Slice 1 proves sync order, producer/per-handler isolation including `[]domain.Recomendacion`, nested publication, nil/no subscriber, panic continuation, context identity, and concurrent registration/publication under `go test -race`. Slice 2 proves ten-run determinism, nil-safe/idempotent registration, normal qualification, `ErrLeadNoCalificable` skip, no duplicate route publication, ASESOR-only ficha, route/time parsing, no plan creation, fixed log fields, and absence of payload/PII/raw errors. Run focused package tests, `go test -race` for both packages, then `go test ./...` and `go build ./...`. Runtime harness: N/A—no process or transport boundary.

## Threat Matrix
All prescribed rows—documentation-like paths, Git repository selection, commit state, push state, and PR commands—are **N/A**: this routes in-memory typed events only and executes no files, shell, VCS, PR, or subprocess operations.

## Delivery / Rollback
1. **Bus work unit:** `memoria.go` + tests, forecast 340 lines; independently usable; rollback both files.
2. **Coordinator work unit:** three `agentes` files, forecast 340 lines; depends only on slice 1/existing ports; rollback them while retaining the bus.

Both auto-chained slices stay under 400 authored lines. No migration or feature flag. Scope excludes ADK/LLM, HTTP/composition wiring, greeting, automatic plans/pause routing, outbox/queue/WebSocket, ports/domain/repositories/migrations, and existing use-case edits. Open questions: none.
