package taskworklog

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/scan/pgxscan"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
	bobgen "github.com/hedgehog-pm/service-tracker/internal/generated/bob/taskworklog"
	"github.com/hedgehog-pm/service-tracker/internal/storage/postgres"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/projectworklogs"
)

func (r *Postgres) ListProjectWorklogs(
	ctx context.Context,
	params projectworklogs.Query,
) (projectworklogs.Page[projectworklogs.Entry], error) {
	q := bobgen.ListProjectWorklogs(
		params.ProjectID,
		params.From,
		params.To,
		params.Search,
		params.Limit,
		params.Offset,
	)

	var page projectworklogs.Page[projectworklogs.Entry]

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return page, postgres.MapError(ctx, fmt.Errorf("build project worklogs query: %w", err))
	}

	rows, err := pgxscan.All(
		ctx,
		postgres.Executor(ctx, r.db),
		q.Scanner,
		query,
		args...,
	)
	if err != nil {
		return page, postgres.MapError(ctx, fmt.Errorf("list project worklogs: %w", err))
	}

	if len(rows) == 0 {
		return page, domain.ErrProjectNotFound
	}

	page.TotalTimeSpentSeconds = int64(rows[0].TotalTimeSpentSeconds)
	page.TotalEstimateSeconds = rows[0].TotalEstimateSeconds
	page.Items = lo.FilterMap(
		rows,
		func(row bobgen.ListProjectWorklogsRow, _ int) (projectworklogs.Entry, bool) {
			if !row.ID.Valid {
				return projectworklogs.Entry{}, false
			}

			return projectworklogs.Entry{
				ID:               row.ID.V,
				TaskID:           row.TaskID.V,
				TaskKey:          row.TaskKey.V,
				TaskTitle:        row.TaskTitle.V,
				StartedAt:        row.StartedAt.V,
				TimeSpentSeconds: row.TimeSpentSeconds.V,
				Comment:          commentToPointer(row.Comment),
			}, true
		},
	)

	return page, nil
}

func (r *Postgres) ListProjectWorklogsByTask(
	ctx context.Context,
	params projectworklogs.Query,
) (projectworklogs.Page[projectworklogs.TaskTotal], error) {
	q := bobgen.ListProjectWorklogsByTask(
		params.ProjectID,
		params.From,
		params.To,
		params.Search,
		params.Limit,
		params.Offset,
	)

	var page projectworklogs.Page[projectworklogs.TaskTotal]

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return page, postgres.MapError(
			ctx,
			fmt.Errorf("build project worklogs by task query: %w", err),
		)
	}

	rows, err := pgxscan.All(
		ctx,
		postgres.Executor(ctx, r.db),
		q.Scanner,
		query,
		args...,
	)
	if err != nil {
		return page, postgres.MapError(ctx, fmt.Errorf("list project worklogs by task: %w", err))
	}

	if len(rows) == 0 {
		return page, domain.ErrProjectNotFound
	}

	page.TotalTimeSpentSeconds = int64(rows[0].TotalTimeSpentSeconds)
	page.TotalEstimateSeconds = rows[0].TotalEstimateSeconds
	page.Items = lo.FilterMap(
		rows,
		func(row bobgen.ListProjectWorklogsByTaskRow, _ int) (projectworklogs.TaskTotal, bool) {
			if !row.TaskID.Valid {
				return projectworklogs.TaskTotal{}, false
			}

			return projectworklogs.TaskTotal{
				TaskID:           row.TaskID.V,
				TaskKey:          row.TaskKey.V,
				TaskTitle:        row.TaskTitle.V,
				WorklogsCount:    int64(row.WorklogsCount.V),
				TimeSpentSeconds: int64(row.TimeSpentSeconds.V),
			}, true
		},
	)

	return page, nil
}
