package api

import (
	"context"

	"github.com/samber/lo"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
	generated "github.com/hedgehog-pm/service-tracker/internal/generated/api/taskworklogs"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/taskworklog"
)

type TaskWorklogHandler struct {
	worklogs *taskworklog.Usecase
}

func NewTaskWorklogHandler(worklogs *taskworklog.Usecase) *TaskWorklogHandler {
	return &TaskWorklogHandler{worklogs: worklogs}
}

func (h *TaskWorklogHandler) GetTaskWorklog(ctx context.Context, request generated.GetTaskWorklogRequestObject) (generated.GetTaskWorklogResponseObject, error) {
	worklog, err := h.worklogs.Get(ctx, request.ProjectID, request.TaskID, request.WorklogID)
	if err != nil {
		return nil, err
	}

	return generated.GetTaskWorklog200JSONResponse(taskWorklogToAPI(worklog)), nil
}

func (h *TaskWorklogHandler) ListTaskWorklogs(ctx context.Context, request generated.ListTaskWorklogsRequestObject) (generated.ListTaskWorklogsResponseObject, error) {
	limit, offset := *request.Params.Limit, *request.Params.Offset

	worklogs, total, err := h.worklogs.List(ctx, request.ProjectID, request.TaskID, limit, offset)
	if err != nil {
		return nil, err
	}

	items := lo.Map(worklogs, func(worklog domain.TaskWorklog, _ int) generated.TaskWorklog {
		return taskWorklogToAPI(worklog)
	})

	return generated.ListTaskWorklogs200JSONResponse(generated.TaskWorklogList{
		Items:                 items,
		Limit:                 limit,
		Offset:                offset,
		TotalTimeSpentSeconds: total,
	}), nil
}

func (h *TaskWorklogHandler) CreateTaskWorklog(ctx context.Context, request generated.CreateTaskWorklogRequestObject) (generated.CreateTaskWorklogResponseObject, error) {
	body := request.Body

	worklog, err := h.worklogs.Create(ctx, request.ProjectID, request.TaskID, domain.TaskWorklogFields{
		StartedAt:        body.StartedAt,
		TimeSpentSeconds: body.TimeSpentSeconds,
		Comment:          body.Comment,
	})
	if err != nil {
		return nil, err
	}

	return generated.CreateTaskWorklog201JSONResponse(taskWorklogToAPI(worklog)), nil
}

func (h *TaskWorklogHandler) UpdateTaskWorklog(ctx context.Context, request generated.UpdateTaskWorklogRequestObject) (generated.UpdateTaskWorklogResponseObject, error) {
	body := request.Body

	worklog, err := h.worklogs.Update(ctx, request.ProjectID, request.TaskID, request.WorklogID, domain.TaskWorklogFields{
		StartedAt:        body.StartedAt,
		TimeSpentSeconds: body.TimeSpentSeconds,
		Comment:          body.Comment,
	})
	if err != nil {
		return nil, err
	}

	return generated.UpdateTaskWorklog200JSONResponse(taskWorklogToAPI(worklog)), nil
}

func (h *TaskWorklogHandler) DeleteTaskWorklog(ctx context.Context, request generated.DeleteTaskWorklogRequestObject) (generated.DeleteTaskWorklogResponseObject, error) {
	if err := h.worklogs.Delete(ctx, request.ProjectID, request.TaskID, request.WorklogID); err != nil {
		return nil, err
	}

	return generated.DeleteTaskWorklog204Response{}, nil
}

func taskWorklogToAPI(worklog domain.TaskWorklog) generated.TaskWorklog {
	return generated.TaskWorklog{
		ID:               worklog.ID,
		TaskID:           worklog.TaskID,
		StartedAt:        worklog.StartedAt,
		TimeSpentSeconds: worklog.TimeSpentSeconds,
		Comment:          worklog.Comment,
		CreatedAt:        worklog.CreatedAt,
		UpdatedAt:        worklog.UpdatedAt,
	}
}
