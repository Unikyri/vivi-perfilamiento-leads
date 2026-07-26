package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain/motor"
)

const reminderSinConsentimiento = "Dejamos la puerta abierta para cuando quieras retomar tu plan. Escríbenos cuando estés listo."

type EntradaCrearPlan struct {
	LeadID         string
	Consintio      bool
	Frecuencia     string
	PrecioObjetivo int64
}

type GestionarPlan struct {
	Leads      LeadRepository
	Planes     PlanRepository
	Reloj      Reloj
	IDs        GeneradorID
	Calendario []domain.EventoCalendario
}

// CrearPlan validates consent and explicit economic inputs before persisting a
// deterministic plan. The plan is persisted before the lead transition so a
// later retry can find and reuse it after a lead-save failure.
func (uc *GestionarPlan) CrearPlan(ctx context.Context, entrada EntradaCrearPlan) (*domain.PlanNutricion, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if strings.TrimSpace(entrada.LeadID) == "" {
		return nil, fmt.Errorf("%w: lead_id vacio", ErrValidacion)
	}
	if uc.Leads == nil || uc.Planes == nil || uc.Reloj == nil || uc.IDs == nil {
		return nil, fmt.Errorf("%w: repositorios y puertos requeridos", ErrValidacion)
	}
	lead, err := uc.Leads.PorID(ctx, entrada.LeadID)
	if err != nil {
		return nil, err
	}
	if lead == nil || strings.TrimSpace(lead.LeadID) == "" {
		return nil, ErrDatosLeadAusentes
	}
	if !entrada.Consintio {
		mensaje := &domain.Mensaje{MensajeID: uc.IDs.Nuevo(), LeadID: lead.LeadID, Autor: domain.AutorMensajeVivi, TipoContenido: domain.TipoContenidoTexto, Texto: reminderSinConsentimiento, CreadoEn: uc.Reloj.Ahora()}
		if err := uc.Leads.AgregarMensaje(ctx, mensaje); err != nil {
			return nil, fmt.Errorf("guardar recordatorio: %w", err)
		}
		return nil, nil
	}
	if entrada.PrecioObjetivo <= 0 || lead.Capacidad == nil {
		return nil, fmt.Errorf("%w: precio objetivo y capacidad requeridos", ErrValidacion)
	}
	if entrada.Frecuencia != "QUINCENAL" && entrada.Frecuencia != "MENSUAL" && entrada.Frecuencia != "TRIMESTRAL" {
		return nil, fmt.Errorf("%w: frecuencia invalida", ErrValidacion)
	}
	if lead.Ruta != domain.RutaNutricion || (lead.Estado != domain.EstadoLeadCalificado && lead.Estado != domain.EstadoLeadEnNutricion) {
		return nil, fmt.Errorf("%w: lead %q fuera de nutricion", ErrLeadNoCalificable, lead.LeadID)
	}

	existing, err := uc.Planes.PorLead(ctx, lead.LeadID)
	if err != nil && !errors.Is(err, ErrNoEncontrado) {
		return nil, fmt.Errorf("consultar plan: %w", err)
	}
	if existing == nil {
		now := uc.Reloj.Ahora()
		gap := entrada.PrecioObjetivo - lead.Capacidad.PresupuestoMax
		if gap < 0 {
			gap = 0
		}
		hitos := motor.DisenarHitos(gap, esConversion(lead), now, uc.Calendario)
		for i := range hitos {
			hitos[i].HitoID = uc.IDs.Nuevo()
		}
		consentAt := now
		plan := &domain.PlanNutricion{PlanID: uc.IDs.Nuevo(), LeadID: lead.LeadID, Estado: domain.EstadoPlanActivo, ConsentimientoEn: &consentAt, Frecuencia: entrada.Frecuencia, MetaMonto: gap, MetaDescripcion: "Cierre de brecha de vivienda", Hitos: hitos}
		if err := uc.Planes.Crear(ctx, plan); err != nil {
			return nil, fmt.Errorf("crear plan: %w", err)
		}
		existing = plan
	}

	if lead.Estado == domain.EstadoLeadCalificado {
		if err := lead.Transicionar(domain.EstadoLeadEnNutricion); err != nil {
			return existing, fmt.Errorf("transicionar lead: %w", err)
		}
		if err := uc.Leads.Guardar(ctx, lead); err != nil {
			return existing, fmt.Errorf("guardar lead: %w", err)
		}
	}
	return existing, nil
}
