# Identifier Generation Specification

## Purpose

Define production behavior for the existing `GeneradorID` port without exposing an identifier encoding as a public contract.

## Requirements

### Requirement: Opaque Concurrent Identifier Generation

`GeneradorID.Nuevo` MUST return a server-owned opaque string that is unique across concurrent calls and no longer than 40 characters, as required by Contract §0. The implementation MAY use ULID internally, but callers SHALL NOT depend on a particular encoding, prefix, timestamp, or sequence. Generation MUST remain safe under concurrent use and MUST NOT require a caller-supplied ID.

#### Scenario: Concurrent uniqueness
- GIVEN many simultaneous calls to `Nuevo`
- WHEN all calls complete
- THEN every returned ID is nonempty, unique, and no longer than 40 characters

#### Scenario: Opaque boundary
- GIVEN a generated ID is stored through an existing repository port
- WHEN a caller later reads that entity
- THEN the same ID is preserved as a string
- AND no format-specific behavior is required by the port

## Bounded Slice Acceptance Evidence

The ID work unit MUST include focused concurrent uniqueness evidence and a stated rollback boundary. It SHOULD remain within the 400 authored-line review budget and MUST NOT alter the `GeneradorID` port.