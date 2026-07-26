package agentes

import (
	"context"
	"log/slog"
	"time"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
)

// ManejadorEvento is an optional observation or orchestration callback.
type ManejadorEvento func(context.Context, usecase.Evento) error

// Calificador is the narrow qualification handoff used by the coordinator.
type Calificador interface {
	Ejecutar(context.Context, usecase.EntradaCalificar) (usecase.SalidaCalificar, error)
}

// Documentadora is the narrow advisor-ficha handoff used by the coordinator.
type Documentadora interface {
	Ejecutar(context.Context, string) (*domain.Ficha, error)
}

// Nutricionista is the narrow clock-driven milestone handoff.
type Nutricionista interface {
	Ejecutar(context.Context, time.Time) (int, error)
}

// Dependencias contains optional coordinator callbacks and use cases.
type Dependencias struct {
	LeadNuevo     ManejadorEvento
	Calificador   Calificador
	Documentadora Documentadora
	Nutricionista Nutricionista
	Logger        *slog.Logger
}
