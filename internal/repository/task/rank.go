package task

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/samber/lo"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/scan/pgxscan"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
	bobgen "github.com/hedgehog-pm/service-tracker/internal/generated/bob/task"
	"github.com/hedgehog-pm/service-tracker/internal/storage/postgres"
)

func (r *Postgres) GetLastRank(ctx context.Context, projectID int) (*int, error) {
	q := bobgen.GetLastRank(projectID)

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return nil, postgres.MapError(ctx, fmt.Errorf("build rank query: %w", err))
	}

	row, err := pgxscan.One(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, postgres.MapError(ctx, fmt.Errorf("read rank: %w", err))
	}

	return &row, nil
}

func (r *Postgres) GetPreviousRank(ctx context.Context, projectID, taskID, rank int) (*int, error) {
	q := bobgen.GetPreviousRank(projectID, taskID, rank)

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return nil, postgres.MapError(ctx, fmt.Errorf("build rank query: %w", err))
	}

	row, err := pgxscan.One(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, postgres.MapError(ctx, fmt.Errorf("read rank: %w", err))
	}

	return &row, nil
}

func (r *Postgres) GetNextRank(ctx context.Context, projectID, taskID, rank int) (*int, error) {
	q := bobgen.GetNextRank(projectID, taskID, rank)

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return nil, postgres.MapError(ctx, fmt.Errorf("build rank query: %w", err))
	}

	row, err := pgxscan.One(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, postgres.MapError(ctx, fmt.Errorf("read rank: %w", err))
	}

	return &row, nil
}

func (r *Postgres) GetLastColumnRank(
	ctx context.Context,
	projectID, boardID, columnID, taskID int,
) (*int, error) {
	q := bobgen.GetLastColumnRank(projectID, boardID, columnID, taskID)

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return nil, postgres.MapError(ctx, fmt.Errorf("build rank query: %w", err))
	}

	row, err := pgxscan.One(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, postgres.MapError(ctx, fmt.Errorf("read rank: %w", err))
	}

	return &row, nil
}

func (r *Postgres) ListTaskOrder(ctx context.Context, projectID int) ([]domain.Task, error) {
	q := bobgen.ListTaskOrder(projectID)

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return nil, postgres.MapError(ctx, fmt.Errorf("build task order query: %w", err))
	}

	rows, err := pgxscan.All(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if err != nil {
		return nil, postgres.MapError(ctx, fmt.Errorf("read task order: %w", err))
	}

	tasks := lo.Map(rows, func(row bobgen.ListTaskOrderRow, _ int) domain.Task {
		return taskToDomain(row)
	})

	return tasks, nil
}

func (r *Postgres) UpdateTaskRanks(ctx context.Context, tasks []domain.Task) error {
	for _, task := range tasks {
		query, args, err := bob.Build(ctx, bobgen.UpdateRank(task.Rank, task.ProjectID, task.ID))
		if err != nil {
			return postgres.MapError(ctx, fmt.Errorf("build update rank query: %w", err))
		}

		if _, err := postgres.Executor(ctx, r.db).Exec(ctx, query, args...); err != nil {
			return postgres.MapError(ctx, fmt.Errorf("update rank: %w", err))
		}
	}

	return nil
}

func (r *Postgres) UpdatePosition(ctx context.Context, task domain.Task) (domain.Task, error) {
	q := bobgen.UpdatePosition(task.Rank, task.StatusID, task.UpdatedAt, task.ProjectID, task.ID)

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return domain.Task{}, postgres.MapError(ctx, fmt.Errorf("build update position query: %w", err))
	}

	row, err := pgxscan.One(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Task{}, domain.ErrTaskNotFound
	}

	if err != nil {
		return domain.Task{}, postgres.MapError(ctx, fmt.Errorf("update position: %w", err))
	}

	return updatedTaskToDomain(task, row), nil
}
