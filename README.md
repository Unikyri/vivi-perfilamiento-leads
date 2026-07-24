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

La documentación de diseño vive en la **[Wiki](../../wiki)** (14 documentos: Visión, BRD, PRD, User Stories, Casos de Uso, NFR, **Contrato v1.1**, Arquitectura, PDD, **Criterios del Motor**, Referencia de Dominio).

> El **Contrato** (doc 10) es la frontera entre bloques y **manda ante cualquier conflicto**. Los **Criterios del Motor** (doc 13) fijan los números exactos que el código debe reproducir.

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
