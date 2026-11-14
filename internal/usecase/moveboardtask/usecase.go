package moveboardtask

import (
	"context"
	"errors"
	"time"

	"github.com/samber/lo"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type (
	Params struct {
		ProjectID    int
		BoardID      int
		TaskID       int
		ColumnID     int
		BeforeTaskID *int
		TransitionID *int
	}
	Usecase struct {
		tasks       taskRepository
		order       taskOrderReader
		boards      boardRepository
		projects    projectLocker
		transitions transitionReader
		transactor  transactor
	}
)

func New(
	tasks taskRepository,
	order taskOrderReader,
	boards boardRepository,
	projects projectLocker,
	transitions transitionReader,
	transactor transactor,
) *Usecase {
	return &Usecase{
		tasks:       tasks,
		order:       order,
		boards:      boards,
		projects:    projects,
		transitions: transitions,
		transactor:  transactor,
	}
}

func (u *Usecase) Do(ctx context.Context, params Params) (domain.Task, error) {
	var task domain.Task

	err := u.transactor.InTx(ctx, func(ctx context.Context) error {
		var err error

		task, err = u.move(ctx, params)

		return err
	})
	if err != nil {
		return domain.Task{}, err
	}

	return task, nil
}

func (u *Usecase) move(ctx context.Context, params Params) (domain.Task, error) {
	if lo.FromPtr(params.BeforeTaskID) == params.TaskID {
		return domain.Task{}, domain.NewError(
			domain.CodeInvalidInput,
			"task cannot be its own position anchor",
		)
	}

	if err := u.projects.LockProject(ctx, params.ProjectID); err != nil {
		return domain.Task{}, err
	}

	board, err := u.boards.GetBoardForUpdate(ctx, params.ProjectID, params.BoardID)
	if err != nil {
		return domain.Task{}, err
	}

	column, err := board.Column(params.ColumnID)
	if err != nil {
		return domain.Task{}, err
	}

	task, err := u.tasks.GetTaskForUpdate(ctx, params.ProjectID, params.TaskID)
	if err != nil {
		return domain.Task{}, err
	}

	if params.TransitionID != nil {
		transition, err := u.transitions.GetTransition(ctx, *params.TransitionID)
		if err != nil {
			return domain.Task{}, err
		}

		if err := task.Transition(transition); err != nil {
			return domain.Task{}, err
		}
	}

	if !column.ContainsStatus(task.StatusID) {
		return domain.Task{}, domain.ErrTransitionUnavailable
	}

	rank, err := u.position(ctx, params, column)
	if errors.Is(err, domain.ErrRankExhausted) {
		if err := u.rebalance(ctx, params.ProjectID); err != nil {
			return domain.Task{}, err
		}

		rank, err = u.position(ctx, params, column)
	}

	if err != nil {
		return domain.Task{}, err
	}

	task.Rank = rank
	task.UpdatedAt = time.Now().UTC()

	return u.tasks.UpdatePosition(ctx, task)
}

func (u *Usecase) position(
	ctx context.Context,
	params Params,
	column domain.BoardColumn,
) (int, error) {
	if params.BeforeTaskID != nil {
		anchor, err := u.tasks.GetTaskForUpdate(ctx, params.ProjectID, *params.BeforeTaskID)
		if errors.Is(err, domain.ErrTaskNotFound) {
			return 0, domain.NewError(
				domain.CodeConflict,
				"position anchor is not available in this project",
			)
		}

		if err != nil {
			return 0, err
		}

		if !column.ContainsStatus(anchor.StatusID) {
			return 0, domain.NewError(
				domain.CodeConflict,
				"position anchor is no longer in the target column",
			)
		}

		previous, err := u.order.GetPreviousRank(ctx, params.ProjectID, params.TaskID, anchor.Rank)
		if err != nil {
			return 0, err
		}

		return domain.RankBetween(previous, &anchor.Rank)
	}

	previous, err := u.order.GetLastColumnRank(
		ctx,
		params.ProjectID,
		params.BoardID,
		params.ColumnID,
		params.TaskID,
	)
	if err != nil {
		return 0, err
	}

	if previous == nil {
		last, err := u.order.GetLastRank(ctx, params.ProjectID)
		if err != nil {
			return 0, err
		}

		return domain.RankBetween(last, nil)
	}

	next, err := u.order.GetNextRank(ctx, params.ProjectID, params.TaskID, *previous)
	if err != nil {
		return 0, err
	}

	return domain.RankBetween(previous, next)
}

func (u *Usecase) rebalance(ctx context.Context, projectID int) error {
	tasks, err := u.order.ListTaskOrder(ctx, projectID)
	if err != nil {
		return err
	}

	if err := domain.RebalanceTaskRanks(tasks); err != nil {
		return err
	}

	return u.tasks.UpdateTaskRanks(ctx, tasks)
}
