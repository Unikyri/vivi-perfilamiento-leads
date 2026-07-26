package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
	"sort"
	"strings"
	"time"
)

// ColaLeads is the advisor queue read model. It never writes or recalculates a lead.
type ColaLeads struct {
	Cupo10 Cupo10       `json:"cupo_10"`
	Leads  []LeadEnCola `json:"leads"`
}
type Cupo10 struct {
	Usados            int `json:"usados"`
	PorcentajeVentana int `json:"porcentaje_ventana"`
}
type LeadEnCola struct {
	LeadID        string            `json:"lead_id"`
	Nombre        string            `json:"nombre"`
	Estado        domain.EstadoLead `json:"estado"`
	Ruta          domain.Ruta       `json:"ruta"`
	Afiliado      bool              `json:"afiliado"`
	Semaforo      domain.Semaforo   `json:"semaforo"`
	Prioridad     float64           `json:"prioridad"`
	Resumen       string            `json:"resumen"`
	ActualizadoEn time.Time         `json:"actualizado_en"`
}

// LeadDetalle contains the persisted Contract fields safe for the advisor API.
// Identity secrets (cedula and telefono) are intentionally not part of this DTO.
type LeadDetalle struct {
	LeadID        string                `json:"lead_id"`
	Nombre        string                `json:"nombre"`
	Estado        domain.EstadoLead     `json:"estado"`
	Ruta          domain.Ruta           `json:"ruta"`
	Afiliado      bool                  `json:"afiliado"`
	Semaforo      domain.Semaforo       `json:"semaforo"`
	Prioridad     float64               `json:"prioridad"`
	ConsumeCupo10 bool                  `json:"consume_cupo_10"`
	Capacidad     *domain.Capacidad     `json:"capacidad"`
	Intencion     *domain.Intencion     `json:"intencion"`
	Perfil        domain.Perfil         `json:"perfil"`
	Plan          *domain.PlanNutricion `json:"plan"`
	CreadoEn      time.Time             `json:"creado_en"`
	ActualizadoEn time.Time             `json:"actualizado_en"`
}

// ConsultarLeads implements the advisor read models over existing ports.
type ConsultarLeads struct {
	Leads  LeadRepository
	Fichas FichaRepository
	Planes PlanRepository
}

