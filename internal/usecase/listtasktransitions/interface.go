package listtasktransitions

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type (
	tasks interface {
		GetTask(ctx context.Context, projectID, taskID int) (domain.Task, error)
	}
	reader interface {
		ListTransitions(ctx context.Context, workflowID, statusID int) ([]domain.Transition, error)
	}
)
