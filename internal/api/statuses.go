package api

import (
	"net/http"

	"github.com/labstack/echo/v5"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
	generated "github.com/hedgehog-pm/service-tracker/internal/generated/api/statuses"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/liststatuses"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/status"
)

type StatusHandler struct {
	statuses   *status.Usecase
	listStatus *liststatuses.Usecase
}

func NewStatusHandler(statuses *status.Usecase, listStatus *liststatuses.Usecase) *StatusHandler {
	return &StatusHandler{statuses: statuses, listStatus: listStatus}
}

func (h *StatusHandler) CreateStatus(ctx *echo.Context, projectID int) error {
	var body generated.CreateStatus
	if err := ctx.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	value, err := h.statuses.Create(ctx.Request().Context(), projectID, body.Name, string(body.Category))
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, statusToAPI(value))
}

func (h *StatusHandler) GetStatus(ctx *echo.Context, projectID int, id int) error {
	value, err := h.statuses.Get(ctx.Request().Context(), projectID, id)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, statusToAPI(value))
}

func (h *StatusHandler) UpdateStatus(ctx *echo.Context, projectID int, id int) error {
	var body generated.UpdateStatus
	if err := ctx.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	value, err := h.statuses.Update(ctx.Request().Context(), projectID, id, body.Name, string(body.Category))
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, statusToAPI(value))
}

func (h *StatusHandler) ListStatuses(ctx *echo.Context, projectID int, params generated.ListStatusesParams) error {
	limit, offset := pagination(params.Limit, params.Offset)

	values, err := h.listStatus.Do(ctx.Request().Context(), projectID, limit, offset)
	if err != nil {
		return err
	}

	items := make([]generated.Status, 0, len(values))
	for _, value := range values {
		items = append(items, statusToAPI(value))
	}

	return ctx.JSON(http.StatusOK, generated.StatusList{Items: items, Limit: limit, Offset: offset})
}

func (h *StatusHandler) DeleteStatus(ctx *echo.Context, projectID, id int) error {
	if err := h.statuses.Delete(ctx.Request().Context(), projectID, id); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}

func statusToAPI(value domain.Status) generated.Status {
	return generated.Status{ID: value.ID, ProjectID: value.ProjectID, Name: value.Name, Category: generated.StatusCategory(value.Category)}
}
