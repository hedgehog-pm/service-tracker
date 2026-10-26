package tasktype

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type (
	projects interface {
		GetProject(ctx context.Context, projectID int) (domain.Project, error)
	}
	reader interface {
		ListTaskTypes(
			ctx context.Context,
			schemeID int,
		) ([]domain.TaskTypeWithInitialStatus, error)
	}
)
