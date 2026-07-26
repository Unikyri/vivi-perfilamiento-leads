package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
	"github.com/jackc/pgx/v5"
)

type scanner interface{ Scan(...any) error }
type LeadRepository struct{ pool pgxPool }

func NuevoLeadRepository(pool pgxPool) *LeadRepository { return &LeadRepository{pool: pool} }

var _ usecase.LeadRepository = (*LeadRepository)(nil)

const leadColumns = `lead_id,nombre,telefono,cedula,fuente,estado,ruta,afiliado,prioridad,consume_cupo_10,perfil,capacidad,intencion,version,creado_en,actualizado_en`

func (r *LeadRepository) Crear(ctx context.Context, lead *domain.Lead) error {
	if lead == nil {
		return fmt.Errorf("lead nil")
	}
	if lead.Version == 0 {
		lead.Version = 1
	}
	perfil, err := encodeJSONB(lead.Perfil)
	if err != nil {
		return err
	}
	capacidad, err := encodeJSONB(lead.Capacidad)
	if err != nil {
		return err
	}
	intencion, err := encodeJSONB(lead.Intencion)
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, `INSERT INTO leads (`+leadColumns+`) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)`,
		lead.LeadID, lead.Nombre, lead.Telefono, lead.Cedula, lead.Fuente, lead.Estado, lead.Ruta, lead.Afiliado, lead.Prioridad, lead.ConsumeCupo10, perfil, capacidad, intencion, lead.Version, lead.CreadoEn, lead.ActualizadoEn)
	return repositoryError("lead", lead.LeadID, err)
}

func (r *LeadRepository) PorID(ctx context.Context, id string) (*domain.Lead, error) {
	row := r.pool.QueryRow(ctx, `SELECT `+leadColumns+` FROM leads WHERE lead_id=$1`, id)
	lead, err := scanLead(row)
	if err != nil {
		return nil, repositoryError("lead", id, err)
	}
	return lead, nil
}

