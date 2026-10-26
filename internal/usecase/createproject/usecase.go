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
	repository repository
	workflows  workflows
	priorities priorities
}

func New(repository repository, workflows workflows, priorities priorities) *Usecase {
	return &Usecase{
		repository: repository,
		workflows:  workflows,
		priorities: priorities,
	}
}

func (u *Usecase) Do(ctx context.Context, key, name string) (domain.Project, error) {
	project := domain.NewProject(key, name)
	if err := project.Validate(); err != nil {
		return domain.Project{}, err
	}

	var err error

	project.WorkflowSchemeID, err = u.workflows.GetDefaultWorkflowScheme(ctx)
	if err != nil {
		return domain.Project{}, err
	}

	project.PrioritySchemeID, err = u.priorities.GetDefaultPriorityScheme(ctx)
	if err != nil {
		return domain.Project{}, err
	}

	return u.repository.InsertProject(ctx, project)
}
