# Buyer-Twin kNN Specification

## Purpose
Define deterministic buyer selection; project ranking is excluded.

## Requirements

### Requirement: Contract-Safe Twin Input
The system MUST accept one request with lead profile, explicit affiliation, buyer records, `CatalogZones` keyed by canonical `ProyectoID`, and `K`. It MUST NOT accept parallel slices, buyer zones, or other zone sources. A missing/non-canonical key SHALL omit zone.

#### Scenario: Catalog zone is keyed safely
- GIVEN a buyer with a matching canonical map key
- WHEN its feature is projected
- THEN zone is the normalized catalog value

#### Scenario: Missing catalog key omits zone
- GIVEN no matching catalog-zone key
- WHEN the buyer is projected
- THEN zone is absent

### Requirement: Lead and Age Projection
The system MUST map `ingreso_hogar` for kNN only: `<=2 SMMLV` to A, `>2` through `<=4` to B, and `>4` to C; missing income SHALL omit category. Historical brackets MUST map `20-35=27.5`, `36-45=40.5`, `46-55=50.5`, and `55+=60.0`; empty/`SIN_DATO` SHALL omit age.

#### Scenario: Lead category boundaries
- GIVEN leads at 2, above 2 through 4, and above 4 SMMLV
- WHEN category is projected
- THEN their categories are A, B, and C respectively

#### Scenario: Open age bracket is stable
- GIVEN `rango_edad` is `55+`
- WHEN age is projected
- THEN its value is exactly `60.0`

### Requirement: Weighted Gower Distance
The system MUST calculate `sum(w*d)/sum(w)` over dimensions present in both points. Weights SHALL be category `.35`, zone `.20`, age `.15`, affiliation `.15`, dependents `.15`; categorical/binary distance is 0 equal, 1 otherwise; age range `32.5`, dependents `10`. Missing dimensions MUST be omitted and weights renormalized; zero dependents is present.

#### Scenario: Missing feature renormalizes
- GIVEN category is absent in one point and others are present
- WHEN distance is calculated
- THEN category is excluded from numerator and denominator
- AND remaining weights are renormalized

#### Scenario: Identical features are symmetric
- GIVEN equal mutually present features
- WHEN distance is calculated in either direction
- THEN both distances are zero and equal

### Requirement: Project-Local Missing Data
The system MUST impute missing historical category with deterministic mode and age with numerical median of bracket representatives from same-project affiliated buyers, only with at least five affiliates. Otherwise it SHALL omit each missing dimension. It MUST NOT use global, cross-project, price, or financial data; imputation is kNN-only.

#### Scenario: Eligible project imputes
- GIVEN a buyer missing category and age in a project with five affiliates
- WHEN features are projected
- THEN same-project affiliated mode and median are used

#### Scenario: Small project falls back to omission
- GIVEN fewer than five project affiliates
- WHEN a required feature is missing
- THEN that dimension is omitted and Gower renormalizes

#### Scenario: Other projects cannot impute
- GIVEN another project has sufficient affiliates
- WHEN this project does not
- THEN no cross-project value is used

### Requirement: Deterministic Neighbor Result
The system MUST return `min(K,n)` neighbors ordered by ascending distance, then ascending buyer `ID`, and SHALL not mutate supplied inputs.

#### Scenario: Default neighbor count
- GIVEN at least 30 buyers and `K=30`
- WHEN neighbors are selected
- THEN exactly 30 are distance-ordered

#### Scenario: Oversized K and ties
- GIVEN fewer buyers than K with equal distances
- WHEN neighbors are selected
- THEN all are returned and ties use ascending ID

### Requirement: Provenance and Issue Boundary
The system MUST NOT use project name, slug, price, or price-derived category as a proxy. This change SHALL NOT modify Issue #8 capacity files or financial capacity, subsidy, or credit behavior.

#### Scenario: Name and price are inert
- GIVEN fixed canonical zones and features but changed name or price
- WHEN distance is recalculated
- THEN it is unchanged

#### Scenario: Capacity is excluded
- GIVEN the Issue #9 change paths
- WHEN inspected
- THEN no Issue #8 capacity path is changed
