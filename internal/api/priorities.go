package api

import (
	"context"

	"github.com/samber/lo"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
	generated "github.com/hedgehog-pm/service-tracker/internal/generated/api/priorities"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/priority"
)

type PriorityHandler struct{ priorities *priority.Usecase }

func NewPriorityHandler(priorities *priority.Usecase) *PriorityHandler {
	return &PriorityHandler{priorities: priorities}
}

func (h *PriorityHandler) ListPriorities(
	ctx context.Context,
	request generated.ListPrioritiesRequestObject,
) (generated.ListPrioritiesResponseObject, error) {
	result, err := h.priorities.List(ctx, request.ProjectID)
	if err != nil {
		return nil, err
	}

	items := lo.Map(result.Items, func(value domain.Priority, _ int) generated.Priority {
		return generated.Priority{
			ID:       value.ID,
			Name:     value.Name,
			Position: value.Position,
		}
	})

	return generated.ListPriorities200JSONResponse(generated.PriorityList{
		Items:             items,
		DefaultPriorityID: result.DefaultPriorityID,
	}), nil
}
