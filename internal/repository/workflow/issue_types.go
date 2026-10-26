package workflow

import (
	"context"
	"fmt"

	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/scan/pgxscan"

	bobgen "github.com/hedgehog-pm/service-tracker/internal/generated/bob/workflow"
	"github.com/hedgehog-pm/service-tracker/internal/storage/postgres"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/issuetype"
)

func (r *Postgres) ListIssueTypes(ctx context.Context, schemeID int) ([]issuetype.Item, error) {
	q := bobgen.ListIssueTypes(schemeID)

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return nil, postgres.MapError(ctx, fmt.Errorf("build issue types query: %w", err))
	}

	rows, err := pgxscan.All(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if err != nil {
		return nil, postgres.MapError(ctx, fmt.Errorf("list issue types: %w", err))
	}

	items := make([]issuetype.Item, 0, len(rows))
	for _, row := range rows {
		items = append(items, issueTypeToItem(row))
	}

	return items, nil
}
