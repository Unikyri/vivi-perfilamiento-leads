## Exploration: issue-10-recommendations-routing

### Current State
Issue #10 provides deterministic recommendations and capacity × intention routing. Application is blocked until #8 is merged; this change does not inspect, use, or edit #8's branch/files. Issue #9 supplies value-only `Vecino`; no recommendation/routing implementation exists.

`domain.Recomendacion`, `Proyecto`, `Capacidad`, `Intencion`, and `Ruta` already provide the shared vocabulary. The motor remains pure: no HTTP, database, LLM, lead/ficha mutation, or milestone creation belongs here.

### Authority Reconciliation
1. **Eligibility:** candidate filtering must implement inclusive `Ratio >= 0.8` with integer money, not float division. For positive prices: `budget >= ceil(0.8*price)`, equivalently `budget >= price - price/5`.
2. **Ratio meanings:** recommendation eligibility is candidate-price based. Matrix ratio follows the Contract/Doc 13 best-affordable-project rule; do not reuse the old fixed median except its pre-candidate fallback.
3. **Route only:** `Matriz2x2` returns only `domain.Ruta`; affiliate milestones and cupo effects are downstream.
4. **#8 isolation:** no capacity code/tests, shared types, Docs, data, pipeline, or #8 branch edits.
5. **Stable ties:** comparator must end with canonical `ProyectoID`, not map/input order.

### Boundary Cases
- Negative/zero budget, non-positive catalog price, missing catalog ID, empty input, and invalid/unknown intention must not produce an advisor route/card.
- Aggregate counts and desistimientos by canonical project ID only; all buyer rows may contribute density, but only catalog entries produce cards.
- Empty/no result returns a non-nil empty slice.

### Approaches
1. **Two pure functions** — catalog-backed neighbor aggregation and separate capacity/intention route selection. Recommended.
2. **Mutate Lead/Ficha or create milestone** — rejected: orchestration side effects.
3. **Float/map-order design** — rejected: boundary drift and non-reproducibility.

### Recommendation
Use typed, side-effect-free motor inputs. Recommendations use exact catalog IDs, integer-safe budget eligibility, and deterministic ranking ending in ID. Matrix returns one route from capacity × intention only. Freeze `NivelMedia` classification and complete ranking ties in proposal/spec. Do not apply until #8 merges.

### Risks
- Candidate affordability and matrix ratio cannot be conflated.
- Invalid money/catalog values need explicit rejection.
- #8 merge is a hard apply/verification dependency.

### Ready for Proposal
Yes — planning may continue; application remains blocked by #8.
