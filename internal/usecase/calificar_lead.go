package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain/motor"
)

var (
	ErrLeadNoCalificable = errors.New("lead no calificable")
	ErrDatosLeadAusentes = errors.New("datos del lead ausentes")
)

type EntradaCalificar struct{ LeadID string }
type EntradaCalificarLead = EntradaCalificar

type SalidaCalificar struct {
	LeadID          string
	Ruta            domain.Ruta
	Estado          domain.EstadoLead
	Capacidad       domain.Capacidad
	Vecinos         []motor.Vecino
	Recomendaciones []domain.Recomendacion
	Conversion      bool
	Prioridad       float64
	Semaforo        domain.Semaforo
	ConsumeCupo10   bool
}

type CalificarLead struct {
	Leads    LeadRepository
	Catalogo CatalogoRepository
	Bus      BusEventos
	Reloj    Reloj
}

func (uc *CalificarLead) Ejecutar(ctx context.Context, entrada EntradaCalificar) (SalidaCalificar, error) {
	if err := ctx.Err(); err != nil {
		return SalidaCalificar{}, err
	}
	if strings.TrimSpace(entrada.LeadID) == "" {
		return SalidaCalificar{}, fmt.Errorf("%w: lead_id vacio", ErrValidacion)
	}
	if uc.Leads == nil || uc.Catalogo == nil {
		return SalidaCalificar{}, fmt.Errorf("%w: repositorios requeridos", ErrValidacion)
	}
	lead, err := uc.Leads.PorID(ctx, entrada.LeadID)
	if err != nil {
		return SalidaCalificar{}, err
	}
	if err = ctx.Err(); err != nil {
		return SalidaCalificar{}, err
	}
	if lead == nil {
		return SalidaCalificar{}, ErrDatosLeadAusentes
	}
	if strings.TrimSpace(lead.LeadID) == "" {
		return SalidaCalificar{}, fmt.Errorf("%w: lead_id vacio", ErrValidacion)
	}
	if lead.Estado != domain.EstadoLeadCalificado {
		return SalidaCalificar{}, fmt.Errorf("%w: lead %q en estado %s", ErrLeadNoCalificable, lead.LeadID, lead.Estado)
	}
	if lead.Intencion == nil {
		return SalidaCalificar{}, fmt.Errorf("%w: intencion ausente", ErrDatosLeadAusentes)
	}
	if lead.Capacidad == nil {
		return SalidaCalificar{}, fmt.Errorf("%w: capacidad ausente", ErrDatosLeadAusentes)
	}
	proyectos, err := uc.Catalogo.Proyectos(ctx)
	if err != nil {
		return SalidaCalificar{}, fmt.Errorf("cargar proyectos: %w", err)
	}
	if err = ctx.Err(); err != nil {
		return SalidaCalificar{}, err
	}
	compradores, err := uc.Catalogo.Compradores(ctx)
	if err != nil {
		return SalidaCalificar{}, fmt.Errorf("cargar compradores: %w", err)
	}
	if err = ctx.Err(); err != nil {
		return SalidaCalificar{}, err
	}
	decision := construirDecision(lead, proyectos, compradores)
	out := salidaCalificar(lead, decision)
	if lead.Ruta == domain.RutaAsesor {
		return out, nil
	}
	lead.Ruta, lead.Prioridad, lead.ConsumeCupo10 = decision.ruta, decision.prioridad, decision.consumeCupo10
	capacidad := copiaCapacidadDecision(decision.capacidad)
	lead.Capacidad = &capacidad
	if uc.Reloj != nil {
		lead.ActualizadoEn = uc.Reloj.Ahora()
	}
	if err = aplicarEstadoRuta(lead, decision.ruta); err != nil {
		return SalidaCalificar{}, err
	}
	if err = ctx.Err(); err != nil {
		return SalidaCalificar{}, err
	}
	if err = uc.Leads.Guardar(ctx, lead); err != nil {
		return SalidaCalificar{}, fmt.Errorf("guardar calificacion: %w", err)
	}
	if uc.Bus != nil {
		uc.Bus.Publicar(ctx, Evento{Tipo: EvRutaDecidida, LeadID: lead.LeadID, Payload: map[string]any{
			"ruta": decision.ruta, "prioridad": decision.prioridad, "semaforo": decision.semaforo,
			"consume_cupo_10": decision.consumeCupo10, "recomendaciones": copiarRecomendaciones(decision.recomendaciones),
		}})
	}
	out.Estado = lead.Estado
	return out, nil
}

type decisionCalificacion struct {
	capacidad       domain.Capacidad
	vecinos         []motor.Vecino
	recomendaciones []domain.Recomendacion
	conversion      bool
	ruta            domain.Ruta
	prioridad       float64
	semaforo        domain.Semaforo
	consumeCupo10   bool
}

