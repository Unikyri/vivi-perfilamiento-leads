-- 001_esquema_inicial.sql
-- Contract §5: Esquema de base de datos — Vivi Perfilamiento Leads
-- Idempotent: all statements use IF NOT EXISTS.

CREATE TABLE IF NOT EXISTS leads (
    lead_id         TEXT PRIMARY KEY,
    nombre          TEXT,
    telefono        TEXT,
    cedula          TEXT,
    fuente          TEXT,
    estado          TEXT,
    ruta            TEXT,
    afiliado        BOOLEAN,
    prioridad       REAL,
    consume_cupo_10 BOOLEAN,
    perfil          JSONB,
    capacidad       JSONB,
    intencion       JSONB,
    version         INT NOT NULL DEFAULT 1,
    creado_en       TIMESTAMPTZ,
    actualizado_en  TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS mensajes (
    mensaje_id     TEXT PRIMARY KEY,
    lead_id        TEXT REFERENCES leads(lead_id),
    autor          TEXT,
    tipo_contenido TEXT,
    texto          TEXT,
    adjunto        JSONB,
    creado_en      TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS planes (
    plan_id           TEXT PRIMARY KEY,
    lead_id           TEXT REFERENCES leads(lead_id),
    estado            TEXT,
    frecuencia        TEXT,
    consentimiento_en TIMESTAMPTZ,
    meta_monto        BIGINT,
    meta_descripcion  TEXT
);

CREATE TABLE IF NOT EXISTS hitos (
    hito_id     TEXT PRIMARY KEY,
    plan_id     TEXT REFERENCES planes(plan_id),
    tipo        TEXT,
    fecha       DATE,
    monto       BIGINT,
    descripcion TEXT,
    estado      TEXT
);

CREATE TABLE IF NOT EXISTS fichas (
    ficha_id    TEXT PRIMARY KEY,
    lead_id     TEXT REFERENCES leads(lead_id) UNIQUE,
    contenido   JSONB,
    generada_en TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS compradores (
    id              INTEGER PRIMARY KEY,
    proyecto_id     TEXT,
    proyecto        TEXT,
    etapa           TEXT,
    afiliado        BOOLEAN,
    categoria       TEXT,
    segmento        TEXT,
    rango_edad      TEXT,
    personas_a_cargo INTEGER,
    piramide        TEXT,
    valor_cop       BIGINT,
    entidad         TEXT,
    medio           TEXT,
    desistio        BOOLEAN,
    fecha_opcion    TEXT
);

CREATE TABLE IF NOT EXISTS demo (
    clave TEXT PRIMARY KEY,
    valor TEXT
);

-- Indexes (Contract §5 implicit + spec requirement: 4 indexes)
CREATE INDEX IF NOT EXISTS idx_mensajes_lead_id ON mensajes(lead_id);
CREATE INDEX IF NOT EXISTS idx_planes_lead_id ON planes(lead_id);
CREATE INDEX IF NOT EXISTS idx_hitos_plan_id ON hitos(plan_id);
CREATE INDEX IF NOT EXISTS idx_compradores_proyecto_id ON compradores(proyecto_id);
