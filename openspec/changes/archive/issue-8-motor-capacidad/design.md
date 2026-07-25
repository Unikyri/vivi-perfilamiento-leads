# Design — `CalcularCapacidad` deterministic engine (Issue #8)

Two new files in `internal/domain/motor/`. One exported function, five unexported helpers, one generic `must` wrapper. All arithmetic in `float64` at full precision with **exactly one** `math.Round` in the whole package. Nothing outside `internal/domain/motor/` is touched.

Proposal decisions D1–D7 and rules S1–S6 are inputs, not open questions. This document records **how**.

---

## 1. File layout — `internal/domain/motor/capacidad.go`

Declaration order (Go convention: consts, vars, exported, unexported in call order).

```go
package motor

import (
    "fmt"
    "math"

    "github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
    "github.com/Unikyri/vivi-perfilamiento-leads/references"
)

const (
    penalizacionCampoNoVerificado = 0.85 // doc 13 §2.1
    confianzaMinima               = 0.50 // doc 13 §2.1 floor
)

var (
    constantes     = mustLoad(domain.LoadConstantes(references.ConstantesJSON))
    tramosSubsidio = mustLoad(domain.LoadSubsidios(references.SubsidiosJSON))
)

// 1. Public entry point.
func CalcularCapacidad(p domain.Perfil, esAfiliado bool, precioCandidato int64) domain.Capacidad

// 2. Test seam + future injection path.
func calcular(p domain.Perfil, esAfiliado bool, precioCandidato int64,
    cts domain.Constantes, tramos []domain.TramoSubsidio) domain.Capacidad

// 3. Pure annuity math.
func factorAnualidad(tasaEA float64, plazoMeses int) float64

// 4. Income guard + four eligibility conditions + bracket lookup.
//    Returns the amount in COP and the derived `regla` string (S5).
func subsidioAplicable(p domain.Perfil, esAfiliado bool,
    cts domain.Constantes, tramos []domain.TramoSubsidio) (int64, string)

// 5. Commutative product over CamposCriticos, floored.
func calcularConfianza(p domain.Perfil) float64

// 6. Three fixed items, provenance derived per S6.
func armarDesglose(p domain.Perfil, cts domain.Constantes,
    credito, subsidio, recursos int64, reglaSubsidio string) []domain.ItemDesglose

// 7. Init-time panic wrapper.
func mustLoad[T any](v T, err error) T
```

### Why each unexported helper earns its place

| Helper | Justification | Verdict |
|---|---|---|
| `calcular` | Only way a test can prove the output *derives* from the loaded constants instead of hardcoded literals: pass a different `cts`, the numbers must move. Also the D1 injection path for #10/#12 without changing the public signature. | Keep |
| `factorAnualidad` | The single most defect-prone expression in the change. Directly assertable against doc 13's published `102.1576657634`, without building a `Perfil`. | Keep |
| `subsidioAplicable` | Branchiest logic in the file: one guard + four conditions + bracket scan. Inlining it would push `calcular` past 60 lines. | Keep |
| `calcularConfianza` | Owns the two doc-13 literals and the floor. Isolated so CAP-8/B-4 asserts one function, not a whole pipeline. | Keep |
| `armarDesglose` | Owns the two rules most likely to drift: S5 (regla strings from constants) and S6 (provenance derivation, which already conflicts with Contract §2.2). One place for the contract note. | Keep — marginal; if it lands under 15 lines with no branches, inlining is acceptable |
| `mustLoad[T]` | One generic wrapper beats two type-specific `mustConstantes` / `mustSubsidios`. Go 1.25 (`go.mod:3`) supports it; `mustLoad(domain.LoadConstantes(x))` is legal — a multi-value call may be the complete argument list. | Keep |

**Deleted before it existed:** no `clonarConAfiliado`, no `ParametrosCredito` type, no `redondearPesos` helper (one call site — `math.Round` inline is shorter and clearer).

---

## 2. Constant loading (D1)

