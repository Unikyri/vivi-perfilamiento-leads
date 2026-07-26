# Design: Provider-free lead qualification and deterministic ficha

## Technical Approach

Add two synchronous application services in `internal/usecase`, backed only by existing repositories, clock/ID/event ports, and pure motor functions. A package-private shared decision helper reconstructs capacity, KNN evidence, and recommendations identically for both use cases. No provider, transaction, or new boundary is introduced.

## Architecture Decisions

| Decision | Choice and rationale | Rejected |
|---|---|---|
| Boundary | `CalificarLead` and `GenerarFicha` orchestrate existing domain/motor/ports; helpers remain package-private. | Coordinator, event subscription, LLM, or port/domain changes. |
| Shared decision | One helper loads `Proyectos`, computes preliminary capacity with candidate `0`, selects the lowest positive affordable `PrecioDesde`, recomputes final capacity, loads `Compradores`, builds zones only from exact `map key == ProyectoID`, derives dependents as `personas_hogar-1` when present, calls `GemeloKNN(K=30)`, then `RecomendarProyectos`. The motor owns income-category projection. | Duplicated formulas, `BandaDeCategoria`, median-first selection, or mutable catalog slices. |
| Persistence | Lead CAS precedes `RutaDecidida`; ficha upsert precedes `CALIFICADO→ENTREGADO`. | Publishing before durability or pretending the two repositories are transactional. |
| Retry | Existing partial ficha supplies stable ID/time; the current deterministic decision is rebuilt and upserted before completing the lead transition. Successful terminal calls reject. A previously qualified `CALIFICADO/ASESOR` call returns reconstructed output without another save/event. | Duplicate ficha rows/metadata/events or compensating deletes. |

## Exact Operation Ordering

```mermaid
sequenceDiagram
  Caller->>CalificarLead: Ejecutar(leadID)
  CalificarLead->>LeadRepository: PorID; guard CALIFICADO + intention
  CalificarLead->>CatalogoRepository: Proyectos; Compradores
  CalificarLead->>Motor: capacity(0); candidate; capacity(candidate); KNN; recommend; matrix
  CalificarLead->>LeadRepository: Guardar(CAS)
  CalificarLead->>BusEventos: RutaDecidida
  Caller->>GenerarFicha: Ejecutar(leadID)
  GenerarFicha->>LeadRepository: PorID; guard CALIFICADO/ASESOR
  GenerarFicha->>FichaRepository: PorLead (repair probe)
  GenerarFicha->>CatalogoRepository: decision helper; reuse existing ID/time on repair
  GenerarFicha->>FichaRepository: Guardar(upsert)
  GenerarFicha->>LeadRepository: Guardar(ENTREGADO, CAS)
```

`CalificarLead` computes conversion only for non-affiliates using exact `INDEPENDIENTE`, `hogar_con_afiliado`, or trimmed nonblank `caja_externa`; calls `Matriz2x2`; sets priority as route weight (`1/.5/.25/.1`) × `min(ratio,1.2)` × confidence; maps semáforo `ASESOR→VERDE`, `NUTRICION→AMBAR`, others→`GRIS`; and sets cupo only for non-affiliate `ASESOR`. State mapping is `ASESOR→CALIFICADO` (unchanged), `NUTRICION→EN_NUTRICION`, `REMARKETING→REMARKETING`, `DESPEDIDA→DESPEDIDO`. The published event payload contains `ruta`, `prioridad`, `semaforo`, `consume_cupo_10`, and recommendations.

`GenerarFicha` copies identification, final capacity, profile, intention, recommendations, and cupo. Below confidence `0.6` it sets exactly `PERFIL PARCIALMENTE DECLARADO — validar campos marcados`. Benefits append `Subsidio de caja 30 SMMLV` when subsidy is positive, then for affiliates `Crédito propio Colsubsidio` and `Acompañamiento PerteneSer`. If rent is positive, the sole argument is `Paga $<arriendo> de arriendo; la cuota estimada es $<cuota>`, using deterministic dotted-COP formatting and `cuota=40%` of income. Alert rate is desistidos/all KNN neighbors (`0` when empty), active only above `0.20`, with nil detail. ID/time come only from existing ports.

## File Changes

| File | Action | Responsibility |
|---|---|---|
| `internal/usecase/calificar_lead.go` | Create | Service, shared decision/candidate helpers, route/priority/semaphore. |
| `internal/usecase/calificar_lead_test.go` | Create | Matrix, boundaries, ordering, event/CAS tests. |
| `internal/usecase/generar_ficha.go` | Create | Deterministic content and partial-write repair. |
| `internal/usecase/generar_ficha_test.go` | Create | Content order, errors, upsert-repair, aliasing tests. |

No implementation file is modified or deleted.

## Errors, Tests, and Durability

Blank IDs, absent intention/capacity, wrong state/route, context cancellation, and repository errors return before downstream writes; errors wrap existing validation/not-found/CAS/domain errors where applicable. The bus has no acknowledgement/outbox: lead state is durable before publication, but event delivery is best-effort and not exactly-once.

Use fixed clock/IDs, existing lead fake, and narrow catalog/ficha/bus fail-at-step fakes. Assert call order, no writes on guards/read/motor failures, all route/conversion/candidate boundaries, cross-use-case recommendation equality, immutable inputs, exact strings/order, partial-write retry, focused tests, then `go test ./...`, `go vet ./...`, and `go build ./...`.

## Threat Matrix, Delivery, and Rollback

All five threat-matrix rows are N/A: this is business route selection, with no path execution, Git selection, commit/push state, PR command composition, shell, subprocess, or process integration.

Slice 1 is qualification plus tests; slice 2 (targeting slice 1) is ficha plus error/repair tests; each must remain ≤400 authored lines. Roll back either additive slice independently. Confirmed unchanged: domain, motor, ports, infrastructure, adapters, HTTP, frontend, migrations, config, Contract, and Wiki.

## Open Questions

None.
