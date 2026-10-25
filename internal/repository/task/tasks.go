package task

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/scan/pgxscan"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
	bobgen "github.com/hedgehog-pm/service-tracker/internal/generated/bob/task"
	"github.com/hedgehog-pm/service-tracker/internal/storage/postgres"
)

func (r *Postgres) InsertTask(ctx context.Context, task domain.Task) (domain.Task, error) {
	q := bobgen.InsertTask(taskToInsertParams(task))

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return domain.Task{}, postgres.MapError(ctx, fmt.Errorf("build insert task query: %w", err))
	}

	row, err := pgxscan.One(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if err != nil {
		return domain.Task{}, postgres.MapError(ctx, fmt.Errorf("insert task: %w", err))
	}

	return insertedTaskToDomain(task, row), nil
}

func (r *Postgres) UpdateTask(ctx context.Context, task domain.Task) (domain.Task, error) {
	q := bobgen.UpdateTask(
		task.Title, task.Description, task.StatusID, pointerToNull(task.ModuleID), task.UpdatedAt, task.ProjectID, task.ID,
	)

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return domain.Task{}, postgres.MapError(ctx, fmt.Errorf("build update task query: %w", err))
	}

	row, err := pgxscan.One(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Task{}, domain.ErrTaskNotFound
	}

	if err != nil {
		return domain.Task{}, postgres.MapError(ctx, fmt.Errorf("update task %d: %w", task.ID, err))
	}

	return updatedTaskToDomain(task, row), nil
}

func (r *Postgres) DeleteTask(ctx context.Context, projectID, taskID int) error {
	query, args, err := bob.Build(ctx, bobgen.DeleteTask(projectID, taskID))
	if err != nil {
		return postgres.MapError(ctx, fmt.Errorf("build delete task query: %w", err))
	}

	result, err := postgres.Executor(ctx, r.db).Exec(ctx, query, args...)
	if err != nil {
		return postgres.MapDeleteError(ctx, fmt.Errorf("delete task %d: %w", taskID, err))
	}

	if result.RowsAffected() == 0 {
		return domain.ErrTaskNotFound
	}

	return nil
}
