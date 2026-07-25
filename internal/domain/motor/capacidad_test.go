package motor

import (
	"math"
	"testing"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
)

// casoCapacidad is one doc 13 §7 fixture. wantRatio/wantConfianza are
// pointers because doc 13 constrains those two floats only for a handful
// of cases; 0.0 is a value a buggy implementation could legitimately
// produce, so it cannot double as "skip".
type casoCapacidad struct {
	nombre          string
	perfil          domain.Perfil
	esAfiliado      bool
	precioCandidato int64

	wantCredito     int64
	wantSubsidio    int64
	wantPresupuesto int64
	wantRatio       *float64
	wantConfianza   *float64

	wantFuenteIngreso  domain.FuenteCampo
	wantFuenteRecursos domain.FuenteCampo
}

func newCampo(valor any, fuente domain.FuenteCampo) domain.CampoPerfil {
	return domain.CampoPerfil{Valor: valor, Fuente: fuente}
}

func f64(v float64) *float64 { return &v }

// assertClose compares floats with the epsilon this package's precision
// contract demands: full float64 through every intermediate, tolerance
// only at the assertion boundary.
func assertClose(t *testing.T, nombre string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("%s = %v, want %v (epsilon 1e-9)", nombre, got, want)
	}
}

func assertCaso(t *testing.T, c casoCapacidad) {
	t.Helper()
	got := CalcularCapacidad(c.perfil, c.esAfiliado, c.precioCandidato)

	if got.CreditoMax != c.wantCredito {
		t.Errorf("CreditoMax = %d, want %d", got.CreditoMax, c.wantCredito)
	}
	if got.SubsidioAplicable != c.wantSubsidio {
		t.Errorf("SubsidioAplicable = %d, want %d", got.SubsidioAplicable, c.wantSubsidio)
	}
	if got.PresupuestoMax != c.wantPresupuesto {
		t.Errorf("PresupuestoMax = %d, want %d", got.PresupuestoMax, c.wantPresupuesto)
	}
	if c.wantRatio != nil {
		assertClose(t, "Ratio", got.Ratio, *c.wantRatio)
	}
	if c.wantConfianza != nil {
		assertClose(t, "Confianza", got.Confianza, *c.wantConfianza)
	}

	if len(got.Desglose) != 3 {
		t.Fatalf("len(Desglose) = %d, want 3", len(got.Desglose))
	}
	wantConceptos := [3]string{"CREDITO", "SUBSIDIO", "RECURSOS_PROPIOS"}
	wantFuentes := [3]domain.FuenteCampo{c.wantFuenteIngreso, c.wantFuenteIngreso, c.wantFuenteRecursos}
	for i, item := range got.Desglose {
		if item.Concepto != wantConceptos[i] {
			t.Errorf("Desglose[%d].Concepto = %q, want %q", i, item.Concepto, wantConceptos[i])
		}
		if item.Fuente != wantFuentes[i] {
			t.Errorf("Desglose[%d].Fuente = %q, want %q", i, item.Fuente, wantFuentes[i])
		}
	}
}

