# Proposal: Issue #1 — Repository Foundation (go.mod, Clean Architecture tree, Makefile)

Issue: [#1 `[F1] Estructura del repositorio, go.mod y Makefile`](https://github.com/Unikyri/vivi-perfilamiento-leads/issues/1) · Branch: `feat/repo-foundation` → `main`

## Intent

The repository holds only documentation and local tooling: no `go.mod`, no command entry points, no `internal/` tree, no build automation. Issue #1 is a hard blocker for every other issue. This change creates the compile-safe Go 1.24 skeleton exactly as the issue specifies, so both blocks (A núcleo, B datos) can start in parallel without inventing layout or build commands.

## Scope

### In Scope (exact Issue #1 deliverables)
- `go.mod`: module `github.com/Unikyri/vivi-perfilamiento-leads`, version line exactly `go 1.24`.
- 8 `doc.go` package docs with the issue's literal contents: `internal/domain`, `internal/domain/motor`, `internal/usecase`, `internal/adapters/http`, `internal/adapters/agentes`, `internal/infrastructure/{postgres,llm,config}`.
- `cmd/servidor/main.go` printing `vivi servidor - placeholder`; `cmd/pipeline/main.go` printing `vivi pipeline - placeholder`.
- `.gitkeep` in `web/`, `data/`, `references/`, `migrations/`, `skills/`.
- `Makefile` with the issue's exact `.PHONY` targets and command bodies (`build test vet run datos limpiar`), TAB-indented.
- `.gitignore`: append `# Binarios compilados` then `bin/` at the end.
- SDD artifacts under `openspec/changes/issue-1-repo-foundation/` and the existing local `.kiro/` configuration, included in this consolidated delivery unchanged in behavior.

### Out of Scope (Non-Goals)
- `Docs/` and `README.md` — must remain byte-identical.
- Any domain, motor, use-case, HTTP, ADK, PostgreSQL, LLM, or config implementation; no dependencies added to `go.mod`.
- Contract v1.1 (doc 10) and motor criteria (doc 13) behavior; API routes, schemas, migrations content, CI workflows, frontend code.
- Directories not listed by Issue #1 (e.g. `docs/`, `contracts/schemas/`) even though doc 11 §8 anticipates them.
- Tests: no Go test files exist yet; `make test` must succeed trivially, not be backfilled with coverage.

## Capabilities

### New Capabilities
- `repository-foundation`: module identity, Clean Architecture package skeleton and layer ownership docs, command placeholders, build automation targets, tracked-directory placeholders, and build-artifact ignore rules.

### Modified Capabilities
- None (`openspec/specs/` is empty; no existing capability changes).

## Architecture Boundary Constraints

| Authority | Constraint this change must honor |
|---|---|
| NFR-M-01 (doc 09 §5) | Dependency flow strictly inward. `internal/domain` imports nothing from outer layers; `internal/usecase` may import `domain` only; adapters/infrastructure sit outermost. Encoded in the `doc.go` comments verbatim. |
| Doc 11 §8 | Package tree matches the planned layout; foundation adds only the subset Issue #1 lists. |
| Doc 10 (Contract v1.1) | A/B boundary is the HTTP API plus data files. This change declares no endpoints and no contract types, so the boundary stays frozen. |
| Doc 13 | `internal/domain/motor` is documented as zero-dependency, zero-I/O, zero-LLM so the deterministic engine can later be implemented pure and table-testable. |
| Repo layout | Nested wiki repo `Docs/vivi-perfilamiento-leads.wiki/` and `.atl/` remain ignored; `.gitignore` is appended to, never rewritten. |

## Approach

Minimal contract-aligned scaffold (exploration option 1): transcribe Issue #1 literally rather than generating idiomatic variants. Two placeholder `main` packages depend only on `fmt`, so `go build ./...` and `go vet ./...` pass with an empty dependency graph. Layer boundaries are enforced socially (package docs) now and mechanically (imports) by later issues.

## Affected Areas

| Area | Impact | Description |
|---|---|---|
| `go.mod` | New | Module path + `go 1.24` |
| `Makefile` | New | 6 targets, exact command bodies |
| `cmd/servidor`, `cmd/pipeline` | New | Compile-safe placeholders |
| `internal/**/doc.go` | New | 8 package docs (layer ownership) |
| `web`, `data`, `references`, `migrations`, `skills` | New | `.gitkeep` only |
| `.gitignore` | Modified | Append `bin/` rule |
| `openspec/changes/issue-1-repo-foundation/`, `.kiro/` | New | SDD artifacts + existing local config, behavior unchanged |
| `Docs/`, `README.md` | Unchanged | Prohibited paths |

## Validation Plan

Issue #1 Definition of Done, run from repo root after implementation (no TDD; `strict_tdd: false`):

1. `go build ./...` → exit 0, no output.
2. `go vet ./...` → exit 0.
3. `make build` → produces `bin/servidor` and `bin/pipeline`.
4. `ls bin/servidor bin/pipeline` → both exist.
5. `go run ./cmd/servidor` → stdout exactly `vivi servidor - placeholder`.
6. `go run ./cmd/pipeline` → stdout exactly `vivi pipeline - placeholder` (supporting check).
7. `make test` → exit 0 with no test files (`go test ./... -count=1`); no coverage claim.
8. `git status` → `bin/` untracked-ignored; `Docs/` and `README.md` show zero diff and remain unstaged (`Docs/` currently holds 45 untracked binaries under `Docs/Indumentos/`, so staging must be path-explicit — never `git add .`).

Verification must report Go toolchain availability explicitly; if `go 1.24` is absent locally, the phase reports the blocker instead of downgrading the version line.

## Delivery and Review-Size Forecast

- Strategy: single PR `feat/repo-foundation` → `main`, per state `delivery.strategy`.
- Authored Go/build surface: ~80 changed lines (16 new source/config files + 3 `.gitignore` lines) — well inside 400.
- Consolidated non-behavioral surface: `.kiro/` ≈ 5,054 lines / 52 files plus `openspec/` ≈ 170 lines, all currently untracked. This dominates the diff.
- `Decision needed before apply: Yes`
- `Chained PRs recommended: Yes`
- `400-line budget risk: High`
- Recommended resolution for the orchestrator/user: either (a) split into slice 1 = Issue #1 Go foundation (~80 lines, reviewable) and slice 2 = `.kiro/` + `openspec/` tooling and SDD artifacts, or (b) accept an explicit `size:exception` on the grounds that the oversized portion is pre-existing agent configuration with no runtime effect. `sdd-tasks` must carry this decision forward.

## Risks

| Risk | Likelihood | Mitigation |
|---|---|---|
| Makefile indented with spaces → `make` fails | Med | Write literal TABs; validate with `make build` in verify |
| Local Go toolchain older than 1.24 | Med | Detect before apply; report blocker, never edit the required version line |
| `.kiro/` bulk makes the PR unreviewable | High | Explicit size decision above before apply |
| Placeholder `main.go` drifting from exact strings | Low | Spec pins exact stdout strings; verify asserts them |
| Directory over-creation beyond Issue #1 list | Med | Spec enumerates the closed file list; extras are a verify failure |
| Accidental `Docs/`/`README.md` edit | Low | Prohibited paths in state; `git status` diff check in validation |
| Blanket `git add .` stages 45 untracked binary files under `Docs/Indumentos/` (pptx/xlsx/docx + `:Zone.Identifier`) | High | Never `git add .`; stage only the explicit Issue #1 paths plus `.kiro/` and `openspec/` |
| Layer boundaries only documented, not enforced | Med | Accepted for F1; note import-lint as follow-up issue |

## Rollback Plan

All work is additive on `feat/repo-foundation`; `main` is untouched until PR merge.

- Pre-merge: `git restore --staged --worktree` the added paths, or delete the branch — no other change depends on it yet.
- Post-merge: revert the single squash-merge commit (`git revert -m 1 <sha>`), then `rm -rf bin/`. Reverting restores the documentation-only baseline; the only externally visible loss is the module declaration, which nothing else consumes yet.
- `.gitignore` rollback removes only the appended two lines.

## Dependencies

- Go 1.24+ toolchain and GNU Make on the implementing machine.
- No upstream issues (Issue #1 depends on nothing; it blocks all others).

## Success Criteria

- [ ] Every file in the Issue #1 list exists with the exact specified contents; nothing extra is created.
- [ ] `go build ./...`, `go vet ./...`, `make build` all exit 0; `bin/servidor` and `bin/pipeline` exist.
- [ ] `go run ./cmd/servidor` prints exactly `vivi servidor - placeholder`.
- [ ] `bin/` is ignored via the appended `# Binarios compilados` block.
- [ ] `Docs/` and `README.md` are unmodified.
- [ ] PR opened to `main` with the size decision recorded.

## Proposal question round

Interactive questioning was not possible in this delegated phase. Assumptions needing user confirmation (answer, correct, or skip):

1. Delivery size: split the `.kiro/`/`openspec/` tooling into its own PR, or accept `size:exception` for one consolidated PR? (Assumed: user wants one consolidated PR → `size:exception`.)
2. Should Issue #1 also create `docs/` and `contracts/schemas/` from doc 11 §8, or stay strictly on the issue's closed list? (Assumed: strictly the issue list.)
3. Is a Go-import boundary linter wanted now to enforce NFR-M-01, or deferred to a later issue? (Assumed: deferred.)
4. Is `make test` passing with zero test files acceptable as the F1 baseline? (Assumed: yes, `strict_tdd: false`.)
