# Proposal: Lead State Machine (Issue #7)

## Intent

`Lead.Estado` is unguarded: any package can write an impossible state (`NUEVO → CERRADO`, reviving `CERRADO`) silently. Contract v1.1 §1, architecture §5, and NFR-M-03.5 require an explicit State machine whose edges live in one file and whose invalid edges are domain errors. Issues #17, #18, #19, #21 depend on it.

## Scope

### In Scope
- `internal/domain/estado.go`: **private** table `transiciones`, `ErrTransicionInvalida{Desde,Hacia}`, `PuedeTransicionar(desde,hacia) bool`, `(*Lead) Transicionar(hacia) error`, `EstadosPosibles(desde) []EstadoLead` returning a **defensive copy**.
- `internal/domain/estado_test.go`: table-driven tests — 11 valid edges, invalid pairs, terminals, failure atomicity, slice aliasing.

### Out of Scope
- ADK agents or any ADK import (pure domain, NFR-M-01)
- usecase, adapters, infrastructure, frontend, Docs
- HTTP 409 mapping; optimistic-lock `version` (Contract §5, NFR-E-01)
- Any change to `enums.go`, `Lead` fields, or JSON tags

## Capabilities

### New Capabilities
- `lead-state-machine`: authoritative lifecycle relation, guarded mutation, reachability query in `internal/domain`.

### Modified Capabilities
None — `domain-contract-types` keeps its exact type and JSON surface.

## Approach

One unexported adjacency map keyed by `EstadoLead`, so no caller can add/remove edges or race on it. `Transicionar` validates **before** assignment and returns `ErrTransicionInvalida` on rejection, leaving `Estado` untouched. `EstadosPosibles` allocates a fresh slice per call in Contract order.

## Affected Areas

| Area | Impact |
|------|--------|
| `internal/domain/estado.go` | New — table, error, guard, queries |
| `internal/domain/estado_test.go` | New — lifecycle tests |
| `enums.go`, `lead.go` | Unchanged — read-only reuse |

## Risks

| Risk | Likelihood | Mitigation |
|------|---|---|
| Exported table/slice corrupts policy | Low | Private map, copy on return |
| Direct `l.Estado = x` bypasses guard | Med | Doc `Transicionar` as sole mutator |
| Issue says `EstadoNuevo`, repo has `EstadoLeadNuevo` | High | Use #6 identifiers |

## Rollback Plan

`git revert` the merge commit, or delete the two new files. No existing file is modified, so the tree returns to a compiling state with zero consumer impact.

## Dependencies

- Issue #6 (`domain-contract-types`) — merged
- Contract v1.1 §1, architecture §5, NFR-M-03.5 — read-only authorities

## Success Criteria

- [ ] `go build ./...` and `go vet ./...` clean
- [ ] `go test ./internal/domain/...` passes, incl. `-run TestTransicion` and `-run TestEstadosTerminales`
- [ ] All 11 Contract §1 edges accepted; every other pair (self-transitions, edges out of `CERRADO`/`DESPEDIDO`, unknown states) returns `ErrTransicionInvalida` with `Desde`/`Hacia` set
- [ ] A failed transition leaves `Lead.Estado` unchanged
- [ ] Table exists in exactly one file (grep-provable)
- [ ] Mutating the `EstadosPosibles` result does not change a later call
- [ ] New files import only `fmt` (NFR-M-01); no ADK/agent code

## Proposal Question Round

> Autonomous delivery authorized — assumptions applied.

| # | Question | Assumption |
|---|---|---|
| 1 | Error value or pointer? | Value, `errors.As`-compatible |
| 2 | `EstadosPosibles`: backing slice or copy? | Copy — aliasing is a defect |
| 3 | `EstadoNuevo` or `EstadoLeadNuevo`? | Repo names |
| 4 | Hide public `Lead.Estado`? | Keep — Contract JSON needs it |
