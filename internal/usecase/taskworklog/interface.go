package taskworklog

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type taskWorklogRepository interface {
	GetTaskWorklog(ctx context.Context, projectID, taskID, worklogID int) (domain.TaskWorklog, error)
	ListTaskWorklogs(ctx context.Context, projectID, taskID, limit, offset int) ([]domain.TaskWorklog, int64, error)
	InsertTaskWorklog(ctx context.Context, projectID int, worklog domain.TaskWorklog) (domain.TaskWorklog, error)
	UpdateTaskWorklog(ctx context.Context, projectID int, worklog domain.TaskWorklog) (domain.TaskWorklog, error)
	DeleteTaskWorklog(ctx context.Context, projectID, taskID, worklogID int) error
}
