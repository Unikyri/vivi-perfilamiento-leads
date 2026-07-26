# Arquitectura de Software — Vivi · Fase 0

Este documento transcribe los niveles C4 1–3 y los flujos críticos de la arquitectura de Fase 0. Los diagramas Mermaid se conservan para renderizado nativo en GitHub.

> **Fuente:** `11‐Arquitectura-de-Software-—-Vivi-·-Fase-0.md`, Wiki de Vivi. Se transcriben sus secciones 1–4; las secciones de estado, despliegue y ADR quedan fuera de este slice documental.

## 1. C4 — Nivel 1: Contexto

```mermaid
flowchart TB
    lead(["👤 Lead<br/>(Ana / Carlos)"])
    marcela(["👤 Asesora comercial<br/>(Marcela)"])
    jonathan(["👤 Gerencia<br/>(Jonathan)"])
    usuario(["👤 Usuario del demo"])

    vivi["🏠 VIVI<br/>Sistema multiagente de<br/>perfilamiento de leads<br/>(Heroku)"]

    gemini["🤖 Google Gemini API<br/>(LLM primario)"]
    qwen["🤖 Qwen API<br/>(LLM respaldo)"]

    meta["📣 Meta Ads<br/>(lead forms)"]
    wa["💬 WhatsApp Cloud API"]
    sf["☁️ Salesforce CRM"]

    lead -->|"conversa (chat / audio)"| vivi
    usuario -->|"recorre el demo"| vivi
    vivi -->|"cola priorizada + fichas"| marcela
    vivi -->|"buyer persona vivo"| jonathan
    vivi -->|"turnos LLM"| gemini
    vivi -.->|"fallback"| qwen
    meta -.->|"Fase 1: webhook real"| vivi
    vivi -.->|"Fase 1: canal real"| wa
    vivi -.->|"Fase 1: fichas → CRM"| sf

    style vivi fill:#1a6b52,color:#fff
    style meta stroke-dasharray: 5 5
    style wa stroke-dasharray: 5 5
    style sf stroke-dasharray: 5 5
```

Línea punteada = fuera de Fase 0 (adaptadores previstos, no implementados). En Fase 0 el "lead form de Meta" se simula con `POST /api/leads` y el canal WhatsApp con el chat web (más sandbox Twilio solo para el video).

## 2. C4 — Nivel 2: Contenedores

```mermaid
flowchart LR
    subgraph heroku["Heroku (una app)"]
        front["📱 SPA Frontend<br/>TypeScript · patrón MVC<br/>vistas: chat / dashboard / gerencia"]
        back["⚙️ Backend Go 1.24<br/>Clean Architecture + ADK Go 2.0<br/>binario único: API + estáticos"]
        db[("🗄️ Heroku Postgres<br/>estado único del sistema<br/>(blackboard)")]
    end

    front -->|"JSON · contrato v1.0<br/>polling 1.5s conv / 5s cola"| back
    back -->|"pgx pool + optimistic lock"| db
    back -->|"LLMProvider (puerto)"| gemini["Gemini API"]
    back -.->|"circuit breaker abierto"| qwen["Qwen API<br/>(OpenAI-compatible)"]
```

El backend sirve los estáticos del front compilado (un solo dyno, cero CORS). Toda comunicación front↔back es el contrato `10-contratos.md`; el front jamás toca la BD ni los LLM.

## 3. C4 — Nivel 3: Componentes del backend (capas Clean)

```mermaid
flowchart TB
    subgraph adapters["internal/adapters — mecanismos de entrega"]
        http["Controladores HTTP<br/>+ Presenters (JSON contrato)"]
        grafo["Grafo ADK (Mediator)<br/>Coordinadora · Investigadora · Asesora<br/>Nutricionista · Documentadora<br/>+ skills/*.md inyectadas"]
    end

    subgraph usecase["internal/usecase — reglas de aplicación"]
        ucs["Casos de uso (1:1 con UC-01..08)<br/>PerfilarLead · ProcesarMensaje · CalificarLead<br/>GenerarFicha · EjecutarHitos · GestionarPlan"]
        puertos["Puertos (interfaces)<br/>LeadRepository · PlanRepository · FichaRepository<br/>LLMProvider · Reloj · BusEventos"]
    end

    subgraph domain["internal/domain — reglas de negocio puras"]
        entidades["Entidades + máquina de estados (State)<br/>Lead · Perfil · PlanNutricion · Ficha"]
        motor["MOTOR (servicios de dominio, sin I/O)<br/>CalcularCapacidad · GemeloKNN<br/>RecomendarProyectos · Matriz2x2"]
    end

    subgraph infra["internal/infrastructure — detalles"]
        pg["postgres/ (repos)"]
        llm["llm/ GeminiProvider · QwenProvider<br/>Decorators: Guardarrailes → CircuitBreaker → Métricas<br/>Factory: NewLLMProvider(cfg)"]
        cfg["config/ (12-factor)"]
    end

    http --> ucs
    grafo --> ucs
    ucs --> entidades
    ucs --> motor
    ucs --> puertos
    pg -. "implementa" .-> puertos
    llm -. "implementa" .-> puertos

    style domain fill:#f5efe0
```

