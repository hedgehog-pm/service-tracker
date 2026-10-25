package api

import (
	"net/http"

	"github.com/labstack/echo/v5"

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

func NewTaskHandler(tasks *task.Usecase, createTask *createtask.Usecase, listTasks *listtasks.Usecase, updateTask *updatetask.Usecase) *TaskHandler {
	return &TaskHandler{tasks: tasks, createTask: createTask, listTasks: listTasks, updateTask: updateTask}
}

func (h *TaskHandler) CreateTask(ctx *echo.Context, projectID int) error {
	var body generated.CreateTask
	if err := ctx.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	fields := domain.TaskFields{Title: body.Title, StatusID: body.StatusID, ModuleID: body.ModuleID}
	if body.Description != nil {
		fields.Description = *body.Description
	}

	task, err := h.createTask.Do(ctx.Request().Context(), projectID, fields)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, taskToAPI(task))
}

func (h *TaskHandler) GetTask(ctx *echo.Context, projectID, taskID int) error {
	task, err := h.tasks.Get(ctx.Request().Context(), projectID, taskID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, taskToAPI(task))
}

func (h *TaskHandler) ListTasks(ctx *echo.Context, projectID int, params generated.ListTasksParams) error {
	limit, offset := pagination(params.Limit, params.Offset)

	tasks, err := h.listTasks.Do(ctx.Request().Context(), listtasks.Params{
		ProjectID: projectID, StatusID: params.StatusID, ModuleID: params.ModuleID, Limit: limit, Offset: offset,
	})
	if err != nil {
		return err
	}

	items := make([]generated.Task, 0, len(tasks))
	for _, task := range tasks {
		items = append(items, taskToAPI(task))
	}

	return ctx.JSON(http.StatusOK, generated.TaskList{Items: items, Limit: limit, Offset: offset})
}

func (h *TaskHandler) UpdateTask(ctx *echo.Context, projectID, taskID int) error {
	var body generated.UpdateTask
	if err := ctx.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	task, err := h.updateTask.Do(ctx.Request().Context(), projectID, taskID, domain.TaskFields{
		Title: body.Title, Description: body.Description, StatusID: body.StatusID, ModuleID: body.ModuleID,
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, taskToAPI(task))
}

func (h *TaskHandler) DeleteTask(ctx *echo.Context, projectID, taskID int) error {
	if err := h.tasks.Delete(ctx.Request().Context(), projectID, taskID); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}
