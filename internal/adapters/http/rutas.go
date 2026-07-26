package http

import (
	"context"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
	"net/http"
	"time"
)

type Perfilador interface {
	Ejecutar(context.Context, usecase.EntradaPerfilar) (usecase.SalidaPerfilar, error)
}
type Dependencias struct {
	Perfilar      Perfilador
	Leads         usecase.LeadRepository
	Fichas        usecase.FichaRepository
	Planes        usecase.PlanRepository
	Catalogo      usecase.CatalogoRepository
	Turnos        TurnoTracker
	Demo          usecase.DemoRepository
	Reloj         usecase.Reloj
	AvanzarDemo   *usecase.AvanzarDemo
	ReiniciarDemo *usecase.ReiniciarDemo
}
type Controlador struct {
	perfilador    Perfilador
	leads         usecase.LeadRepository
	consulta      *usecase.ConsultarLeads
	buyerPersona  *usecase.BuyerPersona
	turnos        TurnoTracker
	avanzarDemo   *usecase.AvanzarDemo
	reiniciarDemo *usecase.ReiniciarDemo
}

func NuevoControlador(d Dependencias) (*Controlador, error) {
	if d.Perfilar == nil || d.Leads == nil {
		return nil, usecase.ErrValidacion
	}
	return &Controlador{
		perfilador:    d.Perfilar,
		leads:         d.Leads,
		consulta:      &usecase.ConsultarLeads{Leads: d.Leads, Fichas: d.Fichas, Planes: d.Planes},
		buyerPersona:  &usecase.BuyerPersona{Catalogo: d.Catalogo},
		turnos:        d.Turnos,
		avanzarDemo:   demoUseCase(d),
		reiniciarDemo: d.ReiniciarDemo,
	}, nil
}
func (c *Controlador) Registrar(mux *http.ServeMux) {
	if mux == nil {
		return
	}
	mux.HandleFunc("POST /api/leads", c.crearLead)
	mux.HandleFunc("GET /api/leads", c.cola)
	mux.HandleFunc("GET /api/leads/{lead_id}", c.detalleLead)
	mux.HandleFunc("GET /api/leads/{lead_id}/ficha", c.ficha)
	mux.HandleFunc("POST /api/leads/{lead_id}/mensajes", c.mensaje)
	mux.HandleFunc("GET /api/leads/{lead_id}/conversacion", c.conversacion)
	mux.HandleFunc("GET /api/gerencia/buyer-persona", c.gerenciaBuyerPersona)
	mux.HandleFunc("POST /api/demo/tiempo", c.avanzarTiempo)
	mux.HandleFunc("POST /api/demo/reiniciar", c.reiniciar)
	mux.HandleFunc("/api/", c.apiNoEncontrada)
}

func demoUseCase(d Dependencias) *usecase.AvanzarDemo {
	if d.AvanzarDemo != nil {
		return d.AvanzarDemo
	}
	if d.Demo == nil || d.Reloj == nil {
		return nil
	}
	return &usecase.AvanzarDemo{Demo: d.Demo, Reloj: d.Reloj}
}
func (c *Controlador) apiNoEncontrada(w http.ResponseWriter, _ *http.Request) {
	writeError(w, &usecase.NotFoundError{Resource: "lead"})
}

type relojSistema struct{}

// relojSistema is a non-persistent real wall clock; Avanzar is intentionally a no-op.

func (relojSistema) Ahora() time.Time         { return time.Now().UTC() }
func (relojSistema) FechaSimulada() time.Time { return time.Now().UTC() }
func (relojSistema) Avanzar(time.Time)        {}
func RelojSistema() usecase.Reloj             { return relojSistema{} }
