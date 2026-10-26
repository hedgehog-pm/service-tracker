package api

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
	generated "github.com/hedgehog-pm/service-tracker/internal/generated/api/tasks"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/createtask"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/listtasks"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/task"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/updatetask"
)

type TaskHandler struct {
	tasks      *task.Usecase
	createTask *createtask.Usecase
	listTasks  *listtasks.Usecase
	updateTask *updatetask.Usecase
}

func NewTaskHandler(
	tasks *task.Usecase,
	createTask *createtask.Usecase,
	listTasks *listtasks.Usecase,
	updateTask *updatetask.Usecase,
) *TaskHandler {
	return &TaskHandler{
		tasks:      tasks,
		createTask: createTask,
		listTasks:  listTasks,
		updateTask: updateTask,
	}
}

func (h *TaskHandler) CreateTask(
	ctx context.Context,
	request generated.CreateTaskRequestObject,
) (generated.CreateTaskResponseObject, error) {
	body := request.Body

	params := createtask.Params{
		ProjectID:  request.ProjectID,
		Title:      body.Title,
		ModuleID:   body.ModuleID,
		TaskTypeID: body.TaskTypeID,
		PriorityID: body.PriorityID,
	}
	if body.Description != nil {
		params.Description = *body.Description
	}

	task, err := h.createTask.Do(ctx, params)
	if err != nil {
		return nil, err
	}

	return generated.CreateTask201JSONResponse(taskToAPI(task)), nil
}

func (h *TaskHandler) GetTask(
	ctx context.Context,
	request generated.GetTaskRequestObject,
) (generated.GetTaskResponseObject, error) {
	task, err := h.tasks.Get(ctx, request.ProjectID, request.TaskID)
	if err != nil {
		return nil, err
	}

	return generated.GetTask200JSONResponse(taskToAPI(task)), nil
}

func (h *TaskHandler) ListTasks(
	ctx context.Context,
	request generated.ListTasksRequestObject,
) (generated.ListTasksResponseObject, error) {
	limit, offset := *request.Params.Limit, *request.Params.Offset

	tasks, err := h.listTasks.Do(ctx, listtasks.Params{
		ProjectID:  request.ProjectID,
		StatusID:   request.Params.StatusID,
		ModuleID:   request.Params.ModuleID,
		TaskTypeID: request.Params.TaskTypeID,
		PriorityID: request.Params.PriorityID,
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		return nil, err
	}

	items := make([]generated.Task, 0, len(tasks))
	for _, task := range tasks {
		items = append(items, taskToAPI(task))
	}

	return generated.ListTasks200JSONResponse(generated.TaskList{
		Items:  items,
		Limit:  limit,
		Offset: offset,
	}), nil
}

func (h *TaskHandler) UpdateTask(
	ctx context.Context,
	request generated.UpdateTaskRequestObject,
) (generated.UpdateTaskResponseObject, error) {
	body := request.Body

	task, err := h.updateTask.Do(ctx, request.ProjectID, request.TaskID, domain.TaskFields{
		Title:       body.Title,
		Description: body.Description,
		PriorityID:  body.PriorityID,
		ModuleID:    body.ModuleID,
	})
	if err != nil {
		return nil, err
	}

	return generated.UpdateTask200JSONResponse(taskToAPI(task)), nil
}

func (h *TaskHandler) DeleteTask(
	ctx context.Context,
	request generated.DeleteTaskRequestObject,
) (generated.DeleteTaskResponseObject, error) {
	if err := h.tasks.Delete(ctx, request.ProjectID, request.TaskID); err != nil {
		return nil, err
	}

	return generated.DeleteTask204Response{}, nil
}
