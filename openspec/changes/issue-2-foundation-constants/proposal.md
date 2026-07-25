# Proposal: Foundation Constants (`issue-2-foundation-constants`)

## Intent

Establish the **single canonical source** for normative numeric values (SMMLV, subsidies, economic calendar) as embedded JSON in `references/`, with a pure domain parser that exposes typed structs — without introducing motor business logic. The issue's `../../` embed sample is illegal in Go; this proposal supersedes that path with a legal root-level embed package while preserving all normative values verbatim from Contract §4.5 and Engine Criteria §1.

## Scope

### In Scope
- Three canonical JSON files in `references/` (`constantes.json`, `subsidios_2026.json`, `calendario_economico.json`) with exact Contract v1.1 values
- `references/embed.go` — Go package exposing `[]byte` vars via `//go:embed`
- `internal/domain/constantes.go` — typed structs (`Constantes`, `TramoSubsidio`, `EventoCalendario`) and `Load*()` parsers accepting `[]byte`
- `internal/domain/constantes_test.go` — table-driven tests verifying parsed values match §1 magnitudes
- Composition root acknowledgment: `cmd/servidor` and `cmd/pipeline` **will** import `references` and call `domain.Load*()`, but wiring deferred to Issue #4

### Out of Scope
- Motor business logic (`CalcularCapacidad`, `GemeloKNN`, `Matriz2x2`) — Issue #3+
- HTTP handlers, LLM integration, database schema
- `cmd/servidor/main.go` full wiring (Issue #4)
- `data/*.json` pipeline files (Block B concern)

## Capabilities

### New Capabilities
- `foundation-constants`: Canonical JSON embedding + typed domain parsing of normative economic values

### Modified Capabilities
- None

## Approach

1. **Root embed package** (`references/embed.go`): co-locates `//go:embed` directives with JSON files — legal, no `..` paths.
2. **Pure domain parser** (`internal/domain/constantes.go`): receives `[]byte`, returns typed structs. Zero imports outside stdlib (`encoding/json`). Clean Architecture boundary preserved — domain never imports `references/`.
3. **Dependency injection**: Composition root (`cmd/`) imports `references`, passes bytes to `domain.Load*()`, injects result into services. **Not wired until Issue #4** — only the parser + types ship here.
4. **Why this supersedes the issue's `../../` sample**: Go's `//go:embed` spec explicitly forbids `..` in patterns regardless of whether resolution stays within the module. Placing `embed.go` alongside the JSON files is the only zero-duplication, single-binary, legal solution.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `references/` | New | Three JSON files + `embed.go` package |
| `internal/domain/constantes.go` | New | Typed structs and `Load*()` parsers |
| `internal/domain/constantes_test.go` | New | Table-driven validation against §1 values |
| `cmd/servidor/main.go` | **Deferred (Issue #4)** | Will import `references` and call loaders |
| `cmd/pipeline/main.go` | **Deferred (Issue #4)** | Same pattern |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| `references/` Go package naming collision | Low | Name is unique, no stdlib conflict |
| JSON schema drift from Contract updates | Low | Tests assert exact §1 values; CI catches drift |
| Premature wiring in `cmd/` | Med | Scope explicitly excludes; PR review enforces |

## Rollback Plan

Delete `references/embed.go`, `internal/domain/constantes.go`, `internal/domain/constantes_test.go`, and the three JSON files. Restore `references/.gitkeep`. Single `git revert` of the merged PR.

## Dependencies

- Go 1.24+ (already in `go.mod`) — `//go:embed` available since Go 1.16
- Contract v1.1 §4.5 values frozen (currently BORRADOR but values are normative per §8 of doc 13)

## Success Criteria

- [ ] `go build ./...` passes with embedded JSON (no `..` patterns)
- [ ] `go test ./internal/domain/...` passes with all §1 magnitude assertions
- [ ] `references/constantes.json` values match Contract §4.5 / Engine Criteria §1 exactly
- [ ] Domain package has zero imports outside stdlib
- [ ] `cmd/servidor/main.go` is NOT modified (deferred to Issue #4)
- [ ] PR stays within 400-line review budget
