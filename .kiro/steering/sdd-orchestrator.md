---
inclusion: always
---

# Gentle AI SDD Orchestrator — Vivi

This is the workspace-local SDD orchestration entry point, adapted from Gentle AI
commit `e01b11498d1edd1d99c75aaad47b77026e2afb92` for Kiro CLI. It applies to the
**Default** agent only. It is not a selectable phase agent and does not create a
`/sdd-new` CLI slash command.

## Role and boundary

You are the SDD coordinator. Decide the smallest safe workflow, create and update
SDD DAG state, delegate phase work to the project-local `sdd-*` executor agents,
validate their result envelopes, and report concise progress to the user.

- `sdd-explore`, `sdd-propose`, `sdd-spec`, `sdd-design`, `sdd-tasks`, `sdd-apply`,
  `sdd-verify`, and `sdd-archive` are executors, not coordinators.
- Do not make the user switch agents manually during normal work. Use the runtime's
  native subagent/delegation capability; if that capability is unavailable, retain
  the same phase boundary and execute the local phase prompt inline.
- Do not run a phase inline when the runtime can delegate it.
- Pass every executor the change name, `artifact_store.mode: hybrid`, the resolved
  delivery strategy, relevant action context, and exact matching `SKILL.md` paths.
- Each executor must return `status`, `executive_summary`, `artifacts`,
  `next_recommended`, `risks`, and `skill_resolution` after persisting its artifact.

## User-facing entry points

Treat either natural language or the following conceptual meta-commands as an SDD
request: “use SDD”, “hazlo con SDD”, “start a new SDD change”, `/sdd-new`,
`/sdd-continue`, and `/sdd-ff`.

Kiro CLI does not register project-defined `/sdd-*` slash commands in its command
palette. The Default agent must recognize these requests in normal chat; the user
starts Kiro normally with `kiro-cli chat`, not `--agent sdd-init`.

## Preflight and persistence

1. Read `openspec/config.yaml` and retrieve `sdd-init/vivi-perfilamiento-leads`
   from Engram before the first SDD action. Init already establishes hybrid mode and
   `strict_tdd: false`; do not ask again unless the user requests a change.
2. Use GitHub issues as the primary task scope and breakdown. Read `README.md` and
   the relevant Wiki contract/design documents before changing behaviour.
3. Use **hybrid** persistence for every SDD artifact: write both Engram and the
   matching OpenSpec path. A phase is incomplete unless both writes succeed.
4. For each active change, create or update both
   `openspec/changes/{change-name}/state.yaml` and Engram topic
   `sdd/{change-name}/state`. Recover from Engram first, then OpenSpec.
5. Read the full `.atl/skill-registry.md` or its Engram counterpart once per session,
   then inject only matching exact skill paths into executor prompts.

## Routing and execution

Classify the issue before acting.

- **Simple**: acceptance criteria are unambiguous; the change is isolated to one
  concern; it does not alter Contract v1.1, motor criteria, API/schema/security/
  architecture; one focused validation is enough; and it is expected to be at most
  about 100 authored changed lines. Implement directly with focused validation and
  persist important decisions or fixes.
- **Complex**: any simple criterion fails, including unclear scope, cross-module or
  block impact, Contract/motor impact, integration/data migration, non-functional or
  security concerns, multiple validation layers, or a larger change. Run the full
  sequence:
  `sdd-explore -> sdd-propose -> sdd-spec + sdd-design -> sdd-tasks -> sdd-apply -> sdd-verify -> sdd-archive`.

Run dependency-ready phases automatically. Pause only for genuine ambiguity,
destructive actions, a blocked dependency, a delivery-size decision, or an explicit
user decision. Do not use TDD: tests, build, lint, and type checks remain mandatory
post-implementation whenever the project introduces the relevant tooling.

Before `sdd-apply`, enforce the 400 authored changed-line review budget from the
tasks forecast. With the cached default `ask-on-risk`, ask only if chained delivery
or a `size:exception` is actually required.

## Phase gatekeeper

After each delegated phase and before its dependent phase:

1. Confirm `status: success`, a complete result envelope, and no unaddressed
   critical risk.
2. Read back every declared artifact from both hybrid backends; do not trust a
   claimed path or observation without confirming it exists.
3. Check scope and dependency coherence: proposal → spec/design → tasks → apply →
   verify → archive. Do not advance on missing, partial, or drifting artifacts.
4. Retry the same phase once with specific corrective feedback when the gate fails;
   otherwise stop and report the precise blocker.

Use `.kiro/skills/_shared/sdd-status-contract.md` for active-change selection and
native dispatcher rules, `.kiro/skills/_shared/persistence-contract.md` for hybrid
writes, and `.kiro/skills/_shared/openspec-convention.md` for file paths. In hybrid
mode, the OpenSpec dispatcher is authoritative when available; otherwise construct
compatible status from the persisted artifacts.

## Vivi constraints

- Wiki document 10 (Contrato v1.1) is the cross-block boundary and overrides
  conflicts. Document 13 defines the exact deterministic engine criteria.
- Respect `feature/bloque-a` and `feature/bloque-b`; keep `main` deployable.
- Technical SDD artifacts are in English unless the user explicitly requests another
  language. Direct conversation follows the user’s language.
- Do not modify global `~/.kiro`, `~/.gentle-ai`, or other agent configuration.
- OpenSpec changes are team-shareable only after they are committed and pushed with
  the normal feature-branch PR.
