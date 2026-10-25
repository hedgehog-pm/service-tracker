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

func (r *Postgres) InsertProject(ctx context.Context, value domain.Project) (domain.Project, error) {
	q := bobgen.InsertProject([]bobgen.InsertProject_Rows{{Key: value.Key, Name: value.Name, NextNumber: value.NextNumber}})

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return domain.Project{}, postgres.MapError(ctx, fmt.Errorf("build insert project query: %w", err))
	}

	row, err := pgxscan.One(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if err != nil {
		return domain.Project{}, postgres.MapError(ctx, fmt.Errorf("insert project: %w", err))
	}

	return projectToDomain(row), nil
}

func (r *Postgres) GetProject(ctx context.Context, id int) (domain.Project, error) {
	q := bobgen.GetProject(id)

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return domain.Project{}, postgres.MapError(ctx, fmt.Errorf("build get project query: %w", err))
	}

	row, err := pgxscan.One(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Project{}, domain.ErrProjectNotFound
	}

	if err != nil {
		return domain.Project{}, postgres.MapError(ctx, fmt.Errorf("get project: %w", err))
	}

	return projectToDomain(row), nil
}

func (r *Postgres) UpdateProject(ctx context.Context, value domain.Project) (domain.Project, error) {
	q := bobgen.UpdateProject(value.Name, value.ID)

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return domain.Project{}, postgres.MapError(ctx, fmt.Errorf("build update project query: %w", err))
	}

	row, err := pgxscan.One(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Project{}, domain.ErrProjectNotFound
	}

	if err != nil {
		return domain.Project{}, postgres.MapError(ctx, fmt.Errorf("update project: %w", err))
	}

	return projectToDomain(row), nil
}

func (r *Postgres) ListProjects(ctx context.Context, limit int, offset int) ([]domain.Project, error) {
	q := bobgen.ListProjects(offset, limit)

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return nil, postgres.MapError(ctx, fmt.Errorf("build projects query: %w", err))
	}

	rows, err := pgxscan.All(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if err != nil {
		return nil, postgres.MapError(ctx, fmt.Errorf("list projects: %w", err))
	}

	items := make([]domain.Project, 0, len(rows))
	for _, row := range rows {
		items = append(items, projectToDomain(row))
	}

	return items, nil
}

func (r *Postgres) DeleteProject(ctx context.Context, id int) error {
	query, args, err := bob.Build(ctx, bobgen.DeleteProject(id))
	if err != nil {
		return postgres.MapError(ctx, fmt.Errorf("build delete project query: %w", err))
	}

	result, err := postgres.Executor(ctx, r.db).Exec(ctx, query, args...)
	if err != nil {
		return postgres.MapDeleteError(ctx, fmt.Errorf("delete project: %w", err))
	}

	if result.RowsAffected() == 0 {
		return domain.ErrProjectNotFound
	}

	return nil
}
