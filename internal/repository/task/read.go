package task

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/scan/pgxscan"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
	bobgen "github.com/hedgehog-pm/service-tracker/internal/generated/bob/task"
	"github.com/hedgehog-pm/service-tracker/internal/storage/postgres"
)

func (r *Postgres) GetTask(ctx context.Context, projectID, taskID int) (domain.Task, error) {
	q := bobgen.GetTask(projectID, taskID)

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return domain.Task{}, postgres.MapError(ctx, fmt.Errorf("build task query: %w", err))
	}

	row, err := pgxscan.One(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Task{}, domain.ErrTaskNotFound
	}

	if err != nil {
		return domain.Task{}, postgres.MapError(ctx, fmt.Errorf("read task %d: %w", taskID, err))
	}

	return taskToDomain(row), nil
}

func (r *Postgres) ListTasks(ctx context.Context, projectID int, statusID, moduleID *int, limit, offset int) ([]domain.Task, error) {
	q := bobgen.ListTasks(projectID).With(
		sm.Limit(psql.Arg(limit)), sm.Offset(psql.Arg(offset)),
	)

	if statusID != nil {
		q.Apply(sm.Where(psql.Quote("t", "status_id").EQ(psql.Arg(*statusID))))
	}

	if moduleID != nil {
		q.Apply(sm.Where(psql.Quote("t", "module_id").EQ(psql.Arg(*moduleID))))
	}

	query, args, err := q.Build(ctx)
	if err != nil {
		return nil, postgres.MapError(ctx, fmt.Errorf("build tasks query: %w", err))
	}

	rows, err := pgxscan.All(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if err != nil {
		return nil, postgres.MapError(ctx, fmt.Errorf("list tasks: %w", err))
	}

	tasks := make([]domain.Task, 0, len(rows))
	for _, row := range rows {
		tasks = append(tasks, taskToDomain(row))
	}

	return tasks, nil
}
