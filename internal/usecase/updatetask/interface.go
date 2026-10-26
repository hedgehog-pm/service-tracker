package updatetask

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type repository interface {
	GetTaskForUpdate(ctx context.Context, projectID, taskID int) (domain.Task, error)
	UpdateTask(ctx context.Context, task domain.Task) (domain.Task, error)
}

type transactor interface {
	InTx(ctx context.Context, fn func(ctx context.Context) error) error
}
