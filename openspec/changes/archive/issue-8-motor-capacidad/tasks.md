# Tasks — `CalcularCapacidad` deterministic engine (Issue #8)

Scope: `internal/domain/motor/capacidad.go` + `capacidad_test.go` only. No other file is created or modified. No implementation happens in this phase — this is the checklist for `sdd-apply`.

## Review Workload Forecast

| Field | Value |
|---|---|
| Estimated changed lines | ~350 Go (capacidad.go ~140 + capacidad_test.go ~210, design §9) |
| SDD artifact accounting | openspec/changes/issue-8-motor-capacidad/{proposal,design,tasks}.md + specs/capacidad/spec.md + state.yaml + exploration.md ≈ 600+ lines, already tracked separately from the Go diff |
| 400-line budget risk | Low for the Go diff alone; High if SDD markdown is bundled into the same PR as the Go diff |
| Chained PRs recommended | No — the Go diff alone fits one PR; recommend a docs/code split, not a review chain |
| Suggested split | Single PR `feat/issue-8-motor-capacidad` → `main` carries only `capacidad.go` + `capacidad_test.go`; SDD markdown lands in its own commit (or its own `docs` PR) so reviewer attention stays on the 350 lines that move money |
| Decision needed before apply | Yes — confirm docs-vs-code bundling before opening the PR (design §9 flags this explicitly) |

```text
Decision needed before apply: Yes
Chained PRs recommended: No
Chain strategy: pending
400-line budget risk: Low
```

### Suggested Work Units

| Unit | Goal | Likely PR | Notes |
|---|---|---|---|
| 1 | `capacidad.go` + `capacidad_test.go` (all phases below) | PR 1 — `feat/issue-8-motor-capacidad` → `main` | ~350 lines, Low risk, self-contained |
| 2 | SDD markdown artifacts (proposal/spec/design/tasks/state) | Separate commit or `docs` PR | Not counted against the 400-line code budget; keep out of Unit 1's diff |

## Implementation Order

