package pgsql

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/qurk0/pr-service/internal/domain/models"
)

func mapErr(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, context.Canceled):
		return models.ErrCanceled

	case errors.Is(err, context.DeadlineExceeded):
		return models.ErrTimeout
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return models.ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return models.ErrInternal
	}

	return models.ErrInternal
}
