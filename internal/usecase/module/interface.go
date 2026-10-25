package module

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type repository interface {
	GetModule(ctx context.Context, projectID, moduleID int) (domain.Module, error)
	InsertModule(ctx context.Context, module domain.Module) (domain.Module, error)
	UpdateModule(ctx context.Context, module domain.Module) (domain.Module, error)
	DeleteModule(ctx context.Context, projectID, moduleID int) error
}
