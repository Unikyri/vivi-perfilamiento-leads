```yaml
schema: gentle-ai.verify-result/v1
evidence_revision: sha256:2b263a3bdf6f2177775dbbdb13594dda5239f68ee90faaa96b1552b85e78b0eb
verdict: pass_with_warnings
blockers: 0
critical_findings: 0
requirements: 8/8
scenarios: 12/12
test_command: go test ./... -count=1
test_exit_code: 0
test_output_hash: sha256:1975bd0b500c9dc5d6f2a3913c33b4c613a9790ffb79fb78a211335e412031f3
build_command: go build ./...
build_exit_code: 0
build_output_hash: sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
```

## Verification Report

**Change**: issue-27-deploy-demo
**Version**: delta specs `demo-control`, `despliegue-heroku`, `estaticos-spa`, `seed-demo` (N/A version tag)
**Mode**: Standard (`strict_tdd: false`)
**Workspace**: `/tmp/vivi-issue-27-archive` at merged `main` 852ba01, candidate tree `50668e718815c7d7206efaa8b9e10dc0872c1cc6`
**Review authority**: lineage `review-90f1e5d054fafb2a`, post-apply gate allow, binding revision `sha256:7ea734d0e6f873ad24547e54b1cf8f4277c654ce73461c9ebe3f98b5e768e4aa`
**Runtime attempt**: ordinal 1, revision `sha256:6d757e74071977f87a41d8519e9e768f1f83d4869cbfc1ec3537def5a5eb014c`, request `issue-27-final-verify-20260726`

### Completeness
| Metric | Value |
|--------|-------|
| Tasks total | 14 |
| Tasks complete | 14 |
| Tasks incomplete | 0 |

Twelve implementation tasks (S1–S5) plus remediation tasks R1 and R2 are checked in `tasks.md` and corroborated by `apply-progress.md`.

### Build & Tests Execution
**Build**: ✅ Passed

```text
$ go build ./...            # exit 0, no output
$ go vet ./...              # exit 0, no output
$ make build-todo           # reused harness: npm heroku-postbuild (vite build + verify:build)
                            # -> make check-static-assets -> go build bin/servidor, bin/pipeline
```

**Tests**: ✅ 14 packages ok / ❌ 0 failed / ⚠️ 0 failed-open (5 packages have no test files)

```text
$ go test ./... -count=1    # exit 0
ok  internal/adapters/http            0.118s
ok  internal/usecase                  0.012s
ok  internal/infrastructure/postgres  0.010s
ok  internal/infrastructure/config    0.008s
ok  cmd/servidor                      0.007s
(+ agentes, domain, domain/motor, bus, ids, llm, reloj, pipeline all ok)

reused harness: go test -race ./... -count=1  # exit 0, all packages ok
```

**Fail-closed asset-gate probes** (executed outside the repository, `/tmp/vivi27-probe`, no repository mutation):

```text
$ npm run verify:build   # index references a missing hashed asset -> exit 1
$ npm run verify:build   # empty index.html                        -> exit 1
```

**Coverage**: ➖ Not available — no coverage threshold configured for this change.

**Repository state after evidence**: `git status --porcelain` empty; tree still `50668e718815c7d7206efaa8b9e10dc0872c1cc6`. No product file was created, modified, or deleted during verification.

