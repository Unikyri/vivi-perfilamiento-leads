# Delta for demo-control

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
