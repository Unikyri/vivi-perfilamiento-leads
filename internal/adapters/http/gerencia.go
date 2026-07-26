package http

import (
	"net/http"
	"strings"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
)

func (c *Controlador) gerenciaBuyerPersona(w http.ResponseWriter, r *http.Request) {
	if c.buyerPersona == nil {
		writeError(w, usecase.ErrValidacion)
		return
	}
	values, supplied := r.URL.Query()["proyecto_id"]
	if supplied && (len(values) != 1 || strings.TrimSpace(values[0]) == "") {
		writeError(w, usecase.ErrValidacion)
		return
	}
	if supplied {
		out, err := c.buyerPersona.Proyecto(r.Context(), strings.TrimSpace(values[0]))
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, out)
		return
	}
	out, err := c.buyerPersona.CatalogoCompleto(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}
