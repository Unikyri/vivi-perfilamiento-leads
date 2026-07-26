# Tasks: Issue #27 Deployable Demo

**Readiness:** S1–S4 are complete and independently revertible. S5 is now safe because the maintainer confirmed the existing Heroku app, native GitHub promotion, single Basic dyno, PostgreSQL, and provider-managed configuration. No provider call or custom deployment workflow is part of this slice.

## Review Workload Forecast

| Field | Value |
|---|---|
| Estimated changed lines | S1 140–180; S2 180–240; S3 280–340; S4 60–90; S5 50–80 |
| 400-line budget risk | Low for this S5 slice |
| Chained PRs recommended | No for this bounded slice |
| Suggested split | S1 → S2 → S3 → S4 → S5 |
| Delivery strategy / chain | ask-on-risk / current S5 slice explicitly scoped |

Decision needed before apply: No
Chained PRs recommended: No
Chain strategy: not applicable to this bounded S5 slice
400-line budget risk: Low

### Suggested Work Units

| Unit | Goal | PR | Focused test | Runtime harness | Rollback |
|---|---|---|---|---|---|
| S1 | Node→Go build | 1 | `make build-todo && go vet ./... && go test ./... -race` | Clean checkout; absent output fails | CI, Makefile, root npm files |
| S2 | Embedded SPA | 2 | `go test ./internal/adapters/http/...` | `/`, asset, browser path, `/salud`, `/api/*` | handler and registration |
| S3 | Seed/reset | 3 | `go test ./internal/usecase/... ./internal/adapters/http/... ./internal/infrastructure/postgres/...` | Approved scratch DB: repeat in <3s | seed/repository/wiring |
| S4 | Demo guide | 4 | documentation review | Real API + `/salud`; mock visual-only | `docs/guion-demo.md` |
| S5 | Native Heroku release metadata and public verification guide | 5 | JSON/schema, docs, secret-pattern, and diff checks | N/A: no external provider call; guide records the public `/salud` check | `app.json`, S5 documentation/task metadata |

## Phase 1: S1 — Packaged build

- [x] 1.1 Update root `package.json`/`package-lock.json` and `Makefile`: build `web` before Go.
- [x] 1.2 Update `.github/workflows/ci.yml`: Go from `go.mod`, entry/hashed-asset check, vet, race test, and build; missing output fails.

## Phase 2: S2 — Embedded SPA

- [x] 2.1 Add RED `internal/adapters/http/estaticos_test.go`: entry, asset, fallback, 404, CSP/nosniff, and `/api`/`/salud` precedence.
- [x] 2.2 Create `internal/adapters/http/estaticos.go`; register file-first routing after existing routes in `cmd/servidor/main.go`.
- [x] 2.3 Use generated local assets only and version the reviewed `estaticos/` embed tree; the build lifecycle verifies it before Go compilation.

## Phase 3: S3 — Controlled seed

- [x] 3.1 Before destructive runtime testing, confirm the reset scope is only `fichas`, `hitos`, `planes`, `mensajes`, and `leads`; preserve buyer data and use no DDL.
- [x] 3.2 Add RED tests in `seed_test.go`, `leads_test.go`, and `demo_repository_test.go`: default-off, repeat, timeout, disabled `ERROR_INTERNO`, atomic reset.
- [x] 3.3 Implement canonical `ana`/`carlos`/`luisa`, port/repository, gated startup/reset in design-listed files; default `.env.example` `DEMO_SEED=false`.

## Phase 4: S4 — Demo guide

- [x] 4.1 Create `docs/guion-demo.md`: both Doc 12 §4 flows, real-API acceptance, `/salud`, config/rollback, and optional ping.
- [x] 4.2 Label `?mock=1` visual-only; check Doc 12 §4 before merge and verify no credential values.

## Phase 5: S5 — Native Heroku promotion reconciliation

- [x] 5.1 Record the maintainer-confirmed Heroku app `vivi-37863aed9d29`, native GitHub `main` auto-deploy after CI, Basic one-dyno plan, PostgreSQL attachment, and provider-managed config without copying secret values.
- [x] 5.2 Add secret-free `app.json` release metadata and document native Heroku GitHub promotion plus public `/salud` verification. Do not add `deploy.yml`, a duplicate GitHub deployment action, provider calls, or secrets.

## Focused verifier remediation

- [x] R1 Make the Heroku Node lifecycle run the Vite build and fail before Go compilation when `index.html`, hashed assets, or any local `src`/`href` reference is missing.
- [x] R2 Version the generated `internal/adapters/http/estaticos/` tree and add regression coverage so an ignored placeholder cannot mask stale or dirty embedded output.
