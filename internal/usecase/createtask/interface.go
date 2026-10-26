package createtask

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type repository interface {
	InsertTask(ctx context.Context, task domain.Task) (domain.Task, error)
}

type projects interface {
	GetProject(ctx context.Context, projectID int) (domain.Project, error)
	ReserveTaskNumber(ctx context.Context, projectID int) (int, string, error)
}

type workflows interface {
	GetInitialStatus(ctx context.Context, schemeID, issueTypeID int) (domain.InitialStatus, error)
}

type priorities interface {
	GetDefaultPriority(ctx context.Context, schemeID int) (int, error)
}

type transactor interface {
	InTx(ctx context.Context, fn func(ctx context.Context) error) error
}
