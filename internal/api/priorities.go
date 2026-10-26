package api

import (
	"context"

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

	items := make([]generated.Priority, 0, len(result.Items))
	for _, value := range result.Items {
		items = append(items, generated.Priority{
			ID:       value.ID,
			Name:     value.Name,
			Position: value.Position,
		})
	}

	return generated.ListPriorities200JSONResponse(generated.PriorityList{
		Items:             items,
		DefaultPriorityID: result.DefaultPriorityID,
	}), nil
}