**Regla de dependencia (verificada en CI con lint de imports):** las flechas solo apuntan hacia adentro. `domain` no importa nada; `usecase` solo importa `domain`; `adapters` e `infrastructure` importan hacia adentro y nunca entre sí. El grafo ADK vive en adapters porque es mecanismo de orquestación/entrega, no regla de negocio: sus tools delegan en casos de uso.

## 4. Diagramas de secuencia (flujos críticos)

### 4.1 Turno conversacional (el flujo del presupuesto <3s, NFR-R-01)

```mermaid
sequenceDiagram
    autonumber
    participant F as Front (chat MVC)
    participant H as Controlador HTTP
    participant UC as ProcesarMensaje
    participant GE as Guardarrail entrada
    participant R as LeadRepository (PG)
    participant L as LLMProvider<br/>(Decorators+Fallback)
    participant M as Motor (dominio)

    F->>H: POST /api/leads/{id}/mensajes
    H-->>F: 202 {turno_en_proceso:true}
    Note over F: poll GET /conversacion cada 1.5s<br/>muestra "escribiendo…"
    H->>UC: ejecutar(lead_id, mensaje)
    UC->>GE: clasificar(texto)
    alt adversario / fuera de dominio
        GE-->>UC: plantilla de redirección (sin LLM)
    else entrada legítima
        UC->>R: perfil + capacidad vigente
        UC->>L: GenerarTurno (UNA llamada, timeout 8s)
        Note over L: 429/timeout ×3 → breaker abre<br/>→ QwenProvider transparente
        L-->>UC: SalidaTurno JSON (contrato §7)
        UC->>UC: guardarrail salida:<br/>cifras ∈ NUMEROS_DEL_MOTOR
        UC->>M: CalcularCapacidad(perfil ∪ extraídos)
        M-->>UC: Capacidad (determinista)
        UC->>R: guardar campos, mensajes, estado (version+1)
        opt accion = PERFIL_COMPLETO
            UC->>UC: publicar PerfilCompleto en bus
        end
    end
    F->>H: GET /conversacion?desde=...
    H-->>F: mensajes nuevos de VIVI
```

### 4.2 Lead nuevo → saludo personalizado (UC-01 pasos 1–5)

```mermaid
sequenceDiagram
    autonumber
    participant F as Front
    participant H as HTTP
    participant CO as Coordinadora (ADK)
    participant IN as Investigadora (ADK)
    participant UC as PerfilarLead
    participant R as Repos (PG)
    participant AS as Asesora (ADK)

    F->>H: POST /api/leads {cedula|precargado_id}
    H-->>F: 201 {lead_id, afiliado_detectado}
    H->>CO: evento LeadNuevo
    CO->>IN: delegar(lead_id)
    IN->>UC: perfilar(lead_id)
    UC->>R: consultar_afiliado(cedula) en afiliados_mock
    alt match (Ana)
        UC->>UC: campos → VERIFICADO_BASE + motor.precalcular
    else sin match (Carlos)
        UC->>UC: marca no_encontrado, rama SIN_DETERMINAR
    end
    UC->>R: guardar perfil (estado PERFILANDO)
    CO->>AS: delegar(saludo, lead_id)
    AS->>R: leer pre-perfil → redactar saludo (skill tono)
    AS->>R: persistir Mensaje VIVI
    Note over F: el saludo aparece por polling
```

### 4.3 TickReloj → hito de nutrición (UC-07)

```mermaid
sequenceDiagram
    autonumber
    participant F as Front (botón "avanzar tiempo")
    participant H as HTTP
    participant CO as Coordinadora
    participant NU as Nutricionista (ADK)
    participant UC as EjecutarHitos
    participant M as Motor
    participant R as Repos

    F->>H: POST /api/demo/tiempo {avanzar_hasta}
    H->>CO: TickReloj(fecha_simulada)
    CO->>NU: delegar
    NU->>UC: ejecutar(fecha)
    UC->>R: hitos PENDIENTES ≤ fecha de planes ACTIVOS
    loop por hito vencido
        UC->>M: recalcular brecha
        alt brecha cerrada / afiliación cumplida
            UC->>CO: re-disparar calificación (EN_NUTRICION→PERFILANDO)
        else continúa el plan
            UC->>NU: redactar mensaje (skill redaccion-con-dignidad)
            NU->>R: persistir Mensaje VIVI (HITO_NUTRICION)
        end
        UC->>R: hito → NOTIFICADO
    end
    H-->>F: 200 {fecha_simulada, hitos_disparados}
```
