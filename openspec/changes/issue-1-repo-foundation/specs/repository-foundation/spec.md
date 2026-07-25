# Repository Foundation Specification

## Purpose
Provide the exact Go repository scaffold required by Issue #1 while preserving Contract v1.1.

## Requirements

### Requirement: Module and Package Declarations
The repository MUST create `go.mod` with module `github.com/Unikyri/vivi-perfilamiento-leads` and a version line exactly `go 1.24`. It MUST create these eight `doc.go` declarations exactly as listed in Issue #1:

| Path | Required declaration text |
|---|---|
| `internal/domain/doc.go` | `Package domain contains the business entities and pure domain services.` / `It must NOT import anything from usecase, adapters or infrastructure.` / `package domain` |
| `internal/domain/motor/doc.go` | `Package motor contains the deterministic decision engine.` / `Zero external dependencies, zero I/O, zero LLM calls.` / `package motor` |
| `internal/usecase/doc.go` | `Package usecase contains application use cases and the ports (interfaces)` / `they depend on. It may import domain, never adapters or infrastructure.` / `package usecase` |
| `internal/adapters/http/doc.go` | `Package http contains HTTP controllers and presenters.` / `package http` |
| `internal/adapters/agentes/doc.go` | `Package agentes contains the ADK agent graph wiring.` / `package agentes` |
| `internal/infrastructure/postgres/doc.go` | `Package postgres implements the repository ports using PostgreSQL.` / `package postgres` |
| `internal/infrastructure/llm/doc.go` | `Package llm implements the LLMProvider port for each vendor.` / `package llm` |
| `internal/infrastructure/config/doc.go` | `Package config loads configuration from environment variables (12-factor).` / `package config` |

Each slash denotes a newline; comment lines MUST retain Go `// ` prefixes.

#### Scenario: Inspect module and package skeleton
- GIVEN a checkout before Issue #1 implementation
- WHEN the required module and eight package files are inspected
- THEN their paths and contents match this requirement exactly
- AND no third-party module dependency is introduced

### Requirement: Placeholder Commands
`cmd/servidor/main.go` and `cmd/pipeline/main.go` MUST be compile-safe `main` packages that import only `fmt` and call `fmt.Println` with, respectively, `vivi servidor - placeholder` and `vivi pipeline - placeholder`.

#### Scenario: Run placeholders
- GIVEN the foundation is built with Go 1.24+
- WHEN `go run ./cmd/servidor` and `go run ./cmd/pipeline` run
- THEN stdout is exactly the respective required placeholder line

### Requirement: Makefile Commands
The root `Makefile` MUST declare `.PHONY: build test vet run datos limpiar`; each target MUST execute exactly:

```makefile
.PHONY: build test vet run datos limpiar

build:
	go build -o bin/servidor ./cmd/servidor
	go build -o bin/pipeline ./cmd/pipeline

test:
	go test ./... -count=1

vet:
	go vet ./...

run:
	go run ./cmd/servidor

datos:
	go run ./cmd/pipeline

limpiar:
	rm -rf bin/
```

Recipe lines MUST be TAB-indented.

#### Scenario: Use build lifecycle targets
- GIVEN GNU Make and Go are available at repository root
- WHEN `make build`, `make test`, `make vet`, `make run`, `make datos`, and `make limpiar` execute
- THEN each invokes its specified command
- AND build creates both binaries while limpiar removes `bin/`

### Requirement: Tracked Directories and Ignore Rule
The repository MUST track only `web/.gitkeep`, `data/.gitkeep`, `references/.gitkeep`, `migrations/.gitkeep`, and `skills/.gitkeep` from the Issue #1 empty-directory list. It MUST append, without rewriting existing rules, `# Binarios compilados` followed by `bin/` as the final `.gitignore` block. It MUST NOT create Issue-omitted `docs/` or `contracts/schemas/` directories.

#### Scenario: Inspect tracked placeholders and ignore behavior
- GIVEN a clean foundation checkout
- WHEN the listed directories and `.gitignore` are inspected after a build
- THEN exactly the five `.gitkeep` paths exist and `bin/` is ignored
- AND omitted directory paths are absent

### Requirement: Scope and Validation Boundary
Implementation MUST NOT modify `Docs/` or `README.md`, define API/contracts, or implement domain, motor, use-case, adapter, infrastructure, or frontend behavior. It MUST pass `go build ./...`, `go vet ./...`, `make build`, `ls bin/servidor bin/pipeline`, and `go run ./cmd/servidor` with the required server output.

#### Scenario: Validate the closed Issue #1 scope
- GIVEN the implemented change is compared with its baseline
- WHEN the required validation commands and path diff are run
- THEN every command exits successfully and both binaries exist
- AND `Docs/` and `README.md` have no diff
