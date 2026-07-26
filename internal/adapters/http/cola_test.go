package http

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type planReadHTTPFake struct{ plan *domain.PlanNutricion }

func (f planReadHTTPFake) Crear(context.Context, *domain.PlanNutricion) error   { return nil }
func (f planReadHTTPFake) Guardar(context.Context, *domain.PlanNutricion) error { return nil }
func (f planReadHTTPFake) HitosVencidos(context.Context, time.Time) ([]usecase.HitoConPlan, error) {
	return nil, nil
}
func (f planReadHTTPFake) MarcarHito(context.Context, string, domain.EstadoHito) error { return nil }
func (f planReadHTTPFake) PorLead(context.Context, string) (*domain.PlanNutricion, error) {
	return f.plan, nil
}

func readRouterWithPlan(repo usecase.LeadRepository, fichas usecase.FichaRepository, plans usecase.PlanRepository) *http.ServeMux {
	c, err := NuevoControlador(Dependencias{Perfilar: &perfiladorHTTPFake{}, Leads: repo, Fichas: fichas, Planes: plans})
	if err != nil {
		panic(err)
	}
	mux := http.NewServeMux()
	c.Registrar(mux)
	return mux
}

type readLeadRepo struct {
	leads  map[string]*domain.Lead
	listed []*domain.Lead
}

func (f *readLeadRepo) Crear(context.Context, *domain.Lead) error { return nil }
func (f *readLeadRepo) Guardar(context.Context, *domain.Lead) error {
	return errors.New("unexpected write")
}
func (f *readLeadRepo) AgregarMensaje(context.Context, *domain.Mensaje) error { return nil }
func (f *readLeadRepo) Conversacion(context.Context, string) ([]domain.Mensaje, error) {
	return nil, nil
}
func (f *readLeadRepo) Listar(context.Context, usecase.FiltroLeads) ([]*domain.Lead, error) {
	return f.listed, nil
}
func (f *readLeadRepo) PorID(_ context.Context, id string) (*domain.Lead, error) {
	if lead, ok := f.leads[id]; ok {
		return lead, nil
	}
	return nil, &usecase.NotFoundError{Resource: "lead", ID: id}
}

type readFichaRepo struct {
	ficha *domain.Ficha
	err   error
}

func (f readFichaRepo) Guardar(context.Context, *domain.Ficha) error {
	return errors.New("unexpected write")
}
func (f readFichaRepo) PorLead(context.Context, string) (*domain.Ficha, error) { return f.ficha, f.err }
func readRouter(repo usecase.LeadRepository, fichas usecase.FichaRepository) *http.ServeMux {
	c, err := NuevoControlador(Dependencias{Perfilar: &perfiladorHTTPFake{}, Leads: repo, Fichas: fichas})
	if err != nil {
		panic(err)
	}
	mux := http.NewServeMux()
	c.Registrar(mux)
	return mux
}
func TestS3QueueContract(t *testing.T) {
	now := time.Date(2026, 7, 26, 7, 0, 0, 0, time.UTC)
	repo := &readLeadRepo{listed: []*domain.Lead{{LeadID: "low", Nombre: "Low", Ruta: domain.RutaAsesor, Prioridad: .2, ActualizadoEn: now}, {LeadID: "high", Nombre: "High", Ruta: domain.RutaAsesor, Afiliado: true, Prioridad: .9, ActualizadoEn: now}}}
	mux := readRouter(repo, nil)
	for _, tc := range []struct {
		path   string
		status int
		code   string
	}{{"/api/leads", 200, ""}, {"/api/leads?afiliado=true", 200, ""}, {"/api/leads?afiliado=yes", 400, "VALIDACION"}, {"/api/leads?ruta=UNKNOWN", 400, "VALIDACION"}} {
		r := httptest.NewRecorder()
		mux.ServeHTTP(r, httptest.NewRequest(http.MethodGet, tc.path, nil))
		if r.Code != tc.status || !strings.HasPrefix(r.Header().Get("Content-Type"), "application/json; charset=utf-8") {
			t.Fatalf("%s response=%d %s", tc.path, r.Code, r.Body)
		}
		if tc.code != "" && !strings.Contains(r.Body.String(), `"codigo":"`+tc.code+`"`) {
			t.Fatalf("%s error=%s", tc.path, r.Body)
		}
		if tc.path == "/api/leads" {
			var got usecase.ColaLeads
			if err := json.NewDecoder(r.Body).Decode(&got); err != nil || got.Leads[0].LeadID != "high" || got.Cupo10.Usados != 1 {
				t.Fatalf("queue=%+v err=%v", got, err)
			}
		}
	}
}
func TestS3DetailAndFichaContract(t *testing.T) {
	lead := &domain.Lead{LeadID: "lead-1", Nombre: "Ana", Cedula: "secret-cedula", Telefono: "secret-phone", Afiliado: true, Estado: domain.EstadoLeadCalificado, Ruta: domain.RutaAsesor, Prioridad: .8}
	repo := &readLeadRepo{leads: map[string]*domain.Lead{"lead-1": lead}}
	ficha := &domain.Ficha{FichaID: "f-1", LeadID: "lead-1"}
	mux := readRouterWithPlan(repo, readFichaRepo{ficha: ficha}, planReadHTTPFake{plan: &domain.PlanNutricion{PlanID: "plan-1", LeadID: "lead-1", Estado: domain.EstadoPlanActivo}})
	get := func(path string) *httptest.ResponseRecorder {
		r := httptest.NewRecorder()
		mux.ServeHTTP(r, httptest.NewRequest(http.MethodGet, path, nil))
		return r
	}
	r := get("/api/leads/lead-1")
	if r.Code != http.StatusOK || strings.Contains(r.Body.String(), "secret-") || !strings.Contains(r.Body.String(), `"semaforo":"VERDE"`) || !strings.Contains(r.Body.String(), `"plan":{"plan_id":"plan-1"`) || strings.Contains(r.Body.String(), `"cedula"`) || strings.Contains(r.Body.String(), `"telefono"`) {
		t.Fatalf("detail=%d %s", r.Code, r.Body)
	}
	r = get("/api/leads/lead-1/ficha")
	if r.Code != http.StatusOK || !strings.Contains(r.Body.String(), `"ficha_id":"f-1"`) {
		t.Fatalf("ficha=%d %s", r.Code, r.Body)
	}
	missing := readRouter(repo, readFichaRepo{err: &usecase.NotFoundError{Resource: "ficha", ID: "lead-1"}})
	r = httptest.NewRecorder()
	missing.ServeHTTP(r, httptest.NewRequest(http.MethodGet, "/api/leads/lead-1/ficha", nil))
	if r.Code != http.StatusNotFound || !strings.Contains(r.Body.String(), `"codigo":"FICHA_NO_DISPONIBLE"`) {
		t.Fatalf("missing ficha=%d %s", r.Code, r.Body)
	}
	r = httptest.NewRecorder()
	missing.ServeHTTP(r, httptest.NewRequest(http.MethodGet, "/api/leads/nope/ficha", nil))
	if r.Code != http.StatusNotFound || !strings.Contains(r.Body.String(), `"codigo":"LEAD_NO_ENCONTRADO"`) {
		t.Fatalf("missing lead=%d %s", r.Code, r.Body)
	}
}
