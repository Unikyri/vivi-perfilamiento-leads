package llm

import (
	"strings"
	"testing"
)

func TestParserContractValidation(t *testing.T) {
	valid := `{"campos_extraidos":[{"campo":"ingreso_hogar","valor":100000,"fuente":"DECLARADO","confianza":0.8,"requiere_confirmacion":true}],"intencion":{"nivel":"ALTA","confianza":"MEDIA","senales":["compra"]},"respuesta":"Presupuesto $100.000","accion":"CONTINUAR"}`
	cases := []struct {
		name, data string
		kind       ErrorKind
	}{{"valid", valid, ""}, {"unknown/redacted", strings.Replace(valid, `"accion"`, `"secret":"secret-key-123","accion"`, 1), KindInvalidOutput}, {"nested unknown", strings.Replace(valid, `"senales"`, `"extra":1,"senales"`, 1), KindInvalidOutput}, {"source", strings.Replace(valid, "DECLARADO", "VERIFICADO_BASE", 1), KindInvalidOutput}, {"confidence", strings.Replace(valid, "0.8", "1.1", 1), KindInvalidOutput}, {"action", strings.Replace(valid, "CONTINUAR", "DELETE", 1), KindInvalidOutput}, {"motor", strings.Replace(valid, "100.000", "999.000", 1), KindInvalidOutput}, {"malformed", "{", KindMalformed}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := ParseSalida([]byte(tc.data), map[string]int64{"presupuesto": 100000})
			if tc.kind == "" && (err != nil || out.Accion != "CONTINUAR") {
				t.Fatalf("valid: %#v %v", out, err)
			}
			if tc.kind != "" && (ErrorKindOf(err) != tc.kind || strings.Contains(err.Error(), "secret-key-123")) {
				t.Fatalf("kind=%q err=%v", ErrorKindOf(err), err)
			}
		})
	}
}

func TestParserRestrictsOnlyCurrencyAmounts(t *testing.T) {
	data := `{"campos_extraidos":[],"intencion":{"nivel":"ALTA","confianza":"MEDIA","senales":["compra"]},"respuesta":"Tienes 2 personas desde 2026-07-25. Presupuesto $100.000","accion":"CONTINUAR"}`
	if _, err := ParseSalida([]byte(data), map[string]int64{"presupuesto": 100000}); err != nil {
		t.Fatalf("non-currency numbers rejected: %v", err)
	}
	bad := strings.Replace(data, "$100.000", "$999.000", 1)
	if _, err := ParseSalida([]byte(bad), map[string]int64{"presupuesto": 100000}); ErrorKindOf(err) != KindInvalidOutput {
		t.Fatalf("currency outside motor kind=%q err=%v", ErrorKindOf(err), err)
	}
}

func TestParserColombianMillionAmounts(t *testing.T) {
	base := `{"campos_extraidos":[],"intencion":{"nivel":"ALTA","confianza":"MEDIA","senales":["compra"]},"respuesta":"Presupuesto 2 millones","accion":"CONTINUAR"}`
	for _, expression := range []string{"2 millones", "$2 millones", "COP 2 millones"} {
		data := strings.Replace(base, "2 millones", expression, 1)
		if _, err := ParseSalida([]byte(data), map[string]int64{"presupuesto": 2000000}); err != nil {
			t.Fatalf("authorized expression %q rejected: %v", expression, err)
		}
	}
	unauthorized := strings.Replace(base, "2 millones", "3 millones", 1)
	if _, err := ParseSalida([]byte(unauthorized), map[string]int64{"presupuesto": 2000000}); ErrorKindOf(err) != KindInvalidOutput {
		t.Fatalf("unauthorized millions kind=%q err=%v", ErrorKindOf(err), err)
	}
	ordinary := strings.Replace(base, "Presupuesto 2 millones", "2 personas desde 2026-07-25", 1)
	if _, err := ParseSalida([]byte(ordinary), map[string]int64{"presupuesto": 2000000}); err != nil {
		t.Fatalf("ordinary count/date rejected: %v", err)
	}
}

