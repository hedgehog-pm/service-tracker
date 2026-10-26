package priority

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type (
	Result struct {
		Items             []domain.Priority
		DefaultPriorityID int
	}
	Usecase struct {
		reader   reader
		projects projects
	}
)

func New(reader reader, projects projects) *Usecase {
	return &Usecase{reader: reader, projects: projects}
}

func (u *Usecase) List(ctx context.Context, projectID int) (Result, error) {
	project, err := u.projects.GetProject(ctx, projectID)
	if err != nil {
		return Result{}, err
	}

	defaultID, err := u.reader.GetDefaultPriority(ctx, project.PrioritySchemeID)
	if err != nil {
		return Result{}, err
	}

	items, err := u.reader.ListPriorities(ctx, project.PrioritySchemeID)
	if err != nil {
		return Result{}, err
	}

	return Result{Items: items, DefaultPriorityID: defaultID}, nil
}
