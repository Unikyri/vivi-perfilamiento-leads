package llm

import (
	"context"
	"encoding/json"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
	"io"
	"sync"
	"time"
)

type metricEvent struct {
	LeadID   string `json:"lead_id"`
	Event    string `json:"event"`
	Latency  int64  `json:"latency_ms"`
	Provider string `json:"provider"`
	Outcome  string `json:"outcome"`
}
type Metricas struct {
	next   usecase.LLMProvider
	writer io.Writer
	clock  Clock
	mu     sync.Mutex
}

func ConMetricas(next usecase.LLMProvider, writer io.Writer, clock Clock) *Metricas {
	if writer == nil {
		writer = io.Discard
	}
	if clock == nil {
		clock = wallClock{}
	}
	return &Metricas{next: next, writer: writer, clock: clock}
}
func (m *Metricas) Nombre() string {
	if m == nil || m.next == nil {
		return "unconfigured"
	}
	return m.next.Nombre()
}
func (m *Metricas) GenerarTurno(ctx context.Context, in usecase.EntradaTurno) (usecase.SalidaTurno, error) {
	start := m.clock.Now()
	if m.next == nil {
		return m.record(in, "generar_turno", start, usecase.SalidaTurno{}, providerError(KindConfig, nil))
	}
	out, err := m.next.GenerarTurno(ctx, in)
	return m.record(in, "generar_turno", start, out, err)
}
func (m *Metricas) ProcesarAudio(ctx context.Context, audio usecase.Audio, in usecase.EntradaTurno) (usecase.SalidaTurno, error) {
	start := m.clock.Now()
	if m.next == nil {
		return m.record(in, "procesar_audio", start, usecase.SalidaTurno{}, providerError(KindConfig, nil))
	}
	out, err := m.next.ProcesarAudio(ctx, audio, in)
	return m.record(in, "procesar_audio", start, out, err)
}
func (m *Metricas) Observe(provider string, kind ErrorKind) {
	outcome := string(kind)
	if outcome == "" {
		outcome = "ok"
	}
	m.emit(metricEvent{Event: "provider_attempt", Provider: provider, Outcome: outcome})
}
func (m *Metricas) CircuitBreakerState() BreakerState {
	if owner, ok := m.next.(breakerStateOwner); ok {
		return owner.CircuitBreakerState()
	}
	return BreakerClosed
}
func (m *Metricas) record(in usecase.EntradaTurno, event string, start time.Time, out usecase.SalidaTurno, err error) (usecase.SalidaTurno, error) {
	outcome := "accepted"
	if err != nil {
		outcome = string(ErrorKindOf(err))
		if outcome == "" {
			outcome = "unknown"
		}
		if ErrorKindOf(err) == KindCircuitOpen {
			outcome = "breaker_open"
		}
	} else if out.Accion == usecase.AccionFueraDeDominio {
		outcome = "rejected"
	}
	latency := m.clock.Now().Sub(start).Milliseconds()
	if latency < 0 {
		latency = 0
	}
	m.emit(metricEvent{LeadID: in.LeadID, Event: event, Latency: latency, Provider: m.Nombre(), Outcome: outcome})
	return out, err
}
func (m *Metricas) emit(event metricEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	_, _ = m.writer.Write(append(data, '\n'))
}
