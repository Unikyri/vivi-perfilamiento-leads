package postgres

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/infrastructure/ids"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
)

func TestJSONBHelpers_PreserveNullAndNestedValues(t *testing.T) {
	raw, err := encodeJSONB(map[string]any{"nested": map[string]any{"items": []any{1.0, "x"}}})
	if err != nil {
		t.Fatal(err)
	}
	var output map[string]any
	if err := decodeJSONB(raw, &output); err != nil || output["nested"].(map[string]any)["items"].([]any)[1] != "x" {
		t.Fatalf("nested JSONB = %#v/%v", output, err)
	}
	if raw, err := encodeJSONB(nil); err != nil || raw != nil {
		t.Fatalf("nil JSONB = %v/%v", raw, err)
	}
}

func TestRepositoryErrorMapping(t *testing.T) {
	err := repositoryError("lead", "missing", pgx.ErrNoRows)
	var notFound *usecase.NotFoundError
	if !errors.As(err, &notFound) || !errors.Is(err, usecase.ErrNoEncontrado) || notFound.Resource != "lead" {
		t.Fatalf("not-found mapping = %v", err)
	}
	if !strings.Contains(repositoryError("lead", "x", errors.New("database down")).Error(), `lead "x"`) {
		t.Fatal("database error lost repository identity")
	}
}

func anyArgs(n int) []interface{} {
	a := make([]interface{}, n)
	for i := range a {
		a[i] = pgxmock.AnyArg()
	}
	return a
}

func leadRows(mock pgxmock.PgxPoolIface, id, name string, priority float64) *pgxmock.Rows {
	now := time.Date(2026, 7, 25, 12, 0, 0, 0, time.UTC)
	return mock.NewRows([]string{"lead_id", "nombre", "telefono", "cedula", "fuente", "estado", "ruta", "afiliado", "prioridad", "consume_cupo_10", "perfil", "capacidad", "intencion", "version", "creado_en", "actualizado_en"}).AddRow(id, name, "300", "1", "web", domain.EstadoLeadNuevo, domain.RutaAsesor, true, priority, false, []byte(`{}`), nil, nil, 3, now, now)
}

func TestLeadRepository_RuntimeBehavior(t *testing.T) {
	ctx := t.Context()
	t.Run("generated opaque ID reaches create", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		if err != nil {
			t.Fatal(err)
		}
		defer mock.Close()
		id := ids.NuevoGeneradorID().Nuevo()
		lead := &domain.Lead{LeadID: id, Estado: domain.EstadoLeadNuevo, Ruta: domain.RutaAsesor, Perfil: domain.Perfil{}, CreadoEn: time.Now(), ActualizadoEn: time.Now()}
		args := append([]interface{}{id}, anyArgs(15)...)
		mock.ExpectExec("INSERT INTO leads").WithArgs(args...).WillReturnResult(pgxmock.NewResult("INSERT", 1))
		if err := NuevoLeadRepository(mock).Crear(ctx, lead); err != nil {
			t.Fatal(err)
		}
		mock.ExpectQuery("SELECT lead_id,nombre").WithArgs(id).WillReturnRows(leadRows(mock, id, "Generated", 0))
		stored, err := NuevoLeadRepository(mock).PorID(ctx, id)
		if err != nil || stored.LeadID != id {
			t.Fatalf("opaque ID read-back = %#v/%v, want %q", stored, err, id)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("CAS success advances caller version", func(t *testing.T) {
		mock, _ := pgxmock.NewPool()
		defer mock.Close()
		lead := &domain.Lead{LeadID: "cas", Version: 3, Perfil: domain.Perfil{}}
		mock.ExpectQuery("UPDATE leads SET").WithArgs(anyArgs(15)...).WillReturnRows(mock.NewRows([]string{"version"}).AddRow(4))
		if err := NuevoLeadRepository(mock).Guardar(ctx, lead); err != nil || lead.Version != 4 {
			t.Fatalf("save/version = %v/%d", err, lead.Version)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})
	for _, tc := range []struct {
		name, id       string
		exists         bool
		wantOptimistic bool
	}{{"missing", "missing", false, false}, {"stale", "stale", true, true}} {
		t.Run(tc.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			defer mock.Close()
			lead := &domain.Lead{LeadID: tc.id, Version: 3, Perfil: domain.Perfil{}}
			mock.ExpectQuery("UPDATE leads SET").WithArgs(anyArgs(15)...).WillReturnRows(mock.NewRows([]string{"version"}))
			mock.ExpectQuery("SELECT EXISTS").WithArgs(tc.id).WillReturnRows(mock.NewRows([]string{"exists"}).AddRow(tc.exists))
			err := NuevoLeadRepository(mock).Guardar(ctx, lead)
			if lead.Version != 3 {
				t.Fatalf("conflict mutated caller version: %d", lead.Version)
			}
			if tc.wantOptimistic && !errors.Is(err, usecase.ErrOptimisticLock) {
				t.Fatalf("error = %v", err)
			}
			if !tc.wantOptimistic {
				var nf *usecase.NotFoundError
				if !errors.As(err, &nf) || nf.Resource != "lead" {
					t.Fatalf("error = %v", err)
				}
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
	t.Run("conjunctive ordered list and chronological conversation", func(t *testing.T) {
		mock, _ := pgxmock.NewPool()
		defer mock.Close()
		affiliate, route := true, domain.RutaAsesor
		mock.ExpectQuery("SELECT .* FROM leads WHERE .*ORDER BY prioridad DESC,lead_id ASC").WithArgs(affiliate, route).WillReturnRows(leadRows(mock, "high", "High", 9).AddRow("low", "Low", "", "", "", domain.EstadoLeadNuevo, domain.RutaAsesor, true, 2.0, false, []byte(`{}`), nil, nil, 1, time.Now(), time.Now()))
		got, err := NuevoLeadRepository(mock).Listar(ctx, usecase.FiltroLeads{Afiliado: &affiliate, Ruta: &route})
		if err != nil || len(got) != 2 || got[0].LeadID != "high" || got[1].LeadID != "low" {
			t.Fatalf("list = %#v/%v", got, err)
		}
		mock.ExpectQuery("SELECT EXISTS").WithArgs("high").WillReturnRows(mock.NewRows([]string{"exists"}).AddRow(true))
		old, recent := time.Date(2026, 7, 24, 8, 0, 0, 0, time.UTC), time.Date(2026, 7, 25, 8, 0, 0, 0, time.UTC)
		mock.ExpectQuery("SELECT mensaje_id.*ORDER BY creado_en ASC,mensaje_id ASC").WithArgs("high").WillReturnRows(mock.NewRows([]string{"mensaje_id", "lead_id", "autor", "tipo_contenido", "texto", "adjunto", "creado_en"}).AddRow("m1", "high", domain.AutorMensajeLead, domain.TipoContenidoTexto, "old", []byte(`{}`), old).AddRow("m2", "high", domain.AutorMensajeVivi, domain.TipoContenidoTexto, "recent", []byte(`{}`), recent))
		messages, err := NuevoLeadRepository(mock).Conversacion(ctx, "high")
		if err != nil || len(messages) != 2 || messages[0].Texto != "old" || !messages[0].CreadoEn.Before(messages[1].CreadoEn) {
			t.Fatalf("conversation = %#v/%v", messages, err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})
}
