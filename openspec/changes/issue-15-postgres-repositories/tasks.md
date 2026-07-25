# Tasks: PostgreSQL Repositories and ULID ID Generator

## Scope Guard
Change only infrastructure, focused tests, and `go.mod`/`go.sum`; never public ports, domain, schema/migrations, data, Docs, adapters, or wiring.

## Review Workload Forecast

| Field | Value |
|---|---|
| Delivery | Three independently mergeable PRs, sequentially to `main` |
| PR estimates | Slice 1 <=360; Slice 2 <=380; Slice 3 <=330 authored additions+deletions |
| Chain strategy | `stacked-to-main`: merge each slice before beginning the next |
| Size exception | Not allowed |

Decision needed before apply: No
Chained PRs recommended: Yes
Chain strategy: stacked-to-main
400-line budget risk: High

### Suggested Work Units

| Unit | Focused tests | Runtime evidence | Rollback |
|---|---|---|---|
| PR 1: JSON/errors, ULID, leads (<=360) | `go test ./internal/infrastructure/ids -run TestULID`; `go test ./internal/infrastructure/postgres -run TestLeadRepository` | N/A: wiring and DB harness are deferred to PR 3 | Revert helpers, lead/ID files, tests, and ULID pin |
| PR 2: plans and fichas (<=380) | `go test ./internal/infrastructure/postgres -run 'TestPlanRepository|TestFichaRepository'` | N/A: no wiring; DB harness is in PR 3 | Revert plan/ficha files and tests |
| PR 3: catalog and integration harness (<=330) | `go test ./internal/infrastructure/postgres -run 'TestCatalogo|TestPostgresIntegration'` | Run `VIVI_TEST_DATABASE_URL=... go test ./internal/infrastructure/postgres -run TestPostgresIntegration`; otherwise record gated skip | Revert catalog and `testdb_test.go` |

## Phase 1: Slice 1 — JSON, Errors, ULID, and Leads

- [ ] 1.1 Create `postgres/jsonb.go` and `errors.go` with null-preserving JSONB and repository error mapping; test null and nested JSON.
- [ ] 1.2 Pin only `github.com/oklog/ulid/v2 v2.1.1` in `go.mod`/`go.sum`; create mutex-protected UTC monotonic `ids/ulid.go`.
- [ ] 1.3 Add `ids/ulid_test.go` for parseability, 26-character opaque output, <=40 length, and concurrent uniqueness.
- [ ] 1.4 Add `lead_repository.go` and tests for create/get, zero-row CAS disambiguation/no mutation, conjunctive queue order, and chronological messages.
- [ ] 1.5 PR 1 handoff: focused/full tests and build; clean-architecture/scope audit, <=360 lines; merge to `main` before Slice 2.

## Phase 2: Slice 2 — Plan and Ficha Transactions

- [ ] 2.1 Add transactional `plan_repository.go` upserting supplied hitos while retaining omitted hitos; never add plan CAS or destructive replacement.
- [ ] 2.2 Add plan tests for rollback, aggregate reconstruction, overdue active/pending ordering, and missing-hito not-found behavior.
- [ ] 2.3 Add `ficha_repository.go` with lead-keyed `ON CONFLICT` replacement and lead-versus-ficha absence disambiguation.
- [ ] 2.4 Add ficha tests for one-row upsert/retrieval and both typed not-found cases.
- [ ] 2.5 PR 2 handoff: focused/full tests and build; clean-architecture/scope audit, <=380 lines; merge to `main` before Slice 3.

## Phase 3: Slice 3 — Catalog Cache and Integration Harness

- [ ] 3.1 Add eager `catalogo_repository.go` `fs.FS` validation, duplicate/identity checks, typed misses, and defensive copies.
- [ ] 3.2 Add `t.TempDir` catalog tests for load failure, no-repeat I/O, mutation isolation, and affiliate/brochure misses.
- [ ] 3.3 Add gated `testdb_test.go`; skip under `-short` or absent `VIVI_TEST_DATABASE_URL`, otherwise prove repository integration.
- [ ] 3.4 PR 3 handoff: focused/full tests and build; record integration result or justified skip, clean-architecture/scope audit, <=330 lines; merge to `main`.
