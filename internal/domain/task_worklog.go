package domain

import "time"

type TaskWorklogFields struct {
	StartedAt        time.Time
	TimeSpentSeconds int
	Comment          *string
}

type TaskWorklog struct {
	ID     int
	TaskID int
	TaskWorklogFields
	CreatedAt time.Time
	UpdatedAt time.Time
}
