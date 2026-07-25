# Lead State Machine Specification

## Purpose

Define the only permitted lifecycle transitions for `Lead.Estado` in the pure `domain` layer while preserving the existing `EstadoLead` literals and JSON contract.

## Requirements

### Requirement: Contract Lifecycle Relation

The domain MUST accept exactly the following ordered `EstadoLead` edges and MUST reject every other pair.

| From | To |
|---|---|
| `EstadoLeadNuevo` | `EstadoLeadPerfilando` |
| `EstadoLeadPerfilando` | `EstadoLeadCalificado` |
| `EstadoLeadCalificado` | `EstadoLeadEntregado` |
| `EstadoLeadCalificado` | `EstadoLeadEnNutricion` |
| `EstadoLeadCalificado` | `EstadoLeadRemarketing` |
| `EstadoLeadCalificado` | `EstadoLeadDespedido` |
| `EstadoLeadEnNutricion` | `EstadoLeadPausado` |
| `EstadoLeadEnNutricion` | `EstadoLeadPerfilando` |
| `EstadoLeadPausado` | `EstadoLeadEnNutricion` |
| `EstadoLeadRemarketing` | `EstadoLeadPerfilando` |
| `EstadoLeadEntregado` | `EstadoLeadCerrado` |

#### Scenario: Accept every Contract edge

- GIVEN a `Lead` in each table source state
- WHEN `Transicionar` receives the corresponding table target
- THEN it MUST return no error
- AND the lead's `Estado` MUST equal that target

### Requirement: Guarded Rejection and Atomicity

`Lead.Transicionar` MUST reject every pair not in the Contract Lifecycle Relation, including self-transitions, pairs from unknown values, and pairs to unknown values. On rejection, it MUST NOT change `Lead.Estado`.

#### Scenario: Reject all non-listed state pairs

- GIVEN a lead for every source and target in the nine Contract states plus `EstadoLead("")` and `EstadoLead("DESCONOCIDO")`, excluding listed edges
- WHEN `Transicionar` is invoked for each pair
- THEN every invocation MUST fail
- AND each lead MUST retain its initial state

#### Scenario: Reject terminal and self transitions

- GIVEN a lead in `EstadoLeadCerrado` or `EstadoLeadDespedido`, or a lead whose target equals its current state
- WHEN `Transicionar` is invoked with any target
- THEN it MUST fail without changing `Estado`

### Requirement: Typed Invalid-Transition Error

Each rejected transition MUST return an `ErrTransicionInvalida` domain error value, usable with `errors.As`, whose `Desde` and `Hacia` fields preserve the exact attempted `EstadoLead` values. The error SHALL represent Contract code `TRANSICION_INVALIDA`; HTTP mapping MUST remain outside the domain.

#### Scenario: Inspect rejection data

- GIVEN a lead in `EstadoLeadNuevo`
- WHEN it transitions to `EstadoLeadCerrado`
- THEN the returned error MUST expose `Desde == EstadoLeadNuevo` and `Hacia == EstadoLeadCerrado`

### Requirement: Single Private Policy and Defensive Query

The transition relation MUST have one private, immutable-in-practice authority. `PuedeTransicionar`, `Transicionar`, and `EstadosPosibles` MUST derive eligibility from it; no mutable policy map or backing slice SHALL be exposed. `EstadosPosibles` MUST return a fresh slice in table order and a zero-length result for terminal or unknown sources.

#### Scenario: Defend possible-state results from aliasing

- GIVEN `EstadosPosibles(EstadoLeadCalificado)` is called twice
- WHEN the caller mutates the first returned slice
- THEN the second result MUST remain `Entregado`, `EnNutricion`, `Remarketing`, `Despedido` in that order

#### Scenario: Query terminal and unknown sources

- GIVEN `EstadoLeadCerrado`, `EstadoLeadDespedido`, or `EstadoLead("DESCONOCIDO")`
- WHEN `EstadosPosibles` and `PuedeTransicionar` are queried
- THEN the possible-state result MUST be zero-length and eligibility MUST be false

### Requirement: Pure Domain Compatibility

This capability MUST remain in `internal/domain`, MUST NOT import ADK, agent, use-case, adapter, infrastructure, or transport packages, and SHOULD use only the Go standard library. It MUST NOT change `EstadoLead` literals, `Lead` fields, or existing JSON tags.

#### Scenario: Preserve domain boundary and contract types

- GIVEN the state-machine package is built with existing domain tests
- WHEN its imports and existing `Lead`/`EstadoLead` serialization are checked
- THEN no forbidden dependency or public type/tag change MUST be present
