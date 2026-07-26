package usecase

import "github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"

var prioridadCampos = []string{
	"ingreso_hogar", "recursos_propios", "zona_deseada",
	"plazo_compra_meses", "arriendo_actual", "personas_hogar",
}

var camposCriticosPerfil = map[string]bool{
	"ingreso_hogar": true, "recursos_propios": true, "zona_deseada": true,
}

// SiguienteMejorPregunta returns the highest-priority missing conversational field.
func SiguienteMejorPregunta(perfil domain.Perfil) string {
	for _, campo := range prioridadCampos {
		valor, ok := perfil[campo]
		if ok && valor.Fuente == domain.FuenteCampoVerificadoBase {
			continue
		}
		if !ok || valor.Valor == nil {
			return campo
		}
	}
	return ""
}

// PerfilEstaCompleto checks the conversational gate without changing motor criteria.
func PerfilEstaCompleto(perfil domain.Perfil) bool {
	for campo := range camposCriticosPerfil {
		valor, ok := perfil[campo]
		if !ok || valor.Valor == nil {
			return false
		}
		if texto, ok := valor.Valor.(string); ok && texto == "" {
			return false
		}
	}
	return true
}
