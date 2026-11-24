package listboardcolumntasks

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type Filters struct {
	Query      string
	StatusID   *int
	ModuleID   *int
	PriorityID *int
	TaskTypeID *int
}

type (
	boardReader interface {
		GetBoard(ctx context.Context, projectID, boardID int) (domain.Board, error)
	}
	taskReader interface {
		ListBoardColumnTasks(
			ctx context.Context,
			projectID, boardID, columnID int,
			filters Filters,
			afterRank, afterNumber *int,
			limit int,
		) ([]domain.Task, error)
	}
)
