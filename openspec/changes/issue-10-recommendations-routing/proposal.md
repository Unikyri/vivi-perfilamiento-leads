# Proposal: Recommendations and 2x2 Routing — Issue #10

## Intent

The motor can find similar buyers (#9) but nothing converts that evidence into advisor-facing project cards or a commercial route. Doc 13 §4/§5 (PRD RF-M2-03/RF-M2-04) close the deterministic engine; Issue #17 (`CalificarLead`) is blocked until both exist.

## Scope

### In Scope
- `internal/domain/motor/recomendar.go`: neighbor aggregation by canonical `ProyectoID`, catalog filter, integer-safe eligibility, frozen ranking, max 3 cards.
- `internal/domain/motor/matriz.go`: pure `Matriz2x2(EntradaMatriz) domain.Ruta`.
- `recomendar_test.go`, `matriz_test.go`: REC-1..3, MAT-1..7, boundary and purity tests.

### Out of Scope
- Issue #8 capacity (`capacidad.go`, `capacidad_test.go`, its branch): not read, not edited, not recomputed.
- `Hito`/`PlanNutricion` creation, `ConsumeCupo10`, `Lead`/`Ficha` mutation, prioridad, usecase/HTTP/ADK wiring, catalog data, `Docs/`.

## Capabilities

### New Capabilities
- `recommendations-routing`: deterministic project recommendation and 2x2 route selection.

### Modified Capabilities
- None. `Recomendacion`, `Proyecto`, `Capacidad`, `Intencion`, `Ruta` are consumed unchanged.

## Approach — Frozen Decisions

| # | Decision | Rule |
|---|---|---|
| 1 | Application gate | Planning only; **apply/verify blocked until #8 is merged to `main`**. Matrix consumes a caller-supplied `Ratio`. |
| 2 | Integer-safe eligibility | `PrecioDesde > 0 && PresupuestoMax >= PrecioDesde - PrecioDesde/5` (= `ceil(0.8·price)`). No float division, no overflow, no `/0`. Replaces the sketch's `float64(precio) > float64(presupuesto)/0.8`. |
| 3 | Invalid data | `PrecioDesde <= 0`, negative budget, empty/non-catalog `ProyectoID`, `K<=0`, empty catalog/neighbors → excluded; result is non-nil empty. Catalog projects with 0 contributing neighbors emit no card. Never merge by display name. |
| 4 | Ranking comparator (total order) | `Vecinos` desc → desistimiento asc (integer cross-product `dᵃ·nᵇ < dᵇ·nᵃ`, no float compare) → `ProyectoID` asc → truncate to 3. Input/map order is never a tie-break. |
| 5 | Route-only matrix | Returns only `domain.Ruta`; no mutation, no milestone, no cupo. Capacity high = `Ratio >= 0.95`; non-affiliate ASESOR needs `Ratio >= 1.05`; conversion route + high intention → `NUTRICION`. Thresholds stay 0.95/1.05 (Doc 13 §5), distinct from the 0.8 candidate filter (Doc 13 §4). |
| 6 | `NivelMedia` | High intention is **`Nivel == ALTA` only** (authoritative 2x2 table). `MEDIA`, `BAJA`, empty and unknown values fall on the low-intention side; `NaN`/absent ratio is not high capacity, so no invalid input reaches `ASESOR`. |
| 7 | No capacity change | Ratio semantics (best affordable catalog project, Contract v1.1 §2.2) are #8's output; #10 neither recomputes nor redefines them. |

## Affected Areas

| Area | Impact | Description |
|---|---|---|
| `internal/domain/motor/recomendar.go` | New | Aggregation, eligibility, cards, ranking |
| `internal/domain/motor/matriz.go` | New | Pure route selection |
| `internal/domain/motor/*_test.go` (2 new) | New | REC/MAT + boundary + purity |
| `internal/domain/{capacidad,ficha,lead}.go` | Read-only | Shared shapes unchanged |
| `internal/domain/motor/knn.go` | Read-only | `Vecino` consumed as-is |

## Risks

| Risk | Likelihood | Mitigation |
|---|---|---|
| Wiki Doc 13/Contract files absent from worktree | High | Decisions 2–7 sourced from persisted authority (#412/#413) + issue sketch; spec must quote exact clauses once the wiki is checked out |
| #8 merge changes `Ratio` meaning | Med | Matrix takes `Ratio` as input; re-verify MAT cases after merge before apply |
| Float ratio conflated with integer money | Med | Eligibility recomputed from `int64` operands only; stale `Capacidad.Ratio` never admits a candidate |
| Non-reproducible ordering | Med | Comparator frozen in decision 4 with `ProyectoID` terminator |
| Scope leak into #8 or downstream orchestration | Low | Diff limited to 4 new files; route-only assertion test |

## Rollback Plan

All four files are new and referenced by no production code. Rollback = `git revert` the work-unit commit (or delete the four files) on the Issue #10 branch. `main`, `#8`, `#9`, contract types, and data files are untouched, so no migration or data repair is needed.

## Dependencies

- **#8 capacity — hard, unmet**: authoritative `Capacidad.PresupuestoMax`/`Ratio`. Apply and verify stay blocked until merged to `main`.
- #6 domain types (merged), #9 `Vecino` (present on this branch).
- Caller supplies catalog keyed by canonical `ProyectoID`, plus `EsAfiliado` and `RutaConversion`.

## Success Criteria

- [ ] REC-1 max 3, REC-2 eligibility filter, REC-3 out-of-catalog ignored.
- [ ] Boundary tests: budget one peso below `ceil(0.8·price)` fails, exactly at it passes; `price<=0` and negative budget excluded.
- [ ] Identical output across shuffled neighbor/catalog input order.
- [ ] MAT-1..7 pass; all four quadrants covered; `MEDIA` routes as non-high; no invalid input yields `ASESOR`.
- [ ] Purity: inputs unmutated, no `Hito`/`ConsumeCupo10`/`Lead` write.
- [ ] `go build ./...` and `go test ./internal/domain/motor/...` green; motor coverage ≥ 90%; no #8 file in the diff.

## Proposal question round

Recorded assumptions for review (autonomous run; answer or correct any):
1. Recommendation cap is exactly 3 (REC-1) even when more projects are eligible.
2. `MEDIA` intention is low intention everywhere in the matrix (decision 6) — confirm no `MEDIA`-specific route is expected.
3. `RutaConversion` is a caller-provided boolean; #10 does not infer affiliation eligibility.
4. `Razon` text stays the sketch's budget phrasing (`"Encaja con tu presupuesto de $NM"`), rounded by integer division.
5. Candidate eligibility (0.8) and matrix capacity (0.95/1.05) remain two separate thresholds.
