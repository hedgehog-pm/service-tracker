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
	workflowSchemes workflowSchemeReader
	prioritySchemes prioritySchemeReader
}

func New(
	projects projectRepository,
	workflowSchemes workflowSchemeReader,
	prioritySchemes prioritySchemeReader,
) *Usecase {
	return &Usecase{
		projects:        projects,
		workflowSchemes: workflowSchemes,
		prioritySchemes: prioritySchemes,
	}
}

func (u *Usecase) Do(ctx context.Context, key, name string) (domain.Project, error) {
	project := domain.NewProject(key, name)
	if err := project.Validate(); err != nil {
		return domain.Project{}, err
	}

	var err error

	project.WorkflowSchemeID, err = u.workflowSchemes.GetDefaultWorkflowScheme(ctx)
	if err != nil {
		return domain.Project{}, err
	}

	project.PrioritySchemeID, err = u.prioritySchemes.GetDefaultPriorityScheme(ctx)
	if err != nil {
		return domain.Project{}, err
	}

	return u.projects.InsertProject(ctx, project)
}
