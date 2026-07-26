package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain/motor"
)

const (
	campoIngreso   = "ingreso_hogar"
	campoCategoria = "categoria"
	campoSegmento  = "segmento"
	campoPersonas  = "personas_hogar"
	campoTipoHogar = "tipo_hogar"
	campoVivienda  = "tiene_vivienda"
	campoSubsidio  = "recibio_subsidio"
	campoHogar     = "hogar_con_afiliado"
	campoFamiliar  = "cedula_familiar_afiliado"
)

var (
	ErrFamiliarNoEncontrado = errors.New("familiar afiliado no encontrado")
	ErrFamiliarYaRegistrado = errors.New("familiar afiliado ya registrado")
)

type EntradaPerfilar struct {
	Nombre   string
	Telefono string
	Cedula   string
	Fuente   string
}

type SalidaPerfilar struct {
	LeadID            string
	Estado            domain.EstadoLead
	AfiliadoDetectado bool
}

type PerfilarLead struct {
	Leads    LeadRepository
	Catalogo CatalogoRepository
	IDs      GeneradorID
	Bus      BusEventos
	Reloj    Reloj
}

func (uc *PerfilarLead) Ejecutar(ctx context.Context, entrada EntradaPerfilar) (SalidaPerfilar, error) {
	if err := ctx.Err(); err != nil {
		return SalidaPerfilar{}, err
	}
	now := uc.Reloj.Ahora()
	a, err := uc.lookupActive(ctx, entrada.Cedula)
	if err != nil {
		return SalidaPerfilar{}, err
	}
	perfil := domain.Perfil{}
	nombre := entrada.Nombre
	if a != nil {
		if nombre == "" {
			nombre = a.Nombre
		}
		perfil = perfilAfiliado(a, now)
	}
	lead := &domain.Lead{
		LeadID: uc.IDs.Nuevo(), Nombre: nombre, Telefono: entrada.Telefono,
		Cedula: entrada.Cedula, Fuente: entrada.Fuente, Estado: domain.EstadoLeadNuevo,
		Afiliado: a != nil, Perfil: perfil, CreadoEn: now, ActualizadoEn: now,
	}
	lead.Capacidad = capacity(lead)
	if err := ctx.Err(); err != nil {
		return salida(lead), err
	}
	if err := lead.Transicionar(domain.EstadoLeadPerfilando); err != nil {
		return salida(lead), fmt.Errorf("transicionar lead: %w", err)
	}
	if err := uc.Leads.Crear(ctx, lead); err != nil {
		return salida(lead), fmt.Errorf("crear lead: %w", err)
	}
	uc.Bus.Publicar(ctx, Evento{Tipo: EvLeadNuevo, LeadID: lead.LeadID,
		Payload: map[string]any{"cedula": lead.Cedula, "nombre": lead.Nombre,
			"telefono": lead.Telefono, "fuente": lead.Fuente}})
	return salida(lead), nil
}

func (uc *PerfilarLead) ReconsultarPorFamiliar(ctx context.Context, leadID, cedula string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	lead, err := uc.Leads.PorID(ctx, leadID)
	if err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if lead.Perfil == nil {
		lead.Perfil = domain.Perfil{}
	}
	if field, ok := lead.Perfil[campoFamiliar]; ok && field.Fuente == domain.FuenteCampoVerificadoBase {
		if value, ok := field.Valor.(string); ok && value == cedula {
			return nil
		}
		return ErrFamiliarYaRegistrado
	}
	a, lookupErr := uc.lookupActive(ctx, cedula)
	if lookupErr != nil {
		return lookupErr
	}
	now := uc.Reloj.Ahora()
	if a == nil {
		lead.Perfil[campoFamiliar] = campo(cedula, domain.FuenteCampoDeclarado, .5, true, now)
		if err := uc.Leads.Guardar(ctx, lead); err != nil {
			return fmt.Errorf("guardar familiar no encontrado: %w", err)
		}
		return ErrFamiliarNoEncontrado
	}
	income, _ := lead.Perfil.Entero(campoIngreso)
	lead.Perfil[campoIngreso] = campo(income+a.IngresoMensual, domain.FuenteCampoVerificadoBase, 1, false, now)
	lead.Perfil[campoHogar] = campo(true, domain.FuenteCampoVerificadoBase, 1, false, now)
	lead.Perfil[campoFamiliar] = campo(cedula, domain.FuenteCampoVerificadoBase, 1, false, now)
	lead.Capacidad = capacity(lead)
	if err := uc.Leads.Guardar(ctx, lead); err != nil {
		return fmt.Errorf("guardar reconsulta familiar: %w", err)
	}
	return nil
}

func (uc *PerfilarLead) lookupActive(ctx context.Context, cedula string) (*Afiliado, error) {
	if cedula == "" {
		return nil, nil
	}
	a, err := uc.Catalogo.AfiliadoPorCedula(ctx, cedula)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}
		return nil, nil
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if a == nil || !a.AfiliadoActivo {
		return nil, nil
	}
	return a, nil
}

func perfilAfiliado(a *Afiliado, now time.Time) domain.Perfil {
	return domain.Perfil{
		campoIngreso:   campo(a.IngresoMensual, domain.FuenteCampoVerificadoBase, 1, false, now),
		campoCategoria: campo(a.Categoria, domain.FuenteCampoVerificadoBase, 1, false, now),
		campoSegmento:  campo(a.Segmento, domain.FuenteCampoVerificadoBase, 1, false, now),
		campoPersonas:  campo(int64(a.PersonasACargo+1), domain.FuenteCampoVerificadoBase, 1, false, now),
		campoTipoHogar: campo(a.TipoHogar, domain.FuenteCampoVerificadoBase, 1, false, now),
		campoVivienda:  campo(false, domain.FuenteCampoVerificadoBase, 1, false, now),
		campoSubsidio:  campo(false, domain.FuenteCampoVerificadoBase, 1, false, now),
	}
}

func campo(value any, source domain.FuenteCampo, confidence float64, confirm bool, now time.Time) domain.CampoPerfil {
	return domain.CampoPerfil{Valor: value, Fuente: source, Confianza: confidence, RequiereConfirmacion: confirm, ActualizadoEn: now}
}

func capacity(lead *domain.Lead) *domain.Capacidad {
	value := motor.CalcularCapacidad(lead.Perfil, lead.Afiliado, 0)
	return &value
}

func salida(lead *domain.Lead) SalidaPerfilar {
	return SalidaPerfilar{LeadID: lead.LeadID, Estado: lead.Estado, AfiliadoDetectado: lead.Afiliado}
}
