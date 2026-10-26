package tasktype

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

// Usecase получения доступных типов задач проекта.
//
// Типы задач выбираются из схемы в [domain.Project.WorkflowSchemeID].
// [domain.TaskTypeWithInitialStatus.InitialStatus] берётся из workflow,
// выбранного для типа задачи в этой схеме.
type Usecase struct {
	reader   reader
	projects projects
}

func New(reader reader, projects projects) *Usecase {
	return &Usecase{
		reader:   reader,
		projects: projects,
	}
}

func (u *Usecase) List(
	ctx context.Context,
	projectID int,
) ([]domain.TaskTypeWithInitialStatus, error) {
	project, err := u.projects.GetProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return u.reader.ListTaskTypes(ctx, project.WorkflowSchemeID)
}
