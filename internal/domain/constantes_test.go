package domain

import (
	"testing"
)

func TestLoadConstantes(t *testing.T) {
	input := []byte(`{"smmlv_2026":1750905,"tasa_ea_default":0.107,"plazo_meses_default":240,"cuota_ingreso_max":0.40,"tope_vis_smmlv":150,"mediana_vis_cop":195000000}`)

	tests := []struct {
		name  string
		field string
		got   func(Constantes) interface{}
		want  interface{}
	}{
		{"SMMLV2026", "smmlv_2026", func(c Constantes) interface{} { return c.SMMLV2026 }, 1750905},
		{"TasaEADefault", "tasa_ea_default", func(c Constantes) interface{} { return c.TasaEADefault }, 0.107},
		{"PlazoMesesDefault", "plazo_meses_default", func(c Constantes) interface{} { return c.PlazoMesesDefault }, 240},
		{"CuotaIngresoMax", "cuota_ingreso_max", func(c Constantes) interface{} { return c.CuotaIngresoMax }, 0.40},
		{"TopeVISSMMLV", "tope_vis_smmlv", func(c Constantes) interface{} { return c.TopeVISSMMLV }, 150},
		{"MedianaVISCOP", "mediana_vis_cop", func(c Constantes) interface{} { return c.MedianaVISCOP }, int64(195000000)},
	}

	c, err := LoadConstantes(input)
	if err != nil {
		t.Fatalf("LoadConstantes() unexpected error: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.got(c)
			if got != tt.want {
				t.Errorf("Constantes.%s = %v (%T), want %v (%T)", tt.name, got, got, tt.want, tt.want)
			}
		})
	}
}

func TestLoadSubsidios(t *testing.T) {
	input := []byte(`[{"ingreso_hasta_smmlv":2,"subsidio_smmlv":30},{"ingreso_hasta_smmlv":4,"subsidio_smmlv":20}]`)

	tests := []struct {
		name         string
		index        int
		wantIngreso  int
		wantSubsidio int
	}{
		{"Tramo1_hasta2SMMLV_subsidio30", 0, 2, 30},
		{"Tramo2_hasta4SMMLV_subsidio20", 1, 4, 20},
	}

	s, err := LoadSubsidios(input)
	if err != nil {
		t.Fatalf("LoadSubsidios() unexpected error: %v", err)
	}
	if len(s) != 2 {
		t.Fatalf("LoadSubsidios() len = %d, want 2", len(s))
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if s[tt.index].IngresoHastaSMMLV != tt.wantIngreso {
				t.Errorf("TramoSubsidio[%d].IngresoHastaSMMLV = %d, want %d", tt.index, s[tt.index].IngresoHastaSMMLV, tt.wantIngreso)
			}
			if s[tt.index].SubsidioSMMLV != tt.wantSubsidio {
				t.Errorf("TramoSubsidio[%d].SubsidioSMMLV = %d, want %d", tt.index, s[tt.index].SubsidioSMMLV, tt.wantSubsidio)
			}
		})
	}
}

func TestLoadCalendario(t *testing.T) {
	input := []byte(`[{"tipo":"CESANTIAS","fecha":"--02-14"},{"tipo":"PRIMA","fecha":"--06-30"},{"tipo":"PRIMA","fecha":"--12-20"}]`)

	tests := []struct {
		name      string
		index     int
		wantTipo  string
		wantFecha string
	}{
		{"CESANTIAS_feb14", 0, "CESANTIAS", "--02-14"},
		{"PRIMA_jun30", 1, "PRIMA", "--06-30"},
		{"PRIMA_dic20", 2, "PRIMA", "--12-20"},
	}

	e, err := LoadCalendario(input)
	if err != nil {
		t.Fatalf("LoadCalendario() unexpected error: %v", err)
	}
	if len(e) != 3 {
		t.Fatalf("LoadCalendario() len = %d, want 3", len(e))
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if e[tt.index].Tipo != tt.wantTipo {
				t.Errorf("EventoCalendario[%d].Tipo = %q, want %q", tt.index, e[tt.index].Tipo, tt.wantTipo)
			}
			if e[tt.index].Fecha != tt.wantFecha {
				t.Errorf("EventoCalendario[%d].Fecha = %q, want %q", tt.index, e[tt.index].Fecha, tt.wantFecha)
			}
		})
	}
}

func TestLoad_MalformedJSON(t *testing.T) {
	malformed := []byte("{{invalid")

	tests := []struct {
		name   string
		loader func([]byte) error
	}{
		{"LoadConstantes_malformed", func(b []byte) error { _, err := LoadConstantes(b); return err }},
		{"LoadSubsidios_malformed", func(b []byte) error { _, err := LoadSubsidios(b); return err }},
		{"LoadCalendario_malformed", func(b []byte) error { _, err := LoadCalendario(b); return err }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.loader(malformed)
			if err == nil {
				t.Errorf("%s: expected non-nil error for malformed JSON, got nil", tt.name)
			}
		})
	}
}

func TestLoad_EmptyInput(t *testing.T) {
	inputs := []struct {
		name string
		data []byte
	}{
		{"nil", nil},
		{"empty_slice", []byte("")},
	}

	loaders := []struct {
		name   string
		loader func([]byte) error
	}{
		{"LoadConstantes", func(b []byte) error { _, err := LoadConstantes(b); return err }},
		{"LoadSubsidios", func(b []byte) error { _, err := LoadSubsidios(b); return err }},
		{"LoadCalendario", func(b []byte) error { _, err := LoadCalendario(b); return err }},
	}

	for _, input := range inputs {
		for _, l := range loaders {
			t.Run(l.name+"_"+input.name, func(t *testing.T) {
				err := l.loader(input.data)
				if err == nil {
					t.Errorf("%s(%s): expected non-nil error, got nil", l.name, input.name)
				}
			})
		}
	}
}
