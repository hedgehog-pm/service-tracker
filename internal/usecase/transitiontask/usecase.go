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
	// Usecase смены статуса задачи через переход workflow.
	//
	// Переход выбирается по [Params.TransitionID].
	//
	//   - [domain.Transition.WorkflowID] должен совпадать с [domain.Task.WorkflowID],
	//     а [domain.Transition.FromStatusID] с текущим [domain.Task.StatusID].
	//     Иначе возвращается [domain.ErrTransitionUnavailable].
	//   - [domain.Task.StatusID] меняется на [domain.Transition.ToStatus.ID],
	//     а [domain.Task.UpdatedAt] обновляется текущим временем в UTC.
	//
	// Проверка перехода и сохранение статуса проходят в одной транзакции под блокировкой задачи.
	Usecase struct {
		repository repository
		workflows  workflows
		transactor transactor
	}
)

func New(repository repository, workflows workflows, transactor transactor) *Usecase {
	return &Usecase{
		repository: repository,
		workflows:  workflows,
		transactor: transactor,
	}
}

func (u *Usecase) Do(ctx context.Context, params Params) (domain.Task, error) {
	var task domain.Task

	err := u.transactor.InTx(ctx, func(ctx context.Context) error {
		var err error

		// Блокируем задачу до проверки перехода, чтобы параллельный запрос
		// не применил переход по устаревшему статусу.
		task, err = u.repository.GetTaskForUpdate(ctx, params.ProjectID, params.TaskID)
		if err != nil {
			return err
		}

		transition, err := u.workflows.GetTransition(ctx, params.TransitionID)
		if err != nil {
			return err
		}

		if err := task.Transition(transition); err != nil {
			return err
		}

		task.UpdatedAt = time.Now().UTC()
		task, err = u.repository.UpdateStatus(ctx, task)

		return err
	})
	if err != nil {
		return domain.Task{}, err
	}

	return task, nil
}
