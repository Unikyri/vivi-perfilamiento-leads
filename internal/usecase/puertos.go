package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
)

var (
	ErrNoEncontrado      = errors.New("recurso no encontrado")
	ErrNotFound          = ErrNoEncontrado
	ErrOptimisticLock    = errors.New("conflicto de version")
	ErrFichaNoDisponible = errors.New("ficha no disponible")
)

// NotFoundError adds the resource identity while preserving errors.Is matching.
type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s %q no encontrado", e.Resource, e.ID)
}

func (e *NotFoundError) Unwrap() error { return ErrNoEncontrado }

// FiltroLeads contains the optional lead-queue filters.
type FiltroLeads struct {
	Afiliado *bool
	Ruta     *domain.Ruta
}

// Afiliado is an affiliate lookup result from the data pipeline.
type Afiliado struct {
	Cedula          string `json:"cedula"`
	Nombre          string `json:"nombre"`
	Categoria       string `json:"categoria"`
	Segmento        string `json:"segmento"`
	IngresoMensual  int64  `json:"ingreso_mensual"`
	PersonasACargo  int    `json:"personas_a_cargo"`
	TipoHogar       string `json:"tipo_hogar"`
	EmpresaPiramide string `json:"empresa_piramide"`
	AfiliadoActivo  bool   `json:"afiliado_activo"`
}

// EntradaTurno is the provider-neutral input for one LLM turn.
type EntradaTurno struct {
	LeadID            string            `json:"lead_id"`
	Nombre            string            `json:"nombre"`
	MensajeUsuario    string            `json:"mensaje_usuario"`
	Perfil            domain.Perfil     `json:"perfil"`
	Capacidad         *domain.Capacidad `json:"capacidad"`
	NumerosDelMotor   map[string]int64  `json:"numeros_del_motor"`
	HistorialReciente []domain.Mensaje  `json:"historial_reciente"`
	EsAfiliado        bool              `json:"es_afiliado"`
}

// CampoExtraido is one Contract v1.1 §7 profile extraction.
type CampoExtraido struct {
	Campo                string             `json:"campo"`
	Valor                any                `json:"valor"`
	Fuente               domain.FuenteCampo `json:"fuente"`
	Confianza            float64            `json:"confianza"`
	RequiereConfirmacion bool               `json:"requiere_confirmacion"`
}

// SalidaTurno is the provider-neutral structured response for one turn.

// Actions valid in SalidaTurno.Accion (Contract v1.1 §7).
const (
	AccionContinuar        = "CONTINUAR"
	AccionPerfilCompleto   = "PERFIL_COMPLETO"
	AccionConsentimientoSi = "CONSENTIMIENTO_SI"
	AccionConsentimientoNo = "CONSENTIMIENTO_NO"
	AccionPausarContacto   = "PAUSAR_CONTACTO"
	AccionFueraDeDominio   = "FUERA_DE_DOMINIO"
	AccionAudioInint       = "AUDIO_ININTELIGIBLE"
)

type SalidaTurno struct {
	CamposExtraidos []CampoExtraido  `json:"campos_extraidos"`
	Intencion       domain.Intencion `json:"intencion"`
	Respuesta       string           `json:"respuesta"`
	Accion          string           `json:"accion"`
}

// Audio is an incoming voice note.
type Audio struct {
	Base64    string `json:"base64"`
	MIME      string `json:"mime"`
	DuracionS int    `json:"duracion_s"`
}

// HitoConPlan identifies an overdue milestone and its owning plan and lead.
type HitoConPlan struct {
	Hito   domain.Hito `json:"hito"`
	PlanID string      `json:"plan_id"`
	LeadID string      `json:"lead_id"`
}

// Evento is an internal domain event.
type Evento struct {
	Tipo    string         `json:"tipo"`
	LeadID  string         `json:"lead_id"`
	Payload map[string]any `json:"payload"`
}

// Event type constants defined by Contract v1.1 §6.
const (
	EvLeadNuevo       = "LeadNuevo"
	EvMensajeEntrante = "MensajeEntrante"
	EvPerfilCompleto  = "PerfilCompleto"
	EvRutaDecidida    = "RutaDecidida"
	EvHitoVencido     = "HitoVencido"
	EvTickReloj       = "TickReloj"
	EvResultadoReal   = "ResultadoReal"
)

// LeadRepository persists leads and their messages.
type LeadRepository interface {
	Crear(context.Context, *domain.Lead) error
	PorID(context.Context, string) (*domain.Lead, error)
	Guardar(context.Context, *domain.Lead) error
	Listar(context.Context, FiltroLeads) ([]*domain.Lead, error)
	AgregarMensaje(context.Context, *domain.Mensaje) error
	Conversacion(context.Context, string) ([]domain.Mensaje, error)
}

// PlanRepository persists nutrition plans and milestones.
type PlanRepository interface {
	Crear(context.Context, *domain.PlanNutricion) error
	PorLead(context.Context, string) (*domain.PlanNutricion, error)
	Guardar(context.Context, *domain.PlanNutricion) error
	HitosVencidos(context.Context, time.Time) ([]HitoConPlan, error)
	MarcarHito(context.Context, string, domain.EstadoHito) error
}

// FichaRepository persists advisor-facing lead fichas.
type FichaRepository interface {
	Guardar(context.Context, *domain.Ficha) error
	PorLead(context.Context, string) (*domain.Ficha, error)
}

// CatalogoRepository exposes read-only data-pipeline outputs.
type CatalogoRepository interface {
	Proyectos(context.Context) (map[string]domain.Proyecto, error)
	Compradores(context.Context) ([]domain.Comprador, error)
	AfiliadoPorCedula(context.Context, string) (*Afiliado, error)
	BrochureMarkdown(context.Context, string) (string, error)
}

// LLMProvider abstracts the model provider. One call is made per turn.
type LLMProvider interface {
	GenerarTurno(context.Context, EntradaTurno) (SalidaTurno, error)
	ProcesarAudio(context.Context, Audio, EntradaTurno) (SalidaTurno, error)
	Nombre() string
}

// MensajeriaGateway is the minimal NFR-M-01 outbound messaging port.
type MensajeriaGateway interface {
	Enviar(context.Context, *domain.Mensaje) error
}

// Reloj provides real and simulated time for deterministic demo behavior.
type Reloj interface {
	Ahora() time.Time
	FechaSimulada() time.Time
	Avanzar(time.Time)
}

// BusEventos publishes and subscribes to internal domain events.
type BusEventos interface {
	Publicar(context.Context, Evento)
	Suscribir(string, func(context.Context, Evento))
}

// GeneradorID creates opaque identifiers owned by the application boundary.
type GeneradorID interface {
	Nuevo() string
}
