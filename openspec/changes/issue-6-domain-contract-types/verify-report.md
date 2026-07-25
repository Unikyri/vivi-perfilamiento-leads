```yaml
schema: gentle-ai.verify-result/v1
evidence_revision: sha256:882d851d40675b949cb4159e0ef8ce8097c5cf70f9d866e7523b914f88522d8a
verdict: pass
blockers: 0
critical_findings: 0
requirements: 6/6
scenarios: 7/7
test_command: go test ./internal/domain/... ./internal/pipeline/... -v -count=1
test_exit_code: 0
test_output_hash: sha256:9a7bc59110faceee779bce5a7b562ba300e365b605a7fef85ec9a29209b5f595
build_command: go build ./...
build_exit_code: 0
build_output_hash: sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
```

## Verification Report

**Change**: issue-6-domain-contract-types (GitHub Issue #6 — `[A1] Domain entities, enums, and contract types`)
**Version**: `openspec/changes/issue-6-domain-contract-types/specs/domain-contract-types/spec.md` (6 requirements, 7 scenarios)
**Mode**: Standard (`strict_tdd: false`)
**Target**: merged `main`, HEAD `73cf631` (`feat(domain): add advisor ficha types`); baseline `3e7d542`
**Toolchain**: `go1.25.0 linux/amd64`; `gentle-ai 1.43.3`
**Artifact store**: hybrid — proposal, spec, design, tasks, apply-progress, and state read from OpenSpec files and Engram (`#452`–`#458`)

### Completeness

| Metric | Value |
|--------|-------|
| Tasks total | 9 |
| Tasks complete | 9 |
| Tasks incomplete | 0 |

Native `gentle-ai sdd-status issue-6-domain-contract-types --json` confirms `taskProgress 9/9`, `applyState: all_done`, `dependencies.verify: ready`, `blockedReasons: []`, and `verifyReport: missing` before this run.

### Build & Tests Execution

**Build**: ✅ Passed

```text
$ go build ./internal/domain/...        # exit 0, no output
$ go build ./...                        # exit 0, no output
$ go vet ./internal/domain/...          # exit 0, no output
$ go vet ./...                          # exit 0, no output
```

**Tests**: ✅ 19 top-level tests / 58 test nodes passed, 0 failed, 0 skipped

```text
$ go test ./internal/domain/... ./internal/pipeline/... -v -count=1
ok  github.com/Unikyri/vivi-perfilamiento-leads/internal/domain     0.007s
?   github.com/Unikyri/vivi-perfilamiento-leads/internal/domain/motor  [no test files]
ok  github.com/Unikyri/vivi-perfilamiento-leads/internal/pipeline   0.005s
```

Issue #6 contract tests: `TestEnumJSONWireValues`, `TestPerfilAccessors` (8 subtests), `TestPerfilVerificationAndKeySets`, `TestPerfilSchemaAndJSON`, `TestCapacityContractTypes`, `TestLeadAndPlanContractTypes`, `TestFichaContractTypes`. Pre-existing `constantes` and `pipeline` suites also pass unchanged.

**Dependency guard**: ✅ No forbidden layer

```text
$ go list -deps ./internal/domain/... | grep -E 'internal/(usecase|adapters|infrastructure)'
(no output, 0 bytes)
$ go list -f '{{.ImportPath}}: {{join .Imports " "}}' ./internal/domain/...
internal/domain: encoding/json fmt time
internal/domain/motor:
```

**Verifier schema probe** (independent, `go test -overlay`, no repository file created or modified): 3 additional tests passed — exact ordered field name/Go type/JSON tag for all 14 contract structs, all 45 enum literals marshalled to their exact wire strings, and nullability (`null` for `*string`, `*int64`, `*time.Time`, nil maps/slices) plus omitted `Lead.Version` and `Mensaje.LeadID`.

**Coverage**: 100.0% of statements (`internal/domain`) / threshold: 0% → ✅ Above

### Spec Compliance Matrix

| Requirement | Scenario | Test | Result |
|-------------|----------|------|--------|
| Contract Enums | Enum wire values | `internal/domain/perfil_test.go > TestEnumJSONWireValues` (+ verifier probe `TestProbeAllEnumLiterals`) | ⚠️ PARTIAL |
| Perfil Schema and Accessors | Accessor conversion | `perfil_test.go > TestPerfilAccessors` | ✅ COMPLIANT |
| Perfil Schema and Accessors | Verification and key sets | `perfil_test.go > TestPerfilVerificationAndKeySets`, `TestPerfilSchemaAndJSON` | ✅ COMPLIANT |
| Capacity and Recommendation Types | Capacity schema | `perfil_test.go > TestCapacityContractTypes` (+ verifier probe `TestProbeExactSchemas`) | ⚠️ PARTIAL |
| Lead and Message Types | Lead/message schema | `perfil_test.go > TestLeadAndPlanContractTypes` (+ verifier probe `TestProbeExactSchemas`, `TestProbeNullabilityMarshal`) | ⚠️ PARTIAL |
| Nutrition and Ficha Types | Nutrition/ficha schema | `perfil_test.go > TestLeadAndPlanContractTypes`, `TestFichaContractTypes` (+ verifier probe `TestProbeNullabilityMarshal`) | ⚠️ PARTIAL |
| Compatibility, Isolation, and Evidence | Relocation and isolation | `go build`/`go vet`/`go test ./internal/pipeline/...`/dependency guard | ✅ COMPLIANT |

**Compliance summary**: 7/7 scenario outcomes proven by passing runtime evidence; 3/7 are fully proven by committed repository tests alone, and 4/7 depend on the verifier probe for the exact-type, exact-JSON-name, and nullability clauses (see WARNING-1). No scenario is FAILING or UNTESTED.

### Correctness (Static Evidence)

| Requirement | Status | Notes |
|------------|--------|-------|
| Contract Enums | ✅ Implemented | `enums.go`: 11 typed `string` groups, 45 literals, exact spec wire values, no extra groups. |
| Perfil Schema and Accessors | ✅ Implemented | `perfil.go`: 5-field `CampoPerfil`, `Perfil map[string]CampoPerfil`, exactly 4 exported methods, `CamposReconocidos` 18 keys, `CamposCriticos` 4 keys. |
| Capacity and Recommendation Types | ✅ Implemented | `capacidad.go`: `ItemDesglose`, `Capacidad`, `Intencion`, `Recomendacion`, plus relocated `Proyecto`/`Comprador` field-identical to spec. |
| Lead and Message Types | ✅ Implemented | `lead.go`: `Lead` 16 fields with `Version:"-"`, `Mensaje` 7 fields with `LeadID:"-"`. |
| Nutrition and Ficha Types | ✅ Implemented | `plan.go` (`Hito`, `PlanNutricion`), `ficha.go` (`AlertaDesistimiento`, `Identificacion`, `Ficha` 14 fields) with exact pointers. |
| Compatibility, Isolation, and Evidence | ✅ Implemented | `comprador.go` and `proyecto.go` deleted; pipeline consumers (`compradores.go`, `proyectos.go`, `distribuciones.go`, tests) compile unchanged; domain imports only `encoding/json`, `fmt`, `time`. |

### Coherence (Design)

| Decision | Followed? | Notes |
|----------|-----------|-------|
| Copy the 11 specified enum groups and exact structs/tags | ✅ Yes | Verified literal-by-literal and field-by-field. |
| `Perfil` is `map[string]CampoPerfil`; `ActualizadoEn` is `time.Time` | ✅ Yes | Confirmed by reflection at runtime. |
| Expose only `Entero`, `Booleano`, `Texto`, `EsVerificado` | ✅ Yes | `reflect.TypeOf(Perfil{}).NumMethod() == 4` asserted in the committed test. |
| Delete legacy files and relocate identical structs | ✅ Yes | `git diff --name-status`: `D comprador.go`, `D proyecto.go`, `A capacidad.go`; no consumer edits. |
| Keep domain production code stdlib-only (NFR-M-01) | ✅ Yes | Import list and dependency guard both clean. |
| Six production files + one test file created | ✅ Yes | `enums.go`, `perfil.go`, `capacidad.go`, `lead.go`, `plan.go`, `ficha.go`, `perfil_test.go`. |
| No pipeline or `Docs/` changes | ✅ Yes | `git diff --name-status 3e7d542..HEAD` touches only `internal/domain/*` and `openspec/changes/issue-6-domain-contract-types/*`. |

### Delivery Guard

| Field | Value |
|---|---|
| Authored changed lines, code | 587 (551 additions + 36 deletions in `internal/`) |
| Authored changed lines, SDD artifacts | 448 additions |
| Total for the change | 1,035 (forecast was 961–1,078) |
| Review budget | 400 lines per PR |
| Delivery | `stacked-to-main`, 7 sequential PRs, each forecast ≤ 210 lines |
| Verdict | ✅ Budget respected per PR; no `size:exception` needed |

### Admission Procedure

The required admission command `gentle-ai sdd-verify-validate --input <path> --requirements 6 --scenarios 7` is **not available** in the installed native binary (`gentle-ai 1.43.3`, `/home/daikyri/.local/bin/gentle-ai`); the CLI answers `unknown command "sdd-verify-validate"`. `gentle-ai sdd-status issue-6-domain-contract-types --cwd . --json` was executed instead and reported `dependencies.verify: ready`, `blockedReasons: []`, and `verifyReport: missing`.

No prior `verify-report` existed in either backend (OpenSpec `artifactPaths.verifyReport: []`; Engram search returned 7 artifacts, none a verify-report), so persistence cannot overwrite or truncate earlier evidence. Admission was therefore self-administered against the strict envelope rules before writing: the envelope is the first non-empty content, every field appears exactly once, requirement and scenario totals were counted from the retrieved spec (`grep -c '^### Requirement:'` → 6, `grep -c '^#### Scenario:'` → 7), and both declared commands were executed with recorded exit codes and SHA-256 output hashes. `evidence_revision` is the SHA-256 of the evidence manifest (HEAD, both declared commands, exit codes, output hashes, vet and dependency-guard results, requirement and scenario counts). This substitution is recorded as WARNING-4 so the orchestrator can re-run native admission after upgrading the binary.

Envelope counts and matrix statuses are reconciled as follows: the envelope numerator counts scenarios whose specified outcome is proven by passing runtime evidence (7), while `⚠️ PARTIAL` in the matrix flags committed-test depth that is thinner than the scenario text and is carried as WARNING-1 rather than a blocker.

### Issues Found

**CRITICAL**: None

**WARNING**:
1. Committed contract tests assert struct field names, field counts, and the presence of a `json` tag, but not the exact Go type, the exact JSON name, or nullability. `TestEnumJSONWireValues` marshals 11 of the 45 enum literals; `TestFichaContractTypes` and `TestLeadAndPlanContractTypes` never marshal `Ficha`, `PlanNutricion`, or `Hito`. The implementation is exact — proven by the verifier overlay probe — but the repository would not detect a future type, tag, or nullability regression. Recommended follow-up: fold the probe assertions into `perfil_test.go`.
2. Engram hybrid drift: `sdd/issue-6-domain-contract-types/apply-progress` (`#458`) is still the PR-4 snapshot listing tasks 1.2, 2.3, 2.4, and 3.1–3.3 as remaining, while `openspec/.../apply-progress.md` records all 9 tasks complete. The upsert was not re-run for PRs 5–7.
3. Engram hybrid drift: `sdd/issue-6-domain-contract-types/state` (`#452`) still reports `Current slice: PR 1/7` and `apply is ready`, while `openspec/.../state.yaml` reports `apply: success`, `verify: pending`. The orchestrator should refresh both topics before archive.
4. Native verify admission validator unavailable in `gentle-ai 1.43.3` (see Admission Procedure). Self-administered strict admission was used with zero risk of destroying prior evidence.

**SUGGESTION**:
1. Pre-existing and out of Issue #6 scope: `go.mod` lists `github.com/xuri/excelize/v2 v2.11.0` under the indirect block although `internal/pipeline` imports it directly; any `-mod=mod` build promotes it to a direct requirement and dirties `go.mod`. Default `-mod=readonly` builds are unaffected. Fix in a pipeline-scoped change, not here.
2. `TipoMensaje` is exported by the Contract but not referenced by any struct in this change; `Mensaje` carries `TipoContenido` only. This matches the spec exactly and is expected to be consumed by the conversation layer.
3. `Semaforo` is likewise defined for the dashboard and not yet consumed.

### Verdict

**PASS WITH WARNINGS**

All 9 tasks are complete, all 6 requirements are implemented exactly as specified, and all 7 scenario outcomes are proven by passing runtime evidence: `go test` (domain + pipeline), `go build ./...`, `go vet ./...`, and the dependency guard all exit 0 with no forbidden layer, domain statement coverage is 100%, the legacy `Comprador`/`Proyecto` relocation left every pipeline consumer compiling unchanged, and no pipeline or `Docs/` file was touched. The only deviations are committed-test depth below the scenario text for four schema scenarios, two stale Engram topics, an unavailable native admission validator, and one pre-existing out-of-scope `go.mod` classification.
