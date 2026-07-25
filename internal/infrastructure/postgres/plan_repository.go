package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
)

type PlanRepository struct{ pool pgxPool }

func NuevoPlanRepository(pool pgxPool) *PlanRepository { return &PlanRepository{pool: pool} }

var _ usecase.PlanRepository = (*PlanRepository)(nil)

const planColumns = "plan_id,lead_id,estado,frecuencia,consentimiento_en,meta_monto,meta_descripcion"
const planInsertSQL = `INSERT INTO planes (` + planColumns + `) VALUES ($1,$2,$3,$4,$5,$6,$7)`
const planUpsertSQL = `INSERT INTO planes (` + planColumns + `) VALUES ($1,$2,$3,$4,$5,$6,$7) ON CONFLICT (plan_id) DO UPDATE SET lead_id=EXCLUDED.lead_id,estado=EXCLUDED.estado,frecuencia=EXCLUDED.frecuencia,consentimiento_en=EXCLUDED.consentimiento_en,meta_monto=EXCLUDED.meta_monto,meta_descripcion=EXCLUDED.meta_descripcion`
const hitoUpsertSQL = `INSERT INTO hitos (hito_id,plan_id,tipo,fecha,monto,descripcion,estado) VALUES ($1,$2,$3,$4,$5,$6,$7) ON CONFLICT (hito_id) DO UPDATE SET plan_id=EXCLUDED.plan_id,tipo=EXCLUDED.tipo,fecha=EXCLUDED.fecha,monto=EXCLUDED.monto,descripcion=EXCLUDED.descripcion,estado=EXCLUDED.estado`
const planHitosQuery = `SELECT hito_id,tipo,fecha::text,monto,descripcion,estado FROM hitos WHERE plan_id=$1 ORDER BY fecha ASC,hito_id ASC`
const overdueHitosQuery = `SELECT h.hito_id,h.tipo,h.fecha::text,h.monto,h.descripcion,h.estado,h.plan_id,p.lead_id FROM hitos h JOIN planes p ON p.plan_id=h.plan_id WHERE p.estado=$1 AND h.estado=$2 AND h.fecha <= $3::date ORDER BY h.fecha ASC,h.hito_id ASC`

func (r *PlanRepository) Crear(ctx context.Context, plan *domain.PlanNutricion) error {
	return r.persist(ctx, plan, planInsertSQL)
}

func (r *PlanRepository) Guardar(ctx context.Context, plan *domain.PlanNutricion) error {
	return r.persist(ctx, plan, planUpsertSQL)
}

func (r *PlanRepository) persist(ctx context.Context, plan *domain.PlanNutricion, planSQL string) error {
	if plan == nil {
		return fmt.Errorf("plan nil")
	}
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return repositoryError("plan", plan.PlanID, err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if _, err = tx.Exec(ctx, planSQL, plan.PlanID, plan.LeadID, plan.Estado, plan.Frecuencia, plan.ConsentimientoEn, plan.MetaMonto, plan.MetaDescripcion); err != nil {
		return repositoryError("plan", plan.PlanID, err)
	}
	for _, hito := range plan.Hitos {
		if _, err = tx.Exec(ctx, hitoUpsertSQL, hito.HitoID, plan.PlanID, hito.Tipo, hito.Fecha, hito.Monto, hito.Descripcion, hito.Estado); err != nil {
			return repositoryError("hito", hito.HitoID, err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return repositoryError("plan", plan.PlanID, err)
	}
	return nil
}

func (r *PlanRepository) PorLead(ctx context.Context, leadID string) (*domain.PlanNutricion, error) {
	var plan domain.PlanNutricion
	var consent *time.Time
	err := r.pool.QueryRow(ctx, `SELECT `+planColumns+` FROM planes WHERE lead_id=$1`, leadID).Scan(&plan.PlanID, &plan.LeadID, &plan.Estado, &plan.Frecuencia, &consent, &plan.MetaMonto, &plan.MetaDescripcion)
	if err != nil {
		return nil, repositoryError("plan", leadID, err)
	}
	plan.ConsentimientoEn = consent
	rows, err := r.pool.Query(ctx, planHitosQuery, plan.PlanID)
	if err != nil {
		return nil, repositoryError("plan", plan.PlanID, err)
	}
	defer rows.Close()
	plan.Hitos = make([]domain.Hito, 0)
	for rows.Next() {
		hito, scanErr := scanHito(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		plan.Hitos = append(plan.Hitos, hito)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *PlanRepository) HitosVencidos(ctx context.Context, at time.Time) ([]usecase.HitoConPlan, error) {
	rows, err := r.pool.Query(ctx, overdueHitosQuery, domain.EstadoPlanActivo, domain.EstadoHitoPendiente, at)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]usecase.HitoConPlan, 0)
	for rows.Next() {
		var item usecase.HitoConPlan
		if err := rows.Scan(&item.Hito.HitoID, &item.Hito.Tipo, &item.Hito.Fecha, &item.Hito.Monto, &item.Hito.Descripcion, &item.Hito.Estado, &item.PlanID, &item.LeadID); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func (r *PlanRepository) MarcarHito(ctx context.Context, hitoID string, state domain.EstadoHito) error {
	tag, err := r.pool.Exec(ctx, `UPDATE hitos SET estado=$2 WHERE hito_id=$1`, hitoID, state)
	if err != nil {
		return repositoryError("hito", hitoID, err)
	}
	if tag.RowsAffected() == 0 {
		return &usecase.NotFoundError{Resource: "hito", ID: hitoID}
	}
	return nil
}

func scanHito(s scanner) (domain.Hito, error) {
	var h domain.Hito
	if err := s.Scan(&h.HitoID, &h.Tipo, &h.Fecha, &h.Monto, &h.Descripcion, &h.Estado); err != nil {
		return domain.Hito{}, err
	}
	return h, nil
}
