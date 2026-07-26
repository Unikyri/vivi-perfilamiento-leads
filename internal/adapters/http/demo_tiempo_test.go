package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/infrastructure/bus"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
)

type demoHTTPRepo struct {
	date  time.Time
	saves int
}

func (r *demoHTTPRepo) FechaSimulada(context.Context) (time.Time, error) { return r.date, nil }
func (r *demoHTTPRepo) GuardarFechaSimulada(_ context.Context, date time.Time) error {
	r.date, r.saves = date, r.saves+1
	return nil
}

type demoHTTPClock struct{ date time.Time }

func (c *demoHTTPClock) Ahora() time.Time         { return c.date }
func (c *demoHTTPClock) FechaSimulada() time.Time { return c.date }
func (c *demoHTTPClock) Avanzar(date time.Time)   { c.date = date }

func demoRouter(repo usecase.DemoRepository, clock usecase.Reloj, events *int) *http.ServeMux {
	b := bus.Nuevo(nil)
	b.Suscribir(usecase.EvTickReloj, func(context.Context, usecase.Evento) { *events++ })
	c, _ := NuevoControlador(Dependencias{Perfilar: &perfiladorHTTPFake{}, Leads: &leadHTTPFake{}, AvanzarDemo: &usecase.AvanzarDemo{Demo: repo, Reloj: clock, Bus: b}})
	mux := http.NewServeMux()
	c.Registrar(mux)
	return mux
}
func TestDemoTiempoAdvancesAndCountsOneTick(t *testing.T) {
	start := time.Date(2026, 7, 26, 0, 0, 0, 0, time.UTC)
	repo, clock, events := &demoHTTPRepo{date: start}, &demoHTTPClock{date: start}, 0
	r := httptest.NewRecorder()
	body := `{"avanzar_dias":2}`
	demoRouter(repo, clock, &events).ServeHTTP(r, httptest.NewRequest(http.MethodPost, "/api/demo/tiempo", strings.NewReader(body)))
	var got tiempoResponse
	raw := r.Body.String()
	if r.Code != http.StatusOK || json.NewDecoder(strings.NewReader(raw)).Decode(&got) != nil || got.HitosDisparados != 0 || got.FechaSimulada != "2026-07-28" || strings.Contains(raw, "T00:00:00") || repo.saves != 1 || events != 1 {
		t.Fatalf("code=%d got=%+v repo=%+v events=%d", r.Code, got, repo, events)
	}
}
func TestDemoTiempoRejectsNeitherOrBoth(t *testing.T) {
	start := time.Date(2026, 7, 26, 0, 0, 0, 0, time.UTC)
	repo, clock := &demoHTTPRepo{date: start}, &demoHTTPClock{date: start}
	for _, body := range []string{`{}`, `{"avanzar_dias":1,"avanzar_hasta":"2026-07-27T00:00:00Z"}`} {
		r := httptest.NewRecorder()
		demoRouter(repo, clock, new(int)).ServeHTTP(r, httptest.NewRequest(http.MethodPost, "/api/demo/tiempo", strings.NewReader(body)))
		if r.Code != http.StatusBadRequest || !strings.Contains(r.Body.String(), `"codigo":"VALIDACION"`) {
			t.Fatalf("body=%s code=%d", body, r.Code)
		}
	}
}
