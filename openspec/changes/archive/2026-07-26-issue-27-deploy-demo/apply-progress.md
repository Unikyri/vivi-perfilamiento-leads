# Apply Progress: Issue #27 — Deployable demo

**Change:** `issue-27-deploy-demo`
**Mode:** Standard (`strict_tdd: false`)
**Artifact store:** Hybrid (OpenSpec + Engram)
**Delivery:** Bounded S5 slice (`s5-native-heroku-promotion`)

## Completed Tasks

- [x] 1.1 Add the root Node manifest/lockfile and make the production build package `web` before compiling Go.
- [x] 1.2 Align CI with `go.mod`, package the frontend before Go, check the entrypoint and hashed assets, then run vet, race tests, and build.
- [x] 2.1 Add `internal/adapters/http/estaticos_test.go` covering embedded entry/assets, extensionless SPA fallback, missing-asset 404, CSP/nosniff, and `/api`/`/salud` precedence.
- [x] 2.2 Add `internal/adapters/http/estaticos.go` with deterministic `go:embed` file-first serving and wire it after existing routes in `cmd/servidor/main.go`.
- [x] 2.3 Keep Vite output generated and versioned: `internal/adapters/http/estaticos/` contains the reviewed embed tree, and the Node lifecycle verifies it before Go compilation.
- [x] 3.1 Confirm the destructive reset scope is exactly `fichas`, `hitos`, `planes`, `mensajes`, and `leads`; buyers remain untouched and no DDL is executed.
- [x] 3.2 Add focused coverage for default-off seed, canonical identities, bounded timeout, repeat reset, generic disabled `ERROR_INTERNO`, and atomic Postgres mock transactions.
- [x] 3.3 Implement canonical deterministic `ana`, `carlos`, and `luisa` seeds; add seed/reset ports and upsert repository transactions; gate startup and reset on `DEMO_SEED`; set `.env.example` to `DEMO_SEED=false`.
- [x] 4.1 Create `docs/guion-demo.md` with two real-API demo/recovery scripts, `/salud`, config, reset/rollback, and optional anti-sleep guidance.
- [x] 4.2 Label `?mock=1` as visual-only, record the unavailable Doc 12 source, and verify that no credential values were added.
- [x] 5.1 Record the maintainer-confirmed Heroku app `vivi-37863aed9d29`, native GitHub `main` auto-deploy after CI, Basic one-dyno plan, PostgreSQL attachment, and provider-managed config without copying secret values.
- [x] 5.2 Add secret-free `app.json` release metadata and document native Heroku GitHub promotion plus public `/salud` verification. Deliberately do not add `deploy.yml`, a duplicate GitHub deployment action, provider calls, or secrets.

## Files Changed in S5

| File | Action | What was done |
|---|---|---|
| `app.json` | Created | Declared Node-before-Go buildpacks, one Basic web dyno, Go package target, and `DEMO_SEED=false`; contains no credentials. |
| `docs/guion-demo.md` | Modified | Added confirmed native Heroku GitHub promotion flow, public `/salud`/SPA checks, rollback guidance, and explicit no-custom-workflow boundary. |
| `openspec/changes/issue-27-deploy-demo/design.md` | Modified | Reconciled S5 from a custom pinned action to the existing native Heroku GitHub integration. |
| `openspec/changes/issue-27-deploy-demo/tasks.md` | Modified | Records S5 as complete with the bounded low-risk scope and no `deploy.yml`. |
| `openspec/changes/issue-27-deploy-demo/state.yaml` | Modified | Records 12/12 tasks, no external blockers, and `sdd-verify` as next. |

## Work Unit Evidence

| Evidence | Result |
|---|---|
| Focused metadata/docs validation | `python3 -m json.tool app.json`, required metadata/docs checks, no custom `deploy.yml`, credential-pattern scan, and `git diff --check` — passed. |
| Go regression | `go test ./... -count=1` — passed across all packages. |
| Static analysis | `go vet ./...` — passed. |
| Production build | `make build-todo` — passed: Node/Vite packaged hashed assets, static presence passed, and both Go binaries built. Existing npm audit output reports one moderate and one high vulnerability; dependency versions were not changed. |
| Runtime harness | N/A for this S5 metadata/docs unit; no provider call, deploy, secret read/write, or runtime behavior change was permitted. The public `/salud` scenario is documented for the operator to run after native promotion. |
| Rollback boundary | Revert `app.json`, the S5 documentation additions, and S5 OpenSpec metadata; leave S1–S4 implementation intact. Native Heroku rollback remains operator-controlled and was not invoked. |

## Deviations from Design

The original S5 design called for `.github/workflows/deploy.yml` and a pinned deployment action. The maintainer confirmed an existing native Heroku GitHub integration that waits for CI and promotes `main`, so S5 was reconciled to use that mechanism instead. No custom workflow or action is added, and no external operation is performed.

## Issues Found

No blocking implementation issues. Doc 12 remains unavailable in the target checkout and reachable Git history; the guide retains the maintainer comparison note. `make build-todo` reports the pre-existing npm audit summary; dependency changes are outside this slice.

## Remaining Tasks

None. All 12 implementation tasks are complete; ready for `sdd-verify`.

**Status:** 12/12 tasks complete. S5 is complete and ready for verification.

## Focused verifier remediation

- [x] R1 Heroku Node lifecycle now runs `npm --prefix web run verify:build` after Vite and before the ordered Go buildpack; it checks the non-empty entry, hashed assets, and every local `src`/`href` reference.
- [x] R2 Removed the ignored static placeholder exception, kept the generated `internal/adapters/http/estaticos/` tree versioned, and added `TestPackagedIndexReferencesExistingEmbeddedFiles`.

### Validation Evidence

| Evidence | Result |
|---|---|
| Focused regression | `npm --prefix web run verify:build` — passed; a temporary missing-asset probe failed closed as expected. `go test ./internal/adapters/http/... -count=1` — passed. |
| Production build | `make build-todo` — passed: Node/Vite packaged and verified the SPA before both Go binaries compiled. |
| Full tests | `go test ./... -count=1` — passed. |
| Race tests | `go test ./... -count=1 -race` — passed. |
| Static analysis | `go vet ./...` — passed. |
| Diff hygiene | `git diff --check` and generated-SPA ignore check — passed. |
| Runtime harness | N/A — no deployment/provider call permitted; Heroku ordering is enforced by `app.json` buildpacks and the root `heroku-postbuild` gate. |
| Rollback boundary | Revert `package.json`, `web/package.json`, `.gitignore`, `internal/adapters/http/estaticos_test.go`, generated SPA files, and remediation SDD metadata; S1–S5 behavior remains otherwise unchanged. |

**Remediation status:** verifier warnings addressed; all 14 tracked implementation/remediation tasks complete; ready for fresh verification.
