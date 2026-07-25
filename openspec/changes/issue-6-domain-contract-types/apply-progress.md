# Apply Progress: PR4 — Enums and Perfil

## Change
`issue-6-domain-contract-types`

## Delivery Boundary
- Strategy: `auto-chain`
- Chain: `stacked-to-main`
- Current slice: PR4 / 7
- Scope: tasks 1.1, 2.1, and 2.2 only

## Mode
Standard mode (`strict_tdd: false` in `openspec/config.yaml`).

## Completed Tasks
- [x] 1.1 Added table-driven profile accessor, verification, and exact key-set tests.
- [x] 2.1 Added the 11 typed-string Contract v1.1 enum groups and exact literals, including `EstadoHito`.
- [x] 2.2 Added `CampoPerfil`, `Perfil`, `CamposReconocidos`, `CamposCriticos`, and only the four specified accessors.

## Remaining Tasks
- [ ] 1.2 Aggregate enum/profile reflection and JSON assertions for the complete exported schema; enum/profile-focused coverage added in this slice, while aggregate-schema assertions remain for later slices.
- [ ] 2.3 Capacity/recommendation declarations and legacy `Comprador`/`Proyecto` relocation.
- [ ] 2.4 Lead, message, nutrition, and ficha declarations.
- [ ] 3.1 Full domain and pipeline tests.
- [ ] 3.2 Build, vet, and dependency-isolation checks.
- [ ] 3.3 Final scope/evidence confirmation.

## Files Changed
- `internal/domain/enums.go` — created 11 typed string enum groups and constants.
- `internal/domain/perfil.go` — created exact profile schema, key sets, and accessors.
- `internal/domain/perfil_test.go` — created enum JSON/underlying-type and profile behavior/reflection/JSON tests.
- `openspec/changes/issue-6-domain-contract-types/tasks.md` — marked 1.1, 2.1, and 2.2 complete only.

## Work Unit Evidence
| Evidence | Result |
|---|---|
| Focused test command and exact result | `go test ./internal/domain/... -run 'Test(Enum|Perfil)' -v` — exit 0; all enum/profile tests passed. |
| Runtime harness command/scenario and exact result | N/A — this slice contains pure standard-library domain declarations with no executable, HTTP, database, LLM, or integration runtime boundary. |
| Rollback boundary | Remove `internal/domain/enums.go`, `internal/domain/perfil.go`, and `internal/domain/perfil_test.go`; restore the three task checkboxes to pending. No unrelated files require rollback. |

## Deviations and Risks
None. Implementation follows the supplied spec/design. Aggregate schema assertions are intentionally deferred with task 1.2 and later domain slices.

## Validation
- `gofmt -w internal/domain/enums.go internal/domain/perfil.go internal/domain/perfil_test.go` — passed.
- `git diff --check` — passed.
