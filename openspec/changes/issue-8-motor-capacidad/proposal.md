# Proposal — `CalcularCapacidad` deterministic engine (Issue #8)

Give Vivi one function that turns a lead profile into the money numbers the whole product quotes: max credit, applicable subsidy, max budget, ratio, and a confidence score. Today the engine package is an empty `doc.go`, so every downstream feature that needs a budget figure is blocked, and any number shown to a lead would have to be invented by the LLM — the exact failure this system exists to prevent.

---

## 1. Intent

| Question | Answer |
|---|---|
| What problem | `internal/domain/motor/` has zero logic. Nothing can compute a lead's budget, so nothing can route, recommend, or quote. |
| Why now | Issue #6 landed all the domain types. #8 is the first consumer and the hard dependency of #9 (kNN twin), #10 (2x2 matrix) and the advisor ficha. |
| Why it matters commercially | Every peso Vivi ever shows a lead must come from this function. If it drifts from wiki doc 13, the product quotes wrong subsidies to real families. |
| What success looks like | `go test ./internal/domain/motor/...` reproduces every doc 13 §2.3 figure **exactly**, with ≥ 90 % package coverage (NFR-M-04). |
| Non-negotiable | Fully deterministic: zero LLM, zero runtime I/O, no imputed data. Only `DECLARADO` / `VERIFICADO_BASE` values from the live lead feed the math. |

---

## 2. Scope

### In scope

- `internal/domain/motor/capacidad.go` — new.
- `internal/domain/motor/capacidad_test.go` — new.

Nothing else. No file outside `internal/domain/motor/` is created, edited, or moved.

### Out of scope

| Excluded | Belongs to |
|---|---|
| `GemeloKNN`, Gower distance, imputation | Issue #9 |
| `RecomendarProyectos`, `Matriz2x2`, route thresholds, B-2 | Issue #10 |
| Picking the "best affordable catalog project" as the ratio denominator (correction C2) | Issue #10 — #8 receives the price as a parameter |
| Adding `Campo*` key constants or accessors to `internal/domain` | Out of scope: would edit `internal/domain/perfil.go` / `constantes.go` |
| Contract v1.1 §2.2 amendments | Contract change protocol, doc 10 §9 |
| Any usecase, adapter, infrastructure, frontend, or Docs change | Later issues |

---

## 3. Approach

### 3.1 Public surface

```go
func CalcularCapacidad(p domain.Perfil, esAfiliado bool, precioCandidato int64) domain.Capacidad
```

One exported function, no error return, no new exported types. It delegates to an unexported core that receives the already-loaded constants, which is the entire test seam and the entire future-injection path:

```go
func calcular(p domain.Perfil, esAfiliado bool, precioCandidato int64,
              cts domain.Constantes, tramos []domain.TramoSubsidio) domain.Capacidad
```

Math, per doc 13 §2.1: `i = (1+tasa_ea)^(1/12) − 1`, `factor = (1 − (1+i)^−n)/i`, everything in full `float64` precision, **one** round-half-up to integer pesos on each final concept — never on intermediate steps.

### 3.2 The seven decisions

#### D1 — Constant loading: package-level vars in `motor`, loaded from `references`

**Position.** `capacidad.go` declares two package-level vars initialized once from the embedded JSON, via `domain.LoadConstantes(references.ConstantesJSON)` and `domain.LoadSubsidios(references.SubsidiosJSON)`, panicking on parse failure. The public function stays error-free; `calcular` takes the values as parameters so tests can drive it directly.

**Rejected alternatives.**

| Option | Why not |
|---|---|
| Inject `cts` / `tramos` through the public signature | Pushes a load-and-handle-error burden onto every future caller (#10, #12) for data that is compile-time constant. Each caller would then build the very `domain.Cts()` singleton that #6 deliberately did not ship — the same code, just duplicated N times. |
| Add `domain.Cts()` / `domain.TramosSubsidio()` accessors | Edits `internal/domain/constantes.go`, which state.yaml puts out of scope. Also re-litigates a #6 decision inside an #8 PR. |
| `sync.Once` | Go already initializes package vars exactly once and race-free. Unrequested machinery. |

**Does importing `references` violate NFR-M-01? No — and the state.yaml wording is wrong, not the design.** The rule is mechanically enforced in `.github/workflows/ci.yml:33`:

```
go list -deps ./internal/domain/... | grep -E "vivi-perfilamiento-leads/internal/(usecase|adapters|infrastructure)"
```

`references` is a root-level leaf package whose only import is `embed`; it matches none of the forbidden prefixes, so the CI gate passes. `state.yaml`'s paraphrase — *"motor imports only internal/domain"* — is stricter than the rule it claims to restate. **Recommendation: correct the paraphrase, do not contort the design around it.** If the orchestrator insists on the literal paraphrase, the fallback is signature injection (D1's first rejected option) and the public signature grows two parameters; call that out before spec if so.

