package usecase

import (
	"context"
	"fmt"
	"time"
)

type EntradaAvanzarDemo struct {
	Hasta *time.Time
	Dias  *int
}

type SalidaAvanzarDemo struct {
	FechaSimulada   time.Time `json:"fecha_simulada"`
	HitosDisparados int       `json:"hitos_disparados"`
}

type AvanzarDemo struct {
	Demo  DemoRepository
	Reloj Reloj
	Bus   BusEventos
}

func (uc *AvanzarDemo) Ejecutar(ctx context.Context, entrada EntradaAvanzarDemo) (SalidaAvanzarDemo, error) {
	if uc == nil || uc.Demo == nil || uc.Reloj == nil {
		return SalidaAvanzarDemo{}, fmt.Errorf("%w: puertos de demo requeridos", ErrValidacion)
	}
	if (entrada.Hasta == nil) == (entrada.Dias == nil) || (entrada.Dias != nil && *entrada.Dias <= 0) {
		return SalidaAvanzarDemo{}, ErrValidacion
	}
	hasta := uc.Reloj.FechaSimulada()
	if entrada.Hasta != nil {
		hasta = entrada.Hasta.UTC()
	} else {
		hasta = hasta.AddDate(0, 0, *entrada.Dias).UTC()
	}
	if hasta.Before(uc.Reloj.FechaSimulada()) {
		return SalidaAvanzarDemo{}, ErrTiempoSimuladoAtras
	}
	if err := uc.Demo.GuardarFechaSimulada(ctx, hasta); err != nil {
		return SalidaAvanzarDemo{}, fmt.Errorf("persistir fecha simulada: %w", err)
	}
	uc.Reloj.Avanzar(hasta)
	result := &ResultadoTick{}
	if uc.Bus != nil {
		uc.Bus.Publicar(ConResultadoTick(ctx, result), Evento{Tipo: EvTickReloj, Payload: map[string]any{"fecha_simulada": hasta}})
		if result.Err != nil {
			return SalidaAvanzarDemo{FechaSimulada: hasta, HitosDisparados: result.HitosDisparados}, result.Err
		}
	}
	return SalidaAvanzarDemo{FechaSimulada: hasta, HitosDisparados: result.HitosDisparados}, nil
}
