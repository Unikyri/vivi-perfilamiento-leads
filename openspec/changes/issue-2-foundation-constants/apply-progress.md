# Apply Progress: issue-2-foundation-constants

**Change**: issue-2-foundation-constants
**Mode**: Standard
**Status**: 18/18 tasks complete. Ready for verify.

## What

Completed all 18 tasks for issue-2-foundation-constants (Phases 1–4).
Establish single canonical source for normative economic values as embedded JSON.

## Files Changed

| File | Action | What Was Done |
|------|--------|---------------|
| `references/constantes.json` | Created | Canonical §4.5 economic constants (SMMLV, tasa, plazo, cuota, tope, mediana) |
| `references/subsidios_2026.json` | Created | Subsidy tiers (2 tramos: ≤2 SMMLV → 30, ≤4 SMMLV → 20) |
| `references/calendario_economico.json` | Created | 3 economic events (CESANTIAS, 2× PRIMA) with recurring month-day |
| `references/embed.go` | Created | Go embed package exposing `ConstantesJSON`, `SubsidiosJSON`, `CalendarioJSON` as `[]byte` |
| `internal/domain/constantes.go` | Created | Pure domain parser: structs + `Load*` functions, stdlib-only (encoding/json + fmt) |
| `internal/domain/constantes_test.go` | Created | Table-driven tests: 5 test functions, 20 subtests covering happy path + malformed + empty |

## Work Unit Evidence

| Evidence | Value |
|----------|-------|
| Focused test command | `go build ./... && go test ./internal/domain/... -v -count=1` → exit 0, 5 tests PASS (20 subtests) |
| Runtime harness | N/A — no runtime boundary; compile + unit tests prove correctness |
| Rollback boundary | `git revert` single commit; restore `references/.gitkeep`; delete 5 new files |

## Completed Tasks: 18/18

### Phase 1 (1.1–1.5): Canonical JSON assets + embed package ✓

- [x] 1.1 Create `references/constantes.json` with exact §4.5 payload
- [x] 1.2 Create `references/subsidios_2026.json` with exact subsidy tiers
- [x] 1.3 Create `references/calendario_economico.json` with 3 economic events
- [x] 1.4 Create `references/embed.go` with `//go:embed` directives (no `..`, no `init()`)
- [x] 1.5 Delete `references/.gitkeep` (superseded by real files)

### Phase 2 (2.1–2.5): Pure domain parser with stdlib-only imports ✓

- [x] 2.1 Define `Constantes`, `TramoSubsidio`, `EventoCalendario` structs
- [x] 2.2 Implement `LoadConstantes([]byte) (Constantes, error)`
- [x] 2.3 Implement `LoadSubsidios([]byte) ([]TramoSubsidio, error)`
- [x] 2.4 Implement `LoadCalendario([]byte) ([]EventoCalendario, error)`
- [x] 2.5 Each `Load*` returns `fmt.Errorf("domain: parse %s: %w", name, err)` on failure

### Phase 3 (3.1–3.8): Tests + build + import validation ✓

- [x] 3.1 `TestLoadConstantes` — all 6 fields match §1 magnitudes
- [x] 3.2 `TestLoadSubsidios` — len=2, exact values per tramo
- [x] 3.3 `TestLoadCalendario` — 3 events with correct tipos/fechas
- [x] 3.4 `TestLoad_MalformedJSON` — non-nil error, no panics
- [x] 3.5 `TestLoad_EmptyInput` — non-nil error on nil/empty
- [x] 3.6 `go build ./...` → exit 0
- [x] 3.7 `go test ./internal/domain/... -v -count=1` → PASS (0.004s)
- [x] 3.8 Import check: only `encoding/json` and `fmt`

### Phase 4 (4.1–4.5): Delivery guards verified ✓

- [x] 4.1 `cmd/servidor/main.go` NOT modified
- [x] 4.2 `cmd/pipeline/main.go` NOT modified
- [x] 4.3 No `Docs/` or `README.md` files created/modified
- [x] 4.4 No `..` in embed directives
- [x] 4.5 PR diff within 400-line review budget

## Validation Results

- `go build ./...` → exit 0
- `go test ./internal/domain/... -v -count=1` → PASS (0.004s)
- `go list -f '{{join .Imports "\n"}}' ./internal/domain/` → encoding/json, fmt (only)
- No `..` in embed.go
- No prohibited paths modified (cmd/, Docs/, README.md)
- gofmt clean on all files

## Deviations from Design

None — implementation matches design.

## Issues Found

None.

## Workload / PR Boundary

- Mode: Single PR
- Current work unit: Canonical assets + embed + domain parser + tests
- Boundary: Full issue scope in one commit
- Estimated review budget impact: ~180–220 lines (within 400 budget)

## Learned

- Go `//go:embed` must be co-located with files in same package directory
- `fmt.Errorf` wrapping provides clean domain errors with context

---

*Persisted: 2026-07-24T21:16:16-05:00*
*Source: Engram observation #443, topic sdd/issue-2-foundation-constants/apply-progress*
*Session: foundation-issues-20260724*
