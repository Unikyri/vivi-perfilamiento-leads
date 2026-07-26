# Design: Issue #27 Deployable Demo

## Technical Approach

Deliver five independently revertible slices: S1 aligns CI and reproduces the Node→Go production build; S2 embeds Vite output without shadowing API routes; S3 loads three stable synthetic leads and restores them transactionally; S4 documents the judged flow and operations; S5 adds a manual, default-no-op promotion gate. No phase deploys, provisions credentials, or adopts ADK.

## Architecture Decisions

| Decision | Choice and rationale | Rejected alternatives |
|---|---|---|
| Build order | Root `package.json`/lockfile exposes `heroku-postbuild`, running `npm --prefix web ci && npm --prefix web run build && npm --prefix web run verify:build`. The Node lifecycle verifies `index.html`, hashed assets, and every local entry reference before the ordered Node→Go buildpacks proceed; CI and Make use the same lifecycle. The generated `estaticos/` tree is versioned so `go:embed` cannot package an ignored placeholder or stale hash. | Ignored generated output permits a clean checkout to embed an entry whose hashed files are absent; Go-first builds can embed the placeholder. |
| Static routing | `go:embed all:estaticos` plus `fs.Sub`; register `GET /` after existing routes. Serve files first, fall back to `index.html` only for extensionless paths, and return 404 for missing assets. Add CSP and `X-Content-Type-Options: nosniff`. Go `ServeMux` specificity preserves `/salud` and `/api/`. | Catch-all fallback for every miss can hide broken assets; proxying Vite is development-only. |
| Demo seed | `usecase` owns canonical `ana/carlos/luisa` DTOs with stable IDs. A narrow Postgres seed port uses conflict-safe inserts; startup calls it only when `DEMO_SEED=true`, with a three-second timeout and fail-fast error. Reset deletes app-owned rows and reinserts seed/date in its existing transaction. | Calling the HTTP handler creates duplicates; LLM/ADK seeding is nondeterministic; a follow-up seed after reset is non-atomic. |
| Promotion | Native Heroku GitHub integration deploys `main` after required CI; `app.json` records Node-before-Go buildpacks, one Basic web dyno, and `DEMO_SEED=false`. Public `/salud` verification is documented. No custom workflow, deployment action, provider call, or credential value is added. | Push-triggered custom workflows or repository credentials can promote accidentally/leak secrets. |

## Data Flow

```mermaid
sequenceDiagram
  participant CI
  participant Node
  participant Go
  participant App
  CI->>Node: npm ci + Vite build
  Node-->>Go: hashed assets in estaticos/
  CI->>Go: presence check + test + embed/build
  App->>App: DEMO_SEED=false => skip
  App->>DB: enabled => idempotent seed (<3s)
  Browser->>App: GET /route
  App-->>Browser: asset or SPA index; API retains precedence
  Operator->>GitHub: manual promote=true
  GitHub->>GitHub: validate + approval + secret/ref gates
  GitHub->>Heroku: exact main SHA; Node buildpack then Go
```

## File Changes

| Slice | Files | Action |
|---|---|---|
| S1 | `.github/workflows/ci.yml`, `Makefile`; root `package.json`, `package-lock.json` | Modify/create: Go version from `go.mod`; shared Node→Go build and asset check. |
| S2 | `internal/adapters/http/estaticos.go`, `estaticos_test.go`; `cmd/servidor/main.go`, `internal/adapters/http/estaticos/` | Create/modify: embed handler, route tests, Bloque B registration, and versioned Vite output whose index references are checked against embedded files. |
| S3 | `internal/usecase/seed.go`, `seed_test.go`, `puertos.go`, `reiniciar_demo.go`; `internal/adapters/http/leads.go`, `leads_test.go`; `internal/infrastructure/postgres/demo_repository.go`, `demo_repository_test.go`; `cmd/servidor/main.go`; `.env.example` | Create/modify: canonical seeds, gated startup, atomic reset. |
| S4 | `docs/guion-demo.md` | Create: both Doc 12 §4 scripts, real-API acceptance, config/rollback/anti-sleep matrix; label `?mock=1` visual-only. |
| S5 | `app.json`, `docs/guion-demo.md` | Create/modify: secret-free release metadata and native Heroku GitHub promotion/public health guidance; no custom deploy workflow. |

No files are deleted. S3 apply requires adding `internal/infrastructure/postgres` to its allowed edit roots.

## Interfaces / Contracts

`CargarSeed.Ejecutar(ctx)` is idempotent; `DemoSeedRepository` owns startup insert and atomic reset. `RegistrarEstaticos(mux) error` rejects a missing entry. Public API/schema remain unchanged.

## Testing Strategy

| Layer | Evidence |
|---|---|
| Build/CI | Clean checkout runs Node→Go; asset presence, `go vet`, race tests, Go build; missing output fails. |
| Unit | `httptest`: root/assets/fallback/404, CSP/nosniff, `/api` and `/salud` precedence; seed disabled/idempotent/timeout. |
| Integration | `pgxmock` proves one reset transaction and three stable inserts; gated Postgres test proves two loads/reset under three seconds. |
| Promotion | Local JSON/schema, documentation, secret-pattern, and diff checks; public `/salud` verification is documented but not invoked during SDD apply. |

## Threat Matrix

| Boundary | Applicability | Safe/failure behavior and RED test |
|---|---|---|
| Documentation-like paths | N/A | No path classifier executes repository files. |
| Git repository selection | N/A | Checkout is fixed to workflow repository/SHA; no `git -C` or path input. |
| Commit state | N/A | No index or commit mutation. |
| Push state | N/A | Pinned action/API owns promotion; no tracking branch, first-push, or refspec logic. |
| PR commands | N/A | Workflow neither creates nor edits PRs. |

## Migration / Rollout

No schema migration. Merge S1–S5 separately; rollback per slice. S5 records the already-confirmed native Heroku GitHub promotion, app metadata, and public `/salud` check without invoking the provider or storing credentials. Anti-sleep is documented operator choice, not automation.

## Open Questions

None blocking architecture. S4 must be checked against unavailable local Doc 12 before merge. S5 uses the maintainer-confirmed Heroku GitHub integration; no custom deployment workflow is required.
