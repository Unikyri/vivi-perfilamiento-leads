package pipeline

import "testing"

func TestSlug(t *testing.T) {
	casos := map[string]string{
		"Agrupación De Vivienda Monguí": "agrupacion_de_vivienda_mongui",
		"INARI":                          "inari",
		"Los Nogales":                    "los_nogales",
	}
	for in, esperado := range casos {
		if got := Slug(in); got != esperado {
			t.Errorf("Slug(%q) = %q, esperado %q", in, got, esperado)
		}
	}
}

func TestValorVivienda(t *testing.T) {
	if got := valorVivienda("1770600000000"); got != 177060000 {
		t.Errorf("valorVivienda = %d, esperado 177060000", got)
	}
}

func TestMapeoCategorias(t *testing.T) {
	casos := map[string]string{"OMEGA": "A", "ETA": "B", "TAU": "C", "CHI": "SIN_DATO", "": "SIN_DATO"}
	for in, esperado := range casos {
		if got := mapear(categoriaMap, in); got != esperado {
			t.Errorf("categoria(%q) = %q, esperado %q", in, got, esperado)
		}
	}
}

func TestNormalizacionEdad(t *testing.T) {
	// las dos grafías deben colapsar al mismo bracket
	if mapear(edadMap, "20 - 35 años") != mapear(edadMap, "20 a 35 años") {
		t.Error("las dos grafías de 20-35 deben normalizar igual")
	}
	if got := mapear(edadMap, "Mayor de 55 años"); got != "55+" {
		t.Errorf("Mayor de 55 = %q, esperado 55+", got)
	}
}
