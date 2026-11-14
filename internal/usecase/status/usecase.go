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
	statuses statusReader
	projects projectReader
}

func New(statuses statusReader, projects projectReader) *Usecase {
	return &Usecase{
		statuses: statuses,
		projects: projects,
	}
}

func (u *Usecase) List(ctx context.Context, projectID int) ([]domain.Status, error) {
	project, err := u.projects.GetProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return u.statuses.ListStatuses(ctx, project.WorkflowSchemeID)
}
