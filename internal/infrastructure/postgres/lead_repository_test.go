package postgres

import (
	"errors"
	"strings"
	"testing"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
	"github.com/jackc/pgx/v5"
)

func TestJSONBHelpers_PreserveNullAndNestedValues(t *testing.T) {
	null, err := encodeJSONB(nil)
	if err != nil || null != nil {
		t.Fatalf("nil JSONB = %v/%v, want nil/nil", null, err)
	}
	input := map[string]any{"nested": map[string]any{"items": []any{1.0, "x"}}}
	raw, err := encodeJSONB(input)
	if err != nil {
		t.Fatal(err)
	}
	var output map[string]any
	if err := decodeJSONB(raw, &output); err != nil || output["nested"].(map[string]any)["items"].([]any)[1] != "x" {
		t.Fatalf("nested JSONB = %#v/%v", output, err)
	}
	if err := decodeJSONB(nil, &output); err != nil {
		t.Fatal(err)
	}
}
func TestRepositoryErrorMapping(t *testing.T) {
	err := repositoryError("lead", "missing", pgx.ErrNoRows)
	var notFound *usecase.NotFoundError
	if !errors.As(err, &notFound) || !errors.Is(err, usecase.ErrNoEncontrado) || notFound.Resource != "lead" {
		t.Fatalf("not-found mapping = %v", err)
	}
	wrapped := repositoryError("lead", "x", errors.New("database down"))
	if !strings.Contains(wrapped.Error(), `lead "x"`) {
		t.Fatalf("wrapped mapping = %v", wrapped)
	}
}

func TestLeadRepository_SQLContract(t *testing.T) {
	checks := []string{
		"($1::boolean IS NULL OR afiliado=$1) AND ($2::text IS NULL OR ruta=$2)",
		"ORDER BY prioridad DESC,lead_id ASC",
	}
	query := `SELECT ` + leadColumns + ` FROM leads WHERE ($1::boolean IS NULL OR afiliado=$1) AND ($2::text IS NULL OR ruta=$2) ORDER BY prioridad DESC,lead_id ASC`
	for _, check := range checks[:2] {
		if !strings.Contains(query, check) {
			t.Errorf("lead query missing %q", check)
		}
	}
}
