# Exploration: CalcularCapacidad motor (Issue #8)

## Current State

`internal/domain/motor/` contains only `doc.go`, which states "Zero external dependencies, zero I/O, zero LLM calls" — consistent with NFR-M-01, though it carries no explicit import allowlist. Neither `capacidad.go` nor `capacidad_test.go` exists yet.

`internal/domain/` (delivered by issue #6) provides every data type this change consumes. `references/embed.go` exposes the raw embedded JSON; there is **no cached or singleton constants accessor anywhere in the codebase**.

## Numeric Verification

Real constant values in `internal/domain/constantes.go`, `references/constantes.json`, and `references/subsidios_2026.json` match wiki doc 13 exactly:

| Constant | Value |
|---|---|
| `SMMLV2026` | 1750905 |
| `TasaEADefault` | 0.107 |
| `PlazoMesesDefault` | 240 |
| `CuotaIngresoMax` | 0.40 |
| `MedianaVISCOP` | 195000000 |
| Subsidy brackets | `[{2, 30}, {4, 20}]` |

Arithmetic reproduced against doc 13 §2.3, using annuity factor `102.1576657634` and round-half-up:

| Case | Quota | Product | Rounded |
|---|---|---|---|
| CAP-1 | 0.40 × 2,600,000 = 1,040,000 | 106,243,972.393936 | **106,243,972** ✓ |
| CAP-2 | 800,000 | 81,726,132.611 | **81,726,133** ✓ |
| CAP-3 | 1,200,000 | 122,589,198.916 | **122,589,199** ✓ |
| CAP-4 | 2,000,000 | 204,315,331.527 | **204,315,332** ✓ |
| CAP-5 | 3,200,000 | 326,904,530.442 | **326,904,530** ✓ |

Every `PresupuestoMax` sum reproduces doc 13 exactly.

**Conclusion: there is no numeric drift between the issue's test table and doc 13.** The real risk in this change is code-symbol drift, not arithmetic.

## Affected Areas

- `internal/domain/motor/capacidad.go`, `internal/domain/motor/capacidad_test.go` — to be created; the only in-scope write targets.
- `internal/domain/capacidad.go`, `perfil.go`, `constantes.go`, `enums.go` — read-only dependencies, all confirmed present.
- `references/embed.go`, `references/*.json` — the motor must load constants itself via `domain.LoadConstantes(references.ConstantesJSON)` and `domain.LoadSubsidios(references.SubsidiosJSON)`.

## Q1 — Symbol drift: issue sketch vs. actual issue #6 code

The issue body was written before #6 landed. Its reference sketch does **not compile** against the delivered domain.

| Symbol used by sketch | Exists? | Reality |
|---|---|---|
| `domain.Perfil`, `domain.CampoPerfil` | yes | exact match |
| `domain.FuenteVerificadoBase` | **NO** | real constant is `domain.FuenteCampoVerificadoBase` (likewise `FuenteCampoDeclarado`, `FuenteCampoInferido`) |
| `domain.Capacidad`, `domain.ItemDesglose` | yes | field and JSON tags match Contract §2.2 exactly |
| `domain.TramoSubsidio` | yes | `IngresoHastaSMMLV int`, `SubsidioSMMLV int` |
| `domain.CamposCriticos` | yes, **wrong shape** | it is `map[string]bool`, not a slice — see Q1b |
| `Entero` / `Booleano` / `Texto` / `EsVerificado` | yes | exact signatures |
| `domain.Cts()` | **NO** | only `LoadConstantes([]byte) (Constantes, error)` |
| `domain.TramosSubsidio()` | **NO** | only `LoadSubsidios([]byte) ([]TramoSubsidio, error)` |
| `domain.Campo*` key constants | **NO** | profile keys exist only as raw string literals |
| `ParametrosCredito` | **NO** | new motor-owned type, to be created by this change |
| `FactorAnualidad(...)` | **NO** | correctly absent; new motor-only math |

### Q1b — `CamposCriticos` iteration bug in the sketch

The sketch's confidence loop is:

```go
for _, campo := range domain.CamposCriticos {
    c, ok := p[campo]
    ...
}
```

`CamposCriticos` is declared as `map[string]bool`. Ranging with `for _, campo` binds `campo` to the **value** (`bool`), not the key, so `p[campo]` indexes a `map[string]CampoPerfil` with a `bool`. This is a compile error, not merely a style issue. The correct form iterates keys: `for campo := range domain.CamposCriticos`.

Because `CamposCriticos` is a map, iteration order is non-deterministic — harmless here since confidence multiplication is commutative, but the design must not introduce order-dependent logic over it.

## Q2 — Constants and arithmetic

Confirmed: existing #6 constants reproduce doc 13 exactly. See the Numeric Verification table above.

## Q3 — The `esAfiliado` magic key

**Reject the magic key. Use a plain `esAfiliado bool` parameter.**

The public `CalcularCapacidad` signature already accepts `esAfiliado bool` explicitly. The `"__afiliado_base"` key exists solely to smuggle that value one level down into the **private** `subsidioAplicable` helper, via a `clonarConAfiliado` function that copies the entire `Perfil` map to inject one fake entry.

This is a design smell on four counts:

1. It pollutes `Perfil` with a key absent from both `CamposReconocidos` and `CamposCriticos`, so the domain type no longer means what its own schema says.
2. It allocates a full map copy on every call to avoid changing a signature.
3. The helper is private and in the same file — changing its signature costs nothing and breaks no caller.
4. The sketch's clone uses the non-existent `domain.FuenteVerificadoBase`, so it would not compile anyway.

Recommended signature:

```go
func subsidioAplicable(p domain.Perfil, esAfiliado bool, smmlv int, tramos []domain.TramoSubsidio) (int64, bool)
```

## Q4 — Test conventions

Two divergent conventions exist in the repo:

- `internal/domain/*_test.go` — strict table-driven, `t.Run(tt.name, ...)`, `t.Helper()` assertion helpers, English names.
- `internal/pipeline/*_test.go` — looser `map[string]string` loops without `t.Run`, Spanish failure messages.

`capacidad_test.go` should follow the **`internal/domain` convention**, with subtests named after doc 13 case IDs (`CAP-1_Ana`, `CAP-2`, …) so the suite maps directly onto doc 13 §7's coverage table.

## Q5 — Contract §2.2 conflicts

None. `Capacidad` and `ItemDesglose` already match the contract exactly, including JSON tags.

## Risks

| Severity | Risk |
|---|---|
| CRITICAL | `domain.FuenteVerificadoBase` does not exist; real name is `domain.FuenteCampoVerificadoBase`. |
| CRITICAL | `domain.Cts()`, `domain.TramosSubsidio()`, and `ParametrosCredito` do not exist anywhere. |
| CRITICAL | `CamposCriticos` is `map[string]bool`; the sketch's `for _, campo := range` does not compile. |
| CRITICAL | `"__afiliado_base"` magic profile key — reject in favour of an explicit `esAfiliado bool` parameter. |
| MEDIUM | The issue's test table stops at CAP-7; doc 13 §7 additionally requires CAP-8 and B-1/B-4/B-5/B-6. |
| MEDIUM | CAP-6 ("ingreso 0" boundary) leaves `afiliado` / `tiene_vivienda` / `recibio_subsidio` ambiguous; the spec must pin them. |
| LOW | No `Campo*` key constants exist, so the motor must use raw string literals for profile keys. |

## Recommendation

Proceed to `sdd-propose`. The proposal and design must resolve: the constant-loading strategy (package-level init vs. injected parameter), the corrected `FuenteCampo*` names, the explicit `esAfiliado bool` parameter, key-based iteration over `CamposCriticos`, and an expanded test plan covering doc 13's full CAP-1..8 plus B-1/B-4/B-5/B-6 set.

**Ready for proposal: yes.**
