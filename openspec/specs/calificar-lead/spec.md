# Calificar Lead Specification

## Purpose
Provider-free deterministic qualification of a profiled lead; no Contract, domain, motor, port, or adapter change.

## Requirements

### Requirement: Guard, Capacity, and Recommendations
`CalificarLead` MUST accept only `CALIFICADO`; another state SHALL fail without save or event. It MUST calculate capacity with candidate `0`, select the lowest positive catalog `PrecioDesde` within that preliminary budget, and recalculate; absent selection MUST retain `0` for motor median fallback. It MUST run `GemeloKNN` with current profile, affiliation, dependents, exact catalog-ID zones, buyers, and `K=30`, then use final budget with `RecomendarProyectos`.

#### Scenario: Guard and candidate
- GIVEN a non-`CALIFICADO` lead or unordered catalog prices
- WHEN qualification runs
- THEN the first case has no write/event and the second uses the lowest affordable positive price
- AND no affordable price uses candidate `0` and median fallback

#### Scenario: Exact recommendation evidence
- GIVEN at least 30 buyers and canonical catalog zones
- WHEN identical inputs run repeatedly
- THEN the same 30 distance-ranked neighbors and recommendation order result
- AND absent/non-canonical zones are not substituted

### Requirement: Conversion, Route, Priority, Cupo
For non-affiliates, conversion MUST be true only for `situacion_laboral == INDEPENDIENTE`, `hogar_con_afiliado == true`, or trimmed nonblank `caja_externa`. The use case MUST delegate to `Matriz2x2`, apply Contract §3.5 priority with ratio capped at `1.2`, derive the Contract semáforo, and set `ConsumeCupo10` only for non-affiliate `ASESOR`.

#### Scenario: Predicates and matrix
- GIVEN each conversion signal alone, blank/absent signals, every matrix quadrant, and ratio above `1.2`
- WHEN qualification runs
- THEN only accepted signals enable conversion; route, capped priority, and semáforo match Contract
- AND only non-affiliate `ASESOR` consumes cupo

### Requirement: Durable Route Event
`ASESOR` MUST remain `CALIFICADO`; `NUTRICION`, `REMARKETING`, and `DESPEDIDA` MUST become `EN_NUTRICION`, `REMARKETING`, and `DESPEDIDO`. The lead MUST save before `RutaDecidida`; save failure MUST publish no event.

#### Scenario: Save ordering
- GIVEN each route and a failed save
- WHEN qualification runs
- THEN successful state persistence precedes one event and failure emits none
