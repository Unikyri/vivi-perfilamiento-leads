package agentes

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/infrastructure/bus"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
)

type fakeCalificador struct {
	ids []string
	err error
}

func (f *fakeCalificador) Ejecutar(_ context.Context, in usecase.EntradaCalificar) (usecase.SalidaCalificar, error) {
	f.ids = append(f.ids, in.LeadID)
	return usecase.SalidaCalificar{}, f.err
}

type fakeDocumentadora struct{ ids []string }

func (f *fakeDocumentadora) Ejecutar(_ context.Context, id string) (*domain.Ficha, error) {
	f.ids = append(f.ids, id)
	return nil, nil
}

type fakeNutricionista struct{ times []time.Time }

func (f *fakeNutricionista) Ejecutar(_ context.Context, at time.Time) (int, error) {
	f.times = append(f.times, at)
	return 0, nil
}

func TestCoordinadoraDeterministicRouting(t *testing.T) {
	qualifier := &fakeCalificador{}
	doc := &fakeDocumentadora{}
	nutrition := &fakeNutricionista{}
	b := bus.Nuevo(nil)
	coord := Nueva(b, Dependencias{Calificador: qualifier, Documentadora: doc, Nutricionista: nutrition})
	coord.Registrar()
	coord.Registrar()
	for i := 0; i < 10; i++ {
		b.Publicar(context.Background(), usecase.Evento{Tipo: usecase.EvPerfilCompleto, LeadID: "lead-1"})
	}
	if len(qualifier.ids) != 10 || len(doc.ids) != 0 || len(nutrition.times) != 0 {
		t.Fatalf("qualification=%v ficha=%v ticks=%v", qualifier.ids, doc.ids, nutrition.times)
	}

	cases := []struct {
		name  string
		route any
		ficha int
	}{
		{"advisor typed route", domain.RutaAsesor, 1}, {"advisor string route", "ASESOR", 1}, {"nutrition no ficha", domain.RutaNutricion, 0}, {"invalid no ficha", "UNKNOWN", 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			before := len(doc.ids)
			b.Publicar(context.Background(), usecase.Evento{Tipo: usecase.EvRutaDecidida, LeadID: "lead-2", Payload: map[string]any{"ruta": tc.route}})
			if got := len(doc.ids) - before; got != tc.ficha {
				t.Fatalf("ficha calls=%d, want %d", got, tc.ficha)
			}
		})
	}
	at := time.Date(2026, 7, 26, 1, 0, 0, 0, time.UTC)
	b.Publicar(context.Background(), usecase.Evento{Tipo: usecase.EvTickReloj, Payload: map[string]any{"hasta": at}})
	b.Publicar(context.Background(), usecase.Evento{Tipo: usecase.EvTickReloj, Payload: map[string]any{"hasta": at.Format(time.RFC3339)}})
	if len(nutrition.times) != 2 || !nutrition.times[0].Equal(at) || !nutrition.times[1].Equal(at) {
		t.Fatalf("ticks=%v", nutrition.times)
	}
}

func TestCoordinadoraSkipsReprofileAndMalformedEvents(t *testing.T) {
	qualifier := &fakeCalificador{err: usecase.ErrLeadNoCalificable}
	doc, nutrition := &fakeDocumentadora{}, &fakeNutricionista{}
	b := bus.Nuevo(nil)
	Nueva(b, Dependencias{Calificador: qualifier, Documentadora: doc, Nutricionista: nutrition}).Registrar()
	b.Publicar(context.Background(), usecase.Evento{Tipo: usecase.EvPerfilCompleto, LeadID: "reprofile"})
	b.Publicar(context.Background(), usecase.Evento{Tipo: usecase.EvRutaDecidida, LeadID: "lead", Payload: nil})
	b.Publicar(context.Background(), usecase.Evento{Tipo: usecase.EvTickReloj, Payload: map[string]any{"hasta": "not-time"}})
	b.Publicar(context.Background(), usecase.Evento{Tipo: usecase.EvMensajeEntrante, LeadID: "ignored"})
	if len(qualifier.ids) != 1 || len(doc.ids) != 0 || len(nutrition.times) != 0 {
		t.Fatalf("qualifier=%v ficha=%v ticks=%v", qualifier.ids, doc.ids, nutrition.times)
	}
}

func TestCoordinadoraLeadNuevoIsObserveOnly(t *testing.T) {
	b := bus.Nuevo(nil)
	observed := 0
	Nueva(b, Dependencias{LeadNuevo: func(context.Context, usecase.Evento) error { observed++; return nil }}).Registrar()
	b.Publicar(context.Background(), usecase.Evento{Tipo: usecase.EvLeadNuevo, LeadID: "lead"})
	if observed != 1 {
		t.Fatalf("observed=%d, want 1", observed)
	}
}

func TestCoordinadoraNilDependenciesAndSafeMetadata(t *testing.T) {
	var logs bytes.Buffer
	qualifier := &fakeCalificador{err: errors.New("cedula 123456 sensitive name")}
	b := bus.Nuevo(slog.New(slog.NewTextHandler(&logs, nil)))
	Nueva(b, Dependencias{Calificador: qualifier, Logger: slog.New(slog.NewTextHandler(&logs, nil))}).Registrar()
	Nueva(nil, Dependencias{}).Registrar()
	b.Publicar(context.Background(), usecase.Evento{Tipo: usecase.EvPerfilCompleto, LeadID: "opaque", Payload: map[string]any{"secret": "payload"}})
	output := logs.String()
	for _, secret := range []string{"123456", "sensitive name", "payload"} {
		if strings.Contains(output, secret) {
			t.Fatalf("unsafe log contains %q: %s", secret, output)
		}
	}
	for _, field := range []string{"tipo=PerfilCompleto", "lead_id=opaque", "handler=calificar_lead", "outcome=error"} {
		if !strings.Contains(output, field) {
			t.Fatalf("log missing %q: %s", field, output)
		}
	}
}
