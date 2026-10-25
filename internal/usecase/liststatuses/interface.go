package liststatuses

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type repository interface {
	ListStatuses(ctx context.Context, projectID int, limit int, offset int) ([]domain.Status, error)
}

type projects interface {
	GetProject(ctx context.Context, projectID int) (domain.Project, error)
}
