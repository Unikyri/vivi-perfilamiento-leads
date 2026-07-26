package usecase

import (
	"context"
	"testing"
	"time"
)

type resetDemoFake struct {
	date   time.Time
	resets int
}

func (f *resetDemoFake) Reiniciar(context.Context) (time.Time, error) {
	f.resets++
	return f.date, nil
}

type resetClockFake struct{ date time.Time }

func (c *resetClockFake) Ahora() time.Time         { return c.date }
func (c *resetClockFake) FechaSimulada() time.Time { return c.date }
func (c *resetClockFake) Avanzar(date time.Time)   { c.date = date }

func TestReiniciarDemoDisabledDoesNotMutate(t *testing.T) {
	date := time.Date(2026, 7, 28, 0, 0, 0, 0, time.UTC)
	repo, clock := &resetDemoFake{date: date}, &resetClockFake{date: date.Add(24 * time.Hour)}
	_, err := (&ReiniciarDemo{Repository: repo, Reloj: clock}).Ejecutar(context.Background())
	if err != ErrDemoDeshabilitado || repo.resets != 0 || !clock.date.Equal(date.Add(24*time.Hour)) {
		t.Fatalf("err=%v resets=%d date=%v", err, repo.resets, clock.date)
	}
}

func TestReiniciarDemoIsIdempotentAndFast(t *testing.T) {
	date := time.Date(2026, 7, 26, 0, 0, 0, 0, time.UTC)
	repo, clock := &resetDemoFake{date: date}, &resetClockFake{date: date.AddDate(0, 0, 4)}
	uc := &ReiniciarDemo{Repository: repo, Reloj: clock, Habilitado: true}
	started := time.Now()
	first, firstErr := uc.Ejecutar(context.Background())
	second, secondErr := uc.Ejecutar(context.Background())
	if firstErr != nil || secondErr != nil || !first.Reiniciado || !second.Reiniciado || !first.FechaSimulada.Equal(second.FechaSimulada) || repo.resets != 2 || !clock.date.Equal(date) {
		t.Fatalf("first=%+v/%v second=%+v/%v resets=%d clock=%v", first, firstErr, second, secondErr, repo.resets, clock.date)
	}
	if elapsed := time.Since(started); elapsed >= 3*time.Second {
		t.Fatalf("reset took %s", elapsed)
	}
}