Two package-level vars, initialized once by the Go runtime, race-free, no `sync.Once`.

**Panic beats returning an error.** The JSON is compiled into the binary by `references/embed.go`, so a parse failure is a **build defect**, not a runtime input error. `internal/domain/constantes_test.go` already fails first in CI if the files break. The alternative — degrading to a zeroed `Constantes` — silently yields `CuotaIngresoMax = 0` and therefore `credito_max = 0` for **every lead**. Silently wrong money is strictly worse than a loud crash at process start.

**Importing `references` does not violate NFR-M-01.** The rule is mechanically enforced at `.github/workflows/ci.yml:33`:

```
go list -deps ./internal/domain/... | grep -E "vivi-perfilamiento-leads/internal/(usecase|adapters|infrastructure)"
```

`references` resolves to `github.com/Unikyri/vivi-perfilamiento-leads/references` — a root-level leaf package whose only import is `embed` (`references/embed.go:5`). It matches none of the three forbidden prefixes; the gate passes. `state.yaml`'s paraphrase ("motor imports only internal/domain") is stricter than the rule it claims to restate. **Correct the paraphrase, do not contort the design.**

---

## 3. Precision contract

This is the section the tests exist to defend.

### Where float64 flows, and where it stops

| Step | Expression | Type | Rounded? |
|---|---|---|---|
| Monthly rate | `i = math.Pow(1+cts.TasaEADefault, 1.0/12.0) - 1` | float64 | no |
| Annuity factor | `(1 - math.Pow(1+i, -float64(n))) / i` | float64 | **no** |
| Quota | `cts.CuotaIngresoMax * float64(ingreso)` | float64 | **no** |
| **CréditoMax** | `int64(math.Round(cuota * factor))` | int64 | **YES — the only round** |
| Subsidy | `int64(tramo.SubsidioSMMLV) * int64(cts.SMMLV2026)` | int64 | n/a — never a float |
| Bracket test | `ingreso <= int64(t.IngresoHastaSMMLV) * int64(cts.SMMLV2026)` | int64 | n/a — never a float |
| PresupuestoMax | `credito + subsidio + recursos` | int64 | n/a — sum of settled ints |
| Ratio | `float64(presupuesto) / float64(denominador)` | float64 | **no** (S4) |
| Confianza | product, then `math.Max(c, confianzaMinima)` | float64 | **no** |

`int64(math.Round(x))` is correct here because every amount is non-negative: Go's `math.Round` rounds half **away from zero**, which for `x >= 0` is exactly the round-half-up mandated by the doc 13 preamble. A negative credit is unreachable — `cuota` is `0.40 × ingreso` with `ingreso > 0` guaranteed by the income guard, and `factor` is positive for any `tasaEA > 0`.

### CAP-1 trace

```
i        = 1.107^(1/12) - 1        = 0.008507088...
factor   = (1 - (1+i)^-240) / i    = 102.1576657634
cuota    = 0.40 x 2_600_000        = 1_040_000        (exact, NOT rounded)
producto = 1_040_000 x factor      = 106_243_972.393936
credito  = round(producto)         = 106_243_972       <- doc 13 §2.3
presupuesto = 106_243_972 + 52_527_150 + 8_000_000 = 166_771_122
```

**The defect these tests catch.** Rounding an intermediate produces wrong pesos. Truncating the factor to `102.15` gives `1_040_000 × 102.15 = 106_236_000` — **7,972 pesos off** on a single lead, and it drifts further as income rises. Pre-rounding `cuota` corrupts every case where `0.40 × ingreso` is not an integer. Neither failure is visible by inspection; only the exact doc 13 integers expose them. That is the entire reason CAP-1..5 assert exact `int64` equality and not a tolerance.

**Subsidy and bracket comparison are integer-only, deliberately.** Doc 13 §2.1 phrases the bracket rule as "`ingreso_hogar` en SMMLV", which invites `float64(ingreso)/float64(smmlv) <= 2.0`. That is float-fragile exactly at the boundary. Multiplying instead of dividing keeps it exact: 2 SMMLV = `3_501_810`, 4 SMMLV = `7_003_620`, both integers. `30 × 1_750_905 = 52_527_150` needs no rounding at all.

