# Issue #24 — Agent Skills Tasks

## Review Workload Forecast

| Field | Value |
|---|---|
| Go runtime/test lines | 157 authored lines maximum for the loader, embed helper, and focused tests |
| Decision needed before apply | No |
| Chained PRs recommended | No |
| 400-line budget risk | Low for the requested Go slice |
| Delivery | One autonomous Issue #24 slice |

## Implementation

- [x] 1.1 Create the eight Issue #24 `SKILL.md` documents with exact frontmatter fields, required sections, examples, edge cases, ownership, and domain rules.
- [x] 1.2 Add a compatible-root Go embed package and `internal/adapters/agentes/skills.go` with the exact ownership map and body-only `CargarSkills` API.
- [x] 1.3 Add focused tests for all eight assets, frontmatter/template sections, every agent loader path, and unknown-agent behavior.
- [x] 1.4 Run formatting and the required focused/full/race/build/vet/diff validation with no live LLM.

## Work-unit boundary
Starts at the empty Issue #24 skill directory and ends at the eight validated skill assets plus the standalone loader/test package. It does not modify providers, prompts, factories, infrastructure skills, ADK, API, domain, or composition root.
