package api

import (
	"context"

	"github.com/samber/lo"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
	generated "github.com/hedgehog-pm/service-tracker/internal/generated/api/statuses"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/status"
)

type StatusHandler struct{ statuses *status.Usecase }

func NewStatusHandler(statuses *status.Usecase) *StatusHandler {
	return &StatusHandler{statuses: statuses}
}

func (h *StatusHandler) ListStatuses(
	ctx context.Context,
	request generated.ListStatusesRequestObject,
) (generated.ListStatusesResponseObject, error) {
	values, err := h.statuses.List(ctx, request.ProjectID)
	if err != nil {
		return nil, err
	}

	items := lo.Map(values, func(value domain.Status, _ int) generated.Status {
		return statusToAPI(value)
	})

	return generated.ListStatuses200JSONResponse(generated.StatusList{Items: items}), nil
}

func statusToAPI(value domain.Status) generated.Status {
	return generated.Status{
		ID:       value.ID,
		Name:     value.Name,
		Category: generated.StatusCategory(value.Category),
	}
}
