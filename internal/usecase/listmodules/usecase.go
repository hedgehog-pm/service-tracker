package listmodules

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type Usecase struct {
	modules  moduleReader
	projects projectReader
}

func New(modules moduleReader, projects projectReader) *Usecase {
	return &Usecase{
		modules:  modules,
		projects: projects,
	}
}

func (u *Usecase) Do(
	ctx context.Context,
	projectID int,
	limit int,
	offset int,
) ([]domain.Module, error) {
	if _, err := u.projects.GetProject(ctx, projectID); err != nil {
		return nil, err
	}

	return u.modules.ListModules(ctx, projectID, limit, offset)
}
