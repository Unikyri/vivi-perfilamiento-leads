```yaml
schema: gentle-ai.verify-result/v1
evidence_revision: sha256:8f1e5126ef573523c80852a139a268f8a0bdbbc797dfd30bf8719849e995e189
verdict: pass
blockers: 0
critical_findings: 0
requirements: 6/6
scenarios: 7/7
test_command: go test ./... -count=1
test_exit_code: 0
test_output_hash: sha256:3b84b79ec6e999dc65fc47ad3e2511e58c17bd3bd3b57601cb3c9902cd555d4d
build_command: go build ./...
build_exit_code: 0
build_output_hash: sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
```

## Verification Report

**Change**: issue-18-perfilar-lead
**Issue**: #18 `[A10] Caso de uso: PerfilarLead (UC-01, US-01/US-02)` (OPEN, label `bloque-a`)
**Worktree / branch**: `/tmp/vivi-issue-18` on `feat/issue-18-perfilar-lead`, HEAD `2ad800a` (= `origin/main`)
**Artifact store**: hybrid (OpenSpec files + Engram observations #548–#555)
**Version**: spec has no version field (N/A)
**Mode**: Standard (`openspec/config.yaml` `strict_tdd: false`; `strict-tdd-verify.md` intentionally not loaded)
**Run**: final independent requirements/runtime verification after remediation, read-only on implementation and tests

### Review Authority (bound, approved)

`./.tools/gentle-ai review validate --gate post-apply --cwd /tmp/vivi-issue-18` → `result: allow`,
`reason: "authoritative transaction, current repository target, and content-bound artifacts match"`.

| Field | Value |
|---|---|
| Lineage | `review-8395941eaa5fec38`, generation 1, terminal state `approved` |
| Risk tier / lenses | high / `review-risk`, `review-resilience`, `review-readability`, `review-reliability` (all zero findings) |
| Base tree | `35b39c067fde0714a39e2fb0096379711e533cac` |
| Candidate tree | `1cf0c727a0def8b78e2da24b01ec6593f794812c` (initial = final) |
| Paths digest | `sha256:3cc769d02614ede96fc641b56991b6ac3b97887094d4594e1c98aed9be4ad9d1` |
| Policy / evidence hash | `sha256:34fb63d7f2…a184f6` / `sha256:fcfddb5ada…d05f9db` |
| Fix delta / ledger hash | empty digest (no correction transaction consumed) |
| Store revision | `sha256:cb6972a8bfb0b7f660241cc42be79b66736e8e6c380831c76244b6feeb27d16e` |
| `sdd-status` reviewGate | `allow` — "explicit bound compact authority exactly matches the current repository" |

The previous report's `authority_only_failure` / `missing_review_authority` preflight denial is therefore
**resolved and withdrawn**: the declared envelope commands were executed for real this time (exit `0`, digests below).

**Toolchain note**: `gentle-ai sdd-verify-validate` does not exist in either installed binary
(`.tools/gentle-ai` 2.1.11 or PATH `gentle-ai` 1.43.3 → `unknown command`). Admission was therefore
self-checked against `sdd-verify/references/report-format.md`: envelope first, every canonical field
exactly once, no recovery fields, counts measured from the retrieved spec (6 requirements / 7 scenarios).
This report is **not** validator-attested.

### Completeness
| Metric | Value |
|--------|-------|
| Tasks total | 12 |
| Tasks complete | 12 |
| Tasks incomplete | 0 |
| Authored implementation lines | `perfilar_lead.go` 178 + `perfilar_lead_test.go` 178 = **356** |
| Hard stop / review budget | 390 / 400 → ✅ 34 lines of headroom |
| Tracked file modifications | none — `git diff` and `git diff --cached` are empty; only untracked new files |
| Implementation scope | exactly the two planned files; no other source, test, fixture, schema, config, or Contract file touched |

### Independent Evidence Executed (read-only)

All commands run by this verifier in `/tmp/vivi-issue-18`. Digests are SHA-256 over exact combined
stdout+stderr bytes.

| Check | Command | Exit | Output digest |
|---|---|---:|---|
| Focused | `go test ./internal/usecase/... -run 'TestPerfilar\|TestReconsulta' -v` | 0 | `sha256:4beb75598ec4ceb94418f45e5d9f91879e5db6e2df6da1dc6fb1f81e5b21e2e0` |
| Full suite | `go test ./... -count=1` | 0 | `sha256:3b84b79ec6e999dc65fc47ad3e2511e58c17bd3bd3b57601cb3c9902cd555d4d` |
| Race | `go test -race ./internal/...` | 0 | `sha256:30fce6c65556756a27eb1d3c6e7b5a038e94dce8ca097c544255f23e6ba7888f` |
| Build | `go build ./...` | 0 | `sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` (empty) |
| Vet | `go vet ./...` | 0 | `sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` (empty) |
| Module verify | `go mod verify` | 0 | `sha256:b4537ed75f533f993f371954de47e42a793b8e5b0587577de7e27fb3e50696bd` (`all modules verified`) |
| gofmt (changed files) | `gofmt -l internal/usecase/perfilar_lead.go internal/usecase/perfilar_lead_test.go` | 0 | empty digest → clean |
| gofmt (repo) | `gofmt -l internal cmd` | 0 | `sha256:f3d3eae605dee7dca8994bdf6a2625ac39b3fa3498266612d6c28c9585bcd9a1` — 2 pre-existing base-only files (see below) |
| Diff / budget | `git status --porcelain`, `git diff --stat`, `wc -l` | 0 | 356 authored lines, no tracked modification |

`evidence_revision` is SHA-256 over the canonical preimage of the six identity fields (change, branch,
HEAD, base tree, candidate tree, lineage) plus the nine `name|command|exit|digest` lines and the
`authored_lines`/`requirements`/`scenarios` totals above.

Focused run detail: `TestPerfilarLeadEjecucion` (4 subtests: `active_affiliate`, `unknown_fallback`,
`inactive_fallback`, `catalog_failure_fallback`), `TestPerfilarLeadErrorsDoNotPublish`,
`TestReconsultarFamiliar`, `TestReconsultarFamiliarNoEncontradoPersisteDeclarado` — 4 test functions,
4 subtests, all PASS, 0 failed, 0 skipped. Full suite: 10 packages `ok`, 4 `[no test files]`, 0 failures.

**Coverage**: ➖ Not available — `openspec/config.yaml` sets `coverage_threshold: 0` and declares no coverage command.

### Pre-existing base-only gofmt drift (explicit classification)

`gofmt -l internal cmd` reports exactly `internal/pipeline/compradores_test.go` and
`internal/pipeline/proyectos_test.go`. Proven **base-only, not attributable to Issue #18**:

- `git diff origin/main --stat -- internal/pipeline/` is empty → both files are byte-identical to the base.
- Their last commit is `64a8022 feat(pipeline): add proyectos.json 16-project catalog (#13)`, before this branch.
- `gofmt -l` on the blobs extracted straight from `origin/main` (`git show origin/main:<path>`) flags both
  files as well, so the drift exists in the protected base tree itself.
- Neither path is in the reviewed candidate path set for this change.

Classification: pre-existing / base-only WARNING, out of scope for #18, follow-up on `main`.
`state.yaml` records this accurately as `gofmt_repo_wide: fail_pre_existing_base_only`; the `gofmt: pass`
field refers to the two changed files only.

### Spec Compliance Matrix
| Requirement | Scenario | Test | Result |
|-------------|----------|------|--------|
| Verified affiliate pre-profile | Ana receives a verified demo pre-profile | `perfilar_lead_test.go > TestPerfilarLeadEjecucion/active_affiliate` | ✅ COMPLIANT |
| Non-affiliate fallback | Unknown or inactive affiliate remains discoverable | `TestPerfilarLeadEjecucion/{unknown_fallback,inactive_fallback,catalog_failure_fallback}` | ✅ COMPLIANT |
| Valid creation, persistence, and event order | Successful pre-profile is durable before notification | `TestPerfilarLeadEjecucion` (persisted `PERFILANDO`, `Version=1`, exactly one `LeadNuevo` with exact payload) | ✅ COMPLIANT |
| Valid creation, persistence, and event order | Create or transition failure has no success event | `TestPerfilarLeadErrorsDoNotPublish` (canceled ctx + `Crear` error → 0 events) | ✅ COMPLIANT |
| Verified household re-consultation | Family match is summed once | `TestReconsultarFamiliar` | ✅ COMPLIANT |
| Unverified family and failure safety | Unknown family requires later confirmation | `TestReconsultarFamiliarNoEncontradoPersisteDeclarado` (declared + `RequiereConfirmacion` + sentinel, plus canceled ctx, `PorID` failure, save/CAS failure) | ✅ COMPLIANT |
| Clean deterministic boundary | Isolated use-case execution | whole focused suite with local `catalogFake`/`busFake`, `NuevoLeadRepoFake`, `NuevoRelojFake(time.Unix(100,0))`, `NuevoIDFake("lead")` | ✅ COMPLIANT |

**Compliance summary**: 7/7 scenarios compliant, 0 PARTIAL, 0 FAILING, 0 UNTESTED (was 5/7 with 2 PARTIAL before remediation).

Remediation closed both prior PARTIALs with runtime assertions:
- Event order/payload: `reflect.DeepEqual(Payload, {cedula, nombre, telefono, fuente})`, `Tipo == EvLeadNuevo`,
  `Evento.LeadID == out.LeadID`, and 0 events on canceled/create-failure paths.
- Family scenario: `hogar_con_afiliado == true`, `EsVerificado(hogar_con_afiliado)`,
  `EsVerificado(cedula_familiar_afiliado)`, `ingreso_hogar == 4_500_000` once,
  recalculated `SubsidioAplicable == 35_018_100`, and an unchanged income/version/capacity/event count on repeat.

The "transition fails" disjunct of scenario 4 is unreachable through the public API: the lead is always minted
in `EstadoLeadNuevo` and `internal/domain/estado.go` allows `NUEVO → PERFILANDO` unconditionally, so the wrapped
`transicionar lead:` branch is dead defensive code, not missing coverage. The `Crear`-failure disjunct is asserted.

### Correctness (Static Evidence, cross-checked against the governing docs)
| Requirement | Status | Notes |
|------------|--------|-------|
| Verified affiliate pre-profile | ✅ Implemented | `perfilAfiliado` stamps exactly the 7 recognized keys with `FuenteCampoVerificadoBase`, `Confianza: 1`, `Reloj` time; `personas_hogar = int64(PersonasACargo+1)` → 3; empty input name backfills `a.Nombre`. Subsidy independently recomputed from `references/`: `smmlv_2026 = 1750905`, brackets `[{≤2 SMMLV → 30}, {≤4 → 20}]`; `2600000 ≤ 3501810` → `30 × 1750905 = 52527150` ✓ (doc 13 CAP-1, issue DoD). |
| Non-affiliate fallback | ✅ Implemented | `lookupActive` returns `(nil, nil)` for empty cedula, catalog miss, `AfiliadoActivo=false`, and non-context lookup error; profile stays `domain.Perfil{}` (non-nil, len 0), `Afiliado=false`, motor income guard → subsidy 0; creation never fails on affiliate availability. |
| Valid creation, persistence, and event order | ✅ Implemented | Minted in `EstadoLeadNuevo`; `motor.CalcularCapacidad(perfil, afiliado, 0)` so the `mediana_vis_cop = 195000000` non-positive-candidate fallback applies; `Lead.Transicionar(EstadoLeadPerfilando)` uses the domain machine; `Crear` strictly precedes the single `Bus.Publicar`; payload matches Contract v1.1 §6 `LeadNuevo { lead_id, cedula?, nombre, telefono, fuente }` with `lead_id` carried by `Evento.LeadID`. Only `AfiliadoPorCedula` is called on `CatalogoRepository` — no catalogue, kNN, or recommendation access. |
| Verified household re-consultation | ✅ Implemented | `PorID` → duplicate guard → `lookupActive` → `ingreso_hogar += IngresoMensual` with verified `hogar_con_afiliado`/`cedula_familiar_afiliado` → recompute → `Guardar`. CAS is real in the shared fake (`stored.Version != lead.Version → ErrOptimisticLock`, clone-on-read/write). Same verified cedula returns `nil` before mutation; a distinct verified cedula returns `ErrFamiliarYaRegistrado`. `4500000 ≤ 4 SMMLV (7003620)` → `20 × 1750905 = 35018100` ✓. |
| Unverified family and failure safety | ✅ Implemented | Unknown/inactive family writes `cedula_familiar_afiliado` as `FuenteCampoDeclarado`, `Confianza .5`, `RequiereConfirmacion: true`, then `Guardar`; save error wrapped and takes precedence over `ErrFamiliarNoEncontrado`; `PorID` and `ctx.Err()` return before mutation; family path publishes nothing. |
| Clean deterministic boundary | ✅ Implemented | Imports only `context`, `errors`, `fmt`, `time`, `internal/domain`, `internal/domain/motor`. No `time.Now()`, no DB/LLM/HTTP/ADK/messaging/kNN reference; time and IDs come from ports. Race suite clean. |

### Coherence (Design)
| Decision | Followed? | Notes |
|----------|-----------|-------|
| Boundary: ports + domain/motor only (NFR-M-01) | ✅ Yes | Import inspection + full in-memory fake execution. |
| Real 3-arg `CalcularCapacidad`, drop issue-snippet `ParametrosCredito{TasaEA}` | ✅ Yes | Documented snippet drift; doc 13 owns the money. |
| Profile vocabulary limited to `domain.CamposReconocidos` | ✅ Yes | 9 unexported constants, all recognized; no invented `no_encontrado` key. |
| Provenance: verified base + synthetic demo eligibility baseline | ✅ Yes | Affiliate hits only, as approved in exploration/proposal/spec. |
| Family idempotency + single-cedula guard | ✅ Yes | Matches the Contract's singular `cedula_familiar_afiliado`. |
| Commit/event order: compute → transition → `Crear` → publish | ✅ Yes | `EvPerfilCompleto` not misused; family saves emit no event. |
| `LeadNuevo` payload conforms to Contract v1.1 §6 | ✅ Yes | Prior WARNING-1 (`{"estado"}` payload) is **fixed**: payload is now exactly `{cedula, nombre, telefono, fuente}` and is asserted by exact map equality. Doc 10 outranks the issue snippet's `{nombre, afiliado, fuente}`. |
| Forecast ≤390 authored lines, single PR | ✅ Yes | Actual 356. |

### Issues Found

**CRITICAL**: None in implementation, tests, or spec compliance.

**Archive-gate condition (process, not code)** — persisting this report changes a receipt-bound path.
`verify-report.md` and `state.yaml` are inside the reviewed genesis path set of `review-8395941eaa5fec38`,
whose receipt binds the exact candidate tree `1cf0c727…`. A reversible probe (append one comment line,
validate, restore byte-identical, re-validate) measured:

- bound content modified → `review validate --gate post-apply` returns `result: invalidated`,
  `allowed: false`, `action: explicit-maintainer-action`, `denial.stage: receipt-discovery`,
  `code: receipt_ambiguous`;
- byte-identical restore → `result: allow` again.

So final verification evidence and the approved receipt cannot both be current for this lineage: the review
snapshot was taken with the SDD artifacts included, while final verification necessarily post-dates review.
Archive requires `reviewGate.result: allow` plus a receipt matching the final tree, so **the orchestrator must
re-bind authority to the new candidate tree (explicit maintainer action) before archive**. The pre-write bytes
are preserved at `/tmp/verify18-evidence/verify-report.orig.md` if reverting is preferred.

**WARNING**:
- **W1 — report is not validator-attested.** `gentle-ai sdd-verify-validate` is absent from both installed
  binaries (2.1.11, 1.43.3), so the mandated admission gate could not be executed; admission was self-checked
  against the format reference.
- **W2 — repo-wide gofmt drift** in `internal/pipeline/compradores_test.go` and
  `internal/pipeline/proyectos_test.go`: **pre-existing / base-only** (proof above). Does not block #18;
  a whole-repo formatting check in CI would fail on `main` independently of this change.
- **W3 — synthetic demo provenance.** `tiene_vivienda=false` / `recibio_subsidio=false` are stamped
  `VERIFICADO_BASE` with no source field behind them. Accepted, demo-scoped, and documented, but real affiliate
  data must supply explicit eligibility fields before this survives Fase 0.
- **W4 — catalog outage is indistinguishable from a genuine non-affiliate.** Intentional (proposal question 2
  answered "silent non-affiliate"), but there is no observability signal; US-01 case B's internal
  `no_encontrado` mark is deliberately absent because it is not a recognized Contract key.

**SUGGESTION**:
- Non-affiliate assertions use `len(lead.Perfil) == 0`, which a nil map also satisfies; the spec's "non-nil
  profile" clause holds in code but is not distinguished at runtime.
- The unknown-family test does not assert the negative half (`ingreso_hogar` absent, `hogar_con_afiliado` absent)
  nor `Confianza == 0.5`; the save-failure case does not assert storage is unchanged.
- `Ejecutar` returns a populated `SalidaPerfilar` alongside transition/create errors, so a caller can read a
  `LeadID` for a lead that was never persisted; the issue snippet returned the zero value.
- The `ana()` fixture uses `Categoria: "B"`, `Segmento: "A"` while `data/afiliados_mock.json` has `"A"`/`"BASICO"`,
  and the family fixture omits `Nombre`. Aligning the fake with the real fixture would raise demo fidelity at no cost.
- `tasks.md` 2.4 still implies transition-failure coverage; the branch is unreachable, so reword rather than test it.

### Verdict

**PASS WITH WARNINGS** (envelope `verdict: pass`, 0 blockers, 0 critical findings).

12/12 tasks complete, 6/6 requirements and 7/7 scenarios compliant with passing runtime evidence, 356/390
authored lines, and every check green: focused, full, race, build, vet, module verify, changed-file gofmt, and
diff/budget. Review authority `review-8395941eaa5fec38` is approved and was `allow`-bound to the verified tree.

**Archive readiness: NOT READY — one process step remains, no code change required.** Persisting this report
mutates a receipt-bound path and flips the gate to `invalidated` (`receipt_ambiguous`); authority must be
re-bound to the final candidate tree before `sdd-archive` may run.

### Post-write gate measurement (confirmed, not a prediction)

After persisting this report and the `state.yaml` verification fields:

- `./.tools/gentle-ai review validate --gate post-apply --cwd /tmp/vivi-issue-18` → `result: invalidated`,
  `allowed: false`, `action: explicit-maintainer-action`, `reason: "multiple terminal review receipts require
  explicit target selection"`, `denial: {stage: receipt-discovery, code: receipt_ambiguous}`.
- `./.tools/gentle-ai sdd-status issue-18-perfilar-lead --json` → `nextRecommended: resolve-review`,
  `blockedReasons: ["bound compact post-apply gate context changed"]`, `dependencies.verify: blocked`,
  `dependencies.archive: blocked`.

Implementation and tests were not modified by this phase: `internal/usecase/perfilar_lead.go`
`sha256:47fb24de91080c63e261b7810bcac5e8ad2085a577d27d4851550ae62098e500`,
`internal/usecase/perfilar_lead_test.go`
`sha256:68888ddeadb2e078bbaa2ad59294f2d4cfa20799077260092a1d81052f186bc0`;
`git status --porcelain` still lists only the two untracked source files and the untracked change folder.
Substantive verification stands as PASS WITH WARNINGS; only review-authority re-binding remains before archive.