---

## 4. Subsidy evaluation order (D6, S1, S2)

```
subsidioAplicable:
  1. ingreso, ok := p.Entero("ingreso_hogar")
     if !ok || ingreso <= 0        -> return 0, regla("sin ingreso declarado")   <-- GUARD FIRST
  2. tiene_vivienda   != false     -> return 0        (absent => not satisfied, S2)
     recibio_subsidio != false     -> return 0
     no affiliate route            -> return 0        (esAfiliado || hogar_con_afiliado || caja_externa)
     ingreso > 4 SMMLV             -> return 0        (B-5)
  3. first tramo where ingreso <= IngresoHastaSMMLV * SMMLV2026
     -> return SubsidioSMMLV * SMMLV2026, derived regla
```

**Why the guard must come first.** CAP-6 / B-1 pins `esAfiliado = true` with `ingreso_hogar` absent-or-zero, and doc 13 still requires `subsidio = 0`, `budget = 3_000_000`. If the four conditions ran first, this profile would satisfy all of them (`tiene_vivienda = false`, `recibio_subsidio = false`, affiliate present, `0 <= 4 SMMLV`) and fall into the first bracket, returning **52,527,150** — a budget of 55,527,150 that appears nowhere in doc 13 and contradicts its authoritative CAP-6 row. The guard is what makes CAP-6 and B-1 one scenario instead of two contradictory ones.

Secondary effect: with no positive income there is no bracket to look up, so the guard also removes any meaningless bracket scan. Division by zero never arises — the annuity uses `i`, never `ingreso`, and the Ratio denominator falls back to `mediana_vis_cop` (B-6).

---

## 5. Confidence (D4, S3)

```go
c := 1.0
for campo := range domain.CamposCriticos {   // KEY iteration
    if !p.EsVerificado(campo) {
        c *= penalizacionCampoNoVerificado
    }
}
return math.Max(c, confianzaMinima)
```

`domain.CamposCriticos` is `map[string]bool` (`internal/domain/perfil.go:40`). The issue sketch's `for _, campo := range domain.CamposCriticos` binds `campo` to the map **value** (`bool`) and then indexes a `map[string]CampoPerfil` with it — a **compile error**, not a style preference.

Map iteration order is non-deterministic. That is safe here **only because** the body is a product of commutative factors: the same multiset of `0.85` factors yields a bit-identical `float64` regardless of visit order (IEEE-754 multiplication is commutative, and all factors are identical anyway).

**Prohibition, binding on implementation and review:** no first-match, no early `break`, no accumulating slice, no index arithmetic, no sorted-key workaround over `CamposCriticos`. Any order-dependent construct here is a defect even if tests pass, because it passes non-deterministically.

`EsVerificado` returns `false` for absent keys (`perfil.go:87-90`), so "absent" is penalised identically to `DECLARADO`/`INFERIDO` — exactly S3. The rejected "absent counts as 1.0" alternative would score an empty profile at confidence 1.0.

The two literals live as named `motor` consts because neither exists in `constantes.json`; adding them there would edit `references/`, which is out of scope. Follow-up, not a blocker.

---

## 6. Test architecture — `capacidad_test.go`

Follows `internal/domain`'s convention (`constantes_test.go`, `perfil_test.go`): table-driven, `t.Run(tt.nombre, ...)`, `t.Helper()` assertion helpers, English helper names over Spanish domain terms.

### Fixture shape

```go
type casoCapacidad struct {
    nombre          string        // doc 13 ID, e.g. "CAP-6_B-1_ingreso_cero"
    perfil          domain.Perfil
    esAfiliado      bool
    precioCandidato int64
    wantCredito     int64
    wantSubsidio    int64
    wantPresupuesto int64
    wantRatio       *float64      // nil => not asserted
    wantConfianza   *float64      // nil => not asserted
}

func newCampo(valor any, fuente domain.FuenteCampo) domain.CampoPerfil  // t.Helper not needed
func assertClose(t *testing.T, nombre string, got, want float64)        // epsilon 1e-9
```

