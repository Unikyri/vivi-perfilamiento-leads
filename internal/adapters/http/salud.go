package http

import (
	"encoding/json"
	"net/http"
)

// EstadoSalud es la respuesta de GET /salud (Contrato v1.1 §3.8).
type EstadoSalud struct {
	Estado             string `json:"estado"`
	Version            string `json:"version"`
	LLMProveedorActivo string `json:"llm_proveedor_activo"`
	CircuitBreaker     string `json:"circuit_breaker"`
	BD                 string `json:"bd"`
	FechaSimulada      string `json:"fecha_simulada"`
}

// ProveedorSalud lo implementan las piezas que reportan su estado.
type ProveedorSalud interface {
	Estado() EstadoSalud
}

// HandlerSalud devuelve el handler de GET /salud.
func HandlerSalud(p ProveedorSalud) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(p.Estado())
	}
}
