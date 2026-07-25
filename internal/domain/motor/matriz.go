package motor

import (
	"math"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
)

// EntradaMatriz contains the caller-supplied facts used by the pure routing
// matrix. Ratio is the authoritative capacity ratio for this decision.
type EntradaMatriz struct {
	Ratio          float64
	Intencion      domain.Nivel
	EsAfiliado     bool
	RutaConversion bool
}

// Matriz2x2 selects a route without recalculating capacity or mutating any
// lead, ficha, cupo, plan, or milestone state.
func Matriz2x2(entrada EntradaMatriz) domain.Ruta {
	altaIntencion := entrada.Intencion == domain.NivelAlta

	// Conversion is preferred before any capacity decision for non-affiliates
	// with high intention.
	if !entrada.EsAfiliado && altaIntencion && entrada.RutaConversion {
		return domain.RutaNutricion
	}

	// Invalid ratios are never high capacity. This also prevents NaN and
	// infinities from reaching an advisor route through threshold comparisons.
	ratioValido := !math.IsNaN(entrada.Ratio) && !math.IsInf(entrada.Ratio, 0) && entrada.Ratio >= 0
	altaCapacidad := ratioValido && entrada.Ratio >= 0.95

	if entrada.EsAfiliado {
		if altaCapacidad {
			if altaIntencion {
				return domain.RutaAsesor
			}
			return domain.RutaRemarketing
		}
		if altaIntencion {
			return domain.RutaNutricion
		}
		return domain.RutaDespedida
	}

	if altaIntencion && ratioValido && entrada.Ratio >= 1.05 {
		return domain.RutaAsesor
	}
	if altaIntencion {
		return domain.RutaNutricion
	}
	return domain.RutaDespedida
}
