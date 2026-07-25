## Exploration: Domain Contract Types (Issue #6)

### Current State

The `internal/domain/` package currently contains 4 source files and 1 test file:

| File | Contents | Purpose |
|------|----------|---------|
| `doc.go` | Package comment only | Documents domain import rule (NFR-M-01) |
| `comprador.go` | `Comprador` struct (15 fields) | Historical buyers for kNN twin matching |
| `proyecto.go` | `Proyecto` struct (8 fields) | Project catalog (16 projects with brochure+360) |
| `constantes.go` | `Constantes`, `TramoSubsidio`, `EventoCalendario` structs + loaders | Economic parameters (Contract §4.5) |
| `constantes_test.go` | Unit tests for loaders | Validates JSON parsing |

**Key observation:** `comprador.go` and `proyecto.go` were created during issue #1 (repo foundation) as lightweight data-transfer structs for the pipeline (`cmd/pipeline/main.go`). They are imported and populated by `internal/pipeline/compradores.go` and `internal/pipeline/proyectos.go`.

Issue #6 (Paso 3 — `capacidad.go`) re-declares **both** `Comprador` and `Proyecto` with the full Contract §4.1/§4.2 schema. This creates a **duplicate declaration conflict** that Go will not compile.

### Reconciliation Analysis

#### `Comprador` — field-by-field comparison

| Field | Existing (`comprador.go`) | Target (`capacidad.go` Paso 3) | Difference |
|-------|---------------------------|-------------------------------|------------|
| ID | `int` / `json:"id"` | `int` / `json:"id"` | None |
| Proyecto | `string` / `json:"proyecto"` | `string` / `json:"proyecto"` | None |
| ProyectoID | `string` / `json:"proyecto_id"` | `string` / `json:"proyecto_id"` | None |
| Etapa | `string` / `json:"etapa"` | `string` / `json:"etapa"` | None |
| FechaOpcion | `string` / `json:"fecha_opcion"` | `string` / `json:"fecha_opcion"` | None |
| Desistio | `bool` / `json:"desistio"` | `bool` / `json:"desistio"` | None |
| Entidad | `string` / `json:"entidad"` | `string` / `json:"entidad"` | None |
| Medio | `string` / `json:"medio"` | `string` / `json:"medio"` | None |
| ValorCOP | `int64` / `json:"valor_cop"` | `int64` / `json:"valor_cop"` | None |
| Afiliado | `bool` / `json:"afiliado"` | `bool` / `json:"afiliado"` | None |
| Segmento | `string` / `json:"segmento"` | `string` / `json:"segmento"` | None |
| Categoria | `string` / `json:"categoria"` | `string` / `json:"categoria"` | None |
| RangoEdad | `string` / `json:"rango_edad"` | `string` / `json:"rango_edad"` | None |
| PersonasACargo | `int` / `json:"personas_a_cargo"` | `int` / `json:"personas_a_cargo"` | None |
| Piramide | `string` / `json:"piramide"` | `string` / `json:"piramide"` | None |

**Result:** Field sets are **identical** in name, type, and JSON tags. Only field ordering and doc comment differ.

#### `Proyecto` — field-by-field comparison

| Field | Existing (`proyecto.go`) | Target (`capacidad.go` Paso 3) | Difference |
|-------|--------------------------|-------------------------------|------------|
| ProyectoID | `string` / `json:"proyecto_id"` | `string` / `json:"proyecto_id"` | None |
| Nombre | `string` / `json:"nombre"` | `string` / `json:"nombre"` | None |
| Zona | `string` / `json:"zona"` | `string` / `json:"zona"` | None |
| PrecioDesde | `int64` / `json:"precio_desde"` | `int64` / `json:"precio_desde"` | None |
| PrecioHasta | `int64` / `json:"precio_hasta"` | `int64` / `json:"precio_hasta"` | None |
| EsVIS | `bool` / `json:"es_vis"` | `bool` / `json:"es_vis"` | None |
| BrochureURL | `string` / `json:"brochure_url"` | `string` / `json:"brochure_url"` | None |
| Recorrido360URL | `string` / `json:"recorrido_360_url"` | `string` / `json:"recorrido_360_url"` | None |

**Result:** Structs are **identical** in all respects.

### Affected Areas

