# Apply Progress: PR4 — Enums and Perfil

## Scope and Mode
- Change: `issue-6-domain-contract-types`; standard mode (`strict_tdd: false`).
- Delivery: `auto-chain`, `stacked-to-main`, PR 4/7.
- Completed: 1.1 profile tests; 2.1 eleven enums (including `EstadoHito`); 2.2 `CampoPerfil`, `Perfil`, exact key sets, and four accessors.
- Remaining: 1.2 aggregate schema tests; 2.3 capacity relocation; 2.4 lead/plan/ficha; 3.1–3.3 final evidence.

## Changed Files
- `internal/domain/enums.go` — exact typed-string Contract enums.
- `internal/domain/perfil.go` — exact profile schema and accessors.
- `internal/domain/perfil_test.go` — enum/profile behavior and schema tests.
- `tasks.md` — only 1.1, 2.1, and 2.2 marked complete.

## Work Unit Evidence
| Evidence | Result |
|---|---|
| Focused test | `go test ./internal/domain/... -run 'Test(Enum|Perfil)' -v` — exit 0, all passed. |
| Runtime harness | N/A — pure standard-library declarations; no runtime boundary. |
| Rollback | Remove the three domain files and restore the three task checkboxes. |

No deviations or risks. `git diff --check` and `gofmt` passed.
