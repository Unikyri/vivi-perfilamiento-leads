# Tasks: Issue #1 — Repository Foundation

## Review Workload Forecast

| Field | Value |
|---|---|
| Estimated changed lines | Product foundation: ~80; consolidated delivery: `.kiro/` ~5,054 plus all SDD/OpenSpec artifacts |
| 400-line budget risk | High (consolidated delivery) |
| Chained PRs recommended | No — single consolidated delivery accepted by explicit `size:exception` |
| Suggested split | Not used; the user authorized one consolidated delivery |
| Delivery strategy | `feat/repo-foundation` → PR → `main`; `size:exception` authorized |
| Chain strategy | not applicable |
| Size decision | `size:exception` explicitly authorized by the user for one consolidated delivery |

Decision needed before apply: No
Chained PRs recommended: No
Chain strategy: not applicable
400-line budget risk: Accepted by explicit `size:exception`

**Apply status: READY** — implementation is limited to the assigned checklist tasks; validation and delivery remain pending.

### Suggested Work Units

| Unit | Goal | Likely PR | Focused test command | Runtime harness | Rollback boundary |
|---|---|---|---|---|---|
| 1 | Issue #1 Go foundation (~80 lines) | `feat/repo-foundation` → `main` | `go build ./... && go vet ./...` | `go run ./cmd/servidor` and `./cmd/pipeline` | Enumerated foundation paths and appended ignore block |
| 2 | Pre-existing `.kiro/` and SDD/OpenSpec artifacts | Same approved consolidated PR | `git diff --check -- .kiro openspec` | N/A — no product runtime behavior | `.kiro/` and `openspec/` only |

## Phase 1: Build Scaffold

- [x] 1.1 Create `go.mod` with only `module github.com/Unikyri/vivi-perfilamiento-leads` and `go 1.24`.
- [x] 1.2 Create `cmd/servidor/main.go` and `cmd/pipeline/main.go` as `fmt`-only placeholders with their exact specified stdout.
- [x] 1.3 Create the eight literal package docs: `internal/domain{/motor}/doc.go`, `internal/usecase/doc.go`, `internal/adapters/{http,agentes}/doc.go`, and `internal/infrastructure/{postgres,llm,config}/doc.go`.
- [x] 1.4 Create only `web/`, `data/`, `references/`, `migrations/`, and `skills/` `.gitkeep` placeholders.

## Phase 2: Apply Ignore/Build Tooling

- [x] 2.1 Create root `Makefile` with the exact six `.PHONY` targets and TAB-indented recipes from the specification.
- [x] 2.2 Append exactly `# Binarios compilados` and `bin/` to `.gitignore`; do not rewrite prior rules.

## Phase 3: Validate

- [x] 3.1 With Go 1.24+, run `go build ./...`, `go vet ./...`, `make build`, `make test`, `make vet`, `make run`, `make datos`, and `make limpiar`; add no tests.
- [x] 3.2 Verify both binary paths, exact placeholder stdout, and `git check-ignore bin/servidor` after `make build`.
- [x] 3.3 Confirm no extra Issue #1 paths and zero diff for `Docs/` and `README.md`.

## Phase 4: Stage/PR Delivery

- [x] 4.1 User decision recorded: explicit `size:exception` for one consolidated delivery. (Complete decision; delivery execution remains pending.)
- [x] 4.2 Path-explicitly staged only approved foundation paths plus `.kiro/` and `openspec/`; `Docs/` and `README.md` are absent from the index and `git add .` was never used.
- [x] 4.3 Opened approved PR [#33](https://github.com/Unikyri/vivi-perfilamiento-leads/pull/33) from `feat/repo-foundation` to `main`, linked to Issue #1 with the `type:feature` label.
