package motor

import (
	"math"
	"reflect"
	"testing"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
)

func TestRecomendarProyectos(t *testing.T) {
	project := func(id string, price int64) domain.Proyecto {
		return domain.Proyecto{
			ProyectoID:      id,
			Nombre:          "Nombre " + id,
			Zona:            "Zona " + id,
			PrecioDesde:     price,
			PrecioHasta:     price + 10,
			EsVIS:           true,
			BrochureURL:     "https://example.test/" + id,
			Recorrido360URL: "https://example.test/360/" + id,
		}
	}

	tests := []struct {
		name        string
		vecinos     []Vecino
		catalogo    map[string]domain.Proyecto
		presupuesto int64
		wantIDs     []string
		wantCounts  []int
		wantRates   []float64
	}{
		{
			name:        "empty input returns non nil empty output",
			vecinos:     nil,
			catalogo:    map[string]domain.Proyecto{},
			presupuesto: 100,
			wantIDs:     []string{},
		},
		{
			name: "aggregates canonical IDs and ignores orphan or empty IDs",
			vecinos: []Vecino{
				{ProyectoID: "P-1", Desistio: true},
				{ProyectoID: "P-1"},
				{ProyectoID: "ORPHAN", Desistio: true},
				{ProyectoID: ""},
				{ProyectoID: "P-2"},
			},
			catalogo: map[string]domain.Proyecto{
				"P-1": project("P-1", 100),
				"P-2": project("P-2", 100),
			},
			presupuesto: 100,
			wantIDs:     []string{"P-1", "P-2"},
			wantCounts:  []int{2, 1},
			wantRates:   []float64{0.5, 0},
		},
		{
			name: "count then exact rate then canonical ID order",
			vecinos: []Vecino{
				{ProyectoID: "P-3", Desistio: true},
				{ProyectoID: "P-1"},
				{ProyectoID: "P-2"},
				{ProyectoID: "P-3"},
				{ProyectoID: "P-1", Desistio: true},
				{ProyectoID: "P-2"},
			},
			catalogo: map[string]domain.Proyecto{
				"P-1": project("P-1", 100),
				"P-2": project("P-2", 100),
				"P-3": project("P-3", 100),
			},
			presupuesto: 100,
			wantIDs:     []string{"P-2", "P-1", "P-3"},
			wantCounts:  []int{2, 2, 2},
			wantRates:   []float64{0, 0.5, 0.5},
		},
		{
			name: "caps output at three with shuffled input stable ordering",
			vecinos: []Vecino{
				{ProyectoID: "P-5"}, {ProyectoID: "P-3"}, {ProyectoID: "P-1"},
				{ProyectoID: "P-4"}, {ProyectoID: "P-2"},
			},
			catalogo: map[string]domain.Proyecto{
				"P-1": project("P-1", 100), "P-2": project("P-2", 100),
				"P-3": project("P-3", 100), "P-4": project("P-4", 100),
				"P-5": project("P-5", 100),
			},
			presupuesto: 100,
			wantIDs:     []string{"P-1", "P-2", "P-3"},
			wantCounts:  []int{1, 1, 1},
			wantRates:   []float64{0, 0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RecomendarProyectos(tt.vecinos, tt.catalogo, tt.presupuesto)
			if got == nil {
				t.Fatal("result is nil, want non-nil slice")
			}
			if len(got) != len(tt.wantIDs) {
				t.Fatalf("len(result) = %d, want %d: %#v", len(got), len(tt.wantIDs), got)
			}
			for i, wantID := range tt.wantIDs {
				if got[i].ProyectoID != wantID {
					t.Errorf("result[%d].ProyectoID = %q, want %q", i, got[i].ProyectoID, wantID)
				}
				if len(tt.wantCounts) > 0 && got[i].Vecinos != tt.wantCounts[i] {
					t.Errorf("result[%d].Vecinos = %d, want %d", i, got[i].Vecinos, tt.wantCounts[i])
				}
				if len(tt.wantRates) > 0 && math.Abs(got[i].TasaDesistimiento-tt.wantRates[i]) > 1e-12 {
					t.Errorf("result[%d].TasaDesistimiento = %v, want %v", i, got[i].TasaDesistimiento, tt.wantRates[i])
				}
			}
		})
	}

	catalogo := map[string]domain.Proyecto{
		"A": project("A", 100), "B": project("B", 100),
		"C": project("C", 100), "D": project("D", 100),
	}
	first := RecomendarProyectos([]Vecino{{ProyectoID: "D"}, {ProyectoID: "B"}, {ProyectoID: "A"}, {ProyectoID: "C"}}, catalogo, 100)
	second := RecomendarProyectos([]Vecino{{ProyectoID: "C"}, {ProyectoID: "A"}, {ProyectoID: "D"}, {ProyectoID: "B"}}, catalogo, 100)
	if !reflect.DeepEqual(first, second) {
		t.Fatalf("shuffled inputs changed recommendations:\nfirst=%#v\nsecond=%#v", first, second)
	}
}

