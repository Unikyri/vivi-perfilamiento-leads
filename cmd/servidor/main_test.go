package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	adapterhttp "github.com/Unikyri/vivi-perfilamiento-leads/internal/adapters/http"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/infrastructure/llm"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
)

type healthFailingProvider struct{}

func (healthFailingProvider) Nombre() string { return "gemini" }
func (healthFailingProvider) GenerarTurno(context.Context, usecase.EntradaTurno) (usecase.SalidaTurno, error) {
	return usecase.SalidaTurno{}, &llm.ProviderError{Kind: llm.KindTimeout}
}
func (healthFailingProvider) ProcesarAudio(context.Context, usecase.Audio, usecase.EntradaTurno) (usecase.SalidaTurno, error) {
	return usecase.SalidaTurno{}, &llm.ProviderError{Kind: llm.KindTimeout}
}

func TestSaludReportsLiveBreakerState(t *testing.T) {
	provider := llm.NewFallbackProvider(healthFailingProvider{}, nil,
		llm.WithPrimaryRateLimiter(llm.RateLimiterFunc(func() bool { return true })))
	for i := 0; i < 3; i++ {
		_, _ = provider.GenerarTurno(context.Background(), usecase.EntradaTurno{})
	}
	handler := adapterhttp.HandlerSalud(salud{proveedorLLM: "gemini", llmProvider: provider, bd: "OK"})
	recorder := httptest.NewRecorder()
	handler(recorder, httptest.NewRequest(http.MethodGet, "/salud", nil))
	var got adapterhttp.EstadoSalud
	if err := json.NewDecoder(recorder.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.CircuitBreaker != "ABIERTO" {
		t.Fatalf("health breaker=%q", got.CircuitBreaker)
	}
	if got.FechaSimulada == "" {
		t.Fatal("health date must be present")
	}
}
