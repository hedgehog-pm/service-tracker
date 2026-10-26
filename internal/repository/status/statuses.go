package status

import (
	"context"
	"fmt"

	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/scan/pgxscan"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
	bobgen "github.com/hedgehog-pm/service-tracker/internal/generated/bob/status"
	"github.com/hedgehog-pm/service-tracker/internal/storage/postgres"
)

func (r *Postgres) ListStatuses(ctx context.Context, schemeID int) ([]domain.Status, error) {
	q := bobgen.ListStatuses(schemeID)

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
