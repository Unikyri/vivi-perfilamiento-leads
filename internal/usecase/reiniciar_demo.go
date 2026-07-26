package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var ErrDemoDeshabilitado = errors.New("demo reset disabled")

// ReiniciarDemo restores the approved demo state only when explicitly enabled.
type ReiniciarDemo struct {
	Repository DemoResetRepository
	Reloj      Reloj
	Habilitado bool
}

type SalidaReiniciarDemo struct {
	Reiniciado    bool
	FechaSimulada time.Time
}

func (uc *ReiniciarDemo) Ejecutar(ctx context.Context) (SalidaReiniciarDemo, error) {
	if uc == nil || uc.Repository == nil || uc.Reloj == nil {
		return SalidaReiniciarDemo{}, fmt.Errorf("%w: puertos de reset requeridos", ErrValidacion)
	}
	if !uc.Habilitado {
		return SalidaReiniciarDemo{}, ErrDemoDeshabilitado
	}
	date, err := uc.Repository.Reiniciar(ctx)
	if err != nil {
		return SalidaReiniciarDemo{}, fmt.Errorf("reiniciar demo: %w", err)
	}
	uc.Reloj.Avanzar(date)
	return SalidaReiniciarDemo{Reiniciado: true, FechaSimulada: date.UTC()}, nil
}
