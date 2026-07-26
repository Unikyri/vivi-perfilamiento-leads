# Tasks: Deterministic PerfilarLead (Issue #18)

## Review Workload Forecast

| File | Work | Forecast |
|---|---|---:|
| `internal/usecase/perfilar_lead.go` | DTOs and deterministic use case | 160 |
| `internal/usecase/perfilar_lead_test.go` | fakes and table-driven acceptance/error tests | 210 |
| **Total authored code/tests** | **single PR** | **370** |

Hard implementation stop: 390 authored changed lines; prune duplicate table setup/assertions before apply—never request `size:exception`.

Decision needed before apply: No
Chained PRs recommended: No
Chain strategy: pending
400-line budget risk: Low

### Suggested Work Units

| Unit | Goal | Likely PR | Focused test command | Runtime harness | Rollback boundary |
|---|---|---|---|---|---|
| 1 | Complete isolated profiling use case and tests | Single PR to `main` | `go test ./internal/usecase/... -run 'TestPerfilar|TestReconsulta' -v` | N/A: deterministic in-memory ports only | Delete the two new files/revert this PR |

## Phase 1: Use-case implementation — `internal/usecase/perfilar_lead.go` (160 lines)

- [x] 1.1 Add `EntradaPerfilar`, `SalidaPerfilar`, `PerfilarLead`, family sentinel errors, and private recognized-profile-key helpers; use only existing ports/domain types.
- [x] 1.2 Implement `Ejecutar`: fixed port ID/time, non-nil `NUEVO` lead, active-affiliate mapping/name backfill and demo baseline; missing, inactive, or non-context catalog failure remains non-affiliate.
- [x] 1.3 Calculate through `motor.CalcularCapacidad(perfil, afiliado, 0)`, transition via `Lead.Transicionar(PERFILANDO)`, `Crear`, then publish exactly one `LeadNuevo` only after create succeeds; cancellation/transition/create errors publish nothing.
- [x] 1.4 Implement `ReconsultarPorFamiliar`: load, same verified cedula no-op, distinct verified cedula rejection, active-family verified mapping/income addition/recalculation/CAS save, and no lifecycle event.
- [x] 1.5 Persist unknown/inactive family cedula as declared plus confirmation-required; return `ErrFamiliarNoEncontrado` after a successful save, while load/save/context errors take precedence.

## Phase 2: Deterministic tests — `internal/usecase/perfilar_lead_test.go` (210 lines)

- [x] 2.1 Add compact interface-complete local `CatalogoFake` and `BusFake`; reuse `NuevoLeadRepoFake`, `NuevoRelojFake`, and `NuevoIDFake` without real services or wall clock.
- [x] 2.2 Add table-driven acceptance tests: Ana verified mapping/demo baseline, `personas_hogar=3`, subsidy `52527150`, `PERFILANDO`, persisted version, and one post-create event; unknown/inactive/catalog-error fallback has empty profile and zero subsidy.
- [x] 2.3 Add family tests: verified income is added once, repeat is idempotent, a distinct verified cedula is rejected, and an unknown family saves declared confirmation-required data then returns the sentinel.
- [x] 2.4 Add error-safety tests: canceled context, transition/create failure, and CAS/save failure preserve failure semantics and emit no success event; consolidate shared setup/assertions before adding cases.

## Phase 3: Validation and delivery

- [x] 3.1 Run `gofmt -w internal/usecase/perfilar_lead.go internal/usecase/perfilar_lead_test.go` and `go test ./internal/usecase/... -run 'TestPerfilar|TestReconsulta' -v`.
- [x] 3.2 Run `go test ./...`, `go build ./...`, `go vet ./...`, and confirm `gofmt -l` is empty.
- [x] 3.3 Check `git diff --stat` and `git diff --numstat`; if authored production/test changes exceed 390, remove duplicated test scaffolding/assertions before review and keep code/tests in one rollbackable commit.
