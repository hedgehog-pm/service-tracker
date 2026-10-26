package priority

import (
	"context"
	"fmt"

	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/scan/pgxscan"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
	bobgen "github.com/hedgehog-pm/service-tracker/internal/generated/bob/priority"
	"github.com/hedgehog-pm/service-tracker/internal/storage/postgres"
)

func (r *Postgres) GetDefaultPriority(ctx context.Context, schemeID int) (int, error) {
	q := bobgen.GetDefaultPriority(schemeID)

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return 0, postgres.MapError(ctx, fmt.Errorf("build default priority query: %w", err))
	}

	row, err := pgxscan.One(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if err != nil {
		return 0, postgres.MapError(ctx, fmt.Errorf("read default priority: %w", err))
	}

	return row, nil
}

func (r *Postgres) ListPriorities(ctx context.Context, schemeID int) ([]domain.Priority, error) {
	q := bobgen.ListPriorities(schemeID)

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return nil, postgres.MapError(ctx, fmt.Errorf("build priorities query: %w", err))
	}

	rows, err := pgxscan.All(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if err != nil {
		return nil, postgres.MapError(ctx, fmt.Errorf("list priorities: %w", err))
	}

	items := make([]domain.Priority, 0, len(rows))
	for _, row := range rows {
		items = append(items, priorityToDomain(row))
	}

	return items, nil
}

func (r *Postgres) GetDefaultPriorityScheme(ctx context.Context) (int, error) {
	q := bobgen.GetDefaultPriorityScheme()

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return 0, postgres.MapError(ctx, fmt.Errorf("build default priority scheme query: %w", err))
	}

	value, err := pgxscan.One(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if err != nil {
		return 0, postgres.MapError(ctx, fmt.Errorf("read default priority scheme: %w", err))
	}

	return value, nil
}
