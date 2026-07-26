package http

import (
	"encoding/json"
	"errors"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
	"net/http"
	"strings"
)

type apiError struct {
	Codigo   string         `json:"codigo"`
	Mensaje  string         `json:"mensaje"`
	Detalles map[string]any `json:"detalles"`
}
type errorEnvelope struct {
	Error apiError `json:"error"`
}

var errorCatalog = map[string]struct {
	status  int
	message string
}{
	"VALIDACION":                  {http.StatusBadRequest, "La solicitud no es válida."},
	"LEAD_NO_ENCONTRADO":          {http.StatusNotFound, "Lead no encontrado."},
	"FICHA_NO_DISPONIBLE":         {http.StatusNotFound, "Ficha no disponible."},
	"TRANSICION_INVALIDA":         {http.StatusConflict, "La transición solicitada no es válida."},
	"AUDIO_INVALIDO":              {http.StatusUnprocessableEntity, "El audio no es válido."},
	"LIMITE_TASA":                 {http.StatusTooManyRequests, "Se alcanzó el límite de solicitudes."},
	"PROVEEDOR_LLM_NO_DISPONIBLE": {http.StatusServiceUnavailable, "El proveedor no está disponible."},
	"ERROR_INTERNO":               {http.StatusInternalServerError, "Ocurrió un error interno."},
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func writeError(w http.ResponseWriter, err error) {
	code := errorCode(err)
	entry, ok := errorCatalog[code]
	if !ok {
		code, entry = "ERROR_INTERNO", errorCatalog["ERROR_INTERNO"]
	}
	writeJSON(w, entry.status, errorEnvelope{Error: apiError{Codigo: code, Mensaje: entry.message, Detalles: map[string]any{}}})
}
func errorCode(err error) string {
	if err == nil {
		return "ERROR_INTERNO"
	}
	if errors.Is(err, ErrTurnoEnProceso) {
		return "LIMITE_TASA"
	}
	if errors.Is(err, usecase.ErrValidacion) {
		return "VALIDACION"
	}
	if errors.Is(err, usecase.ErrAudioInvalido) {
		return "AUDIO_INVALIDO"
	}
	if errors.Is(err, usecase.ErrLeadNoPerfilando) {
		return "TRANSICION_INVALIDA"
	}
	var transition domain.ErrTransicionInvalida
	if errors.As(err, &transition) {
		return "TRANSICION_INVALIDA"
	}
	var notFound *usecase.NotFoundError
	if errors.As(err, &notFound) && strings.EqualFold(notFound.Resource, "ficha") {
		return "FICHA_NO_DISPONIBLE"
	}
	if errors.As(err, &notFound) && strings.EqualFold(notFound.Resource, "lead") {
		return "LEAD_NO_ENCONTRADO"
	}
	if errors.Is(err, usecase.ErrTiempoSimuladoAtras) {
		return "VALIDACION"
	}
	if errors.Is(err, usecase.ErrNoEncontrado) {
		return "LEAD_NO_ENCONTRADO"
	}
	return "ERROR_INTERNO"
}
