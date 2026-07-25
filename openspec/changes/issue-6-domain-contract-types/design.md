# Design: Domain Contract Types (Issue #6)

## Technical Approach

Implement the corrected issue #6/spec vocabulary as pure Go domain declarations. Create six production files and one standard-library test file; delete the two legacy declaration files. `Comprador` and `Proyecto` move field-identically into `capacidad.go`, preserving `domain.Comprador`/`domain.Proyecto` and the pipeline API without pipeline edits.

## Architecture Decisions

| Decision | Rejected alternative | Rationale |
|---|---|---|
| Copy the 11 specified typed-string enum groups and exact structs/tags | Additional enum groups or Ficha fields | Contract wire values and corrected spec are authoritative. |
| `Perfil` is `map[string]CampoPerfil`; `ActualizadoEn` is `time.Time` | `map[CampoPerfil]any`, pointer timestamp, flat struct | Matches the exact dynamic JSON object and required timestamp. |
| Expose only `Entero`, `Booleano`, `Texto`, `EsVerificado` | Any additional profile methods | Issue #6 defines the complete method surface. |
| Delete legacy files and relocate identical structs | Duplicate declarations or aliases | One declaration compiles while keeping consumer-visible names, fields, types, and tags unchanged. |
| Keep domain production code stdlib-only | HTTP/database/LLM or outer-layer imports | Enforces the domain dependency boundary. |

## Data Flow

```text
XLSX/mapa -> internal/pipeline -> domain.Comprador/domain.Proyecto -> compradores.json/proyectos.json
                                      |
Contract enums + Perfil -> Lead/Plan/Capacidad -> Ficha JSON
```

## File Changes

| File | Action | Exact responsibility |
|---|---|---|
| `internal/domain/enums.go` | Create | Eleven enum groups and literals below. |
| `internal/domain/perfil.go` | Create | Profile keys, sets, value type, and four methods. |
| `internal/domain/capacidad.go` | Create | Capacity/recommendation types plus relocated `Comprador`/`Proyecto`. |
| `internal/domain/lead.go` | Create | `Lead`, `Mensaje`. |
| `internal/domain/plan.go` | Create | `Hito`, `PlanNutricion`. |
| `internal/domain/ficha.go` | Create | `AlertaDesistimiento`, `Identificacion`, `Ficha`. |
| `internal/domain/perfil_test.go` | Create | Table-driven contract and accessor tests using `testing`. |
| `internal/domain/comprador.go` | Delete | Declaration relocates field-identically. |
| `internal/domain/proyecto.go` | Delete | Declaration relocates field-identically. |

`internal/pipeline/{compradores,proyectos}.go`, their tests, `cmd/pipeline/main.go`, and existing domain constants remain unchanged.

## Interfaces / Contracts

Enums: `EstadoLead=NUEVO|PERFILANDO|CALIFICADO|ENTREGADO|EN_NUTRICION|PAUSADO|REMARKETING|DESPEDIDO|CERRADO`; `Ruta=ASESOR|NUTRICION|REMARKETING|DESPEDIDA`; `FuenteCampo=VERIFICADO_BASE|DECLARADO|INFERIDO`; `Nivel=ALTA|MEDIA|BAJA`; `TipoMensaje=TEXTO|AUDIO`; `TipoContenido=TEXTO|TARJETAS_PROYECTOS|OFERTA_PLAN|HITO_NUTRICION|SISTEMA`; `AutorMensaje=LEAD|VIVI|SISTEMA`; `EstadoPlan=ACTIVO|PAUSADO|COMPLETADO|CANCELADO`; `TipoHito=AFILIACION|AHORRO|CESANTIAS|PRIMA|REEVALUACION`; `EstadoHito=PENDIENTE|NOTIFICADO|CUMPLIDO|OMITIDO`; `Semaforo=VERDE|AMBAR|GRIS`.

`CampoPerfil{Valor:any/valor,Fuente:FuenteCampo/fuente,Confianza:float64/confianza,RequiereConfirmacion:bool/requiere_confirmacion,ActualizadoEn:time.Time/actualizado_en}`. `Perfil map[string]CampoPerfil` has only `Entero(string)(int64,bool)` (`int64`, `int`, JSON `float64`), `Booleano`, `Texto`, and `EsVerificado`.

`CamposReconocidos map[string]bool` has exactly: `ingreso_hogar, recursos_propios, personas_hogar, tipo_hogar, edad, zona_deseada, plazo_compra_meses, arriendo_actual, tiene_vivienda, recibio_subsidio, reporte_credito, situacion_laboral, caja_externa, hogar_con_afiliado, cedula_familiar_afiliado, cesantias_entidad, categoria, segmento`. `CamposCriticos` is exactly `ingreso_hogar, recursos_propios, tiene_vivienda, recibio_subsidio`.

Exact layouts (`Field:Type/json`):

