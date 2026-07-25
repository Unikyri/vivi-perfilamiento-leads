# Exploration: Foundation Constants (`issue-2-foundation-constants`)

## Problem Statement

Issue #2 requires that normative numeric values (SMMLV, subsidy table, economic calendar) live in **exactly one canonical location** (`references/*.json` at the repo root) while the Go binary remains self-contained (Heroku single-binary deploy). The issue's sample code proposes:

```go
//go:embed ../../references/constantes.json
var constantesJSON []byte
```

**This is illegal in Go.** The `//go:embed` directive resolves patterns **relative to the package directory** and explicitly forbids paths that traverse parent directories (`..`). The Go compiler will reject it at build time.

## Current State

- `references/` exists at the repo root (currently has only `.gitkeep`).
- `internal/domain/` is the domain layer (package `domain`); depends on nothing outside itself.
- `internal/domain/motor/` is a sub-package for the deterministic engine (package `motor`); also zero external deps.
- Clean Architecture rule: domain MUST NOT import infrastructure, adapters, or usecase.
- Module: `github.com/Unikyri/vivi-perfilamiento-leads`, Go 1.24.

## Affected Areas

- `references/` — canonical JSON files (source of truth for all blocks)
- `internal/domain/` — typed constants exposure (consumed by motor, kNN, nutrición)
- `internal/domain/motor/` — the deterministic engine that reads constants
- `cmd/servidor/main.go` — composition root (where embed can legally live)
- `cmd/pipeline/main.go` — Block B composition root (also consumes constants)

## Approaches

### 1. **Symlink / Copy at Build** — rejected

- Description: A Makefile step copies `references/*.json` into `internal/domain/` before build.
- Pros: Embed just works with local paths.
- Cons: **Duplicates** the canonical data (violates single-source-of-truth). Race conditions. CI must remember the step. Git-tracked copies are confusing.
- Effort: Low
- **Verdict: REJECTED** — violates the "one place" requirement.

### 2. **Embed at Composition Root (`cmd/`) + Inject via `Init()` Function**

- Description: Place `//go:embed` in `cmd/servidor/` (or a shared internal package at the module root level) which *can* legally embed `../../references/*.json`... **No — still illegal.** `cmd/servidor/` is at depth 2 from root; `references/` is at depth 1. The pattern `../../references/constantes.json` goes above the module root, which embed forbids. **But:** `cmd/servidor/` can use a relative path only if it resolves *within* the module tree. Let's check: from `cmd/servidor/`, `../../references/` = repo root's `references/` — that IS within the module tree. However, Go embed still forbids `..` in patterns entirely (the spec says: "Patterns must not contain '.' or '..' or empty path elements").

- **Verdict: REJECTED** — Go's embed spec forbids `..` regardless of whether the resolved path stays within the module.

### 3. **Dedicated Embed Package at Repo Root** ✅ RECOMMENDED

- Description: Create a new package `references/` (turning the directory into a Go package) with an `embed.go` that embeds the co-located JSON files. Domain code accesses the raw bytes via a thin, dependency-free interface.
- Layout:
  ```
  references/
  ├── constantes.json          (canonical data — unchanged)
  ├── subsidios_2026.json      (canonical data — unchanged)
  ├── calendario_economico.json(canonical data — unchanged)
  └── embed.go                 (Go package that embeds the above)
  ```
- `embed.go`:
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
- Domain access pattern:
  ```go
  // internal/domain/constantes.go
  package domain

  import "encoding/json"

  type Constantes struct { ... }
  type TramoSubsidio struct { ... }

  // LoadConstantes parses raw JSON into typed values.
  // Called once by the composition root with bytes from the embed package.
  func LoadConstantes(raw []byte) (Constantes, error) { ... }
  func LoadSubsidios(raw []byte) ([]TramoSubsidio, error) { ... }
  ```
- Composition root wires it:
  ```go
  // cmd/servidor/main.go
  import "github.com/Unikyri/vivi-perfilamiento-leads/references"
  import "github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"

  func main() {
      cts, _ := domain.LoadConstantes(references.ConstantesJSON)
      // inject cts into services...
  }
  ```
