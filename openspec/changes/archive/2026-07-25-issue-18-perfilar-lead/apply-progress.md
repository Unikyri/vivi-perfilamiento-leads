# Apply Progress: Deterministic PerfilarLead (Issue #18)

## Mode

Standard mode (`strict_tdd: false`), single-PR work unit, no size exception. This batch is a focused remediation of the independent verifier warnings; native review was not started.

## Completed Tasks

- [x] 1.1 Added `EntradaPerfilar`, `SalidaPerfilar`, `PerfilarLead`, family sentinel errors, and recognized profile-key helpers.
- [x] 1.2 Implemented deterministic lead creation with fixed ID/time ports, non-nil profile, active-affiliate mapping/name backfill, and demo eligibility baseline; misses, inactive records, and non-context catalog errors use the empty non-affiliate path.
- [x] 1.3 Calculated through `motor.CalcularCapacidad(perfil, afiliado, 0)`, transitioned through `Lead.Transicionar`, created before publishing one Contract v1.1 `EvLeadNuevo` event with exactly `cedula`, `nombre`, `telefono`, and `fuente` in `Payload` while retaining `LeadID` in `Evento`.
- [x] 1.4 Implemented family re-consultation load, idempotent same verified cedula handling, distinct verified cedula rejection, verified income addition/recalculation, and CAS save without lifecycle events.
- [x] 1.5 Persisted unknown/inactive family cedulas as declared confirmation-required data and returned `ErrFamiliarNoEncontrado`; context and save errors take precedence.
- [x] 2.1 Added compact interface-complete local catalog and bus fakes; reused the existing repository, clock, and ID fakes.
- [x] 2.2 Added table-driven active, unknown, inactive, and catalog-error acceptance coverage, including Ana's verified profile, household count, subsidy, state, version, and exact Contract event payload; failure cases assert zero events.
- [x] 2.3 Added family acceptance coverage for `hogar_con_afiliado=true` with `VERIFICADO_BASE`, household income `4,500,000`, recalculated subsidy `35,018,100`, repeat idempotence, unchanged capacity, and no additional event.
- [x] 2.4 Added canceled-context, `PorID` failure, create failure, and save/CAS failure coverage with no success event on failure.
- [x] 3.1 Formatted files and passed the focused profiling/re-consultation suite.
- [x] 3.2 Passed full tests, race tests, build, vet, module verification, and changed-file gofmt checks.
- [x] 3.3 Confirmed 356 authored implementation/test lines, below the 390-line hard stop.

## Work Unit Evidence

| Evidence | Result |
|---|---|
| Focused test command and exact result | `go test ./internal/usecase/... -run 'TestPerfilar|TestReconsulta' -v` — PASS; 4 test functions and all subcases passed, including exact event payload, family recalculation/idempotence, canceled context, and `PorID` failure. |
| Runtime harness command/scenario and exact result | N/A: this is a provider-free usecase boundary exercised entirely with deterministic in-memory ports; no runtime integration boundary is in scope. |
| Rollback boundary | Revert/delete only `internal/usecase/perfilar_lead.go` and `internal/usecase/perfilar_lead_test.go`; no existing production files, schema, data, or wiring are changed. |

## Validation Evidence

- `go test ./internal/usecase/... -run 'TestPerfilar|TestReconsulta' -v` — PASS.
- `go test ./... -count=1` — PASS.
- `go test -race ./internal/...` — PASS; no races.
- `go build ./...` — PASS.
- `go vet ./...` — PASS.
- `go mod verify` — PASS (`all modules verified`).
- `gofmt -l internal/usecase/perfilar_lead.go internal/usecase/perfilar_lead_test.go` — empty output.
- `gofmt -l internal cmd` — reports only pre-existing `internal/pipeline/compradores_test.go` and `internal/pipeline/proyectos_test.go`; neither is in Issue #18 scope.
- Authored line count excluding OpenSpec: production 178 + tests 178 = 356; hard stop 390.
- Scope check: exactly the two planned source/test files are implementation files; OpenSpec artifacts are the only additional SDD files.

## Remediation Coverage

- `LeadNuevo` success events are asserted as exact map equality: `cedula`, `nombre`, `telefono`, and `fuente`; `Tipo` and `Evento.LeadID` are also asserted. Canceled execution and create failure assert zero events.
- Active family re-consultation asserts `hogar_con_afiliado=true`, verified provenance for household and family cedula, one-time income addition to `4,500,000`, subsidy `35,018,100`, repeat unchanged income/version/capacity, and exactly the original one event.
- Re-consultation canceled context and `PorID` errors are asserted before mutation; no event is emitted by the family path.

## Implementation Notes

- Affiliate lookup errors are treated as non-affiliate only when they are not context cancellation/deadline errors; context failures remain errors.
- The event payload now follows Contract v1.1 §6 and contains no state-only or profile-internal fields.
- Family re-consultation uses `lead.Afiliado` plus the verified `hogar_con_afiliado` profile field when recalculating through the existing motor.

## Status

All 12 tasks remain complete. Focused remediation is complete and ready for the required native review; do not archive or perform final verification until review authority is established.