func construirDecision(lead *domain.Lead, proyectos map[string]domain.Proyecto, compradores []domain.Comprador) decisionCalificacion {
	preliminar := motor.CalcularCapacidad(lead.Perfil, lead.Afiliado, 0)
	candidato := menorPrecioAsequible(proyectos, preliminar.PresupuestoMax)
	final := motor.CalcularCapacidad(lead.Perfil, lead.Afiliado, candidato)
	var dependientes *int
	if personas, ok := lead.Perfil.Entero(campoPersonas); ok {
		n := int(personas) - 1
		dependientes = &n
	}
	zonas := make(map[string]string, len(proyectos))
	for key, proyecto := range proyectos {
		if key == proyecto.ProyectoID {
			zonas[key] = proyecto.Zona
		}
	}
	vecinos := motor.GemeloKNN(motor.EntradaGemelo{Perfil: lead.Perfil, Afiliado: lead.Afiliado, PersonasACargo: dependientes, Compradores: compradores, ZonasCatalogo: zonas, K: 30})
	recomendaciones := motor.RecomendarProyectos(vecinos, proyectos, final.PresupuestoMax)
	conversion := esConversion(lead)
	ruta := motor.Matriz2x2(motor.EntradaMatriz{Ratio: final.Ratio, Intencion: lead.Intencion.Nivel, EsAfiliado: lead.Afiliado, RutaConversion: conversion})
	ratio := final.Ratio
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1.2 {
		ratio = 1.2
	}
	return decisionCalificacion{capacidad: final, vecinos: vecinos, recomendaciones: recomendaciones, conversion: conversion, ruta: ruta, prioridad: pesoRuta(ruta) * ratio * final.Confianza, semaforo: semaforoRuta(ruta), consumeCupo10: !lead.Afiliado && ruta == domain.RutaAsesor}
}

func menorPrecioAsequible(proyectos map[string]domain.Proyecto, presupuesto int64) int64 {
	var menor int64
	for _, proyecto := range proyectos {
		if proyecto.PrecioDesde > 0 && proyecto.PrecioDesde <= presupuesto && (menor == 0 || proyecto.PrecioDesde < menor) {
			menor = proyecto.PrecioDesde
		}
	}
	return menor
}
func esConversion(lead *domain.Lead) bool {
	if lead.Afiliado {
		return false
	}
	if laboral, ok := lead.Perfil.Texto("situacion_laboral"); ok && laboral == "INDEPENDIENTE" {
		return true
	}
	if hogar, ok := lead.Perfil.Booleano(campoHogar); ok && hogar {
		return true
	}
	caja, ok := lead.Perfil.Texto("caja_externa")
	return ok && strings.TrimSpace(caja) != ""
}
func pesoRuta(ruta domain.Ruta) float64 {
	switch ruta {
	case domain.RutaAsesor:
		return 1
	case domain.RutaNutricion:
		return .5
	case domain.RutaRemarketing:
		return .25
	case domain.RutaDespedida:
		return .1
	default:
		return 0
	}
}
func semaforoRuta(ruta domain.Ruta) domain.Semaforo {
	switch ruta {
	case domain.RutaAsesor:
		return domain.SemaforoVerde
	case domain.RutaNutricion:
		return domain.SemaforoAmbar
	default:
		return domain.SemaforoGris
	}
}
func aplicarEstadoRuta(lead *domain.Lead, ruta domain.Ruta) error {
	switch ruta {
	case domain.RutaAsesor:
		return nil
	case domain.RutaNutricion:
		return lead.Transicionar(domain.EstadoLeadEnNutricion)
	case domain.RutaRemarketing:
		return lead.Transicionar(domain.EstadoLeadRemarketing)
	case domain.RutaDespedida:
		return lead.Transicionar(domain.EstadoLeadDespedido)
	default:
		return fmt.Errorf("%w: ruta %q", ErrLeadNoCalificable, ruta)
	}
}
func estadoTrasRuta(actual domain.EstadoLead, ruta domain.Ruta) domain.EstadoLead {
	switch ruta {
	case domain.RutaNutricion:
		return domain.EstadoLeadEnNutricion
	case domain.RutaRemarketing:
		return domain.EstadoLeadRemarketing
	case domain.RutaDespedida:
		return domain.EstadoLeadDespedido
	default:
		return actual
	}
}
func salidaCalificar(lead *domain.Lead, d decisionCalificacion) SalidaCalificar {
	return SalidaCalificar{LeadID: lead.LeadID, Ruta: d.ruta, Estado: estadoTrasRuta(lead.Estado, d.ruta), Capacidad: copiaCapacidadDecision(d.capacidad), Vecinos: append([]motor.Vecino(nil), d.vecinos...), Recomendaciones: copiarRecomendaciones(d.recomendaciones), Conversion: d.conversion, Prioridad: d.prioridad, Semaforo: d.semaforo, ConsumeCupo10: d.consumeCupo10}
}
func copiaCapacidadDecision(c domain.Capacidad) domain.Capacidad {
	c.Desglose = append([]domain.ItemDesglose(nil), c.Desglose...)
	return c
}
func copiarRecomendaciones(r []domain.Recomendacion) []domain.Recomendacion {
	return append([]domain.Recomendacion(nil), r...)
}