`*float64` rather than a `0.0` sentinel: doc 13 constrains provenance only for CAP-1 and CAP-8, so most rows must **not** assert confidence — and `0.0` is a value a buggy implementation can legitimately produce, so it cannot double as "skip".

### Ten subtests, twelve doc 13 §7 IDs

| Subtest | IDs | Asserts |
|---|---|---|
| `CAP-1_Ana` | CAP-1 | credit 106 243 972 · subsidy 52 527 150 · budget 166 771 122 · ratio vs 156 470 000 · confianza 0.85 |
| `CAP-2` | CAP-2 | 81 726 133 / 52 527 150 / 134 253 283 |
| `CAP-3` | CAP-3 | 122 589 199 / 52 527 150 / 175 116 349 |
| `CAP-4_tramo_medio` | CAP-4 | 204 315 332 / 35 018 100 / 239 333 432 |
| `CAP-5_mas_4_smmlv` | CAP-5 | 326 904 530 / 0 / 326 904 530 |
| `CAP-6_B-1_ingreso_cero` | CAP-6, **B-1** | 0 / 0 / 3 000 000 — income guard, no division by zero |
| `CAP-7_tiene_vivienda` | CAP-7 | subsidy 0, budget = credit + savings |
| `CAP-8_B-4_confianza_minima` | CAP-8, **B-4** | four criticals `DECLARADO` → 0.85⁴ = **0.52200625** (minimum reachable; 0.5 floor is a defensive bound that cannot bite — doc 13 v1.1 correction) |
| `B-5_ingreso_supera_tope` | B-5 | subsidy 0 with every other condition satisfied |
| `B-6_precio_candidato_cero` | B-6 | ratio denominator = `mediana_vis_cop` (195 000 000) |

CAP-6 = B-1 and CAP-8 = B-4 are the same input with the same expectation stated twice in doc 13 (once as a value row, once as a property). Each merged subtest is **named for both IDs** so a reviewer maps §7 line by line without opening the doc.

### Comparison discipline

- **Integers** (`CreditoMax`, `SubsidioAplicable`, `PresupuestoMax`, `RecursosPropios`): exact `!=`. A tolerance here would hide precisely the intermediate-rounding defect from §3.
- **Floats** (`Ratio`, `Confianza`): `assertClose` with epsilon `1e-9`. Doc 13's `1.07` for Ana is a 2-decimal *presentation* of `1.06583…`; the fixture pins the full-precision expectation, not the printed one.
- **Desglose**: assert `len == 3` and the three `Concepto` values on every row (S5 guarantees fixed length even at amount 0).

### Two extra subtests for NFR-M-04 (≥ 90 % coverage)

1. `factorAnualidad(0.107, 240)` → `102.1576657634` within `1e-9`. Pins doc 13 §1's published constant directly.
2. `mustLoad` panic path: `defer recover()` around `mustLoad(domain.LoadConstantes([]byte("{{")))`. Six lines that turn the only otherwise-uncoverable branch green.

---

## 7. Rejected alternatives

