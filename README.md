# Vivi — Perfilamiento Inteligente de Leads

> **Vivi convierte leads en vecinos.**

## Sobre el proyecto

Muchas familias afiliadas a Colsubsidio ya tienen una oportunidad real de vivienda y nunca la descubren: el primer contacto llega tarde, la información es difícil de entender, o el asesor debe empezar cada conversación desde cero. Un lead se enfría antes de saber si tiene subsidios, capacidad de compra o un proyecto posible.

**Vivi no es un chatbot ni una única asesora digital: es un equipo comercial digital dedicado a cada lead.** Mientras una asesora conversa con calidez por WhatsApp, un equipo de agentes especializados perfila la necesidad, calcula la capacidad financiera, identifica subsidios, recomienda proyectos por similitud (gemelo comprador kNN) y prepara el mejor momento para que intervenga un asesor humano.

> No reemplazamos al asesor: le entregamos una conversación que ya tiene contexto, prioridad y una ruta posible.

## Arquitectura multiagente

Cada lead avanza por un equipo de agentes coordinados mediante un **bus de eventos en memoria** y una **Coordinadora** (patrón Mediator, `internal/adapters/agentes`) — sin acoplar a los agentes entre sí:

| Agente | Rol | Dónde vive |
|---|---|---|
| **Saludo** | Primer contacto y consentimiento de datos (Ley 1581) | `usecase.SaludarLead` |
| **Perfilador** | Convierte la conversación en una ficha estructurada | `usecase.PerfilarLead` / `ProcesarMensaje` |
| **Analista de capacidad** | Calcula crédito, subsidio y presupuesto — **motor determinista, cero LLM** | `internal/domain/motor` |
| **Especialista en subsidios** | Identifica los subsidios aplicables y por qué | `internal/domain/motor` |
| **Recomendador** | Encuentra el proyecto más compatible (kNN + matriz 2×2) | `internal/domain/motor` |
| **Nutricionista** | Si el lead aún no está listo, lo acompaña con un plan y calendario | `usecase.GestionarPlan` |

La IA conversa, entiende intención y coordina el equipo; **no inventa decisiones financieras**. Los cálculos y subsidios salen de un motor auditable con reglas y fuentes explícitas (`VERIFICADO_BASE` / `DECLARADO` / `INFERIDO`), no de un LLM. Cada agente además carga su propio *skill* (`skills/*/SKILL.md`: tono Colsubsidio, dominio de caja, explicación financiera humana, siguiente mejor pregunta, entre otros) para mantener un comportamiento consistente y auditable.

## Construido con

- **Backend:** Go 1.24+ · Clean Architecture · bus de eventos en memoria + Mediator (sin framework externo de agentes)
- **Frontend:** TypeScript · patrón MVC · SPA sin dependencias de UI
- **Base de datos:** PostgreSQL
- **LLM:** Gemini (primario), con conmutación automática a Qwen ante fallas (circuit breaker + fallback en cascada)
- **Despliegue:** Heroku · CI/CD con GitHub Actions

## Empezar

### Prerrequisitos

- Go 1.24+
- Node.js y npm
- Docker Compose (para PostgreSQL local)

### Instalación

```bash
cp .env.example .env
make db-up
make build-todo
make run
curl --fail --silent http://127.0.0.1:8080/salud
```

`make run` mantiene el servidor en primer plano; ejecutá el `curl` desde otra terminal. `.env.example` sólo trae valores locales de ejemplo, nunca credenciales reales.

## Demo

**[vivi-37863aed9d29.herokuapp.com](https://vivi-37863aed9d29.herokuapp.com/)**

La documentación de diseño completa (Visión, BRD, PRD, User Stories, Casos de Uso, NFR, **Contrato v1.1**, Arquitectura, PDD, **Criterios del Motor**, Referencia de Dominio) vive en la **[Wiki](https://github.com/Unikyri/vivi-perfilamiento-leads/wiki)**. El **Contrato** (doc 10) es la frontera entre bloques y manda ante cualquier conflicto; los **Criterios del Motor** (doc 13) fijan los números exactos que el código reproduce.

## Mapa RF → paquete

| Módulo | Paquetes principales |
|---|---|
| Pipeline de datos | `cmd/pipeline`, `internal/pipeline`, `data/`, `references/` |
| Motor determinista | `internal/domain/motor`, `internal/domain` |
| Agentes y orquestación | `internal/adapters/agentes`, `internal/usecase`, `skills/` |
| Conversación (WhatsApp) | `internal/usecase`, `internal/adapters/http` |
| Nutrición | `internal/usecase`, `internal/domain` |
| Ficha y dashboard | `internal/usecase`, `web/src/views`, `web/src/controllers` |
| Infraestructura y CI/CD | `internal/infrastructure`, `.github/`, `cmd/servidor` |

Descripción completa de capas y flujos en [docs/arquitectura.md](docs/arquitectura.md).

## Hoja de ruta

El siguiente paso es un **piloto con afiliados y asesores reales** para medir tiempo de respuesta, calidad de perfilamiento y tasa de entrega comercial. En el corto plazo:

- Cerrar la vista de Nutrición y el catálogo de Proyectos en el dashboard del asesor.
- Prueba de carga formal y rate limiting distribuido (hoy el límite es local al proceso).
- Integración con los canales oficiales de Colsubsidio más allá del demo de WhatsApp simulado.

## Licencia

MIT — ver [LICENSE](LICENSE).

## Contacto

Juan Arango · Nicolás Lozano · David Hernández Ortiz

## Agradecimientos

Hackathon Colsubsidio × 30X.
