package listmodules

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type (
	moduleReader interface {
		ListModules(
			ctx context.Context,
			projectID int,
			limit int,
			offset int,
		) ([]ModuleWithProgress, error)
	}
	projectReader interface {
		GetProject(ctx context.Context, projectID int) (domain.Project, error)
	}
)
