package transitiontask

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type repository interface {
	GetTaskForUpdate(ctx context.Context, projectID, taskID int) (domain.Task, error)
	UpdateStatus(ctx context.Context, task domain.Task) (domain.Task, error)
}

type workflows interface {
	GetTransition(ctx context.Context, transitionID int) (domain.Transition, error)
}

type transactor interface {
	InTx(ctx context.Context, fn func(ctx context.Context) error) error
}
