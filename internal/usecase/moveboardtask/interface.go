package moveboardtask

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type (
	taskRepository interface {
		GetTaskForUpdate(ctx context.Context, projectID, taskID int) (domain.Task, error)
		UpdateTaskRanks(ctx context.Context, tasks []domain.Task) error
		UpdatePosition(ctx context.Context, task domain.Task) (domain.Task, error)
	}
	taskOrderReader interface {
		GetLastRank(ctx context.Context, projectID int) (*int, error)
		GetPreviousRank(ctx context.Context, projectID, taskID, rank int) (*int, error)
		GetNextRank(ctx context.Context, projectID, taskID, rank int) (*int, error)
		GetLastColumnRank(ctx context.Context, projectID, boardID, columnID, taskID int) (*int, error)
		ListTaskOrder(ctx context.Context, projectID int) ([]domain.Task, error)
	}
	boardRepository interface {
		GetBoardForUpdate(ctx context.Context, projectID, boardID int) (domain.Board, error)
	}
	projectLocker interface {
		LockProject(ctx context.Context, projectID int) error
	}
	transitionReader interface {
		GetTransition(ctx context.Context, transitionID int) (domain.Transition, error)
	}
	transactor interface {
		InTx(ctx context.Context, fn func(ctx context.Context) error) error
	}
)
