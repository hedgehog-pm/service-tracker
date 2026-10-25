package createtask

import (
	"context"
	"time"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type Usecase struct {
	repository repository
	projects   projects
	transactor transactor
}

func New(repository repository, projects projects, transactor transactor) *Usecase {
	return &Usecase{repository: repository, projects: projects, transactor: transactor}
}

func (u *Usecase) Do(ctx context.Context, projectID int, fields domain.TaskFields) (domain.Task, error) {
	if err := fields.Validate(); err != nil {
		return domain.Task{}, err
	}

	var task domain.Task

	err := u.transactor.InTx(ctx, func(ctx context.Context) error {
		number, key, err := u.projects.ReserveTaskNumber(ctx, projectID)
		if err != nil {
			return err
		}

		now := time.Now().UTC()
		task, err = u.repository.InsertTask(ctx, domain.Task{
			ProjectID: projectID, ProjectKey: key, Number: number,
			TaskFields: fields, CreatedAt: now, UpdatedAt: now,
		})

		return err
	})
	if err != nil {
		return domain.Task{}, err
	}

	return task, nil
}
