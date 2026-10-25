package task

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type repository interface {
	GetTask(ctx context.Context, projectID, taskID int) (domain.Task, error)
	DeleteTask(ctx context.Context, projectID, taskID int) error
}
