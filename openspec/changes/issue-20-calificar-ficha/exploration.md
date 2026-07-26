# Exploration: issue-20-calificar-ficha

## Current State
Issue #20 is the post-conversation qualification boundary. `PerfilarLead` creates leads in `PERFILANDO`; `ProcesarMensaje` merges recognised profile fields, recalculates capacity, transitions complete conversations to `CALIFICADO`, and emits `PerfilCompleto` only after durable writes. No qualification or ficha orchestration exists.

Existing ports are sufficient: `LeadRepository`, `FichaRepository`, `CatalogoRepository`, `GeneradorID`, `Reloj`, and `BusEventos`. The domain already supplies `Lead`, `Ficha`, `Capacidad`, `Intencion`, `Proyecto`, `Comprador`, state transitions, and the deterministic motor functions `CalcularCapacidad`, `GemeloKNN`, `RecomendarProyectos`, and `Matriz2x2`. Ficha persistence is a lead-keyed upsert; lead persistence is optimistic CAS; the pair is not transactional.

## Scope and Decisions
- `CalificarLead` accepts only a `CALIFICADO` lead. It calculates preliminary capacity with candidate `0`, selects the lowest positive catalog `PrecioDesde` within preliminary budget, recalculates capacity with that candidate (the motor falls back to the median for zero), invokes `GemeloKNN` with `K=30` using exact catalog zones, and obtains recommendations through `RecomendarProyectos`.
- Conversion is true only for a non-affiliate with `situacion_laboral == "INDEPENDIENTE"`, `hogar_con_afiliado == true`, or a nonblank `caja_externa`. Routing delegates to `Matriz2x2`; priority follows Contract §3.5 with ratio capped at `1.2`; cupo 10% is set only for a non-affiliate `ASESOR` route.
- An `ASESOR` lead remains `CALIFICADO` for `GenerarFicha`. `NUTRICION`, `REMARKETING`, and `DESPEDIDA` transition to `EN_NUTRICION`, `REMARKETING`, and `DESPEDIDO` respectively. `RutaDecidida` is published only after the lead save succeeds.
- `GenerarFicha` accepts only `CALIFICADO`/`ASESOR`, recomputes deterministic recommendations from catalog and buyers, creates the existing Contract `Ficha`, persists it, then transitions the lead to `ENTREGADO` and saves it. There is no LLM or narrative field/port in scope; all ficha output is deterministic.
- The low-confidence warning must be exactly `PERFIL PARCIALMENTE DECLARADO — validar campos marcados`. Benefits, sales arguments, and alert preserve fixed Contract order.

## Affected Areas
- `internal/usecase/calificar_lead.go` and deterministic qualification tests.
- `internal/usecase/generar_ficha.go` and deterministic ficha tests.
- Hybrid OpenSpec and Engram SDD artifacts.

Explicitly excluded: domain, deterministic motor, ports, infrastructure, adapters, HTTP, frontend, migrations, config, Contract/Wiki documents, LLM wiring, and narrative ports.

## Delivery Strategy
The estimated 500–700 authored code/test lines exceed the hard 400-line review budget. Deliver two chained, rollbackable slices without a size exception: qualification first; ficha generation and its error/retry paths second.

## Risks
- Separate ficha and lead writes can leave a persisted ficha while the lead stays `CALIFICADO`; retry is repairable via ficha upsert and no event is emitted before lead durability.
- Candidate price must be the lowest positive affordable catalog price; median is fallback only.
- Recommendation reconstruction must use identical catalog zones, buyers, and `K=30` in both use cases.
- Conversion must include the three accepted signals and ignore blank `caja_externa`.

## Recommendation
Proceed with two focused, provider-free application use cases and preserve all existing boundaries. No user decision is required before proposal.
