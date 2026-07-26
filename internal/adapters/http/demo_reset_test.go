package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
)

type resetHTTPFake struct {
	date   time.Time
	resets int
}

func (f *resetHTTPFake) Reiniciar(context.Context) (time.Time, error) {
	f.resets++
	return f.date, nil
}

type resetHTTPClock struct{ date time.Time }

func (c *resetHTTPClock) Ahora() time.Time         { return c.date }
func (c *resetHTTPClock) FechaSimulada() time.Time { return c.date }
func (c *resetHTTPClock) Avanzar(date time.Time)   { c.date = date }

func resetRouter(repo usecase.DemoResetRepository, clock usecase.Reloj, enabled bool) *http.ServeMux {
	controller, _ := NuevoControlador(Dependencias{
		Perfilar: &perfiladorHTTPFake{}, Leads: &leadHTTPFake{},
		ReiniciarDemo: &usecase.ReiniciarDemo{Repository: repo, Reloj: clock, Habilitado: enabled},
	})
	mux := http.NewServeMux()
	controller.Registrar(mux)
	return mux
}

func TestDemoResetDisabledIsGenericAndReadOnly(t *testing.T) {
	repo := &resetHTTPFake{date: time.Date(2026, 7, 26, 0, 0, 0, 0, time.UTC)}
	r := httptest.NewRecorder()
	resetRouter(repo, &resetHTTPClock{date: repo.date}, false).ServeHTTP(r, httptest.NewRequest(http.MethodPost, "/api/demo/reiniciar", strings.NewReader(`{}`)))
	var body errorEnvelope
	_ = json.NewDecoder(r.Body).Decode(&body)
	if r.Code != http.StatusInternalServerError || body.Error.Codigo != "ERROR_INTERNO" || repo.resets != 0 || strings.Contains(r.Body.String(), "DEMO_SEED") {
		t.Fatalf("code=%d body=%s resets=%d", r.Code, r.Body, repo.resets)
	}
}

func TestDemoResetReturnsSameDateTwice(t *testing.T) {
	date := time.Date(2026, 7, 26, 0, 0, 0, 0, time.UTC)
	repo, clock := &resetHTTPFake{date: date}, &resetHTTPClock{date: date.AddDate(0, 0, 2)}
	mux := resetRouter(repo, clock, true)
	started := time.Now()
	var responses []reiniciarResponse
	for range 2 {
		r := httptest.NewRecorder()
		mux.ServeHTTP(r, httptest.NewRequest(http.MethodPost, "/api/demo/reiniciar", strings.NewReader(`{}`)))
		var response reiniciarResponse
		if r.Code != http.StatusOK || json.NewDecoder(r.Body).Decode(&response) != nil || !response.Reiniciado {
			t.Fatalf("code=%d body=%s", r.Code, r.Body)
		}
		responses = append(responses, response)
	}
	if responses[0].FechaSimulada != "2026-07-26" || responses[0].FechaSimulada != responses[1].FechaSimulada || repo.resets != 2 || time.Since(started) >= 3*time.Second {
		t.Fatalf("responses=%+v resets=%d", responses, repo.resets)
	}
}