### Spec Compliance Matrix
| Requirement | Scenario | Test | Result |
|-------------|----------|------|--------|
| demo-control: Confirmed config-gated reset | Repeat | `internal/adapters/http/demo_reset_test.go > TestDemoResetReturnsSameDateTwice`; `internal/usecase/reiniciar_demo_test.go > TestReiniciarDemoIsIdempotentAndFast`; `internal/infrastructure/postgres/demo_repository_test.go > TestDemoRepositoryResetWithSeedCommitsDateAndLeadsTogether` | ✅ COMPLIANT |
| demo-control: Confirmed config-gated reset | Disabled | `internal/adapters/http/demo_reset_test.go > TestDemoResetDisabledIsGenericAndReadOnly`; `internal/usecase/reiniciar_demo_test.go > TestReiniciarDemoDisabledDoesNotMutate` | ✅ COMPLIANT |
| despliegue-heroku: Reproducible packaged build | Ordered build | `make build-todo` runtime (npm `heroku-postbuild` -> `check-static-assets` -> `go build`); `.github/workflows/ci.yml` `go-version-file: go.mod`; `go build ./...` exit 0 | ✅ COMPLIANT |
| despliegue-heroku: Reproducible packaged build | Missing assets | fail-closed `verify:build` probes (exit 1 for missing hashed asset and empty entry); `make check-static-assets`; `internal/adapters/http/estaticos_test.go > TestPackagedIndexReferencesExistingEmbeddedFiles` | ✅ COMPLIANT |
| despliegue-heroku: Inert secret-safe configuration | No promotion | `internal/infrastructure/config/config_test.go > TestCargarDemoSeedRequiresExplicitTrue`; `app.json` `DEMO_SEED=false`; no `deploy.yml` in `.github/workflows/`; credential-pattern scan clean | ✅ COMPLIANT |
| despliegue-heroku: Real acceptance and health guidance | Acceptance and ping | `cmd/servidor/main_test.go > TestSaludReportsLiveBreakerState`; `docs/guion-demo.md` real-API acceptance rule, `?mock=1` labeled visual-only, `/salud` verification, optional external ping | ✅ COMPLIANT (documentation dimension by inspection — see W2) |
| estaticos-spa: Packaged SPA | Entry and asset | `internal/adapters/http/estaticos_test.go > TestRegistrarEstaticosServesSPAAssetsAndPreservesRoutes/entry` and `/hashed asset` (byte-equal embedded content) | ✅ COMPLIANT |
| estaticos-spa: Packaged SPA | Missing assets | `make check-static-assets` and `verify:build` fail-closed probes; `TestPackagedIndexReferencesExistingEmbeddedFiles` | ✅ COMPLIANT |
| estaticos-spa: Safe fallback | Precedence | `TestRegistrarEstaticosServesSPAAssetsAndPreservesRoutes/api precedence`, `/health precedence`, `/browser fallback`, `/missing asset` plus `nosniff` and CSP assertions | ✅ COMPLIANT |
| seed-demo: Explicit demo-seed gate | Default startup | `internal/usecase/seed_test.go > TestCargarSeedIsDefaultOffAndCanonical` (disabled path, zero repository calls); `TestCargarDemoSeedRequiresExplicitTrue`; `cmd/servidor/main.go` gates on `cfg.DemoSeed` | ✅ COMPLIANT |
| seed-demo: Explicit demo-seed gate | Manual enablement | `TestCargarSeedIsDefaultOffAndCanonical` (enabled path loads exactly `ana`, `carlos`, `luisa` with `Fuente=DEMO`); `TestCargarDemoSeedRequiresExplicitTrue` explicit-true branch | ✅ COMPLIANT |
| seed-demo: Idempotent synthetic seed | Repeat | `TestReiniciarDemoIsIdempotentAndFast` (two runs, six seeded leads, under three seconds); `internal/infrastructure/postgres/demo_repository_test.go > TestDemoRepositorySeedsCanonicalLeadsInOneTransaction`; `internal/usecase/seed_test.go > TestCargarSeedHonorsContextTimeout` | ✅ COMPLIANT |

**Compliance summary**: 12/12 scenarios compliant

