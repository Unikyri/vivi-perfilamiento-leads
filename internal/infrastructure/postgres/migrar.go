package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Unikyri/vivi-perfilamiento-leads/migrations"
)

// Migrar applies the canonical schema DDL (Contract §5) to the database.
// Idempotent: uses IF NOT EXISTS for all tables and indexes.
func Migrar(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, migrations.EsquemaInicial)
	if err != nil {
		return fmt.Errorf("postgres.Migrar: %w", err)
	}
	return nil
}
