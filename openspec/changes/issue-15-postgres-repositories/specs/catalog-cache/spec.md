# Catalog Cache Specification

## Purpose

Define the existing `CatalogoRepository` port as a file-backed, read-only in-memory catalog for Contract v1.1 data inputs.

## Requirements

### Requirement: Eager File-Backed Catalog

The repository MUST load and validate projects, buyers, affiliate records, and project brochures from the Contract §4 file inputs before serving catalog requests. It MUST retain a successful in-memory snapshot and MUST NOT require PostgreSQL or repeat disk reads for subsequent project, buyer, affiliate, or brochure queries. Initialization failure MUST be returned deterministically and the repository SHALL NOT serve a partial catalog.

#### Scenario: Valid eager snapshot
- GIVEN all required data files conform to their Contract §4 shapes
- WHEN the catalog repository is initialized
- THEN all catalog query methods are available from one validated snapshot
- AND repeated reads perform no additional file or database access

#### Scenario: Invalid source material
- GIVEN a required file is absent, malformed, or violates required identity relationships
- WHEN the catalog repository is initialized
- THEN initialization returns an error
- AND no partial catalog is available

### Requirement: Defensive and Explicit Catalog Results

The repository MUST return defensive results for mutable project maps, buyer collections, and affiliate records so a caller cannot modify the cached snapshot. `AfiliadoPorCedula` and `BrochureMarkdown` MUST return the existing not-found vocabulary for an unmatched identity; a successful brochure lookup SHALL return the source markdown for the requested project.

#### Scenario: Caller mutation isolation
- GIVEN a loaded catalog
- WHEN a caller mutates returned project, buyer, or affiliate data
- THEN a later query returns the original cached values

#### Scenario: Missing affiliate or brochure
- GIVEN a loaded catalog without the requested cédula or project brochure
- WHEN the corresponding lookup is made
- THEN it returns a not-found error

## Bounded Slice Acceptance Evidence

Catalog delivery MUST provide focused evidence for eager validation, no-repeat I/O, defensive-result isolation, and both lookup-miss paths. The work unit MUST state its rollback boundary and SHOULD remain within the 400 authored-line review budget.