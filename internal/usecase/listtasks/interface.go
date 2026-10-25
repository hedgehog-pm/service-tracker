package listtasks

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type reader interface {
	GetProjectKey(ctx context.Context, projectID int) (string, error)
	ListTasks(ctx context.Context, projectID int, statusID, moduleID *int, limit, offset int) ([]domain.Task, error)
}