func (uc *ConsultarLeads) Ejecutar(ctx context.Context, filtro FiltroLeads) (ColaLeads, error) {
	if uc == nil || uc.Leads == nil {
		return ColaLeads{}, fmt.Errorf("%w: repositorio de leads requerido", ErrValidacion)
	}
	leads, err := uc.Leads.Listar(ctx, filtro)
	if err != nil {
		return ColaLeads{}, err
	}
	sort.SliceStable(leads, func(i, j int) bool {
		if leads[i].Prioridad != leads[j].Prioridad {
			return leads[i].Prioridad > leads[j].Prioridad
		}
		return leads[i].LeadID < leads[j].LeadID
	})
	allLeads := leads
	if filtro.Afiliado != nil || filtro.Ruta != nil {
		allLeads, err = uc.Leads.Listar(ctx, FiltroLeads{})
		if err != nil {
			return ColaLeads{}, err
		}
	}
	out := ColaLeads{Cupo10: Cupo10{PorcentajeVentana: 10}, Leads: make([]LeadEnCola, 0, len(leads))}
	for _, lead := range allLeads {
		if !lead.Afiliado && lead.Ruta == domain.RutaAsesor {
			out.Cupo10.Usados++
		}
	}
	for _, lead := range leads {
		out.Leads = append(out.Leads, LeadEnCola{LeadID: lead.LeadID, Nombre: lead.Nombre, Estado: lead.Estado, Ruta: lead.Ruta, Afiliado: lead.Afiliado, Semaforo: semaforoCola(lead), Prioridad: lead.Prioridad, Resumen: resumenCola(lead), ActualizadoEn: lead.ActualizadoEn})
	}
	return out, nil
}
func (uc *ConsultarLeads) Detalle(ctx context.Context, id string) (LeadDetalle, error) {
	lead, err := uc.lead(ctx, id)
	if err != nil {
		return LeadDetalle{}, err
	}
	var plan *domain.PlanNutricion
	if uc.Planes != nil {
		plan, err = uc.Planes.PorLead(ctx, id)
		if err != nil && !planNotFoundError(err) {
			return LeadDetalle{}, err
		}
		if planNotFoundError(err) {
			plan = nil
		}
	}
	return LeadDetalle{LeadID: lead.LeadID, Nombre: lead.Nombre, Estado: lead.Estado, Ruta: lead.Ruta, Afiliado: lead.Afiliado, Semaforo: semaforoCola(lead), Prioridad: lead.Prioridad, ConsumeCupo10: lead.ConsumeCupo10, Perfil: lead.Perfil, Capacidad: lead.Capacidad, Intencion: lead.Intencion, Plan: plan, CreadoEn: lead.CreadoEn, ActualizadoEn: lead.ActualizadoEn}, nil
}
func (uc *ConsultarLeads) Ficha(ctx context.Context, id string) (*domain.Ficha, error) {
	if uc == nil || uc.Leads == nil || uc.Fichas == nil {
		return nil, fmt.Errorf("%w: repositorios de leads y fichas requeridos", ErrValidacion)
	}
	if _, err := uc.lead(ctx, id); err != nil {
		return nil, err
	}
	ficha, err := uc.Fichas.PorLead(ctx, id)
	if err != nil {
		if fichaNotFoundError(err) {
			return nil, &NotFoundError{Resource: "ficha", ID: id}
		}
		return nil, err
	}
	if ficha == nil {
		return nil, &NotFoundError{Resource: "ficha", ID: id}
	}
	return ficha, nil
}
func (uc *ConsultarLeads) lead(ctx context.Context, id string) (*domain.Lead, error) {
	if uc == nil || uc.Leads == nil {
		return nil, fmt.Errorf("%w: repositorio de leads requerido", ErrValidacion)
	}
	if strings.TrimSpace(id) == "" || strings.Contains(id, "/") {
		return nil, &NotFoundError{Resource: "lead", ID: id}
	}
	lead, err := uc.Leads.PorID(ctx, id)
	if err != nil {
		return nil, err
	}
	if lead == nil {
		return nil, &NotFoundError{Resource: "lead", ID: id}
	}
	return lead, nil
}
func fichaNotFoundError(err error) bool {
	var notFound *NotFoundError
	return errors.Is(err, ErrNoEncontrado) && (!errors.As(err, &notFound) || strings.EqualFold(notFound.Resource, "ficha"))
}
func planNotFoundError(err error) bool {
	var notFound *NotFoundError
	return errors.Is(err, ErrNoEncontrado) && errors.As(err, &notFound) && strings.EqualFold(notFound.Resource, "plan")
}
func semaforoCola(lead *domain.Lead) domain.Semaforo {
	switch {
	case lead.Ruta == domain.RutaAsesor && lead.Afiliado:
		return domain.SemaforoVerde
	case lead.Ruta == domain.RutaAsesor || lead.Ruta == domain.RutaNutricion:
		return domain.SemaforoAmbar
	default:
		return domain.SemaforoGris
	}
}
func resumenCola(lead *domain.Lead) string {
	parts := make([]string, 0, 3)
	if lead.Afiliado {
		category := "N/A"
		if lead.Perfil != nil {
			if value, ok := lead.Perfil.Texto("categoria"); ok && strings.TrimSpace(value) != "" {
				category = value
			}
		}
		parts = append(parts, "Afiliada cat. "+category)
	} else {
		parts = append(parts, "No afiliado")
	}
	if lead.Capacidad != nil {
		parts = append(parts, fmt.Sprintf("presupuesto $%.1fM", float64(lead.Capacidad.PresupuestoMax)/1_000_000))
	}
	if lead.Intencion != nil && lead.Intencion.Nivel != "" {
		parts = append(parts, "intención "+strings.ToLower(string(lead.Intencion.Nivel)))
	}
	return strings.Join(parts, " · ")
}
