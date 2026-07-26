package usecase

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain/motor"
)

const advertenciaPerfilParcial = "PERFIL PARCIALMENTE DECLARADO — validar campos marcados"

type GenerarFicha struct {
	Leads    LeadRepository
	Fichas   FichaRepository
	Catalogo CatalogoRepository
	IDs      GeneradorID
	Reloj    Reloj
}

func (uc *GenerarFicha) Ejecutar(ctx context.Context, leadID string) (*domain.Ficha, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if strings.TrimSpace(leadID) == "" {
		return nil, fmt.Errorf("%w: lead_id vacio", ErrValidacion)
	}
	if uc.Leads == nil || uc.Fichas == nil || uc.Catalogo == nil || uc.IDs == nil || uc.Reloj == nil {
		return nil, fmt.Errorf("%w: repositorios y puertos requeridos", ErrValidacion)
	}
	lead, err := uc.Leads.PorID(ctx, leadID)
	if err != nil {
		return nil, err
	}
	if err = ctx.Err(); err != nil {
		return nil, err
	}
	if lead == nil || strings.TrimSpace(lead.LeadID) == "" {
		return nil, ErrDatosLeadAusentes
	}
	if lead.Estado != domain.EstadoLeadCalificado || lead.Ruta != domain.RutaAsesor || lead.Intencion == nil || lead.Capacidad == nil {
		return nil, fmt.Errorf("%w: ficha requiere lead CALIFICADO/ASESOR con intencion y capacidad", ErrLeadNoCalificable)
	}

	existing, err := uc.Fichas.PorLead(ctx, lead.LeadID)
	if err != nil && !fichaNoEncontrada(err) {
		return nil, fmt.Errorf("consultar ficha: %w", err)
	}
	if err = ctx.Err(); err != nil {
		return nil, err
	}
	proyectos, err := uc.Catalogo.Proyectos(ctx)
	if err != nil {
		return nil, fmt.Errorf("cargar proyectos: %w", err)
	}
	if err = ctx.Err(); err != nil {
		return nil, err
	}
	compradores, err := uc.Catalogo.Compradores(ctx)
	if err != nil {
		return nil, fmt.Errorf("cargar compradores: %w", err)
	}
	if err = ctx.Err(); err != nil {
		return nil, err
	}

	decision := construirDecision(lead, proyectos, compradores)
	ficha := construirFicha(lead, decision, existing, uc.IDs, uc.Reloj)
	if err = ctx.Err(); err != nil {
		return nil, err
	}
	if err = uc.Fichas.Guardar(ctx, &ficha); err != nil {
		return nil, fmt.Errorf("guardar ficha: %w", err)
	}

	toSave := copiaLeadFicha(lead)
	if err = toSave.Transicionar(domain.EstadoLeadEntregado); err != nil {
		return &ficha, fmt.Errorf("transicionar lead: %w", err)
	}
	if err = ctx.Err(); err != nil {
		return &ficha, err
	}
	toSave.ActualizadoEn = uc.Reloj.Ahora()
	if err = uc.Leads.Guardar(ctx, toSave); err != nil {
		return &ficha, fmt.Errorf("entregar lead: %w", err)
	}
	return &ficha, nil
}

func fichaNoEncontrada(err error) bool {
	if !errors.Is(err, ErrNoEncontrado) {
		return false
	}
	var nf *NotFoundError
	return !errors.As(err, &nf) || nf.Resource == "ficha"
}

func construirFicha(lead *domain.Lead, d decisionCalificacion, existing *domain.Ficha, ids GeneradorID, reloj Reloj) domain.Ficha {
	id, generatedAt := "", time.Time{}
	if existing != nil {
		id, generatedAt = existing.FichaID, existing.GeneradaEn
	} else {
		id, generatedAt = ids.Nuevo(), reloj.Ahora()
	}
	ficha := domain.Ficha{FichaID: id, LeadID: lead.LeadID, GeneradaEn: generatedAt, ConfianzaPerfil: d.capacidad.Confianza,
		Identificacion: domain.Identificacion{Nombre: lead.Nombre, Afiliada: lead.Afiliado, Categoria: textoPerfil(lead.Perfil, "categoria"), Telefono: lead.Telefono},
		Capacidad:      copiaCapacidadDecision(d.capacidad), Perfil: copiaPerfilFicha(lead.Perfil), Intencion: copiaIntencionFicha(*lead.Intencion), Recomendaciones: copiarRecomendaciones(d.recomendaciones),
		Beneficios: beneficiosFicha(lead, d.capacidad), ArgumentosVenta: argumentosFicha(lead.Perfil), AlertaDesistimiento: alertaFicha(d.vecinos), ConsumeCupo10: lead.ConsumeCupo10}
	if ficha.ConfianzaPerfil < .6 {
		ficha.BandaAdvertencia = strptrFicha(advertenciaPerfilParcial)
	}
	return ficha
}

func beneficiosFicha(lead *domain.Lead, capacidad domain.Capacidad) []string {
	beneficios := make([]string, 0, 3)
	if capacidad.SubsidioAplicable > 0 {
		beneficios = append(beneficios, "Subsidio de caja 30 SMMLV")
	}
	if lead.Afiliado {
		beneficios = append(beneficios, "Crédito propio Colsubsidio", "Acompañamiento PerteneSer")
	}
	return beneficios
}

func argumentosFicha(perfil domain.Perfil) []string {
	renta, ok := perfil.Entero("arriendo_actual")
	if !ok || renta <= 0 {
		return []string{}
	}
	ingreso, _ := perfil.Entero("ingreso_hogar")
	return []string{fmt.Sprintf("Paga $%s de arriendo; la cuota estimada es $%s", formatoCOP(renta), formatoCOP(ingreso*40/100))}
}

func alertaFicha(vecinos []motor.Vecino) domain.AlertaDesistimiento {
	desistidos := 0
	for _, vecino := range vecinos {
		if vecino.Desistio {
			desistidos++
		}
	}
	rate := 0.0
	if len(vecinos) > 0 {
		rate = float64(desistidos) / float64(len(vecinos))
	}
	return domain.AlertaDesistimiento{Activa: rate > .20, TasaVecinos: rate}
}

func textoPerfil(perfil domain.Perfil, key string) string {
	value, _ := perfil.Texto(key)
	return value
}
func strptrFicha(value string) *string { return &value }
func formatoCOP(value int64) string {
	negative := value < 0
	if negative {
		value = -value
	}
	digits := strconv.FormatInt(value, 10)
	for i := len(digits) - 3; i > 0; i -= 3 {
		digits = digits[:i] + "." + digits[i:]
	}
	if negative {
		return "-" + digits
	}
	return digits
}

func copiaIntencionFicha(in domain.Intencion) domain.Intencion {
	in.Senales = append([]string(nil), in.Senales...)
	return in
}
func copiaPerfilFicha(in domain.Perfil) domain.Perfil {
	if in == nil {
		return nil
	}
	out := make(domain.Perfil, len(in))
	for key, field := range in {
		out[key] = field
	}
	return out
}
func copiaLeadFicha(in *domain.Lead) *domain.Lead {
	out := *in
	out.Perfil = copiaPerfilFicha(in.Perfil)
	if in.Capacidad != nil {
		capacity := copiaCapacidadDecision(*in.Capacidad)
		out.Capacidad = &capacity
	}
	if in.Intencion != nil {
		intention := copiaIntencionFicha(*in.Intencion)
		out.Intencion = &intention
	}
	return &out
}
