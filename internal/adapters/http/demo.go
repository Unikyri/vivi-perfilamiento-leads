package http

import (
	"net/http"
	"strings"
	"time"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
)

type tiempoRequest struct {
	AvanzarHasta *string `json:"avanzar_hasta"`
	AvanzarDias  *int    `json:"avanzar_dias"`
}

type tiempoResponse struct {
	FechaSimulada   string `json:"fecha_simulada"`
	HitosDisparados int    `json:"hitos_disparados"`
}

func (c *Controlador) avanzarTiempo(w http.ResponseWriter, r *http.Request) {
	if c.avanzarDemo == nil {
		writeError(w, usecase.ErrValidacion)
		return
	}
	var req tiempoRequest
	if err := decodeJSON(w, r, &req); err != nil {
		return
	}
	entrada := usecase.EntradaAvanzarDemo{Dias: req.AvanzarDias}
	if req.AvanzarHasta != nil {
		parsed, err := parseDemoTime(strings.TrimSpace(*req.AvanzarHasta))
		if err != nil {
			writeError(w, usecase.ErrValidacion)
			return
		}
		entrada.Hasta = &parsed
	}
	out, err := c.avanzarDemo.Ejecutar(r.Context(), entrada)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, tiempoResponse{FechaSimulada: out.FechaSimulada.UTC().Format("2006-01-02"), HitosDisparados: out.HitosDisparados})
}

func parseDemoTime(value string) (time.Time, error) {
	if parsed, err := time.Parse(time.RFC3339, value); err == nil {
		return parsed, nil
	}
	return time.ParseInLocation("2006-01-02", value, time.UTC)
}

type reiniciarResponse struct {
	Reiniciado    bool   `json:"reiniciado"`
	FechaSimulada string `json:"fecha_simulada"`
}

func (c *Controlador) reiniciar(w http.ResponseWriter, r *http.Request) {
	if c.reiniciarDemo == nil {
		writeError(w, usecase.ErrDemoDeshabilitado)
		return
	}
	out, err := c.reiniciarDemo.Ejecutar(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, reiniciarResponse{Reiniciado: out.Reiniciado, FechaSimulada: out.FechaSimulada.UTC().Format("2006-01-02")})
}
