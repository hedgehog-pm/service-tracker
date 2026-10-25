package createtask

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type repository interface {
	InsertTask(ctx context.Context, task domain.Task) (domain.Task, error)
}

type projects interface {
	ReserveTaskNumber(ctx context.Context, projectID int) (int, string, error)
}

type transactor interface {
	InTx(ctx context.Context, fn func(ctx context.Context) error) error
}
