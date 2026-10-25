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

func (r *Postgres) GetProjectKey(ctx context.Context, projectID int) (string, error) {
	q := bobgen.GetProjectKey(projectID)

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return "", postgres.MapError(ctx, fmt.Errorf("build project query: %w", err))
	}

	key, err := pgxscan.One(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return "", domain.ErrProjectNotFound
	}

	if err != nil {
		return "", postgres.MapError(ctx, fmt.Errorf("read project %d: %w", projectID, err))
	}

	return key, nil
}
