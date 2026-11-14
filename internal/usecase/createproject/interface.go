package createproject

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type (
	projectRepository interface {
		InsertProject(ctx context.Context, project domain.Project) (domain.Project, error)
	}
	workflowSchemeReader interface {
		GetDefaultWorkflowScheme(ctx context.Context) (int, error)
	}
	prioritySchemeReader interface {
		GetDefaultPriorityScheme(ctx context.Context) (int, error)
	}
)
