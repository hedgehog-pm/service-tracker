package api

import (
	"context"

	generated "github.com/hedgehog-pm/service-tracker/internal/generated/api/transitions"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/listtasktransitions"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/transitiontask"
)

type TransitionHandler struct {
	list       *listtasktransitions.Usecase
	transition *transitiontask.Usecase
}

func NewTransitionHandler(
	list *listtasktransitions.Usecase,
	transition *transitiontask.Usecase,
) *TransitionHandler {
	return &TransitionHandler{
		list:       list,
		transition: transition,
	}
}

func (h *TransitionHandler) ListTaskTransitions(
	ctx context.Context,
	request generated.ListTaskTransitionsRequestObject,
) (generated.ListTaskTransitionsResponseObject, error) {
	values, err := h.list.Do(ctx, request.ProjectID, request.TaskID)
	if err != nil {
		return nil, err
	}

	items := make([]generated.Transition, 0, len(values))
	for _, value := range values {
		items = append(items, generated.Transition{
			ID:   value.ID,
			Name: value.Name,
			To:   statusToAPI(value.ToStatus),
		})
	}

	return generated.ListTaskTransitions200JSONResponse(generated.TransitionList{Items: items}), nil
}

func (h *TransitionHandler) TransitionTask(
	ctx context.Context,
	request generated.TransitionTaskRequestObject,
) (generated.TransitionTaskResponseObject, error) {
	value, err := h.transition.Do(ctx, transitiontask.Params{
		ProjectID:    request.ProjectID,
		TaskID:       request.TaskID,
		TransitionID: request.Body.TransitionID,
	})
	if err != nil {
		return nil, err
	}

	return generated.TransitionTask200JSONResponse(taskToAPI(value)), nil
}