**Malformed JSON.** `must`-style panic at package init. Justification: the JSON is compiled into the binary, so a parse failure is a *build* defect, not a runtime input error; `internal/domain/constantes_test.go` already fails first in CI if the files break. A panic is loud; the alternative — a zeroed `Constantes` — would silently produce `credito_max = 0` for every lead. Silently wrong money is the worst possible outcome here.

#### D2 — Corrected symbol names

Use `domain.FuenteCampoVerificadoBase`, `domain.FuenteCampoDeclarado`, `domain.FuenteCampoInferido` (`internal/domain/enums.go:32-34`). The issue sketch's `domain.FuenteVerificadoBase` / `domain.FuenteDeclarado` do not exist and will not compile. The sketch predates #6; **the delivered code wins over the sketch on every symbol conflict**, and the same applies to the phantom `domain.Cts()`, `domain.TramosSubsidio()`, and `ParametrosCredito`.

#### D3 — `esAfiliado` is a plain parameter; the `"__afiliado_base"` magic key is rejected

```go
func subsidioAplicable(p domain.Perfil, esAfiliado bool, cts domain.Constantes,
                       tramos []domain.TramoSubsidio) (int64, string)
```

No `clonarConAfiliado`, no synthetic profile entry, no map copy. Reasons for the rejection, recorded for the audit trail:

1. `"__afiliado_base"` is absent from `CamposReconocidos` and `CamposCriticos`, so the `Perfil` map would stop meaning what its own schema declares — anything that ranges the full profile (serialization, the advisor ficha, a "what we know about you" report) would leak a fake field to a human.
2. It clones the whole map on every call purely to avoid changing a private helper's signature.
3. `subsidioAplicable` is private and lives in the same file. Adding one `bool` parameter costs nothing and breaks no caller.
4. The sketch's clone uses `domain.FuenteVerificadoBase`, which does not exist — it never compiled anyway.

#### D4 — Key-based iteration over `CamposCriticos`

`domain.CamposCriticos` is `map[string]bool`, so the sketch's `for _, campo := range domain.CamposCriticos` binds `campo` to a `bool` and then indexes `Perfil` with it — a compile error. Correct form:

```go
for campo := range domain.CamposCriticos { ... }
```

Map iteration order is non-deterministic. Safe here **only because** the confidence rule is a product of commutative factors. The implementation MUST NOT introduce any order-dependent logic over this map (no first-match, no accumulating slice, no early break).

#### D5 — Test scope: doc 13 §7 in full, not the issue's truncated table

The issue lists CAP-1..7. Doc 13 §7 maps `CalcularCapacidad` to **CAP-1..8, B-1, B-4, B-5, B-6** — twelve IDs, all in scope. Two pairs collapse to a single fixture each (CAP-6 = B-1, CAP-8 = B-4), so twelve IDs are covered by ten subtests; each merged subtest is named for both IDs so doc 13 §7 stays traceable line by line.

| Case | Assertion | Note |
|---|---|---|
| CAP-1 (Ana) | credit 106 243 972, subsidy 52 527 150, budget 166 771 122, confidence 0.85 | ratio vs Monguí 156 470 000 |
| CAP-2 | credit 81 726 133, subsidy 52 527 150, budget 134 253 283 | |
| CAP-3 | credit 122 589 199, subsidy 52 527 150, budget 175 116 349 | |
| CAP-4 | credit 204 315 332, subsidy 35 018 100, budget 239 333 432 | middle bracket |
| CAP-5 | credit 326 904 530, subsidy 0, budget 326 904 530 | > 4 SMMLV |
| CAP-6 / B-1 | credit 0, subsidy 0, budget 3 000 000 | inputs pinned in D6; same scenario, one subtest named for both IDs |
| CAP-7 | `tiene_vivienda = true` → subsidy 0, budget = credit + savings | |
| CAP-8 / B-4 | all four critical fields `DECLARADO` → `0.85⁴ = 0.52200625` (minimum reachable; the 0.5 floor cannot bite with only four critical fields — doc 13 v1.1 correction) | same assertion, one subtest named for both IDs |
| B-5 | income > 4 SMMLV with every other condition satisfied → subsidy 0 | proves income alone disqualifies |
| B-6 | `precioCandidato <= 0` → ratio denominator falls back to `mediana_vis_cop` (195 000 000) | |

