package domain

import (
	"encoding/json"
	"fmt"
)

// Constantes holds the normative economic parameters from Contract §4.5.
type Constantes struct {
	SMMLV2026         int     `json:"smmlv_2026"`
	TasaEADefault     float64 `json:"tasa_ea_default"`
	PlazoMesesDefault int     `json:"plazo_meses_default"`
	CuotaIngresoMax   float64 `json:"cuota_ingreso_max"`
	TopeVISSMMLV      int     `json:"tope_vis_smmlv"`
	MedianaVISCOP     int64   `json:"mediana_vis_cop"`
}

// TramoSubsidio represents a subsidy bracket by income range.
type TramoSubsidio struct {
	IngresoHastaSMMLV int `json:"ingreso_hasta_smmlv"`
	SubsidioSMMLV     int `json:"subsidio_smmlv"`
}

// EventoCalendario represents a recurring annual economic event.
type EventoCalendario struct {
	Tipo  string `json:"tipo"`
	Fecha string `json:"fecha"`
}

// LoadConstantes parses raw JSON bytes into Constantes.
func LoadConstantes(data []byte) (Constantes, error) {
	if len(data) == 0 {
		return Constantes{}, fmt.Errorf("domain: parse constantes: %w", fmt.Errorf("empty input"))
	}
	var c Constantes
	if err := json.Unmarshal(data, &c); err != nil {
		return Constantes{}, fmt.Errorf("domain: parse constantes: %w", err)
	}
	return c, nil
}

// LoadSubsidios parses raw JSON bytes into a slice of TramoSubsidio.
func LoadSubsidios(data []byte) ([]TramoSubsidio, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("domain: parse subsidios: %w", fmt.Errorf("empty input"))
	}
	var s []TramoSubsidio
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("domain: parse subsidios: %w", err)
	}
	return s, nil
}

// LoadCalendario parses raw JSON bytes into a slice of EventoCalendario.
func LoadCalendario(data []byte) ([]EventoCalendario, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("domain: parse calendario: %w", fmt.Errorf("empty input"))
	}
	var e []EventoCalendario
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, fmt.Errorf("domain: parse calendario: %w", err)
	}
	return e, nil
}
