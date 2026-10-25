package project

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/scan/pgxscan"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
	bobgen "github.com/hedgehog-pm/service-tracker/internal/generated/bob/project"
	"github.com/hedgehog-pm/service-tracker/internal/storage/postgres"
)

func (r *Postgres) ReserveTaskNumber(ctx context.Context, projectID int) (int, string, error) {
	q := bobgen.ReserveTaskNumber(projectID)

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return 0, "", postgres.MapError(ctx, fmt.Errorf("build task number query: %w", err))
	}

	row, err := pgxscan.One(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, "", domain.ErrProjectNotFound
	}

	if err != nil {
		return 0, "", postgres.MapError(ctx, fmt.Errorf("reserve task number in project %d: %w", projectID, err))
	}

	return row.Number, row.Key, nil
}
