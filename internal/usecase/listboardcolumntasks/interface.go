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

type Query struct {
	ProjectID   int
	BoardID     int
	ColumnID    int
	Filters     Filters
	AfterRank   *int
	AfterNumber *int
	Limit       int
}

type (
	boardReader interface {
		GetBoard(ctx context.Context, projectID, boardID int) (domain.Board, error)
	}
	taskReader interface {
		ListBoardColumnTasks(ctx context.Context, query Query) ([]domain.Task, error)
	}
)