- `internal/domain/comprador.go` — **DELETE**: this file's `Comprador` struct moves into `capacidad.go` (per issue #6 Paso 3)
- `internal/domain/proyecto.go` — **DELETE**: this file's `Proyecto` struct moves into `capacidad.go` (per issue #6 Paso 3)
- `internal/domain/capacidad.go` — **CREATE**: new file per Paso 3, contains `Comprador`, `Proyecto`, `Capacidad`, `Intencion`, `Recomendacion`, `ItemDesglose`
- `internal/domain/enums.go` — **CREATE**: all Contract §1 enums (Paso 1)
- `internal/domain/perfil.go` — **CREATE**: `Perfil`, `CampoPerfil`, accessors (Paso 2)
- `internal/domain/lead.go` — **CREATE**: `Lead`, `Mensaje` (Paso 4)
- `internal/domain/plan.go` — **CREATE**: `PlanNutricion`, `Hito` (Paso 5)
- `internal/domain/ficha.go` — **CREATE**: `Ficha`, `AlertaDesistimiento`, `Identificacion` (Paso 6)
- `internal/domain/perfil_test.go` — **CREATE**: tests for `Perfil` accessors (Paso 7)
- `internal/pipeline/compradores.go` — **NO CHANGE** (imports `domain.Comprador`, which remains same API)
- `internal/pipeline/proyectos.go` — **NO CHANGE** (imports `domain.Proyecto`, which remains same API)
- `internal/pipeline/compradores_test.go` — **NO CHANGE**
- `internal/pipeline/proyectos_test.go` — **NO CHANGE**
- `cmd/pipeline/main.go` — **NO CHANGE** (uses pipeline functions, not domain types directly)

### Approaches

1. **Delete-and-Replace (recommended)** — Delete `comprador.go` and `proyecto.go`; issue #6 code in `capacidad.go` becomes the canonical declaration
   - Pros: Matches issue spec exactly; single source of truth; field ordering follows Contract grouping convention; no dead files
   - Cons: Git blame history for those types moves; requires confirming pipeline still compiles (trivial — same API surface)
   - Effort: Low

2. **Keep existing files, consolidate references** — Keep `comprador.go` and `proyecto.go` as-is; remove the `Comprador`/`Proyecto` declarations from `capacidad.go`
   - Pros: No file deletions; preserves git blame
   - Cons: **Violates the issue spec** which places them in `capacidad.go` with grouped domain context; doc comments in existing files are less informative than the target
   - Effort: Low

3. **Type alias bridging** — Keep old files with `type Comprador = capacidad.Comprador` aliases
   - Pros: No import breakage at all
   - Cons: Unnecessary complexity; Go doesn't need this (same package); violates "no external imports" since they'd reference each other within the same package
   - Effort: Medium (pointless complexity)

### Recommendation

**Approach 1 (Delete-and-Replace)** is the correct path:

1. Delete `internal/domain/comprador.go` and `internal/domain/proyecto.go`
2. Create the 6 new files exactly as specified in issue #6 (Pasos 1–6)
3. Create the test file (Paso 7)
4. Verify compilation and all downstream consumers (pipeline package) still work — they will, because the exported API surface (`domain.Comprador` and `domain.Proyecto`) is structurally identical

**Why this is safe:**
- Same package name (`domain`)
- Identical field names, types, and JSON tags
- Only differences: field order (irrelevant to Go), and doc comments (improved)
- All importers (`internal/pipeline`) reference by field name, never by positional initialization
- Confirmed: pipeline uses `c.Proyecto`, `c.ProyectoID`, `c.ValorCOP`, etc. — named field access only

### Risks

- **Low:** If any code elsewhere uses positional struct initialization (`domain.Comprador{1, "x", ...}`) it would break with field reordering. **Verified not present** — all usage is named-field (`c.Proyecto = ...`, `domain.Proyecto{ProyectoID: ...}`).
- **Low:** `constantes.go` already occupies the domain package and has its own types. No naming conflicts with issue #6 types (verified: no overlap in struct names or exported constants).
- **None:** NFR-M-01 compliance — all new files are pure domain (no imports outside stdlib `time`, `encoding/json` is only in `constantes.go`).

### Ready for Proposal

Yes — the reconciliation approach is unambiguous. The orchestrator should proceed to `sdd-propose`. Key decision: delete the two legacy files and replace with the issue #6 specification verbatim. Zero API breakage.
