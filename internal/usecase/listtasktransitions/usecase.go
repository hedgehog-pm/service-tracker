package listtasktransitions

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

// Usecase получения доступных переходов задачи.
//
// Переходы выбираются по [domain.Task.WorkflowID] и текущему [domain.Task.StatusID].
// Список отражает состояние на момент чтения и не резервирует переход.
type Usecase struct {
	tasks  tasks
	reader reader
}

func New(tasks tasks, reader reader) *Usecase {
	return &Usecase{
		tasks:  tasks,
		reader: reader,
	}
}

func (u *Usecase) Do(ctx context.Context, projectID, taskID int) ([]domain.Transition, error) {
	task, err := u.tasks.GetTask(ctx, projectID, taskID)
	if err != nil {
		return nil, err
	}

	return u.reader.ListTransitions(ctx, task.WorkflowID, task.StatusID)
}
