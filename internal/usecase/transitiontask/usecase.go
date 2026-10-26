package transitiontask

import (
	"context"
	"time"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type (
	Params struct {
		ProjectID    int
		TaskID       int
		TransitionID int
	}
	Usecase struct {
		repository repository
		workflows  workflows
		transactor transactor
	}
)

func New(repository repository, workflows workflows, transactor transactor) *Usecase {
	return &Usecase{repository: repository, workflows: workflows, transactor: transactor}
}

func (u *Usecase) Do(ctx context.Context, params Params) (domain.Task, error) {
	var task domain.Task

	err := u.transactor.InTx(ctx, func(ctx context.Context) error {
		var err error

		task, err = u.repository.GetTaskForUpdate(ctx, params.ProjectID, params.TaskID)
		if err != nil {
			return err
		}

		transition, err := u.workflows.GetTransition(ctx, params.TransitionID)
		if err != nil {
			return err
		}

		if err := task.Transition(transition, time.Now().UTC()); err != nil {
			return err
		}

		task, err = u.repository.UpdateStatus(ctx, task)

		return err
	})
	if err != nil {
		return domain.Task{}, err
	}

	return task, nil
}