func TestCalcularCapacidad(t *testing.T) {
	casos := []casoCapacidad{
		{
			nombre: "CAP-1_Ana",
			perfil: domain.Perfil{
				"ingreso_hogar":    newCampo(int64(2_600_000), domain.FuenteCampoVerificadoBase),
				"recursos_propios": newCampo(int64(8_000_000), domain.FuenteCampoDeclarado),
				"tiene_vivienda":   newCampo(false, domain.FuenteCampoVerificadoBase),
				"recibio_subsidio": newCampo(false, domain.FuenteCampoVerificadoBase),
			},
			esAfiliado:         true,
			precioCandidato:    156_470_000,
			wantCredito:        106_243_972,
			wantSubsidio:       52_527_150,
			wantPresupuesto:    166_771_122,
			wantRatio:          f64(1.0658344858439317),
			wantConfianza:      f64(0.85),
			wantFuenteIngreso:  domain.FuenteCampoVerificadoBase,
			wantFuenteRecursos: domain.FuenteCampoDeclarado,
		},
		{
			nombre: "CAP-2",
			perfil: domain.Perfil{
				"ingreso_hogar":    newCampo(int64(2_000_000), domain.FuenteCampoVerificadoBase),
				"recursos_propios": newCampo(int64(0), domain.FuenteCampoVerificadoBase),
				"tiene_vivienda":   newCampo(false, domain.FuenteCampoVerificadoBase),
				"recibio_subsidio": newCampo(false, domain.FuenteCampoVerificadoBase),
			},
			esAfiliado:         true,
			precioCandidato:    100_000_000,
			wantCredito:        81_726_133,
			wantSubsidio:       52_527_150,
			wantPresupuesto:    134_253_283,
			wantFuenteIngreso:  domain.FuenteCampoVerificadoBase,
			wantFuenteRecursos: domain.FuenteCampoVerificadoBase,
		},
		{
			nombre: "CAP-3",
			perfil: domain.Perfil{
				"ingreso_hogar":    newCampo(int64(3_000_000), domain.FuenteCampoVerificadoBase),
				"recursos_propios": newCampo(int64(0), domain.FuenteCampoVerificadoBase),
				"tiene_vivienda":   newCampo(false, domain.FuenteCampoVerificadoBase),
				"recibio_subsidio": newCampo(false, domain.FuenteCampoVerificadoBase),
			},
			esAfiliado:         true,
			precioCandidato:    100_000_000,
			wantCredito:        122_589_199,
			wantSubsidio:       52_527_150,
			wantPresupuesto:    175_116_349,
			wantFuenteIngreso:  domain.FuenteCampoVerificadoBase,
			wantFuenteRecursos: domain.FuenteCampoVerificadoBase,
		},
		{
			nombre: "CAP-4_tramo_medio",
			perfil: domain.Perfil{
				"ingreso_hogar":    newCampo(int64(5_000_000), domain.FuenteCampoVerificadoBase),
				"recursos_propios": newCampo(int64(0), domain.FuenteCampoVerificadoBase),
				"tiene_vivienda":   newCampo(false, domain.FuenteCampoVerificadoBase),
				"recibio_subsidio": newCampo(false, domain.FuenteCampoVerificadoBase),
			},
			esAfiliado:         true,
			precioCandidato:    100_000_000,
			wantCredito:        204_315_332,
			wantSubsidio:       35_018_100,
			wantPresupuesto:    239_333_432,
			wantFuenteIngreso:  domain.FuenteCampoVerificadoBase,
			wantFuenteRecursos: domain.FuenteCampoVerificadoBase,
		},
		{
			nombre: "CAP-5_mas_4_smmlv",
			perfil: domain.Perfil{
				"ingreso_hogar":    newCampo(int64(8_000_000), domain.FuenteCampoVerificadoBase),
				"recursos_propios": newCampo(int64(0), domain.FuenteCampoVerificadoBase),
				"tiene_vivienda":   newCampo(false, domain.FuenteCampoVerificadoBase),
				"recibio_subsidio": newCampo(false, domain.FuenteCampoVerificadoBase),
			},
			esAfiliado:         true,
			precioCandidato:    100_000_000,
			wantCredito:        326_904_530,
			wantSubsidio:       0,
			wantPresupuesto:    326_904_530,
			wantFuenteIngreso:  domain.FuenteCampoVerificadoBase,
			wantFuenteRecursos: domain.FuenteCampoVerificadoBase,
		},
		{
			// CAP-6 (value row) and B-1 (zero-income boundary property) are
			// one scenario per design D6: the income guard zeroes the
			// subsidy, not affiliation. ingreso_hogar is present with value
			// 0 (DECLARADO), never absent, so every emitted Fuente stays
			// on-enum (F1).
			nombre: "CAP-6_B-1_ingreso_cero",
			perfil: domain.Perfil{
				"ingreso_hogar":    newCampo(int64(0), domain.FuenteCampoDeclarado),
				"recursos_propios": newCampo(int64(3_000_000), domain.FuenteCampoDeclarado),
				"tiene_vivienda":   newCampo(false, domain.FuenteCampoDeclarado),
				"recibio_subsidio": newCampo(false, domain.FuenteCampoDeclarado),
			},
			esAfiliado:         true,
			precioCandidato:    100_000_000,
			wantCredito:        0,
			wantSubsidio:       0,
			wantPresupuesto:    3_000_000,
			wantFuenteIngreso:  domain.FuenteCampoDeclarado,
			wantFuenteRecursos: domain.FuenteCampoDeclarado,
		},
		{
			nombre: "CAP-7_tiene_vivienda",
			perfil: domain.Perfil{
				"ingreso_hogar":    newCampo(int64(2_000_000), domain.FuenteCampoVerificadoBase),
				"recursos_propios": newCampo(int64(0), domain.FuenteCampoVerificadoBase),
				"tiene_vivienda":   newCampo(true, domain.FuenteCampoVerificadoBase),
				"recibio_subsidio": newCampo(false, domain.FuenteCampoVerificadoBase),
			},
			esAfiliado:         true,
			precioCandidato:    100_000_000,
			wantCredito:        81_726_133,
			wantSubsidio:       0,
			wantPresupuesto:    81_726_133,
			wantFuenteIngreso:  domain.FuenteCampoVerificadoBase,
			wantFuenteRecursos: domain.FuenteCampoVerificadoBase,
		},
		{
			// CAP-8 (value row) and B-4 (confidence-floor property) share
			// Ana's CAP-1 inputs with every critical field re-stamped
			// DECLARADO. With exactly four critical fields the minimum
			// reachable product is 0.85^4 = 0.52200625, which sits above
			// the 0.5 floor — the floor is a defensive bound that current
			// parameters can never trigger (doc 13 v1.1 correction).
			nombre: "CAP-8_B-4_confianza_minima",
			perfil: domain.Perfil{
				"ingreso_hogar":    newCampo(int64(2_600_000), domain.FuenteCampoDeclarado),
				"recursos_propios": newCampo(int64(8_000_000), domain.FuenteCampoDeclarado),
				"tiene_vivienda":   newCampo(false, domain.FuenteCampoDeclarado),
				"recibio_subsidio": newCampo(false, domain.FuenteCampoDeclarado),
			},
			esAfiliado:         true,
			precioCandidato:    156_470_000,
			wantCredito:        106_243_972,
			wantSubsidio:       52_527_150,
			wantPresupuesto:    166_771_122,
			wantConfianza:      f64(0.52200625),
			wantFuenteIngreso:  domain.FuenteCampoDeclarado,
			wantFuenteRecursos: domain.FuenteCampoDeclarado,
		},
		{
			nombre: "B-5_ingreso_supera_tope",
			perfil: domain.Perfil{
				"ingreso_hogar":    newCampo(int64(7_500_000), domain.FuenteCampoVerificadoBase),
				"recursos_propios": newCampo(int64(0), domain.FuenteCampoVerificadoBase),
				"tiene_vivienda":   newCampo(false, domain.FuenteCampoVerificadoBase),
				"recibio_subsidio": newCampo(false, domain.FuenteCampoVerificadoBase),
			},
			esAfiliado:         true,
			precioCandidato:    100_000_000,
			wantCredito:        306_472_997,
			wantSubsidio:       0,
			wantPresupuesto:    306_472_997,
			wantFuenteIngreso:  domain.FuenteCampoVerificadoBase,
			wantFuenteRecursos: domain.FuenteCampoVerificadoBase,
		},
		{
			// Reliability review finding #1: a negative declared income must
			// yield zero credit, never a negative amount in a lead-facing
			// Capacidad.
			nombre: "R-1_ingreso_negativo",
			perfil: domain.Perfil{
				"ingreso_hogar":    newCampo(int64(-5_000_000), domain.FuenteCampoDeclarado),
				"recursos_propios": newCampo(int64(3_000_000), domain.FuenteCampoDeclarado),
				"tiene_vivienda":   newCampo(false, domain.FuenteCampoDeclarado),
				"recibio_subsidio": newCampo(false, domain.FuenteCampoDeclarado),
			},
			esAfiliado:         true,
			precioCandidato:    100_000_000,
			wantCredito:        0,
			wantSubsidio:       0,
			wantPresupuesto:    3_000_000,
			wantFuenteIngreso:  domain.FuenteCampoDeclarado,
			wantFuenteRecursos: domain.FuenteCampoDeclarado,
		},
		{
			// Reliability review finding #2: negative declared savings are
			// clamped to zero instead of subtracting from the budget.
			nombre: "R-2_recursos_negativos",
			perfil: domain.Perfil{
				"ingreso_hogar":    newCampo(int64(2_000_000), domain.FuenteCampoVerificadoBase),
				"recursos_propios": newCampo(int64(-1_000_000), domain.FuenteCampoDeclarado),
				"tiene_vivienda":   newCampo(false, domain.FuenteCampoVerificadoBase),
				"recibio_subsidio": newCampo(false, domain.FuenteCampoVerificadoBase),
			},
			esAfiliado:         true,
			precioCandidato:    100_000_000,
			wantCredito:        81_726_133,
			wantSubsidio:       52_527_150,
			wantPresupuesto:    134_253_283,
			wantFuenteIngreso:  domain.FuenteCampoVerificadoBase,
			wantFuenteRecursos: domain.FuenteCampoDeclarado,
		},
		{
			// Reliability review finding #5: numeric fields decoded from JSON
			// arrive as float64, the dominant real ingest shape. CAP-1 must
			// reproduce identically with float64-stored values.
			nombre: "R-5_CAP-1_valores_float64",
			perfil: domain.Perfil{
				"ingreso_hogar":    newCampo(float64(2_600_000), domain.FuenteCampoVerificadoBase),
				"recursos_propios": newCampo(float64(8_000_000), domain.FuenteCampoDeclarado),
				"tiene_vivienda":   newCampo(false, domain.FuenteCampoVerificadoBase),
				"recibio_subsidio": newCampo(false, domain.FuenteCampoVerificadoBase),
			},
			esAfiliado:         true,
			precioCandidato:    156_470_000,
			wantCredito:        106_243_972,
			wantSubsidio:       52_527_150,
			wantPresupuesto:    166_771_122,
			wantRatio:          f64(1.0658344858439317),
			wantConfianza:      f64(0.85),
			wantFuenteIngreso:  domain.FuenteCampoVerificadoBase,
			wantFuenteRecursos: domain.FuenteCampoDeclarado,
		},
		{
			// Reliability review finding #7: INFERIDO provenance must pass
			// through to the Desglose untouched.
			nombre: "R-7_fuente_inferida",
			perfil: domain.Perfil{
				"ingreso_hogar":    newCampo(int64(2_000_000), domain.FuenteCampoInferido),
				"recursos_propios": newCampo(int64(0), domain.FuenteCampoInferido),
				"tiene_vivienda":   newCampo(false, domain.FuenteCampoVerificadoBase),
				"recibio_subsidio": newCampo(false, domain.FuenteCampoVerificadoBase),
			},
			esAfiliado:         true,
			precioCandidato:    100_000_000,
			wantCredito:        81_726_133,
			wantSubsidio:       52_527_150,
			wantPresupuesto:    134_253_283,
			wantFuenteIngreso:  domain.FuenteCampoInferido,
			wantFuenteRecursos: domain.FuenteCampoInferido,
		},
		{
			nombre: "B-6_precio_candidato_cero",
			perfil: domain.Perfil{
				"ingreso_hogar":    newCampo(int64(2_000_000), domain.FuenteCampoVerificadoBase),
				"recursos_propios": newCampo(int64(0), domain.FuenteCampoVerificadoBase),
				"tiene_vivienda":   newCampo(false, domain.FuenteCampoVerificadoBase),
				"recibio_subsidio": newCampo(false, domain.FuenteCampoVerificadoBase),
			},
			esAfiliado:         true,
			precioCandidato:    0,
			wantCredito:        81_726_133,
			wantSubsidio:       52_527_150,
			wantPresupuesto:    134_253_283,
			wantRatio:          f64(134_253_283.0 / 195_000_000.0),
			wantFuenteIngreso:  domain.FuenteCampoVerificadoBase,
			wantFuenteRecursos: domain.FuenteCampoVerificadoBase,
		},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			assertCaso(t, c)
		})
	}
}

