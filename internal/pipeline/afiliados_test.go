package pipeline

import (
	"encoding/json"
	"os"
	"regexp"
	"testing"
)

type afiliado struct {
	Cedula         string `json:"cedula"`
	Nombre         string `json:"nombre"`
	Categoria      string `json:"categoria"`
	IngresoMensual int64  `json:"ingreso_mensual"`
}

func cargar(t *testing.T) []afiliado {
	t.Helper()
	b, err := os.ReadFile("../../data/afiliados_mock.json")
	if err != nil {
		t.Fatalf("no se pudo leer afiliados_mock.json: %v", err)
	}
	var as []afiliado
	if err := json.Unmarshal(b, &as); err != nil {
		t.Fatalf("JSON inválido: %v", err)
	}
	return as
}

func TestAfiliadosMinimo10(t *testing.T) {
	if as := cargar(t); len(as) < 10 {
		t.Errorf("hay %d afiliados, mínimo 10 (RF-M1-03)", len(as))
	}
}

func TestAfiliadosObligatorios(t *testing.T) {
	as := cargar(t)
	requeridas := map[string]string{
		"1032456789": "Ana",
		"1015789456": "", // esposa de Carlos (UC-03)
		"1098765432": "Luisa",
	}
	encontradas := map[string]bool{}
	for _, a := range as {
		encontradas[a.Cedula] = true
	}
	for ced := range requeridas {
		if !encontradas[ced] {
			t.Errorf("falta la cédula obligatoria %s", ced)
		}
	}
}

func TestCedulasSinteticasValidas(t *testing.T) {
	re := regexp.MustCompile(`^\d{6,12}$`)
	for _, a := range cargar(t) {
		if !re.MatchString(a.Cedula) {
			t.Errorf("cédula %q no cumple el formato 6-12 dígitos (Contrato §0)", a.Cedula)
		}
	}
}

func TestAnaTieneLosDatosDelCasoDeReferencia(t *testing.T) {
	for _, a := range cargar(t) {
		if a.Cedula == "1032456789" {
			if a.IngresoMensual != 2600000 {
				t.Errorf("ingreso de Ana = %d, esperado 2600000 (doc 13 §2.2)", a.IngresoMensual)
			}
			if a.Categoria != "A" {
				t.Errorf("categoría de Ana = %q, esperada A", a.Categoria)
			}
			return
		}
	}
	t.Fatal("no se encontró a Ana")
}
