package listboardcolumntasks

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type (
	boardReader interface {
		GetBoard(ctx context.Context, projectID, boardID int) (domain.Board, error)
	}
	taskReader interface {
		ListBoardColumnTasks(
			ctx context.Context,
			projectID, boardID, columnID int,
			afterRank, afterNumber *int,
			limit int,
		) ([]domain.Task, error)
	}
)
