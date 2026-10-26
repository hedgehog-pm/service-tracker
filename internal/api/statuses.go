package api

import (
	"net/http"

	"github.com/labstack/echo/v5"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
	generated "github.com/hedgehog-pm/service-tracker/internal/generated/api/statuses"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/status"
)

type StatusHandler struct{ statuses *status.Usecase }

func NewStatusHandler(statuses *status.Usecase) *StatusHandler {
	return &StatusHandler{statuses: statuses}
}

func (h *StatusHandler) ListStatuses(ctx *echo.Context, projectID int) error {
	values, err := h.statuses.List(ctx.Request().Context(), projectID)
	if err != nil {
		return err
	}

	items := make([]generated.Status, 0, len(values))
	for _, value := range values {
		items = append(items, statusToAPI(value))
	}

	return ctx.JSON(http.StatusOK, generated.StatusList{Items: items})
}

func statusToAPI(value domain.Status) generated.Status {
	return generated.Status{ID: value.ID, Name: value.Name, Category: generated.StatusCategory(value.Category)}
}
