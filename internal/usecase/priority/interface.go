package priority

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type (
	priorityReader interface {
		ListPriorities(ctx context.Context, schemeID int) ([]domain.Priority, error)
		GetDefaultPriority(ctx context.Context, schemeID int) (int, error)
	}
	projectReader interface {
		GetProject(ctx context.Context, projectID int) (domain.Project, error)
	}
)
