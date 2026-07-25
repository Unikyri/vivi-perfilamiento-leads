# Recommendations Routing Specification

## Purpose

Specify recommendations and routing.

## Requirements

### Requirement: Integer-Safe Candidate Eligibility

The system MUST evaluate each catalog project independently using only integer `PresupuestoMax` and positive integer `PrecioDesde`. A project is eligible only when `PresupuestoMax >= PrecioDesde - PrecioDesde/5`, the inclusive `ceil(0.8 × PrecioDesde)` boundary. It MUST NOT use floating-point division, a stale `Capacidad.Ratio`, an invalid price, or a negative budget to admit a project.

#### Scenario: Inclusive affordability boundary
- GIVEN a catalog project with positive `PrecioDesde`
- WHEN the budget equals `PrecioDesde - PrecioDesde/5`
- THEN the project is eligible
- AND a budget one peso lower is ineligible

#### Scenario: Invalid monetary data
- GIVEN a negative budget or a catalog price less than or equal to zero
- WHEN eligibility is evaluated
- THEN no recommendation card is produced

### Requirement: Catalog-Backed Recommendation Aggregation

The system MUST aggregate neighbor count and desistimiento evidence by exact canonical `ProyectoID`. It SHALL create cards only for catalog IDs with at least one contributing neighbor, preserve catalog identity and fields, and ignore empty or non-catalog neighbor IDs. Empty inputs, `K <= 0`, or no qualifying evidence MUST yield a non-nil empty result. Display names MUST NOT be used as identities.

#### Scenario: Orphan neighbor evidence
- GIVEN neighbors for a catalog project and an ID absent from the catalog
- WHEN recommendations are built
- THEN only the catalog project produces a card
- AND its evidence excludes the orphan ID

### Requirement: Deterministic Ranked Recommendations

The system MUST rank eligible cards by neighbor count descending, desistimiento rate ascending, then canonical `ProyectoID` ascending, using exact deterministic comparison. It MUST NOT rely on map or input order, SHALL return at most three cards, and each card MUST expose its contributing neighbor count and desistimiento rate.

#### Scenario: Stable tie and cap
- GIVEN more than three eligible catalog projects with tied counts and rates
- WHEN recommendations are produced from differently ordered inputs
- THEN both results contain the same first three IDs in ascending ID tie order
- AND no result contains more than three cards

### Requirement: Pure 2x2 Route Selection

The system MUST return only `domain.Ruta` from caller-supplied ratio, affiliation, conversion availability, and intention. High intention SHALL mean `Nivel == ALTA` only; `MEDIA`, `BAJA`, empty, and unknown levels are low intention. High capacity SHALL require a finite ratio `>= 0.95`; a non-affiliate requires ratio `>= 1.05` for `ASESOR`. Invalid or negative ratios MUST NOT yield `ASESOR`.

#### Scenario: Thresholds and quadrants
- GIVEN an affiliate with ratio exactly `0.95` and `ALTA` intention
- WHEN the route is selected
- THEN it is `ASESOR`
- AND high-capacity/non-`ALTA`, low-capacity/`ALTA`, and low-capacity/non-`ALTA` select `REMARKETING`, `NUTRICION`, and `DESPEDIDA` respectively

#### Scenario: Non-affiliate threshold
- GIVEN a non-affiliate with `ALTA` intention and no conversion route
- WHEN ratio is exactly `1.05`
- THEN the route is `ASESOR`
- AND ratio from `0.95` inclusive to below `1.05` routes to `NUTRICION`

### Requirement: Conversion Preference and Isolation

The system MUST prefer `NUTRICION` for a non-affiliate with `ALTA` intention and an available conversion route, even when ratio is at least `1.05`. The matrix MUST NOT mutate `Lead` or `Ficha`, create milestones/plans, set cupo fields, or recompute capacity.

#### Scenario: Conversion before advisor allocation
- GIVEN a non-affiliate with ratio `1.10`, `ALTA` intention, and a conversion route
- WHEN the route is selected
- THEN it is `NUTRICION`
- AND no side effect is produced

### Requirement: Dependency and Scope Precondition

Implementation and verification MUST remain blocked until Issue #8 capacity is merged to `main`. This capability MUST NOT modify Issue #8 capacity files or tests, shared domain types, Docs, data, pipeline, HTTP, or use-case layers.

#### Scenario: Unmet capacity dependency
- GIVEN Issue #8 is not merged to `main`
- WHEN apply or verification is requested
- THEN the request is blocked
- AND planning artifacts remain valid without implementation changes
