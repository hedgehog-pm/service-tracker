package task

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type Usecase struct {
	tasks taskRepository
}

func New(tasks taskRepository) *Usecase {
	return &Usecase{tasks: tasks}
}

func (u *Usecase) Get(ctx context.Context, projectID, taskID int) (domain.Task, error) {
	return u.tasks.GetTask(ctx, projectID, taskID)
}

func (u *Usecase) Delete(ctx context.Context, projectID, taskID int) error {
	return u.tasks.DeleteTask(ctx, projectID, taskID)
}
