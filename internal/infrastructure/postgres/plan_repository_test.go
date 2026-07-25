package postgres

import (
	"errors"
	"testing"
	"time"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
	"github.com/pashagolub/pgxmock/v4"
)

func TestPlanRepository_RuntimeBehavior(t *testing.T) {
	ctx := t.Context()
	plan := &domain.PlanNutricion{PlanID: "p1", LeadID: "l1", Estado: domain.EstadoPlanActivo, Frecuencia: "MENSUAL", MetaMonto: 100}
	monto := int64(50)
	plan.Hitos = []domain.Hito{{HitoID: "h1", Tipo: domain.TipoHitoAhorro, Fecha: "2026-07-01", Monto: &monto, Descripcion: "save", Estado: domain.EstadoHitoPendiente}}
	t.Run("commits plan and supplied hito", func(t *testing.T) {
		mock, _ := pgxmock.NewPool()
		defer mock.Close()
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO planes").WithArgs(plan.PlanID, plan.LeadID, plan.Estado, plan.Frecuencia, plan.ConsentimientoEn, plan.MetaMonto, plan.MetaDescripcion).WillReturnResult(pgxmock.NewResult("INSERT", 1))
		mock.ExpectExec("INSERT INTO hitos").WithArgs("h1", "p1", domain.TipoHitoAhorro, "2026-07-01", &monto, "save", domain.EstadoHitoPendiente).WillReturnResult(pgxmock.NewResult("INSERT", 1))
		mock.ExpectCommit()
		if err := NuevoPlanRepository(mock).Guardar(ctx, plan); err != nil {
			t.Fatal(err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("rolls back when hito write fails", func(t *testing.T) {
		mock, _ := pgxmock.NewPool()
		defer mock.Close()
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO planes").WithArgs(anyArgs(7)...).WillReturnResult(pgxmock.NewResult("INSERT", 1))
		mock.ExpectExec("INSERT INTO hitos").WithArgs(anyArgs(7)...).WillReturnError(errors.New("constraint"))
		mock.ExpectRollback()
		if err := NuevoPlanRepository(mock).Guardar(ctx, plan); err == nil {
			t.Fatal("expected rollback error")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("reconstructs ordered aggregate and selects overdue active pending hitos", func(t *testing.T) {
		mock, _ := pgxmock.NewPool()
		defer mock.Close()
		at := time.Date(2026, 7, 25, 0, 0, 0, 0, time.UTC)
		mock.ExpectQuery("SELECT plan_id,lead_id.*FROM planes WHERE lead_id").WithArgs("l1").WillReturnRows(mock.NewRows([]string{"plan_id", "lead_id", "estado", "frecuencia", "consentimiento_en", "meta_monto", "meta_descripcion"}).AddRow("p1", "l1", domain.EstadoPlanActivo, "MENSUAL", nil, int64(100), "goal"))
		mock.ExpectQuery("SELECT hito_id,tipo,fecha::text.*FROM hitos WHERE plan_id").WithArgs("p1").WillReturnRows(mock.NewRows([]string{"hito_id", "tipo", "fecha", "monto", "descripcion", "estado"}).AddRow("h1", domain.TipoHitoAhorro, "2026-07-01", &monto, "save", domain.EstadoHitoPendiente).AddRow("h2", domain.TipoHitoPrima, "2026-08-01", nil, "later", domain.EstadoHitoPendiente))
		got, err := NuevoPlanRepository(mock).PorLead(ctx, "l1")
		if err != nil || len(got.Hitos) != 2 || got.Hitos[0].HitoID != "h1" {
			t.Fatalf("aggregate = %#v/%v", got, err)
		}
		mock.ExpectQuery("SELECT h.hito_id.*FROM hitos h JOIN planes").WithArgs(domain.EstadoPlanActivo, domain.EstadoHitoPendiente, at).WillReturnRows(mock.NewRows([]string{"hito_id", "tipo", "fecha", "monto", "descripcion", "estado", "plan_id", "lead_id"}).AddRow("h1", domain.TipoHitoAhorro, "2026-07-01", &monto, "save", domain.EstadoHitoPendiente, "p1", "l1"))
		due, err := NuevoPlanRepository(mock).HitosVencidos(ctx, at)
		if err != nil || len(due) != 1 || due[0].PlanID != "p1" {
			t.Fatalf("overdue = %#v/%v", due, err)
		}
		mock.ExpectExec("UPDATE hitos SET estado").WithArgs("missing", domain.EstadoHitoCumplido).WillReturnResult(pgxmock.NewResult("UPDATE", 0))
		err = NuevoPlanRepository(mock).MarcarHito(ctx, "missing", domain.EstadoHitoCumplido)
		var nf *usecase.NotFoundError
		if !errors.As(err, &nf) || nf.Resource != "hito" {
			t.Fatalf("missing hito = %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
	})
}
