package reloj

import (
	"context"
	"testing"
	"time"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
)

type operationalLeadRepository struct {
	lead    *domain.Lead
	message *domain.Mensaje
}

func (r *operationalLeadRepository) Crear(context.Context, *domain.Lead) error { return nil }
func (r *operationalLeadRepository) PorID(_ context.Context, id string) (*domain.Lead, error) {
	if r.lead == nil || r.lead.LeadID != id {
		return nil, usecase.ErrNoEncontrado
	}
	return r.lead, nil
}
func (r *operationalLeadRepository) Guardar(context.Context, *domain.Lead) error { return nil }
func (r *operationalLeadRepository) Listar(context.Context, usecase.FiltroLeads) ([]*domain.Lead, error) {
	return nil, nil
}
func (r *operationalLeadRepository) AgregarMensaje(_ context.Context, message *domain.Mensaje) error {
	r.message = message
	return nil
}
func (r *operationalLeadRepository) Conversacion(context.Context, string) ([]domain.Mensaje, error) {
	return nil, nil
}

type operationalIDs struct{}

func (operationalIDs) Nuevo() string { return "message-1" }

func TestPostgresClockKeepsOperationalWritesOnWallTime(t *testing.T) {
	start := time.Date(2026, 7, 26, 0, 0, 0, 0, time.UTC)
	clock, err := NuevoPostgres(context.Background(), &memoryDemo{date: start})
	if err != nil {
		t.Fatal(err)
	}
	advanced := time.Date(2099, 3, 4, 0, 0, 0, 0, time.UTC)
	clock.Avanzar(advanced)

	repo := &operationalLeadRepository{lead: &domain.Lead{LeadID: "lead-1", Nombre: "Ana"}}
	before := time.Now().UTC()
	err = (&usecase.SaludarLead{Leads: repo, IDs: operationalIDs{}, Reloj: clock}).Ejecutar(
		context.Background(), usecase.Evento{LeadID: "lead-1"},
	)
	after := time.Now().UTC()
	if err != nil {
		t.Fatal(err)
	}
	if repo.message == nil {
		t.Fatal("operational write did not persist a message")
	}
	created := repo.message.CreadoEn
	if created.Before(before) || created.After(after) || !created.Equal(created.UTC()) {
		t.Fatalf("created=%v is outside current UTC wall-time window", created)
	}
	if !created.Before(advanced) {
		t.Fatalf("created=%v incorrectly follows advanced simulated date %v", created, advanced)
	}
	if !clock.FechaSimulada().Equal(advanced) {
		t.Fatalf("simulated=%v, want %v", clock.FechaSimulada(), advanced)
	}
}
