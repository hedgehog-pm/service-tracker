package status

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type (
	statusReader interface {
		ListStatuses(ctx context.Context, schemeID int) ([]domain.Status, error)
	}
	projectReader interface {
		GetProject(ctx context.Context, projectID int) (domain.Project, error)
	}
)
