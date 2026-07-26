package usecase

import (
	"context"
	"fmt"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
	"strings"
)

// SaludarLead persists the deterministic first Vivi message for a new lead.
type SaludarLead struct {
	Leads LeadRepository
	IDs   GeneradorID
	Reloj Reloj
}

func (uc *SaludarLead) Ejecutar(ctx context.Context, event Evento) error {
	if uc == nil || uc.Leads == nil || uc.IDs == nil || uc.Reloj == nil {
		return fmt.Errorf("%w: dependencias del saludo requeridas", ErrValidacion)
	}
	if strings.TrimSpace(event.LeadID) == "" {
		return fmt.Errorf("%w: lead_id vacio", ErrValidacion)
	}
	lead, err := uc.Leads.PorID(ctx, event.LeadID)
	if err != nil {
		return err
	}
	if lead == nil {
		return ErrNoEncontrado
	}
	texto := fmt.Sprintf("¡Hola %s! 👋 Estamos aquí para ayudarte a encontrar tu vivienda ideal. ¿Sueñas con comprar este año?", lead.Nombre)
	if lead.Afiliado {
		texto = fmt.Sprintf("¡Hola %s! 👋 Como afiliada tienes un subsidio de hasta $52,5M. ¿Sueñas con comprar este año?", lead.Nombre)
	}
	return uc.Leads.AgregarMensaje(ctx, &domain.Mensaje{
		MensajeID: uc.IDs.Nuevo(), LeadID: lead.LeadID, Autor: domain.AutorMensajeVivi,
		TipoContenido: domain.TipoContenidoTexto, Texto: texto, CreadoEn: uc.Reloj.Ahora().UTC(),
	})
}
