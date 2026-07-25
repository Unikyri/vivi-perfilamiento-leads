# Capacidad Specification

## Purpose

Defines the deterministic financial-capacity engine (`CalcularCapacidad`) that turns a lead's `domain.Perfil` into `domain.Capacidad`: max credit, applicable subsidy, max budget, ratio, confidence, and breakdown. Every figure MUST reproduce wiki doc 13 §2.3 exactly. Zero LLM, zero I/O, zero imputed data — only `DECLARADO` / `VERIFICADO_BASE` / `INFERIDO` values already present on the live lead.

## Requirements

### Requirement: Public API surface

The system MUST expose exactly one exported function:

```go
func CalcularCapacidad(p domain.Perfil, esAfiliado bool, precioCandidato int64) domain.Capacidad
```

It MUST delegate to an unexported seam that accepts pre-loaded constants as parameters, so tests drive the math directly without touching package-level state:

```go
func calcular(p domain.Perfil, esAfiliado bool, precioCandidato int64,
              cts domain.Constantes, tramos []domain.TramoSubsidio) domain.Capacidad
```

No new exported types, no error return.

#### Scenario: Public function delegates to the seam

- GIVEN a `domain.Perfil`, `esAfiliado`, and `precioCandidato`
- WHEN `CalcularCapacidad` is called
- THEN it loads package-level `Constantes`/`[]TramoSubsidio` and returns `calcular(...)`'s result unchanged

### Requirement: Annuity credit calculation

`CreditoMax` MUST use `i = (1+tasa_ea)^(1/12) − 1`, `factor = (1 − (1+i)^−n)/i`, full `float64` precision through every intermediate step, with exactly ONE round-half-up to integer pesos applied to the final `CreditoMax` value only. Reference factor at 10.7% EA / 240 months = `102.1576657634`.

`CreditoMax` MUST be 0 when `ingreso_hogar` is absent, zero, or negative — a lead-facing `Capacidad` SHALL never carry a negative amount. `recursos_propios` MUST be clamped to 0 when negative before entering `PresupuestoMax`. *(Added after the pre-PR reliability review: the original draft let a negative declared income produce a negative credit.)*

#### Scenario: Negative declared amounts never produce negative money

- GIVEN `ingreso_hogar = -5,000,000` (DECLARADO) and `recursos_propios = 3,000,000` (DECLARADO)
- WHEN `CalcularCapacidad` is called
- THEN `CreditoMax = 0`, `SubsidioAplicable = 0`, `PresupuestoMax = 3,000,000`

#### Scenario: Credit reproduces doc 13 §2.3 across income levels (CAP-1..5)

- GIVEN `tasa_ea = 0.107`, `n = 240`, and the incomes below (all `tiene_vivienda=false`, `recibio_subsidio=false`, affiliated)
- WHEN `CalcularCapacidad` is called
- THEN `CreditoMax` and `PresupuestoMax` match exactly:

| Case | ingreso_hogar | ahorro | CreditoMax | Subsidio | PresupuestoMax |
|---|---|---|---|---|---|
| CAP-1 (Ana) | 2,600,000 | 8,000,000 | 106,243,972 | 52,527,150 | 166,771,122 |
| CAP-2 | 2,000,000 | 0 | 81,726,133 | 52,527,150 | 134,253,283 |
| CAP-3 | 3,000,000 | 0 | 122,589,199 | 52,527,150 | 175,116,349 |
| CAP-4 | 5,000,000 | 0 | 204,315,332 | 35,018,100 | 239,333,432 |
| CAP-5 | 8,000,000 | 0 | 326,904,530 | 0 | 326,904,530 |

### Requirement: Subsidy eligibility and bracket lookup

`subsidioAplicable` MUST guard on income FIRST: if `ingreso_hogar` is absent or `<= 0`, it MUST return 0 before evaluating any eligibility condition. Only if income is positive MUST it check, in any order, all four doc 13 §2.1 conditions — `tiene_vivienda == false`, `recibio_subsidio == false`, household has an affiliate (`esAfiliado` OR `hogar_con_afiliado` OR `caja_externa` present), `ingreso_hogar <= 4 SMMLV` — treating an absent boolean condition field as NOT satisfied. If all four hold, it MUST look up the bracket: `<= 2 SMMLV → 30 SMMLV (52,527,150)`; `> 2 and <= 4 SMMLV → 20 SMMLV (35,018,100)`; `> 4 SMMLV → 0`.

#### Scenario: CAP-6_B-1 — zero income zeroes subsidy via the income guard, not affiliation

- GIVEN `ingreso_hogar` **present with value 0** (DECLARADO), `recursos_propios = 3,000,000` (DECLARADO), `tiene_vivienda = false` (DECLARADO), `recibio_subsidio = false` (DECLARADO), `esAfiliado = true`
- WHEN `CalcularCapacidad` is called
- THEN `CreditoMax = 0`, `SubsidioAplicable = 0` (returned by the income guard, before any eligibility condition runs), `PresupuestoMax = 3,000,000`

