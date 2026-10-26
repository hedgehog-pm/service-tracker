package workflow

import (
	"context"
	"fmt"

	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/scan/pgxscan"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
	bobgen "github.com/hedgehog-pm/service-tracker/internal/generated/bob/workflow"
	"github.com/hedgehog-pm/service-tracker/internal/storage/postgres"
)

// ListTaskTypes возвращает типы задач проекта с их начальными статусами в workflow.
func (r *Postgres) ListTaskTypes(
	ctx context.Context,
	schemeID int,
) ([]domain.TaskTypeWithInitialStatus, error) {
	q := bobgen.ListTaskTypes(schemeID)

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return nil, postgres.MapError(ctx, fmt.Errorf("build task types query: %w", err))
	}

	rows, err := pgxscan.All(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if err != nil {
		return nil, postgres.MapError(ctx, fmt.Errorf("list task types: %w", err))
	}

	items := make([]domain.TaskTypeWithInitialStatus, 0, len(rows))
	for _, row := range rows {
		items = append(items, taskTypeToDomain(row))
	}

	return items, nil
}
