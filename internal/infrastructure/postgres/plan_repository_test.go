package postgres

import (
	"strings"
	"testing"
)

func TestPlanRepository_SQLContract(t *testing.T) {
	if strings.Contains(hitoUpsertSQL, "DELETE") || !strings.Contains(hitoUpsertSQL, "ON CONFLICT (hito_id) DO UPDATE") {
		t.Fatalf("hito save must be non-destructive upsert: %s", hitoUpsertSQL)
	}
	if !strings.Contains(planUpsertSQL, "ON CONFLICT (plan_id) DO UPDATE") {
		t.Fatalf("plan save must upsert without CAS: %s", planUpsertSQL)
	}
	if strings.Contains(planUpsertSQL, "version") || strings.Contains(hitoUpsertSQL, "version") {
		t.Fatal("plan persistence must not introduce version/CAS columns")
	}
}

func TestPlanRepository_QueryContracts(t *testing.T) {
	if !strings.Contains(planHitosQuery, "ORDER BY fecha ASC,hito_id ASC") {
		t.Fatalf("aggregate query must reconstruct deterministic hito order: %s", planHitosQuery)
	}
	for _, clause := range []string{"p.estado=$1", "h.estado=$2", "h.fecha <= $3::date", "ORDER BY h.fecha ASC,h.hito_id ASC"} {
		if !strings.Contains(overdueHitosQuery, clause) {
			t.Errorf("overdue query missing %q: %s", clause, overdueHitosQuery)
		}
	}
}