Out of scope: B-2 (ratio threshold → #10), B-3 (kNN → #9).

Convention: `internal/domain`'s style — table-driven with `t.Run`, `t.Helper()` assertion helpers, English messages, subtests named after doc 13 case IDs so a reviewer maps the suite onto §7 line by line.

#### D6 — CAP-6 inputs, pinned (CAP-6 and B-1 are one scenario)

Doc 13's CAP-6 row is *"income 0, savings 3 000 000 → subsidy 0, credit 0, budget 3 000 000"*. Doc 13 §6 B-1 describes the same input (`ingreso_hogar = 0`) as a boundary property: *"CréditoMax = 0; PresupuestoMax = subsidio + ahorro; sin división por cero."* **These are the same case**, stated twice — once as a value row, once as a property. Doc 13 §7 lists both IDs against `CalcularCapacidad` for traceability, not because two fixtures exist.

Pinned inputs:

```
ingreso_hogar    absent / 0   (DECLARADO)
recursos_propios = 3_000_000  (DECLARADO)
tiene_vivienda   = false      (DECLARADO)
recibio_subsidio = false      (DECLARADO)
esAfiliado       = true
```

**The subsidy is zeroed by the income guard, not by affiliation.** With no positive `ingreso_hogar` there is no income bracket to look up, so `subsidioAplicable` returns 0 before any eligibility condition is evaluated. This is corroborated twice: doc 13's CAP-6 row gives budget `3 000 000` (= 0 credit + 0 subsidy + 3 000 000 savings), and the issue's own fixture runs CAP-6 with `esAfiliado: true` and still expects `subsidioEsp: 0`.

**Rejected: pinning `esAfiliado = false`.** It reaches the documented `0` by the wrong mechanism, and it forces B-1 to be invented as a separate "eligible variant" (`esAfiliado = true` → subsidy 52 527 150, budget 55 527 150). That budget figure appears nowhere in doc 13 and contradicts CAP-6, which is the authoritative value row for zero income. B-1's *"PresupuestoMax = subsidio + ahorro"* is satisfied by CAP-6 with `subsidio = 0`; it does not assert that the subsidy is non-zero.

Implementation consequence: `subsidioAplicable` MUST guard on income before evaluating eligibility conditions, and MUST return 0 when `ingreso_hogar` is absent or `<= 0`. One subtest covers both IDs, named `CAP-6_B-1_ingreso_cero`.

#### D7 — Rollback plan

New isolated package, **zero callers**. Rollback is `git revert` of the single PR, or deleting the two new files. Nothing regresses because nothing depends on it yet, no data is migrated, no schema changes, no config is touched, and no existing behavior is modified. No feature flag, no staged rollout, no ceremony — stating otherwise would be theatre. This is precisely the moment to get the numbers right, because after #9/#10 land, changing this function *will* be a breaking change.

### 3.3 Supporting business rules the spec must carry

These are not in the seven, but the implementation cannot be written without a position on them.

| # | Rule | Position |
|---|---|---|
| S1 | Eligibility conditions | Follow **doc 13 §2.1's four conditions**: `tiene_vivienda == false`, `recibio_subsidio == false`, household has an affiliate (`esAfiliado` OR `hogar_con_afiliado` OR `caja_externa` present), income ≤ 4 SMMLV. Doc 14 §1 adds `personas_hogar ≥ 1` and the 150-SMMLV property cap; both are **excluded** here — enforcing `personas_hogar ≥ 1` would zero the subsidy in CAP-1..4 (those fixtures do not set it) and break doc 13's authoritative figures, and the property cap depends on a candidate project, which is #10's job. Doc 14 itself defers to doc 13 on conflicts. |
| S2 | Absent boolean condition field | Treated as **not satisfied** → subsidy 0. Conservative by design: Vivi never invents a favourable assumption about a family's eligibility. Confidence separately signals the uncertainty. |
| S3 | Confidence | Start at 1.0; multiply by 0.85 for **each** critical field that is not `VERIFICADO_BASE` (covers `DECLARADO`, `INFERIDO`, and absent); floor 0.5 kept as a defensive bound (unreachable: minimum is 0.85⁴ = 0.52200625 with four critical fields — doc 13 v1.1 correction). Reproduces CAP-1 (0.85) and CAP-8 (0.52200625). The "absent counts as 1.0" alternative is rejected: an empty profile would score confidence 1.0. Reuses the existing `Perfil.EsVerificado` accessor. |
| S4 | `Ratio` precision | Stored as full-precision `float64`, not pre-rounded. Doc 13's `1.07` for Ana is a two-decimal *presentation* of `1.0658…`. The round-half-up rule in the doc 13 preamble is explicitly scoped to *"redondeo a pesos enteros"*, and #10's MAT-3 requires an exact `Ratio >= 0.95` comparison that pre-rounding would corrupt. Tests assert with an epsilon. |
| S5 | `Desglose` | Always emit all three items (`CREDITO`, `SUBSIDIO`, `RECURSOS_PROPIOS`) even at amount 0, so length is predictable. `regla` strings derived from the loaded constants — never hardcoded — because correction C1 already moved the rate once (12 % → 10.7 %) and a stale literal is exactly the drift this change exists to kill. |
| S6 | `ItemDesglose.Fuente` | Derived from the driving profile field's provenance (`CREDITO` and `SUBSIDIO` ← `ingreso_hogar`; `RECURSOS_PROPIOS` ← `recursos_propios`). **Flagged conflict:** Contract §2.2's Ana example shows `CREDITO` with `fuente: DECLARADO` and `SUBSIDIO` with `VERIFICADO_BASE`, although both derive from the same `ingreso_hogar` field, which Ana holds as `VERIFICADO_BASE`. The example is internally inconsistent. Derivation is chosen because §2.1 defines `FuenteCampo` as the provenance of a *profile value*, and the entire confidence mechanism depends on provenance being honest. Not a blocker for #8; raise as a contract note. |

---

## 4. Proposal question round

The executor could not reach the user directly. These need a yes/no or a correction before spec freezes:

1. **D1 / NFR-M-01.** Confirm that correcting `state.yaml`'s paraphrase is acceptable, rather than distorting the public signature to satisfy a wording that CI does not enforce. *(Assumed: yes — `references` import approved.)*
2. **S1 eligibility.** Confirm that `personas_hogar ≥ 1` and the 150-SMMLV property cap stay out of `CalcularCapacidad`. *(Assumed: yes — enforcing them would contradict doc 13's authoritative CAP figures.)*
3. **S2 absent fields.** Confirm the conservative rule: an unknown `tiene_vivienda` / `recibio_subsidio` means "no subsidy yet", not "assume eligible". Product impact: an early-stage lead sees a lower budget until they answer. *(Assumed: yes — Vivi does not invent favourable assumptions.)*
4. ~~**D6 CAP-6.** Confirm `esAfiliado = false` as the pinned reason CAP-6's subsidy is 0.~~ **Resolved by the orchestrator against the sources, not assumed.** CAP-6 and B-1 are one scenario; the subsidy is zeroed by the income guard with `esAfiliado = true`. See D6. The originally proposed `esAfiliado = false` split would have introduced budget `55 527 150`, a figure absent from doc 13 and contradicting CAP-6.
5. **S6 desglose provenance.** Confirm the derived-provenance rule and that the Contract §2.2 example inconsistency is logged as a contract note rather than blocking #8. *(Assumed: yes.)*

Answer, skip, or request a second round — skipping means the assumptions above become the spec.

---

## 5. Risks

| Severity | Risk | Mitigation |
|---|---|---|
| MEDIUM | `state.yaml`'s NFR-M-01 paraphrase contradicts the CI-enforced rule; a reviewer could reject the `references` import on the paraphrase alone. | D1 documents the evidence (`ci.yml:33`) and names the fallback design. |
| MEDIUM | Contract §2.2's `desglose.fuente` example is internally inconsistent, so tests will assert values that differ from the literal contract example. | S6 records the rationale; raise via the doc 10 §9 contract protocol. |
| MEDIUM | Doc 13 §2.3 does not state the source (`fuente`) of every field in CAP-2..CAP-7, only amounts. Confidence assertions for those rows are unconstrained. | Spec pins each fixture's provenance explicitly; only CAP-1 and CAP-8 assert confidence, per doc 13. |
| LOW | No `Campo*` key constants exist, so `capacidad.go` uses raw string literals (`"ingreso_hogar"`, …). A typo produces a silent zero, not a compile error. | Tests over the exact doc 13 values catch any mistyped key immediately. |
| LOW | Non-deterministic `CamposCriticos` map order. | D4 forbids order-dependent logic; confidence is a commutative product. |

---

## 6. Acceptance checklist

- [ ] `internal/domain/motor/capacidad.go` and `capacidad_test.go` are the only files added; nothing else changed.
- [ ] All twelve doc 13 §7 IDs (CAP-1..8, B-1, B-4, B-5, B-6) pass with doc 13's exact integers, across ten subtests (CAP-6 = B-1, CAP-8 = B-4).
- [ ] `subsidioAplicable` guards on income **before** eligibility conditions; absent or `<= 0` income returns 0.
- [ ] No `"__afiliado_base"` key, no `clonarConAfiliado`, no `Perfil` mutation anywhere.
- [ ] Only `domain.FuenteCampo*` names are used; the package compiles.
- [ ] `for campo := range domain.CamposCriticos` — key iteration, no order dependence.
- [ ] `go vet ./...` clean; no `internal/{usecase,adapters,infrastructure}` import from `motor`.
- [ ] Package coverage ≥ 90 % (NFR-M-04).
- [ ] Changed lines under the 400-line review budget.

## 7. Next step

`sdd-spec` and `sdd-design` may run in parallel. Spec pins the Given/When/Then fixtures for all twelve cases (RFC 2119); design records D1's loading mechanism and the annuity-factor precision contract.
