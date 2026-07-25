package postgres

import (
	"errors"
	"fmt"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func repositoryError(resource, id string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return &usecase.NotFoundError{Resource: resource, ID: id}
	}
	return fmt.Errorf("%s %q: %w", resource, id, err)
}

func foreignKeyError(resource, id string, err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23503" {
		return &usecase.NotFoundError{Resource: resource, ID: id}
	}
	return repositoryError(resource, id, err)
}
