package postgres

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

// MapError переводит известные ошибки PostgreSQL в доменные и сохраняет исходные доменные ошибки.
// Неизвестную ошибку логирует и заменяет на ErrInternalServerError.
func MapError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	if _, ok := errors.AsType[*domain.Error](err); ok {
		return err
	}

	if problem, ok := errors.AsType[*pgconn.PgError](err); ok {
		switch problem.Code {
		case "23505", "23P01":
			return domain.ErrAlreadyExists
		case "23503":
			return domain.ErrInvalidReference
		case "23502", "23514", "22P02", "22001":
			return domain.NewError(domain.CodeInvalidInput, "invalid input")
		case "22003", "2200H":
			return domain.NewError(domain.CodeConflict, "numeric range exhausted")
		case "40001", "40P01":
			return domain.NewError(
				domain.CodeConflict,
				"concurrent modification, retry the request",
			)
		}
	}

	slog.ErrorContext(ctx, "database operation failed", "error", err)

	return domain.ErrInternalServerError
}

// MapDeleteError трактует нарушение внешнего ключа при удалении как ErrResourceInUse.
// Остальные ошибки обрабатывает по правилам MapError.
func MapDeleteError(ctx context.Context, err error) error {
	mapped := MapError(ctx, err)
	if errors.Is(mapped, domain.ErrInvalidReference) {
		return domain.ErrResourceInUse
	}

	return mapped
}
