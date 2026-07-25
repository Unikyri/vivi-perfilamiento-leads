# Usecase Ports Specification

## Purpose

Define application-owned usecase ports and test-double semantics without changing domain or infrastructure contracts.

## Requirements

### Requirement: Clean, Compatible Port Boundary

`internal/usecase` MUST own context-aware ports and import only the standard library and `internal/domain`; it MUST NOT import HTTP, SQL, ADK, provider SDKs, adapters, or infrastructure. `LeadRepository` SHALL retain pointer-entity `Crear`, `PorID`, `Guardar`, and `Listar` operations using `FiltroLeads`. `PlanRepository`, `FichaRepository`, and `CatalogoRepository` MUST be declared with pointer entities. `LLMProvider` MUST expose `Nombre`, text/audio turns, and Contract §7-aligned DTOs. `MensajeriaGateway`, `Reloj`, `BusEventos`, and `GeneradorID` MUST be declared.

#### Scenario: Application consumer compiles
- GIVEN domain types from Issues #6 and #7
- WHEN a consumer imports `internal/usecase`
- THEN it receives pointer ports and application DTOs
- AND no outer-layer dependency is required

#### Scenario: Provider identity is available
- GIVEN an LLM implementation
- WHEN health or metrics request its identity
- THEN `Nombre` returns it without SDK types

### Requirement: Explicit Repository Results

The package MUST expose `ErrNoEncontrado` and `ErrOptimisticLock`, matchable with `errors.Is`. A missing single record MUST return a nil pointer and an error wrapping `ErrNoEncontrado`; missing ficha MUST remain distinguishable from missing lead. Duplicate lead creation MUST fail without replacing storage.

#### Scenario: Absent lead
- GIVEN no lead exists for an ID
- WHEN `PorID` is invoked
- THEN it returns nil and `ErrNoEncontrado`

#### Scenario: Missing ficha
- GIVEN an existing lead without a ficha
- WHEN its ficha is requested
- THEN the result identifies missing ficha, not missing lead

### Requirement: Lead Fake CAS and Isolation

`FakeLeadRepository.Guardar` MUST compare incoming and stored `Version`. A mismatch SHALL return `ErrOptimisticLock` without mutation. A success MUST return and store a committed defensive copy incremented exactly once; stale writes MUST NOT increment. Every fake store/return boundary MUST recursively clone mutable profiles (including nested maps/slices), capacity, intention, and pointers while preserving concrete values.

#### Scenario: Stale save
- GIVEN stored version 2 and incoming version 1
- WHEN `Guardar` runs
- THEN it returns `ErrOptimisticLock`
- AND storage remains version 2 and unchanged

#### Scenario: Boundary mutation
- GIVEN a lead with nested mutable data
- WHEN a caller mutates an input or returned value
- THEN a fresh read retains independent stored data

### Requirement: Deterministic Lead Queue

`FiltroLeads` MUST conjunctively apply non-nil `Afiliado` and `Ruta`. `Listar` MUST return a non-nil empty slice when unmatched and order matches by `Prioridad` descending, then `LeadID` ascending.

#### Scenario: Stable filtered list
- GIVEN leads inserted in arbitrary order
- WHEN `Listar` receives populated filters
- THEN only conjunctive matches appear in the required order

#### Scenario: Empty list
- GIVEN no lead matches the filters
- WHEN `Listar` is called
- THEN it returns a non-nil empty slice

### Requirement: Narrow Fakes and Scoped Delivery

Only the lead fake and minimal LLM, clock, and deterministic-ID doubles MUST be implemented. Plan, ficha, and catalog fakes MUST be deferred; plan CAS MUST NOT be invented. Work SHALL use sequential ports/shape-test then lead-fake/behavior-test slices, each within the 400-line review budget. Implementation/tests MUST stay in `internal/usecase`; no domain, adapters, infrastructure, data, or Docs files may change.

#### Scenario: Deferred ports
- GIVEN declared plan, ficha, and catalog ports
- WHEN no approved usecase scenario uses them
- THEN no corresponding fake or plan versioning is added

#### Scenario: Scope review
- GIVEN the completed two-slice diff
- WHEN its paths are inspected
- THEN only permitted usecase and Issue #11 OpenSpec paths are present
