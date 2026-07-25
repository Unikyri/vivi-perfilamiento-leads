package motor

import (
	"math"
	"reflect"
	"testing"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
)

func TestMatriz2x2(t *testing.T) {
	tests := []struct {
		name  string
		input EntradaMatriz
		want  domain.Ruta
	}{
		{
			name:  "affiliate high capacity and high intention goes to advisor",
			input: EntradaMatriz{Ratio: 0.95, Intencion: domain.NivelAlta, EsAfiliado: true},
			want:  domain.RutaAsesor,
		},
		{
			name:  "affiliate high capacity and low intention goes to remarketing",
			input: EntradaMatriz{Ratio: 0.95, Intencion: domain.NivelBaja, EsAfiliado: true},
			want:  domain.RutaRemarketing,
		},
		{
			name:  "affiliate low capacity and high intention goes to nutrition",
			input: EntradaMatriz{Ratio: 0.949999, Intencion: domain.NivelAlta, EsAfiliado: true},
			want:  domain.RutaNutricion,
		},
		{
			name:  "affiliate low capacity and low intention goes to despedida",
			input: EntradaMatriz{Ratio: 0.949999, Intencion: domain.NivelBaja, EsAfiliado: true},
			want:  domain.RutaDespedida,
		},
		{
			name:  "affiliate media intention is low",
			input: EntradaMatriz{Ratio: 2, Intencion: domain.NivelMedia, EsAfiliado: true},
			want:  domain.RutaRemarketing,
		},
		{
			name:  "affiliate unknown intention is low",
			input: EntradaMatriz{Ratio: 2, Intencion: domain.Nivel("UNKNOWN"), EsAfiliado: true},
			want:  domain.RutaRemarketing,
		},
		{
			name:  "affiliate empty intention is low",
			input: EntradaMatriz{Ratio: 2, EsAfiliado: true},
			want:  domain.RutaRemarketing,
		},
		{
			name:  "affiliate one point zero seven and high intention goes to advisor",
			input: EntradaMatriz{Ratio: 1.07, Intencion: domain.NivelAlta, EsAfiliado: true},
			want:  domain.RutaAsesor,
		},
		{
			name:  "affiliate one point zero zero and high intention goes to advisor",
			input: EntradaMatriz{Ratio: 1.00, Intencion: domain.NivelAlta, EsAfiliado: true},
			want:  domain.RutaAsesor,
		},
		{
			name:  "non affiliate one point zero zero and high intention goes to nutrition",
			input: EntradaMatriz{Ratio: 1.00, Intencion: domain.NivelAlta},
			want:  domain.RutaNutricion,
		},
		{
			name:  "non affiliate exact point nine five and high intention goes to nutrition",
			input: EntradaMatriz{Ratio: 0.95, Intencion: domain.NivelAlta},
			want:  domain.RutaNutricion,
		},
		{
			name:  "non affiliate exact advisor threshold",
			input: EntradaMatriz{Ratio: 1.05, Intencion: domain.NivelAlta},
			want:  domain.RutaAsesor,
		},
		{
			name:  "non affiliate below advisor threshold goes to nutrition",
			input: EntradaMatriz{Ratio: math.Nextafter(1.05, 0), Intencion: domain.NivelAlta},
			want:  domain.RutaNutricion,
		},
		{
			name:  "non affiliate low intention goes to despedida at high ratio",
			input: EntradaMatriz{Ratio: 2, Intencion: domain.NivelBaja},
			want:  domain.RutaDespedida,
		},
		{
			name:  "conversion preference overrides advisor eligibility",
			input: EntradaMatriz{Ratio: 1.10, Intencion: domain.NivelAlta, RutaConversion: true},
			want:  domain.RutaNutricion,
		},
		{
			name:  "conversion preference also applies at invalid ratio",
			input: EntradaMatriz{Ratio: math.NaN(), Intencion: domain.NivelAlta, RutaConversion: true},
			want:  domain.RutaNutricion,
		},
		{
			name:  "affiliate invalid ratio and high intention never reaches advisor",
			input: EntradaMatriz{Ratio: math.NaN(), Intencion: domain.NivelAlta, EsAfiliado: true},
			want:  domain.RutaNutricion,
		},
		{
			name:  "affiliate positive infinity is low capacity",
			input: EntradaMatriz{Ratio: math.Inf(1), Intencion: domain.NivelAlta, EsAfiliado: true},
			want:  domain.RutaNutricion,
		},
		{
			name:  "affiliate negative infinity is low capacity",
			input: EntradaMatriz{Ratio: math.Inf(-1), Intencion: domain.NivelAlta, EsAfiliado: true},
			want:  domain.RutaNutricion,
		},
		{
			name:  "negative ratio is never advisor",
			input: EntradaMatriz{Ratio: -1, Intencion: domain.NivelAlta},
			want:  domain.RutaNutricion,
		},
		{
			name:  "non affiliate invalid ratio and low intention is despedida",
			input: EntradaMatriz{Ratio: math.NaN(), Intencion: domain.NivelMedia},
			want:  domain.RutaDespedida,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Matriz2x2(tt.input); got != tt.want {
				t.Fatalf("Matriz2x2(%+v) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestMatriz2x2IsPureAndRouteOnly(t *testing.T) {
	input := EntradaMatriz{
		Ratio:          1.10,
		Intencion:      domain.NivelAlta,
		EsAfiliado:     false,
		RutaConversion: false,
	}
	before := input

	first := Matriz2x2(input)
	second := Matriz2x2(input)
	if first != domain.RutaAsesor || second != first {
		t.Fatalf("repeated route selection = %q, %q; want stable %q", first, second, domain.RutaAsesor)
	}
	if !reflect.DeepEqual(input, before) {
		t.Fatalf("matrix mutated its input: before=%+v after=%+v", before, input)
	}
}
