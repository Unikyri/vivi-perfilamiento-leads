package usecase

import (
	"context"
	"reflect"
	"testing"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
)

type buyerPersonaCatalogFake struct {
	projects map[string]domain.Proyecto
	buyers   []domain.Comprador
}

func (f *buyerPersonaCatalogFake) Proyectos(context.Context) (map[string]domain.Proyecto, error) {
	return f.projects, nil
}
func (f *buyerPersonaCatalogFake) Compradores(context.Context) ([]domain.Comprador, error) {
	return f.buyers, nil
}
func (f *buyerPersonaCatalogFake) AfiliadoPorCedula(context.Context, string) (*Afiliado, error) {
	return nil, nil
}
func (f *buyerPersonaCatalogFake) BrochureMarkdown(context.Context, string) (string, error) {
	return "", nil
}

func TestBuyerPersonaProjectIsDeterministicAndReadOnly(t *testing.T) {
	buyers := []domain.Comprador{
		{ID: 1, ProyectoID: "p1", Afiliado: true, Categoria: "A", RangoEdad: "20-35", FechaOpcion: "2024-01-01"},
		{ID: 2, ProyectoID: "p1", Categoria: "B", RangoEdad: "36-45", Desistio: true, FechaOpcion: "2024-02-01"},
		{ID: 3, ProyectoID: "p2", Categoria: "C", RangoEdad: "55+"},
	}
	before := append([]domain.Comprador(nil), buyers...)
	uc := &BuyerPersona{Catalogo: &buyerPersonaCatalogFake{
		projects: map[string]domain.Proyecto{"p1": {ProyectoID: "p1", Nombre: "Uno"}},
		buyers:   buyers,
	}}
	got, err := uc.Proyecto(context.Background(), "p1")
	if err != nil {
		t.Fatal(err)
	}
	if got.Muestras != 2 || got.Afiliacion.Afiliados != .5 || got.Afiliacion.NoAfiliados != .5 || got.TasaDesistimiento != .5 {
		t.Fatalf("summary=%+v", got)
	}
	if got.Categoria["A"] != .5 || got.Categoria["B"] != .5 || got.RangoEdad["20-35"] != .5 || got.RangoEdad["36-45"] != .5 {
		t.Fatalf("proportions=%+v %+v", got.Categoria, got.RangoEdad)
	}
	if got.ActualizadoEn != "2024-02-01T00:00:00Z" || !reflect.DeepEqual(before, buyers) {
		t.Fatalf("date or source mutation: %q %#v", got.ActualizadoEn, buyers)
	}
	again, err := uc.Proyecto(context.Background(), "p1")
	if err != nil || !reflect.DeepEqual(got, again) {
		t.Fatalf("non-deterministic result: %+v %+v %v", got, again, err)
	}
}

func TestBuyerPersonaCatalogIsSortedAndIncludesEmptyProjects(t *testing.T) {
	uc := &BuyerPersona{Catalogo: &buyerPersonaCatalogFake{
		projects: map[string]domain.Proyecto{
			"z": {ProyectoID: "z", Nombre: "Z"}, "a": {ProyectoID: "a", Nombre: "A"},
		},
		buyers: []domain.Comprador{{ID: 1, ProyectoID: "z", Categoria: "SIN_DATO", RangoEdad: "SIN_DATO"}},
	}}
	got, err := uc.CatalogoCompleto(context.Background())
	if err != nil || len(got.Proyectos) != 2 || got.Proyectos[0].ProyectoID != "a" || got.Proyectos[1].Muestras != 1 {
		t.Fatalf("catalog=%+v err=%v", got, err)
	}
	if got.Proyectos[0].ActualizadoEn != "0001-01-01T00:00:00Z" || got.Proyectos[0].Categoria["A"] != 0 {
		t.Fatalf("empty aggregate=%+v", got.Proyectos[0])
	}
}
