# Tasks: Foundation Constants

## Review Workload Forecast

| Field | Value |
|-------|-------|
| Estimated changed lines | ~180–220 (authored) |
| 400-line budget risk | Low |
| Chained PRs recommended | No |
| Suggested split | Single PR `feat/foundation-constants` → main |
| Delivery strategy | feature-branch-pr-to-main |
| Chain strategy | size-exception |

Decision needed before apply: No
Chained PRs recommended: No
Chain strategy: size-exception
400-line budget risk: Low

### Suggested Work Units

| Unit | Goal | Likely PR | Focused test command | Runtime harness | Rollback boundary |
|------|------|-----------|----------------------|-----------------|-------------------|
| 1 | Canonical assets + embed + domain parser + tests | Single PR | `go build ./... && go test ./internal/domain/... -v -count=1` | N/A — no runtime; compile + unit tests prove correctness | `git revert` single commit; restore `references/.gitkeep` |

## Phase 1: Canonical Assets & Embed Package

- [x] 1.1 Create `references/constantes.json` with exact §4.5 payload: `{"smmlv_2026":1750905,"tasa_ea_default":0.107,"plazo_meses_default":240,"cuota_ingreso_max":0.40,"tope_vis_smmlv":150,"mediana_vis_cop":195000000}`
- [x] 1.2 Create `references/subsidios_2026.json` with exact payload: `[{"ingreso_hasta_smmlv":2,"subsidio_smmlv":30},{"ingreso_hasta_smmlv":4,"subsidio_smmlv":20}]`
- [x] 1.3 Create `references/calendario_economico.json` with exact payload: `[{"tipo":"CESANTIAS","fecha":"--02-14"},{"tipo":"PRIMA","fecha":"--06-30"},{"tipo":"PRIMA","fecha":"--12-20"}]`
- [x] 1.4 Create `references/embed.go`: `package references`; import `_ "embed"`; three `//go:embed` directives exposing `ConstantesJSON`, `SubsidiosJSON`, `CalendarioJSON` as exported `[]byte`. No parsing, no `init()`, no types. Reject any `..` in embed patterns.
- [x] 1.5 Delete `references/.gitkeep` (superseded by real files)

## Phase 2: Pure Domain Parser

- [x] 2.1 Create `internal/domain/constantes.go`: define structs `Constantes`, `TramoSubsidio`, `EventoCalendario` with `json:"snake_case"` tags per design
- [x] 2.2 Implement `LoadConstantes([]byte) (Constantes, error)` using only `encoding/json` + `fmt`
- [x] 2.3 Implement `LoadSubsidios([]byte) ([]TramoSubsidio, error)` — same stdlib-only constraint
- [x] 2.4 Implement `LoadCalendario([]byte) ([]EventoCalendario, error)` — same stdlib-only constraint
- [x] 2.5 Each `Load*` returns `fmt.Errorf("domain: parse %s: %w", name, err)` on failure; no panics

## Phase 3: Tests & Validation

- [x] 3.1 Create `internal/domain/constantes_test.go` with table-driven `TestLoadConstantes` asserting all 6 fields match §1 magnitudes exactly
- [x] 3.2 Add `TestLoadSubsidios` asserting len=2, exact `ingreso_hasta_smmlv` and `subsidio_smmlv` per tramo
- [x] 3.3 Add `TestLoadCalendario` asserting 3 events: CESANTIAS/--02-14, PRIMA/--06-30, PRIMA/--12-20
- [x] 3.4 Add `TestLoad_MalformedJSON` verifying each `Load*` returns non-nil error and does not panic on `[]byte("{{invalid")`
- [x] 3.5 Add `TestLoad_EmptyInput` verifying each `Load*` returns non-nil error on `nil` and `[]byte("")`
- [x] 3.6 Run `go build ./...` — must pass with no `..` embed patterns
- [x] 3.7 Run `go test ./internal/domain/... -v -count=1` — all assertions pass
- [x] 3.8 Run import check: `go list -f '{{join .Imports "\n"}}' ./internal/domain/` — only `encoding/json` and `fmt` allowed

## Phase 4: Delivery Guards

- [x] 4.1 Verify `cmd/servidor/main.go` is NOT modified (deferred to Issue #4)
- [x] 4.2 Verify `cmd/pipeline/main.go` is NOT modified (deferred to Issue #4)
- [x] 4.3 Verify no `Docs/` or `README.md` files created or modified
- [x] 4.4 Verify `references/embed.go` contains no `..` in any embed directive
- [x] 4.5 Confirm PR diff stays within 400-line review budget
