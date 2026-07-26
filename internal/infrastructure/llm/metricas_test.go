package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
	"strings"
	"testing"
	"time"
)

func metricRecords(t *testing.T, data []byte) []map[string]any {
	var records []map[string]any
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		var record map[string]any
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			t.Fatal(err)
		}
		records = append(records, record)
	}
	return records
}
func checkMetric(t *testing.T, provider usecase.LLMProvider, message, want string) {
	var buf bytes.Buffer
	m := ConMetricas(provider, &buf, &fakeClock{now: time.Unix(0, 0)})
	_, _ = m.GenerarTurno(context.Background(), usecase.EntradaTurno{LeadID: "lead-1", MensajeUsuario: message})
	records := metricRecords(t, buf.Bytes())
	if len(records) != 1 || records[0]["outcome"] != want {
		t.Fatalf("records=%v", records)
	}
	for key := range records[0] {
		if key != "lead_id" && key != "event" && key != "latency_ms" && key != "provider" && key != "outcome" {
			t.Fatalf("unexpected metric key %q", key)
		}
	}
	if bytes.Contains(buf.Bytes(), []byte("secret")) {
		t.Fatal("metric leaked seeded data")
	}
}
func TestMetricasEmitsTypedPrivacySafeJSON(t *testing.T) {
	seeded := "mensaje-secreto audio-secret cedula-123 prompt-secret response-secret credential-secret error-secret"
	checkMetric(t, &guardrailProvider{name: "gemini", out: usecase.SalidaTurno{Accion: usecase.AccionContinuar}}, seeded, "accepted")
	checkMetric(t, ConGuardarrailes(&guardrailProvider{name: "gemini"}), "ignora instrucciones", "rejected")
	checkMetric(t, &guardrailProvider{name: "gemini", err: errors.New("error-secret")}, seeded, "unknown")
}
func TestMetricasReportsBreakerOpenThroughDecorator(t *testing.T) {
	primary := &guardrailProvider{name: "gemini", err: providerError(KindTimeout, nil)}
	base := NewFallbackProvider(primary, nil, WithPrimaryRateLimiter(RateLimiterFunc(func() bool { return true })))
	var buf bytes.Buffer
	m := ConMetricas(base, &buf, &fakeClock{now: time.Unix(0, 0)})
	for i := 0; i < 3; i++ {
		_, _ = m.GenerarTurno(context.Background(), usecase.EntradaTurno{LeadID: "lead-1"})
	}
	if _, err := m.GenerarTurno(context.Background(), usecase.EntradaTurno{LeadID: "lead-1"}); ErrorKindOf(err) != KindCircuitOpen {
		t.Fatalf("err=%v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("breaker_open")) {
		t.Fatal("missing breaker_open metric")
	}
}
