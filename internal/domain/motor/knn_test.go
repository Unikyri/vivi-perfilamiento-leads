package motor

import (
	"math"
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

func TestGemeloKNNImputation(t *testing.T) {
	buyers := []domain.Comprador{
		{ProyectoID: "P-1", Afiliado: true, Categoria: "A", RangoEdad: "20-35"},
		{ProyectoID: "P-1", Afiliado: true, Categoria: "B", RangoEdad: "36-45"},
		{ProyectoID: "P-1", Afiliado: true, Categoria: "A", RangoEdad: "46-55"},
		{ProyectoID: "P-1", Afiliado: true, Categoria: "B", RangoEdad: "55+"},
		{ProyectoID: "P-1", Afiliado: true, Categoria: "C", RangoEdad: "20-35"},
		{ProyectoID: "P-2", Afiliado: true, Categoria: "C", RangoEdad: "20-35"},
		{ProyectoID: "P-2", Afiliado: true, Categoria: "C", RangoEdad: "20-35"},
		{ProyectoID: "P-2", Afiliado: true, Categoria: "A", RangoEdad: "36-45"},
		{ProyectoID: "P-2", Afiliado: true, Categoria: "B", RangoEdad: "46-55"},
		{ProyectoID: "P-2", Afiliado: true, Categoria: "B", RangoEdad: "55+"},
		{ProyectoID: "P-2", Afiliado: true, Categoria: "A", RangoEdad: "55+"},
	}
	for i := 0; i < 4; i++ {
		buyers = append(buyers, domain.Comprador{
			ProyectoID: "P-3", Afiliado: true, Categoria: "C", RangoEdad: "55+",
		})
	}
	for i := 0; i < 5; i++ {
		buyers = append(buyers, domain.Comprador{
			ProyectoID: "P-4", Afiliado: true, Categoria: "B", RangoEdad: "SIN_DATO",
		})
		buyers = append(buyers, domain.Comprador{
			ProyectoID: "P-5", Afiliado: true, Categoria: "SIN_DATO", RangoEdad: "46-55",
		})
	}

	stats := buildProjectAffiliateStats(buyers)
	if got := stats["P-1"]; !got.hasCategory || got.categoryMode != "A" || !got.hasAge || got.ageMedian != 40.5 {
		t.Fatalf("P-1 statistics = %#v, want deterministic mode A and odd median 40.5", got)
	}
	if got := stats["P-2"]; !got.hasAge || got.ageMedian != 45.5 {
		t.Fatalf("P-2 statistics = %#v, want even median 45.5", got)
	}
	if got := stats["P-3"]; got.affiliateCount != 4 || got.hasCategory || got.hasAge {
		t.Fatalf("below-threshold statistics = %#v, want omission", got)
	}

	imputed := projectBuyerWithStats(domain.Comprador{ProyectoID: "P-1"}, nil, stats)
	if !imputed.tieneCategoria || imputed.categoria != "A" || !imputed.tieneEdad || imputed.edad != 40.5 {
		t.Fatalf("eligible project fallback = %#v, want category A and age 40.5", imputed)
	}
	belowThreshold := projectBuyerWithStats(domain.Comprador{ProyectoID: "P-3"}, nil, stats)
	if belowThreshold.tieneCategoria || belowThreshold.tieneEdad {
		t.Fatalf("below-threshold fallback = %#v, want both dimensions omitted", belowThreshold)
	}
	categoryOnly := projectBuyerWithStats(domain.Comprador{ProyectoID: "P-4"}, nil, stats)
	if !categoryOnly.tieneCategoria || categoryOnly.categoria != "B" || categoryOnly.tieneEdad {
		t.Fatalf("independent category fallback = %#v, want age omitted", categoryOnly)
	}
	ageOnly := projectBuyerWithStats(domain.Comprador{ProyectoID: "P-5"}, nil, stats)
	if ageOnly.tieneCategoria || !ageOnly.tieneEdad || ageOnly.edad != 50.5 {
		t.Fatalf("independent age fallback = %#v, want category omitted", ageOnly)
	}
	crossProject := projectBuyerWithStats(domain.Comprador{ProyectoID: "P-6"}, nil, stats)
	if crossProject.tieneCategoria || crossProject.tieneEdad {
		t.Fatalf("cross-project fallback = %#v, must not use P-1/P-2 statistics", crossProject)
	}
}

func TestGemeloKNNGower(t *testing.T) {
	allDifferent := featureVector{
		categoria: "A", tieneCategoria: true,
		zona: "N", tieneZona: true,
		edad: 27.5, tieneEdad: true,
		afiliado: false, tieneAfiliado: true,
		personasACargo: 0, tienePersonasACargo: true,
	}
	other := featureVector{
		categoria: "B", tieneCategoria: true,
		zona: "S", tieneZona: true,
		edad: 60, tieneEdad: true,
		afiliado: true, tieneAfiliado: true,
		personasACargo: 10, tienePersonasACargo: true,
	}
	if got := gowerDistance(allDifferent, other); got != 1 {
		t.Fatalf("maximum Gower distance = %.12f, want 1", got)
	}
	if got := gowerDistance(allDifferent, allDifferent); got != 0 {
		t.Fatalf("identity Gower distance = %.12f, want 0", got)
	}
	if got, reverse := gowerDistance(allDifferent, other), gowerDistance(other, allDifferent); got != reverse {
		t.Fatalf("Gower is not symmetric: forward=%.12f reverse=%.12f", got, reverse)
	}

	missingCategory := allDifferent
	missingCategory.tieneCategoria = false
	missingCategory.zona = other.zona
	wantRenormalized := (gowerAgeWeight + gowerAffiliationWeight + gowerDependentsWeight) /
		(gowerZoneWeight + gowerAgeWeight + gowerAffiliationWeight + gowerDependentsWeight)
	if got := gowerDistance(missingCategory, other); math.Abs(got-wantRenormalized) > 1e-12 {
		t.Fatalf("renormalized Gower = %.12f, want %.12f", got, wantRenormalized)
	}

	zeroDependents := other
	zeroDependents.personasACargo = 0
	if got := gowerDistance(zeroDependents, zeroDependents); got != 0 {
		t.Fatalf("equal zero dependents were not treated as present: %.12f", got)
	}
	changedDependents := zeroDependents
	changedDependents.personasACargo = 10
	if got := gowerDistance(zeroDependents, changedDependents); got != gowerDependentsWeight {
		t.Fatalf("dependent distance with 0 versus 10 = %.12f, want %.12f", got, gowerDependentsWeight)
	}
}

func TestGemeloKNNSelection(t *testing.T) {
	t.Run("returns 30 deterministic distance-then-ID ordered neighbors", func(t *testing.T) {
		buyers := make([]domain.Comprador, 0, 35)
		// Twenty buyers have the closest age bracket; their input order is
		// intentionally reversed so ID ordering is observable for ties.
		for id := 20; id >= 1; id-- {
			buyers = append(buyers, domain.Comprador{
				ID: id, ProyectoID: "P-1", RangoEdad: "36-45",
			})
		}
		for id := 35; id >= 21; id-- {
			buyers = append(buyers, domain.Comprador{
				ID: id, ProyectoID: "P-1", RangoEdad: "20-35",
			})
		}

		got := GemeloKNN(EntradaGemelo{
			Perfil:      domain.Perfil{"edad": {Valor: int64(40)}},
			Compradores: buyers,
			K:           30,
		})
		if len(got) != 30 {
			t.Fatalf("neighbor count = %d, want 30", len(got))
		}
		for i, neighbor := range got {
			if neighbor.ID != i+1 {
				t.Fatalf("neighbor[%d].ID = %d, want %d; got %#v", i, neighbor.ID, i+1, got)
			}
			if i > 0 && got[i-1].Distancia > neighbor.Distancia {
				t.Fatalf("distance order regressed at index %d: %#v", i, got)
			}
		}
		for i := 1; i < 20; i++ {
			if got[i-1].Distancia != got[i].Distancia {
				t.Fatalf("close-distance tie changed at index %d: %#v", i, got)
			}
		}
		if got[19].Distancia >= got[20].Distancia {
			t.Fatalf("close buyers were not strictly nearer than distant buyers: %#v", got)
		}
	})

	t.Run("returns all buyers when K exceeds N and orders equal distances by ID", func(t *testing.T) {
		got := GemeloKNN(EntradaGemelo{
			Compradores: []domain.Comprador{
				{ID: 9, ProyectoID: "P-9"},
				{ID: 2, ProyectoID: "P-2"},
				{ID: 7, ProyectoID: "P-7"},
			},
			K: 99,
		})
		wantIDs := []int{2, 7, 9}
		if len(got) != len(wantIDs) {
			t.Fatalf("neighbor count = %d, want %d", len(got), len(wantIDs))
		}
		for i, wantID := range wantIDs {
			if got[i].ID != wantID {
				t.Fatalf("neighbor[%d].ID = %d, want %d", i, got[i].ID, wantID)
			}
		}
	})

	t.Run("returns non-nil empty output for non-positive K and empty input", func(t *testing.T) {
		cases := map[string]EntradaGemelo{
			"zero K":       {K: 0, Compradores: []domain.Comprador{{ID: 1}}},
			"negative K":   {K: -1, Compradores: []domain.Comprador{{ID: 1}}},
			"empty buyers": {K: 1},
		}
		for name, in := range cases {
			t.Run(name, func(t *testing.T) {
				got := GemeloKNN(in)
				if got == nil {
					t.Fatal("empty result is nil, want non-nil empty slice")
				}
				if len(got) != 0 {
					t.Fatalf("empty result length = %d, want 0", len(got))
				}
			})
		}
	})

	t.Run("is repeatable and does not mutate buyers or catalog", func(t *testing.T) {
		buyers := []domain.Comprador{
			{ID: 3, ProyectoID: "P-1", RangoEdad: "36-45", Proyecto: "Original", ValorCOP: 100},
			{ID: 1, ProyectoID: "P-2", RangoEdad: "20-35", Proyecto: "Other", ValorCOP: 200},
		}
		zones := map[string]string{"P-1": "Bogotá - Bosa", "P-2": "Cali - Sur"}
		in := EntradaGemelo{
			Perfil:        domain.Perfil{"edad": {Valor: int64(40)}},
			Compradores:   buyers,
			ZonasCatalogo: zones,
			K:             2,
		}
		beforeBuyers := append([]domain.Comprador(nil), in.Compradores...)
		beforeZones := map[string]string{"P-1": zones["P-1"], "P-2": zones["P-2"]}

		first := GemeloKNN(in)
		second := GemeloKNN(in)
		if !reflect.DeepEqual(first, second) {
			t.Fatalf("repeated calls differ: first=%#v second=%#v", first, second)
		}
		if !reflect.DeepEqual(in.Compradores, beforeBuyers) || !reflect.DeepEqual(in.ZonasCatalogo, beforeZones) {
			t.Fatalf("selection mutated inputs: before buyers=%#v zones=%#v after=%#v zones=%#v", beforeBuyers, beforeZones, in.Compradores, in.ZonasCatalogo)
		}
	})

	t.Run("does not use project name or price when catalog zones are fixed", func(t *testing.T) {
		base := domain.Comprador{
			ID: 1, ProyectoID: "P-1", Proyecto: "Bogotá - Bosa", ValorCOP: 250_000_000,
			Categoria: "A", RangoEdad: "36-45", PersonasACargo: 1,
		}
		changed := base
		changed.Proyecto = "Cartagena - Manga"
		changed.ValorCOP = 900_000_000
		input := EntradaGemelo{
			Perfil:      domain.Perfil{"edad": {Valor: int64(40)}},
			Compradores: []domain.Comprador{base},
			ZonasCatalogo: map[string]string{
				"P-1": "Bogotá - Bosa",
			},
			K: 1,
		}
		changedInput := input
		changedInput.Compradores = []domain.Comprador{changed}
		if got, changedGot := GemeloKNN(input), GemeloKNN(changedInput); !reflect.DeepEqual(got, changedGot) {
			t.Fatalf("name/price changed fixed-catalog result: base=%#v changed=%#v", got, changedGot)
		}
	})
}

func TestZonasCoinciden(t *testing.T) {
	casos := []struct {
		nombre   string
		lead     string
		catalogo string
		want     bool
	}{
		{"municipio suelto contra nombre comercial", "Soacha", "Ciudadela Maiporé - Soacha", true},
		{"acentos plegados", "maipore", "Ciudadela Maiporé - Soacha", true},
		{"mayusculas y espacios sobrantes", "  SOACHA  ", "Ciudadela Maiporé - Soacha", true},
		{"bogota contra ciudadela calle 80", "Bogotá", "Ciudadela Calle 80 - Bogotá", true},
		{"numero de via", "Calle 80", "Ciudadela Calle 80 - Bogotá", true},
		{"zona de un solo token", "Ricaurte", "Ricaurte", true},
		{"municipios distintos NO coinciden", "Soacha", "Ciudadela Calle 80 - Bogotá", false},
		{"palabra vacia no basta", "Ciudadela", "Ciudadela Maiporé - Soacha", false},
		{"conjunto no basta", "conjunto residencial", "Ciudadela Maiporé - Soacha", false},
		{"lead vacio", "", "Ciudadela Maiporé - Soacha", false},
		{"catalogo vacio", "Soacha", "", false},
		{"SIN_DATO se trata como vacio", "SIN_DATO", "Ciudadela Maiporé - Soacha", false},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			if got := zonasCoinciden(c.lead, c.catalogo); got != c.want {
				t.Errorf("zonasCoinciden(%q, %q) = %v, want %v", c.lead, c.catalogo, got, c.want)
			}
		})
	}
}

