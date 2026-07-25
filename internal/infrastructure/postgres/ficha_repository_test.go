package postgres

import (
	"strings"
	"testing"
)

func TestFichaRepository_SQLContract(t *testing.T) {
	if !strings.Contains(fichaUpsertSQL, "ON CONFLICT (lead_id) DO UPDATE") {
		t.Fatalf("ficha must be keyed by lead_id: %s", fichaUpsertSQL)
	}
	if strings.Contains(fichaUpsertSQL, "ON CONFLICT (ficha_id)") {
		t.Fatal("ficha upsert must not permit duplicate fichas per lead")
	}
	if !strings.Contains(fichaUpsertSQL, "contenido=EXCLUDED.contenido") {
		t.Fatal("ficha replacement must update content")
	}
}
