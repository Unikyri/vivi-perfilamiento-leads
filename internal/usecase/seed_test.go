package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
)

type seedFake struct {
	calls int
	leads []domain.Lead
	block bool
}

func (f *seedFake) Sembrar(ctx context.Context, leads []domain.Lead) error {
	f.calls++
	if f.block {
		<-ctx.Done()
		return ctx.Err()
	}
	f.leads = append(f.leads, leads...)
	return nil
}

func TestCargarSeedIsDefaultOffAndCanonical(t *testing.T) {
	fake := &seedFake{}
	if err := (&CargarSeed{Repository: fake}).Ejecutar(context.Background()); err != nil {
		t.Fatal(err)
	}
	if fake.calls != 0 {
		t.Fatalf("disabled seed calls=%d", fake.calls)
	}
	if err := (&CargarSeed{Repository: fake, Habilitado: true}).Ejecutar(context.Background()); err != nil {
		t.Fatal(err)
	}
	if fake.calls != 1 || len(fake.leads) != 3 {
		t.Fatalf("calls=%d leads=%d", fake.calls, len(fake.leads))
	}
	for i, want := range []string{"ana", "carlos", "luisa"} {
		if fake.leads[i].LeadID != want || fake.leads[i].Fuente != "DEMO" {
			t.Fatalf("seed[%d]=%+v", i, fake.leads[i])
		}
	}
}

func TestCargarSeedHonorsContextTimeout(t *testing.T) {
	fake := &seedFake{block: true}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	started := time.Now()
	err := (&CargarSeed{Repository: fake, Habilitado: true}).Ejecutar(ctx)
	if !errors.Is(err, context.DeadlineExceeded) || time.Since(started) >= time.Second {
		t.Fatalf("err=%v elapsed=%s", err, time.Since(started))
	}
}
