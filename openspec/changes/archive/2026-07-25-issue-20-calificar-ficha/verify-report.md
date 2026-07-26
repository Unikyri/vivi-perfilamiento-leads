```yaml
schema: gentle-ai.verify-result/v1
evidence_revision: sha256:defc2383caf6f10de695bf896254b23201d457882190ceec0cee07483aa615ea
verdict: pass
blockers: 0
critical_findings: 0
requirements: 6/6
scenarios: 7/7
tasks: 10/10
test_command: go test ./... -count=1
test_exit_code: 0
test_output_hash: sha256:588ee9ea9dd89b695355a28e63cf5545d97c017eb462c3531c1639c8697aefac
build_command: go build ./...
build_exit_code: 0
build_output_hash: sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
authority_only_failure: false
missing_review_authority: false
substantive_failure: false
command_failed: false
observed_authority_revision: sha256:af0ef565cb7571ae04d4b06d9ea4be86e0c795c77389acb16bb19dc88a7b84ab
review_lineage: review-90c09c25bf720d01
review_gate_post_apply: allow
review_authority_scope: code-and-artifact-tree-at-commit-23f4a2e
validator_available: false
admission: self_checked_validator_unavailable
validator_note: "gentle-ai sdd-verify-validate is absent from both installed binaries (.tools/gentle-ai 2.1.11 and PATH gentle-ai 1.43.3 answer 'unknown command'). This report was NOT validator-attested; it is self-checked field-by-field against sdd-verify/references/report-format.md, exactly as the issue-18 and issue-19 reports were."
tasks_artifact_recorded: 9/10
tasks_checkbox_stale: true
tasks_evidence: commit-23f4a2e-closes-3.2
archive_ready: false
archive_precondition: "openspec tasks.md 3.2 checkbox and state.yaml task_progress must be closed in a committed docs work unit with its own review receipt; editing those reviewed tracked paths in place forfeits the current allow gate (proved below)."
```

# Verification Report: Issue #20 CalificarLead + GenerarFicha (final, pre-archive)

