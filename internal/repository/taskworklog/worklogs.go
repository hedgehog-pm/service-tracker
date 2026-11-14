package taskworklog

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/samber/lo"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/scan/pgxscan"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
	bobgen "github.com/hedgehog-pm/service-tracker/internal/generated/bob/taskworklog"
	"github.com/hedgehog-pm/service-tracker/internal/storage/postgres"
)

func (r *Postgres) GetTaskWorklog(ctx context.Context, projectID, taskID, worklogID int) (domain.TaskWorklog, error) {
	q := bobgen.GetTaskWorklog(projectID, taskID, worklogID)

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return domain.TaskWorklog{}, postgres.MapError(ctx, fmt.Errorf("build get task worklog query: %w", err))
	}

	row, err := pgxscan.One(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.TaskWorklog{}, domain.ErrTaskWorklogNotFound
	}

	if err != nil {
		return domain.TaskWorklog{}, postgres.MapError(ctx, fmt.Errorf("get task worklog: %w", err))
	}

	return worklogToDomain(row), nil
}

func (r *Postgres) ListTaskWorklogs(ctx context.Context, projectID, taskID, limit, offset int) ([]domain.TaskWorklog, int64, error) {
	q := bobgen.ListTaskWorklogs(projectID, taskID, offset, limit)

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return nil, 0, postgres.MapError(ctx, fmt.Errorf("build task worklogs query: %w", err))
	}

	rows, err := pgxscan.All(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if err != nil {
		return nil, 0, postgres.MapError(ctx, fmt.Errorf("list task worklogs: %w", err))
	}

	if len(rows) == 0 {
		return nil, 0, domain.ErrTaskNotFound
	}

	// При пустой странице LEFT JOIN оставляет строку с суммой, но без ID записи.
	// Сохраняем её сумму, а саму строку исключаем из списка записей.
	worklogs := lo.FilterMap(rows, func(row bobgen.ListTaskWorklogsRow, _ int) (domain.TaskWorklog, bool) {
		if !row.ID.Valid {
			return domain.TaskWorklog{}, false
		}

		return domain.TaskWorklog{
			ID:               row.ID.V,
			TaskID:           row.TaskID.V,
			StartedAt:        row.StartedAt.V,
			TimeSpentSeconds: row.TimeSpentSeconds.V,
			Comment:          commentToPointer(row.Comment),
			CreatedAt:        row.CreatedAt.V,
			UpdatedAt:        row.UpdatedAt.V,
		}, true
	})

	return worklogs, int64(rows[0].TotalTimeSpentSeconds), nil
}

func (r *Postgres) InsertTaskWorklog(ctx context.Context, projectID int, worklog domain.TaskWorklog) (domain.TaskWorklog, error) {
	q := bobgen.InsertTaskWorklog(
		worklog.StartedAt,
		worklog.TimeSpentSeconds,
		commentToNull(worklog.Comment),
		worklog.CreatedAt,
		worklog.UpdatedAt,
		projectID,
		worklog.TaskID,
	)

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return domain.TaskWorklog{}, postgres.MapError(ctx, fmt.Errorf("build insert task worklog query: %w", err))
	}

	row, err := pgxscan.One(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.TaskWorklog{}, domain.ErrTaskNotFound
	}

	if err != nil {
		return domain.TaskWorklog{}, postgres.MapError(ctx, fmt.Errorf("insert task worklog: %w", err))
	}

	return worklogToDomain(row), nil
}

func (r *Postgres) UpdateTaskWorklog(ctx context.Context, projectID int, worklog domain.TaskWorklog) (domain.TaskWorklog, error) {
	q := bobgen.UpdateTaskWorklog(
		worklog.StartedAt,
		worklog.TimeSpentSeconds,
		commentToNull(worklog.Comment),
		worklog.UpdatedAt,
		worklog.TaskID,
		worklog.ID,
		projectID,
	)

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return domain.TaskWorklog{}, postgres.MapError(ctx, fmt.Errorf("build update task worklog query: %w", err))
	}

	row, err := pgxscan.One(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.TaskWorklog{}, domain.ErrTaskWorklogNotFound
	}

	if err != nil {
		return domain.TaskWorklog{}, postgres.MapError(ctx, fmt.Errorf("update task worklog: %w", err))
	}

	return worklogToDomain(row), nil
}

func (r *Postgres) DeleteTaskWorklog(ctx context.Context, projectID, taskID, worklogID int) error {
	query, args, err := bob.Build(ctx, bobgen.DeleteTaskWorklog(taskID, worklogID, projectID))
	if err != nil {
		return postgres.MapError(ctx, fmt.Errorf("build delete task worklog query: %w", err))
	}

	result, err := postgres.Executor(ctx, r.db).Exec(ctx, query, args...)
	if err != nil {
		return postgres.MapDeleteError(ctx, fmt.Errorf("delete task worklog: %w", err))
	}

	if result.RowsAffected() == 0 {
		return domain.ErrTaskWorklogNotFound
	}

	return nil
}
