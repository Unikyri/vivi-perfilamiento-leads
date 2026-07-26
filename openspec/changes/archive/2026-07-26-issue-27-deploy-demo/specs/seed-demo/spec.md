# seed-demo Specification

## Requirements

### Requirement: Explicit demo-seed gate
`DEMO_SEED` MUST default to `false`. The service MUST NOT load demo leads unless an operator manually enables it for the controlled demo.

#### Scenario: Default startup
- GIVEN the value is unset or false
- WHEN the service starts
- THEN no demo lead changes.

#### Scenario: Manual enablement
- GIVEN the controlled demo sets it true
- WHEN data initializes
- THEN seed loading is permitted.

### Requirement: Idempotent synthetic seed
When enabled, the system MUST converge within three seconds to singular, usable synthetic leads `ana`, `carlos`, and `luisa`.

#### Scenario: Repeat
- GIVEN seed loading is enabled
- WHEN it runs twice
- THEN the same three usable leads remain without duplicates.
