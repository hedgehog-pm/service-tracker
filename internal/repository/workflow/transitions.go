package workflow

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/scan/pgxscan"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
	bobgen "github.com/hedgehog-pm/service-tracker/internal/generated/bob/workflow"
	"github.com/hedgehog-pm/service-tracker/internal/storage/postgres"
)

func (r *Postgres) GetTransition(ctx context.Context, transitionID int) (domain.Transition, error) {
	q := bobgen.GetTransition(transitionID)

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return domain.Transition{}, postgres.MapError(
			ctx,
			fmt.Errorf("build transition query: %w", err),
		)
	}

	row, err := pgxscan.One(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Transition{}, domain.ErrTransitionNotFound
	}

	if err != nil {
		return domain.Transition{}, postgres.MapError(ctx, fmt.Errorf("read transition: %w", err))
	}

	return transitionToDomain(row), nil
}

func (r *Postgres) ListTransitions(
	ctx context.Context,
	workflowID, statusID int,
) ([]domain.Transition, error) {
	q := bobgen.ListTransitions(workflowID, statusID)

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return nil, postgres.MapError(ctx, fmt.Errorf("build transitions query: %w", err))
	}

	rows, err := pgxscan.All(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if err != nil {
		return nil, postgres.MapError(ctx, fmt.Errorf("list transitions: %w", err))
	}

	items := make([]domain.Transition, 0, len(rows))
	for _, row := range rows {
		items = append(items, transitionToDomain(row))
	}

	return items, nil
}

// GetInitialStatus выбирает workflow по правилу для типа задачи в схеме
// и возвращает его начальный статус.
// Если правило не найдено, возвращает ErrIssueTypeUnavailable.
func (r *Postgres) GetInitialStatus(
	ctx context.Context,
	schemeID, issueTypeID int,
) (domain.InitialStatus, error) {
	q := bobgen.GetInitialStatus(schemeID, issueTypeID)

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return domain.InitialStatus{}, postgres.MapError(
			ctx,
			fmt.Errorf("build initial status query: %w", err),
		)
	}

	row, err := pgxscan.One(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.InitialStatus{}, domain.ErrIssueTypeUnavailable
	}

	if err != nil {
		return domain.InitialStatus{}, postgres.MapError(
			ctx,
			fmt.Errorf("read initial status: %w", err),
		)
	}

	return initialStatusToDomain(row), nil
}

func (r *Postgres) GetDefaultWorkflowScheme(ctx context.Context) (int, error) {
	q := bobgen.GetDefaultWorkflowScheme()

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return 0, postgres.MapError(ctx, fmt.Errorf("build default workflow scheme query: %w", err))
	}

	value, err := pgxscan.One(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if err != nil {
		return 0, postgres.MapError(ctx, fmt.Errorf("read default workflow scheme: %w", err))
	}

	return value, nil
}
