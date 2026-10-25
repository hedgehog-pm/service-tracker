package status

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/scan/pgxscan"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
	bobgen "github.com/hedgehog-pm/service-tracker/internal/generated/bob/status"
	"github.com/hedgehog-pm/service-tracker/internal/storage/postgres"
)

func (r *Postgres) InsertStatus(ctx context.Context, value domain.Status) (domain.Status, error) {
	q := bobgen.InsertStatus(value.Name, value.Category, value.ProjectID)

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return domain.Status{}, postgres.MapError(ctx, fmt.Errorf("build insert status query: %w", err))
	}

	row, err := pgxscan.One(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Status{}, domain.ErrProjectNotFound
	}

	if err != nil {
		return domain.Status{}, postgres.MapError(ctx, fmt.Errorf("insert status: %w", err))
	}

	return statusToDomain(row), nil
}

func (r *Postgres) GetStatus(ctx context.Context, projectID, id int) (domain.Status, error) {
	q := bobgen.GetStatus(projectID, id)

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return domain.Status{}, postgres.MapError(ctx, fmt.Errorf("build get status query: %w", err))
	}

	row, err := pgxscan.One(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Status{}, domain.ErrStatusNotFound
	}

	if err != nil {
		return domain.Status{}, postgres.MapError(ctx, fmt.Errorf("get status: %w", err))
	}

	return statusToDomain(row), nil
}

func (r *Postgres) UpdateStatus(ctx context.Context, value domain.Status) (domain.Status, error) {
	q := bobgen.UpdateStatus(value.Name, value.Category, value.ProjectID, value.ID)

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return domain.Status{}, postgres.MapError(ctx, fmt.Errorf("build update status query: %w", err))
	}

	row, err := pgxscan.One(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Status{}, domain.ErrStatusNotFound
	}

	if err != nil {
		return domain.Status{}, postgres.MapError(ctx, fmt.Errorf("update status: %w", err))
	}

	return statusToDomain(row), nil
}

func (r *Postgres) ListStatuses(ctx context.Context, projectID int, limit int, offset int) ([]domain.Status, error) {
	q := bobgen.ListStatuses(projectID, offset, limit)

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return nil, postgres.MapError(ctx, fmt.Errorf("build statuses query: %w", err))
	}

	rows, err := pgxscan.All(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if err != nil {
		return nil, postgres.MapError(ctx, fmt.Errorf("list statuses: %w", err))
	}

	items := make([]domain.Status, 0, len(rows))
	for _, row := range rows {
		items = append(items, statusToDomain(row))
	}

	return items, nil
}

func (r *Postgres) DeleteStatus(ctx context.Context, projectID, id int) error {
	query, args, err := bob.Build(ctx, bobgen.DeleteStatus(projectID, id))
	if err != nil {
		return postgres.MapError(ctx, fmt.Errorf("build delete status query: %w", err))
	}

	result, err := postgres.Executor(ctx, r.db).Exec(ctx, query, args...)
	if err != nil {
		return postgres.MapDeleteError(ctx, fmt.Errorf("delete status: %w", err))
	}

	if result.RowsAffected() == 0 {
		return domain.ErrStatusNotFound
	}

	return nil
}
