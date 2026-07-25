package domain

import "fmt"

// ErrTransicionInvalida indicates that a lead cannot move between the requested states.
type ErrTransicionInvalida struct {
	Desde EstadoLead
	Hacia EstadoLead
}

func (e ErrTransicionInvalida) Error() string {
	return fmt.Sprintf("TRANSICION_INVALIDA: no se puede transicionar de %s a %s", e.Desde, e.Hacia)
}

// transiciones is the single domain-owned lifecycle policy. Its maps are initialized
// once and are never mutated after package initialization.
var transiciones = map[EstadoLead]map[EstadoLead]struct{}{
	EstadoLeadNuevo: {
		EstadoLeadPerfilando: {},
	},
	EstadoLeadPerfilando: {
		EstadoLeadCalificado: {},
	},
	EstadoLeadCalificado: {
		EstadoLeadEntregado:   {},
		EstadoLeadEnNutricion: {},
		EstadoLeadRemarketing: {},
		EstadoLeadDespedido:   {},
	},
	EstadoLeadEnNutricion: {
		EstadoLeadPausado:    {},
		EstadoLeadPerfilando: {},
	},
	EstadoLeadPausado: {
		EstadoLeadEnNutricion: {},
	},
	EstadoLeadRemarketing: {
		EstadoLeadPerfilando: {},
	},
	EstadoLeadEntregado: {
		EstadoLeadCerrado: {},
	},
}

// estadosOrdenados keeps query results deterministic while allowing each call to
// return a fresh slice. The order follows the Contract's lifecycle destinations.
var estadosOrdenados = []EstadoLead{
	EstadoLeadEntregado,
	EstadoLeadEnNutricion,
	EstadoLeadRemarketing,
	EstadoLeadDespedido,
	EstadoLeadPausado,
	EstadoLeadPerfilando,
	EstadoLeadCalificado,
	EstadoLeadCerrado,
	EstadoLeadNuevo,
}

// PuedeTransicionar reports whether the requested lifecycle edge is allowed.
func PuedeTransicionar(desde, hacia EstadoLead) bool {
	destinos, ok := transiciones[desde]
	if !ok {
		return false
	}
	_, ok = destinos[hacia]
	return ok
}

// EstadosPosibles returns the allowed destinations in Contract order.
func EstadosPosibles(desde EstadoLead) []EstadoLead {
	destinos := transiciones[desde]
	posibles := make([]EstadoLead, 0, len(destinos))
	for _, estado := range estadosOrdenados {
		if _, ok := destinos[estado]; ok {
			posibles = append(posibles, estado)
		}
	}
	return posibles
}

// Transicionar moves the lead only after validating the requested lifecycle edge.
func (l *Lead) Transicionar(hacia EstadoLead) error {
	desde := l.Estado
	if !PuedeTransicionar(desde, hacia) {
		return ErrTransicionInvalida{Desde: desde, Hacia: hacia}
	}
	l.Estado = hacia
	return nil
}
