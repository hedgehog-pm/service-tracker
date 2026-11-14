package transitiontask

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type (
	taskRepository interface {
		GetTaskForUpdate(ctx context.Context, projectID, taskID int) (domain.Task, error)
		UpdateStatus(ctx context.Context, task domain.Task) (domain.Task, error)
	}
	transitionReader interface {
		GetTransition(ctx context.Context, transitionID int) (domain.Transition, error)
	}
	transactor interface {
		InTx(ctx context.Context, fn func(ctx context.Context) error) error
	}
)
