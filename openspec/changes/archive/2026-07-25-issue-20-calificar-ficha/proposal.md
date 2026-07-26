# Proposal: Calificar lead and generate deterministic advisor ficha (Issue #20)

## Intent

Conversations reach `CALIFICADO` but nothing qualifies the lead or produces the advisor ficha. Advisors get no route, priority, or decision sheet, so profiled leads stall. Add the two missing application use cases so every completed profile becomes an auditable, deterministic route plus advisor ficha.

## Scope

### In Scope
- `internal/usecase/calificar_lead.go` — `CalificarLead` qualification orchestration.
- `internal/usecase/generar_ficha.go` — `GenerarFicha` deterministic ficha output.
- Deterministic tests for both (`internal/usecase/*_test.go`), reusing existing fakes.
- SDD artifacts (hybrid OpenSpec + Engram).

### Out of Scope
- `internal/domain`, `internal/domain/motor`, `internal/usecase/puertos.go`.
- Infrastructure, adapters, HTTP, frontend, migrations, config, Contract/Wiki docs.
- LLM/narrative generation, new ports, repository transactions.

## Locked Behavior

`CalificarLead` (accepts only `CALIFICADO`):
1. Capacity with candidate `0` → preliminary budget.
2. Candidate = lowest positive catalog `PrecioDesde` ≤ preliminary budget; recalculate capacity with it (motor median fallback only when candidate is `0`).
3. `GemeloKNN` with `K=30` (profile, affiliation, dependents, exact catalog zones) → `RecomendarProyectos` with final budget.
4. Conversion = non-affiliate AND (`INDEPENDIENTE` OR `hogar_con_afiliado` OR nonblank `caja_externa`).
5. Route via `Matriz2x2`; Contract §3.5 priority with capacity ratio clamped at `1.2`; semáforo via private helper; `ConsumeCupo10` only for non-affiliate `ASESOR`.
6. `ASESOR` stays `CALIFICADO`; `NUTRICION`/`REMARKETING`/`DESPEDIDA` → `EN_NUTRICION`/`REMARKETING`/`DESPEDIDO`.
7. Save lead (CAS) **then** publish `RutaDecidida`.

`GenerarFicha` (accepts only `CALIFICADO` + `ASESOR`): recompute the identical deterministic recommendation set, build the Contract `Ficha` (no LLM) with fixed-order `Beneficios`, `ArgumentosVenta`, `AlertaDesistimiento`, and the verbatim low-confidence band `PERFIL PARCIALMENTE DECLARADO — validar campos marcados`; persist ficha **then** transition lead to `ENTREGADO` and save it.

## Capabilities

### New Capabilities
- `calificar-lead`: qualification, candidate price, conversion, routing, priority, cupo, state/event ordering.
- `generar-ficha`: deterministic Contract ficha content, ordering, persistence and handoff to `ENTREGADO`.

### Modified Capabilities
None.

## Approach

Two provider-free use cases reusing existing ports (`LeadRepository`, `FichaRepository`, `CatalogoRepository`, `GeneradorID`, `Reloj`, `BusEventos`) and pure motor functions. No motor formula is reimplemented.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `internal/usecase/calificar_lead.go` | New | Qualification and routing |
| `internal/usecase/generar_ficha.go` | New | Deterministic ficha |
| `internal/usecase/*_test.go` | New | Deterministic table tests |

## Delivery (feature-branch-chain, no size exception)

| Slice | Content | Target | Budget |
|-------|---------|--------|--------|
| PR #1 | `calificar_lead.go` + tests | `feature/bloque-a` | ≤400 authored lines |
| PR #2 | `generar_ficha.go` + tests, error/retry paths | PR #1 branch | ≤400 authored lines |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Ficha persisted, later lead save fails | Med | Ficha upsert makes retry repair idempotent; no event before durable save |
| Median fallback used despite affordable candidate | Med | Explicit lowest-positive selection test |
| Recommendation drift between use cases | Med | Same zones, buyers, `K=30`, motor ranking; cross-use-case equality test |
| Blank `caja_externa` treated as conversion | Low | Trim/blank test cases |
| Slice 1+2 merged as one PR | Low | Chain enforced at apply |

## Rollback Plan

Both files are new and unreferenced by existing call sites. Revert per slice: `git revert` the slice commit (or delete the new file plus its test) — no schema, port, or Contract change to undo; earlier slice remains valid on its own.

## Dependencies

Issue #18 (`PerfilarLead`) and #19 (`ProcesarMensaje`) archived; existing motor and ports unchanged.

## Success Criteria

- [ ] `CalificarLead` reproduces Contract §3.5 route, priority, semáforo, and cupo for all matrix quadrants.
- [ ] `GenerarFicha` emits byte-stable Contract ficha with fixed order and verbatim warning.
- [ ] State/event ordering verified: save before event; ficha before `ENTREGADO`.
- [ ] Rejected states/routes return errors without writes.
- [ ] Two slices each ≤400 authored lines; `go build ./... && go test ./...` pass.

## Proposal question round

Asked non-interactively; assumptions taken from the verified exploration and Issue #20. Review if any is wrong:
1. `CalificarLead` is invoked explicitly after `PerfilCompleto` (no event subscription in scope).
2. Non-`ASESOR` routes never produce a ficha in this change.
3. Partial-write repair stays manual/retry-based; no transactional port is added now.
4. Fixed benefit/argument text comes from the Issue #20 constants, not product copy review.
