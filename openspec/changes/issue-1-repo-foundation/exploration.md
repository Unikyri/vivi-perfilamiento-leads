## Exploration: Issue #1 — Repository Foundation

### Current State
The repository is intentionally at a documentation-and-configuration baseline. The root contains `LICENSE`, `.gitignore`, `README.md`, `.kiro/`, `.atl/`, and `openspec/`; there is currently no `go.mod`, `Makefile`, Go command entry point, `internal/` package tree, `web/`, `skills/`, `data/`, `references/`, `migrations/`, or Go test suite. The existing `.gitignore` covers common Go outputs, environment files, the local `.atl/` directory, and the nested wiki repository, but does not yet ignore `bin/`.

Issue #1 is therefore a repository-foundation change, not a behavior change. Its implementation should create the Go 1.24 module declaration, the exact Makefile and command placeholders requested by the issue, package documentation for the planned Clean Architecture layers, the required empty directories, and the `bin/` ignore rule. It must not change `Docs/` or `README.md`.

The governing documents establish the forward-compatibility constraints. Wiki document 11 section 8 defines the planned layout: `cmd/servidor`, `cmd/pipeline`, `internal/domain`, `internal/usecase`, `internal/adapters/http`, `internal/adapters/agentes`, `internal/infrastructure/postgres`, `internal/infrastructure/llm`, `web`, `skills`, `data`, `references`, `migrations`, `docs`, `contracts/schemas`, and CI configuration. NFR-M-01 requires strict inward dependency flow: domain has no outward dependencies; usecase depends only on domain; adapters and infrastructure implement inward-facing ports. Document 10 is the future A/B boundary and must remain untouched. Document 13 governs later deterministic motor behavior, but this foundation must leave room for its pure, I/O-free implementation and table-driven tests.

### Affected Areas
- `go.mod` — new Go 1.24 module foundation; module path and toolchain declaration must match the issue's exact acceptance criteria.
- `Makefile` — new canonical developer commands and placeholders; exact command bodies must come from Issue #1 rather than being inferred in this exploration.
- `cmd/servidor/main.go` — placeholder composition root for the HTTP/API binary; should compile without prematurely implementing application behavior.
- `cmd/pipeline/main.go` — placeholder pipeline entry point for the future B→A data contract.
- `internal/domain/`, `internal/usecase/`, `internal/adapters/`, `internal/infrastructure/` — package documentation and directory skeleton preserving NFR-M-01's dependency direction.
- `web/`, `skills/`, `data/`, `references/`, `migrations/`, `docs/`, `contracts/schemas/` — required project directories, likely retained with `.gitkeep` or equivalent non-behavioral placeholders where Git would otherwise omit empty directories.
- `.gitignore` — add the scoped `bin/` build-artifact rule without disturbing existing local configuration or wiki exclusions.
- `openspec/changes/issue-1-repo-foundation/exploration.md` — this exploration artifact; the only SDD artifact permitted in the change folder during this phase.
- `openspec/changes/issue-1-repo-foundation/state.yaml` — phase DAG state must advance from `explore` pending to completed and recommend proposal.

### Approaches
1. **Minimal contract-aligned scaffold** — create only the module, exact Makefile targets, compile-safe command placeholders, package documentation, required directories, and `bin/` ignore rule.
   - Pros: smallest reviewable change; preserves `main` deployability; avoids inventing behavior or dependencies; directly traceable to Issue #1; keeps later SDD phases responsible for specifications and design details.
   - Cons: command placeholders provide little runtime functionality; some directories need placeholder files solely for Git tracking.
   - Effort: Low/Medium.

2. **Foundation plus initial Clean Architecture implementation** — create the scaffold and begin domain entities, ports, HTTP wiring, configuration, or database adapters.
   - Pros: produces an immediately executable skeleton; may reduce later bootstrapping work.
   - Cons: exceeds Issue #1; creates premature API and contract decisions; risks violating the frozen A/B boundary; substantially expands review and validation scope; would mix foundation with behavior owned by later issues.
   - Effort: High.

3. **Tooling-first repository setup** — prioritize a richer Makefile, CI, lint configuration, dependency manifests, and generated project metadata before adding the requested directory/package skeleton.
   - Pros: can standardize developer workflows early and expose validation failures sooner.
   - Cons: not the stated Issue #1 scope; likely invents exact commands and tools not present in the current repository; can create infrastructure coupling before the package boundaries are specified.
   - Effort: Medium/High.

### Recommendation
Use the minimal contract-aligned scaffold. Proposal and downstream specification phases should transcribe Issue #1's exact Makefile targets and placeholder expectations, then implement only those artifacts. Keep command binaries compile-safe and behavior-free, document the Clean Architecture package ownership without introducing cross-layer imports, and add only the required `bin/` ignore rule. Do not edit `Docs/`, `README.md`, Contract v1.1, or motor criteria. Because the current repository has no Go tooling or tests, apply should establish the smallest validation baseline available (`go fmt`, `go vet`, `go test ./...`, and the issue's Makefile checks if available); verification must report unavailable commands explicitly rather than inventing test coverage.

The estimated authored change is below the 400-line review budget if package documentation and Makefile placeholders remain concise. The feature branch delivery strategy is appropriate as a single PR to `main`, subject to the tasks-phase forecast.

### Risks
- The exact Issue #1 Makefile contract is not present in the checked repository context; proposal/spec must use the authoritative GitHub issue text and must not guess target names or command bodies.
- An incorrect Go module path would break future imports and command builds; it must be confirmed from repository/module conventions before implementation.
- Placeholder `main.go` files can accidentally introduce dependencies or fake behavior; they should remain compile-safe and minimal.
- Empty directory placeholders can create noisy or misleading package structure; use them only where the issue explicitly requires directory presence, and document package ownership clearly.
- Future contributors could bypass Clean Architecture boundaries if package documentation is vague; docs should state inward dependency direction and keep framework/database concerns out of `internal/domain` and `internal/usecase`.
- The local wiki is a nested repository and must remain untouched; `.gitignore` changes must preserve its existing exclusion.

### Ready for Proposal
Yes. The change is sufficiently understood to proceed to `sdd-propose`, with one required input check: obtain the exact GitHub Issue #1 acceptance text for the Makefile, module path, command placeholder contents, and directory list before finalizing the proposal.
