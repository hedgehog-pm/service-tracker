package createtask

import (
	"context"
	"time"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type Params struct {
	ProjectID   int
	IssueTypeID int
	PriorityID  *int
	Title       string
	Description string
	ModuleID    *int
}

// Usecase создания задачи.
//
//   - [domain.Task.WorkflowID] выбирается по [Params.IssueTypeID] в workflow-схеме проекта.
//   - [domain.Task.StatusID] берётся из начального статуса выбранного workflow.
//   - [domain.Task.PriorityID] берётся из [Params.PriorityID]. Если [Params.PriorityID] == nil,
//     используется приоритет по умолчанию из схемы приоритетов проекта.
type Usecase struct {
	repository repository
	projects   projects
	workflows  workflows
	priorities priorities
	transactor transactor
}

func New(
	repository repository,
	projects projects,
	workflows workflows,
	priorities priorities,
	transactor transactor,
) *Usecase {
	return &Usecase{
		repository: repository,
		projects:   projects,
		workflows:  workflows,
		priorities: priorities,
		transactor: transactor,
	}
}

func (u *Usecase) Do(ctx context.Context, params Params) (domain.Task, error) {
	fields := domain.TaskFields{
		Title:       params.Title,
		Description: params.Description,
		ModuleID:    params.ModuleID,
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
			params.IssueTypeID,
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

		now := time.Now().UTC()
		task, err = u.repository.InsertTask(ctx, domain.Task{
			ProjectID:   params.ProjectID,
			ProjectKey:  key,
			Number:      number,
			IssueTypeID: params.IssueTypeID,
			WorkflowID:  initial.WorkflowID,
			StatusID:    initial.Status.ID,
			TaskFields:  fields,
			CreatedAt:   now,
			UpdatedAt:   now,
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
