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
	// Usecase получения приоритетов проекта.
	//
	// [Result.Items] и [Result.DefaultPriorityID] берутся из одной схемы,
	// указанной в [domain.Project.PrioritySchemeID].
	Usecase struct {
		priorities priorityReader
		projects   projectReader
	}
)

func New(priorities priorityReader, projects projectReader) *Usecase {
	return &Usecase{
		priorities: priorities,
		projects:   projects,
	}
}

func (u *Usecase) List(ctx context.Context, projectID int) (Result, error) {
	project, err := u.projects.GetProject(ctx, projectID)
	if err != nil {
		return Result{}, err
	}

	defaultID, err := u.priorities.GetDefaultPriority(ctx, project.PrioritySchemeID)
	if err != nil {
		return Result{}, err
	}

	items, err := u.priorities.ListPriorities(ctx, project.PrioritySchemeID)
	if err != nil {
		return Result{}, err
	}

	return Result{
		Items:             items,
		DefaultPriorityID: defaultID,
	}, nil
}
