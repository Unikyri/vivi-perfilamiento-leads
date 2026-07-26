package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
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
var _ usecase.DemoSeedRepository = (*DemoRepository)(nil)
var _ usecase.DemoResetSeedRepository = (*DemoRepository)(nil)

const demoLeadUpsertSQL = `INSERT INTO leads (lead_id,nombre,telefono,cedula,fuente,estado,ruta,afiliado,prioridad,consume_cupo_10,perfil,capacidad,intencion,version,creado_en,actualizado_en) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16) ON CONFLICT (lead_id) DO UPDATE SET nombre=EXCLUDED.nombre,telefono=EXCLUDED.telefono,cedula=EXCLUDED.cedula,fuente=EXCLUDED.fuente,estado=EXCLUDED.estado,ruta=EXCLUDED.ruta,afiliado=EXCLUDED.afiliado,prioridad=EXCLUDED.prioridad,consume_cupo_10=EXCLUDED.consume_cupo_10,perfil=EXCLUDED.perfil,capacidad=EXCLUDED.capacidad,intencion=EXCLUDED.intencion,version=EXCLUDED.version,creado_en=EXCLUDED.creado_en,actualizado_en=EXCLUDED.actualizado_en`

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

func (r *DemoRepository) Sembrar(ctx context.Context, leads []domain.Lead) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return repositoryError("demo seed", "", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	for i := range leads {
		if err := insertDemoLead(ctx, tx, leads[i]); err != nil {
			return err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return repositoryError("demo seed", "commit", err)
	}
	return nil
}

func (r *DemoRepository) Reiniciar(ctx context.Context) (time.Time, error) {
	return r.reiniciar(ctx, nil)
}

func (r *DemoRepository) ReiniciarConSeed(ctx context.Context, leads []domain.Lead) (time.Time, error) {
	return r.reiniciar(ctx, leads)
}

func (r *DemoRepository) reiniciar(ctx context.Context, leads []domain.Lead) (time.Time, error) {
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
	for i := range leads {
		if err := insertDemoLead(ctx, tx, leads[i]); err != nil {
			return time.Time{}, err
		}
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

func insertDemoLead(ctx context.Context, tx pgx.Tx, lead domain.Lead) error {
	perfil, err := encodeJSONB(lead.Perfil)
	if err != nil {
		return err
	}
	capacidad, err := encodeJSONB(lead.Capacidad)
	if err != nil {
		return err
	}
	intencion, err := encodeJSONB(lead.Intencion)
	if err != nil {
		return err
	}
	version := lead.Version
	if version == 0 {
		version = 1
	}
	if _, err := tx.Exec(ctx, demoLeadUpsertSQL, lead.LeadID, lead.Nombre, lead.Telefono, lead.Cedula, lead.Fuente, lead.Estado, lead.Ruta, lead.Afiliado, lead.Prioridad, lead.ConsumeCupo10, perfil, capacidad, intencion, version, lead.CreadoEn, lead.ActualizadoEn); err != nil {
		return repositoryError("demo lead", lead.LeadID, err)
	}
	return nil
}
