# Vivi — Perfilamiento Inteligente de Leads

> **Vivi convierte leads en vecinos.**

Asesora digital de vivienda por WhatsApp que **conoce al prospecto antes de saludarlo**: cruza el lead con los datos de afiliados de Colsubsidio, calcula su capacidad de compra con un motor determinista y auditable, recomienda proyectos por similitud (gemelo comprador kNN) y entrega al asesor una ficha inteligente con score, subsidios y prioridad. Hackathon Colsubsidio × 30X.

## Stack

- **Backend:** Go 1.24+ · ADK Go 2.0 · Clean Architecture
- **Frontend:** TypeScript · patrón MVC
- **BD:** PostgreSQL · **LLM:** Gemini (primario) con conmutación a Qwen/Anthropic
- **Deploy:** Heroku · CI/CD desde GitHub

## Estructura (planeada)

```
cmd/          # servidor (composition root) + pipeline de datos
internal/     # domain (motor) · usecase · adapters · infrastructure
web/          # SPA MVC: chat (WhatsApp) + dashboard asesor
skills/       # SKILL.md por skill de agente
data/  references/  migrations/
```

## Documentación

La documentación de diseño vive en la **[Wiki](https://github.com/Unikyri/vivi-perfilamiento-leads/wiki)** (14 documentos: Visión, BRD, PRD, User Stories, Casos de Uso, NFR, **Contrato v1.1**, Arquitectura, PDD, **Criterios del Motor**, Referencia de Dominio).

> El **Contrato** (doc 10) es la frontera entre bloques y **manda ante cualquier conflicto**. Los **Criterios del Motor** (doc 13) fijan los números exactos que el código debe reproducir.

## Inicio local

Requisitos: Go 1.24+, Node/npm y Docker Compose para PostgreSQL. Desde la raíz del repositorio:

```bash
cp .env.example .env
make db-up
make build-todo
make run
curl --fail --silent http://127.0.0.1:8080/salud
```

`make run` mantiene el servidor en primer plano; ejecuta el último comando desde otra terminal. El `.env.example` solo contiene valores locales de ejemplo y no credenciales.

## Mapa RF → paquete

| Módulo (RF) | Paquetes y fuentes principales |
|---|---|
| M1 · Pipeline de datos | `cmd/pipeline`, `internal/pipeline`, `data/`, `references/`, `migrations/` |
| M2 · Motor determinista | `internal/domain/motor`, `internal/domain` |
| M3 · Orquestación y agentes | `internal/adapters/agentes`, `internal/usecase`, `skills/` |
| M4 · Conversación | `internal/usecase`, `internal/adapters/http` |
| M5 · Nutrición | `internal/usecase`, `internal/domain` |
| M6 · Ficha y dashboard | `internal/usecase`, `web/src/views`, `web/src/controllers` |
| M7 · Frontend chat | `web/src` |
| M8 · Infraestructura y CI/CD | `internal/infrastructure`, `.github/`, `cmd/servidor` |
| M9 · Activos del demo | `docs/`, `tests/carga/`, `internal/usecase` |

La API HTTP y la SPA se encuentran en `internal/adapters/http/` y `web/`, respectivamente. La descripción completa de capas y flujos está en [docs/arquitectura.md](docs/arquitectura.md).

## Límites operativos y evidencia local

- El límite HTTP aplica a `/api` y `/api/*`: 30 solicitudes aceptadas por identidad en una ventana fija de 60 segundos. `/salud` y los estáticos no consumen cuota. El estado es **local al proceso**: no es un contador distribuido y cada proceso tendría su propia ventana.
- `Forwarded` y `X-Forwarded-For` se ignoran por defecto porque `TRUSTED_PROXY_CIDRS` está vacío. Solo una red de proxy explícitamente configurada puede aportar la identidad reenviada; el harness de pruebas confía únicamente en `127.0.0.1/32` y `::1/128`.
- La política de carga reproducible y sus resultados medidos viven en [tests/carga/README.md](tests/carga/README.md). El harness es **solo de pruebas**, escucha en loopback, usa respuestas deterministas en proceso y no llama PostgreSQL, proveedores LLM, credenciales ni endpoints públicos. No apuntar k6 a una URL pública ni incluir credenciales en `BASE_URL`.

## Git flow (hackathon — simple)

Dos desarrolladores, dos bloques paralelos que solo se tocan por el Contrato:

| Rama | Bloque | Alcance |
|---|---|---|
| `main` | — | Siempre desplegable. Solo entra por PR con CI verde. |
| `feature/bloque-a` | **A · Núcleo de decisión** | Motor, agentes (ADK), conversación, nutrición |
| `feature/bloque-b` | **B · Datos y experiencia** | Pipeline de datos, dashboard, frontend, CI/CD |

**Flujo:** cada quien trabaja en su rama de bloque → PR a `main` → merge tras CI verde. Cambios al Contrato (doc 10) van en PR aparte etiquetado `contrato`, con aprobación de ambos. Nadie importa código del otro bloque: la integración es la API HTTP + los archivos de datos.

```bash
git checkout feature/bloque-a   # o bloque-b
# ...trabajar, commitear...
git push
# abrir PR hacia main
```
