# Proposal: Issue #27 — Deployable demo (single dyno, seed, delivery script)

## Intent

Judges must open one URL and complete the flow unaided (US-17/US-18). Today `Procfile` points to an unbuilt `bin/servidor`, no route serves the SPA, `estaticos/` holds only a placeholder, CI Go 1.24.0 conflicts with `go.mod` 1.25.0, and there is no seed or demo script. Deployment is therefore unprovable.

## Scope

### In Scope (repository-only, no external calls)
- Align CI/Heroku Go toolchain with `go.mod`; add `build-todo` (web build before Go build) and an asset-presence check.
- `internal/adapters/http/estaticos.go`: `go:embed all:estaticos`, SPA fallback, `nosniff` + CSP (NFR-S-03), never shadowing `/api` or `/salud`; wire in the Bloque B slot of `main.go`.
- `internal/usecase/seed.go`: idempotent `CargarSeed` for `ana`/`carlos`/`luisa` (<3 s, NFR-D-02); reset restores seed.
- `docs/guion-demo.md` (Doc 12 §4 scripts) plus operator config matrix and anti-sleep options (NFR-D-01).
- `app.json` + `.github/workflows/deploy.yml` authored as inert, secret-referencing artifacts only.

### Out of Scope
- Running any deploy, Heroku/Gemini/Qwen call, or secret provisioning.
- ADK Go adoption, multi-dyno scaling, queue rework, auth/rate limiting, DB-backed catalog.
- Treating `?mock=1` as backend acceptance.

## Capabilities

### New Capabilities
- `estaticos-spa`: embedded SPA serving, fallback, security headers, API precedence.
- `seed-demo`: idempotent demo seed and its timing/enablement rules.
- `despliegue-heroku`: build order, config matrix, CI-gated promotion, `/salud` verification, no-secret policy.

### Modified Capabilities
- `demo-control`: reset must restore the seed, not only the approved date.

## Approach

Exploration Approach A, in five bounded slices, each ≤400 authored runtime/test lines: S1 CI/toolchain+build order → S2 embedded static serving → S3 seed+reset → S4 demo/ops docs → S5 credential-gated promotion (`app.json`, `deploy.yml`). S1–S4 are safe and mergeable now; S5 stays unexecuted until the maintainer confirms Heroku app ownership, Postgres attachment, and GitHub secrets (`HEROKU_API_KEY`, `HEROKU_APP_NAME`, `HEROKU_EMAIL`, `GEMINI_API_KEY`). `DEMO_SEED` defaults to `false` in code, `.env.example`, and `app.json`; only the operator enables it on the controlled demo app.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `.github/workflows/ci.yml`, `go.mod` | Modified | `go-version-file: go.mod`; single toolchain |
| `Makefile`, `.gitignore`, `web/vite.config.ts` | Modified | Deterministic embed input |
| `internal/adapters/http/estaticos.go`, `rutas.go` | New/Modified | SPA + headers |
| `cmd/servidor/main.go` | Modified | Bloque B wiring, optional seed |
| `internal/usecase/seed.go`, `reiniciar_demo.go` | New/Modified | Idempotent seed |
| `app.json`, `.github/workflows/deploy.yml`, `docs/guion-demo.md` | New | Ops artifacts |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Embed input empty (assets gitignored) | High | Explicit packaging decision + build-time asset check before S5 |
| Static route shadows API | Med | Route-precedence tests |
| Public unauthenticated demo + destructive reset | Med | One controlled app, `DEMO_SEED=false` default, synthetic data |
| Secret leakage | Low | Secret references only; gitleaks green |
| Deadline (Sun 11:00) | High | Ship S1–S4 independently |

## Rollback Plan

Per-slice revert: S1 restore `go-version: '1.24.0'` and Makefile target; S2 drop the `mux.Handle("/")` line and `estaticos.go` (API/`/salud` unaffected); S3 stop calling `CargarSeed` (reset falls back to date-only); S4 delete docs; S5 delete `deploy.yml`/`app.json` and redeploy the previous Heroku release, then re-check `/salud`. No migration or data change to reverse.

## Dependencies

- Issues #4, #5, #23, #26 (merged).
- Maintainer-provided Heroku/provider credentials for S5 (external, blocking).
- Reproducible `web` build from `package-lock.json`.

## Success Criteria

- [ ] `go vet`, `go test ./... -race`, and `make build-todo` pass on one toolchain.
- [ ] Local binary serves `/`, hashed assets, SPA fallback, `/salud`, and `/api/*` unchanged.
- [ ] Reset restores the three seed leads idempotently in <3 s when enabled.
- [ ] `docs/guion-demo.md` covers both Doc 12 §4 scripts and operator steps.
- [ ] `DEMO_SEED=false` default; gitleaks green; zero secrets in repo/artifacts.
- [ ] S5 remains inert until credentials are confirmed; no deploy performed in this phase.

## Proposal question round

1. Packaging decision: commit built SPA assets, or rely on the `heroku/nodejs`+`heroku/go` buildpack order with a root `heroku-postbuild`? (Blocks S2/S5.)
2. Confirm `DEMO_SEED=false` in `app.json` with a manual operator enable, overriding the issue snippet's `true`.
3. Should startup seed loading be gated by `DEMO_SEED`, keeping today's on-demand lead creation as the default path?
4. Is judge-facing acceptance the real API demo only, with `?mock=1` labeled visual fallback?
5. Anti-sleep choice: paid non-sleeping plan or external `/salud` ping?
