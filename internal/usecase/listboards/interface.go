package listboards

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type (
	boardReader interface {
		ListBoards(ctx context.Context, projectID, limit, offset int) ([]domain.Board, error)
	}
	projectReader interface {
		GetProject(ctx context.Context, projectID int) (domain.Project, error)
	}
)
