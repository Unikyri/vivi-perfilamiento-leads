package motor

import (
	"fmt"
	"math"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
	"github.com/Unikyri/vivi-perfilamiento-leads/references"
)

const (
	// penalizacionCampoNoVerificado is the confidence multiplier applied for
	// each critical field that is not VERIFICADO_BASE (doc 13 §2.1).
	penalizacionCampoNoVerificado = 0.85
	// confianzaMinima is the floor a confidence score can never drop below.
	confianzaMinima = 0.50
)

var (
	// constantes holds the normative economic parameters, loaded once from
	// the embedded reference JSON.
	constantes = mustLoad(domain.LoadConstantes(references.ConstantesJSON))
	// tramosSubsidio holds the subsidy brackets, loaded once from the
	// embedded reference JSON.
	tramosSubsidio = mustLoad(domain.LoadSubsidios(references.SubsidiosJSON))
)

// CalcularCapacidad computes a lead's financial capacity: max credit,
// applicable subsidy, max budget, ratio against a candidate price, and a
// confidence score. It is fully deterministic: no I/O, no LLM calls, no
// imputed data — only values already present on the given Perfil.
func CalcularCapacidad(p domain.Perfil, esAfiliado bool, precioCandidato int64) domain.Capacidad {
	return calcular(p, esAfiliado, precioCandidato, constantes, tramosSubsidio)
}

// calcular is the test seam: it receives already-loaded constants so tests
// can drive the math directly, and it is the injection path for future
// callers that may need alternate constants.
func calcular(p domain.Perfil, esAfiliado bool, precioCandidato int64,
	cts domain.Constantes, tramos []domain.TramoSubsidio) domain.Capacidad {
	ingreso, _ := p.Entero("ingreso_hogar")
	recursos, _ := p.Entero("recursos_propios")
	if recursos < 0 {
		// A negative declared saving is invalid input, never negative money.
		recursos = 0
	}

	// Credit requires a positive income; zero or negative income yields zero
	// credit, never a negative amount flowing into a lead-facing Capacidad.
	var credito int64
	if ingreso > 0 {
		factor := factorAnualidad(cts.TasaEADefault, cts.PlazoMesesDefault)
		cuota := cts.CuotaIngresoMax * float64(ingreso)
		credito = int64(math.Round(cuota * factor))
	}

	subsidio, reglaSubsidio := subsidioAplicable(p, esAfiliado, cts, tramos)

	presupuesto := credito + subsidio + recursos

	denominador := precioCandidato
	if denominador <= 0 {
		denominador = cts.MedianaVISCOP
	}
	ratio := float64(presupuesto) / float64(denominador)

	return domain.Capacidad{
		PresupuestoMax:    presupuesto,
		CreditoMax:        credito,
		SubsidioAplicable: subsidio,
		RecursosPropios:   recursos,
		Ratio:             ratio,
		Confianza:         calcularConfianza(p),
		Desglose:          armarDesglose(p, cts, credito, subsidio, recursos, reglaSubsidio),
	}
}

// factorAnualidad returns the annuity factor for an effective annual rate
// and a term in months, at full float64 precision. Never rounded — see the
// package's precision contract.
func factorAnualidad(tasaEA float64, plazoMeses int) float64 {
	i := math.Pow(1+tasaEA, 1.0/12.0) - 1
	n := float64(plazoMeses)
	return (1 - math.Pow(1+i, -n)) / i
}

