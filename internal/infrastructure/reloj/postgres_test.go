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

func TestNuevoPostgresLoadsOrPersistsFallback(t *testing.T) {
	initial := time.Date(2026, 7, 26, 8, 0, 0, 0, time.FixedZone("x", -5*3600))
	before := time.Now().UTC()
	tests := []struct {
		name      string
		date      time.Time
		expected  time.Time
		savedOnce int
	}{
		{name: "loaded UTC date", date: initial, expected: initial.UTC()},
		{name: "zero value fallback", expected: time.Time{}, savedOnce: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &memoryDemo{date: tt.date}
			clock, err := NuevoPostgres(context.Background(), repo)
			if err != nil {
				t.Fatal(err)
			}
			if repo.saved != tt.savedOnce {
				t.Fatalf("saved=%d, want %d", repo.saved, tt.savedOnce)
			}
			if tt.date.IsZero() {
				if repo.date.IsZero() || repo.date.Before(before) || repo.date.After(time.Now().UTC()) {
					t.Fatalf("fallback=%v is outside construction window", repo.date)
				}
				return
			}
			if !clock.FechaSimulada().Equal(tt.expected) {
				t.Fatalf("simulated=%v, want %v", clock.FechaSimulada(), tt.expected)
			}
		})
	}
}

func TestPostgresAdvanceChangesOnlySimulatedTime(t *testing.T) {
	initial := time.Date(2026, 7, 26, 8, 0, 0, 0, time.UTC)
	repo := &memoryDemo{date: initial}
	clock, err := NuevoPostgres(context.Background(), repo)
	if err != nil {
		t.Fatal(err)
	}
	advanced := time.Date(2099, 1, 2, 3, 4, 5, 0, time.FixedZone("demo", -5*3600))
	before := time.Now().UTC()
	clock.Avanzar(advanced)
	wall := clock.Ahora()
	after := time.Now().UTC()
	if !clock.FechaSimulada().Equal(advanced.UTC()) {
		t.Fatalf("simulated=%v, want %v", clock.FechaSimulada(), advanced.UTC())
	}
	if wall.Before(before) || wall.After(after) || !wall.Equal(wall.UTC()) {
		t.Fatalf("wall=%v is outside current UTC window", wall)
	}
}
