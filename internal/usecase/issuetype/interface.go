package issuetype

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type (
	projects interface {
		GetProject(ctx context.Context, projectID int) (domain.Project, error)
	}
	reader interface {
		ListIssueTypes(ctx context.Context, schemeID int) ([]Item, error)
	}
)
