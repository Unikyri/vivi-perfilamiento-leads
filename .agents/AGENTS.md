# Reglas y Guías de Agente — Proyecto Vivi

## 1. Git Flow y Gestión de Ramas (Hackathon)
- **`main`**: Rama de producción / estable. Solo se actualiza por PR verificado en CI/CD.
- **`feature/bloque-b`**: Rama base para el **Bloque B (Datos, Dashboard, Frontend, Pipeline)**.
- **Flujo de Trabajo por Issue**:
  - Para cada nuevo issue del Bloque B, derivar una rama desde `feature/bloque-b` (ejemplo: `feat/issue-XX-nombre`).
  - Desarrollar, probar y verificar con tests/linter.
  - Fusionar mediante merge/PR de regreso a `feature/bloque-b`.

## 2. Desarrollo Guiado por Especificaciones (OpenSpec / SDD)
- Toda nueva característica o cambio importante se planifica en `openspec/changes/<issue-id-nombre>/`:
  - `proposal.md`: Propuesta de valor y alcance.
  - `design.md`: Diseño técnico y arquitectura.
  - `specs/`: Requisitos detallados.
  - `tasks.md`: Lista paso a paso de ejecución.
- Utilizar las Skills cargadas desde `.kiro/skills/` (`sdd-propose`, `sdd-design`, `sdd-spec`, `sdd-tasks`, `sdd-apply`, `sdd-verify`).

## 3. Arquitectura y Código
- **Backend (Go 1.24+)**: Clean Architecture (`cmd/`, `internal/domain`, `internal/usecase`, `internal/adapters`, `internal/infrastructure`).
- **Frontend (TypeScript)**: SPA MVC en `web/` usando Vite + Vanilla CSS / TypeScript.
- **Base de Datos**: PostgreSQL con migraciones SQL en `migrations/` embebidas vía `go:embed`.
- **TDD Pragmático y Flexible**: Escribir pruebas unitarias esenciales y verificar la compilación y tests (`go test ./...`, `web/`) de forma eficiente, evitando loops innecesarios para optimizar consumo de tokens.
