# Apply Progress: Slice 3 — Catalog Cache and Integration Harness

**Change**: `issue-15-postgres-repositories`
**Mode**: Standard (strict TDD disabled by `openspec/config.yaml`)
**Work unit**: PR 3 / Slice 3, tasks 3.1–3.3 only
**Boundary**: Eager catalog cache, focused catalog tests, and gated PostgreSQL integration harness; Slice 3 handoff, wiring, ports, domain, schema, data, Docs, and adapters remain out of scope.

## Cumulative Completed Tasks
- [x] 1.1 Added nullable JSONB encode/decode helpers and wrapped PostgreSQL/not-found errors.
- [x] 1.2 Added exactly `github.com/oklog/ulid/v2 v2.1.1`; implemented UTC monotonic generation under a mutex.
- [x] 1.3 Added concurrent ULID parseability, opacity, length, and uniqueness tests.
- [x] 1.4 Added lead create/read/CAS/list/message SQL with explicit stale-versus-missing CAS disambiguation, conjunctive filters, deterministic ordering, and tests for helper/query invariants.
- [x] 2.1 Added transactional plan create/save; supplied hitos use non-destructive ID upserts, omitted hitos remain, and no plan version/CAS is introduced.
- [x] 2.2 Added plan query-contract tests for atomic transaction structure, aggregate hito reconstruction, active/pending overdue filters, deterministic ordering, and missing-hito behavior.
- [x] 2.3 Added lead-keyed ficha upsert and retrieval with lead-versus-ficha absence disambiguation.
- [x] 2.4 Added ficha upsert query-contract tests for replacement and unique lead ownership.

## Files Changed in Slice 2
- `internal/infrastructure/postgres/plan_repository.go`
- `internal/infrastructure/postgres/plan_repository_test.go`
- `internal/infrastructure/postgres/ficha_repository.go`
- `internal/infrastructure/postgres/ficha_repository_test.go`
- `openspec/changes/issue-15-postgres-repositories/tasks.md`
- `openspec/changes/issue-15-postgres-repositories/apply-progress.md`

## Work Unit Evidence
| Evidence | Result |
|---|---|
| Focused tests | `go test ./internal/infrastructure/postgres -run 'TestPlanRepository|TestFichaRepository'` — passed. |
| Full tests | `go test ./...` — passed. |
| Build | `go build ./...` — passed. |
| Static validation | `go vet ./...`, `gofmt -l` on Slice 2 files, and `git diff --check` — passed. |
| Runtime harness | N/A: Slice 2 has no wiring and the PostgreSQL integration harness is explicitly deferred to Slice 3. |
| Scope/budget | Slice 2 implementation and tests are 220 authored lines (118 plan + 29 plan tests + 55 ficha + 18 ficha tests), below the <=380 target; no dependencies added. |
| Rollback boundary | Revert only `plan_repository.go`, `plan_repository_test.go`, `ficha_repository.go`, and `ficha_repository_test.go`; no schema/data/public contract/wiring rollback is needed. |

## Deviations and Risks
- No production design deviation: plan writes use one transaction, upsert only supplied hitos, preserve omitted hitos, reconstruct ordered aggregates, and query active pending overdue hitos; ficha writes use the unique `lead_id` conflict key.
- Focused tests are SQL-contract tests because the real PostgreSQL harness is intentionally deferred to Slice 3; runtime database execution remains pending that harness.
- Task 1.5 is complete from the merged Slice 1 handoff. Task 2.5 remains unchecked for the orchestrator's PR 2 merge gate.

## Slice 3 Completed Tasks
- [x] 3.1 Added `NuevoCatalogo(fs.FS)` with eager JSON/brochure loading, duplicate and identity validation, typed affiliate/brochure misses, and defensive results.
- [x] 3.2 Added `t.TempDir` tests for missing source failure, no-repeat filesystem reads, mutation isolation, and both typed miss paths.
- [x] 3.3 Added `testdb_test.go`, gated by `testing.Short()` or absent `VIVI_TEST_DATABASE_URL`, exercising lead/plan/ficha repositories against the existing schema when configured.

## Files Changed in Slice 3
- `internal/infrastructure/postgres/catalogo_repository.go`
- `internal/infrastructure/postgres/catalogo_repository_test.go`
- `internal/infrastructure/postgres/testdb_test.go`
- `openspec/changes/issue-15-postgres-repositories/tasks.md`
- `openspec/changes/issue-15-postgres-repositories/apply-progress.md`

## Slice 3 Work Unit Evidence
| Evidence | Result |
|---|---|
| Focused tests | `go test ./internal/infrastructure/postgres -run 'TestCatalogoRepository' -count=1` — passed; real `data/` eager-load smoke test — passed. |
| Runtime harness | `go test -v ./internal/infrastructure/postgres -run '^TestPostgresIntegration$' -count=1` — passed with explicit skip: `VIVI_TEST_DATABASE_URL is not set`; `-short` gate also passed. |
| Full tests | `go test ./...` — passed. |
| Build | `go build ./...` — passed. |
| Static validation | `go vet ./...`, `gofmt -l` on Slice 3 files, and `git diff --check` — passed. |
| Scope/budget | 293 authored lines across the three Slice 3 Go files, below the <=330 target; no dependency changes. Scope audit found only the three permitted infrastructure files plus SDD evidence changes. |
| Rollback boundary | Revert only `catalogo_repository.go`, `catalogo_repository_test.go`, and `testdb_test.go`; no schema/data/public contract/wiring rollback is needed. |

## Deviations and Risks
- No production design deviation. Brochure files must expose `proyecto_id` front matter matching their filename and project entry; this validates the Contract §4 identity relationship eagerly.
- PostgreSQL runtime execution was gated and skipped because `VIVI_TEST_DATABASE_URL` is absent; the harness is ready to exercise existing lead/plan/ficha tables when configured.