func TestRecomendarProyectosEligibilityBoundaries(t *testing.T) {
	cases := []struct {
		name       string
		price      int64
		budget     int64
		wantResult bool
	}{
		{name: "exact boundary for even price", price: 100, budget: 80, wantResult: true},
		{name: "one peso below even boundary", price: 100, budget: 79, wantResult: false},
		{name: "ceil boundary for odd price", price: 101, budget: 81, wantResult: true},
		{name: "one peso below odd boundary", price: 101, budget: 80, wantResult: false},
		{name: "price one requires one peso", price: 1, budget: 0, wantResult: false},
		{name: "zero price is invalid", price: 0, budget: 100, wantResult: false},
		{name: "negative price is invalid", price: -1, budget: 100, wantResult: false},
		{name: "negative budget is invalid", price: 100, budget: -1, wantResult: false},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			catalogo := map[string]domain.Proyecto{
				"P": {ProyectoID: "P", PrecioDesde: tt.price},
			}
			got := RecomendarProyectos([]Vecino{{ProyectoID: "P"}}, catalogo, tt.budget)
			if got != nil && !tt.wantResult && len(got) != 0 {
				t.Fatalf("got invalid recommendation: %#v", got)
			}
			if (len(got) == 1) != tt.wantResult {
				t.Fatalf("eligible = %t, want %t; result=%#v", len(got) == 1, tt.wantResult, got)
			}
		})
	}
}

func TestRecomendarProyectosCopiesCatalogDataAndDoesNotMutateInputs(t *testing.T) {
	vecinos := []Vecino{{ProyectoID: "P", Desistio: true}}
	catalogo := map[string]domain.Proyecto{
		"P": {
			ProyectoID: "P", Nombre: "Proyecto Original", Zona: "Bogotá",
			PrecioDesde: 100, PrecioHasta: 120, EsVIS: true,
			BrochureURL: "brochure", Recorrido360URL: "tour",
		},
	}
	vecinosBefore := append([]Vecino(nil), vecinos...)
	catalogoBefore := map[string]domain.Proyecto{"P": catalogo["P"]}

	got := RecomendarProyectos(vecinos, catalogo, 100)
	if len(got) != 1 {
		t.Fatalf("len(result) = %d, want 1", len(got))
	}
	card := got[0]
	if card.ProyectoID != "P" || card.Nombre != "Proyecto Original" || card.Zona != "Bogotá" ||
		card.PrecioDesde != 100 || card.BrochureURL != "brochure" || card.Recorrido360URL != "tour" {
		t.Fatalf("catalog data was not copied: %#v", card)
	}
	if card.Razon == "" {
		t.Fatal("recommendation rationale is empty")
	}
	if !reflect.DeepEqual(vecinos, vecinosBefore) || !reflect.DeepEqual(catalogo, catalogoBefore) {
		t.Fatalf("inputs mutated: vecinos=%#v catalogo=%#v", vecinos, catalogo)
	}
}

func TestFractionLessUint64AvoidsOverflow(t *testing.T) {
	const max = ^uint64(0)
	cases := []struct {
		name                       string
		leftNumerator, leftDenom   uint64
		rightNumerator, rightDenom uint64
		want                       bool
	}{
		{name: "large products preserve order", leftNumerator: max - 1, leftDenom: max, rightNumerator: max, rightDenom: max, want: true},
		{name: "equal fractions compare equal", leftNumerator: max / 3, leftDenom: max, rightNumerator: max / 3, rightDenom: max, want: false},
		{name: "cross products use both high words", leftNumerator: max - 100, leftDenom: max - 1, rightNumerator: max - 99, rightDenom: max, want: true},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := fractionLessUint64(tt.leftNumerator, tt.leftDenom, tt.rightNumerator, tt.rightDenom); got != tt.want {
				t.Fatalf("fractionLessUint64() = %t, want %t", got, tt.want)
			}
		})
	}
}