func (r *LeadRepository) Guardar(ctx context.Context, lead *domain.Lead) error {
	if lead == nil {
		return fmt.Errorf("lead nil")
	}
	perfil, err := encodeJSONB(lead.Perfil)
	if err != nil {
		return err
	}
	capacidad, err := encodeJSONB(lead.Capacidad)
	if err != nil {
		return err
	}
	intencion, err := encodeJSONB(lead.Intencion)
	if err != nil {
		return err
	}
	var version int
	err = r.pool.QueryRow(ctx, `UPDATE leads SET nombre=$2,telefono=$3,cedula=$4,fuente=$5,estado=$6,ruta=$7,afiliado=$8,prioridad=$9,consume_cupo_10=$10,perfil=$11,capacidad=$12,intencion=$13,actualizado_en=$14,version=version+1 WHERE lead_id=$1 AND version=$15 RETURNING version`,
		lead.LeadID, lead.Nombre, lead.Telefono, lead.Cedula, lead.Fuente, lead.Estado, lead.Ruta, lead.Afiliado, lead.Prioridad, lead.ConsumeCupo10, perfil, capacidad, intencion, lead.ActualizadoEn, lead.Version).Scan(&version)
	if err == nil {
		lead.Version = version
		return nil
	}
	if !isNoRows(err) {
		return repositoryError("lead", lead.LeadID, err)
	}
	var exists bool
	if err = r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM leads WHERE lead_id=$1)`, lead.LeadID).Scan(&exists); err != nil {
		return repositoryError("lead", lead.LeadID, err)
	}
	if !exists {
		return &usecase.NotFoundError{Resource: "lead", ID: lead.LeadID}
	}
	return fmt.Errorf("lead %q: %w", lead.LeadID, usecase.ErrOptimisticLock)
}

func (r *LeadRepository) Listar(ctx context.Context, filter usecase.FiltroLeads) ([]*domain.Lead, error) {
	var affiliate, route any
	if filter.Afiliado != nil {
		affiliate = *filter.Afiliado
	}
	if filter.Ruta != nil {
		route = *filter.Ruta
	}
	rows, err := r.pool.Query(ctx, `SELECT `+leadColumns+` FROM leads WHERE ($1::boolean IS NULL OR afiliado=$1) AND ($2::text IS NULL OR ruta=$2) ORDER BY prioridad DESC,lead_id ASC`, affiliate, route)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]*domain.Lead, 0)
	for rows.Next() {
		lead, scanErr := scanLead(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		result = append(result, lead)
	}
	return result, rows.Err()
}

func (r *LeadRepository) AgregarMensaje(ctx context.Context, message *domain.Mensaje) error {
	if message == nil {
		return fmt.Errorf("mensaje nil")
	}
	exists, err := r.leadExists(ctx, message.LeadID)
	if err != nil {
		return repositoryError("lead", message.LeadID, err)
	}
	if !exists {
		return &usecase.NotFoundError{Resource: "lead", ID: message.LeadID}
	}
	attachment, err := encodeJSONB(message.Adjunto)
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, `INSERT INTO mensajes (mensaje_id,lead_id,autor,tipo_contenido,texto,adjunto,creado_en) VALUES ($1,$2,$3,$4,$5,$6,$7)`, message.MensajeID, message.LeadID, message.Autor, message.TipoContenido, message.Texto, attachment, message.CreadoEn)
	return repositoryError("mensaje", message.MensajeID, err)
}

func (r *LeadRepository) Conversacion(ctx context.Context, leadID string) ([]domain.Mensaje, error) {
	exists, err := r.leadExists(ctx, leadID)
	if err != nil {
		return nil, repositoryError("lead", leadID, err)
	}
	if !exists {
		return nil, &usecase.NotFoundError{Resource: "lead", ID: leadID}
	}
	rows, err := r.pool.Query(ctx, `SELECT mensaje_id,lead_id,autor,tipo_contenido,texto,adjunto,creado_en FROM mensajes WHERE lead_id=$1 ORDER BY creado_en ASC,mensaje_id ASC`, leadID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	messages := make([]domain.Mensaje, 0)
	for rows.Next() {
		var m domain.Mensaje
		var attachment []byte
		if err := rows.Scan(&m.MensajeID, &m.LeadID, &m.Autor, &m.TipoContenido, &m.Texto, &attachment, &m.CreadoEn); err != nil {
			return nil, err
		}
		if err := decodeJSONB(attachment, &m.Adjunto); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, rows.Err()
}

func (r *LeadRepository) leadExists(ctx context.Context, id string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM leads WHERE lead_id=$1)`, id).Scan(&exists)
	return exists, err
}

func isNoRows(err error) bool { return err == pgx.ErrNoRows }

func scanLead(s scanner) (*domain.Lead, error) {
	var l domain.Lead
	var perfil, capacidad, intencion []byte
	var created, updated time.Time
	err := s.Scan(&l.LeadID, &l.Nombre, &l.Telefono, &l.Cedula, &l.Fuente, &l.Estado, &l.Ruta, &l.Afiliado, &l.Prioridad, &l.ConsumeCupo10, &perfil, &capacidad, &intencion, &l.Version, &created, &updated)
	if err != nil {
		return nil, err
	}
	l.CreadoEn, l.ActualizadoEn = created, updated
	if err := decodeJSONB(perfil, &l.Perfil); err != nil {
		return nil, err
	}
	if len(capacidad) > 0 {
		l.Capacidad = new(domain.Capacidad)
		if err := decodeJSONB(capacidad, l.Capacidad); err != nil {
			return nil, err
		}
	}
	if len(intencion) > 0 {
		l.Intencion = new(domain.Intencion)
		if err := decodeJSONB(intencion, l.Intencion); err != nil {
			return nil, err
		}
	}
	return &l, nil
}