## Delivery Completion
- [x] 1.5 PR #69 Slice 1 handoff and merge gate — CI-green merge `5d68e9b`.
- [x] 2.5 PR #70 Slice 2 handoff and merge gate — CI-green merge `6344829`.
- [x] 3.4 PR #71 Slice 3 handoff and merge gate — CI-green merge `5aae166`.
- [x] PR #72 runtime-coverage remediation — CI-green merge `def131c`.
- [x] PR #73 opaque-ID boundary coverage — CI-green merge `e75a163`.

## Final Verification Candidate
- Independent verification on merged `main` SHA `e75a1638b5e13a302c4000483a42c4e3b6b9f13f` found 6/6 requirements, 13/13 scenarios, and 14/14 tasks COMPLIANT; full/race tests, build, vet, module checks, import/scope audits, and main CI run `30168376184` passed.
- `TestPostgresIntegration` correctly skips when `VIVI_TEST_DATABASE_URL` is absent. Docker is unavailable in this WSL distro, so real-engine write behavior remains a documented environment-gated warning rather than a failed requirement.
- The candidate verification report is not persisted: installed `gentle-ai 1.43.3` lacks `sdd-verify-validate`, and the persistence contract forbids report admission without that validator. This blocks formal verification receipt and archive only; it does not invalidate the executed evidence above.

## Remediation Progress: Runtime Repository Coverage

**Change**: `issue-15-postgres-repositories`
**Mode**: Standard (strict TDD disabled by `openspec/config.yaml`)
**Work unit**: Verification remediation on `test/issue-15-runtime-coverage`
**Base**: merged main `5aae16615bb662b8c2b718427e4d4796301c1dd6`
**Boundary**: private pgx-shaped pool seam, test-only pgxmock dependency, and executing repository tests only. Public ports, domain, usecase contracts, adapters, composition root, migrations, data, Docs, and Issue #11 evidence remain untouched.

### Remediation Completed
- Added exact `github.com/pashagolub/pgxmock/v4 v4.9.0`; existing `github.com/jackc/pgx/v5 v5.7.6` was not upgraded.
- Added pgxmock and the missing ULID transitive `github.com/pborman/getopt` checksums to `go.sum`.
- Added private `pgxPool` in `internal/infrastructure/postgres/pool.go`; `*pgxpool.Pool` remains production-compatible and pgxmock is imported only by tests.
- Changed only lead, plan, and ficha repository constructors/fields to accept that private interface.
- Replaced tautological lead/plan/ficha SQL-literal tests with executing pgxmock tests covering generated opaque ID propagation, successful CAS, missing/stale CAS with caller-version immutability, conjunctive list filters and ordered reconstruction, chronological conversation, plan commit/rollback, aggregate reconstruction, overdue active/pending selection, missing hito, ficha upsert/retrieval, and lead-versus-ficha absence.

### Files Changed in Remediation
- `internal/infrastructure/postgres/pool.go`
- `internal/infrastructure/postgres/lead_repository.go`
- `internal/infrastructure/postgres/lead_repository_test.go`
- `internal/infrastructure/postgres/plan_repository.go`
- `internal/infrastructure/postgres/plan_repository_test.go`
- `internal/infrastructure/postgres/ficha_repository.go`
- `internal/infrastructure/postgres/ficha_repository_test.go`
- `go.mod`
- `go.sum`

### Work Unit Evidence
| Evidence | Result |
|---|---|
| Focused repository tests | `go test ./internal/infrastructure/postgres -run 'TestLeadRepository|TestPlanRepository|TestFichaRepository|TestRepositoryErrorMapping|TestJSONBHelpers' -count=1` — passed. |
| Focused ID tests | `go test ./internal/infrastructure/ids -run 'TestULID' -count=1` — passed. |
| Focused catalog tests | `go test ./internal/infrastructure/postgres -run 'TestCatalogoRepository' -count=1` — passed. |
| Runtime harness | `go test -v ./internal/infrastructure/postgres -run '^TestPostgresIntegration$' -count=1` — passed with explicit skip: `VIVI_TEST_DATABASE_URL is not set`. |
| Full tests | `go test ./...` — passed. |
| Build | `go build ./...` — passed. |
| Static validation | `go vet ./...`, `gofmt`, and `git diff --check` — passed. |
| Dependency proof | `go list -m -versions github.com/pashagolub/pgxmock/v4` lists v4.9.0; `go mod verify` — all modules verified. |
| Scope/budget | Tracked remediation diff is 222 additions and 61 deletions; new `pool.go` adds 17 lines, for 239 additions and 61 deletions (300 authored changed lines), below 400. No Docs, migration, data, public contract, adapter, composition-root, task, or state edits were made by remediation. |
| Rollback boundary | Revert only the private pool seam, three repository constructor type changes, three executing test replacements, and pgxmock/checksum dependency metadata; no schema/data/public-contract rollback is needed. |

### Task Relationship
The main `tasks.md` checklist and prior Slice 3 evidence were preserved unchanged. This remediation closes verification evidence gaps only; it does not alter completed main-task marks or claim the original integration harness ran against PostgreSQL.

### Risks and Deviations
- No production behavior or public architecture deviation: the interface is private and exposes only methods already used by repositories.
- PostgreSQL integration remains gated and skipped because `VIVI_TEST_DATABASE_URL` is absent; pgxmock now supplies executable repository behavior coverage without importing pgxmock in production source.

### Remediation Status
All requested remediation work is complete and ready for a fresh `sdd-verify` attempt. No files were staged or committed, and no PR was opened.
