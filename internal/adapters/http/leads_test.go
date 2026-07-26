package http

import (
	"context"
	"encoding/json"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type perfiladorHTTPFake struct {
	salida  usecase.SalidaPerfilar
	entrada usecase.EntradaPerfilar
}

func (f *perfiladorHTTPFake) Ejecutar(_ context.Context, in usecase.EntradaPerfilar) (usecase.SalidaPerfilar, error) {
	f.entrada = in
	return f.salida, nil
}

type leadHTTPFake struct {
	usecase.LeadRepository
	lead     *domain.Lead
	messages []domain.Mensaje
}

func (f *leadHTTPFake) PorID(_ context.Context, id string) (*domain.Lead, error) {
	if f.lead == nil || f.lead.LeadID != id {
		return nil, &usecase.NotFoundError{Resource: "lead", ID: id}
	}
	return f.lead, nil
}
func (f *leadHTTPFake) Conversacion(context.Context, string) ([]domain.Mensaje, error) {
	return f.messages, nil
}
func newRouter(p Perfilador, repo usecase.LeadRepository) *http.ServeMux {
	c, _ := NuevoControlador(Dependencias{Perfilar: p, Leads: repo})
	mux := http.NewServeMux()
	c.Registrar(mux)
	return mux
}
func TestS1CreateLeadContract(t *testing.T) {
	cases := []struct {
		name, body string
		status     int
		code       string
	}{{"seed", `{"precargado_id":"ana","nombre":"private","cedula":"secret"}`, 201, ""}, {"invalid source", `{"nombre":"Ana","fuente":"UNKNOWN"}`, 400, "VALIDACION"}, {"strict JSON", `{"nombre":"Ana","extra":"secret"}`, 400, "VALIDACION"}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := &perfiladorHTTPFake{salida: usecase.SalidaPerfilar{LeadID: "lead-1", Estado: domain.EstadoLeadPerfilando, AfiliadoDetectado: true}}
			r := httptest.NewRecorder()
			newRouter(p, &leadHTTPFake{}).ServeHTTP(r, httptest.NewRequest(http.MethodPost, "/api/leads", strings.NewReader(tc.body)))
			if r.Code != tc.status || !strings.HasPrefix(r.Header().Get("Content-Type"), "application/json; charset=utf-8") {
				t.Fatalf("response=%d %s", r.Code, r.Body)
			}
			body := r.Body.String()
			if tc.code != "" {
				var e errorEnvelope
				_ = json.NewDecoder(strings.NewReader(body)).Decode(&e)
				if e.Error.Codigo != tc.code || strings.Contains(body, "secret") {
					t.Fatalf("unsafe error: %s", body)
				}
			} else if p.entrada.Cedula != "1032456789" {
				t.Fatalf("seed cedula=%q", p.entrada.Cedula)
			}
		})
	}
}
func TestS1ConversationContract(t *testing.T) {
	repo := &leadHTTPFake{lead: &domain.Lead{LeadID: "lead-1", Estado: domain.EstadoLeadPerfilando}, messages: []domain.Mensaje{{MensajeID: "m1", Autor: domain.AutorMensajeVivi, TipoContenido: domain.TipoContenidoTexto, Texto: "hola", CreadoEn: time.Unix(1, 0).UTC()}}}
	r := httptest.NewRecorder()
	newRouter(&perfiladorHTTPFake{}, repo).ServeHTTP(r, httptest.NewRequest(http.MethodGet, "/api/leads/lead-1/conversacion", nil))
	var got conversacionResponse
	if r.Code != 200 || json.NewDecoder(r.Body).Decode(&got) != nil || got.LeadID != "lead-1" || got.TurnoEnProceso || len(got.Mensajes) != 1 {
		t.Fatalf("conversation=%d %+v", r.Code, got)
	}
}
func TestS1ConversationMissingIsPrivate404(t *testing.T) {
	r := httptest.NewRecorder()
	newRouter(&perfiladorHTTPFake{}, &leadHTTPFake{}).ServeHTTP(r, httptest.NewRequest(http.MethodGet, "/api/leads/nope/conversacion", nil))
	if r.Code != 404 || strings.Contains(r.Body.String(), "nope") {
		t.Fatalf("private 404: %d %s", r.Code, r.Body)
	}
}
