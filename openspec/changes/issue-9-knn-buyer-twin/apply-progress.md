# Apply Progress: Buyer-Twin kNN — PR 3/5

## Scope
- `auto-chain` / `stacked-to-main`; standard mode.
- Completed 2.1 and 2.2 only: safe contract projection and its RED/GREEN tests.
- Remaining: PR 4 statistics/Gower (3.1–3.2), then PR 5 selection/evidence (4.1–4.3).

## Delivered
- `knn.go`: typed `EntradaGemelo`/`Vecino` and pure lead/buyer projection.
- `knn_test.go`: exact catalog zones, missing zones, income A/B/C, all age brackets including `55+=60`, normalization/missingness, zero dependents, explicit false affiliation, purity, and price/name non-proxy tests.
- Only task checkboxes 2.1/2.2 changed; no capacity, pipeline, data, Docs, Gower, statistics, or selection edits.

## Decisions and Evidence
- Zones use non-empty exact `ProyectoID` map lookup; missing/sentinel values omit the dimension.
- `SIN_DATO`/unknown values omit optional features; zero dependents remains present; name/slug/price are never read.
- `go test ./internal/domain/motor/... -run TestGemeloKNNProjection -count=1` — PASS.
- Runtime: N/A — pure, unwired domain service.
- Rollback: remove the two kNN files and restore two task checkboxes.