```text
ItemDesglose{Concepto:string/concepto,Monto:int64/monto,Regla:string/regla,Fuente:FuenteCampo/fuente}
Capacidad{PresupuestoMax:int64/presupuesto_max,CreditoMax:int64/credito_max,SubsidioAplicable:int64/subsidio_aplicable,RecursosPropios:int64/recursos_propios,Ratio:float64/ratio,Confianza:float64/confianza,Desglose:[]ItemDesglose/desglose}
Intencion{Nivel:Nivel/nivel,Confianza:Nivel/confianza,Senales:[]string/senales}
Recomendacion{ProyectoID:string/proyecto_id,Nombre:string/nombre,Zona:string/zona,PrecioDesde:int64/precio_desde,Razon:string/razon,Vecinos:int/vecinos,TasaDesistimiento:float64/tasa_desistimiento,BrochureURL:string/brochure_url,Recorrido360URL:string/recorrido_360_url}
Proyecto{ProyectoID:string/proyecto_id,Nombre:string/nombre,Zona:string/zona,PrecioDesde:int64/precio_desde,PrecioHasta:int64/precio_hasta,EsVIS:bool/es_vis,BrochureURL:string/brochure_url,Recorrido360URL:string/recorrido_360_url}
Comprador{ID:int/id,ProyectoID:string/proyecto_id,Proyecto:string/proyecto,Etapa:string/etapa,Afiliado:bool/afiliado,Categoria:string/categoria,Segmento:string/segmento,RangoEdad:string/rango_edad,PersonasACargo:int/personas_a_cargo,Piramide:string/piramide,ValorCOP:int64/valor_cop,Entidad:string/entidad,Medio:string/medio,Desistio:bool/desistio,FechaOpcion:string/fecha_opcion}
Lead{LeadID:string/lead_id,Nombre:string/nombre,Telefono:string/telefono,Cedula:string/cedula,Fuente:string/fuente,Estado:EstadoLead/estado,Ruta:Ruta/ruta,Afiliado:bool/afiliado,Prioridad:float64/prioridad,ConsumeCupo10:bool/consume_cupo_10,Perfil:Perfil/perfil,Capacidad:*Capacidad/capacidad,Intencion:*Intencion/intencion,Version:int/-,CreadoEn:time.Time/creado_en,ActualizadoEn:time.Time/actualizado_en}
Mensaje{MensajeID:string/mensaje_id,LeadID:string/-,Autor:AutorMensaje/autor,TipoContenido:TipoContenido/tipo_contenido,Texto:string/texto,CreadoEn:time.Time/creado_en,Adjunto:map[string]any/adjunto}
Hito{HitoID:string/hito_id,Tipo:TipoHito/tipo,Fecha:string/fecha,Monto:*int64/monto,Descripcion:string/descripcion,Estado:EstadoHito/estado}
PlanNutricion{PlanID:string/plan_id,LeadID:string/lead_id,Estado:EstadoPlan/estado,ConsentimientoEn:*time.Time/consentimiento_en,Frecuencia:string/frecuencia,MetaMonto:int64/meta_monto,MetaDescripcion:string/meta_descripcion,Hitos:[]Hito/hitos}
AlertaDesistimiento{Activa:bool/activa,TasaVecinos:float64/tasa_vecinos,Detalle:*string/detalle}
Identificacion{Nombre:string/nombre,Afiliada:bool/afiliada,Categoria:string/categoria,Telefono:string/telefono}
Ficha{FichaID:string/ficha_id,LeadID:string/lead_id,GeneradaEn:time.Time/generada_en,ConfianzaPerfil:float64/confianza_perfil,BandaAdvertencia:*string/banda_advertencia,Identificacion:Identificacion/identificacion,Capacidad:Capacidad/capacidad,Perfil:Perfil/perfil,Intencion:Intencion/intencion,Recomendaciones:[]Recomendacion/recomendaciones,Beneficios:[]string/beneficios,ArgumentosVenta:[]string/argumentos_venta,AlertaDesistimiento:AlertaDesistimiento/alerta_desistimiento,ConsumeCupo10:bool/consume_cupo_10}
```

## Testing Strategy

Use table-driven `testing` cases for accessor conversions/failures, verification, exact 18/four sets, enum JSON literals, and reflected field types/tags/nullability. Run `go test ./internal/domain/... -v`, existing `go test ./internal/pipeline/... -v`, `go build ./internal/domain/...`, `go build ./...`, and `go vet ./internal/domain/...`. Guard isolation with `go list -deps ./internal/domain/... | grep -E 'internal/(usecase|adapters|infrastructure)'`, expecting no output.

## Threat Matrix

N/A — no routing, shell, subprocess, VCS/PR automation, executable-file classification, or process-integration boundary.

## Migration / Rollout

No migration required; no data, schema, API route, or feature-flag transition. Rollback restores the deleted files and removes the new declarations.

## Open Questions

None.
