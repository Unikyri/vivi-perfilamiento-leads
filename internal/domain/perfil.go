package domain

import "time"

// CampoPerfil stores a recognized profile value and its provenance.
type CampoPerfil struct {
	Valor                any         `json:"valor"`
	Fuente               FuenteCampo `json:"fuente"`
	Confianza            float64     `json:"confianza"`
	RequiereConfirmacion bool        `json:"requiere_confirmacion"`
	ActualizadoEn        time.Time   `json:"actualizado_en"`
}

// Perfil contains dynamic profile fields keyed by their Contract names.
type Perfil map[string]CampoPerfil

// CamposReconocidos is the complete set of profile keys accepted by the Contract.
var CamposReconocidos = map[string]bool{
	"ingreso_hogar":            true,
	"recursos_propios":         true,
	"personas_hogar":           true,
	"tipo_hogar":               true,
	"edad":                     true,
	"zona_deseada":             true,
	"plazo_compra_meses":       true,
	"arriendo_actual":          true,
	"tiene_vivienda":           true,
	"recibio_subsidio":         true,
	"reporte_credito":          true,
	"situacion_laboral":        true,
	"caja_externa":             true,
	"hogar_con_afiliado":       true,
	"cedula_familiar_afiliado": true,
	"cesantias_entidad":        true,
	"categoria":                true,
	"segmento":                 true,
}

// CamposCriticos contains the profile keys required for the initial assessment.
var CamposCriticos = map[string]bool{
	"ingreso_hogar":    true,
	"recursos_propios": true,
	"tiene_vivienda":   true,
	"recibio_subsidio": true,
}

// Entero returns an integer profile value from native or JSON-decoded numeric forms.
func (p Perfil) Entero(clave string) (int64, bool) {
	campo, ok := p[clave]
	if !ok {
		return 0, false
	}

	switch valor := campo.Valor.(type) {
	case int64:
		return valor, true
	case int:
		return int64(valor), true
	case float64:
		return int64(valor), true
	default:
		return 0, false
	}
}

// Booleano returns a boolean profile value.
func (p Perfil) Booleano(clave string) (bool, bool) {
	campo, ok := p[clave]
	if !ok {
		return false, false
	}
	valor, ok := campo.Valor.(bool)
	return valor, ok
}

// Texto returns a string profile value.
func (p Perfil) Texto(clave string) (string, bool) {
	campo, ok := p[clave]
	if !ok {
		return "", false
	}
	valor, ok := campo.Valor.(string)
	return valor, ok
}

// EsVerificado reports whether a profile field came from the verified base.
func (p Perfil) EsVerificado(clave string) bool {
	campo, ok := p[clave]
	return ok && campo.Fuente == FuenteCampoVerificadoBase
}
