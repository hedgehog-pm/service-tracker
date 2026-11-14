package listtasktransitions

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type (
	taskReader interface {
		GetTask(ctx context.Context, projectID, taskID int) (domain.Task, error)
	}
	transitionReader interface {
		ListTransitions(ctx context.Context, workflowID, statusID int) ([]domain.Transition, error)
	}
)
