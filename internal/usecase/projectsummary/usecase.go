package projectsummary

import (
	"context"
	"time"
)

type Summary struct {
	TasksTotal            int64
	TasksOpen             int64
	TasksDueSoon          int64
	TasksOverdue          int64
	TotalTimeSpentSeconds int64
	TotalEstimateSeconds  *int64
	UpcomingTasks         []UpcomingTask
}

type UpcomingTask struct {
	ID         int
	Key        string
	Title      string
	StatusID   int
	PriorityID int
	DueDate    time.Time
}

type Usecase struct {
	projects projectReader
}

func New(projects projectReader) *Usecase {
	return &Usecase{projects: projects}
}

func (u *Usecase) Do(ctx context.Context, projectID int) (Summary, error) {
	today := time.Now().UTC().Truncate(24 * time.Hour)

	return u.projects.GetProjectSummary(ctx, projectID, today)
}
