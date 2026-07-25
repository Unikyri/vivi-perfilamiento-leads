## Exploration: issue-7-lead-state-machine

### Current State
The repository has the Contract v1.1 typed `EstadoLead` enum and `Lead.Estado`, but no transition table, domain transition error, mutation guard, or `EstadosPosibles` symbol. Issue #7 is scoped to the pure domain: an explicit lifecycle table, domain error/state mutation guard, and domain tests; it excludes ADK agents, use cases, adapters, infrastructure, frontend, and documentation changes. Contract v1.1 §1 is authoritative and defines 11 valid transitions: `NUEVO→PERFILANDO`, `PERFILANDO→CALIFICADO`, `CALIFICADO→ENTREGADO|EN_NUTRICION|REMARKETING|DESPEDIDO`, `EN_NUTRICION→PAUSADO|PERFILANDO`, `PAUSADO→EN_NUTRICION`, `REMARKETING→PERFILANDO`, and `ENTREGADO→CERRADO`. Every other pair, including self-transitions and any transition out of `CERRADO`, must be rejected as `TRANSICION_INVALIDA`.

The safest interpretation is a standard-library-only implementation in `internal/domain/estado.go`. It adds behavior around existing types without changing enum literals, JSON tags, struct layout, or pipeline types. The ADK Go 2.0 policy is intentionally not applicable: this change must not add agents, ADK imports, or outer-layer dependencies.

### Affected Areas
- `internal/domain/enums.go` — preserve `EstadoLead` and its literals.
- `internal/domain/lead.go` — add a guarded method without changing its public field or JSON contract.
- `internal/domain/estado.go` — single home for relation, error, and queries.
- `internal/domain/estado_test.go` — valid/invalid/atomicity/aliasing tests.
- `Docs/` — read-only authority; never modify.

### Approaches
1. **Private adjacency table plus guarded entity method** — private relation with `PuedeTransicionar`, defensive `EstadosPosibles`, and `Lead.Transicionar`.
   - Pros: one source of truth, no external mutation, atomic failure, pure domain.
   - Cons: adds a small public method.
   - Effort: Low.
2. **Exported mutable table** — callers validate/mutate themselves.
   - Pros: minimal helpers.
   - Cons: mutable aliases corrupt policy and race; rejected.
   - Effort: Low, high risk.
3. **Per-state objects** — state-specific objects.
   - Pros: extensible behavior.
   - Cons: over-engineered static relation; rejected.
   - Effort: Medium/High.

### Recommendation
Use approach 1. The private table is the sole relation. `EstadosPosibles` returns a fresh, deterministic slice; mutating it cannot change later results. `Lead.Transicionar` validates before assignment, returns a domain error, and leaves state unchanged on failure. Keep optimistic locking outside domain.

### Risks
- Exported maps or shared slices would allow policy corruption; use private data and copies.
- HTTP 409 mapping stays outside domain.
- The public `Lead.Estado` field permits direct writes for JSON compatibility; document `Transicionar` as the guarded path.
- ADK Go 2.0 is not applicable to this no-agent domain issue.

### Ready for Proposal
Yes — lifecycle edges, error semantics, encapsulation, compatibility, and tests are clear.