// subsidioAplicable computes the applicable subsidy in COP and the rule
// string that explains the decision. The income guard runs FIRST: an
// absent or non-positive ingreso_hogar zeroes the subsidy before any
// eligibility condition is evaluated (doc 13 §2.1, CAP-6/B-1).
func subsidioAplicable(p domain.Perfil, esAfiliado bool, cts domain.Constantes,
	tramos []domain.TramoSubsidio) (int64, string) {
	ingreso, ok := p.Entero("ingreso_hogar")
	if !ok || ingreso <= 0 {
		return 0, "ingreso_hogar ausente o no positivo: sin ingreso declarado, no aplica subsidio"
	}

	// S2: an absent condition field means the condition is NOT satisfied.
	// The engine never assumes a favorable answer the lead did not give.
	tieneVivienda, ok := p.Booleano("tiene_vivienda")
	if !ok {
		return 0, "tiene_vivienda sin dato: la condicion no se asume favorable"
	}
	if tieneVivienda {
		return 0, "tiene_vivienda = true: el hogar ya cuenta con vivienda"
	}

	recibioSubsidio, ok := p.Booleano("recibio_subsidio")
	if !ok {
		return 0, "recibio_subsidio sin dato: la condicion no se asume favorable"
	}
	if recibioSubsidio {
		return 0, "recibio_subsidio = true: el hogar ya recibio subsidio antes"
	}

	hogarConAfiliado, _ := p.Booleano("hogar_con_afiliado")
	_, cajaExternaPresente := p["caja_externa"]
	if !esAfiliado && !hogarConAfiliado && !cajaExternaPresente {
		return 0, "sin afiliado en el hogar (esAfiliado, hogar_con_afiliado, caja_externa ausentes)"
	}

	topeSMMLV := 0
	for _, t := range tramos {
		if t.IngresoHastaSMMLV > topeSMMLV {
			topeSMMLV = t.IngresoHastaSMMLV
		}
	}
	topePesos := int64(topeSMMLV) * int64(cts.SMMLV2026)
	if ingreso > topePesos {
		return 0, fmt.Sprintf("ingreso_hogar supera el tope de %d SMMLV (%d)", topeSMMLV, topePesos)
	}

	for _, t := range tramos {
		limite := int64(t.IngresoHastaSMMLV) * int64(cts.SMMLV2026)
		if ingreso <= limite {
			monto := int64(t.SubsidioSMMLV) * int64(cts.SMMLV2026)
			return monto, fmt.Sprintf("ingreso_hogar <= %d SMMLV (%d): subsidio %d SMMLV (%d)",
				t.IngresoHastaSMMLV, limite, t.SubsidioSMMLV, monto)
		}
	}
	return 0, "ingreso_hogar no encaja en ningun tramo de subsidio disponible"
}

// calcularConfianza scores how much of the profile's critical data is
// verified. It starts at 1.0 and multiplies by penalizacionCampoNoVerificado
// for each key in domain.CamposCriticos that is not VERIFICADO_BASE, floored
// at confianzaMinima. Iteration MUST stay key-based and order-independent:
// domain.CamposCriticos is a map[string]bool, and the product of identical
// commutative factors is bit-identical regardless of visit order.
func calcularConfianza(p domain.Perfil) float64 {
	c := 1.0
	for campo := range domain.CamposCriticos {
		if !p.EsVerificado(campo) {
			c *= penalizacionCampoNoVerificado
		}
	}
	return math.Max(c, confianzaMinima)
}

// armarDesglose builds the fixed three-line capacity breakdown. Fuente is
// derived from each line's driving profile field: CREDITO and SUBSIDIO from
// ingreso_hogar, RECURSOS_PROPIOS from recursos_propios. When the driving
// field is absent, indexing Perfil yields the zero value domain.FuenteCampo(""),
// which is emitted as-is — never substituted with an invented provenance.
func armarDesglose(p domain.Perfil, cts domain.Constantes,
	credito, subsidio, recursos int64, reglaSubsidio string) []domain.ItemDesglose {
	fuenteIngreso := p["ingreso_hogar"].Fuente
	fuenteRecursos := p["recursos_propios"].Fuente

	reglaCredito := fmt.Sprintf("cuota_ingreso_max=%.2f x ingreso_hogar, tasa_ea=%.3f, plazo_meses=%d",
		cts.CuotaIngresoMax, cts.TasaEADefault, cts.PlazoMesesDefault)
	reglaRecursos := "recursos_propios declarados por el hogar"

	return []domain.ItemDesglose{
		{Concepto: "CREDITO", Monto: credito, Regla: reglaCredito, Fuente: fuenteIngreso},
		{Concepto: "SUBSIDIO", Monto: subsidio, Regla: reglaSubsidio, Fuente: fuenteIngreso},
		{Concepto: "RECURSOS_PROPIOS", Monto: recursos, Regla: reglaRecursos, Fuente: fuenteRecursos},
	}
}

// mustLoad panics on a load error instead of returning a zeroed value. The
// embedded reference JSON is compiled into the binary, so a parse failure
// is a build defect, not a runtime input error; silently degrading to a
// zeroed Constantes would produce credito_max = 0 for every lead.
func mustLoad[T any](v T, err error) T {
	if err != nil {
		panic(fmt.Sprintf("motor: failed to load embedded reference data: %v", err))
	}
	return v
}
