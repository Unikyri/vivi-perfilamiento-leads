package postgres

import (
	"context"
	"fmt"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
)

type FichaRepository struct{ pool pgxPool }

func NuevoFichaRepository(pool pgxPool) *FichaRepository { return &FichaRepository{pool: pool} }

var _ usecase.FichaRepository = (*FichaRepository)(nil)

const fichaUpsertSQL = `INSERT INTO fichas (ficha_id,lead_id,contenido,generada_en) VALUES ($1,$2,$3,$4) ON CONFLICT (lead_id) DO UPDATE SET ficha_id=EXCLUDED.ficha_id,contenido=EXCLUDED.contenido,generada_en=EXCLUDED.generada_en`

func (r *FichaRepository) Guardar(ctx context.Context, ficha *domain.Ficha) error {
	if ficha == nil {
		return fmt.Errorf("ficha nil")
	}
	content, err := encodeJSONB(ficha)
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, fichaUpsertSQL, ficha.FichaID, ficha.LeadID, content, ficha.GeneradaEn)
	return foreignKeyError("lead", ficha.LeadID, err)
}

func (r *FichaRepository) PorLead(ctx context.Context, leadID string) (*domain.Ficha, error) {
	var ficha domain.Ficha
	var content []byte
	err := r.pool.QueryRow(ctx, `SELECT ficha_id,contenido,generada_en FROM fichas WHERE lead_id=$1`, leadID).Scan(&ficha.FichaID, &content, &ficha.GeneradaEn)
	if err == nil {
		fichaID, generatedAt := ficha.FichaID, ficha.GeneradaEn
		if err := decodeJSONB(content, &ficha); err != nil {
			return nil, err
		}
		ficha.FichaID, ficha.LeadID, ficha.GeneradaEn = fichaID, leadID, generatedAt
		return &ficha, nil
	}
	if !isNoRows(err) {
		return nil, repositoryError("ficha", leadID, err)
	}
	var leadExists bool
	if err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM leads WHERE lead_id=$1)`, leadID).Scan(&leadExists); err != nil {
		return nil, repositoryError("lead", leadID, err)
	}
	if !leadExists {
		return nil, &usecase.NotFoundError{Resource: "lead", ID: leadID}
	}
	return nil, &usecase.NotFoundError{Resource: "ficha", ID: leadID}
}
