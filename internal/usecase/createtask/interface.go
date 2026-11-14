package createtask

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type (
	taskRepository interface {
		GetLastRank(ctx context.Context, projectID int) (*int, error)
		ListTaskOrder(ctx context.Context, projectID int) ([]domain.Task, error)
		UpdateTaskRanks(ctx context.Context, tasks []domain.Task) error
		InsertTask(ctx context.Context, task domain.Task) (domain.Task, error)
	}
	projectRepository interface {
		GetProject(ctx context.Context, projectID int) (domain.Project, error)
		ReserveTaskNumber(ctx context.Context, projectID int) (int, string, error)
	}
	workflowReader interface {
		GetInitialStatus(ctx context.Context, schemeID, taskTypeID int) (domain.InitialStatus, error)
	}
	priorityReader interface {
		GetDefaultPriority(ctx context.Context, schemeID int) (int, error)
	}
	transactor interface {
		InTx(ctx context.Context, fn func(ctx context.Context) error) error
	}
)
