package usecase

import (
	"context"
	"testing"
	"time"
)

type demoUsecaseRepo struct {
	date  time.Time
	saves int
}

func (r *demoUsecaseRepo) FechaSimulada(context.Context) (time.Time, error) { return r.date, nil }
func (r *demoUsecaseRepo) GuardarFechaSimulada(_ context.Context, date time.Time) error {
	r.date, r.saves = date, r.saves+1
	return nil
}

type demoUsecaseBus struct{ event Evento }

func (b *demoUsecaseBus) Publicar(_ context.Context, event Evento)      { b.event = event }
func (*demoUsecaseBus) Suscribir(string, func(context.Context, Evento)) {}

func TestAvanzarDemoPersistsAndPublishesOnce(t *testing.T) {
	start := time.Date(2026, 7, 26, 0, 0, 0, 0, time.UTC)
	repo, clock, bus := &demoUsecaseRepo{date: start}, NuevoRelojFake(start), &demoUsecaseBus{}
	out, err := (&AvanzarDemo{Demo: repo, Reloj: clock, Bus: bus}).Ejecutar(context.Background(), EntradaAvanzarDemo{Dias: ptrInt(2)})
	if err != nil || !out.FechaSimulada.Equal(start.AddDate(0, 0, 2)) || repo.saves != 1 || bus.event.Tipo != EvTickReloj || !clock.FechaSimulada().Equal(out.FechaSimulada) {
		t.Fatalf("out=%+v repo=%+v event=%+v clock=%v err=%v", out, repo, bus.event, clock.FechaSimulada(), err)
	}
}
func ptrInt(v int) *int { return &v }

func TestAvanzarDemoRejectsBackwardWithoutPersistence(t *testing.T) {
	start := time.Date(2026, 7, 26, 0, 0, 0, 0, time.UTC)
	repo, clock := &demoUsecaseRepo{date: start}, NuevoRelojFake(start)
	backward := start.AddDate(0, 0, -1)

	_, err := (&AvanzarDemo{Demo: repo, Reloj: clock}).Ejecutar(
		context.Background(), EntradaAvanzarDemo{Hasta: &backward},
	)
	if err != ErrTiempoSimuladoAtras {
		t.Fatalf("err=%v, want %v", err, ErrTiempoSimuladoAtras)
	}
	if repo.saves != 0 {
		t.Fatalf("saves=%d, want 0 after rejected backward advance", repo.saves)
	}
	if !repo.date.Equal(start) || !clock.FechaSimulada().Equal(start) {
		t.Fatalf("backward rejection changed state: repo=%v clock=%v", repo.date, clock.FechaSimulada())
	}
}