- Pros:
  - **Single source of truth**: JSON lives in exactly one place.
  - **Self-contained binary**: embed bakes bytes into the binary at compile time.
  - **Legal Go**: embed is in the same directory as the files.
  - **Clean Architecture**: domain never imports `references/`; it receives `[]byte` and parses.
  - **Both blocks share**: Block B's `cmd/pipeline/` imports the same `references` package.
  - **Testable**: domain tests pass raw JSON literals or `references.ConstantesJSON` — no I/O.
- Cons:
  - `references/` becomes a Go package (adds `embed.go` alongside the JSON). Minor aesthetic change.
  - Domain layer needs an explicit `LoadConstantes([]byte)` call rather than implicit `init()`.
- Effort: **Low**

### 4. **Internal Embed Package (`internal/embed/references/`)**

- Description: Create `internal/embed/references/` that copies/symlinks or uses `go:generate` to pull from root `references/`. Domain imports `internal/embed/references`.
- Pros: Keeps `references/` free of Go code.
- Cons: Requires duplication or a generation step. More complex. Domain importing an `internal/` sibling is fine architecturally but the data duplication is the problem.
- Effort: Medium
- **Verdict: INFERIOR** to Option 3.

### 5. **`go generate` JSON→Go Code**

- Description: A generator reads `references/*.json` and emits `internal/domain/constantes_generated.go` with typed constants as Go literals.
- Pros: Zero runtime JSON parsing. Pure Go constants.
- Cons: **Duplicates** values (JSON + generated Go). Must re-run generator on any change. CI complexity.
- Effort: Medium
- **Verdict: REJECTED** — duplicates the canonical values.

## Recommendation

**Option 3: Dedicated embed package at `references/`.**

This is the simplest layout that satisfies all constraints simultaneously:

1. ✅ **One canonical source** — the JSON files at `references/`.
2. ✅ **Self-contained binary** — `//go:embed` bakes them in at compile time.
3. ✅ **Legal Go** — embed pattern is a local file, no `..` needed.
4. ✅ **Clean Architecture preserved** — `internal/domain` never imports `references/`; it exports `LoadConstantes([]byte)` and the composition root passes in the raw bytes.
5. ✅ **Cross-block sharing** — both `cmd/servidor/` and `cmd/pipeline/` (or any future consumer) just `import .../references"`.
6. ✅ **No numeric duplication** — values exist only in JSON; Go struct fields are populated at startup by JSON unmarshal.

### Downstream Access Pattern (no outer-layer imports in domain)

```
references/embed.go          → exposes []byte vars (pure data, no logic)
         ↑ imported by
cmd/servidor/main.go         → passes bytes down
         ↓ calls
internal/domain.LoadConstantes([]byte) → returns typed Constantes struct
         ↓ injected into
internal/domain/motor/       → receives Constantes as function param or struct field
```

Domain tests import `references` directly for convenience (test binaries aren't production; alternatively they embed test fixtures or inline JSON). This is acceptable because `references/` contains **only** raw data with zero logic — it's effectively a static asset package, not a violation of the dependency rule.

### Alternative for Test Purity

If strict layering forbids even test imports of `references/`, domain tests can:
- Inline the expected JSON as a `[]byte` literal in the test file.
- Or use `go:embed` within a `_test.go` file (test files may import anything).

## Risks

- **Package naming collision**: The directory `references/` becomes a Go package. Its `package references` name is unambiguous and descriptive. No collision with stdlib.
- **Future non-Go consumers**: If Block B's TypeScript frontend or pipeline needs the JSON, they still read from `references/*.json` directly — the added `embed.go` doesn't interfere.
- **`init()` vs explicit load**: The issue's sample uses `init()` with package-level `var`. The recommended approach uses explicit `LoadConstantes()` which is more testable and explicit, but requires the composition root to wire it. A hybrid (package-level singleton set by `init()` in a thin loader) is viable if the team prefers the simpler API from the issue.

## Ready for Proposal

**Yes.** The core technical question (how to legally embed `references/` into the binary while keeping domain clean) is answered. The proposal phase should formalize:
1. Exact file list and package layout.
2. Whether to use `init()` (global singleton) vs explicit `Load()` (dependency injection).
3. The test strategy (domain tests import `references` or inline JSON).
4. The `calendario_economico.json` type definitions.
