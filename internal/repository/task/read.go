package task

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/samber/lo"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/scan/pgxscan"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
	bobgen "github.com/hedgehog-pm/service-tracker/internal/generated/bob/task"
	"github.com/hedgehog-pm/service-tracker/internal/storage/postgres"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/listboardcolumntasks"
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

// GetTaskForUpdate блокирует строку задачи до завершения транзакции из ctx.
// Строка проекта, из которой читается ключ, не блокируется.
// Вызывай метод внутри InTx, чтобы блокировка сохранялась после возврата.
func (r *Postgres) GetTaskForUpdate(
	ctx context.Context,
	projectID, taskID int,
) (domain.Task, error) {
	q := bobgen.GetTaskForUpdate(projectID, taskID)

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

func (r *Postgres) ListTasks(
	ctx context.Context,
	projectID int,
	statusID, moduleID, taskTypeID, priorityID *int,
	limit, offset int,
) ([]domain.Task, error) {
	q := bobgen.ListTasks(projectID).With(
		sm.Limit(psql.Arg(limit)), sm.Offset(psql.Arg(offset)),
	)

	for _, predicate := range taskFilterPredicates(statusID, moduleID, taskTypeID, priorityID) {
		q.Apply(sm.Where(predicate))
	}

	query, args, err := q.Build(ctx)
	if err != nil {
		return nil, postgres.MapError(ctx, fmt.Errorf("build tasks query: %w", err))
	}

	rows, err := pgxscan.All(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if err != nil {
		return nil, postgres.MapError(ctx, fmt.Errorf("list tasks: %w", err))
	}

	tasks := lo.Map(rows, func(row bobgen.ListTasksRow, _ int) domain.Task {
		return taskToDomain(row)
	})

	return tasks, nil
}

func (r *Postgres) ListBoardColumnTasks(
	ctx context.Context,
	params listboardcolumntasks.Query,
) ([]domain.Task, error) {
	q := bobgen.ListBoardColumnTasks(
		params.ProjectID,
		params.BoardID,
		params.ColumnID,
		pointerToNull(params.AfterRank),
		pointerToNull(params.AfterNumber),
		params.Limit,
	).With()

	filters := params.Filters
	for _, predicate := range taskFilterPredicates(
		filters.StatusID,
		filters.ModuleID,
		filters.TaskTypeID,
		filters.PriorityID,
	) {
		q.Apply(sm.Where(predicate))
	}

	if filters.Query != "" {
		q.Apply(sm.Where(psql.Raw(
			`(strpos(lower(t.title), lower(?)) > 0
			OR strpos(lower(p.key || '-' || t.number::text), lower(?)) > 0)`,
			filters.Query,
			filters.Query,
		)))
	}

	query, args, err := q.Build(ctx)
	if err != nil {
		return nil, postgres.MapError(ctx, fmt.Errorf("build column tasks query: %w", err))
	}

	rows, err := pgxscan.All(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if err != nil {
		return nil, postgres.MapError(ctx, fmt.Errorf("read column tasks: %w", err))
	}

	tasks := lo.Map(rows, func(row bobgen.ListBoardColumnTasksRow, _ int) domain.Task {
		return taskToDomain(row)
	})

	return tasks, nil
}

func taskFilterPredicates(statusID, moduleID, taskTypeID, priorityID *int) []bob.Expression {
	var predicates []bob.Expression

	for _, filter := range [...]struct {
		column string
		id     *int
	}{
		{column: "status_id", id: statusID},
		{column: "module_id", id: moduleID},
		{column: "task_type_id", id: taskTypeID},
		{column: "priority_id", id: priorityID},
	} {
		if filter.id != nil {
			predicates = append(
				predicates,
				psql.Quote("t", filter.column).EQ(psql.Arg(*filter.id)),
			)
		}
	}

	return predicates
}
