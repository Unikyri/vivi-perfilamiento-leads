package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Conectar creates a pgxpool connection and validates it with a ping.
// Returns a ready-to-use pool or a wrapped error.
func Conectar(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("postgres.Conectar: pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("postgres.Conectar: ping: %w", err)
	}
	return pool, nil
}
