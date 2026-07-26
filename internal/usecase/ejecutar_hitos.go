package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
	"time"
)

var ErrTiempoSimuladoAtras = errors.New("tiempo simulado no puede retroceder")

const plantillaHito = "Tu próximo paso de vivienda está listo: %s.\nSi prefieres pausar estos mensajes, responde PAUSAR."

type EjecutarHitos struct {
	Leads   LeadRepository
	Planes  PlanRepository
	Gateway MensajeriaGateway
	Reloj   Reloj
	IDs     GeneradorID
	Bus     BusEventos
}

func (uc *EjecutarHitos) Ejecutar(ctx context.Context, hasta time.Time) (int, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	if uc.Leads == nil || uc.Planes == nil || uc.Gateway == nil || uc.Reloj == nil || uc.IDs == nil {
		return 0, fmt.Errorf("%w: puertos requeridos", ErrValidacion)
	}
	if hasta.Before(uc.Reloj.FechaSimulada()) {
		return 0, fmt.Errorf("%w: %s", ErrTiempoSimuladoAtras, hasta.Format(time.RFC3339))
	}
	uc.Reloj.Avanzar(hasta)
	vencidos, err := uc.Planes.HitosVencidos(ctx, hasta)
	if err != nil {
		return 0, fmt.Errorf("consultar hitos: %w", err)
	}
	var errs []error
	entregados := 0
	handoff := make(map[string]bool)
	for _, item := range vencidos {
		lead, err := uc.Leads.PorID(ctx, item.LeadID)
		if err != nil {
			errs = append(errs, fmt.Errorf("hito %s consultar lead: %w", item.Hito.HitoID, err))
			continue
		}
		if lead == nil || lead.Estado == domain.EstadoLeadPausado {
			continue
		}
		mensaje := &domain.Mensaje{MensajeID: uc.IDs.Nuevo(), LeadID: lead.LeadID, Autor: domain.AutorMensajeVivi, TipoContenido: domain.TipoContenidoHitoNutricion, Texto: textoHito(item.Hito), CreadoEn: uc.Reloj.Ahora()}
		if err := uc.Gateway.Enviar(ctx, mensaje); err != nil {
			errs = append(errs, fmt.Errorf("hito %s enviar: %w", item.Hito.HitoID, err))
			continue
		}
		if err := uc.Leads.AgregarMensaje(ctx, mensaje); err != nil {
			errs = append(errs, fmt.Errorf("hito %s guardar mensaje: %w", item.Hito.HitoID, err))
			continue
		}
		if err := uc.Planes.MarcarHito(ctx, item.Hito.HitoID, domain.EstadoHitoNotificado); err != nil {
			errs = append(errs, fmt.Errorf("hito %s marcar: %w", item.Hito.HitoID, err))
			continue
		}
		entregados++
		if handoff[item.LeadID] {
			continue
		}
		plan, err := uc.Planes.PorLead(ctx, item.LeadID)
		if err != nil {
			errs = append(errs, fmt.Errorf("hito %s consultar plan: %w", item.Hito.HitoID, err))
			continue
		}
		if !requiereReperfilado(plan, item.Hito) {
			continue
		}
		if lead.Estado == domain.EstadoLeadEnNutricion {
			if err := lead.Transicionar(domain.EstadoLeadPerfilando); err != nil {
				errs = append(errs, fmt.Errorf("hito %s transicionar lead: %w", item.Hito.HitoID, err))
				continue
			}
			if err := uc.Leads.Guardar(ctx, lead); err != nil {
				errs = append(errs, fmt.Errorf("hito %s guardar lead: %w", item.Hito.HitoID, err))
				continue
			}
		} else if lead.Estado != domain.EstadoLeadPerfilando {
			continue
		}
		if uc.Bus != nil {
			uc.Bus.Publicar(ctx, Evento{Tipo: EvPerfilCompleto, LeadID: item.LeadID})
		}
		handoff[item.LeadID] = true
	}
	return entregados, errors.Join(errs...)
}
func textoHito(hito domain.Hito) string { return fmt.Sprintf(plantillaHito, hito.Descripcion) }
func requiereReperfilado(plan *domain.PlanNutricion, actual domain.Hito) bool {
	if plan == nil {
		return false
	}
	if actual.Tipo == domain.TipoHitoReevaluacion {
		return true
	}
	if plan.MetaMonto <= 0 {
		return false
	}
	var total int64
	presente := false
	for _, hito := range plan.Hitos {
		if hito.HitoID == actual.HitoID {
			presente = true
		}
		if hito.Estado == domain.EstadoHitoNotificado && hito.Monto != nil {
			total += *hito.Monto
		}
	}
	if !presente && actual.Monto != nil {
		total += *actual.Monto
	}
	return total >= plan.MetaMonto
}
