package http

import (
	"encoding/json"
	"errors"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
	"io"
	"net/http"
	"strings"
)

type crearLeadRequest struct {
	Nombre       string `json:"nombre"`
	Telefono     string `json:"telefono"`
	Cedula       string `json:"cedula"`
	Fuente       string `json:"fuente"`
	PrecargadoID string `json:"precargado_id"`
}
type crearLeadResponse struct {
	LeadID            string            `json:"lead_id"`
	Estado            domain.EstadoLead `json:"estado"`
	AfiliadoDetectado bool              `json:"afiliado_detectado"`
}
type conversacionResponse struct {
	LeadID         string            `json:"lead_id"`
	Estado         domain.EstadoLead `json:"estado"`
	TurnoEnProceso bool              `json:"turno_en_proceso"`
	Mensajes       []domain.Mensaje  `json:"mensajes"`
}
type seedLead struct{ nombre, cedula, telefono string }

var seeds = map[string]seedLead{
	"ana":    {"Ana Rodríguez", "1032456789", "+57 300 123 4567"},
	"carlos": {"Carlos Martínez", "1000000000", "+57 311 987 6543"},
	"luisa":  {"Luisa Gómez", "1098765432", "+57 300 000 0000"},
}

func (c *Controlador) crearLead(w http.ResponseWriter, r *http.Request) {
	var req crearLeadRequest
	if err := decodeJSON(w, r, &req); err != nil {
		return
	}
	if strings.TrimSpace(req.PrecargadoID) != "" {
		seed, ok := seeds[strings.ToLower(strings.TrimSpace(req.PrecargadoID))]
		if !ok {
			writeError(w, usecase.ErrValidacion)
			return
		}
		req.Nombre, req.Cedula, req.Telefono = seed.nombre, seed.cedula, seed.telefono
	}
	if strings.TrimSpace(req.Nombre) == "" {
		writeError(w, usecase.ErrValidacion)
		return
	}
	if strings.TrimSpace(req.Fuente) == "" {
		req.Fuente = "DEMO"
	}
	if req.Fuente != "META_LEAD_FORM" && req.Fuente != "CLICK_TO_WHATSAPP" && req.Fuente != "DEMO" {
		writeError(w, usecase.ErrValidacion)
		return
	}
	out, err := c.perfilador.Ejecutar(r.Context(), usecase.EntradaPerfilar{Nombre: req.Nombre, Telefono: req.Telefono, Cedula: req.Cedula, Fuente: req.Fuente})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, crearLeadResponse{LeadID: out.LeadID, Estado: out.Estado, AfiliadoDetectado: out.AfiliadoDetectado})
}
func (c *Controlador) conversacion(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("lead_id"))
	if id == "" || strings.Contains(id, "/") {
		writeError(w, usecase.ErrNoEncontrado)
		return
	}
	lead, err := c.leads.PorID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	messages, err := c.leads.Conversacion(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	if messages == nil {
		messages = []domain.Mensaje{}
	}
	turnoActivo := false
	if c.turnos != nil {
		turnoActivo = c.turnos.Activo(id)
	}
	writeJSON(w, http.StatusOK, conversacionResponse{LeadID: lead.LeadID, Estado: lead.Estado, TurnoEnProceso: turnoActivo, Mensajes: messages})
}
func decodeJSON(w http.ResponseWriter, r *http.Request, target any) error {
	return decodeJSONLimit(w, r, target, 1<<20)
}
func decodeJSONLimit(w http.ResponseWriter, r *http.Request, target any, limit int64) error {
	r.Body = http.MaxBytesReader(w, r.Body, limit)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		writeError(w, usecase.ErrValidacion)
		return err
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		writeError(w, usecase.ErrValidacion)
		return usecase.ErrValidacion
	}
	return nil
}
