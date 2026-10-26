package api

import (
	"net/http"

	"github.com/labstack/echo/v5"

	generated "github.com/hedgehog-pm/service-tracker/internal/generated/api/transitions"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/listtasktransitions"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/transitiontask"
)

type TransitionHandler struct {
	list       *listtasktransitions.Usecase
	transition *transitiontask.Usecase
}

func NewTransitionHandler(list *listtasktransitions.Usecase, transition *transitiontask.Usecase) *TransitionHandler {
	return &TransitionHandler{list: list, transition: transition}
}

func (h *TransitionHandler) ListTaskTransitions(ctx *echo.Context, projectID, taskID int) error {
	values, err := h.list.Do(ctx.Request().Context(), projectID, taskID)
	if err != nil {
		return err
	}

	items := make([]generated.Transition, 0, len(values))
	for _, value := range values {
		items = append(items, generated.Transition{ID: value.ID, Name: value.Name, To: statusToAPI(value.ToStatus)})
	}

	return ctx.JSON(http.StatusOK, generated.TransitionList{Items: items})
}

func (h *TransitionHandler) TransitionTask(ctx *echo.Context, projectID, taskID int) error {
	var body generated.TransitionTask
	if err := ctx.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	value, err := h.transition.Do(ctx.Request().Context(), transitiontask.Params{ProjectID: projectID, TaskID: taskID, TransitionID: body.TransitionID})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, taskToAPI(value))
}
