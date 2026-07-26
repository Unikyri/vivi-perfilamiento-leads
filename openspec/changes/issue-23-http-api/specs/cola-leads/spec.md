# cola-leads Specification

## Requirements

### Requirement: Deterministic advisor queue
`GET /api/leads` MUST return `{cupo_10:{usados,porcentaje_ventana},leads:[...]}` ordered by persisted `prioridad` descending. It MUST accept `afiliado=true|false` and Contract `ruta`; invalid filters MUST return 400 `VALIDACION`. `usados` MUST count only non-affiliated `ASESOR` leads. `semaforo` and `resumen` MUST derive deterministically from persisted Contract fields; this read MUST NOT recalculate priority or mutate leads.

#### Scenario: Ranking
- GIVEN leads with distinct priorities and affiliations
- WHEN an `ASESOR` queue is requested
- THEN matches are priority-descending
- AND `cupo_10` counts only non-affiliated advisor leads.

#### Scenario: Invalid filter
- GIVEN an invalid filter
- WHEN the queue is requested
- THEN it returns 400 `VALIDACION`.
