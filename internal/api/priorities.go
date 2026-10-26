package api

import (
	"net/http"

	"github.com/labstack/echo/v5"

	generated "github.com/hedgehog-pm/service-tracker/internal/generated/api/priorities"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/priority"
)

type PriorityHandler struct{ priorities *priority.Usecase }

func NewPriorityHandler(priorities *priority.Usecase) *PriorityHandler {
	return &PriorityHandler{priorities: priorities}
}

func (h *PriorityHandler) ListPriorities(ctx *echo.Context, projectID int) error {
	result, err := h.priorities.List(ctx.Request().Context(), projectID)
	if err != nil {
		return err
	}

	items := make([]generated.Priority, 0, len(result.Items))
	for _, value := range result.Items {
		items = append(items, generated.Priority{ID: value.ID, Name: value.Name, Position: value.Position})
	}

	return ctx.JSON(http.StatusOK, generated.PriorityList{Items: items, DefaultPriorityID: result.DefaultPriorityID})
}
