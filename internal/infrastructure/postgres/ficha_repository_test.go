package postgres

import (
	"errors"
	"testing"
	"time"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
	"github.com/pashagolub/pgxmock/v4"
)

func TestFichaRepository_RuntimeBehavior(t *testing.T) {
	ctx := t.Context()
	ficha := &domain.Ficha{FichaID: "f1", LeadID: "l1", GeneradaEn: time.Date(2026, 7, 25, 12, 0, 0, 0, time.UTC), Identificacion: domain.Identificacion{Nombre: "Buyer"}}
	t.Run("upserts and retrieves replacement", func(t *testing.T) {
		mock, _ := pgxmock.NewPool()
		defer mock.Close()
		mock.ExpectExec("INSERT INTO fichas").WithArgs(ficha.FichaID, ficha.LeadID, pgxmock.AnyArg(), ficha.GeneradaEn).WillReturnResult(pgxmock.NewResult("INSERT", 1))
		if err := NuevoFichaRepository(mock).Guardar(ctx, ficha); err != nil {
			t.Fatal(err)
		}
		content, _ := encodeJSONB(ficha)
		mock.ExpectQuery("SELECT ficha_id,contenido,generada_en FROM fichas").WithArgs("l1").WillReturnRows(mock.NewRows([]string{"ficha_id", "contenido", "generada_en"}).AddRow("f1", content, ficha.GeneradaEn))
		got, err := NuevoFichaRepository(mock).PorLead(ctx, "l1")
		if err != nil || got.FichaID != "f1" || got.Identificacion.Nombre != "Buyer" {
			t.Fatalf("ficha = %#v/%v", got, err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})
	for _, tc := range []struct {
		name, id, resource string
		leadExists         bool
	}{{"unknown lead", "unknown", "lead", false}, {"lead without ficha", "no-ficha", "ficha", true}} {
		t.Run(tc.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()
			defer mock.Close()
			mock.ExpectQuery("SELECT ficha_id,contenido,generada_en FROM fichas").WithArgs(tc.id).WillReturnRows(mock.NewRows([]string{"ficha_id", "contenido", "generada_en"}))
			mock.ExpectQuery("SELECT EXISTS").WithArgs(tc.id).WillReturnRows(mock.NewRows([]string{"exists"}).AddRow(tc.leadExists))
			err := error(nil)
			_, err = NuevoFichaRepository(mock).PorLead(ctx, tc.id)
			var nf *usecase.NotFoundError
			if !errors.As(err, &nf) || nf.Resource != tc.resource {
				t.Fatalf("error = %v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
