package http

import (
	"context"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
)

type ErrTurnoEnProcesoType struct{}

func (ErrTurnoEnProcesoType) Error() string { return "turno en proceso" }

var ErrTurnoEnProceso error = ErrTurnoEnProcesoType{}

type TurnoProcesador interface {
	Ejecutar(context.Context, usecase.EntradaMensaje) error
}
type TurnoTracker interface {
	Aceptar(context.Context, usecase.EntradaMensaje) (TurnoAceptado, error)
	Activo(string) bool
}
type TurnoAceptado struct {
	MensajeID  string
	RecibidoEn time.Time
}

type EjecutorTurnos struct {
	procesador TurnoProcesador
	leads      usecase.LeadRepository
	ids        usecase.GeneradorID
	reloj      usecase.Reloj
	ctx        context.Context
	cancel     context.CancelFunc
	mu         sync.RWMutex
	activos    map[string]struct{}
	cerrado    bool
	wg         sync.WaitGroup
}

func NuevoEjecutorTurnos(p TurnoProcesador, ids usecase.GeneradorID, reloj usecase.Reloj, leads ...usecase.LeadRepository) *EjecutorTurnos {
	ctx, cancel := context.WithCancel(context.Background())
	e := &EjecutorTurnos{procesador: p, ids: ids, reloj: reloj, ctx: ctx, cancel: cancel, activos: make(map[string]struct{})}
	if len(leads) > 0 {
		e.leads = leads[0]
	}
	return e
}

func (e *EjecutorTurnos) Aceptar(ctx context.Context, entrada usecase.EntradaMensaje) (TurnoAceptado, error) {
	if err := ctx.Err(); err != nil {
		return TurnoAceptado{}, err
	}
	if err := usecase.ValidarEntradaMensaje(entrada); err != nil {
		return TurnoAceptado{}, err
	}
	if strings.TrimSpace(entrada.LeadID) == "" || e.procesador == nil || e.ids == nil || e.reloj == nil {
		return TurnoAceptado{}, usecase.ErrValidacion
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.cerrado || e.ctx.Err() != nil {
		return TurnoAceptado{}, context.Canceled
	}
	if _, ok := e.activos[entrada.LeadID]; ok {
		return TurnoAceptado{}, ErrTurnoEnProceso
	}
	if entrada.MensajeID == "" {
		entrada.MensajeID = e.ids.Nuevo()
	}
	if entrada.RecibidoEn.IsZero() {
		entrada.RecibidoEn = e.reloj.Ahora()
	}
	entrada.RecibidoEn = entrada.RecibidoEn.UTC()
	e.activos[entrada.LeadID] = struct{}{}
	e.wg.Add(1)
	go e.ejecutar(entrada)
	return TurnoAceptado{MensajeID: entrada.MensajeID, RecibidoEn: entrada.RecibidoEn}, nil
}

func (e *EjecutorTurnos) ejecutar(entrada usecase.EntradaMensaje) {
	defer e.wg.Done()
	defer func() {
		e.mu.Lock()
		delete(e.activos, entrada.LeadID)
		e.mu.Unlock()
	}()
	if err := e.procesador.Ejecutar(e.ctx, entrada); err != nil {
		log.Printf("[turnos] error procesando turno para lead %s: %v — guardando fallback", entrada.LeadID, err)
		e.guardarFallback(entrada)
	}
}

func (e *EjecutorTurnos) guardarFallback(entrada usecase.EntradaMensaje) {
	if e.leads == nil {
		return
	}
	lead, err := e.leads.PorID(e.ctx, entrada.LeadID)
	if err != nil || lead == nil {
		return
	}
	recibidoEn := entrada.RecibidoEn
	if recibidoEn.IsZero() {
		recibidoEn = e.reloj.Ahora()
	}
	msgID := entrada.MensajeID
	if msgID == "" {
		msgID = e.ids.Nuevo()
	}
	inbound := &domain.Mensaje{
		MensajeID: msgID, LeadID: lead.LeadID, Autor: domain.AutorMensajeLead,
		TipoContenido: domain.TipoContenidoTexto, Texto: entrada.Texto, CreadoEn: recibidoEn.UTC(),
	}
	if entrada.Tipo == domain.TipoMensajeAudio && entrada.Audio != nil {
		inbound.Adjunto = map[string]any{"audio_original": true}
	}
	_ = e.leads.AgregarMensaje(e.ctx, inbound)

	textoRespuesta := fallbackRespuestaTexto(entrada.Texto, lead.Afiliado)
	_ = e.leads.AgregarMensaje(e.ctx, &domain.Mensaje{
		MensajeID: e.ids.Nuevo(), LeadID: lead.LeadID, Autor: domain.AutorMensajeVivi,
		TipoContenido: domain.TipoContenidoTexto, Texto: textoRespuesta, CreadoEn: e.reloj.Ahora().UTC(),
	})
}

func (e *EjecutorTurnos) Activo(leadID string) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	_, ok := e.activos[leadID]
	return ok
}

