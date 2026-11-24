package projectsummary

import (
	"context"
	"time"
)

type projectReader interface {
	GetProjectSummary(ctx context.Context, projectID int, today time.Time) (Summary, error)
}
