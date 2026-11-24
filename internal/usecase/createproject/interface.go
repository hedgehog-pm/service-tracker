package createproject

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type (
	projectRepository interface {
		InsertProject(ctx context.Context, project domain.Project) (domain.Project, error)
	}
	boardRepository interface {
		InsertDefaultBoard(ctx context.Context, projectID int) error
	}
	workflowSchemeReader interface {
		GetDefaultWorkflowScheme(ctx context.Context) (int, error)
	}
	prioritySchemeReader interface {
		GetDefaultPriorityScheme(ctx context.Context) (int, error)
	}
	transactor interface {
		InTx(ctx context.Context, fn func(ctx context.Context) error) error
	}
)
