package project

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type repository interface {
	GetProject(ctx context.Context, projectID int) (domain.Project, error)
	ListProjects(ctx context.Context, limit int, offset int) ([]domain.Project, error)
	UpdateProject(ctx context.Context, project domain.Project) (domain.Project, error)
	DeleteProject(ctx context.Context, projectID int) error
}
