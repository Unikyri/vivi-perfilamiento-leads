# Domain Contract Types Specification

## Purpose
The domain package SHALL provide the Contract v1.1 shared type vocabulary without external-layer dependencies.

## Requirements

### Requirement: Contract Enums
The package MUST export 11 typed `string` enums with exactly these literals: `EstadoLead=NUEVO|PERFILANDO|CALIFICADO|ENTREGADO|EN_NUTRICION|PAUSADO|REMARKETING|DESPEDIDO|CERRADO`; `Ruta=ASESOR|NUTRICION|REMARKETING|DESPEDIDA`; `FuenteCampo=VERIFICADO_BASE|DECLARADO|INFERIDO`; `Nivel=ALTA|MEDIA|BAJA`; `TipoMensaje=TEXTO|AUDIO`; `TipoContenido=TEXTO|TARJETAS_PROYECTOS|OFERTA_PLAN|HITO_NUTRICION|SISTEMA`; `AutorMensaje=LEAD|VIVI|SISTEMA`; `EstadoPlan=ACTIVO|PAUSADO|COMPLETADO|CANCELADO`; `TipoHito=AFILIACION|AHORRO|CESANTIAS|PRIMA|REEVALUACION`; `EstadoHito=PENDIENTE|NOTIFICADO|CUMPLIDO|OMITIDO`; `Semaforo=VERDE|AMBAR|GRIS`.

#### Scenario: Enum wire values
- GIVEN each enum literal
- WHEN JSON-marshalled
- THEN its exact listed literal is emitted

### Requirement: Perfil Schema and Accessors
`CampoPerfil` MUST expose `Valor any` (`valor`), `Fuente FuenteCampo` (`fuente`), `Confianza float64` (`confianza`), `RequiereConfirmacion bool` (`requiere_confirmacion`), and `ActualizadoEn time.Time` (`actualizado_en`). `Perfil` MUST be `map[string]CampoPerfil` and expose `Entero(clave string) (int64, bool)`, accepting `int64`, `int`, and JSON `float64`; `Booleano(clave string) (bool, bool)`; `Texto(clave string) (string, bool)`; and `EsVerificado(clave string) bool`.

`CamposReconocidos map[string]bool` MUST contain exactly `ingreso_hogar,recursos_propios,personas_hogar,tipo_hogar,edad,zona_deseada,plazo_compra_meses,arriendo_actual,tiene_vivienda,recibio_subsidio,reporte_credito,situacion_laboral,caja_externa,hogar_con_afiliado,cedula_familiar_afiliado,cesantias_entidad,categoria,segmento`. `CamposCriticos` MUST contain exactly `ingreso_hogar,recursos_propios,tiene_vivienda,recibio_subsidio`.

#### Scenario: Accessor conversion
- GIVEN `int64`, `int`, `float64`, boolean, string, absent, and incompatible values
- WHEN the matching accessor is called
- THEN valid values return `ok=true`; all others return zero and `false`

#### Scenario: Verification and key sets
- GIVEN a `VERIFICADO_BASE` field and exported sets
- WHEN inspected
- THEN only that field verifies, the map has 18 keys, and the critical set has the stated four

### Requirement: Capacity and Recommendation Types
For this requirement, `Field:Type/json_name` specifies the exact Go field, type, and JSON field. The package MUST export `ItemDesglose{Concepto:string/concepto,Monto:int64/monto,Regla:string/regla,Fuente:FuenteCampo/fuente}`; `Capacidad{PresupuestoMax:int64/presupuesto_max,CreditoMax:int64/credito_max,SubsidioAplicable:int64/subsidio_aplicable,RecursosPropios:int64/recursos_propios,Ratio:float64/ratio,Confianza:float64/confianza,Desglose:[]ItemDesglose/desglose}`; `Intencion{Nivel:Nivel/nivel,Confianza:Nivel/confianza,Senales:[]string/senales}`; `Recomendacion{ProyectoID:string/proyecto_id,Nombre:string/nombre,Zona:string/zona,PrecioDesde:int64/precio_desde,Razon:string/razon,Vecinos:int/vecinos,TasaDesistimiento:float64/tasa_desistimiento,BrochureURL:string/brochure_url,Recorrido360URL:string/recorrido_360_url}`; `Proyecto{ProyectoID:string/proyecto_id,Nombre:string/nombre,Zona:string/zona,PrecioDesde:int64/precio_desde,PrecioHasta:int64/precio_hasta,EsVIS:bool/es_vis,BrochureURL:string/brochure_url,Recorrido360URL:string/recorrido_360_url}`; and `Comprador{ID:int/id,ProyectoID:string/proyecto_id,Proyecto:string/proyecto,Etapa:string/etapa,Afiliado:bool/afiliado,Categoria:string/categoria,Segmento:string/segmento,RangoEdad:string/rango_edad,PersonasACargo:int/personas_a_cargo,Piramide:string/piramide,ValorCOP:int64/valor_cop,Entidad:string/entidad,Medio:string/medio,Desistio:bool/desistio,FechaOpcion:string/fecha_opcion}`.

