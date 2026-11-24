package api

import (
	"context"

	"github.com/samber/lo"

	generated "github.com/hedgehog-pm/service-tracker/internal/generated/api/taskworklogs"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/projectworklogs"
)

func (h *TaskWorklogHandler) ListProjectWorklogs(
	ctx context.Context,
	request generated.ListProjectWorklogsRequestObject,
) (generated.ListProjectWorklogsResponseObject, error) {
	params := request.Params
	limit, offset := *params.Limit, *params.Offset

	page, err := h.projectWorklogs.List(
		ctx,
		projectworklogs.Query{
			ProjectID: request.ProjectID,
			From:      params.From,
			To:        params.To,
			Search:    lo.FromPtr(params.Q),
			Limit:     limit,
			Offset:    offset,
		},
	)
	if err != nil {
		return nil, err
	}

	items := lo.Map(page.Items, func(entry projectworklogs.Entry, _ int) generated.ProjectWorklog {
		return generated.ProjectWorklog{
			ID:               entry.ID,
			TaskID:           entry.TaskID,
			TaskKey:          entry.TaskKey,
			TaskTitle:        entry.TaskTitle,
			StartedAt:        entry.StartedAt,
			TimeSpentSeconds: entry.TimeSpentSeconds,
			Comment:          entry.Comment,
		}
	})

	return generated.ListProjectWorklogs200JSONResponse(generated.ProjectWorklogList{
		Items:                 items,
		Limit:                 limit,
		Offset:                offset,
		TotalTimeSpentSeconds: page.TotalTimeSpentSeconds,
		TotalEstimateSeconds:  page.TotalEstimateSeconds,
	}), nil
}

func (h *TaskWorklogHandler) ListProjectWorklogsByTask(
	ctx context.Context,
	request generated.ListProjectWorklogsByTaskRequestObject,
) (generated.ListProjectWorklogsByTaskResponseObject, error) {
	params := request.Params
	limit, offset := *params.Limit, *params.Offset

	page, err := h.projectWorklogs.ListByTask(
		ctx,
		projectworklogs.Query{
			ProjectID: request.ProjectID,
			From:      params.From,
			To:        params.To,
			Search:    lo.FromPtr(params.Q),
			Limit:     limit,
			Offset:    offset,
		},
	)
	if err != nil {
		return nil, err
	}

	items := lo.Map(
		page.Items,
		func(task projectworklogs.TaskTotal, _ int) generated.ProjectWorklogByTask {
			return generated.ProjectWorklogByTask{
				TaskID:           task.TaskID,
				TaskKey:          task.TaskKey,
				TaskTitle:        task.TaskTitle,
				WorklogsCount:    task.WorklogsCount,
				TimeSpentSeconds: task.TimeSpentSeconds,
			}
		},
	)

	return generated.ListProjectWorklogsByTask200JSONResponse(generated.ProjectWorklogByTaskList{
		Items:                 items,
		Limit:                 limit,
		Offset:                offset,
		TotalTimeSpentSeconds: page.TotalTimeSpentSeconds,
		TotalEstimateSeconds:  page.TotalEstimateSeconds,
	}), nil
}
