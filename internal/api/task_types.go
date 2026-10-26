package api

import (
	"context"

	generated "github.com/hedgehog-pm/service-tracker/internal/generated/api/tasktypes"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/tasktype"
)

type TaskTypeHandler struct{ taskTypes *tasktype.Usecase }

func NewTaskTypeHandler(taskTypes *tasktype.Usecase) *TaskTypeHandler {
	return &TaskTypeHandler{taskTypes: taskTypes}
}

func (h *TaskTypeHandler) ListTaskTypes(
	ctx context.Context,
	request generated.ListTaskTypesRequestObject,
) (generated.ListTaskTypesResponseObject, error) {
	values, err := h.taskTypes.List(ctx, request.ProjectID)
	if err != nil {
		return nil, err
	}

	items := make([]generated.TaskType, 0, len(values))
	for _, value := range values {
		items = append(items, generated.TaskType{
			ID:            value.ID,
			Name:          value.Name,
			Position:      value.Position,
			InitialStatus: statusToAPI(value.InitialStatus),
		})
	}

	return generated.ListTaskTypes200JSONResponse(generated.TaskTypeList{Items: items}), nil
}