#### Scenario: Capacity schema
- GIVEN all listed types
- WHEN types and JSON fields are reflected
- THEN each exactly matches this requirement

### Requirement: Lead and Message Types
The package MUST export `Lead{LeadID:string/lead_id,Nombre:string/nombre,Telefono:string/telefono,Cedula:string/cedula,Fuente:string/fuente,Estado:EstadoLead/estado,Ruta:Ruta/ruta,Afiliado:bool/afiliado,Prioridad:float64/prioridad,ConsumeCupo10:bool/consume_cupo_10,Perfil:Perfil/perfil,Capacidad:*Capacidad/capacidad,Intencion:*Intencion/intencion,Version:int/-,CreadoEn:time.Time/creado_en,ActualizadoEn:time.Time/actualizado_en}` and `Mensaje{MensajeID:string/mensaje_id,LeadID:string/-,Autor:AutorMensaje/autor,TipoContenido:TipoContenido/tipo_contenido,Texto:string/texto,CreadoEn:time.Time/creado_en,Adjunto:map[string]any/adjunto}`.

#### Scenario: Lead/message schema
- GIVEN both structs
- WHEN their metadata is inspected
- THEN every field and omitted (`-`) JSON field matches

### Requirement: Nutrition and Ficha Types
The package MUST export `Hito{HitoID:string/hito_id,Tipo:TipoHito/tipo,Fecha:string/fecha,Monto:*int64/monto,Descripcion:string/descripcion,Estado:EstadoHito/estado}`; `PlanNutricion{PlanID:string/plan_id,LeadID:string/lead_id,Estado:EstadoPlan/estado,ConsentimientoEn:*time.Time/consentimiento_en,Frecuencia:string/frecuencia,MetaMonto:int64/meta_monto,MetaDescripcion:string/meta_descripcion,Hitos:[]Hito/hitos}`; `AlertaDesistimiento{Activa:bool/activa,TasaVecinos:float64/tasa_vecinos,Detalle:*string/detalle}`; `Identificacion{Nombre:string/nombre,Afiliada:bool/afiliada,Categoria:string/categoria,Telefono:string/telefono}`; and `Ficha{FichaID:string/ficha_id,LeadID:string/lead_id,GeneradaEn:time.Time/generada_en,ConfianzaPerfil:float64/confianza_perfil,BandaAdvertencia:*string/banda_advertencia,Identificacion:Identificacion/identificacion,Capacidad:Capacidad/capacidad,Perfil:Perfil/perfil,Intencion:Intencion/intencion,Recomendaciones:[]Recomendacion/recomendaciones,Beneficios:[]string/beneficios,ArgumentosVenta:[]string/argumentos_venta,AlertaDesistimiento:AlertaDesistimiento/alerta_desistimiento,ConsumeCupo10:bool/consume_cupo_10}`.

#### Scenario: Nutrition/ficha schema
- GIVEN populated and nil optional fields
- WHEN JSON-marshalled
- THEN field names and nullability match the Contract

### Requirement: Compatibility, Isolation, and Evidence
Identical public `Comprador` and `Proyecto` structs MUST relocate to `capacidad.go` without API changes. Domain production code MUST NOT depend on `internal/usecase`, `internal/adapters`, `internal/infrastructure`, HTTP, databases, or LLMs. Focused runtime evidence MUST pass `go build ./internal/domain/...`, `go test ./internal/domain/... -v`, `go vet ./internal/domain/...`, and a dependency check with no forbidden internal layer.

#### Scenario: Relocation and isolation
- GIVEN the relocated declarations
- WHEN the focused commands run
- THEN they pass and report no forbidden dependency
