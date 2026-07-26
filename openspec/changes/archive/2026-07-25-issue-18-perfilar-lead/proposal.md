# Proposal: PerfilarLead (Issue #18)

## Intent

Vivi must know the prospect before greeting, but UC-01 steps 1–3 and UC-03 re-consultation have no implementation: no lead reaches `PERFILANDO` with verified affiliate data or precomputed capacity. Deliver one deterministic use case (US-01, US-02, RF-M3-05/06).

## Scope

### In Scope
- `internal/usecase/perfilar_lead.go`: `EntradaPerfilar`, `SalidaPerfilar`, `Ejecutar`, `ReconsultarPorFamiliar`, `ErrFamiliarNoEncontrado`.
- Affiliate hit: map `ingreso_hogar`, `categoria`, `segmento`, `personas_hogar` (`PersonasACargo+1`), `tipo_hogar` as `VERIFICADO_BASE`, set `Afiliado=true`, backfill an empty name. Miss: empty profile, `Afiliado=false`.
- Demo eligibility baseline, affiliate hit only: `tiene_vivienda=false`, `recibio_subsidio=false` as `VERIFICADO_BASE`. The motor treats absent booleans as unmet, so this baseline is what yields Ana's mandated 52,527,150.
- Create in `NUEVO`, capacity via `motor.CalcularCapacidad(perfil, afiliado, 0)`, `Lead.Transicionar(PERFILANDO)`, `Leads.Crear`, then `LeadNuevo` published only after a successful create.
- Re-consultation: `hogar_con_afiliado`/`cedula_familiar_afiliado` verified, family income added to `ingreso_hogar`, capacity recalculated, CAS `Guardar`. Unknown or inactive family cedula becomes `DECLARADO` + `RequiereConfirmacion` and returns `ErrFamiliarNoEncontrado`.
- `internal/usecase/perfilar_lead_test.go`: the four issue scenarios with a local `CatalogoFake`/`BusFake` plus existing `NuevoLeadRepoFake`, `NuevoRelojFake`, `NuevoIDFake`.

### Out of Scope
LLM turns, greeting/consent, HTTP routes, ADK graph, kNN candidate price, recommendations, ficha/plan, composition wiring, and any domain/motor/port/schema/data edit. The issue creates exactly two files, its Definition of Done never names them, and its tests must pass without database or LLM; adding them would cross NFR-M-01 and the 400-line budget.

## Capabilities

### New Capabilities
- `perfilamiento-lead`: deterministic pre-profile, affiliate provenance, household re-consultation.

### Modified Capabilities
- None. `capacidad` is consumed unchanged.

## Approach

Issue snippet drift resolved toward merged code (doc 13 owns money):

| Snippet | Resolution |
|---|---|
| `motor.ParametrosCredito{TasaEA}` | Real 3-argument `CalcularCapacidad`; drop `TasaEA`. Rate/term/ratio are reference constants. |
| `domain.Campo*`, `EstadoNuevo`, `FuenteVerificadoBase` | Contract keys as unexported usecase constants, `EstadoLead*`, `FuenteCampo*`. |
| `&RelojFake{T:…}`, `&IDFake{}` | `NuevoRelojFake(fixed)`, `NuevoIDFake(prefix)`. |

Idempotence and errors: each `Ejecutar` mints a new lead ID, so no re-entrant duplicate exists, and the call is retry-safe because nothing persists or publishes before `Crear` succeeds. Catalog miss, inactive affiliate, and catalog error all take the non-affiliate branch, so creation never fails on affiliate-data availability. Transition, `Crear`, and `Guardar` failures abort with wrapped errors and no event. `Perfil` is always non-nil.

## Affected Areas

| Area | Impact | Description |
|---|---|---|
| `internal/usecase/perfilar_lead.go` | New | Use case, verified mapping, re-consultation |
| `internal/usecase/perfilar_lead_test.go` | New | Four scenarios, local catalog/bus fakes |
| `internal/usecase/puertos.go`, `internal/domain/**` | Read-only | Consumed unchanged |

## Risks

| Risk | Likelihood | Mitigation |
|---|---|---|
| Baseline stamped `VERIFICADO_BASE` with no source field | Med | Synthetic affiliate hits only, documented in code and spec; real data must supply explicit fields |
| Catalog outage silently reads as "not affiliate" | Med | Required by the issue's discovery path; observability is a follow-up |
| Nil profile map or stale CAS version | Low | Initialize `Perfil`, reload via `PorID`, assert versions |

## Rollback Plan

Both files are new and nothing imports them (`cmd/servidor` keeps its unwired Block A seam). Revert the PR commit or delete `feat/issue-18-perfilar-lead`; no migration, data, config, or Contract change to undo. Depends only on merged #7, #8, #11, #15.

## Success Criteria

- [ ] Ana: verified fields, `ingreso_hogar` 2600000, subsidy 52527150, `PERFILANDO`, one `LeadNuevo`.
- [ ] Unknown cedula: empty profile, subsidy 0, `PERFILANDO`, `LeadNuevo` published.
- [ ] Family hit sums income with `hogar_con_afiliado` verified; family miss returns `ErrFamiliarNoEncontrado`.
- [ ] `go test ./internal/usecase/... -run 'TestPerfilar|TestReconsulta' -v`, then `go test ./...`, `go build ./...`, `go vet ./...`, `gofmt -l` clean.
- [ ] Authored implementation plus tests under 400 lines; no DB, LLM, HTTP, or wall clock in tests.

## Proposal question round

1. Is the demo eligibility baseline acceptable as demo-scoped, or should the affiliate fixture gain real fields in a `contrato` PR first?
2. Should a catalog outage be distinguishable from a genuine non-affiliate in `SalidaPerfilar`?
3. Confirm `personas_hogar = personas_a_cargo + 1` counts the titular.
4. Should re-consultation on an already-affiliate lead be rejected to avoid double-counting household income?

Assumed if unanswered: yes to 1 and 3, silent non-affiliate for 2, permissive re-consultation per the issue snippet.
