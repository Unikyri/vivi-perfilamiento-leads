# Generar Ficha Specification

## Purpose
Provider-free deterministic Contract `Ficha` generation after advisor routing; no LLM, narrative field, transaction, or new port.

## Requirements

### Requirement: Eligibility and Recommendation Reconstruction
`GenerarFicha` MUST accept only `CALIFICADO` plus `ASESOR`; all other state/routes SHALL fail without ficha or lead writes. It MUST reconstruct recommendations from current catalog/buyers, exact catalog-ID zones, profile, affiliation, dependents, `GemeloKNN(K=30)`, and final-budget `RecomendarProyectos`. It MUST NOT call an LLM.

#### Scenario: Guard and parity
- GIVEN an ineligible lead or unchanged eligible inputs
- WHEN generation runs
- THEN the ineligible case persists nothing
- AND the eligible set equals qualification's ordered recommendations

### Requirement: Ordered Contract Content
The ficha MUST populate declared Contract order: identity/time; confidence/warning; identification; capacity; profile; intention; recommendations; benefits; sales arguments; withdrawal alert; cupo. `Beneficios` and `ArgumentosVenta` MUST retain fixed Contract order. Low confidence MUST use exactly `PERFIL PARCIALMENTE DECLARADO — validar campos marcados`. Positive qualifying `arriendo_actual` MUST add the fixed rental argument at its Contract position; missing/nonpositive rent MUST not. Alert MUST activate only when the neighbor withdrawal rate is strictly greater than `0.20`; a rate equal to or below `0.20` MUST remain inactive with deterministic rate/detail.

#### Scenario: Stable content and threshold
- GIVEN repeated low-confidence, positive-rent inputs and evidence at/below or above `0.20`
- WHEN a ficha is generated
- THEN all fields, ordered slices, warning, rental argument, and alert are byte-stable
- AND only above-threshold evidence activates alert

### Requirement: Ficha-First Handoff and Retry
The ficha MUST save before lead transition to `ENTREGADO` and lead save. Ficha-save failure MUST leave the lead unchanged. Later lead-save failure MUST leave it `CALIFICADO`/`ASESOR`; retry MUST upsert the ficha then complete without duplication. No event MUST precede durable lead save.

#### Scenario: Partial-write repair
- GIVEN successful ficha save followed by failed lead save
- WHEN generation retries
- THEN ficha upserts and the lead reaches `ENTREGADO` only after its save
- AND no premature event exists
