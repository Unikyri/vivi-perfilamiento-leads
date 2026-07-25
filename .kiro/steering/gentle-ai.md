---
inclusion: always
---

# Gentle AI SDD — Vivi (workspace-local Kiro CLI adaptation)

This project uses its existing workspace-local SDD adaptation plus the checksum-verified official
Gentle AI `v1.46.0` executable at `./.tools/gentle-ai` (upstream source commit
`b22a7eb8730e0e255c7a6d142aedfc606cbb020e`). The project-local phase agents are JSON configs
under `.kiro/agents/` and their prompts are under `.kiro/prompts/gentle-ai/`. The corresponding
source-of-truth skills are under `.kiro/skills/`.

For native dispatcher inspections, invoke `./.tools/gentle-ai`, never the global executable. The
stable v1.46.0 CLI provides `sdd-status` and `sdd-continue`; it does not provide
`sdd-verify-validate`, `sdd-attempt`, or `review`. Do not fabricate the unavailable admission or
review-receipt lifecycle artifacts, and do not treat this version upgrade as authority to archive
or close a gated SDD change.

## Existing integrations

- Engram is already owned by this project through `.kiro/settings/mcp.json` and the `engram`
  binary. Do not install, replace, or duplicate Engram.
- Context7 is already configured in the same MCP file. Do not add a second Context7 server.
- Do not modify global `~/.kiro`, `~/.gentle-ai`, Claude, Pi, OpenCode, or other agent configs
  for work in this repository.

## SDD routing

Use the smallest workflow that fits the change:

- Small, known change: work directly with focused validation.
- Unclear or multi-file change: `sdd-explore` → `sdd-propose` → `sdd-spec` → `sdd-design`
  → `sdd-tasks` → `sdd-apply` → `sdd-verify` → `sdd-archive`.
- Use `sdd-onboard` for a guided first walkthrough and `sdd-init` before the first SDD change
  when project testing/persistence context has not been initialized.

The Default Kiro CLI agent is the SDD orchestrator: it loads this steering file and
`.kiro/steering/sdd-orchestrator.md`, then delegates phase work to the local
`sdd-*` executor agents. Start normal work with `kiro-cli chat` and say, for example,
“use SDD for issue #42”; do not normally switch agents manually. Selecting a phase
agent with `/agent sdd-explore` or starting `kiro-cli chat --agent sdd-apply` is
manual recovery/debugging only.

Phase agents are executors, not orchestrators: they must follow their local `SKILL.md`, use
Engram through the existing MCP when persistence is needed, and avoid spawning more agents.
The Default orchestrator runs dependency-ready SDD phases automatically. It pauses only for
ambiguity, destructive actions, blocked dependencies, review-size decisions, or decisions
that require explicit user input. It always runs verification before archive.

## Vivi project constraints

- Read `README.md` and the relevant `Docs/vivi-perfilamiento-leads.wiki/` contract/design docs
  before changing behavior.
- Doc 10 (Contrato v1.1) is the boundary between blocks; doc 13 fixes exact motor criteria.
- Respect the two-branch flow (`feature/bloque-a` / `feature/bloque-b`) and keep `main`
  deployable.
- Backend work is Go 1.24+; use the project-local `go-testing` skill for tests.
- Keep technical SDD artifacts in English unless the user explicitly requests another language.
