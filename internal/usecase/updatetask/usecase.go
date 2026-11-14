package updatetask

import (
	"context"
	"time"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type Usecase struct {
	tasks      taskRepository
	transactor transactor
}

func New(tasks taskRepository, transactor transactor) *Usecase {
	return &Usecase{
		tasks:      tasks,
		transactor: transactor,
	}
}

func (u *Usecase) Do(
	ctx context.Context,
	projectID, taskID int,
	fields domain.TaskFields,
) (domain.Task, error) {
	if err := fields.Validate(); err != nil {
		return domain.Task{}, err
	}

	var task domain.Task

	err := u.transactor.InTx(ctx, func(ctx context.Context) error {
		var err error

		task, err = u.tasks.GetTaskForUpdate(ctx, projectID, taskID)
		if err != nil {
			return err
		}

		task.TaskFields = fields
		task.UpdatedAt = time.Now().UTC()
		task, err = u.tasks.UpdateTask(ctx, task)

		return err
	})
	if err != nil {
		return domain.Task{}, err
	}

	return task, nil
}
