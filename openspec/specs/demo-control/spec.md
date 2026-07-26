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

## ADDED Requirements

### Requirement: Simulated time owns only demo decisions

Simulated time MUST govern `AvanzarDemo`, `TickReloj`, milestone due-date comparisons, reset, and `/salud.fecha_simulada`. It MUST NOT supply timestamps for newly created or updated leads, messages, consent records, fichas, or accepted HTTP turns; those records MUST use `Reloj.Ahora()` wall time.

#### Scenario: Advanced demo with operational write
- GIVEN simulated time has advanced beyond the current date
- WHEN a new operational record is written
- THEN its timestamp is current UTC wall time
- AND the demo date remains available through `FechaSimulada()`

### Requirement: Buyer-persona boundary is preserved

The buyer-persona filtered and catalog-wide responses MUST remain deterministic and derived only from `data/compradores.json`. This change MUST NOT delete or rename `compradores`, alter its migration, or add SQL aggregation or a dashboard view.

#### Scenario: Existing buyer-persona requests
- GIVEN a valid `proyecto_id`, no filter, and an invalid filter
- WHEN buyer-persona endpoints are requested
- THEN their existing response shapes remain deterministic and the invalid filter returns `VALIDACION`
