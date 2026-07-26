# Bus Events Specification

## Purpose
Define the provider-free, in-memory delivery contract for internal lead events.

## Requirements

### Requirement: F1 In-Process Boundary
The bus MUST deliver events in-process and MUST NOT require an LLM, ADK, HTTP, database, queue, retry service, or composition-root wiring.

#### Scenario: Standalone delivery
- GIVEN a constructed bus and one subscriber
- WHEN an event is published
- THEN the subscriber receives it without an external dependency

### Requirement: F3 Synchronous Ordered Dispatch
`Publicar` MUST complete each matching callback synchronously in subscription order. Unknown types and zero subscribers MUST succeed silently. A callback MAY publish or subscribe without blocking delivery.

#### Scenario: Ordered nested publication
- GIVEN two subscribers in registration order
- WHEN the first publishes another event
- THEN both original callbacks run in order
- AND the nested publication completes without deadlock

### Requirement: F4 Isolated Event Delivery
The bus MUST preserve the publisher's event snapshot and MUST give every callback an independent deep copy, including nested maps and slices. It MUST accept a nil payload; nil callbacks SHALL be ignored.

#### Scenario: Mutation isolation
- GIVEN two callbacks and a payload containing a nested slice
- WHEN the first callback mutates its event
- THEN the second receives the original nested values
- AND the publisher's later mutation is not observed

### Requirement: F5 Panic Containment
The bus MUST recover a panic from each callback, MUST continue with later callbacks, and MUST NOT propagate callback failure to the publisher.

#### Scenario: Later callback survives
- GIVEN a panicking callback followed by a recording callback
- WHEN an event is published
- THEN publication returns normally
- AND the recording callback receives the event

### Requirement: F6 Privacy-Safe Observability
The bus MUST log only event type, opaque lead ID, handler identity, and a fixed outcome category. It MUST NOT log payloads, message text, names, cedulas, or raw panic/error values.

#### Scenario: Safe panic record
- GIVEN a callback that panics with sensitive text
- WHEN an event is published
- THEN one fixed failure outcome is observable
- AND the sensitive text is absent from logs
