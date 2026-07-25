package pipeline

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
)

// TopeVIS = 150 SMMLV con SMMLV 2026 = 1.750.905.
const TopeVIS int64 = 262635750

// EntradaMapa es una fila de data/mapa_proyectos.json.
type EntradaMapa struct {
	ProyectoID      string `json:"proyecto_id"`
	Nombre          string `json:"nombre"`
	NombreV2        string `json:"nombre_v2"`
	Zona            string `json:"zona"`
	BrochureURL     string `json:"brochure_url"`
	Recorrido360URL string `json:"recorrido_360_url"`
}

// CargarMapa lee el mapa escrito a mano.
func CargarMapa(ruta string) ([]EntradaMapa, error) {
	b, err := os.ReadFile(ruta)
	if err != nil {
		return nil, fmt.Errorf("leyendo mapa: %w", err)
	}
	var m []EntradaMapa
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, fmt.Errorf("mapa inválido: %w", err)
	}
	return m, nil
}

// ConstruirProyectos arma el catálogo cruzando el mapa con los compradores.
func ConstruirProyectos(mapa []EntradaMapa, compradores []domain.Comprador) ([]domain.Proyecto, error) {
	// index: nombre_v2 -> lista de precios
	precios := map[string][]int64{}
	for _, c := range compradores {
		if c.ValorCOP > 0 {
			precios[c.Proyecto] = append(precios[c.Proyecto], c.ValorCOP)
		}
	}

	var out []domain.Proyecto
	var sinCompradores []string
	for _, m := range mapa {
		ps := precios[m.NombreV2]
		if len(ps) == 0 {
			sinCompradores = append(sinCompradores, m.Nombre+" (busca: "+m.NombreV2+")")
			continue
		}
		min, max := ps[0], ps[0]
		for _, p := range ps {
			if p < min {
				min = p
			}
			if p > max {
				max = p
			}
		}
		out = append(out, domain.Proyecto{
			ProyectoID:      m.ProyectoID,
			Nombre:          m.Nombre,
			Zona:            m.Zona,
			PrecioDesde:     min,
			PrecioHasta:     max,
			EsVIS:           min <= TopeVIS,
			BrochureURL:     m.BrochureURL,
			Recorrido360URL: m.Recorrido360URL,
		})
	}
	if len(sinCompradores) > 0 {
		return nil, fmt.Errorf("proyectos del mapa sin compradores en v2 (revisar nombre_v2): %v", sinCompradores)
	}
	return out, nil
}

// AsignarProyectoID reescribe compradores[].proyecto_id usando el mapa,
// para que los IDs del kNN coincidan con los del catálogo.
// Los compradores de proyectos fuera del catálogo conservan su slug
// (cuentan para densidad, pero no generan tarjeta — doc 13 REC-3).
func AsignarProyectoID(mapa []EntradaMapa, cs []domain.Comprador) []domain.Comprador {
	porNombre := map[string]string{}
	for _, m := range mapa {
		porNombre[m.NombreV2] = m.ProyectoID
	}
	for i := range cs {
		if id, ok := porNombre[cs[i].Proyecto]; ok {
			cs[i].ProyectoID = id
		}
	}
	return cs
}
