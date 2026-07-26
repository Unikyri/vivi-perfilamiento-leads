# Exploration: Buyer Persona, Production Clock, and Buyers Table (Issue #30)

## Current state
Current main already exposes deterministic buyer-persona summaries from the file-backed catalog, supports filtered and catalog-wide output, and persists `fecha_simulada` via the demo repository. The `compradores` JSON, port, kNN consumers, and SQL table are Contract-declared. The production clock currently aliases operational time to simulated time.

## Decision
Preserve the existing buyer-persona contract and the `compradores` table. Removing the table would require a separate Contract §9 PR and both-block approval, so it is explicitly out of scope. Make the clock's wall-clock timestamps independent from demo time while retaining persisted demo time, non-regression behavior, reset behavior, and milestone comparisons.

## Evidence
Issue #30 body; Contract evidence recorded in Engram exploration #631; current buyer-persona endpoint/tests, catalog cache, migration, demo repository, and `internal/infrastructure/reloj/postgres.go`.