// TestSubsidioAplicable drives the eligibility helper directly for the
// branches the doc 13 fixtures never reach: exact SMMLV bracket edges
// (reliability review finding #4), the affiliate OR-routes (#3), and the
// S2 rule that an absent condition field is never assumed favorable (#6).
func TestSubsidioAplicable(t *testing.T) {
	base := func(ingreso int64) domain.Perfil {
		return domain.Perfil{
			"ingreso_hogar":    newCampo(ingreso, domain.FuenteCampoVerificadoBase),
			"tiene_vivienda":   newCampo(false, domain.FuenteCampoVerificadoBase),
			"recibio_subsidio": newCampo(false, domain.FuenteCampoVerificadoBase),
		}
	}
	dosSMMLV := int64(2) * int64(constantes.SMMLV2026)    // 3_501_810
	cuatroSMMLV := int64(4) * int64(constantes.SMMLV2026) // 7_003_620

	casos := []struct {
		nombre     string
		perfil     domain.Perfil
		esAfiliado bool
		want       int64
	}{
		{"borde_exacto_2_smmlv", base(dosSMMLV), true, 52_527_150},
		{"borde_exacto_4_smmlv", base(cuatroSMMLV), true, 35_018_100},
		{"un_peso_sobre_4_smmlv", base(cuatroSMMLV + 1), true, 0},
		{
			nombre: "ruta_hogar_con_afiliado",
			perfil: func() domain.Perfil {
				p := base(2_000_000)
				p["hogar_con_afiliado"] = newCampo(true, domain.FuenteCampoDeclarado)
				return p
			}(),
			esAfiliado: false,
			want:       52_527_150,
		},
		{
			nombre: "ruta_caja_externa",
			perfil: func() domain.Perfil {
				p := base(2_000_000)
				p["caja_externa"] = newCampo("comfama", domain.FuenteCampoDeclarado)
				return p
			}(),
			esAfiliado: false,
			want:       52_527_150,
		},
		{"sin_ninguna_ruta_de_afiliacion", base(2_000_000), false, 0},
		{
			nombre: "tiene_vivienda_sin_dato_no_favorable",
			perfil: domain.Perfil{
				"ingreso_hogar":    newCampo(int64(2_000_000), domain.FuenteCampoVerificadoBase),
				"recibio_subsidio": newCampo(false, domain.FuenteCampoVerificadoBase),
			},
			esAfiliado: true,
			want:       0,
		},
		{
			nombre: "recibio_subsidio_sin_dato_no_favorable",
			perfil: domain.Perfil{
				"ingreso_hogar":  newCampo(int64(2_000_000), domain.FuenteCampoVerificadoBase),
				"tiene_vivienda": newCampo(false, domain.FuenteCampoVerificadoBase),
			},
			esAfiliado: true,
			want:       0,
		},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			got, _ := subsidioAplicable(c.perfil, c.esAfiliado, constantes, tramosSubsidio)
			if got != c.want {
				t.Errorf("subsidioAplicable = %d, want %d", got, c.want)
			}
		})
	}
}

// TestFactorAnualidad pins doc 13 §1's published annuity factor directly,
// without building a Perfil.
func TestFactorAnualidad(t *testing.T) {
	assertClose(t, "factorAnualidad(0.107, 240)", factorAnualidad(0.107, 240), 102.1576657634)
}

// TestMustLoadPanicsOnMalformedJSON exercises the only otherwise-
// uncoverable branch: a parse failure at load time must panic loudly
// instead of silently degrading to a zeroed Constantes.
func TestMustLoadPanicsOnMalformedJSON(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("mustLoad did not panic on malformed JSON")
		}
	}()
	mustLoad(domain.LoadConstantes([]byte("{{")))
}
