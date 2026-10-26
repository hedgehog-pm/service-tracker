package status

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

// Usecase получения статусов проекта.
//
// Возвращает статусы всех workflow из схемы в [domain.Project.WorkflowSchemeID].
// Общие для нескольких workflow статусы не повторяются.
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

func (u *Usecase) List(ctx context.Context, projectID int) ([]domain.Status, error) {
	project, err := u.projects.GetProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return u.reader.ListStatuses(ctx, project.WorkflowSchemeID)
}
