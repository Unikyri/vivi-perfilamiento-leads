# Issue #24 — Agent Skills Design

## Asset embedding
Go `embed` rejects patterns containing `..`, so a directive in `internal/adapters/agentes` cannot directly embed the repository-root `skills/` directory. The compatible solution is a minimal `skills` package with `//go:embed */SKILL.md`; `internal/adapters/agentes/skills.go` aliases that embedded `embed.FS` and retains the requested `CargarSkills(string) (string, error)` API.

## Loader behavior
The ownership map is a package-local map of agent names to ordered skill names. `CargarSkills` reads each asset from the embedded root, emits `=== SKILL: name ===`, strips the leading YAML block, and returns the accumulated body. No prompt builder, provider, factory, ADK, API, domain, or composition-root integration is introduced.

## Validation
Focused tests cover all mapped assets, required frontmatter and sections, body-only loading for every agent, and unknown-agent behavior. The tests never invoke an LLM. Go formatting, package tests, full tests, race tests, build, vet, and diff inspection are required before handoff.

## Rollback
Revert the Issue #24 skill documents, `skills/embed.go`, `internal/adapters/agentes/skills.go`, and `internal/adapters/agentes/skills_test.go`; no existing runtime integration is changed.
