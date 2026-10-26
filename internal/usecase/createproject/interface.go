package createproject

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type (
	repository interface {
		InsertProject(ctx context.Context, project domain.Project) (domain.Project, error)
	}
	workflows interface {
		GetDefaultWorkflowScheme(ctx context.Context) (int, error)
	}
	priorities interface {
		GetDefaultPriorityScheme(ctx context.Context) (int, error)
	}
)