func (e *EjecutorTurnos) Cerrar() {
	e.mu.Lock()
	if !e.cerrado {
		e.cerrado = true
		e.cancel()
	}
	e.mu.Unlock()
	e.wg.Wait()
}

// mensajeRequest is the Contract §3.2 transport shape; audio bytes never leave this adapter.
type mensajeRequest struct {
	Tipo        domain.TipoMensaje `json:"tipo"`
	Texto       string             `json:"texto"`
	AudioBase64 string             `json:"audio_base64"`
	MIME        string             `json:"mime"`
	DuracionS   int                `json:"duracion_s"`
}
type mensajeResponse struct {
	MensajeID      string    `json:"mensaje_id"`
	RecibidoEn     time.Time `json:"recibido_en"`
	TurnoEnProceso bool      `json:"turno_en_proceso"`
}

func (c *Controlador) mensaje(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("lead_id"))
	if id == "" || strings.Contains(id, "/") {
		writeError(w, usecase.ErrNoEncontrado)
		return
	}
	if c.turnos == nil {
		writeError(w, errors.New("turnos no configurados"))
		return
	}
	var req mensajeRequest
	if err := decodeJSONLimit(w, r, &req, 4<<20); err != nil {
		return
	}
	entrada := usecase.EntradaMensaje{LeadID: id, Tipo: req.Tipo, Texto: req.Texto}
	if req.Tipo == domain.TipoMensajeAudio {
		entrada.Audio = &usecase.Audio{Base64: req.AudioBase64, MIME: req.MIME, DuracionS: req.DuracionS}
	}
	if _, err := c.leads.PorID(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	accepted, err := c.turnos.Aceptar(r.Context(), entrada)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusAccepted, mensajeResponse{MensajeID: accepted.MensajeID, RecibidoEn: accepted.RecibidoEn, TurnoEnProceso: true})
}

var (
	reProyecto = regexp.MustCompile(`(?i)(proyecto|comprar|vivienda|casa|apto|apartamento)`)
	reIngreso  = regexp.MustCompile(`(?i)(ingreso|gano|salario|sueldo|\$\d)`)
	reZona     = regexp.MustCompile(`(?i)(zona|sector|barrio|soacha|kennedy|bosa|engativ)`)
)

func fallbackRespuestaTexto(texto string, afiliado bool) string {
	switch {
	case reProyecto.MatchString(texto):
		return "¡Genial! Tenemos varios proyectos que se ajustan a tu perfil. ¿Podrías contarme cuánto es tu ingreso mensual familiar para calcular tu capacidad?"
	case reIngreso.MatchString(texto):
		if afiliado {
			return "Perfecto, con eso puedo estimar tu capacidad. Como afiliada a Colsubsidio tienes acceso a subsidios especiales. ¿Tienes algún ahorro o recursos propios para la cuota inicial?"
		}
		return "Gracias por la información. ¿Tienes algún ahorro o recursos propios para la cuota inicial?"
	case reZona.MatchString(texto):
		return "Excelente ubicación. Tenemos proyectos muy buenos en esa zona. ¿Cuántas personas conforman tu hogar? Esto nos ayuda a recomendarte el tipo de vivienda ideal."
	default:
		return "¡Qué bueno saberlo! Para ayudarte mejor, ¿podrías contarme cuál es tu ingreso mensual familiar o si tienes ahorros preparados para la vivienda?"
	}
}
