package createproject

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

// Usecase создания проекта.
//
//   - [domain.Project.WorkflowSchemeID] берётся из схемы workflow по умолчанию.
//   - [domain.Project.PrioritySchemeID] берётся из схемы приоритетов по умолчанию.
//
// В проекте сохраняются ID схем, поэтому смена схем по умолчанию не меняет созданный проект.
type Usecase struct {
	projects        projectRepository
	boards          boardRepository
	workflowSchemes workflowSchemeReader
	prioritySchemes prioritySchemeReader
	transactor      transactor
}

func New(
	projects projectRepository,
	boards boardRepository,
	workflowSchemes workflowSchemeReader,
	prioritySchemes prioritySchemeReader,
	transactor transactor,
) *Usecase {
	return &Usecase{
		projects:        projects,
		boards:          boards,
		workflowSchemes: workflowSchemes,
		prioritySchemes: prioritySchemes,
		transactor:      transactor,
	}
}

func (u *Usecase) Do(ctx context.Context, key, name string) (domain.Project, error) {
	project := domain.NewProject(key, name)
	if err := project.Validate(); err != nil {
		return domain.Project{}, err
	}

	err := u.transactor.InTx(ctx, func(ctx context.Context) error {
		var err error

		project, err = u.insertProject(ctx, project)
		if err != nil {
			return err
		}

		return u.boards.InsertDefaultBoard(ctx, project.ID)
	})
	if err != nil {
		return domain.Project{}, err
	}

	return project, nil
}

func (u *Usecase) insertProject(
	ctx context.Context,
	project domain.Project,
) (domain.Project, error) {
	workflowSchemeID, err := u.workflowSchemes.GetDefaultWorkflowScheme(ctx)
	if err != nil {
		return domain.Project{}, err
	}

	prioritySchemeID, err := u.prioritySchemes.GetDefaultPriorityScheme(ctx)
	if err != nil {
		return domain.Project{}, err
	}

	project.WorkflowSchemeID = workflowSchemeID
	project.PrioritySchemeID = prioritySchemeID

	return u.projects.InsertProject(ctx, project)
}
