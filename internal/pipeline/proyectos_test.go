package pipeline

import (
	"testing"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
)

func TestConstruirProyectosCalculaPrecios(t *testing.T) {
	mapa := []EntradaMapa{{ProyectoID: "mongui", Nombre: "Monguí", NombreV2: "Agrupación De Vivienda Monguí"}}
	cs := []domain.Comprador{
		{Proyecto: "Agrupación De Vivienda Monguí", ValorCOP: 156470000},
		{Proyecto: "Agrupación De Vivienda Monguí", ValorCOP: 204194000},
	}
	got, err := ConstruirProyectos(mapa, cs)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if got[0].PrecioDesde != 156470000 || got[0].PrecioHasta != 204194000 {
		t.Errorf("precios = %d–%d, esperado 156470000–204194000", got[0].PrecioDesde, got[0].PrecioHasta)
	}
	if !got[0].EsVIS {
		t.Error("Monguí debe ser VIS")
	}
}

func TestProyectoNoVIS(t *testing.T) {
	mapa := []EntradaMapa{{ProyectoID: "araucaria", Nombre: "Araucaria", NombreV2: "ARAUCARIA"}}
	cs := []domain.Comprador{{Proyecto: "ARAUCARIA", ValorCOP: 562000000}}
	got, _ := ConstruirProyectos(mapa, cs)
	if got[0].EsVIS {
		t.Error("Araucaria ($562M) NO es VIS (tope 262.635.750)")
	}
}

func TestFallaSiNombreV2NoExiste(t *testing.T) {
	mapa := []EntradaMapa{{ProyectoID: "x", Nombre: "X", NombreV2: "NOMBRE QUE NO EXISTE"}}
	if _, err := ConstruirProyectos(mapa, []domain.Comprador{}); err == nil {
		t.Fatal("debe fallar si un nombre_v2 no aparece en el dataset")
	}
}




