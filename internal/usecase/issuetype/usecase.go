package issuetype

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type (
	Item struct {
		domain.IssueType
		InitialStatus domain.Status
	}
	Usecase struct {
		reader   reader
		projects projects
	}
)

func New(reader reader, projects projects) *Usecase {
	return &Usecase{reader: reader, projects: projects}
}

func (u *Usecase) List(ctx context.Context, projectID int) ([]Item, error) {
	project, err := u.projects.GetProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return u.reader.ListIssueTypes(ctx, project.WorkflowSchemeID)
}