| Rejected | Why not |
|---|---|
| `"__afiliado_base"` magic profile key | Injects a key absent from `CamposReconocidos`/`CamposCriticos`, so `Perfil` stops meaning what its own schema declares — and anything that ranges the full profile (serialization, advisor ficha) leaks a fake field to a human. Also clones the whole map per call to avoid changing a private signature. |
| `sync.Once` around constant loading | Go already initializes package vars exactly once, race-free, before `main`. Unrequested machinery. |
| Injecting `cts`/`tramos` through the **public** signature | Pushes a load-and-handle-error burden onto every future caller (#9, #10, #12) for data that is compile-time constant, duplicating the singleton #6 deliberately did not ship. `calcular` already provides the seam without the cost. |
| `domain.Cts()` / `domain.TramosSubsidio()` accessors | Edits `internal/domain/constantes.go`, which is out of scope, and re-litigates a #6 decision inside an #8 PR. |
| A `redondearPesos(float64) int64` helper | One call site. `int64(math.Round(x))` is shorter and states the rule inline. |
| Two `mustConstantes` / `mustSubsidios` helpers | One generic `mustLoad[T]` covers both call sites on Go 1.25. |
| Float bracket comparison `ingreso/smmlv <= 2.0` | Float-fragile exactly at the 2 and 4 SMMLV boundaries. Integer multiply is exact. |

---

## 8. Findings the implementation must not paper over

**F1 — `ItemDesglose.Fuente` has no value for an absent driving field (new, not in the proposal).** S6 derives each item's provenance from its driving profile field (`CREDITO`/`SUBSIDIO` ← `ingreso_hogar`, `RECURSOS_PROPIOS` ← `recursos_propios`). `Perfil` exposes no provenance accessor, so the read is `p["ingreso_hogar"].Fuente` — legal Go for a map read, but a **missing key yields the zero value `FuenteCampo("")`**, which is off-enum for Contract §2.1.

Position: **emit the empty `FuenteCampo`, do not substitute one.** The amount on that line is 0 because no source field exists; stamping it `DECLARADO` or `INFERIDO` would invent provenance for data the lead never gave — the exact failure this engine exists to prevent. An empty `fuente` on a `monto: 0` line is self-consistent. Log as a contract note alongside S6; it does not block #8. Affects CAP-6/B-1 only if the spec pins `ingreso_hogar` as **absent** rather than present-with-value-0; the spec agent should prefer present-with-0-`DECLARADO` so the fixture stays on-enum.

**F2 — S6 conflicts with Contract §2.2's Ana example.** The example shows `CREDITO` as `DECLARADO` and `SUBSIDIO` as `VERIFICADO_BASE` although both derive from the same `ingreso_hogar`, which Ana holds as `VERIFICADO_BASE`. The example is internally inconsistent. Tests assert the derived value. Raise via doc 10 §9.

**F3 — no `Campo*` key constants exist.** `capacidad.go` uses raw string literals (`"ingreso_hogar"`, `"recursos_propios"`, `"tiene_vivienda"`, `"recibio_subsidio"`, `"hogar_con_afiliado"`, `"caja_externa"`). A typo produces a silent zero, not a compile error. Mitigated only by the exact doc 13 integer assertions. Adding constants to `internal/domain` is out of scope for #8.

Nothing in the proposal fails to compile or is numerically wrong. The exploration's four CRITICAL symbol-drift findings are all resolved in this design (`FuenteCampo*` names, no `domain.Cts()`, key iteration, no magic key), and the arithmetic reproduces doc 13 §2.3 exactly — verified above for CAP-1.

---

## 9. Size budget

| Artifact | Est. lines |
|---|---|
| `internal/domain/motor/capacidad.go` | ~140 |
| `internal/domain/motor/capacidad_test.go` | ~210 |
| **Go total** | **~350 — under the 400-line review budget** |

`400-line budget risk: Low` for code. If the `openspec/changes/issue-8-motor-capacidad/` markdown artifacts ship in the same PR they add ~600 more lines; recommend labelling them `docs` or landing them separately so the reviewer's attention stays on the 350 lines of Go that move money.

## 10. Open questions

- [x] **F1 — resolved by the orchestrator.** Both positions accepted and written into the spec: the empty `domain.FuenteCampo("")` is emitted for an absent driving field (never substituted), and CAP-6/B-1 pins `ingreso_hogar` as **present with value 0** (DECLARADO) rather than absent, so every emitted `Fuente` stays on-enum. The income guard MUST still handle the absent case (`!ok`) even though no fixture exercises it.
- [ ] Exact Spanish wording of the three `regla` strings — spec's call. Design mandates only that every numeric token inside them comes from `cts`/`tramos`, never a literal (S5).
