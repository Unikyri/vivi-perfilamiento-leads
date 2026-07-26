# Perfilamiento Lead Specification

## Purpose

Define Issue #18 pre-profiling (UC-01/UC-03, US-01/US-02) without Contract, motor, port, data, or adapter changes.

## Requirements

### Requirement: Verified affiliate pre-profile

An active affiliate match MUST initialize a non-nil profile and set `Afiliado=true`. It MUST record `ingreso_hogar`, `categoria`, `segmento`, `personas_hogar`, and `tipo_hogar` as `VERIFICADO_BASE`; `personas_hogar` SHALL equal `personas_a_cargo + 1`. For demo only, it MUST record `tiene_vivienda=false` and `recibio_subsidio=false` as `VERIFICADO_BASE`. An empty input name MUST use the affiliate name.

#### Scenario: Ana receives a verified demo pre-profile

- GIVEN active affiliate Ana (`1032456789`) with income `2600000` and two dependants
- WHEN a new lead is pre-profiled with an empty name
- THEN the profile contains verified fields, `personas_hogar=3`, and name `Ana`
- AND capacity has subsidy `52527150` under the existing motor

### Requirement: Non-affiliate fallback

A missing/inactive affiliate MUST yield an empty, non-nil profile and `Afiliado=false`. It MUST NOT stamp verified or demo-eligibility fields; capacity SHALL have zero subsidy.

#### Scenario: Unknown or inactive affiliate remains discoverable

- GIVEN a lead cedula is absent or inactive
- WHEN the lead is pre-profiled
- THEN it remains a non-affiliate with an empty profile and zero subsidy
- AND profiling continues without a lookup failure response

### Requirement: Valid creation, persistence, and event order

The system MUST create the lead in `NUEVO`, calculate capacity through the existing deterministic motor with its non-positive-candidate median fallback, transition through the domain lifecycle to `PERFILANDO`, and persist it. It MUST publish exactly one minimal `LeadNuevo` event only after successful creation. It MUST NOT query project catalogues, select ratio candidates, or produce recommendations in this use case.

#### Scenario: Successful pre-profile is durable before notification

- GIVEN valid new-lead input and a successful repository create
- WHEN pre-profiling completes
- THEN the persisted lead is `PERFILANDO` with its calculated capacity
- AND one `LeadNuevo` event is published after persistence succeeds

#### Scenario: Create or transition failure has no success event

- GIVEN lead creation or the required state transition fails
- WHEN pre-profiling is attempted
- THEN the operation returns the failure
- AND no `LeadNuevo` event is published

### Requirement: Verified household re-consultation

For an active family affiliate, the system MUST mark `hogar_con_afiliado=true` and `cedula_familiar_afiliado` as `VERIFICADO_BASE`, add the family income to `ingreso_hogar`, recalculate capacity, and persist through the existing save/CAS boundary. Re-consulting the same already-verified family cedula MUST be idempotent and MUST NOT add its income again.

#### Scenario: Family match is summed once

- GIVEN a persisted lead and active family affiliate `1015789456` with income `1900000`
- WHEN the family cedula is re-consulted twice
- THEN the first result has verified household affiliation and adds `1900000` once
- AND the second result leaves household income and capacity unchanged

### Requirement: Unverified family and failure safety

For a missing or inactive family cedula, the system MUST persist that cedula as `DECLARADO` with `requiere_confirmacion=true` and return `ErrFamiliarNoEncontrado`; it MUST NOT add income or mark household affiliation verified. A load, save/CAS, or context failure MUST return an error and MUST NOT publish a success event.

#### Scenario: Unknown family requires later confirmation

- GIVEN a persisted lead and an unknown family cedula
- WHEN the cedula is re-consulted
- THEN the declared confirmation-required record is saved and `ErrFamiliarNoEncontrado` is returned
- AND no income, verified household flag, or event is added

### Requirement: Clean deterministic boundary

The use case MUST depend only on existing domain behavior and application ports. It MUST NOT call a real database, LLM, HTTP service, ADK graph, messaging service, wall clock, kNN, or recommendation service; focused tests SHALL use deterministic in-memory fakes.

#### Scenario: Isolated use-case execution

- GIVEN deterministic port fakes and a fixed clock/ID source
- WHEN all pre-profile and re-consultation scenarios run
- THEN results are reproducible without external services
- AND the capacity figures originate only from the existing motor
