# Apply Progress: Domain Contract Types

## Completion
All 9 SDD tasks are complete across PRs 4–7 under `stacked-to-main`.

## Delivered
- Exact 11 Contract enums; Perfil schema, 18 recognised keys, four critical keys, and accessors.
- Capacidad/recommendation types plus compatible Comprador/Proyecto relocation.
- Lead, Mensaje, Hito, PlanNutricion, and Ficha structures with exact JSON tags.
- Contract reflection/JSON tests and retained pipeline compatibility.

## Evidence
- `go test ./internal/domain/... -v` — pass.
- `go test ./internal/pipeline/... -v` — pass.
- `go build ./internal/domain/...`, `go build ./...`, `go vet ./internal/domain/...` — pass.
- Dependency guard — no `usecase`, `adapters`, or `infrastructure` dependency.
- `Docs/` and pipeline sources unchanged.
