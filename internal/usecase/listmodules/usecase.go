package listmodules

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type Usecase struct {
	repository repository
	projects   projects
}

func New(repository repository, projects projects) *Usecase {
	return &Usecase{repository: repository, projects: projects}
}

func (u *Usecase) Do(ctx context.Context, projectID int, limit int, offset int) ([]domain.Module, error) {
	if _, err := u.projects.GetProject(ctx, projectID); err != nil {
		return nil, err
	}

	return u.repository.ListModules(ctx, projectID, limit, offset)
}
