package updatetask

import (
	"context"
	"time"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type Usecase struct {
	repository repository
	transactor transactor
}

func New(repository repository, transactor transactor) *Usecase {
	return &Usecase{repository: repository, transactor: transactor}
}

func (u *Usecase) Do(ctx context.Context, projectID, taskID int, fields domain.TaskFields) (domain.Task, error) {
	if err := fields.Validate(); err != nil {
		return domain.Task{}, err
	}

	var task domain.Task

	err := u.transactor.InTx(ctx, func(ctx context.Context) error {
		var err error

		task, err = u.repository.GetTask(ctx, projectID, taskID)
		if err != nil {
			return err
		}

		task.TaskFields = fields
		task.UpdatedAt = time.Now().UTC()
		task, err = u.repository.UpdateTask(ctx, task)

		return err
	})
	if err != nil {
		return domain.Task{}, err
	}

	return task, nil
}
