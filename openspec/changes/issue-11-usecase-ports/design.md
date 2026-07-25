# Design: Usecase Ports and In-Memory Lead Fake — Issue #11

## Technical Approach

Add application-owned ports in `internal/usecase` so use cases depend only on standard-library types and `internal/domain`. Preserve the Spanish, pointer-based vocabulary already used by dependent issue sketches, while treating Contract v1.1, NFR-M-01/M-04, and Architecture §3 as authoritative. Production ports are separated from test-only fakes; no adapter, domain, data, or Wiki change is required.

## Architecture Decisions

| Decision | Alternatives considered | Choice and rationale |
|---|---|---|
| Port ownership | Adapter-owned interfaces; generic repository | Define explicit interfaces in `usecase`; this enforces inward dependencies and preserves resource-specific behavior. |
| Compatibility surface | English/value-return API | Retain `Crear`/`PorID`/`Guardar`/`Listar`, pointer repository entities, `FiltroLeads`, `LLMProvider.Nombre`, `Reloj`, `BusEventos`, and `MensajeriaGateway`; consumers already share this vocabulary and authority does not conflict. |
| Lead concurrency | Unconditional increment; last-write-wins | `Guardar` performs atomic CAS against `Lead.Version`; stale writes return `ErrOptimisticLock` without changing storage or caller. Success increments stored and caller version exactly once. |
| Fake isolation | Shallow copies; JSON round-trip | Clone every mutable boundary recursively. Explicit domain cloning plus recursive map/slice/pointer cloning preserves concrete numeric types and prevents aliases. |
| Deferred behavior | Fake every declared port; add plan version | Implement only the lead fake and minimal LLM/clock/ID doubles. Declare Plan/Ficha/Catalog ports, but defer their fakes and plan CAS until a consuming scenario defines them. |

## Data Flow

```mermaid
sequenceDiagram
    participant UC as Usecase/Test
    participant F as LeadRepoFake
    participant S as Private cloned storage
    UC->>F: PorID(id)
    F->>S: read under RLock
    F-->>UC: recursive clone
    UC->>F: Guardar(lead, expected Version)
    F->>S: lock + compare version
    alt stale or missing
        F-->>UC: ErrOptimisticLock / ErrNoEncontrado
    else match
        F->>S: clone, increment once, replace
        F-->>UC: nil; caller Version committed
    end
```

`Listar` clones matching leads, applies non-nil `Afiliado` and `Ruta` filters conjunctively, then sorts by `Prioridad` descending and `LeadID` ascending. It returns a non-nil empty slice when unmatched. Conversation results are cloned and chronological.

## File Changes

| File | Action | Review slice | Description |
|---|---|---:|---|
| `internal/usecase/puertos.go` | Create | 1 | Errors, DTOs, filters, event constants, and all port interfaces. |
| `internal/usecase/puertos_test.go` | Create | 1 | Compile-time method-shape and Contract §7 DTO tests. |
| `internal/usecase/fakes_test.go` | Create | 2 | Mutex-protected lead fake, recursive clone helpers, and minimal LLM/clock/ID doubles. |
| `internal/usecase/fakes_behavior_test.go` | Create | 2 | Table-driven CAS, absence, isolation, filtering, ordering, message, and interface tests. |

No files are modified or deleted.

## Interfaces / Contracts

`LeadRepository` exposes `Crear`, `PorID`, `Guardar`, `Listar`, `AgregarMensaje`, and `Conversacion`; `FiltroLeads` contains optional `Afiliado` and `Ruta`. Repository absence is matched with `errors.Is(err, ErrNoEncontrado)`; stale `Guardar` is matched with `ErrOptimisticLock`. `Crear` normalizes version 0 to 1, accepts 1, rejects duplicates/invalid initial versions, and stores a clone.

`PlanRepository` exposes create/by-lead/save/overdue-milestone/mark-milestone operations without CAS. `FichaRepository` exposes save/by-lead. `CatalogoRepository` exposes projects, buyers, affiliate lookup, and brochure Markdown. DTOs `Afiliado`, `EntradaTurno`, `CampoExtraido`, `SalidaTurno`, `Audio`, `Evento`, and `HitoConPlan` follow Contract §§0/4.3/6/7. `LLMProvider` includes text, audio, and `Nombre`; `MensajeriaGateway` sends a `*domain.Mensaje`; `Reloj`, `BusEventos`, and `GeneradorID` retain the issue shapes.

## Testing Strategy

| Layer | Proof | Approach |
|---|---|---|
| Shape | Ports and DTOs compile without outer-layer imports | Compile assertions, DTO JSON tests, `go list -deps`. |
| Fake behavior | Atomic CAS, unchanged stale state, explicit absence, recursive isolation | Table-driven tests; mutate inputs/outputs and re-read. |
| Collections | AND filters, stable ties, non-nil empty list, chronological cloned messages | Shuffled insertion and repeated-result tests. |
| Package | No regressions | `go test ./internal/usecase/... -race`, `go vet ./internal/usecase/...`, then `go test ./...`. |

## Threat Matrix

N/A — no routing, shell, subprocess, VCS/PR automation, executable-file classification, or process-integration boundary.

## Migration / Rollout

No migration or feature flag. Deliver two sequential review slices: ports plus shape tests, then lead fake plus behavior tests. Each slice is independently revertible; slice 2 targets slice 1.

## Open Questions

None.
