# PostgreSQL Repositories Specification

## Purpose

Define production persistence behavior for the existing lead, plan, and ficha repository ports. Public ports, domain types, schema, migrations, and API behavior remain unchanged.

## Requirements

### Requirement: Lead Repository Consistency and Determinism

The repository MUST create and retrieve leads and messages through the existing `LeadRepository` port. A lead save MUST be atomic against the supplied version: a missing lead MUST return the existing not-found vocabulary, and an existing lead with a nonmatching version MUST return `ErrOptimisticLock`; neither outcome SHALL mutate persisted state. A successful save MUST advance the version. Lead filters MUST be conjunctive, lists MUST order by descending priority then ascending lead ID, and conversations MUST order chronologically.

#### Scenario: Successful compare-and-swap
- GIVEN a persisted lead at version 3
- WHEN `Guardar` receives its intended version 3
- THEN the lead is persisted at the next version
- AND the observable replacement is atomic

#### Scenario: Missing or stale compare-and-swap
- GIVEN a missing lead or a persisted lead at version 4
- WHEN `Guardar` receives a version-3 lead
- THEN the missing case returns a not-found error and the stale case returns `ErrOptimisticLock`
- AND neither case changes stored data

#### Scenario: Deterministic queue and conversation
- GIVEN leads matching optional affiliate and route filters and messages with distinct timestamps
- WHEN callers list leads and read a conversation
- THEN lead filters are applied together and results order by priority descending then ID ascending
- AND messages order by creation time ascending

### Requirement: Transactional Nutrition Plan Aggregate

The repository MUST create, retrieve, and save the existing `PlanRepository` aggregate consistently. A plan save MUST atomically persist the plan and supplied milestones; it MUST NOT introduce compare-and-swap behavior because the port has no version. Omitted milestones SHALL NOT be inferred as deletions. `PorLead` MUST reconstruct persisted milestones; `HitosVencidos` MUST return pending milestones of active plans due on or before the requested date; `MarcarHito` MUST report a missing milestone with the existing not-found vocabulary.

#### Scenario: Atomic aggregate save
- GIVEN a plan with milestones and a persistence failure during save
- WHEN `Guardar` returns an error
- THEN no partial plan or milestone replacement is observable

#### Scenario: Due milestone selection and absence
- GIVEN active and inactive plans with pending and non-pending milestones
- WHEN overdue milestones are requested or an unknown milestone is marked
- THEN only pending milestones due on or before the date are returned
- AND the unknown milestone returns not-found

### Requirement: Unique Ficha per Lead

The repository MUST save fichas through the existing `FichaRepository` port with at most one ficha for each lead. Saving a ficha for an existing lead MUST replace its stored ficha rather than create a duplicate. `PorLead` MUST discriminate a missing lead from an existing lead whose ficha has not been generated, using the existing not-found error vocabulary and resource identity.

#### Scenario: Ficha upsert
- GIVEN an existing lead with a stored ficha
- WHEN a replacement ficha for that lead is saved
- THEN one ficha remains for the lead
- AND retrieval returns the replacement

#### Scenario: Discriminated ficha absence
- GIVEN one unknown lead and one existing lead without a ficha
- WHEN each ficha is requested
- THEN both failures are not-found errors
- AND their resource identities distinguish lead absence from ficha absence

## Bounded Slice Acceptance Evidence

Implementation SHALL be delivered without changing public ports or schema. Slice 1 evidence MUST cover lead CAS, ordering, conversation, and ID behavior; Slice 2 MUST cover atomic plan/ficha persistence and absence handling; Slice 3 MUST cover catalog caching. Each slice MUST include focused automated results, an integration/runtime result or justified N/A, and a rollback boundary; no slice SHOULD exceed the 400 authored-line review budget.