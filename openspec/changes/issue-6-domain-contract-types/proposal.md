# Proposal: Domain Contract Types (Issue #6)

## Intent

Domain layer lacks the type vocabulary for Contract v1.1 §1–§5. Engine, agent, and use-case packages are blocked until these exist.

## Scope

### In Scope
- Contract §1 enums (`enums.go`)
- `Perfil` + `CampoPerfil` with typed accessors (`perfil.go`)
- `Capacidad`, `Intencion`, `Recomendacion`, `ItemDesglose`, consolidated `Comprador`/`Proyecto` (`capacidad.go`)
- `Lead`, `Mensaje` (`lead.go`)
- `PlanNutricion`, `Hito` (`plan.go`)
- `Ficha`, `AlertaDesistimiento`, `Identificacion` (`ficha.go`)
- Unit tests for `Perfil` accessors (`perfil_test.go`)
- Deletion of redundant `comprador.go` and `proyecto.go`

### Out of Scope
- Engine logic, persistence adapters, API handlers, frontend types

## Capabilities

### New Capabilities
- `domain-contract-types`: Pure Go structs, enums, and value objects implementing Contract v1.1 §1–§5.

### Modified Capabilities
None

## Approach

**Delete-and-Replace** — remove `comprador.go`/`proyecto.go`; create 6 files (Pasos 1–6) + 1 test (Paso 7) per issue #6. Zero-dependency (stdlib only). Pipeline keeps identical API surface.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `internal/domain/comprador.go` | Removed | Moves to `capacidad.go` |
| `internal/domain/proyecto.go` | Removed | Moves to `capacidad.go` |
| `internal/domain/enums.go` | New | Contract §1 enums |
| `internal/domain/perfil.go` | New | Perfil value object |
| `internal/domain/capacidad.go` | New | Capacidad + Comprador + Proyecto |
| `internal/domain/lead.go` | New | Lead + Mensaje |
| `internal/domain/plan.go` | New | PlanNutricion + Hito |
| `internal/domain/ficha.go` | New | Ficha + AlertaDesistimiento |
| `internal/domain/perfil_test.go` | New | Accessor tests |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Positional struct init in consumers | Low (verified none) | grep before delete |
| NFR-M-01 violation (external import) | Low | Review: only `time` allowed |

## Rollback Plan

`git revert` the merge commit. Or restore deleted files from `HEAD~1` and remove new files. Pipeline compiles identically.

## Dependencies

- Issue #6 spec (GitHub) — field names, types, JSON tags
- Contract v1.1 §1–§5 (Wiki doc 10) — authoritative schema

## Success Criteria

- [ ] `go build ./...` passes
- [ ] `go vet ./...` clean
- [ ] `go test ./internal/domain/...` passes
- [ ] No imports outside stdlib in new domain files (NFR-M-01)
- [ ] All Contract §1–§5 types present and exported
- [ ] Pipeline compiles without modification

## Proposal Question Round

> Autonomous mode — questions with supported assumptions.

| # | Question | Assumption |
|---|----------|------------|
| 1 | `Perfil` accessors: `(value, ok)` or panic? | `(value, ok)` — Go idiom |
| 2 | Enum types: `string` or `int`? | `string` — matches JSON tags |
| 3 | Optional `time.Time` as pointer? | Yes — `*time.Time` |
| 4 | `Ficha.Score` type? | `float64` — decimal precision |
