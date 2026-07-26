# Proposal: Buyer Persona, Production Clock, and Buyers Table Decision (Issue #30)

## Scope
Buyer-persona aggregation and its API already meet the issue's deterministic data requirements. Preserve them, add focused endpoint coverage, and correct the production clock: `Ahora()` returns real UTC wall time while `FechaSimulada()` and `Avanzar()` retain the persisted demo date, non-regression behavior, reset, milestones, and `/salud` behavior.

Add `docs/decisiones/0001-conservar-tabla-compradores.md`, recording that the Contract-declared `compradores` table is retained. It is part of the buyer JSON/port/kNN/buyer-persona data boundary. Removal requires a separately labelled Contract §9 PR and both-block approval.

## Constraints
No migration, Contract change, SQL aggregation, dashboard UI expansion, or `Reloj` interface change. Existing injected fakes remain valid.

## Evidence and verification
Test the wall-vs-simulated clock split, persisted non-regressing demo date, and buyer-persona filtered/catalog contract shape. Run the complete Go suite/build/vet. Forecast: under 160 authored lines.

## Rollback
Revert the isolated adapter/documentation PR; no stored data needs migration or cleanup.