### Correctness (Static Evidence)
| Requirement | Status | Notes |
|------------|--------|-------|
| demo-control: Confirmed config-gated reset | ✅ Implemented | `ReiniciarDemo` returns `ErrDemoDeshabilitado` mapped to generic 500 `ERROR_INTERNO`; enabled path deletes only `fichas`, `hitos`, `planes`, `mensajes`, `leads`, reinserts approved date and the three seeds in one transaction; no DDL; buyers untouched. |
| despliegue-heroku: Reproducible packaged build | ✅ Implemented | Root `heroku-postbuild` runs `npm --prefix web ci && build && verify:build`; `Makefile` and CI share the Node-before-Go order; CI uses `go-version-file: go.mod`. |
| despliegue-heroku: Inert secret-safe configuration | ✅ Implemented | `app.json` declares Node-then-Go buildpacks, one Basic web dyno, `DEMO_SEED=false`; no credential values; no deploy workflow added. |
| despliegue-heroku: Real acceptance and health guidance | ✅ Implemented | `docs/guion-demo.md` documents both real-API scripts, `/salud` checks, config/rollback matrix, optional ping, and the visual-only mock label. |
| estaticos-spa: Packaged SPA | ✅ Implemented | `go:embed all:estaticos` with `fs.Sub`; `/` and hashed assets return embedded bytes; embed tree is versioned and reference-checked. |
| estaticos-spa: Safe fallback | ✅ Implemented | Extensionless paths fall back to the entry, missing assets return 404, static responses set `nosniff` and the NFR-S-03 CSP, and `/api/`+`/salud` keep their handlers. |
| seed-demo: Explicit demo-seed gate | ✅ Implemented | `config.Cargar` requires explicit `true`; `.env.example` and `app.json` default to `false`; startup seed is gated. |
| seed-demo: Idempotent synthetic seed | ✅ Implemented | `usecase.SemillasDemo()` yields canonical `ana`/`carlos`/`luisa`; repository uses conflict-safe inserts in one transaction with a bounded context. |

### Coherence (Design)
| Decision | Followed? | Notes |
|----------|-----------|-------|
| Node-before-Go build order verified by the Node lifecycle | ✅ Yes | `heroku-postbuild` chains `verify:build`; `Makefile` adds `check-static-assets` before `go build`. |
| `go:embed all:estaticos` + `fs.Sub`, route registered after existing routes | ✅ Yes | `RegistrarEstaticos` registers `GET /` last and rejects a nil mux / missing entry. |
| Versioned generated embed tree | ✅ Yes | `internal/adapters/http/estaticos/` is tracked; CI fails when it drifts from `web/src`. |
| Usecase-owned canonical seed, narrow Postgres port, atomic reset | ✅ Yes | `CargarSeed`, `DemoSeedRepository`, and `ReiniciarConSeed` verified by mock transaction ordering. |
| Promotion via native Heroku GitHub integration, no custom workflow | ⚠️ Documented deviation | Design was reconciled from the original pinned-action `deploy.yml` to the maintainer-confirmed native integration; `.github/workflows/` contains only `ci.yml`. Recorded in `design.md` and `apply-progress.md`. |

### Issues Found
**CRITICAL**: None

**WARNING**:
- W1 — No live PostgreSQL run: `internal/infrastructure/postgres` proves SQL shape, delete scope, and single-transaction behavior with `pgxmock`, and `TestPostgresIntegration` is skipped without `DATABASE_URL`. The three-second reset/seed budget is proven at the usecase and HTTP seams with fakes, not against a real database.
- W2 — `despliegue-heroku: Real acceptance and health guidance` is documentation-governed. Its executable part (`/salud`) has a passing test, but the acceptance policy and ping guidance were verified by inspecting `docs/guion-demo.md`; there is no automated documentation test.
- W3 — Public `/salud` verification against the deployed Heroku release was not executed: SDD phases forbid external provider calls. It remains an operator step in `docs/guion-demo.md`.
- W4 — The reused build harness reports two pre-existing `web` dependency vulnerabilities (1 moderate, 1 high) from `npm ci`. Dependency bumps are out of scope for this change.
- W5 — Wiki Doc 12 §4 is absent from the checkout and reachable Git history; the guide records this, and a maintainer must confirm the exact script names before relying on doc parity.

**SUGGESTION**:
- S1 — Add a gated live-Postgres job (or a documented manual runbook step) so the reset/seed three-second budget is measured against a real database before the next demo.
- S2 — Consider a small documentation lint (required sections in `docs/guion-demo.md`) so the acceptance-and-ping scenario gains automated coverage.

### Verdict
PASS WITH WARNINGS
All 14 tasks are complete, 8/8 requirements and 12/12 scenarios have passing runtime evidence on the merged candidate tree, `go test ./... -count=1` and `go build ./...` exit 0, and the five open warnings are non-blocking scope/environment limits that do not contradict any spec requirement.
