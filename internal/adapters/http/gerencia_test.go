package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
)

type buyerPersonaHTTPFake struct {
	projects map[string]domain.Proyecto
	buyers   []domain.Comprador
}

func (f buyerPersonaHTTPFake) Proyectos(context.Context) (map[string]domain.Proyecto, error) {
	return f.projects, nil
}
func (f buyerPersonaHTTPFake) Compradores(context.Context) ([]domain.Comprador, error) {
	return f.buyers, nil
}
func (f buyerPersonaHTTPFake) AfiliadoPorCedula(context.Context, string) (*usecase.Afiliado, error) {
	return nil, nil
}
func (f buyerPersonaHTTPFake) BrochureMarkdown(context.Context, string) (string, error) {
	return "", nil
}

func buyerPersonaRouter(catalog usecase.CatalogoRepository) *http.ServeMux {
	controller, _ := NuevoControlador(Dependencias{Perfilar: &perfiladorHTTPFake{}, Leads: &leadHTTPFake{}, Catalogo: catalog})
	mux := http.NewServeMux()
	controller.Registrar(mux)
	return mux
}

func TestGerenciaBuyerPersonaProjectAndCatalogContract(t *testing.T) {
	catalog := buyerPersonaHTTPFake{
		projects: map[string]domain.Proyecto{"p1": {ProyectoID: "p1", Nombre: "Proyecto Uno"}, "p2": {ProyectoID: "p2", Nombre: "Proyecto Dos"}},
		buyers:   []domain.Comprador{{ID: 1, ProyectoID: "p1", Afiliado: true, Categoria: "A", RangoEdad: "20-35"}},
	}
	mux := buyerPersonaRouter(catalog)
	for _, path := range []string{"/api/gerencia/buyer-persona?proyecto_id=p1", "/api/gerencia/buyer-persona"} {
		r := httptest.NewRecorder()
		mux.ServeHTTP(r, httptest.NewRequest(http.MethodGet, path, nil))
		if r.Code != http.StatusOK {
			t.Fatalf("%s response=%d %s", path, r.Code, r.Body)
		}
		if path[len(path)-1] == '1' {
			var got usecase.BuyerPersonaResumen
			if err := json.NewDecoder(r.Body).Decode(&got); err != nil || got.ProyectoID != "p1" || got.Afiliacion.Afiliados != 1 {
				t.Fatalf("project=%+v err=%v", got, err)
			}
		} else {
			var got usecase.BuyerPersonaCatalogo
			if err := json.NewDecoder(r.Body).Decode(&got); err != nil || len(got.Proyectos) != 2 || got.Proyectos[0].ProyectoID != "p1" {
				t.Fatalf("catalog=%+v err=%v", got, err)
			}
		}
	}
}

func TestGerenciaBuyerPersonaRejectsInvalidFilter(t *testing.T) {
	mux := buyerPersonaRouter(buyerPersonaHTTPFake{})
	for _, path := range []string{"/api/gerencia/buyer-persona?proyecto_id=", "/api/gerencia/buyer-persona?proyecto_id=p1&proyecto_id=p2"} {
		r := httptest.NewRecorder()
		mux.ServeHTTP(r, httptest.NewRequest(http.MethodGet, path, nil))
		if r.Code != http.StatusBadRequest || !strings.Contains(r.Body.String(), `"codigo":"VALIDACION"`) {
			t.Fatalf("invalid filter=%s response=%d %s", path, r.Code, r.Body)
		}
	}
}
