# Delta for demo-control

## MODIFIED Requirements

### Requirement: Confirmed config-gated reset
`POST /api/demo/reiniciar` MUST mutate only when `DEMO_SEED=true` (default false); otherwise it MUST not mutate and return generic 500 `ERROR_INTERNO` without exposing configuration. When enabled, it MUST restore the approved date and synthetic `ana`, `carlos`, `luisa` seed; delete only `fichas`, `hitos`, `planes`, `mensajes`, `leads`; preserve buyer data; perform no DDL; finish under three seconds; and repeatedly return identical 200 `{reiniciado:true,fecha_simulada}`. It MUST NOT enter apply without maintainer confirmation of deletion scope and enabled configuration.

(Previously: Reset restored an approved seed/date without the complete controlled seed or default-off rule.)

#### Scenario: Repeat
- GIVEN confirmation and enabled gate
- WHEN reset is posted twice
- THEN each restores date and three leads within three seconds
- AND no buyer or schema changes.

#### Scenario: Disabled
- GIVEN the gate is unset or false
- WHEN reset is posted
- THEN data is unchanged and generic 500 `ERROR_INTERNO` returns.
