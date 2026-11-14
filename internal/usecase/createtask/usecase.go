package createtask

import (
	"context"
	"errors"
	"time"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type Params struct {
	ProjectID       int
	TaskTypeID      int
	PriorityID      *int
	Title           string
	Description     string
	ModuleID        *int
	EstimateSeconds *int
	DueDate         *time.Time
}

// Usecase создания задачи.
//
//   - [domain.Task.WorkflowID] выбирается по [Params.TaskTypeID] в workflow-схеме проекта.
//   - [domain.Task.StatusID] берётся из начального статуса выбранного workflow.
//   - [domain.Task.PriorityID] берётся из [Params.PriorityID]. Если [Params.PriorityID] == nil,
//     используется приоритет по умолчанию из схемы приоритетов проекта.
type Usecase struct {
	tasks      taskRepository
	projects   projectRepository
	workflows  workflowReader
	priorities priorityReader
	transactor transactor
}

func New(
	tasks taskRepository,
	projects projectRepository,
	workflows workflowReader,
	priorities priorityReader,
	transactor transactor,
) *Usecase {
	return &Usecase{
		tasks:      tasks,
		projects:   projects,
		workflows:  workflows,
		priorities: priorities,
		transactor: transactor,
	}
}

func (u *Usecase) Do(ctx context.Context, params Params) (domain.Task, error) {
	fields := domain.TaskFields{
		Title:           params.Title,
		Description:     params.Description,
		ModuleID:        params.ModuleID,
		EstimateSeconds: params.EstimateSeconds,
		DueDate:         params.DueDate,
	}
	if err := fields.Validate(); err != nil {
		return domain.Task{}, err
	}

	var task domain.Task

	// Резервирование номера и вставка задачи должны откатываться вместе,
	// чтобы ошибка создания не расходовала номер в проекте.
	err := u.transactor.InTx(ctx, func(ctx context.Context) error {
		project, err := u.projects.GetProject(ctx, params.ProjectID)
		if err != nil {
			return err
		}

		initial, err := u.workflows.GetInitialStatus(
			ctx,
			project.WorkflowSchemeID,
			params.TaskTypeID,
		)
		if err != nil {
			return err
		}

		fields.PriorityID, err = u.priorityID(ctx, params.PriorityID, project.PrioritySchemeID)
		if err != nil {
			return err
		}

		number, key, err := u.projects.ReserveTaskNumber(ctx, params.ProjectID)
		if err != nil {
			return err
		}

		rank, err := u.nextRank(ctx, params.ProjectID)
		if err != nil {
			return err
		}

		now := time.Now().UTC()
		task, err = u.tasks.InsertTask(ctx, domain.Task{
			ProjectID:  params.ProjectID,
			ProjectKey: key,
			Number:     number,
			Rank:       rank,
			TaskTypeID: params.TaskTypeID,
			WorkflowID: initial.WorkflowID,
			StatusID:   initial.Status.ID,
			TaskFields: fields,
			CreatedAt:  now,
			UpdatedAt:  now,
		})

		return err
	})
	if err != nil {
		return domain.Task{}, err
	}

	return task, nil
}

func (u *Usecase) priorityID(ctx context.Context, requested *int, schemeID int) (int, error) {
	if requested != nil {
		return *requested, nil
	}

	return u.priorities.GetDefaultPriority(ctx, schemeID)
}

func (u *Usecase) nextRank(ctx context.Context, projectID int) (int, error) {
	last, err := u.tasks.GetLastRank(ctx, projectID)
	if err != nil {
		return 0, err
	}

	rank, err := domain.RankBetween(last, nil)
	if !errors.Is(err, domain.ErrRankExhausted) {
		return rank, err
	}

	tasks, err := u.tasks.ListTaskOrder(ctx, projectID)
	if err != nil {
		return 0, err
	}

	if err := domain.RebalanceTaskRanks(tasks); err != nil {
		return 0, err
	}

	if err := u.tasks.UpdateTaskRanks(ctx, tasks); err != nil {
		return 0, err
	}

	if len(tasks) == 0 {
		return domain.RankBetween(nil, nil)
	}

	return domain.RankBetween(&tasks[len(tasks)-1].Rank, nil)
}