func TestZonasCoincidenEsSimetrica(t *testing.T) {
	pares := [][2]string{
		{"Soacha", "Ciudadela Maiporé - Soacha"},
		{"Bogotá", "Ciudadela Calle 80 - Bogotá"},
		{"Soacha", "Ricaurte"},
	}
	for _, p := range pares {
		if zonasCoinciden(p[0], p[1]) != zonasCoinciden(p[1], p[0]) {
			t.Errorf("asimetría entre %q y %q", p[0], p[1])
		}
	}
}

// TestZonaParcialNoVaciaRecomendaciones fija el defecto que motivó este
// cambio: una zona escrita por el lead no puede dejar el catálogo en cero.
func TestZonaParcialNoVaciaRecomendaciones(t *testing.T) {
	catalogo := map[string]domain.Proyecto{
		"mongui":   {ProyectoID: "mongui", Nombre: "Monguí", Zona: "Ciudadela Maiporé - Soacha", PrecioDesde: 156_470_000},
		"macarena": {ProyectoID: "macarena", Nombre: "La Macarena", Zona: "Ciudadela Maiporé - Soacha", PrecioDesde: 128_340_000},
	}
	zonas := map[string]string{}
	for id, p := range catalogo {
		zonas[id] = p.Zona
	}

	compradores := []domain.Comprador{
		{ID: 1, ProyectoID: "mongui", Afiliado: true, Categoria: "A", RangoEdad: "20-35", PersonasACargo: 2},
		{ID: 2, ProyectoID: "macarena", Afiliado: true, Categoria: "A", RangoEdad: "20-35", PersonasACargo: 2},
		// Comprador fuera del catálogo: sin zona resoluble.
		{ID: 3, ProyectoID: "la_arboleda", Afiliado: true, Categoria: "A", RangoEdad: "20-35", PersonasACargo: 2},
	}

	personas := 2
	perfil := domain.Perfil{
		"ingreso_hogar": domain.CampoPerfil{Valor: int64(2_600_000), Fuente: domain.FuenteCampoVerificadoBase},
		"edad":          domain.CampoPerfil{Valor: int64(32), Fuente: domain.FuenteCampoVerificadoBase},
		"zona_deseada":  domain.CampoPerfil{Valor: "Soacha", Fuente: domain.FuenteCampoDeclarado},
	}

	vecinos := GemeloKNN(EntradaGemelo{
		Perfil: perfil, Afiliado: true, PersonasACargo: &personas,
		Compradores: compradores, ZonasCatalogo: zonas, K: 3,
	})

	enCatalogo := 0
	for _, v := range vecinos {
		if _, ok := catalogo[v.ProyectoID]; ok {
			enCatalogo++
		}
	}
	if enCatalogo == 0 {
		t.Fatal("regresión: una zona parcial dejó cero vecinos en catálogo")
	}

	recs := RecomendarProyectos(vecinos, catalogo, 166_771_122)
	if len(recs) == 0 {
		t.Fatal("regresión: cero recomendaciones con zona_deseada=\"Soacha\"")
	}

	// El comprador del catálogo cuya zona coincide debe quedar MÁS CERCA
	// que el de fuera del catálogo, no más lejos.
	if vecinos[0].ProyectoID == "la_arboleda" {
		t.Errorf("el comprador fuera de catálogo quedó primero: la zona sigue penalizando al catálogo")
	}
}
