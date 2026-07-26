package reloj

import (
	"context"
	"sync"
	"time"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
)

type Postgres struct {
	mu  sync.RWMutex
	now time.Time
}

func NuevoPostgres(ctx context.Context, repository usecase.DemoRepository) (*Postgres, error) {
	if repository == nil {
		return nil, usecase.ErrValidacion
	}
	now, err := repository.FechaSimulada(ctx)
	if err != nil {
		return nil, err
	}
	if now.IsZero() {
		now = time.Now().UTC()
		if err := repository.GuardarFechaSimulada(ctx, now); err != nil {
			return nil, err
		}
	}
	return &Postgres{now: now.UTC()}, nil
}

func (r *Postgres) Ahora() time.Time {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.now
}

func (r *Postgres) FechaSimulada() time.Time { return r.Ahora() }

func (r *Postgres) Avanzar(at time.Time) {
	r.mu.Lock()
	r.now = at.UTC()
	r.mu.Unlock()
}

var _ usecase.Reloj = (*Postgres)(nil)
