package motor

import (
	"reflect"
	"testing"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
)

func TestGemeloKNNProjection(t *testing.T) {
	t.Run("resolves catalog zones by exact project ID and omits absent keys", func(t *testing.T) {
		zones := map[string]string{"P-001": " Bogotá   -   Bosa "}
		withZone := projectBuyer(domain.Comprador{
			ProyectoID: "P-001",
			Proyecto:   "Bogotá - Bosa",
			ValorCOP:   999,
		}, zones)
		withoutZone := projectBuyer(domain.Comprador{
			ProyectoID: "P-002",
			Proyecto:   "Bogotá - Bosa",
			ValorCOP:   1,
		}, zones)

		if !withZone.tieneZona || withZone.zona != "Bogotá - Bosa" {
			t.Fatalf("matching catalog zone = %#v, want normalized zone", withZone)
		}
		if withoutZone.tieneZona {
			t.Fatalf("missing catalog key unexpectedly supplied zone: %#v", withoutZone)
		}
	})

	t.Run("does not infer zone from project name or price", func(t *testing.T) {
		buyer := domain.Comprador{
			ProyectoID: "P-001",
			Proyecto:   "Bogotá - Bosa",
			ValorCOP:   250_000_000,
		}
		projected := projectBuyer(buyer, nil)
		if projected.tieneZona {
			t.Fatalf("zone inferred without catalog entry: %#v", projected)
		}

		changed := buyer
		changed.Proyecto = "Cartagena - Manga"
		changed.ValorCOP = 900_000_000
		if got := projectBuyer(changed, map[string]string{"P-001": "Bogotá - Bosa"}); got != projectBuyer(buyer, map[string]string{"P-001": "Bogotá - Bosa"}) {
			t.Fatalf("project name/price changed projected features: got %#v", got)
		}
	})

	t.Run("maps lead income at exact SMMLV boundaries", func(t *testing.T) {
		const smmlv = int64(1_750_905)
		cases := []struct {
			name   string
			income int64
			want   string
		}{
			{name: "two SMMLV is A", income: 2 * smmlv, want: "A"},
			{name: "just above two is B", income: 2*smmlv + 1, want: "B"},
			{name: "four SMMLV is B", income: 4 * smmlv, want: "B"},
			{name: "just above four is C", income: 4*smmlv + 1, want: "C"},
		}

		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				projected := projectLead(EntradaGemelo{Perfil: domain.Perfil{
					"ingreso_hogar": {Valor: tt.income},
				}})
				if !projected.tieneCategoria || projected.categoria != tt.want {
					t.Fatalf("income %d projected as %#v, want category %q", tt.income, projected, tt.want)
				}
			})
		}
	})

	t.Run("omits missing or negative lead income", func(t *testing.T) {
		for name, profile := range map[string]domain.Perfil{
			"missing":    {},
			"negative":   {"ingreso_hogar": {Valor: int64(-1)}},
			"wrong type": {"ingreso_hogar": {Valor: "3500000"}},
		} {
			t.Run(name, func(t *testing.T) {
				if got := projectLead(EntradaGemelo{Perfil: profile}); got.tieneCategoria {
					t.Fatalf("invalid income projected as category: %#v", got)
				}
			})
		}
	})

	t.Run("maps every historical age bracket including open ended 55 plus", func(t *testing.T) {
		cases := []struct {
			label string
			want  float64
		}{
			{label: "20-35", want: 27.5},
			{label: "36-45", want: 40.5},
			{label: "46-55", want: 50.5},
			{label: "55+", want: 60.0},
			{label: " 55+ ", want: 60.0},
		}

		for _, tt := range cases {
			t.Run(tt.label, func(t *testing.T) {
				projected := projectBuyer(domain.Comprador{RangoEdad: tt.label}, nil)
				if !projected.tieneEdad || projected.edad != tt.want {
					t.Fatalf("age %q projected as %#v, want %.1f", tt.label, projected, tt.want)
				}
			})
		}
	})

	t.Run("omits empty, sentinel, and unknown age brackets", func(t *testing.T) {
		for _, label := range []string{"", " ", "SIN_DATO", "  sin_dato  ", "18-19"} {
			t.Run(label, func(t *testing.T) {
				if got := projectBuyer(domain.Comprador{RangoEdad: label}, nil); got.tieneEdad {
					t.Fatalf("age %q unexpectedly projected: %#v", label, got)
				}
			})
		}
	})

	t.Run("normalizes desired zone and preserves optional-value presence", func(t *testing.T) {
		zero := 0
		negative := -1
		projected := projectLead(EntradaGemelo{
			Perfil: domain.Perfil{
				"zona_deseada": {Valor: " Bogotá   -   Bosa "},
				"edad":         {Valor: int64(36)},
			},
			Afiliado:       false,
			PersonasACargo: &zero,
		})
		if !projected.tieneZona || projected.zona != "Bogotá - Bosa" {
			t.Fatalf("desired zone = %#v, want normalized catalog spelling", projected)
		}
		if !projected.tieneEdad || projected.edad != 36 {
			t.Fatalf("lead age = %#v, want 36", projected)
		}
		if !projected.tieneAfiliado || projected.afiliado {
			t.Fatalf("explicit false affiliation was not preserved: %#v", projected)
		}
		if !projected.tienePersonasACargo || projected.personasACargo != 0 {
			t.Fatalf("zero dependents were treated as missing: %#v", projected)
		}

		missing := projectLead(EntradaGemelo{PersonasACargo: &negative})
		if missing.tieneZona || missing.tieneEdad || missing.tienePersonasACargo {
			t.Fatalf("missing optional dimensions were projected: %#v", missing)
		}
	})

	t.Run("normalizes categorical buyer fields and omits sentinel values", func(t *testing.T) {
		projected := projectBuyer(domain.Comprador{
			Categoria:      "  a ",
			RangoEdad:      " 36-45 ",
			PersonasACargo: 0,
		}, map[string]string{"": "  SIN_DATO  "})
		if projected.categoria != "A" || !projected.tieneCategoria {
			t.Fatalf("category = %#v, want normalized A", projected)
		}
		if projected.edad != 40.5 || !projected.tieneEdad {
			t.Fatalf("age = %#v, want normalized 40.5", projected)
		}
		if !projected.tienePersonasACargo || projected.personasACargo != 0 {
			t.Fatalf("zero buyer dependents were treated as missing: %#v", projected)
		}

		missing := projectBuyer(domain.Comprador{
			Categoria: "SIN_DATO",
			RangoEdad: "SIN_DATO",
		}, map[string]string{"": "SIN_DATO"})
		if missing.tieneCategoria || missing.tieneEdad || missing.tieneZona {
			t.Fatalf("sentinel values projected as present: %#v", missing)
		}
	})

	t.Run("projection is pure and does not mutate supplied records or catalog", func(t *testing.T) {
		dependents := 2
		input := EntradaGemelo{
			Perfil:         domain.Perfil{"zona_deseada": {Valor: "Bogotá - Bosa"}},
			PersonasACargo: &dependents,
			Compradores:    []domain.Comprador{{ProyectoID: "P-001", Categoria: "A"}},
			ZonasCatalogo:  map[string]string{"P-001": "Bogotá - Bosa"},
		}
		before := input
		before.Compradores = append([]domain.Comprador(nil), input.Compradores...)
		before.ZonasCatalogo = map[string]string{"P-001": "Bogotá - Bosa"}

		_ = projectLead(input)
		_ = projectBuyer(input.Compradores[0], input.ZonasCatalogo)
		if !reflect.DeepEqual(input, before) {
			t.Fatalf("projection mutated input: before=%#v after=%#v", before, input)
		}
	})
}
