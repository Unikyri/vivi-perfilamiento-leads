package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
)

const (
	demoClockKey     = "fecha_simulada"
	approvedDemoDate = "2026-07-26T00:00:00Z"
)

type DemoRepository struct{ pool pgxPool }

func NuevoDemoRepository(pool pgxPool) *DemoRepository { return &DemoRepository{pool: pool} }

var _ usecase.DemoRepository = (*DemoRepository)(nil)
var _ usecase.DemoResetRepository = (*DemoRepository)(nil)

func (r *DemoRepository) FechaSimulada(ctx context.Context) (time.Time, error) {
	var value string
	err := r.pool.QueryRow(ctx, `SELECT valor FROM demo WHERE clave=$1`, demoClockKey).Scan(&value)
	if err != nil {
		if err == pgx.ErrNoRows {
			return time.Time{}, nil
		}
		return time.Time{}, repositoryError("demo", demoClockKey, err)
	}
	if value == "" {
		return time.Time{}, nil
	}
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		parsed, err = time.Parse("2006-01-02", value)
	}
	if err != nil {
		return time.Time{}, repositoryError("demo", demoClockKey, err)
	}
	return parsed.UTC(), nil
}

func (r *DemoRepository) GuardarFechaSimulada(ctx context.Context, value time.Time) error {
	if value.IsZero() {
		return usecase.ErrValidacion
	}
	_, err := r.pool.Exec(ctx, `INSERT INTO demo (clave,valor) VALUES ($1,$2) ON CONFLICT (clave) DO UPDATE SET valor=EXCLUDED.valor`, demoClockKey, value.UTC().Format(time.RFC3339Nano))
	return repositoryError("demo", demoClockKey, err)
}

// Reiniciar atomically clears only demo-owned lead data and restores the
// canonical simulated date. Catalog buyers and schema objects are untouched.
func (r *DemoRepository) Reiniciar(ctx context.Context) (time.Time, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return time.Time{}, repositoryError("demo reset", "", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	for _, table := range []string{"fichas", "hitos", "planes", "mensajes", "leads"} {
		if _, err := tx.Exec(ctx, "DELETE FROM "+table); err != nil {
			return time.Time{}, repositoryError("demo reset", table, err)
		}
	}
	if _, err := tx.Exec(ctx, `INSERT INTO demo (clave,valor) VALUES ($1,$2) ON CONFLICT (clave) DO UPDATE SET valor=EXCLUDED.valor`, demoClockKey, approvedDemoDate); err != nil {
		return time.Time{}, repositoryError("demo reset", demoClockKey, err)
	}
	if err := tx.Commit(ctx); err != nil {
		return time.Time{}, repositoryError("demo reset", "commit", err)
	}
	date, err := time.Parse(time.RFC3339, approvedDemoDate)
	if err != nil {
		return time.Time{}, err
	}
	return date, nil
}
