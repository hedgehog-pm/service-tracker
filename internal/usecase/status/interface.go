package status

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type repository interface {
	GetStatus(ctx context.Context, projectID, statusID int) (domain.Status, error)
	InsertStatus(ctx context.Context, status domain.Status) (domain.Status, error)
	UpdateStatus(ctx context.Context, status domain.Status) (domain.Status, error)
	DeleteStatus(ctx context.Context, projectID, statusID int) error
}
