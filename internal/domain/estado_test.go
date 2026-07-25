package domain

import (
	"errors"
	"strings"
	"testing"
)

type transicionEsperada struct {
	desde EstadoLead
	hacia EstadoLead
}

func TestPuedeTransicionarAceptaEdgesDelContrato(t *testing.T) {
	edges := []transicionEsperada{
		{EstadoLeadNuevo, EstadoLeadPerfilando},
		{EstadoLeadPerfilando, EstadoLeadCalificado},
		{EstadoLeadCalificado, EstadoLeadEntregado},
		{EstadoLeadCalificado, EstadoLeadEnNutricion},
		{EstadoLeadCalificado, EstadoLeadRemarketing},
		{EstadoLeadCalificado, EstadoLeadDespedido},
		{EstadoLeadEnNutricion, EstadoLeadPausado},
		{EstadoLeadEnNutricion, EstadoLeadPerfilando},
		{EstadoLeadPausado, EstadoLeadEnNutricion},
		{EstadoLeadRemarketing, EstadoLeadPerfilando},
		{EstadoLeadEntregado, EstadoLeadCerrado},
	}

	for _, edge := range edges {
		t.Run(string(edge.desde)+"_a_"+string(edge.hacia), func(t *testing.T) {
			if !PuedeTransicionar(edge.desde, edge.hacia) {
				t.Fatalf("PuedeTransicionar(%q, %q) = false, want true", edge.desde, edge.hacia)
			}

			lead := Lead{Estado: edge.desde}
			if err := lead.Transicionar(edge.hacia); err != nil {
				t.Fatalf("Transicionar(%q) returned error: %v", edge.hacia, err)
			}
			if lead.Estado != edge.hacia {
				t.Fatalf("lead.Estado = %q, want %q", lead.Estado, edge.hacia)
			}
		})
	}
}

func TestEstadosPosiblesRespetaOrdenDelContrato(t *testing.T) {
	tests := []struct {
		desde EstadoLead
		hacia []EstadoLead
	}{
		{EstadoLeadNuevo, []EstadoLead{EstadoLeadPerfilando}},
		{EstadoLeadPerfilando, []EstadoLead{EstadoLeadCalificado}},
		{EstadoLeadCalificado, []EstadoLead{
			EstadoLeadEntregado,
			EstadoLeadEnNutricion,
			EstadoLeadRemarketing,
			EstadoLeadDespedido,
		}},
		{EstadoLeadEnNutricion, []EstadoLead{EstadoLeadPausado, EstadoLeadPerfilando}},
		{EstadoLeadPausado, []EstadoLead{EstadoLeadEnNutricion}},
		{EstadoLeadRemarketing, []EstadoLead{EstadoLeadPerfilando}},
		{EstadoLeadEntregado, []EstadoLead{EstadoLeadCerrado}},
	}

	for _, tt := range tests {
		t.Run(string(tt.desde), func(t *testing.T) {
			got := EstadosPosibles(tt.desde)
			if !estadosIguales(got, tt.hacia) {
				t.Fatalf("EstadosPosibles(%q) = %v, want %v", tt.desde, got, tt.hacia)
			}
		})
	}
}

func TestTransicionarRechazaParesNoListadosDeFormaAtomica(t *testing.T) {
	estados := []EstadoLead{
		EstadoLeadNuevo,
		EstadoLeadPerfilando,
		EstadoLeadCalificado,
		EstadoLeadEntregado,
		EstadoLeadEnNutricion,
		EstadoLeadPausado,
		EstadoLeadRemarketing,
		EstadoLeadDespedido,
		EstadoLeadCerrado,
		EstadoLead(""),
		EstadoLead("DESCONOCIDO"),
	}

	for _, desde := range estados {
		for _, hacia := range estados {
			if PuedeTransicionar(desde, hacia) {
				continue
			}
			t.Run(string(desde)+"_a_"+string(hacia), func(t *testing.T) {
				lead := Lead{Estado: desde}
				err := lead.Transicionar(hacia)
				if err == nil {
					t.Fatal("Transicionar returned nil for invalid edge")
				}

				var invalid ErrTransicionInvalida
				if !errors.As(err, &invalid) {
					t.Fatalf("errors.As(%T) = false, want ErrTransicionInvalida", err)
				}
				if invalid.Desde != desde || invalid.Hacia != hacia {
					t.Fatalf("error = {%q, %q}, want {%q, %q}", invalid.Desde, invalid.Hacia, desde, hacia)
				}
				if !strings.Contains(err.Error(), "TRANSICION_INVALIDA") {
					t.Fatalf("error = %q, want TRANSICION_INVALIDA", err)
				}
				if lead.Estado != desde {
					t.Fatalf("lead.Estado = %q after rejection, want unchanged %q", lead.Estado, desde)
				}
			})
		}
	}
}

func TestEstadosTerminalesYDesconocidosSonInalcanzables(t *testing.T) {
	for _, estado := range []EstadoLead{
		EstadoLeadCerrado,
		EstadoLeadDespedido,
		EstadoLead("DESCONOCIDO"),
	} {
		t.Run(string(estado), func(t *testing.T) {
			if got := EstadosPosibles(estado); len(got) != 0 {
				t.Fatalf("EstadosPosibles(%q) = %v, want zero-length result", estado, got)
			}
			if PuedeTransicionar(estado, EstadoLeadNuevo) {
				t.Fatalf("PuedeTransicionar(%q, %q) = true, want false", estado, EstadoLeadNuevo)
			}
		})
	}
}

func TestEstadosPosiblesNoExponeLaPolitica(t *testing.T) {
	primera := EstadosPosibles(EstadoLeadCalificado)
	primera[0] = EstadoLeadCerrado

	segunda := EstadosPosibles(EstadoLeadCalificado)
	esperada := []EstadoLead{
		EstadoLeadEntregado,
		EstadoLeadEnNutricion,
		EstadoLeadRemarketing,
		EstadoLeadDespedido,
	}
	if !estadosIguales(segunda, esperada) {
		t.Fatalf("second result = %v after mutating first, want %v", segunda, esperada)
	}
}

func estadosIguales(got, want []EstadoLead) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}
