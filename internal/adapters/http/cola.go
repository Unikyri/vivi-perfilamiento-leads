package http

import (
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
	"net/http"
	"net/url"
	"strings"
)

func (c *Controlador) cola(w http.ResponseWriter, r *http.Request) {
	filtro, err := filtroLeads(r.URL.Query())
	if err != nil {
		writeError(w, err)
		return
	}
	if c.consulta == nil {
		writeError(w, usecase.ErrValidacion)
		return
	}
	out, err := c.consulta.Ejecutar(r.Context(), filtro)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}
func (c *Controlador) detalleLead(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("lead_id"))
	if c.consulta == nil {
		writeError(w, usecase.ErrValidacion)
		return
	}
	out, err := c.consulta.Detalle(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}
func (c *Controlador) ficha(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("lead_id"))
	if c.consulta == nil {
		writeError(w, usecase.ErrValidacion)
		return
	}
	out, err := c.consulta.Ficha(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}
func filtroLeads(query url.Values) (usecase.FiltroLeads, error) {
	filtro := usecase.FiltroLeads{}
	if raw, ok := query["afiliado"]; ok {
		if len(raw) != 1 || (raw[0] != "true" && raw[0] != "false") {
			return filtro, usecase.ErrValidacion
		}
		value := raw[0] == "true"
		filtro.Afiliado = &value
	}
	if raw, ok := query["ruta"]; ok {
		if len(raw) != 1 {
			return filtro, usecase.ErrValidacion
		}
		route := domain.Ruta(raw[0])
		switch route {
		case domain.RutaAsesor, domain.RutaNutricion, domain.RutaRemarketing, domain.RutaDespedida:
			filtro.Ruta = &route
		default:
			return filtro, usecase.ErrValidacion
		}
	}
	return filtro, nil
}