func TestParserBareGroupedAmountsRequireMonetaryContext(t *testing.T) {
	base := `{"campos_extraidos":[],"intencion":{"nivel":"ALTA","confianza":"MEDIA","senales":["compra"]},"respuesta":"RESPONSE","accion":"CONTINUAR"}`
	motor := map[string]int64{"presupuesto": 3500000}
	for _, response := range []string{
		"Ingreso mensual: 3.500.000",
		"3.500.000 mensuales",
		"Cuota de 3.500.000",
		"Salario 3.500.000",
		"Presupuesto 3.500.000",
		"Precio 3.500.000",
		"Valor: 3.500.000.",
	} {
		data := strings.Replace(base, "RESPONSE", response, 1)
		if _, err := ParseSalida([]byte(data), motor); err != nil {
			t.Errorf("authorized contextual amount %q rejected: %v", response, err)
		}
	}

	unauthorized := strings.Replace(base, "RESPONSE", "Presupuesto: 3.500.001.", 1)
	if _, err := ParseSalida([]byte(unauthorized), motor); ErrorKindOf(err) != KindInvalidOutput {
		t.Fatalf("unauthorized contextual amount kind=%q err=%v", ErrorKindOf(err), err)
	}

	for _, response := range []string{
		"3.500.000 personas desde 2026-07-25",
		"3.500.000",
	} {
		data := strings.Replace(base, "RESPONSE", response, 1)
		if _, err := ParseSalida([]byte(data), motor); err != nil {
			t.Errorf("ordinary non-monetary text %q rejected: %v", response, err)
		}
	}
}

func TestParserSlice7ColombianDecimalAndUngroupedAmounts(t *testing.T) {
	base := `{"campos_extraidos":[],"intencion":{"nivel":"ALTA","confianza":"MEDIA","senales":["compra"]},"respuesta":"RESPONSE","accion":"CONTINUAR"}`
	motor := map[string]int64{"presupuesto": 3500000}
	for _, response := range []string{
		"Presupuesto: 3.500.000,00",
		"Presupuesto para vivienda es 3500000.",
		"3500000 de presupuesto",
	} {
		data := strings.Replace(base, "RESPONSE", response, 1)
		if _, err := ParseSalida([]byte(data), motor); err != nil {
			t.Errorf("authorized amount %q rejected: %v", response, err)
		}
	}
	for _, response := range []string{
		"Presupuesto: 3.500.000,50",
		"Presupuesto: 3500001",
	} {
		if ErrorKindOf(ParseError(response, motor)) != KindInvalidOutput {
			t.Errorf("unauthorized amount %q was accepted", response)
		}
	}
}

func TestParserSlice7ContextDoesNotCrossSentenceOrLinkDistantNumber(t *testing.T) {
	base := `{"campos_extraidos":[],"intencion":{"nivel":"ALTA","confianza":"MEDIA","senales":["compra"]},"respuesta":"RESPONSE","accion":"CONTINUAR"}`
	motor := map[string]int64{"presupuesto": 3500000}
	for _, response := range []string{
		"Presupuesto: 3.500.000. Tengo 3.500.001 personas.",
		"Mi presupuesto está en revisión y luego 3500000 personas.",
		"Presupuesto para el año 2026. Tengo 3.500.000 personas.",
	} {
		data := strings.Replace(base, "RESPONSE", response, 1)
		if _, err := ParseSalida([]byte(data), motor); err != nil {
			t.Errorf("distant ordinary number %q was linked to context: %v", response, err)
		}
	}
}

func ParseError(response string, motor map[string]int64) error {
	base := `{"campos_extraidos":[],"intencion":{"nivel":"ALTA","confianza":"MEDIA","senales":["compra"]},"respuesta":"RESPONSE","accion":"CONTINUAR"}`
	_, err := ParseSalida([]byte(strings.Replace(base, "RESPONSE", response, 1)), motor)
	return err
}

func TestParserInt64OverflowIsInvalidWithoutWrap(t *testing.T) {
	response := "Presupuesto: 9223372036854775808"
	// The wrapped value would be MinInt64; keeping it in the fake motor makes
	// accidental signed overflow observable instead of merely out of range.
	if got := ErrorKindOf(ParseError(response, map[string]int64{"presupuesto": -9223372036854775808})); got != KindInvalidOutput {
		t.Fatalf("overflow kind=%q; want %q without signed wrap", got, KindInvalidOutput)
	}
}
