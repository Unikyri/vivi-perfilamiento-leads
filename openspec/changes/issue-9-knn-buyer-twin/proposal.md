# Proposal: Buyer-Twin kNN (Gower) — Issue #9

## Intent

The motor cannot yet answer "which real buyers resemble this lead". Doc 13 §3 requires a deterministic Gower kNN so recommendations come from comparable buyers instead of invented rules, and 27% of buyer records (non-affiliates) lack `categoria`/`rango_edad`. Issue #10 (`RecomendarProyectos`) is blocked until this exists.

## Scope

### In Scope
- `internal/domain/motor/knn.go`: safe input model, feature projection, project-local imputation, Gower distance, deterministic neighbors.
- `internal/domain/motor/knn_test.go`: KNN-1..5 plus zone-provenance and no-proxy tests.
- Frozen deterministic age-bracket representatives, including `55+`.

### Out of Scope
- Issue #8 capacity code (`internal/domain/motor/capacidad.go`, `capacidad_test.go`) — no edits.
- `RecomendarProyectos` (#10), usecase/HTTP wiring (#17), pipeline, `data/`, `Docs/`.
- Adding a zone field to `domain.Comprador`; any financial use of imputed data.

## Capabilities

### New Capabilities
- `buyer-twin-knn`: Gower feature projection, missing-data policy, distance, and deterministic neighbor selection.

### Modified Capabilities
- None.

## Approach

| Decision | Rule |
|---|---|
| Input | Buyer records + `CatalogZones map[proyectoID]zona` + `Perfil` + explicit lead `Afiliado` + `K`. **No independently indexed parallel slices.** |
| Zone | Only `CatalogZones[buyer.ProyectoID]`. Unknown/non-canonical ID → zone dimension absent. Never from project name, slug, or price. |
| Lead | `ingreso_hogar` → band (≤2 SMMLV A, >2–4 B, >4 C) for kNN only; `zona_deseada` as lead zone; absent income → category absent. |
| Imputation | Mode(`categoria`) / median(`rango_edad`) of **affiliated buyers of the same `ProyectoID`**, only when that project has ≥5 affiliates; otherwise omit each dimension and renormalize. `""`/`SIN_DATO` = absent. No price or global fallback. Imputed values never reach the financial motor. |
| Distance | Weights 0.35/0.20/0.15/0.15/0.15 over dimensions present in both points; `Σ(wᵢdᵢ)/Σwᵢ`; age range 32.5, dependents range 10. |
| Output | Sort distance ascending, tie-break `ID` ascending, return `min(k,n)`; inputs never mutated. |
| Delivery | One focused work unit on `feat/issue-9-knn-buyer-twin`; no parallel slices. |

### Age representative (decision)
`20-35`→27.5, `36-45`→40.5, `46-55`→50.5, **`55+`→60.0**. Doc 13 says "punto medio del rango" but leaves the open bracket undefined. 60.0 is chosen because: (1) it closes `55+` at 65, the practical upper age for mortgage origination, giving midpoint 60; (2) it reproduces the documented global range 60−27.5 = 32.5 — any other value silently changes the denominator; (3) it matches Issue #9's normative sketch, so tests stay byte-reproducible. `SIN_DATO`/empty → age absent, never substituted.

## Affected Areas

| Area | Impact | Description |
|---|---|---|
| `internal/domain/motor/knn.go` | New | Pure kNN implementation |
| `internal/domain/motor/knn_test.go` | New | Behavior + determinism tests |
| `internal/domain/capacidad.go` types | Read-only | `Comprador`, `Proyecto`, `Perfil` consumed unchanged |
| `internal/domain/motor/capacidad*.go` | Untouched | Issue #8 boundary |
| `internal/pipeline/*`, `data/*`, `Docs/*` | Read-only | Input contracts (Contract §4.1/§4.2) |

## Risks

| Risk | Likelihood | Mitigation |
|---|---|---|
| Unfrozen `55+` value breaks reproducibility | Med | Representative frozen above; constant test |
| Zone proxy drift (name/slug/price) | Med | Test proves price/name changes cannot alter distance; catalog zone changes can |
| Cross-project or global imputation leak | Med | Statistics keyed by canonical `ProyectoID`, affiliates only, ≥5 threshold test |
| Positional misalignment | Low | Record + keyed-map input; slice-parallel API rejected |
| Issue #8 collision | Low | Scope excludes capacity files; verify diff paths |

## Rollback Plan

Both files are new and unreferenced by production code. Rollback = `git revert` the work-unit commit (or delete `knn.go`/`knn_test.go`) on `feat/issue-9-knn-buyer-twin`; `main` and Issue #8 are unaffected because no existing file, contract type, data file, or wiring changes.

## Dependencies

- Issue #6 (`domain.Comprador`, `Proyecto`, `Perfil`) — merged.
- Caller must supply catalog zones by canonical `ProyectoID` (Issue #10/#17).

## Success Criteria

- [ ] KNN-1..KNN-5 pass (30 neighbors ordered, `k>n`, ID tie-break, imputation ≥5 affiliates, symmetry and `d(x,x)=0`).
- [ ] Absent dimensions omitted with renormalized weights.
- [ ] Test proves price/name cannot influence distance; catalog zone can.
- [ ] `go build ./...` and `go test ./internal/domain/motor/...` green; no capacity or pipeline file modified.

## Proposal question round

Autonomous delivery was authorized, so these assumptions are recorded for review instead of asked:
1. `55+` = 60.0 (range 32.5) is accepted as the frozen product decision.
2. `personas_a_cargo` present-vs-absent: `0` is treated as a real value, absent only when the source marks it missing.
3. Zone comparison uses the raw catalog zone string (e.g. `"Bogotá - Bosa"`) with case/space normalization, not a coarser city grouping.
4. The kNN returns neighbors only; desistimiento rates and project ranking stay in Issue #10.