**Change**: `issue-20-calificar-ficha`
**Issue**: #20 `[A12] Casos de uso: CalificarLead + GenerarFicha (UC-06, US-08/09/15)` — labels `bloque-a`, `contrato`, `status:approved`; issue state still OPEN
**Worktree / branch**: `/tmp/vivi-issue-20-verify` on `test/issue-20-qualification-edges`, HEAD `23f4a2e` (tree `7f45571`), one commit ahead of `origin/main` `38acf2d` (tree `b3d5aa9`)
**Merged functional slices**: `4b08a81` (qualification, PR #81) and `21751e1` (ficha, PR #82) are both already in `origin/main`
**Artifact store**: hybrid — OpenSpec files + Engram observations #570–#577
**Mode**: Standard (`openspec/config.yaml` `strict_tdd: false`; `strict-tdd-verify.md` deliberately not loaded)
**Run**: the single independent requirements/runtime final verification, read-only on all production and test code

## Final Result

**PASS WITH WARNINGS — 6/6 requirements, 7/7 scenarios, 10/10 tasks; every declared command exits 0; review gate `post-apply` returns `allow`.**

No behavioral defect, no failing test, and no missing scenario coverage remains. The three conditions that
were `PARTIAL` in the previous provisional report are now proved by passing discriminating tests, and the
previously reported authority blocker is resolved: lineage `review-90c09c25bf720d01` is `approved` and binds
exactly this candidate tree. What is left is bookkeeping, not verification: the OpenSpec `tasks.md` checkbox
for 3.2 and `state.yaml` `task_progress` still read `9/10`, and closing them in place would forfeit the
receipt binding (proof below). That is an orchestrator/archive decision, recorded here as
`archive_ready: false`.

## Validator admission (transparent)

`gentle-ai sdd-verify-validate` does not exist in either installed binary:

| Probe | Result |
|---|---|
| `./.tools/gentle-ai --version` | `gentle-ai 2.1.11` |
| `./.tools/gentle-ai sdd-verify-validate --help` | `Error: unknown command "sdd-verify-validate"` |
| `gentle-ai --version` (PATH `/home/daikyri/.local/bin/gentle-ai`) | `gentle-ai 1.43.3` |
| `gentle-ai sdd-verify-validate --help` | `Error: unknown command "sdd-verify-validate"` |

Consequence, stated plainly: this report carries **no validator attestation**. Every envelope field was
self-checked against `sdd-verify/references/report-format.md` — envelope first, each field once, counts taken
from the two retrieved delta specs, commands and exit codes actually executed. The prior report was replaced
only because the validator is unavailable and the orchestrator explicitly authorized a self-checked fallback.

## Review authority (code authority only)

```
./.tools/gentle-ai review validate --gate post-apply --cwd /tmp/vivi-issue-20-verify
result: allow | allowed: true | action: continue
reason: "authoritative transaction, current repository target, and content-bound artifacts match"
lineage_id: review-90c09c25bf720d01   generation: 1
store_revision / genesis_revision / chain_identity: sha256:af0ef565cb7571ae04d4b06d9ea4be86e0c795c77389acb16bb19dc88a7b84ab
base_tree: b3d5aa907e21a2ff06867ebaff1ac363e0e37746   candidate_tree: 7f45571aed843a4104045e7d4aad803c123587c1
paths_digest: sha256:700622b23ce994dc19915a2ce3244c2019d39559cda140c0642bcc92412828cf
policy_hash: sha256:34fb63d7f29f8613cd4431382b1057398a4816f8a4c20fc34677fffc80a184f6
evidence_hash: sha256:23db77040349b93f89c3170119d1f8947f739efd17f078c4c834560243b3c1e9
gate JSON bytes: sha256:9276d4b3cea0122f4e715494303bed6e24279e0b15dae90f68e2d3c2531a8063
```

`review-90c09c25bf720d01` (risk `medium`, lens `review-reliability`, zero findings, original changed lines
`126`, correction budget `63`) reviewed exactly four paths: `internal/usecase/calificar_lead_edges_test.go`,
`apply-progress.md`, `state.yaml`, `tasks.md`. Its `final_candidate_tree` equals the current HEAD tree.

**Scope limit — do not over-read this receipt.** It attests the code and artifact bytes committed at
`23f4a2e`. It does **not** attest this verify report, and it cannot attest any archive-time document that
does not exist yet. Two facts were measured rather than assumed:

| Repository state | Gate result |
|---|---|
| Clean tree, untracked `verify-report.md` present | `allow` (this report's presence is gate-neutral: untracked and outside `paths`) |
| Any edit to a reviewed tracked artifact (`state.yaml` probe) | `invalidated` / `receipt-discovery` / `receipt_ambiguous` |
| Same edit plus `--lineage review-90c09c25bf720d01` | `scope-changed` / `receipt-binding` / `candidate-or-paths-mismatch` |

The probe edit was appended and then restored byte-identically
(`sha256 f6cf7ee6f0eb74e53eb6d5780628583a93c299b555fe9be73aee2a21635c5437` before and after,
`git status` clean apart from this untracked report), and `allow` was re-confirmed afterwards. This is the
precise reason the earlier report saw `receipt_ambiguous`: with a dirty reviewed path, discovery cannot
resolve one receipt among the 37 approved lineages in the store.

## Executed evidence

Hashing method: combined stdout+stderr captured by command substitution (trailing newline dropped), then
`printf '%s' "$output" | sha256sum`. `sha256:e3b0c442…b855` is the digest of empty output.

| Command | Exit | Output SHA-256 |
|---|---:|---|
| `go test ./internal/usecase -run 'TestCalificarLead(PriorityCapsRatioAndUsesPartialConfidence\|NonCanonicalCatalogKeyOmitsKNNZone\|RutaDecididaPayloadContainsCopiedDecisionData)$' -count=1 -v` | 0 | `94b3bd8a6ad6df4f224b917c703dc0aee43ad1586a48c6034e3792f34fee7fba` |
| `go test ./internal/usecase -run TestCalificarLead -count=1 -v` | 0 | `fb9d22bcd1f2ce29d74b2f86e17065437516c3f16673f4874f7d77fbda540931` |
| `go test ./internal/usecase -run TestGenerarFicha -count=1 -v` | 0 | `81e13338c6a8c771a54126320c66741a12452141d5a4989308349b11232d6feb` |
| `go test ./... -count=1` | 0 | `588ee9ea9dd89b695355a28e63cf5545d97c017eb462c3531c1639c8697aefac` |
| `go test -race ./... -count=1` | 0 | `9a70c90441f1b8af47f43e945194cc576687b33a13d81bc41318e4f6c996ba1d` |
| `go build ./...` | 0 | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` (empty) |
| `go vet ./...` | 0 | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` (empty) |
| `go mod verify` | 0 | `fd3282825b01f3034bf26af0f73af983a6fef5279109c45d0444bdf1289e8996` (`all modules verified`) |
| `go mod tidy -diff` | 0 | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` (empty) |
| `go list ./...` | 0 | `55b45095557afac01ee6a0bb1a1e442d79205dc91002b19d031b9c3867fd9804` |
| `gofmt -l internal/usecase/` | 0 | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` (no output) |
| `git diff --check` | 0 | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` (no output) |

Toolchain `go1.25.0 linux/amd64`. Package results: `internal/usecase` 27 top-level tests plus 42 subtests,
0 failures, 0 skips; whole-module `go test ./...` green with no `[no test files]` package regression.
Coverage `internal/usecase` = **86.1% of statements** (`coverage_threshold: 0`, so informational).

Runtime harness: **N/A and justified.** Both use cases are synchronous, provider-free application services
with no HTTP, shell, subprocess, or external boundary. The durability boundary is exercised in-process by
fake lead/ficha/catalog repositories, a fixed `Reloj`/`GeneradorID`, ordering probes, and a fail-at-step CAS
harness (`TestGenerarFichaFichaFirstFailureAndRetry`, `TestCalificarLeadRoutesConversionPriorityAndDurability`).

## The three formerly PARTIAL conditions — now proved

All three run in `internal/usecase/calificar_lead_edges_test.go` (63 lines, test-only, committed at `23f4a2e`)
and all three passed by name in this run:

| Former gap | Test | What makes it discriminating | Result |
|---|---|---|---|
| Ratio cap only touched at exactly `1.2`, confidence always `1.0` | `TestCalificarLeadPriorityCapsRatioAndUsesPartialConfidence` | asserts `Ratio > 1.2` **and** `Confianza < 1` (income source forced to `FuenteCampoDeclarado`), then `Prioridad == 1.2 * Confianza && Prioridad < 1.2`; dropping either the cap or the `* Confianza` factor now fails | ✅ PASS |
| Non-canonical catalog zone never exercised | `TestCalificarLeadNonCanonicalCatalogKeyOmitsKNNZone` | compares `construirDecision` with map key `"display-name"` vs `"p1"` for the same project; expects the Gower distance renormalized over `0.80` participating weight (`0.15*(7.5/32.5)/0.8`) versus `0.15*(7.5/32.5)` with the zone present, and `withoutZone > withZone`. The constants match `motor/knn.go` (`gowerZoneWeight 0.20`, `gowerAgeWeight 0.15`, `gowerAgeRange 32.5`), so the assertion is derived from motor criteria, not tautological | ✅ PASS |
| `RutaDecidida` payload `prioridad`/`recomendaciones` unasserted, defensive copy unproved | `TestCalificarLeadRutaDecididaPayloadContainsCopiedDecisionData` | asserts `ruta`, `semaforo`, `consume_cupo_10`, `prioridad`, and `reflect.DeepEqual` on `recomendaciones`, then mutates `recommendations[0].Nombre` and requires the use-case output to be unchanged — proving `copiarRecomendaciones` really breaks aliasing between payload and output | ✅ PASS |

```text
--- PASS: TestCalificarLeadPriorityCapsRatioAndUsesPartialConfidence (0.00s)
--- PASS: TestCalificarLeadNonCanonicalCatalogKeyOmitsKNNZone (0.00s)
--- PASS: TestCalificarLeadRutaDecididaPayloadContainsCopiedDecisionData (0.00s)
ok  github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase  0.003s
```

## Completeness

| Metric | Value |
|---|---|
| Tasks total | 10 |
| Tasks complete (substantive evidence) | 10 |
| Tasks incomplete (substantive evidence) | 0 |
| Tasks as recorded in `tasks.md` | 9/10 — `3.2` still `[ ]` with the obsolete "instructed not to commit" rationale |
| Measured requirements / scenarios | 6 / 7 (calificar-lead 3+4, generar-ficha 3+3) |
| Slice 1 authored source | `calificar_lead.go` 238 + `calificar_lead_test.go` 162 = 400 (exactly at the hard stop) |
| Slice 2 authored source | `generar_ficha.go` 199 + `generar_ficha_test.go` 189 = 388 |
| Edge-hardening slice | `calificar_lead_edges_test.go` 63 authored test lines, 0 runtime lines |
| Candidate diff vs `origin/main` | `63/0` test lines + `21/21`, `15/4`, `1/1` doc lines = 100 insertions, 26 deletions |

Task 3.2 has two halves. The validation half is re-proved above. The commit/work-unit closure half is now
satisfied by `23f4a2e`, a single test-only commit whose reviewed scope is exactly the edge suite plus this
change's artifacts — hence `tasks: 10/10` on substance. The stale checkbox is a documentation defect, listed
under WARNING, not a verification failure.

## Spec compliance matrix

| Spec | Requirement | Scenario | Runtime evidence | Result |
|---|---|---|---|---|
| calificar-lead | Guard, Capacity, and Recommendations | Guard and candidate | `TestCalificarLeadGuardsAndReadFailuresDoNotWrite` (7 cases: blank id, wrong state, missing intention, missing capacity, projects failure, buyers failure, cancellation → zero saves, zero events) + `TestCalificarLeadCandidateKNNAndDeterminism` (candidate `80` chosen from `{80,90}`; unaffordable catalog → candidate `0` and `MedianaVISCOP` fallback ratio equality) | ✅ COMPLIANT |
| calificar-lead | Guard, Capacity, and Recommendations | Exact recommendation evidence | 31 buyers → `len(Vecinos) == 30`; repeated `construirDecision` gives `reflect.DeepEqual` recommendations; `TestCalificarLeadNonCanonicalCatalogKeyOmitsKNNZone` proves a non-exact catalog key is not substituted as a zone | ✅ COMPLIANT |
| calificar-lead | Conversion, Route, Priority, Cupo | Predicates and matrix | 8-case table over all four quadrants and all three conversion signals, blank `caja_externa` rejected, `Matriz2x2` delegation, semáforo, cupo only for non-affiliate `ASESOR`; ratio strictly above `1.2` with confidence `< 1` in `TestCalificarLeadPriorityCapsRatioAndUsesPartialConfidence` | ✅ COMPLIANT |
| calificar-lead | Durable Route Event | Save ordering | `order == ["save","event"]` for all four routes, exactly one event; injected CAS failure → `["save"]`, zero events; payload proved complete and non-aliasing | ✅ COMPLIANT |
| generar-ficha | Eligibility and Recommendation Reconstruction | Guard and parity | `TestGenerarFichaGuardsAndReadFailuresDoNotWrite` (10 cases incl. wrong route, ficha read failure, ficha save failure, cancellation → no ficha, no lead write); recommendation parity vs `construirDecision`; no provider field and no LLM reference in `generar_ficha.go` | ✅ COMPLIANT |
| generar-ficha | Ordered Contract Content | Stable content and threshold | fixed clock/ID run asserts exact `Beneficios` order, rent argument `Paga $1.200.000 de arriendo; la cuota estimada es $1.200.000`, warning exactly `PERFIL PARCIALMENTE DECLARADO — validar campos marcados` at confidence `< 0.6`, `1/5 = .20` → `Activa false`, `2/5 = .40` → `Activa true`, `Detalle nil`, non-aliasing of `Perfil`/`Intencion`, unmodified `domain.Ficha` field order | ✅ COMPLIANT |
| generar-ficha | Ficha-First Handoff and Retry | Partial-write repair | `TestGenerarFichaFichaFirstFailureAndRetry`: ficha stored then CAS failure leaves lead `CALIFICADO`/`ASESOR`; retry reuses `ficha-1` and `GeneradaEn`, `fichas.calls` = read/save/read/save, lead order `["save","save"]`, `ENTREGADO` only after save, `ActualizadoEn` from `Reloj` | ✅ COMPLIANT |

**Compliance summary: 6/6 requirements, 7/7 scenarios COMPLIANT. Zero `FAILING`, zero `UNTESTED`, zero `PARTIAL`.**

## Correctness against Issue #20 and the governing criteria

| Governing rule (Contract §3.5 / doc 13 / RF-M6-01 as quoted in Issue #20) | Status | Evidence |
|---|---|---|
| `prioridad = pesoRuta × min(ratio,1.2) × confianza`, weights `1/.5/.25/.1` | ✅ Implemented and now fully tested | `construirDecision` clamps ratio to `[0,1.2]`; route table plus the new cap/confidence test |
| Ratio denominator is the lowest affordable project, median fallback at `0` | ✅ | `menorPrecioAsequible` picks `80` from `{80,90}`; empty selection → `0` → motor median |
| Semáforo `VERDE/AMBAR/GRIS` (RF-M6-02) | ✅ | `semaforoRuta` + table |
| Conversion `!afiliado && (INDEPENDIENTE ‖ hogar_con_afiliado ‖ caja_externa != "")` | ✅ | `esConversion`, blank/whitespace caja rejected |
| `ConsumeCupo10` only for non-affiliate `ASESOR` | ✅ | affiliate `ASESOR` → false; non-affiliate `ASESOR` → true |
| State mapping `ASESOR` stays `CALIFICADO`; others → `EN_NUTRICION`/`REMARKETING`/`DESPEDIDO` | ✅ | `aplicarEstadoRuta` + table |
| Ficha sections in fixed RF-M6-01 order, numbers only from the motor | ✅ | unmodified `domain.Ficha`; all figures from `CalcularCapacidad`/`GemeloKNN`/`RecomendarProyectos` |
| Warning band below confidence `0.6`, exact string | ✅ | asserted literal |
| Withdrawal alert strictly above `0.20` | ✅ | `.20` inactive, `.40` active |
| Rent-vs-installment argument when rent declared | ✅ | dotted-COP formatting, `cuota = ingreso*40/100` |
| `GemeloKNN(K=30)` with exact catalog-ID zones and `personas_hogar-1` dependents | ✅ | `construirDecision` + zone edge test |
| Public API `Ejecutar(ctx, leadID) (*domain.Ficha, error)` with `Leads/Fichas/Catalogo/IDs/Reloj` | ✅ | matches Issue #20 sketch; legacy entry/exit types removed |
| Accepted deviations from the Issue code sketch | ⚠️ Documented | `BandaDeCategoria` not reimplemented (the motor owns income-category projection, pre-approved in `design.md`); `menorPrecioAsequible` additionally requires `PrecioDesde > 0`, which the delta spec explicitly demands |

## Design coherence

| Design decision | Followed? | Notes |
|---|---|---|
| One package-private `construirDecision` shared by both use cases | ✅ Yes | capacity(0) → candidate → capacity(candidate) → buyers → exact zones → `GemeloKNN(K=30)` → `RecomendarProyectos` |
| Zones only from exact `map key == ProyectoID` | ✅ Yes | now runtime-proved, no longer inspection-only |
| Dependents = `personas_hogar - 1` when present | ✅ Yes | |
| Lead CAS before `RutaDecidida`; ficha upsert before `CALIFICADO→ENTREGADO` | ✅ Yes | ordering probes on both use cases |
| Retry reuses existing ficha ID/time; no duplicate event | ✅ Yes | retry test |
| No new port, transaction, domain, motor, adapter, or boundary change | ✅ Yes | candidate diff touches only `calificar_lead_edges_test.go` plus this change's artifacts |
| Best-effort event delivery documented, not exactly-once | ✅ Yes | design states it; no outbox claimed |

## Scope

Candidate against `origin/main`: `internal/usecase/calificar_lead_edges_test.go` (new, test-only) plus this
change's `apply-progress.md`, `state.yaml`, and `tasks.md`. Forbidden paths all unchanged: `internal/domain/**`,
`internal/domain/motor/**`, `internal/usecase/puertos.go`, adapters, infrastructure, HTTP, frontend,
migrations, config, Contract, and Wiki. Every slice stayed within the 400 authored-line budget
(400 / 388 / 63); no `size:exception` was used.

## Issues Found

### CRITICAL
None. No failing test, no uncovered required scenario, no forbidden-path change, no missing review authority.

### WARNING
1. **`tasks.md` 3.2 and `state.yaml` `task_progress` are stale.** They still read `9/10` with the obsolete
   "instructed not to commit" rationale even though `23f4a2e` closed that work unit. Native
   `gentle-ai sdd-status issue-20-calificar-ficha --json` therefore reports `taskProgress 9/10`,
   `dependencies.verify: blocked`, `dependencies.archive: blocked`, `nextRecommended: apply`. Archive routing
   cannot succeed until this is closed.
2. **Closing warning 1 in place forfeits the receipt.** Measured, not assumed: any edit to `state.yaml` or
   `tasks.md` flips `review validate --gate post-apply` from `allow` to `invalidated` (`receipt_ambiguous`),
   and to `scope-changed` (`candidate-or-paths-mismatch`) even with `--lineage` override, because the receipt
   binds `candidate_tree 7f45571` and those exact paths. This verifier therefore did **not** edit those two
   tracked files; that is an orchestrator decision with a cost either way.
3. **Deliberate hybrid divergence for this run.** Engram `sdd/issue-20-calificar-ficha/{verify-report,state,tasks}`
   were updated with the 6/6, 7/7, 10/10 result; the OpenSpec `state.yaml` and `tasks.md` were left byte-identical
   to the reviewed tree for the reason in warning 2. Only the untracked `verify-report.md` was overwritten
   (proved gate-neutral). Hybrid parity for `state.yaml`/`tasks.md` must be restored by the docs work unit below.
4. **Engram content parity gaps (pre-existing).** `sdd/issue-20-calificar-ficha/spec` (#573) holds a short
   correction note rather than the full delta-spec text, and `sdd/.../state` (#570) previously recorded
   `393/204` authored/runtime lines against the physical `388/199`. Engram alone cannot rebuild the specs.
5. **Delivery not finished at the repository level.** Issue #20 is still OPEN and
   `test/issue-20-qualification-edges` is one commit ahead of `origin/main` with no pushed branch or PR, so
   the edge suite is not yet on `main`. Issue DoD item "PR abierto a main" is satisfied only for the two
   functional slices (PR #81, PR #82).
6. **Untested minor branch.** `argumentosFicha` returns no argument for both a missing `arriendo_actual` key
   and a non-positive value; only the `0` case is exercised, so the `!ok` half of `!ok || renta <= 0` has no
   discriminating test. Behaviorally correct by inspection, and outside every scenario clause.

### SUGGESTION
- `cuota` uses integer `ingreso*40/100` while the Issue sketch uses `int64(0.40 * float64(ingreso))`; results
  agree on tested inputs but can differ by one peso. Worth pinning with a table case.
- `estadoTrasRuta` duplicates the mapping in `aplicarEstadoRuta`; one source would remove drift risk.
- Slice 1 sits exactly at 400 authored lines, so any future slice-1 edit needs its own slice.

## Verdict

**PASS WITH WARNINGS.** Requirements 6/6, scenarios 7/7, tasks 10/10 on substantive evidence, all twelve
declared checks exit 0, coverage 86.1% in the touched package, and code authority `review-90c09c25bf720d01`
validates `allow` for `post-apply`. Not archive-ready yet: the OpenSpec task/state bookkeeping must be closed
in a way that also re-establishes a bound receipt.

## Next

Recommended, in order:

1. Commit one small docs-closure work unit on this branch: `tasks.md` 3.2 → `[x]` with the real closure
   rationale, `state.yaml` `task_progress 10/10` plus this verify block, and this `verify-report.md`.
2. Run `review start` → `review finalize` for that docs-only slice (low risk, no production code), then
   confirm `review validate --gate post-apply` returns `allow` against the new tree.
3. Push the branch, open the PR to `main`, and let CI confirm the same evidence.
4. Archive only after structured status shows `taskProgress 10/10` and `reviewGate.result: allow` for the
   receipt bound to the final tree. Do not treat `review-90c09c25bf720d01` as covering any archive-time
   document: it attests the tree at `23f4a2e` and nothing created afterwards.
