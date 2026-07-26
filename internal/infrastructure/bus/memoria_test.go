package bus

import (
	"bytes"
	"context"
	"log/slog"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
)

func TestEnMemoriaDelivery(t *testing.T) {
	tests := []struct {
		name  string
		setup func(*EnMemoria, *testing.T) ([]string, context.Context)
		want  []string
	}{
		{
			name: "standalone registration order and synchronous context",
			setup: func(b *EnMemoria, _ *testing.T) ([]string, context.Context) {
				ctx := context.WithValue(context.Background(), "marker", "same")
				var got []string
				b.Suscribir(usecase.EvLeadNuevo, func(received context.Context, _ usecase.Evento) {
					if received == ctx {
						got = append(got, "first")
					}
				})
				b.Suscribir(usecase.EvLeadNuevo, func(received context.Context, _ usecase.Evento) {
					if received.Value("marker") == "same" {
						got = append(got, "second")
					}
				})
				b.Publicar(ctx, usecase.Evento{Tipo: usecase.EvLeadNuevo, LeadID: "opaque-1"})
				return got, ctx
			},
			want: []string{"first", "second"},
		},
		{
			name: "unknown and zero subscribers are silent",
			setup: func(b *EnMemoria, _ *testing.T) ([]string, context.Context) {
				b.Publicar(context.Background(), usecase.Evento{Tipo: "Unknown"})
				b.Publicar(context.Background(), usecase.Evento{Tipo: usecase.EvLeadNuevo})
				return nil, nil
			},
		},
		{
			name: "nil handler is ignored",
			setup: func(b *EnMemoria, _ *testing.T) ([]string, context.Context) {
				b.Suscribir(usecase.EvLeadNuevo, nil)
				return nil, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := tt.setup(Nuevo(nil), t)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("callbacks = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnMemoriaNestedPublishAndSubscribe(t *testing.T) {
	b := Nuevo(nil)
	var got []string
	nested := false
	b.Suscribir(usecase.EvLeadNuevo, func(ctx context.Context, _ usecase.Evento) {
		got = append(got, "outer-first")
		if nested {
			return
		}
		nested = true
		b.Suscribir(usecase.EvLeadNuevo, func(context.Context, usecase.Evento) { got = append(got, "new") })
		b.Publicar(ctx, usecase.Evento{Tipo: usecase.EvLeadNuevo, LeadID: "nested"})
	})
	b.Suscribir(usecase.EvLeadNuevo, func(context.Context, usecase.Evento) { got = append(got, "outer-second") })
	done := make(chan struct{})
	go func() {
		b.Publicar(context.Background(), usecase.Evento{Tipo: usecase.EvLeadNuevo, LeadID: "outer"})
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("nested publication deadlocked")
	}
	want := []string{"outer-first", "outer-first", "outer-second", "new", "outer-second"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("callbacks = %v, want %v", got, want)
	}
}

func TestEnMemoriaDeepCopiesProducerAndSubscribers(t *testing.T) {
	b := Nuevo(nil)
	var second usecase.Evento
	b.Suscribir(usecase.EvRutaDecidida, func(_ context.Context, event usecase.Evento) {
		event.Payload["nested"].(map[string]any)["name"] = "changed"
		event.Payload["recommendations"].([]domain.Recomendacion)[0].Nombre = "changed"
		event.Payload["nested"].(map[string]any)["added"] = true
	})
	b.Suscribir(usecase.EvRutaDecidida, func(_ context.Context, event usecase.Evento) { second = event })
	original := usecase.Evento{Tipo: usecase.EvRutaDecidida, LeadID: "opaque-2", Payload: map[string]any{
		"nested":          map[string]any{"name": "original"},
		"recommendations": []domain.Recomendacion{{Nombre: "project"}},
	}}
	b.Publicar(context.Background(), original)
	if second.Payload["nested"].(map[string]any)["name"] != "original" || second.Payload["nested"].(map[string]any)["added"] != nil {
		t.Fatalf("second handler observed first handler mutation: %#v", second.Payload)
	}
	if second.Payload["recommendations"].([]domain.Recomendacion)[0].Nombre != "project" {
		t.Fatal("second handler observed nested slice mutation")
	}
	if original.Payload["nested"].(map[string]any)["name"] != "original" || original.Payload["recommendations"].([]domain.Recomendacion)[0].Nombre != "project" {
		t.Fatal("producer event was mutated")
	}
}

func TestEnMemoriaRecoversPerHandlerAndLogsSafeMetadata(t *testing.T) {
	var logs bytes.Buffer
	b := Nuevo(slog.New(slog.NewTextHandler(&logs, nil)))
	called := false
	b.Suscribir(usecase.EvLeadNuevo, func(context.Context, usecase.Evento) {
		panic("cedula 123456 nombre Ana")
	})
	b.Suscribir(usecase.EvLeadNuevo, func(context.Context, usecase.Evento) { called = true })
	b.Publicar(context.Background(), usecase.Evento{Tipo: usecase.EvLeadNuevo, LeadID: "opaque-3", Payload: map[string]any{
		"nombre": "Ana", "cedula": "123456", "mensaje": "sensitive text",
	}})
	if !called {
		t.Fatal("handler after panic was not called")
	}
	output := logs.String()
	for _, secret := range []string{"123456", "Ana", "sensitive text", "cedula"} {
		if strings.Contains(output, secret) {
			t.Fatalf("log contains sensitive value %q: %s", secret, output)
		}
	}
	for _, field := range []string{"tipo=LeadNuevo", "lead_id=opaque-3", "handler=0", "outcome=panic", "outcome=ok"} {
		if !strings.Contains(output, field) {
			t.Fatalf("log missing safe field %q: %s", field, output)
		}
	}
}
