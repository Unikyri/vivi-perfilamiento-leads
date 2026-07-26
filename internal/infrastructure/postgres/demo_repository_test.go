package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v4"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
)

func TestDemoRepositoryReadsAndWritesRFC3339(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()
	at := time.Date(2026, 7, 28, 12, 30, 0, 0, time.UTC)
	mock.ExpectQuery("SELECT valor FROM demo").WithArgs(demoClockKey).WillReturnRows(mock.NewRows([]string{"valor"}).AddRow(at.Format(time.RFC3339Nano)))
	repo := NuevoDemoRepository(mock)
	got, err := repo.FechaSimulada(context.Background())
	if err != nil || !got.Equal(at) {
		t.Fatalf("date=%v err=%v", got, err)
	}
	mock.ExpectExec("INSERT INTO demo").WithArgs(demoClockKey, at.Format(time.RFC3339Nano)).WillReturnResult(pgxmock.NewResult("INSERT", 1))
	if err := repo.GuardarFechaSimulada(context.Background(), at); err != nil {
		t.Fatal(err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestDemoRepositoryResetDeletesOnlyApplicationRowsInOneTransaction(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()
	mock.ExpectBegin()
	for _, table := range []string{"fichas", "hitos", "planes", "mensajes", "leads"} {
		mock.ExpectExec("DELETE FROM " + table).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	}
	mock.ExpectExec("INSERT INTO demo").WithArgs(demoClockKey, approvedDemoDate).WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectCommit()
	got, err := NuevoDemoRepository(mock).Reiniciar(context.Background())
	want, _ := time.Parse(time.RFC3339, approvedDemoDate)
	if err != nil || !got.Equal(want) {
		t.Fatalf("date=%v err=%v", got, err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestDemoRepositorySeedsCanonicalLeadsInOneTransaction(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()
	mock.ExpectBegin()
	for _, lead := range usecase.SemillasDemo() {
		mock.ExpectExec("INSERT INTO leads").WithArgs(append([]interface{}{lead.LeadID}, anyArgs(15)...)...).WillReturnResult(pgxmock.NewResult("INSERT", 1))
	}
	mock.ExpectCommit()
	if err := NuevoDemoRepository(mock).Sembrar(context.Background(), usecase.SemillasDemo()); err != nil {
		t.Fatal(err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestDemoRepositoryResetWithSeedCommitsDateAndLeadsTogether(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()
	mock.ExpectBegin()
	for _, table := range []string{"fichas", "hitos", "planes", "mensajes", "leads"} {
		mock.ExpectExec("DELETE FROM " + table).WillReturnResult(pgxmock.NewResult("DELETE", 1))
	}
	mock.ExpectExec("INSERT INTO demo").WithArgs(demoClockKey, approvedDemoDate).WillReturnResult(pgxmock.NewResult("INSERT", 1))
	for _, lead := range usecase.SemillasDemo() {
		mock.ExpectExec("INSERT INTO leads").WithArgs(append([]interface{}{lead.LeadID}, anyArgs(15)...)...).WillReturnResult(pgxmock.NewResult("INSERT", 1))
	}
	mock.ExpectCommit()
	if _, err := NuevoDemoRepository(mock).ReiniciarConSeed(context.Background(), usecase.SemillasDemo()); err != nil {
		t.Fatal(err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
