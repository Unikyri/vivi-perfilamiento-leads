# Apply Progress: Issue 30 — Buyer Clock

## Completed Tasks

### Original implementation

- [x] 1.1 Added table-driven Postgres clock coverage for persisted UTC loading, zero-value fallback persistence, simulated advancement, and current UTC wall-time bounds.
- [x] 1.2 Strengthened buyer-persona HTTP characterization for filtered and catalog-wide response shapes with byte-for-byte repeatability, while retaining empty/repeated `proyecto_id` validation coverage.
- [x] 1.3 Verified demo-time HTTP coverage persists one forward advance and exposes the updated simulated date; added an assertion against the simulated clock.
- [x] 2.1 Split `Postgres.Ahora()` to return `time.Now().UTC()` while `FechaSimulada()` and `Avanzar()` remain mutex-protected persisted simulation state.
- [x] 2.2 Documented the non-persistent `relojSistema` fallback and intentional no-op `Avanzar`.
- [x] 3.1 Added `docs/decisiones/0001-conservar-tabla-compradores.md` in Spanish, preserving the buyers boundary and requiring a separate Contract §9 PR with both-block approval for removal.
- [x] 3.2 Confirmed no changes to migrations, buyer data, buyer-persona use case, buyer endpoints, buyer schema, or the `usecase.Reloj` interface.
- [x] 4.1 Focused adapter and HTTP tests pass.
- [x] 4.2 Full tests, vet, and build pass.

### Focused remediation for verify report #651

- [x] R1 Added a committed health test with an injected simulated date (`2099-03-04`) and exact `/salud.fecha_simulada` assertion.
- [x] R2 Added a committed `AvanzarDemo` backward-date test asserting `ErrTiempoSimuladoAtras`, zero persistence saves, and unchanged repository/clock state.
- [x] R3 Added a committed cross-package test using the real `reloj.Postgres` with `SaludarLead`; after advancing simulation to 2099, the operational message timestamp remains within the current UTC wall-time window and before the advanced simulated date.
- [x] R4 Created this OpenSpec `apply-progress.md` mirror from the persisted Engram apply evidence and updated hybrid state.

## Work Unit Evidence

| Evidence | Result |
|---|---|
| Focused test command and exact result | `go test -count=1 -v ./cmd/servidor -run TestSaludReportsInjectedSimulatedDate && go test -count=1 -v ./internal/usecase -run TestAvanzarDemoRejectsBackwardWithoutPersistence && go test -count=1 -v ./internal/infrastructure/reloj -run TestPostgresClockKeepsOperationalWritesOnWallTime` — exit 0; all three named tests passed. |
| Runtime harness command/scenario and exact result | `TestPostgresClockKeepsOperationalWritesOnWallTime` binds the real `reloj.Postgres` to `usecase.SaludarLead` and records an operational message after simulated time advances to 2099; exit 0. The health test exercises the HTTP handler through `httptest`; no external service exists. |
| Rollback boundary | Revert `cmd/servidor/main_test.go`, `internal/usecase/avanzar_demo_test.go`, `internal/infrastructure/reloj/postgres_operational_test.go`, this apply-progress mirror, and the hybrid state update. No production clock interface/call-site, buyer/schema, migration, data, or API behavior rollback is required. |

## Broader Validation

- `go test ./... -count=1` — exit 0; all Go packages passed (packages without tests reported as such).
- `go vet ./...` — exit 0.
- `go build ./...` — exit 0.
- `gofmt -l cmd/servidor/main_test.go internal/usecase/avanzar_demo_test.go internal/infrastructure/reloj/postgres_operational_test.go` — no output; all remediation test files are formatted.

## Scope and Design Compliance

- No production source files were changed by this remediation.
- The `usecase.Reloj` interface and all clock call sites remain unchanged.
- Buyer behavior, buyer schema, migrations, data files, HTTP response shapes, and Contract surfaces remain unchanged.
- The remediation authored 116 test lines, below the 400-line review budget.
- The original implementation remains adapter-only and follows the persisted demo-clock design.

## Deviations

None. The missing committed evidence identified in verify report #651 was added without changing the implementation design.

## Delivery

Mode: single-pr.
Work unit: focused remediation for verify report #651 warnings 1–4.
Next phase: `sdd-verify` after the bounded review gate is satisfied.
