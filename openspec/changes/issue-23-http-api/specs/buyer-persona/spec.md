# buyer-persona Specification

## Requirements

### Requirement: Deterministic catalog aggregates
`GET /api/gerencia/buyer-persona?proyecto_id=...` MUST return Contract identity, `muestras`, affiliation, category and age-band proportions, `tasa_desistimiento`, and `actualizado_en`. Without `proyecto_id`, it MUST return `{proyectos:[...]}` with that summary per project. Values MUST derive only from `data/compradores.json`, preserve Contract normalization, be deterministic for unchanged data, and MUST NOT mutate catalog or leads.

#### Scenario: Project
- GIVEN buyers for one project
- WHEN its view is requested
- THEN its counts, proportions, and abandonment rate are returned.

#### Scenario: Catalog
- GIVEN multiple projects
- WHEN no filter is supplied
- THEN every summary returns without source mutation.
