package postgres

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
)

func TestPostgresIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("PostgreSQL integration disabled by -short")
	}
	url := os.Getenv("VIVI_TEST_DATABASE_URL")
	if url == "" {
		t.Skip("VIVI_TEST_DATABASE_URL is not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pool, err := Conectar(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	if _, err := NuevoLeadRepository(pool).Listar(ctx, usecase.FiltroLeads{}); err != nil {
		t.Fatalf("lead repository against existing schema: %v", err)
	}
	if _, err := NuevoPlanRepository(pool).HitosVencidos(ctx, time.Now()); err != nil {
		t.Fatalf("plan repository against existing schema: %v", err)
	}
	_, err = NuevoFichaRepository(pool).PorLead(ctx, "integration-missing-lead")
	if !errors.Is(err, usecase.ErrNoEncontrado) {
		t.Fatalf("ficha repository missing-lead result = %v", err)
	}
}
