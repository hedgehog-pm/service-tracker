package project

import (
	"context"
	"fmt"
	"time"

	"github.com/samber/lo"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/scan/pgxscan"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
	bobgen "github.com/hedgehog-pm/service-tracker/internal/generated/bob/project"
	"github.com/hedgehog-pm/service-tracker/internal/storage/postgres"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/projectsummary"
)

func (r *Postgres) GetProjectSummary(
	ctx context.Context,
	projectID int,
	today time.Time,
) (projectsummary.Summary, error) {
	q := bobgen.GetProjectSummary(projectID, today)

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return projectsummary.Summary{}, postgres.MapError(
			ctx,
			fmt.Errorf("build project summary query: %w", err),
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
		return projectsummary.Summary{}, postgres.MapError(
			ctx,
			fmt.Errorf("get project summary: %w", err),
		)
	}

	if len(rows) == 0 {
		return projectsummary.Summary{}, domain.ErrProjectNotFound
	}

	totals := rows[0]

	summary := projectsummary.Summary{
		TasksTotal:            int64(totals.TasksTotal),
		TasksOpen:             int64(totals.TasksOpen),
		TasksDueSoon:          int64(totals.TasksDueSoon),
		TasksOverdue:          int64(totals.TasksOverdue),
		TotalTimeSpentSeconds: int64(totals.TotalTimeSpentSeconds),
		TotalEstimateSeconds:  totals.TotalEstimateSeconds,
	}

	summary.UpcomingTasks = lo.FilterMap(
		rows,
		func(row bobgen.GetProjectSummaryRow, _ int) (projectsummary.UpcomingTask, bool) {
			if !row.ID.Valid {
				return projectsummary.UpcomingTask{}, false
			}

			task := domain.Task{
				ProjectKey: row.ProjectKey,
				Number:     row.Number.V,
			}

			return projectsummary.UpcomingTask{
				ID:         row.ID.V,
				Key:        task.Key(),
				Title:      row.Title.V,
				StatusID:   row.StatusID.V,
				PriorityID: row.PriorityID.V,
				DueDate:    row.DueDate.V,
			}, true
		},
	)

	return summary, nil
}
