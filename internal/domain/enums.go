package domain

// EstadoLead identifies a lead's lifecycle state.
type EstadoLead string

const (
	EstadoLeadNuevo       EstadoLead = "NUEVO"
	EstadoLeadPerfilando  EstadoLead = "PERFILANDO"
	EstadoLeadCalificado  EstadoLead = "CALIFICADO"
	EstadoLeadEntregado   EstadoLead = "ENTREGADO"
	EstadoLeadEnNutricion EstadoLead = "EN_NUTRICION"
	EstadoLeadPausado     EstadoLead = "PAUSADO"
	EstadoLeadRemarketing EstadoLead = "REMARKETING"
	EstadoLeadDespedido   EstadoLead = "DESPEDIDO"
	EstadoLeadCerrado     EstadoLead = "CERRADO"
)

// Ruta identifies the next route for a lead.
type Ruta string

const (
	RutaAsesor      Ruta = "ASESOR"
	RutaNutricion   Ruta = "NUTRICION"
	RutaRemarketing Ruta = "REMARKETING"
	RutaDespedida   Ruta = "DESPEDIDA"
)

// FuenteCampo identifies how a profile value was obtained.
type FuenteCampo string

const (
	FuenteCampoVerificadoBase FuenteCampo = "VERIFICADO_BASE"
	FuenteCampoDeclarado      FuenteCampo = "DECLARADO"
	FuenteCampoInferido       FuenteCampo = "INFERIDO"
)

// Nivel expresses a qualitative level.
type Nivel string

const (
	NivelAlta  Nivel = "ALTA"
	NivelMedia Nivel = "MEDIA"
	NivelBaja  Nivel = "BAJA"
)

// TipoMensaje identifies the transport type of a message.
type TipoMensaje string

const (
	TipoMensajeTexto TipoMensaje = "TEXTO"
	TipoMensajeAudio TipoMensaje = "AUDIO"
)

// TipoContenido identifies the semantic content of a message.
type TipoContenido string

const (
	TipoContenidoTexto             TipoContenido = "TEXTO"
	TipoContenidoTarjetasProyectos TipoContenido = "TARJETAS_PROYECTOS"
	TipoContenidoOfertaPlan        TipoContenido = "OFERTA_PLAN"
	TipoContenidoHitoNutricion     TipoContenido = "HITO_NUTRICION"
	TipoContenidoSistema           TipoContenido = "SISTEMA"
)

// AutorMensaje identifies who authored a message.
type AutorMensaje string

const (
	AutorMensajeLead    AutorMensaje = "LEAD"
	AutorMensajeVivi    AutorMensaje = "VIVI"
	AutorMensajeSistema AutorMensaje = "SISTEMA"
)

// EstadoPlan identifies a nutrition plan's lifecycle state.
type EstadoPlan string

const (
	EstadoPlanActivo     EstadoPlan = "ACTIVO"
	EstadoPlanPausado    EstadoPlan = "PAUSADO"
	EstadoPlanCompletado EstadoPlan = "COMPLETADO"
	EstadoPlanCancelado  EstadoPlan = "CANCELADO"
)

// TipoHito identifies a nutrition milestone.
type TipoHito string

const (
	TipoHitoAfiliacion   TipoHito = "AFILIACION"
	TipoHitoAhorro       TipoHito = "AHORRO"
	TipoHitoCesantias    TipoHito = "CESANTIAS"
	TipoHitoPrima        TipoHito = "PRIMA"
	TipoHitoReevaluacion TipoHito = "REEVALUACION"
)

// EstadoHito identifies a milestone's state.
type EstadoHito string

const (
	EstadoHitoPendiente  EstadoHito = "PENDIENTE"
	EstadoHitoNotificado EstadoHito = "NOTIFICADO"
	EstadoHitoCumplido   EstadoHito = "CUMPLIDO"
	EstadoHitoOmitido    EstadoHito = "OMITIDO"
)

// Semaforo identifies the dashboard traffic-light state.
type Semaforo string

const (
	SemaforoVerde Semaforo = "VERDE"
	SemaforoAmbar Semaforo = "AMBAR"
	SemaforoGris  Semaforo = "GRIS"
)
