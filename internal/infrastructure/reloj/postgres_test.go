package reloj

import (
	"context"
	"testing"
	"time"
)

type memoryDemo struct {
	date  time.Time
	saved int
}

func (m *memoryDemo) FechaSimulada(context.Context) (time.Time, error) { return m.date, nil }
func (m *memoryDemo) GuardarFechaSimulada(_ context.Context, date time.Time) error {
	m.date, m.saved = date, m.saved+1
	return nil
}

func TestNuevoPostgresLoadsAndAdvancesSafely(t *testing.T) {
	initial := time.Date(2026, 7, 26, 8, 0, 0, 0, time.FixedZone("x", -5*3600))
	repo := &memoryDemo{date: initial}
	clock, err := NuevoPostgres(context.Background(), repo)
	if err != nil || !clock.FechaSimulada().Equal(initial.UTC()) || repo.saved != 0 {
		t.Fatalf("clock=%v saved=%d err=%v", clock, repo.saved, err)
	}
	clock.Avanzar(initial.AddDate(0, 0, 1))
	if !clock.Ahora().Equal(initial.UTC().AddDate(0, 0, 1)) {
		t.Fatalf("now=%v", clock.Ahora())
	}
}
