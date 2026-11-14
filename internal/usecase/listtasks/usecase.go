package listtasks

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type Params struct {
	ProjectID  int
	StatusID   *int
	ModuleID   *int
	TaskTypeID *int
	PriorityID *int
	Limit      int
	Offset     int
}

type Usecase struct {
	tasks taskReader
}

func New(tasks taskReader) *Usecase {
	return &Usecase{tasks: tasks}
}

func (u *Usecase) Do(ctx context.Context, params Params) ([]domain.Task, error) {
	if _, err := u.tasks.GetProjectKey(ctx, params.ProjectID); err != nil {
		return nil, err
	}

	return u.tasks.ListTasks(
		ctx,
		params.ProjectID,
		params.StatusID,
		params.ModuleID,
		params.TaskTypeID,
		params.PriorityID,
		params.Limit,
		params.Offset,
	)
}
