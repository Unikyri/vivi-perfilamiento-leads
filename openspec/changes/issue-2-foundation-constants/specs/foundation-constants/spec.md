# Foundation Constants Specification

## Purpose

Single-source normative economic values (Contract §4.5, Engine Criteria §1) embedded at compile time and exposed as typed domain structs via pure parsing APIs. The authoritative root JSON files in `references/` are the single source of truth; the embed package (`references/embed.go`) only exposes their raw bytes — it adds no logic, no transformation, and no independent values.

## Requirements

### Requirement: Canonical JSON Files

The system MUST provide three JSON files at `references/` with exact normative content:

| File | Payload |
|------|---------|
| `constantes.json` | `{"smmlv_2026":1750905,"tasa_ea_default":0.107,"plazo_meses_default":240,"cuota_ingreso_max":0.40,"tope_vis_smmlv":150,"mediana_vis_cop":195000000}` |
| `subsidios_2026.json` | `[{"ingreso_hasta_smmlv":2,"subsidio_smmlv":30},{"ingreso_hasta_smmlv":4,"subsidio_smmlv":20}]` |
| `calendario_economico.json` | `[{"tipo":"CESANTIAS","fecha":"--02-14"},{"tipo":"PRIMA","fecha":"--06-30"},{"tipo":"PRIMA","fecha":"--12-20"}]` |

No normative numeric value SHALL be hard-coded in any `.go` source file. The JSON files are the sole canonical store.

#### Scenario: JSON content matches normative sources

- GIVEN the three JSON files exist in `references/`
- WHEN their content is read
- THEN each value MUST match Contract §4.5 / Engine Criteria §1 exactly

### Requirement: Legal Root Package Embedding

`references/embed.go` MUST declare `package references` and use `//go:embed` directives to expose each JSON file as an exported `[]byte` variable (`ConstantesJSON`, `SubsidiosJSON`, `CalendarioJSON`). The embed package MUST NOT contain parsing logic, type definitions, or `init()` functions. It SHALL only expose raw bytes.

#### Scenario: Embed compiles without illegal paths

- GIVEN `embed.go` co-located with JSON files in `references/`
- WHEN `go build ./...` is executed
- THEN compilation MUST succeed with no `..` embed patterns

### Requirement: Pure Domain Parsing APIs

`internal/domain/constantes.go` MUST export: `LoadConstantes([]byte) (Constantes, error)`, `LoadSubsidios([]byte) ([]TramoSubsidio, error)`, `LoadCalendario([]byte) ([]EventoCalendario, error)`. Each function accepts raw JSON bytes and returns typed structs.

Struct fields:
- `Constantes`: SMMLV (int), TasaEADefault (float64), PlazoMesesDefault (int), CuotaIngresoMax (float64), TopeVISSMMLV (int), MedianaVISCOP (int64)
- `TramoSubsidio`: IngresoHastaSMMLV (int), SubsidioSMMLV (int)
- `EventoCalendario`: Tipo (string), Fecha (string, `--MM-DD` annual recurrence format)

#### Scenario: Successful parse of valid JSON

- GIVEN valid `constantes.json` bytes
- WHEN `LoadConstantes(bytes)` is called
- THEN it MUST return a `Constantes` with SMMLV=1750905, TasaEADefault=0.107, PlazoMesesDefault=240, CuotaIngresoMax=0.40, TopeVISSMMLV=150, MedianaVISCOP=195000000
- AND error MUST be nil

#### Scenario: Malformed JSON returns error

- GIVEN bytes containing invalid JSON (e.g., truncated or non-JSON)
- WHEN any `Load*()` function is called
- THEN it MUST return a non-nil error wrapping the parse failure
- AND MUST NOT panic

#### Scenario: Calendar parsing

- GIVEN valid `calendario_economico.json` bytes
- WHEN `LoadCalendario(bytes)` is called
- THEN it MUST return 3 `EventoCalendario` entries: CESANTIAS/--02-14, PRIMA/--06-30, PRIMA/--12-20

### Requirement: Clean Architecture Import Constraints

The domain package (`internal/domain/`) MUST NOT import `references/` or any package outside the Go standard library. Only composition roots (`cmd/`) MAY import both `references/` and `internal/domain/`.

#### Scenario: Domain has no outer imports

- GIVEN the source of `internal/domain/constantes.go`
- WHEN its import list is inspected
- THEN it MUST contain only `encoding/json` (and optionally `fmt` or `errors`)

### Requirement: Composition Root Wiring Deferred

This change MUST NOT modify `cmd/servidor/main.go` or `cmd/pipeline/main.go`. Wiring is deferred to Issue #4.

#### Scenario: cmd packages unchanged

- GIVEN the repository state after this change
- WHEN `cmd/` directory is inspected
- THEN no files SHALL have been added or modified
