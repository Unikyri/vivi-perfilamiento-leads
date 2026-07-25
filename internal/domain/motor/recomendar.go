package motor

import (
	"fmt"
	"math/bits"
	"sort"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
)

const maxRecommendations = 3

type recommendationStats struct {
	neighbors int
	desistios int
}

type recommendationCandidate struct {
	card  domain.Recomendacion
	stats recommendationStats
}

// RecomendarProyectos returns up to three affordable catalog projects ranked by
// canonical neighbor evidence. Neighbor project IDs are treated as opaque
// canonical identifiers; display names never participate in aggregation.
func RecomendarProyectos(vecinos []Vecino, catalogo map[string]domain.Proyecto, presupuesto int64) []domain.Recomendacion {
	result := make([]domain.Recomendacion, 0, maxRecommendations)
	if presupuesto < 0 || len(vecinos) == 0 || len(catalogo) == 0 {
		return result
	}

	statsByID := make(map[string]recommendationStats)
	for _, vecino := range vecinos {
		projectID := vecino.ProyectoID
		if projectID == "" {
			continue
		}
		project, ok := catalogo[projectID]
		if !ok || project.ProyectoID != projectID {
			continue
		}

		stats := statsByID[projectID]
		stats.neighbors++
		if vecino.Desistio {
			stats.desistios++
		}
		statsByID[projectID] = stats
	}

	candidates := make([]recommendationCandidate, 0, len(statsByID))
	for projectID, stats := range statsByID {
		project := catalogo[projectID]
		if project.PrecioDesde <= 0 || presupuesto < project.PrecioDesde-project.PrecioDesde/5 {
			continue
		}

		rate := float64(stats.desistios) / float64(stats.neighbors)
		candidates = append(candidates, recommendationCandidate{
			card: domain.Recomendacion{
				ProyectoID:        project.ProyectoID,
				Nombre:            project.Nombre,
				Zona:              project.Zona,
				PrecioDesde:       project.PrecioDesde,
				Razon:             fmt.Sprintf("%d vecinos similares; tasa de desistimiento %.2f%%", stats.neighbors, rate*100),
				Vecinos:           stats.neighbors,
				TasaDesistimiento: rate,
				BrochureURL:       project.BrochureURL,
				Recorrido360URL:   project.Recorrido360URL,
			},
			stats: stats,
		})
	}

	sort.Slice(candidates, func(i, j int) bool {
		left, right := candidates[i], candidates[j]
		if left.stats.neighbors != right.stats.neighbors {
			return left.stats.neighbors > right.stats.neighbors
		}
		if fractionLessUint64(uint64(left.stats.desistios), uint64(left.stats.neighbors), uint64(right.stats.desistios), uint64(right.stats.neighbors)) {
			return true
		}
		if fractionLessUint64(uint64(right.stats.desistios), uint64(right.stats.neighbors), uint64(left.stats.desistios), uint64(left.stats.neighbors)) {
			return false
		}
		return left.card.ProyectoID < right.card.ProyectoID
	})

	if len(candidates) > maxRecommendations {
		candidates = candidates[:maxRecommendations]
	}
	for _, candidate := range candidates {
		result = append(result, candidate.card)
	}
	return result
}

// fractionLessUint64 compares leftNumerator/leftDenominator and
// rightNumerator/rightDenominator without converting the ordering operation to
// floating point or overflowing a 64-bit cross-product.
func fractionLessUint64(leftNumerator, leftDenominator, rightNumerator, rightDenominator uint64) bool {
	leftHigh, leftLow := bits.Mul64(leftNumerator, rightDenominator)
	rightHigh, rightLow := bits.Mul64(rightNumerator, leftDenominator)
	if leftHigh != rightHigh {
		return leftHigh < rightHigh
	}
	return leftLow < rightLow
}
