# Design: Foundation Constants

## Technical Approach

Root-level embed package (`references/`) co-locates `//go:embed` with canonical JSON files, exposing raw `[]byte` vars. Pure domain parsers in `internal/domain/constantes.go` accept `[]byte` and return typed structs — no outer imports, no `init()`, full DI via composition root (wiring deferred to Issue #4). This resolves the illegal `../../` embed pattern from the issue sample without duplicating data.

## Architecture Decisions

| Decision | Alternatives | Rationale |
|----------|-------------|-----------|
| Embed package at `references/` | Copy-at-build, internal embed pkg, go:generate | Only legal single-source zero-duplication option; `//go:embed` forbids `..` |
| Explicit `Load*([]byte)` over `init()` global | Package-level singleton set by `init()` | Testable, explicit dependency, no hidden coupling |
| `fmt.Errorf` wrapping for parse errors | Custom error types | Minimal; errors are startup-fatal, not classified at runtime |
| Domain imports only `encoding/json` + `fmt` | None considered | Clean Architecture constraint from `doc.go`; spec requirement |
| JSON struct tags use `snake_case` Spanish | camelCase | Matches Contract §0 field naming convention |

## Data Flow

```
references/*.json ──[go:embed]──→ references/embed.go ([]byte vars)
                                         │
            ┌────────────────────────────┘
            │ (imported by cmd/ in Issue #4)
            ▼
 cmd/servidor or cmd/pipeline (composition root)
            │ passes []byte
            ▼
 domain.LoadConstantes([]byte) → Constantes
 domain.LoadSubsidios([]byte)  → []TramoSubsidio
 domain.LoadCalendario([]byte) → []EventoCalendario
            │ injected into
            ▼
 motor / usecase services (future issues)
```

**Data ownership**: `references/*.json` is the single source of truth owned by both blocks. `references/embed.go` is infrastructure glue (Block A). `internal/domain/constantes.go` is domain-layer owned by Block A.

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `references/constantes.json` | Create | Exact §4.5 values |
| `references/subsidios_2026.json` | Create | Exact §4.5 subsidy table |
| `references/calendario_economico.json` | Create | Exact §4.5 calendar |
| `references/embed.go` | Create | `package references`; 3 `//go:embed` + exported `[]byte` vars |
| `references/.gitkeep` | Delete | Superseded by real files |
| `internal/domain/constantes.go` | Create | Structs + `Load*()` parsers |
| `internal/domain/constantes_test.go` | Create | Table-driven validation |

## Interfaces / Contracts

```go
package references

import _ "embed"

//go:embed constantes.json
var ConstantesJSON []byte

//go:embed subsidios_2026.json
var SubsidiosJSON []byte

//go:embed calendario_economico.json
var CalendarioJSON []byte
```

```go
package domain

import (
    "encoding/json"
    "fmt"
)

type Constantes struct {
    SMMLV2026         int     `json:"smmlv_2026"`
    TasaEADefault     float64 `json:"tasa_ea_default"`
    PlazoMesesDefault int     `json:"plazo_meses_default"`
    CuotaIngresoMax   float64 `json:"cuota_ingreso_max"`
    TopeVISSMMLV      int     `json:"tope_vis_smmlv"`
    MedianaVISCOP     int64   `json:"mediana_vis_cop"`
}

type TramoSubsidio struct {
    IngresoHastaSMMLV int `json:"ingreso_hasta_smmlv"`
    SubsidioSMMLV     int `json:"subsidio_smmlv"`
}

type EventoCalendario struct {
    Tipo  string `json:"tipo"`
    Fecha string `json:"fecha"`
}

func LoadConstantes(data []byte) (Constantes, error)
func LoadSubsidios(data []byte) ([]TramoSubsidio, error)
func LoadCalendario(data []byte) ([]EventoCalendario, error)
```

Each `Load*` returns `fmt.Errorf("domain: parse %s: %w", name, err)` on failure.

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit | `LoadConstantes` returns exact §1 magnitudes | Table-driven: assert every field equals spec value |
| Unit | `LoadSubsidios` returns 2 tramos with exact values | Assert len=2, each `ingreso_hasta_smmlv` and `subsidio_smmlv` |
| Unit | `LoadCalendario` returns 3 events in correct order | Assert types and `--MM-DD` dates |
| Unit | Malformed JSON returns wrapped error, no panic | Pass `[]byte("{{invalid")` to each `Load*` |
| Unit | Empty slice input returns error | Pass `[]byte("")` or `nil` |
| Build | `go build ./...` succeeds | CI validation command |
| Build | Domain has no non-stdlib imports | `go list -f '{{.Imports}}'` check |

**Validation commands:**
```bash
go build ./...
go test ./internal/domain/... -v -count=1
go list -f '{{join .Imports "\n"}}' ./internal/domain/ | grep -v "^encoding/json$" | grep -v "^fmt$" | grep -v "^$" && echo "FAIL: unexpected import" && exit 1 || true
```

## Threat Matrix

N/A — no routing, shell, subprocess, VCS/PR automation, executable-file classification, or process-integration boundary.

## Migration / Rollout

No migration required. Rollback: `git revert` the single commit; restore `references/.gitkeep`.

## Open Questions

None — all technical decisions resolved in exploration and proposal.