Single file, single author, one call-dependency chain — no parallel tracks. Phases 1→6 are strictly sequential (each helper's signature or return value feeds the next). Phase 7 (tests) requires all of Phase 1–6 to compile first — do not write the full table-driven suite before `calcular`/`CalcularCapacidad` exist. Phase 8 is the final gate.

## Phase 1: Constants and loaded state (`capacidad.go`)

- [x] 1.1 Declare consts `penalizacionCampoNoVerificado = 0.85`, `confianzaMinima = 0.50` (design §1, §5). Traces: Spec "Confidence score".
- [x] 1.2 Implement `mustLoad[T any](v T, err error) T` panic wrapper (design §1 item 7, §2 — panic beats silent zeroed constants). Traces: Proposal D1.
- [x] 1.3 Declare package vars `constantes = mustLoad(domain.LoadConstantes(references.ConstantesJSON))`, `tramosSubsidio = mustLoad(domain.LoadSubsidios(references.SubsidiosJSON))` (design §2). Traces: Proposal D1; state.yaml `architecture_authority` (references import is NFR-M-01-compliant).

## Phase 2: Annuity math

- [x] 2.1 Implement `factorAnualidad(tasaEA float64, plazoMeses int) float64`: `i = (1+tasaEA)^(1/12)-1`, `factor = (1-(1+i)^-n)/i`, zero intermediate rounding (design §3). Traces: Spec "Annuity credit calculation".

## Phase 3: Subsidy eligibility

- [x] 3.1 Implement `subsidioAplicable(p domain.Perfil, esAfiliado bool, cts domain.Constantes, tramos []domain.TramoSubsidio) (int64, string)`: income guard FIRST (absent or `<=0` → return `0, regla`), then the four S1/S2 conditions in any order, then the SMMLV-integer bracket lookup (never float division) (design §4). Traces: Spec "Subsidy eligibility and bracket lookup"; Proposal D3, D6.

## Phase 4: Confidence

- [x] 4.1 Implement `calcularConfianza(p domain.Perfil) float64`: `for campo := range domain.CamposCriticos` (key iteration only), multiply by `penalizacionCampoNoVerificado` when `!p.EsVerificado(campo)`, floor via `math.Max(c, confianzaMinima)` (design §5). Traces: Spec "Confidence score"; Proposal D4/S3. No order-dependent construct over the map (early break, sorted keys, accumulation).

## Phase 5: Desglose

- [x] 5.1 Implement `armarDesglose(p domain.Perfil, cts domain.Constantes, credito, subsidio, recursos int64, reglaSubsidio string) []domain.ItemDesglose`: always 3 items (`CREDITO`, `SUBSIDIO`, `RECURSOS_PROPIOS`); `Regla` strings built from `cts`/`tramos` values, never a literal; `Fuente` derived from the driving field's provenance, emitting empty `domain.FuenteCampo("")` when that field is absent — never substituted (design §6, F1). Traces: Spec "Desglose breakdown"; Proposal S5/S6.

## Phase 6: Public surface

- [x] 6.1 Implement `calcular(p domain.Perfil, esAfiliado bool, precioCandidato int64, cts domain.Constantes, tramos []domain.TramoSubsidio) domain.Capacidad`: wire `factorAnualidad` → single `int64(math.Round(cuota*factor))` for `CreditoMax`, `subsidioAplicable`, `calcularConfianza`, `Ratio` with `precioCandidato<=0` fallback to `cts.MedianaVISCOP`, `armarDesglose`. Traces: Spec "Ratio computation".
- [x] 6.2 Implement `CalcularCapacidad(p domain.Perfil, esAfiliado bool, precioCandidato int64) domain.Capacidad` delegating to `calcular` with package-level `constantes`/`tramosSubsidio`. Traces: Spec "Public API surface" scenario.

## Phase 7: Tests (`capacidad_test.go`)

- [x] 7.1 Add `casoCapacidad` fixture struct, `newCampo(valor any, fuente domain.FuenteCampo)`, `assertClose(t, nombre string, got, want float64)` (epsilon `1e-9`), following `internal/domain`'s table-driven/`t.Helper()` convention (design §6).
- [x] 7.2 Write ten `t.Run` subtests covering all twelve doc 13 §7 IDs: `CAP-1_Ana`, `CAP-2`, `CAP-3`, `CAP-4_tramo_medio`, `CAP-5_mas_4_smmlv`, `CAP-6_B-1_ingreso_cero`, `CAP-7_tiene_vivienda`, `CAP-8_B-4_confianza_minima`, `B-5_ingreso_supera_tope`, `B-6_precio_candidato_cero`. Exact `int64` equality on `CreditoMax`/`SubsidioAplicable`/`PresupuestoMax`; `assertClose` on `Ratio`/`Confianza`; assert `len(Desglose)==3` plus `Concepto`/`Fuente` per row. Traces: every Spec scenario; Proposal D5/D6.
- [x] 7.3 Add `TestFactorAnualidad`: `factorAnualidad(0.107, 240)` equals `102.1576657634` within `1e-9`. Traces: NFR-M-04.
- [x] 7.4 Add the `mustLoad` panic-path test: `defer recover()` around `mustLoad(domain.LoadConstantes([]byte("{{")))`. Traces: NFR-M-04 (only otherwise-uncoverable branch).

## Phase 8: Verification gate

- [x] 8.1 `go build ./... && go vet ./...` — clean.
- [x] 8.2 `go test ./internal/domain/motor/... -count=1 -v` — all subtests pass.
- [x] 8.3 `go test ./internal/domain/motor/... -cover` — ≥ 90% (NFR-M-04).
- [x] 8.4 `go list -deps ./internal/domain/... | grep -E "internal/(usecase|adapters|infrastructure)"` — empty output (NFR-M-01).
