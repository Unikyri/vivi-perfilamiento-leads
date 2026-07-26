# Issue #24 — Agent Skills Specification

## Requirement: Skill documents
The change MUST provide exactly the eight Issue #24 skill documents under `skills/{name}/SKILL.md`: `tono-colsubsidio`, `normalizacion-de-declarados`, `siguiente-mejor-pregunta`, `explicacion-financiera-humana`, `dominio-caja`, `redaccion-con-dignidad`, `presentacion-de-proyectos`, and `ficha-comercial`.

Each document MUST use the Issue #24 frontmatter fields `name`, `description`, `agente`, and `fuente_de_verdad`, and MUST contain `Por qué existe`, `Instrucciones`, `Ejemplos` with a limit case, and `Criterios de aceptación`.

## Requirement: Ownership map
The loader MUST preserve the Issue #24 ownership map: asesora loads five skills, investigadora loads `dominio-caja`, nutricionista loads tone/dignity/financial explanation, and documentadora loads `ficha-comercial`.

## Requirement: Skill loading
`CargarSkills` MUST concatenate the selected skill bodies in map order, include a visible skill delimiter, omit YAML frontmatter, and return an empty string without error for an unknown agent. It MUST load repository assets without any live LLM or provider call.

### Scenario: Complete skill set
- **Given** the repository is built
- **When** the agent skill package is tested
- **Then** all eight documents are readable and have valid required metadata and sections.

### Scenario: Body-only prompt injection
- **Given** `CargarSkills("asesora")`
- **When** it returns the concatenated prompt material
- **Then** it contains each owned skill body and no `fuente_de_verdad` or other YAML frontmatter.

### Scenario: Unknown agent
- **Given** an agent name absent from the ownership map
- **When** `CargarSkills` is called
- **Then** it returns an empty string and no error.