> The fixture pins `ingreso_hogar` as **present with value 0**, not absent. An absent key would make the S6 provenance read `p["ingreso_hogar"].Fuente` yield the zero value `FuenteCampo("")`, which is off-enum for Contract §2.1. Present-with-0 keeps every emitted `Fuente` on-enum while still tripping the income guard. The guard itself MUST still handle the absent case (`!ok`) — see the requirement above.

#### Scenario: CAP-7 — home ownership disqualifies regardless of income bracket

- GIVEN `ingreso_hogar = 2,000,000`, `tiene_vivienda = true`
- WHEN `CalcularCapacidad` is called
- THEN `SubsidioAplicable = 0` and `PresupuestoMax = CreditoMax + recursos_propios`

#### Scenario: B-5 — income above 4 SMMLV disqualifies even when every other condition is satisfied

- GIVEN `ingreso_hogar` > 4 SMMLV, `tiene_vivienda = false`, `recibio_subsidio = false`, household affiliated
- WHEN `CalcularCapacidad` is called
- THEN `SubsidioAplicable = 0`

### Requirement: Confidence score

Confidence MUST start at `1.0` and MUST multiply by `0.85` for each key in `domain.CamposCriticos` whose profile value is absent OR whose `Fuente` is not `domain.FuenteCampoVerificadoBase` (i.e. `DECLARADO`, `INFERIDO`, or missing all count). The result MUST floor at `0.5`. The implementation MUST iterate `for campo := range domain.CamposCriticos` (key-based) and MUST NOT introduce order-dependent logic over this map.

#### Scenario: CAP-1 — one DECLARADO critical field

- GIVEN only `recursos_propios` is `DECLARADO`; the other three critical fields are `VERIFICADO_BASE`
- WHEN `CalcularCapacidad` is called
- THEN `Confianza = 0.85`

#### Scenario: CAP-8_B-4 — all four critical fields DECLARADO yields the minimum reachable confidence

- GIVEN all four `CamposCriticos` entries are `DECLARADO`
- WHEN `CalcularCapacidad` is called
- THEN `Confianza = 1.0 * 0.85^4 = 0.52200625` (epsilon 1e-9)

> The `0.5` floor MUST remain in the implementation as a defensive lower bound, but it is **unreachable** with the current parameters: exactly four critical fields exist, so the minimum product is `0.85^4 = 0.52200625 > 0.5`. Doc 13's original CAP-8 line ("0.522 → floor 0.5 applies") contradicted its own arithmetic — a floor cannot lower 0.522 to 0.5 — and was corrected in the wiki (doc 13 v1.1). Asserting `0.5` here would require inventing a clamp rule that exists nowhere.

### Requirement: Ratio computation

`Ratio` MUST be `PresupuestoMax / precioCandidato` stored as full-precision `float64`, never pre-rounded. WHEN `precioCandidato <= 0`, the denominator MUST fall back to `cts.MedianaVISCOP` (195,000,000).

#### Scenario: B-6 — non-positive candidate price falls back to the median VIS

- GIVEN `precioCandidato = 0`
- WHEN `CalcularCapacidad` is called
- THEN `Ratio = PresupuestoMax / 195,000,000`, computed at full `float64` precision

### Requirement: Desglose breakdown

`Desglose` MUST always contain exactly three `ItemDesglose` entries, in this concept order: `CREDITO`, `SUBSIDIO`, `RECURSOS_PROPIOS` — even when an amount is 0. Each `Regla` string MUST be derived from the loaded `Constantes`/`[]TramoSubsidio` values, never a hardcoded literal. `Fuente` MUST be derived from the driving profile field's provenance: `CREDITO` and `SUBSIDIO` from `ingreso_hogar`'s `Fuente`; `RECURSOS_PROPIOS` from `recursos_propios`'s `Fuente`.

WHEN the driving profile field is absent, the emitted `Fuente` MUST be the empty `domain.FuenteCampo("")`. The implementation MUST NOT substitute `DECLARADO`, `INFERIDO`, or any other value: the line's amount is 0 precisely because no source field exists, and stamping a provenance would invent an origin for data the lead never supplied. An empty `Fuente` on a `monto: 0` line is self-consistent. This is logged as a contract note against Contract §2.1's enum and does not block this change.

#### Scenario: Three items always emitted with derived rule text and provenance

- GIVEN any valid `Perfil` (including CAP-6_B-1 where `SubsidioAplicable = 0`)
- WHEN `CalcularCapacidad` is called
- THEN `len(Desglose) == 3`, `Regla` strings embed the current `tasa_ea`/`cuota_ingreso_max`/bracket values, and each `Fuente` matches its driving field's `Fuente` on the input `Perfil`
