# Apply Progress — Issue #24 Agent Skills

- Change: `issue-24-agent-skills`
- Mode: Standard (`strict_tdd: false`)
- Delivery: one autonomous slice
- Status: implementation and validation complete; ready for independent verification

## Completed tasks
- [x] 1.1 Eight Issue #24 skill documents created with the required frontmatter, sections, examples, edge cases, and exact ownership metadata.
- [x] 1.2 Compatible-root embed package and `CargarSkills` loader created; no prohibited runtime areas changed.
- [x] 1.3 Focused tests cover all eight documents, frontmatter/template sections, every mapped agent, and unknown-agent behavior.
- [x] 1.4 Focused tests passed; full/race/build/vet/diff are required handoff checks.

## Work Unit Evidence

| Evidence | Result |
|---|---|
| Focused test | `go test ./internal/adapters/agentes/... -run 'Test(Skills|CargarSkills)' -v` — exit 0; all Issue #24 tests passed |
| Runtime harness | N/A — this slice is static prompt assets and an offline embed loader; no live LLM or runtime boundary is permitted |
| Rollback boundary | Revert the eight `skills/*/SKILL.md` files, `skills/embed.go`, `internal/adapters/agentes/skills.go`, and `skills_test.go` |

## Design deviation
The issue's illustrative ancestor `go:embed` pattern is invalid in Go. `skills/embed.go` is the minimal compatible-root embed package; the public `agentes.CargarSkills` API and ownership map remain unchanged.
