# demo-control Specification

## Requirements

### Requirement: Simulated time
`POST /api/demo/tiempo` MUST accept exactly one of `avanzar_hasta` or `avanzar_dias`; neither or both MUST return 400 `VALIDACION`. A valid request MUST persist the date, emit `TickReloj`, execute due milestones, and return 200 `{fecha_simulada,hitos_disparados}`. `GET /salud` MUST expose that persisted date, not wall-clock time.

#### Scenario: Advance time
- GIVEN a simulated date and due active-plan milestones
- WHEN time advances validly
- THEN it returns the new date and exact dispatched count.

### Requirement: Confirmed config-gated reset
`POST /api/demo/reiniciar` MUST mutate only with `DEMO_SEED=true`; otherwise it MUST not mutate and return generic 500 `ERROR_INTERNO` without exposing configuration. When enabled, it MUST restore approved seed/date; delete only `fichas`, `hitos`, `planes`, `mensajes`, and `leads`; preserve buyer data; perform no DDL; finish under three seconds; and repeatedly return 200 `{reiniciado:true,fecha_simulada}` identically. This destructive final slice MUST NOT enter apply absent explicit maintainer confirmation of deletion scope and enabled configuration.

#### Scenario: Confirmed repeat
- GIVEN confirmation and an enabled gate
- WHEN reset is posted twice
- THEN both restore the same seed within three seconds
- AND no buyer record or schema changes.
