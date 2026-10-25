package task

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type Usecase struct {
	repository repository
}

func New(repository repository) *Usecase {
	return &Usecase{repository: repository}
}

func (u *Usecase) Get(ctx context.Context, projectID, taskID int) (domain.Task, error) {
	return u.repository.GetTask(ctx, projectID, taskID)
}

func (u *Usecase) Delete(ctx context.Context, projectID, taskID int) error {
	return u.repository.DeleteTask(ctx, projectID, taskID)
}
