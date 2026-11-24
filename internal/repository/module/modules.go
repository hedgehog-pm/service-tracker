package module

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/samber/lo"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/scan/pgxscan"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
	bobgen "github.com/hedgehog-pm/service-tracker/internal/generated/bob/module"
	"github.com/hedgehog-pm/service-tracker/internal/storage/postgres"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/listmodules"
)

func (r *Postgres) InsertModule(ctx context.Context, value domain.Module) (domain.Module, error) {
	q := bobgen.InsertModule(value.Name, value.ProjectID)

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return domain.Module{}, postgres.MapError(
			ctx,
			fmt.Errorf("build insert module query: %w", err),
		)
	}

	row, err := pgxscan.One(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Module{}, domain.ErrProjectNotFound
	}

	if err != nil {
		return domain.Module{}, postgres.MapError(ctx, fmt.Errorf("insert module: %w", err))
	}

	return moduleToDomain(row), nil
}

func (r *Postgres) GetModule(ctx context.Context, projectID, id int) (domain.Module, error) {
	q := bobgen.GetModule(projectID, id)

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return domain.Module{}, postgres.MapError(
			ctx,
			fmt.Errorf("build get module query: %w", err),
		)
	}

	row, err := pgxscan.One(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Module{}, domain.ErrModuleNotFound
	}

	if err != nil {
		return domain.Module{}, postgres.MapError(ctx, fmt.Errorf("get module: %w", err))
	}

	return moduleToDomain(row), nil
}

func (r *Postgres) UpdateModule(ctx context.Context, value domain.Module) (domain.Module, error) {
	q := bobgen.UpdateModule(value.Name, value.ProjectID, value.ID)

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return domain.Module{}, postgres.MapError(
			ctx,
			fmt.Errorf("build update module query: %w", err),
		)
	}

	row, err := pgxscan.One(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Module{}, domain.ErrModuleNotFound
	}

	if err != nil {
		return domain.Module{}, postgres.MapError(ctx, fmt.Errorf("update module: %w", err))
	}

	return moduleToDomain(row), nil
}

func (r *Postgres) ListModules(
	ctx context.Context,
	projectID int,
	limit int,
	offset int,
) ([]listmodules.ModuleWithProgress, error) {
	q := bobgen.ListModules(projectID, offset, limit)

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return nil, postgres.MapError(ctx, fmt.Errorf("build modules query: %w", err))
	}

	rows, err := pgxscan.All(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if err != nil {
		return nil, postgres.MapError(ctx, fmt.Errorf("list modules: %w", err))
	}

	items := lo.Map(rows, func(row bobgen.ListModulesRow, _ int) listmodules.ModuleWithProgress {
		return listmodules.ModuleWithProgress{
			ID:         row.ID,
			ProjectID:  row.ProjectID,
			Name:       row.Name,
			TasksTotal: int64(row.TasksTotal),
			TasksDone:  int64(row.TasksDone),
		}
	})

	return items, nil
}

func (r *Postgres) DeleteModule(ctx context.Context, projectID, id int) error {
	query, args, err := bob.Build(ctx, bobgen.DeleteModule(projectID, id))
	if err != nil {
		return postgres.MapError(ctx, fmt.Errorf("build delete module query: %w", err))
	}

	result, err := postgres.Executor(ctx, r.db).Exec(ctx, query, args...)
	if err != nil {
		return postgres.MapDeleteError(ctx, fmt.Errorf("delete module: %w", err))
	}

	if result.RowsAffected() == 0 {
		return domain.ErrModuleNotFound
	}

	return nil
}
